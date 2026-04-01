"""Persistence-backed helpers for the LXDR router."""

from __future__ import annotations

from sqlalchemy import select
from sqlalchemy.orm import Session, sessionmaker

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
