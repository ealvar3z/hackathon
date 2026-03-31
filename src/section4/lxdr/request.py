"""LXDR core request container types."""

from __future__ import annotations

from dataclasses import dataclass

from section4.lxdr.header import LXDRHeader
from section4.lxdr.segments.base import LXDRSegment
from section4.lxdr.segments.registry import (
    parse_segment_text,
    segment_from_dict,
)


@dataclass(slots=True)
class LXDRRequestContainer:
    """Core LXDR request container: header plus child segments."""

    header: LXDRHeader
    segments: list[LXDRSegment]

    def validate(self) -> None:
        """Validate header and segment count coherence."""

        self.header.validate()
        if len(self.segments) != self.header.request_segments:
            raise ValueError(
                "Header request segment count does not match payload"
            )

    def to_dict(self) -> dict[str, object]:
        """Return a structured request container representation."""

        self.validate()
        return {
            "header": self.header.to_dict(),
            "segments": [segment.to_dict() for segment in self.segments],
        }

    @classmethod
    def from_dict(cls, data: dict[str, object]) -> LXDRRequestContainer:
        """Create a request container from structured data."""

        header_data = data.get("header")
        segments_data = data.get("segments")
        if not isinstance(header_data, dict):
            raise ValueError("Request container header must be a mapping")
        if not isinstance(segments_data, list):
            raise ValueError("Request container segments must be a list")

        return cls(
            header=LXDRHeader.from_dict(header_data),
            segments=[
                segment_from_dict(segment)
                for segment in segments_data
                if isinstance(segment, dict)
            ],
        )

    @classmethod
    def parse_text(cls, text: str) -> LXDRRequestContainer:
        """Parse a canonical text burst into a request container."""

        lines = [line.strip() for line in text.splitlines() if line.strip()]
        if not lines:
            raise ValueError("Request container text is empty")

        header = LXDRHeader.parse_text(lines[0])
        segments = [parse_segment_text(line) for line in lines[1:]]
        return cls(header=header, segments=segments)
