"""Public codec helpers for LXDR-Link serialization and parsing."""

from __future__ import annotations

from section4.lxdr.link import LXDRLinkFrame


def serialize_link_frame(frame: LXDRLinkFrame) -> bytes:
    """Serialize a link frame."""

    return frame.serialize()


def parse_link_frame(data: bytes) -> LXDRLinkFrame:
    """Parse a serialized link frame."""

    return LXDRLinkFrame.parse(data)
