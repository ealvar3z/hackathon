"""Persistence-backed helpers for the LXDR router."""

from __future__ import annotations

from sqlalchemy import select
from sqlalchemy.orm import Session, sessionmaker

from section4.lxdr.link import LXDRLinkFrame
from section4.lxdr.request import LXDRRequestContainer
from section4.lxdr.router import LXDRRouter
from section4.storage.tables import LXDRInboundFrame, LXDROutboundFrame


class PersistentLXDRRouter:
    """Bridge the in-memory router with persisted inbox/outbox state."""

    def __init__(
        self,
        session_factory: sessionmaker[Session],
        memory_router: LXDRRouter,
    ) -> None:
        self.session_factory = session_factory
        self.memory_router = memory_router

    def queue_request(
        self,
        request: LXDRRequestContainer,
        recipient_id: str,
        created_at_local: str,
    ):
        """Queue a request in memory and persist the outbound frame."""

        frame = self.memory_router.queue_request(
            request=request,
            recipient_id=recipient_id,
            created_at_local=created_at_local,
        )
        with self.session_factory() as session:
            session.add(
                LXDROutboundFrame(
                    link_message_id=frame.link_message_id,
                    sender_id=frame.sender_id,
                    recipient_id=frame.recipient_id,
                    request_unique_identification_local=(
                        request.header.request_unique_identification_local
                    ),
                    request_unique_identification_sync=(
                        request.header.request_unique_identification_sync
                    ),
                    delivery_method=frame.delivery_method.value,
                    representation=frame.representation.value,
                    state=frame.state.value,
                    created_at_local=frame.created_at_local,
                    attempt_count=frame.attempt_count,
                    correlation_id=frame.correlation_id,
                    payload_json=frame.to_dict(),
                )
            )
            session.commit()
        return frame

    def record_inbound_frame(
        self,
        link_message_id: str,
        sender_id: str,
        recipient_id: str,
        payload_count: int,
    ) -> None:
        """Persist a received inbound frame identity for dedupe."""

        with self.session_factory() as session:
            existing = session.scalar(
                select(LXDRInboundFrame).where(
                    LXDRInboundFrame.link_message_id == link_message_id
                )
            )
            if existing is not None:
                return
            session.add(
                LXDRInboundFrame(
                    link_message_id=link_message_id,
                    sender_id=sender_id,
                    recipient_id=recipient_id,
                    payload_count=payload_count,
                )
            )
            session.commit()

    def apply_sync_update(
        self,
        local_request_id: str,
        sync_request_id: str,
    ) -> int:
        """Apply a sync identifier in memory and persist it."""

        request = self.memory_router.request_by_local_id(local_request_id)
        if request is None:
            return 0

        request.header.apply_sync_identifier(sync_request_id)

        with self.session_factory() as session:
            row = session.scalar(
                select(LXDROutboundFrame).where(
                    LXDROutboundFrame.request_unique_identification_local
                    == local_request_id
                )
            )
            if row is None:
                return 0
            row.request_unique_identification_sync = sync_request_id
            row.state = "SYNCED"
            row.payload_json = request.to_dict()
            session.commit()

        return 1

    def mark_frame_sending(
        self,
        link_message_id: str,
        attempted_at: str,
    ) -> int:
        """Mark an outbound frame as actively sending."""

        with self.session_factory() as session:
            row = self._get_outbound_row(session, link_message_id)
            if row is None:
                return 0
            row.state = "SENDING"
            row.attempt_count += 1
            row.last_attempt_at = attempted_at
            row.last_error = None
            self._update_payload_state(
                row,
                state="SENDING",
                attempt_count=row.attempt_count,
            )
            session.commit()
        return 1

    def mark_frame_sent(self, link_message_id: str) -> int:
        """Mark an outbound frame as sent."""

        with self.session_factory() as session:
            row = self._get_outbound_row(session, link_message_id)
            if row is None:
                return 0
            row.state = "SENT"
            self._update_payload_state(row, state="SENT")
            session.commit()
        return 1

    def mark_frame_failed(
        self,
        link_message_id: str,
        error_message: str,
    ) -> int:
        """Mark an outbound frame as failed with an error message."""

        with self.session_factory() as session:
            row = self._get_outbound_row(session, link_message_id)
            if row is None:
                return 0
            row.state = "FAILED"
            row.last_error = error_message
            self._update_payload_state(row, state="FAILED")
            session.commit()
        return 1

    def retryable_outbound_frames(self) -> list[LXDRLinkFrame]:
        """Return queued or failed frames suitable for retry scheduling."""

        with self.session_factory() as session:
            rows = session.scalars(
                select(LXDROutboundFrame).where(
                    LXDROutboundFrame.state.in_(("QUEUED", "FAILED"))
                )
            ).all()
        return [
            LXDRLinkFrame.from_dict(row.payload_json)
            for row in rows
        ]

    @staticmethod
    def _get_outbound_row(
        session: Session,
        link_message_id: str,
    ) -> LXDROutboundFrame | None:
        """Fetch one persisted outbound row by link message ID."""

        return session.scalar(
            select(LXDROutboundFrame).where(
                LXDROutboundFrame.link_message_id == link_message_id
            )
        )

    @staticmethod
    def _update_payload_state(
        row: LXDROutboundFrame,
        *,
        state: str,
        attempt_count: int | None = None,
    ) -> None:
        """Keep the stored payload metadata aligned with row state."""

        payload = dict(row.payload_json)
        payload["state"] = state
        if attempt_count is not None:
            payload["attempt_count"] = attempt_count
        row.payload_json = payload
