"""Urwid-based s4net console for local section4 operations."""

from __future__ import annotations

import asyncio
from collections.abc import Callable
from dataclasses import dataclass
from pathlib import Path

import urwid
from sqlalchemy.orm import Session, sessionmaker

from section4.storage import create_all, create_session_factory
from section4.storage.tables import Capability, Incident, Job
from section4.transport import (
    PyTAKRuntime,
    build_capability_cot,
    build_incident_cot,
    build_job_cot,
    record_published_cot,
)
from section4.tui.data import (
    BrowserItem,
    DashboardData,
    load_page_items,
)


@dataclass(frozen=True)
class PageDefinition:
    """A navigable page in the s4net operator console."""

    key: str
    title: str


class FocusListBox(urwid.ListBox):
    """ListBox with vim-style navigation and focus callbacks."""

    def __init__(
        self,
        body: urwid.ListWalker,
        on_focus: Callable[[object, int | None], None],
    ) -> None:
        super().__init__(body)
        self._on_focus = on_focus

    def keypress(self, size: tuple[int, int], key: str) -> str | None:
        """Map vim keys onto ListBox navigation."""

        keymap = {
            "j": "down",
            "k": "up",
            "ctrl d": "page down",
            "ctrl u": "page up",
            "g": "home",
            "G": "end",
        }
        mapped = keymap.get(key, key)
        before = self.get_focus()
        result = super().keypress(size, mapped)
        after = self.get_focus()
        if before != after:
            self._on_focus(after[0], after[1])
        return result


PALETTE = [
    ("header", "black", "light gray", "bold"),
    ("footer", "black", "dark cyan"),
    ("nav", "light gray", "black"),
    ("nav_focus", "black", "dark cyan", "bold"),
    ("panel_title", "light cyan", "black", "bold"),
    ("section", "yellow", "black", "bold"),
    ("body", "light gray", "black"),
    ("muted", "dark gray", "black"),
    ("success", "dark green", "black", "bold"),
    ("error", "light red", "black", "bold"),
]

PAGES = [
    PageDefinition("dashboard", "COP"),
    PageDefinition("incidents", "Requests"),
    PageDefinition("capabilities", "Capabilities"),
    PageDefinition("jobs", "Tasks"),
    PageDefinition("artifacts", "Artifacts"),
    PageDefinition("events", "Sync Log"),
]

PUBLISH_LABELS = {
    "incidents": "incident",
    "capabilities": "capability",
    "jobs": "job",
}


