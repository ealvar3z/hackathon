"""Base LXDR segment types and specifications."""

from __future__ import annotations

from dataclasses import dataclass, field

from section4.lxdr.types import ExchangeBehavior, Usage


@dataclass(frozen=True, slots=True)
class SegmentFieldSpec:
    """Specification for one segment field."""

    field_name: str
    usage: Usage
    exchange: ExchangeBehavior
    min_length: int
    max_length: int
    repeats: bool = False


@dataclass(frozen=True, slots=True)
class SegmentSpec:
    """Specification for a segment type."""

    name: str
    request_type: str | None
    source_reference: str
    fields: tuple[SegmentFieldSpec, ...]


@dataclass(slots=True)
class LXDRSegment:
    """Concrete segment instance."""

    spec: SegmentSpec
    values: dict[str, str] = field(default_factory=dict)

    def pass_field_names(self) -> list[str]:
        """Return fields that belong in the canonical text burst."""

        return [
            field.field_name
            for field in self.spec.fields
            if field.exchange is ExchangeBehavior.PASS
        ]

    def to_dict(self) -> dict[str, object]:
        """Return a structured segment representation."""

        return {
            "segment_name": self.spec.name,
            "request_type": self.spec.request_type,
            "source_reference": self.spec.source_reference,
            "values": self.values,
        }
