"""Draft LXDR-Link envelope types."""

from __future__ import annotations

import json
from dataclasses import dataclass, field
from typing import Any

from section4.lxdr.codec.text_burst import render_request_container
from section4.lxdr.request import LXDRRequestContainer
from section4.lxdr.types import DeliveryMethod, LinkRepresentation, LinkState


@dataclass(slots=True)
class LXDRLinkFrame:
    """Transport-agnostic link frame for one or more LXDR core payloads."""

    link_message_id: str
    sender_id: str
    recipient_id: str
    created_at_local: str
    delivery_method: DeliveryMethod
    representation: LinkRepresentation
    payload_count: int
    payloads: list[str] = field(default_factory=list)
    payload_refs: list[str] = field(default_factory=list)
    fragment_index: int = 1
    fragment_count: int = 1
    correlation_id: str | None = None
    ack_of: str | None = None
    sync_of: str | None = None
    attempt_count: int = 0
    integrity_metadata: dict[str, Any] = field(default_factory=dict)
    confidentiality_metadata: dict[str, Any] = field(default_factory=dict)
    state: LinkState = LinkState.CREATED

    def validate(self) -> None:
        """Validate basic frame coherence."""

        if self.payload_count < 1:
            raise ValueError("Link frame must carry at least one payload")
        if self.payload_count != len(self.payloads) + len(self.payload_refs):
            raise ValueError(
                "Link frame payload_count does not match payload contents"
            )
        if self.fragment_index < 1 or self.fragment_count < 1:
            raise ValueError("Fragment indices must be positive")
        if self.fragment_index > self.fragment_count:
            raise ValueError("Fragment index cannot exceed fragment count")

    def to_dict(self) -> dict[str, object]:
        """Return a structured LXDR-Link frame representation."""

        self.validate()
        return {
            "link_message_id": self.link_message_id,
            "sender_id": self.sender_id,
            "recipient_id": self.recipient_id,
            "created_at_local": self.created_at_local,
            "delivery_method": self.delivery_method.value,
            "representation": self.representation.value,
            "payload_count": self.payload_count,
            "payloads": self.payloads,
            "payload_refs": self.payload_refs,
            "fragment_index": self.fragment_index,
            "fragment_count": self.fragment_count,
            "correlation_id": self.correlation_id,
            "ack_of": self.ack_of,
            "sync_of": self.sync_of,
            "attempt_count": self.attempt_count,
            "integrity_metadata": self.integrity_metadata,
            "confidentiality_metadata": self.confidentiality_metadata,
            "state": self.state.value,
        }

    @classmethod
    def from_dict(cls, data: dict[str, object]) -> LXDRLinkFrame:
        """Create a link frame from a structured representation."""

        payloads = data.get("payloads", [])
        payload_refs = data.get("payload_refs", [])
        if not isinstance(payloads, list):
            raise ValueError("Link frame payloads must be a list")
        if not isinstance(payload_refs, list):
            raise ValueError("Link frame payload_refs must be a list")

        frame = cls(
            link_message_id=str(data["link_message_id"]),
            sender_id=str(data["sender_id"]),
            recipient_id=str(data["recipient_id"]),
            created_at_local=str(data["created_at_local"]),
            delivery_method=DeliveryMethod(str(data["delivery_method"])),
            representation=LinkRepresentation(str(data["representation"])),
            payload_count=int(data["payload_count"]),
            payloads=[str(payload) for payload in payloads],
            payload_refs=[str(payload_ref) for payload_ref in payload_refs],
            fragment_index=int(data.get("fragment_index", 1)),
            fragment_count=int(data.get("fragment_count", 1)),
            correlation_id=(
                str(data["correlation_id"])
                if data.get("correlation_id") is not None
                else None
            ),
            ack_of=(
                str(data["ack_of"])
                if data.get("ack_of") is not None
                else None
            ),
            sync_of=(
                str(data["sync_of"])
                if data.get("sync_of") is not None
                else None
            ),
            attempt_count=int(data.get("attempt_count", 0)),
            integrity_metadata=_require_mapping(
                data.get("integrity_metadata", {}),
                "integrity_metadata",
            ),
            confidentiality_metadata=_require_mapping(
                data.get("confidentiality_metadata", {}),
                "confidentiality_metadata",
            ),
            state=LinkState(str(data.get("state", LinkState.CREATED.value))),
        )
        frame.validate()
        return frame

    def serialize(self) -> bytes:
        """Serialize the link frame into deterministic JSON bytes."""

        return json.dumps(
            self.to_dict(),
            sort_keys=True,
            separators=(",", ":"),
        ).encode("utf-8")

    @classmethod
    def parse(cls, data: bytes) -> LXDRLinkFrame:
        """Parse deterministic JSON bytes into a link frame."""

        parsed = json.loads(data.decode("utf-8"))
        if not isinstance(parsed, dict):
            raise ValueError("Link frame payload must decode to an object")
        return cls.from_dict(parsed)

    def embedded_requests(self) -> list[LXDRRequestContainer]:
        """Decode embedded request payloads into request containers."""

        requests: list[LXDRRequestContainer] = []
        for payload in self.payloads:
            if self.representation is LinkRepresentation.INLINE_TEXT:
                requests.append(LXDRRequestContainer.parse_text(payload))
            elif self.representation is LinkRepresentation.INLINE_STRUCTURED:
                requests.append(_decode_structured_request_payload(payload))
            elif self.representation is LinkRepresentation.BUNDLE:
                payload_data = json.loads(payload)
                if not isinstance(payload_data, dict):
                    raise ValueError("Bundle payload must be an object")
                bundle_requests = payload_data.get("requests")
                if not isinstance(bundle_requests, list):
                    raise ValueError(
                        "Bundle payload requests must be a list"
                    )
                requests.extend(
                    _decode_request_dict(request_payload)
                    for request_payload in bundle_requests
                )
            else:
                raise ValueError(
                    "Embedded request decoding is not defined for this "
                    "representation"
                )
        return requests