class S4NetTUIApplication:
    """Local-first s4net operator console over section4 persistence."""

    def __init__(
        self,
        session_factory: sessionmaker[Session],
        cot_url: str,
    ) -> None:
        self.session_factory = session_factory
        self.cot_url = cot_url
        self.current_page_key = "dashboard"
        self.current_items: list[BrowserItem] = []
        self.current_item_index = 0
        self.footer_text = urwid.Text("", align="center")
        self.nav_walker = urwid.SimpleFocusListWalker([])
        self.list_walker = urwid.SimpleFocusListWalker([])
        self.detail_walker = urwid.SimpleFocusListWalker([])
        self.page_title = urwid.Text("", align="left")
        self.detail_title = urwid.Text("", align="left")
        self.header = urwid.AttrMap(
            urwid.Text(" s4net ", align="center"),
            "header",
        )
        self.nav_listbox = FocusListBox(
            self.nav_walker,
            self._on_nav_focus,
        )
        self.record_listbox = FocusListBox(
            self.list_walker,
            self._on_record_focus,
        )
        self.detail_listbox = FocusListBox(
            self.detail_walker,
            self._on_detail_focus,
        )
        self.frame = urwid.Frame(
            header=self.header,
            body=self._build_body(),
            footer=urwid.AttrMap(self.footer_text, "footer"),
        )
        self.loop = urwid.MainLoop(
            self.frame,
            PALETTE,
            unhandled_input=self._unhandled_input,
        )
        self._populate_nav()
        self._set_footer()
        self._refresh_page()

    def _build_body(self) -> urwid.Widget:
        """Build the static TUI layout."""

        nav_box = urwid.LineBox(self.nav_listbox, title="Views")
        list_box = urwid.Frame(
            body=urwid.LineBox(self.record_listbox),
            header=urwid.AttrMap(self.page_title, "panel_title"),
        )
        detail_box = urwid.Frame(
            body=urwid.LineBox(self.detail_listbox),
            header=urwid.AttrMap(self.detail_title, "panel_title"),
        )
        content = urwid.Columns(
            [
                ("weight", 2, list_box),
                ("weight", 3, detail_box),
            ],
            dividechars=1,
            focus_column=0,
        )
        return urwid.Columns(
            [
                ("weight", 1, nav_box),
                ("weight", 5, content),
            ],
            dividechars=1,
            focus_column=0,
        )

    def _populate_nav(self) -> None:
        """Build the navigation column."""

        self.nav_walker.clear()
        for page in PAGES:
            self.nav_walker.append(
                urwid.AttrMap(
                    urwid.SelectableIcon(page.title, cursor_position=0),
                    "nav",
                    "nav_focus",
                )
            )

    def _set_footer(
        self,
        message: str | None = None,
        *,
        style: str = "footer",
    ) -> None:
        """Render keyboard help plus an optional status message."""

        help_text = " h/l panes  j/k move  g/G top/bottom  enter open  "
        help_text += "p publish ATAK  r refresh  q quit "
        if message:
            self.footer_text.set_text(
                [(style, message), ("footer", " | "), help_text]
            )
            return
        self.footer_text.set_text(help_text)

    def _on_nav_focus(self, _: object, position: int | None) -> None:
        """Track navigation focus changes."""

        if position is None:
            return
        page = PAGES[position]
        self.current_page_key = page.key
        self._refresh_page()

    def _on_record_focus(self, _: object, position: int | None) -> None:
        """Render detail when the focused list item changes."""

        if position is None:
            return
        if 0 <= position < len(self.current_items):
            self.current_item_index = position
            self._show_item_detail(self.current_items[position])

    def _on_detail_focus(self, _: object, __: int | None) -> None:
        """Ignore detail focus changes."""

    def _refresh_page(self) -> None:
        """Reload the active page from SQLite."""

        page_data = load_page_items(
            self.session_factory,
            self.current_page_key,
        )
        if isinstance(page_data, DashboardData):
            self.current_items = []
            self.current_item_index = 0
            self._render_dashboard(page_data)
            self._set_footer()
            return

        title, items = page_data
        self.current_items = items
        self.current_item_index = 0
        self.page_title.set_text(title)
        self.detail_title.set_text("Detail")
        self.list_walker.clear()
        self.detail_walker.clear()

        if not items:
            self.list_walker.append(
                urwid.Text(("muted", "No records available."))
            )
            self.detail_walker.append(
                urwid.Text(("muted", "Select a record when available."))
            )
            self._set_footer()
            return

        for item in items:
            self.list_walker.append(
                urwid.AttrMap(
                    urwid.SelectableIcon(item.label, cursor_position=0),
                    "body",
                    "nav_focus",
                )
            )

        self.record_listbox.focus_position = 0
        self._show_item_detail(items[0])
        self._set_footer()

    def _render_dashboard(self, dashboard: DashboardData) -> None:
        """Render the dashboard summary page."""

        self.page_title.set_text("Logistics COP")
        self.detail_title.set_text("Sync and Event Log")
        self.list_walker.clear()
        self.detail_walker.clear()

        for line in dashboard.summary_lines:
            self.list_walker.append(urwid.Text((line.style, line.text)))

        for line in dashboard.recent_event_lines:
            self.detail_walker.append(urwid.Text((line.style, line.text)))

    def _show_item_detail(self, item: BrowserItem) -> None:
        """Show detail lines for a record."""

        self.detail_title.set_text(item.label)
        self.detail_walker.clear()
        for line in item.detail_lines:
            self.detail_walker.append(urwid.Text((line.style, line.text)))
        if len(self.detail_walker) > 0:
            self.detail_listbox.focus_position = 0

    def _move_focus_horizontal(self, direction: int) -> None:
        """Move focus between nav, list, and detail panes."""

        widget = self.frame.body
        if not isinstance(widget, urwid.Columns):
            return
        current = widget.focus_position
        next_pos = max(0, min(len(widget.contents) - 1, current + direction))
        widget.focus_position = next_pos
        if next_pos == 1:
            content = widget.contents[1][0]
            if isinstance(content, urwid.Columns):
                if direction > 0:
                    content.focus_position = 0
        if next_pos == 0:
            self.nav_listbox.focus_position = min(
                len(PAGES) - 1,
                next(
                    (
                        index
                        for index, page in enumerate(PAGES)
                        if page.key == self.current_page_key
                    ),
                    0,
                ),
            )

    def _activate_focused(self) -> None:
        """Activate the current focus target."""

        body = self.frame.body
        if not isinstance(body, urwid.Columns):
            return
        if body.focus_position == 0:
            _, position = self.nav_listbox.get_focus()
            self._on_nav_focus(None, position)
            return
        content = body.contents[1][0]
        if not isinstance(content, urwid.Columns):
            return
        if content.focus_position == 0:
            _, position = self.record_listbox.get_focus()
            self._on_record_focus(None, position)

    def _selected_item(self) -> BrowserItem | None:
        """Return the currently selected list item, if any."""

        if not self.current_items:
            return None
        if self.current_item_index >= len(self.current_items):
            return None
        return self.current_items[self.current_item_index]

    async def _publish_selected_async(self) -> str:
        """Publish the selected record over CoT."""

        selected = self._selected_item()
        if selected is None:
            raise RuntimeError("No selectable record on the current page")

        runtime = PyTAKRuntime(self.cot_url)
        with self.session_factory() as session:
            if self.current_page_key == "incidents":
                incident = session.get(Incident, selected.item_id)
                if incident is None:
                    raise RuntimeError("Selected incident no longer exists")
                payload = build_incident_cot(incident)
                summary = f"Published incident {incident.id} as CoT"
                kwargs = {
                    "incident_id": incident.id,
                    "payload": {
                        "uid": incident.external_uid
                        or f"section4-incident-{incident.id}",
                        "cot_url": self.cot_url,
                        "kind": "incident",
                    },
                }
            elif self.current_page_key == "capabilities":
                capability = session.get(Capability, selected.item_id)
                if capability is None:
                    raise RuntimeError("Selected capability no longer exists")
                payload = build_capability_cot(capability)
                summary = f"Published capability {capability.id} as CoT"
                kwargs = {
                    "capability_id": capability.id,
                    "payload": {
                        "uid": f"section4-capability-{capability.id}",
                        "cot_url": self.cot_url,
                        "kind": "capability",
                    },
                }
            elif self.current_page_key == "jobs":
                job = session.get(Job, selected.item_id)
                if job is None:
                    raise RuntimeError("Selected job no longer exists")
                payload = build_job_cot(job)
                summary = f"Published job {job.id} as CoT"
                kwargs = {
                    "incident_id": job.incident_id,
                    "job_id": job.id,
                    "payload": {
                        "uid": f"section4-job-{job.id}",
                        "cot_url": self.cot_url,
                        "kind": "job",
                    },
                }
            else:
                raise RuntimeError(
                    "Publish is supported on incidents, capabilities, and jobs"
                )

            await runtime.start()
            try:
                await runtime.publish(payload)
            finally:
                await runtime.stop()

            record_published_cot(
                session,
                summary=summary,
                **kwargs,
            )
        record_label = PUBLISH_LABELS.get(self.current_page_key, "record")
        return f"Published {record_label} {selected.item_id}"

    def _publish_selected(self) -> None:
        """Publish the selected item and show status."""

        try:
            message = asyncio.run(self._publish_selected_async())
        except Exception as exc:  # pragma: no cover - UI reporting path
            self._set_footer(str(exc), style="error")
            return
        self._refresh_page()
        self._set_footer(message, style="success")

    def _unhandled_input(self, key: str) -> None:
        """Handle global keyboard shortcuts."""

        if key in {"q", "Q"}:
            raise urwid.ExitMainLoop()
        if key in {"r", "R"}:
            self._refresh_page()
            return
        if key in {"l", "right", "tab"}:
            self._move_focus_horizontal(1)
            return
        if key in {"h", "left", "shift tab"}:
            self._move_focus_horizontal(-1)
            return
        if key == "enter":
            self._activate_focused()
            return
        if key in {"p", "P"}:
            self._publish_selected()

    def run(self) -> None:
        """Start the urwid main loop."""

        self.loop.run()


def run_tui(db_path: str | Path, cot_url: str) -> None:
    """Start the s4net console for the given SQLite database."""

    create_all(db_path)
    session_factory = create_session_factory(db_path)
    S4NetTUIApplication(session_factory, cot_url).run()
