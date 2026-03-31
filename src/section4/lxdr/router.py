"""In-memory LXDR router primitives for DDIL request exchange."""

from __future__ import annotations

import json
from dataclasses import dataclass

from section4.lxdr.link import (
    LXDRLinkFrame,
    embed_bundle,
    embed_inline_structured,
)
from section4.lxdr.request import LXDRRequestContainer
from section4.lxdr.types import DeliveryMethod, LinkRepresentation, LinkState


@dataclass(slots=True)
class OutboxEntry:
    """A queued outbound request and its current link frame."""

    request: LXDRRequestContainer
    frame: LXDRLinkFrame


class LXDRRouter:
    """Minimal in-memory router for outbound and inbound LXDR traffic."""

    def __init__(self, sender_id: str) -> None:
        self.sender_id = sender_id
        self._message_counter = 0
        self._outbox: list[OutboxEntry] = []
        self._inbox_seen: set[str] = set()

    @property
    def outbox_size(self) -> int:
        """Return the number of queued outbound requests."""

        return len(self._outbox)

    def queue_request(
        self,
        request: LXDRRequestContainer,
        recipient_id: str,
        created_at_local: str,
    ) -> LXDRLinkFrame:
        """Queue a request as a structured outbound link frame."""

        frame = LXDRLinkFrame(
            link_message_id=self._next_message_id(),
            sender_id=self.sender_id,
            recipient_id=recipient_id,
            created_at_local=created_at_local,
            delivery_method=DeliveryMethod.DIRECT,
            representation=LinkRepresentation.INLINE_STRUCTURED,
            payload_count=1,
            state=LinkState.QUEUED,
        )
        frame = embed_inline_structured(frame, [request])
        self._outbox.append(OutboxEntry(request=request, frame=frame))
        return frame

    def pending_sync_requests(self) -> list[LXDRRequestContainer]:
        """Return queued requests that do not yet have sync identifiers."""

        return [
            entry.request
            for entry in self._outbox
            if entry.request.header.request_unique_identification_sync is None
        ]

    def request_by_local_id(
        self,
        local_request_id: str,
    ) -> LXDRRequestContainer | None:
        """Return a queued request by its local ADRIAN request identifier."""

        for entry in self._outbox:
            if (
                entry.request.header.request_unique_identification_local
                == local_request_id
            ):
                return entry.request
        return None

    def build_sync_bundle(
        self,
        recipient_id: str,
        created_at_local: str,
        correlation_id: str,
    ) -> LXDRLinkFrame:
        """Build a synchronization bundle from unsynchronized requests."""

        requests = self.pending_sync_requests()
        if not requests:
            raise ValueError("No unsynchronized requests available")

        frame = LXDRLinkFrame(
            link_message_id=self._next_message_id(),
            sender_id=self.sender_id,
            recipient_id=recipient_id,
            created_at_local=created_at_local,
            delivery_method=DeliveryMethod.SYNCHRONIZATION,
            representation=LinkRepresentation.BUNDLE,
            payload_count=1,
            correlation_id=correlation_id,
            state=LinkState.QUEUED,
        )
        return embed_bundle(frame, requests)

    def accept_inbound_frame(self, frame: LXDRLinkFrame) -> bool:
        """Accept an inbound frame unless it has already been seen."""

        if frame.link_message_id in self._inbox_seen:
            return False
        self._inbox_seen.add(frame.link_message_id)
        return True

    def apply_sync_response(self, frame: LXDRLinkFrame) -> int:
        """Apply synchronized request identifiers from an inbound frame."""

        updates = 0
        for payload in frame.payloads:
            payload_data = json.loads(payload)
            if not isinstance(payload_data, dict):
                raise ValueError("Sync response payload must be an object")

            local_request_id = payload_data.get("local_request_id")
            sync_request_id = payload_data.get("sync_request_id")
            if not isinstance(local_request_id, str):
                raise ValueError("Sync response local_request_id is required")
            if not isinstance(sync_request_id, str):
                raise ValueError("Sync response sync_request_id is required")

            for entry in self._outbox:
                if (
                    entry.request.header.request_unique_identification_local
                    == local_request_id
                ):
                    entry.request.header.apply_sync_identifier(sync_request_id)
                    entry.frame.state = LinkState.SYNCED
                    updates += 1
                    break

        return updates

    def _next_message_id(self) -> str:
        """Generate a deterministic local message identity."""

        self._message_counter += 1
        return f"{self.sender_id}-msg-{self._message_counter:06d}"