def embed_inline_text(
    frame: LXDRLinkFrame,
    requests: list[LXDRRequestContainer],
) -> LXDRLinkFrame:
    """Return a copy-like frame populated with canonical text payloads."""

    payloads = [render_request_container(request) for request in requests]
    return LXDRLinkFrame(
        link_message_id=frame.link_message_id,
        sender_id=frame.sender_id,
        recipient_id=frame.recipient_id,
        created_at_local=frame.created_at_local,
        delivery_method=frame.delivery_method,
        representation=LinkRepresentation.INLINE_TEXT,
        payload_count=len(payloads) + len(frame.payload_refs),
        payloads=payloads,
        payload_refs=list(frame.payload_refs),
        fragment_index=frame.fragment_index,
        fragment_count=frame.fragment_count,
        correlation_id=frame.correlation_id,
        ack_of=frame.ack_of,
        sync_of=frame.sync_of,
        attempt_count=frame.attempt_count,
        integrity_metadata=dict(frame.integrity_metadata),
        confidentiality_metadata=dict(frame.confidentiality_metadata),
        state=frame.state,
    )


def embed_inline_structured(
    frame: LXDRLinkFrame,
    requests: list[LXDRRequestContainer],
) -> LXDRLinkFrame:
    """Return a copy-like frame populated with structured JSON payloads."""

    payloads = [
        json.dumps(
            request.to_dict(),
            sort_keys=True,
            separators=(",", ":"),
        )
        for request in requests
    ]
    return LXDRLinkFrame(
        link_message_id=frame.link_message_id,
        sender_id=frame.sender_id,
        recipient_id=frame.recipient_id,
        created_at_local=frame.created_at_local,
        delivery_method=frame.delivery_method,
        representation=LinkRepresentation.INLINE_STRUCTURED,
        payload_count=len(payloads) + len(frame.payload_refs),
        payloads=payloads,
        payload_refs=list(frame.payload_refs),
        fragment_index=frame.fragment_index,
        fragment_count=frame.fragment_count,
        correlation_id=frame.correlation_id,
        ack_of=frame.ack_of,
        sync_of=frame.sync_of,
        attempt_count=frame.attempt_count,
        integrity_metadata=dict(frame.integrity_metadata),
        confidentiality_metadata=dict(frame.confidentiality_metadata),
        state=frame.state,
    )


def embed_bundle(
    frame: LXDRLinkFrame,
    requests: list[LXDRRequestContainer],
) -> LXDRLinkFrame:
    """Return a sync bundle carrying one or more structured requests."""

    bundle_payload = json.dumps(
        {
            "bundle_type": "LXDR_SYNC_BUNDLE",
            "request_count": len(requests),
            "requests": [request.to_dict() for request in requests],
        },
        sort_keys=True,
        separators=(",", ":"),
    )
    return LXDRLinkFrame(
        link_message_id=frame.link_message_id,
        sender_id=frame.sender_id,
        recipient_id=frame.recipient_id,
        created_at_local=frame.created_at_local,
        delivery_method=frame.delivery_method,
        representation=LinkRepresentation.BUNDLE,
        payload_count=1 + len(frame.payload_refs),
        payloads=[bundle_payload],
        payload_refs=list(frame.payload_refs),
        fragment_index=frame.fragment_index,
        fragment_count=frame.fragment_count,
        correlation_id=frame.correlation_id,
        ack_of=frame.ack_of,
        sync_of=frame.sync_of,
        attempt_count=frame.attempt_count,
        integrity_metadata=dict(frame.integrity_metadata),
        confidentiality_metadata=dict(frame.confidentiality_metadata),
        state=frame.state,
    )


def _require_mapping(
    value: object,
    field_name: str,
) -> dict[str, Any]:
    """Validate and normalize mapping-like metadata values."""

    if not isinstance(value, dict):
        raise ValueError(f"Link frame {field_name} must be a mapping")
    return {str(key): value[key] for key in value}


def _decode_structured_request_payload(payload: str) -> LXDRRequestContainer:
    """Decode one structured request payload."""

    payload_data = json.loads(payload)
    return _decode_request_dict(payload_data)


def _decode_request_dict(payload_data: object) -> LXDRRequestContainer:
    """Decode one request mapping into a request container."""

    if not isinstance(payload_data, dict):
        raise ValueError("Structured request payload must be an object")
    return LXDRRequestContainer.from_dict(payload_data)
