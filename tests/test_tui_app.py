"""Tests for the s4net terminal UI contract."""

from __future__ import annotations

from pathlib import Path

from section4.lxdr.router import LXDRRouter
from section4.lxdr.router_store import PersistentLXDRRouter
from section4.storage import create_all, create_session_factory
from section4.storage.seed import seed_demo_data
from section4.storage.tables import LogisticsStatusReport
from section4.tui.app import PAGES, S4NetTUIApplication
from tests.test_lxdr_router_store import build_sample_request


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

    assert app.header.original_widget.text == " s4net "
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
    dashboard_texts = [widget.text for widget in app.list_walker]
    assert "SUPPORTABILITY" in dashboard_texts
    assert "Status reports: 0" in dashboard_texts
    assert "LXDR LIFECYCLE" in dashboard_texts
    assert "Outbound frames: 0" in dashboard_texts
    assert "Synced: 0" in dashboard_texts


def test_request_page_uses_adrian_request_language(tmp_path: Path) -> None:
    """The incident browser should now read as an ADRIAN request browser."""

    db_path = tmp_path / "section4-test.db"
    create_all(db_path)
    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        seed_demo_data(session)
        session.add(
            LogisticsStatusReport(
                node_id="atak-node-3",
                callsign="NODE3",
                report_type="on_hand_usage",
                item_reference="BRKT-4421",
                on_hand_quantity=0,
                usage_rate_per_day=1.0,
                required_quantity=1,
                required_by_local="2027OCT15",
                delivery_method="direct",
                transport_mode="ground",
                payload_json={},
            )
        )
        session.commit()

    router = PersistentLXDRRouter(
        session_factory=session_factory,
        memory_router=LXDRRouter(sender_id="NODE1"),
    )
    router.queue_request(
        request=build_sample_request("3838JBNM5X"),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )

    app = S4NetTUIApplication(session_factory, "udp://239.2.3.1:6969")

    app._on_nav_focus(None, 1)

    assert app.page_title.text == "Requests"
    detail_texts = [widget.text for widget in app.detail_walker]
    assert any(text == "Request ID: 3838JBNM5X" for text in detail_texts)
    assert any(
        text == "Request type: PM" for text in detail_texts
    )
    assert any(text == "SUPPORTABILITY" for text in detail_texts)
    assert any(
        text == "Recommended COA: fabricate" for text in detail_texts
    )
    assert "PM 3838JBNM5X P02" in app.current_items[0].label


def test_sync_log_page_shows_persisted_lxdr_frames(tmp_path: Path) -> None:
    """The sync log page should show persisted outbound frame lifecycle."""

    db_path = tmp_path / "section4-test.db"
    create_all(db_path)
    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        seed_demo_data(session)

    router = PersistentLXDRRouter(
        session_factory=session_factory,
        memory_router=LXDRRouter(sender_id="NODE1"),
    )
    frame = router.queue_request(
        request=build_sample_request(),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )
    router.mark_frame_sending(
        link_message_id=frame.link_message_id,
        attempted_at="2027OCT13T15480352",
    )

    app = S4NetTUIApplication(session_factory, "udp://239.2.3.1:6969")
    app._on_nav_focus(None, 5)

    assert app.page_title.text == "Sync Log"
    assert app.current_items
    assert app.current_items[0].label.startswith("OUT ")
    detail_texts = [widget.text for widget in app.detail_walker]
    assert any(text == "State: SENDING" for text in detail_texts)
    assert any(text == "Attempt count: 1" for text in detail_texts)


def test_request_page_surfaces_sync_identifier_clearly(tmp_path: Path) -> None:
    """A synchronized ADRIAN request should show the enterprise ID."""

    db_path = tmp_path / "section4-test.db"
    create_all(db_path)
    session_factory = create_session_factory(db_path)
    with session_factory() as session:
        seed_demo_data(session)

    router = PersistentLXDRRouter(
        session_factory=session_factory,
        memory_router=LXDRRouter(sender_id="NODE1"),
    )
    router.queue_request(
        request=build_sample_request("3838JBNM5X"),
        recipient_id="ALOC",
        created_at_local="2027OCT13T15470352",
    )
    router.apply_sync_update(
        local_request_id="3838JBNM5X",
        sync_request_id="ABCD1234EFGH",
    )

    app = S4NetTUIApplication(session_factory, "udp://239.2.3.1:6969")
    app._on_nav_focus(None, 1)

    assert "-> ABCD1234EFGH" in app.current_items[0].label
    detail_texts = [widget.text for widget in app.detail_walker]
    assert any(
        text == "Sync request ID: ABCD1234EFGH"
        for text in detail_texts
    )
