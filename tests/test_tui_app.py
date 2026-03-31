"""Tests for the s4net terminal UI contract."""

from __future__ import annotations

from pathlib import Path

from section4.storage import create_all, create_session_factory
from section4.storage.seed import seed_demo_data
from section4.tui.app import PAGES, S4NetTUIApplication


def build_seeded_session_factory(tmp_path: Path):
    """Create a seeded SQLite session factory for TUI tests."""

    db_path = tmp_path / "section4-test.db"
    create_all(db_path)
    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        seed_demo_data(session)
    return session_factory


def test_s4net_header_and_page_labels_match_current_contract(
    tmp_path: Path,
) -> None:
    """The TUI should expose the s4net and ADRIAN-oriented labels."""

    session_factory = build_seeded_session_factory(tmp_path)
    app = S4NetTUIApplication(session_factory, "udp://239.2.3.1:6969")

    assert app.header.original_widget.text == " s4net / section4 "
    assert [page.title for page in PAGES] == [
        "COP",
        "Requests",
        "Capabilities",
        "Tasks",
        "Artifacts",
        "Sync Log",
    ]


def test_dashboard_uses_cop_and_sync_framing(tmp_path: Path) -> None:
    """The default dashboard should render the COP-oriented titles."""

    session_factory = build_seeded_session_factory(tmp_path)
    app = S4NetTUIApplication(session_factory, "udp://239.2.3.1:6969")

    assert app.page_title.text == "Logistics COP"
    assert app.detail_title.text == "Sync and Event Log"
    assert app.list_walker[0].text == "COP SUMMARY"
    assert app.detail_walker[0].text == "LATEST EVENTS"


def test_request_page_uses_adrian_request_language(tmp_path: Path) -> None:
    """The incident browser should now read as an ADRIAN request browser."""

    session_factory = build_seeded_session_factory(tmp_path)
    app = S4NetTUIApplication(session_factory, "udp://239.2.3.1:6969")

    app._on_nav_focus(None, 1)

    assert app.page_title.text == "Requests"
    detail_texts = [widget.text for widget in app.detail_walker]
    assert any(text.startswith("Request ID: ") for text in detail_texts)
    assert any(
        text == "Request UID: incident-alpha-001" for text in detail_texts
    )
    assert "MRZR suspension bracket failure" in app.current_items[0].label
