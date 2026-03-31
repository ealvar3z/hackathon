"""Shared LXDR enums derived from the current draft protocol."""

from __future__ import annotations

from enum import Enum


class Usage(str, Enum):
    """Field usage constraints from ADRIAN tables."""

    MUST_USE = "Must Use"
    OPTIONAL = "Optional"


class ExchangeBehavior(str, Enum):
    """Wire exchange behavior from ADRIAN tables."""

    PASS = "Pass"
    NOT_PASSED = "Not Passed"
    SYNC_RESPONSE = "Sync Response"


class DeliveryMethod(str, Enum):
    """Draft LXDR-Link delivery methods."""

    OPPORTUNISTIC = "OPPORTUNISTIC"
    DIRECT = "DIRECT"
    RELAYED = "RELAYED"
    SYNCHRONIZATION = "SYNCHRONIZATION"
    EXPORT = "EXPORT"


class LinkRepresentation(str, Enum):
    """Draft LXDR-Link payload representations."""

    INLINE_TEXT = "INLINE_TEXT"
    INLINE_STRUCTURED = "INLINE_STRUCTURED"
    BUNDLE = "BUNDLE"


class LinkState(str, Enum):
    """Draft LXDR-Link delivery lifecycle states."""

    CREATED = "CREATED"
    QUEUED = "QUEUED"
    SENDING = "SENDING"
    SENT = "SENT"
    ACKNOWLEDGED = "ACKNOWLEDGED"
    DELIVERED = "DELIVERED"
    SYNCED = "SYNCED"
    REJECTED = "REJECTED"
    FAILED = "FAILED"
    EXPIRED = "EXPIRED"
    CANCELLED = "CANCELLED"
