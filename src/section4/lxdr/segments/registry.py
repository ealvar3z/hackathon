"""Initial ADRIAN-derived LXDR segment registry."""

from __future__ import annotations

from section4.lxdr.segments.base import (
    LXDRSegment,
    SegmentFieldSpec,
    SegmentSpec,
)
from section4.lxdr.types import ExchangeBehavior, Usage


def _field(
    field_name: str,
    usage: Usage,
    exchange: ExchangeBehavior,
    min_length: int,
    max_length: int,
) -> SegmentFieldSpec:
    """Create one segment field specification."""

    return SegmentFieldSpec(
        field_name=field_name,
        usage=usage,
        exchange=exchange,
        min_length=min_length,
        max_length=max_length,
    )


PAX_MOVEMENT = SegmentSpec(
    name="mobility_pax",
    request_type="PM",
    source_reference="ADRIAN 3.4.3.1 / Table 13 / Table 14",
    fields=(
        _field("segment_number", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 1),
        _field("request_type", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field("request_priority", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field("zap_or_edipi", Usage.MUST_USE, ExchangeBehavior.PASS, 8, 10),
        _field(
            "earliest_departure_date",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            9,
            9,
        ),
        _field(
            "latest_departure_date",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            9,
            9,
        ),
        _field(
            "departure_location",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            14,
            14,
        ),
        _field(
            "destination_location",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            14,
            14,
        ),
        _field(
            "estimated_baggage_weight_lbs",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            3,
            3,
        ),
        _field(
            "hazardous_material_type",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            1,
            1,
        ),
    ),
)


CARGO_MOVEMENT = SegmentSpec(
    name="mobility_cargo",
    request_type="CM",
    source_reference="ADRIAN 3.4.3.2 / Table 15 / Table 16",
    fields=(
        _field("segment_number", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 1),
        _field("request_type", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field("request_priority", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field("item_by_niin", Usage.MUST_USE, ExchangeBehavior.PASS, 9, 9),
        _field("item_quantity", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 5),
        _field("serial_number", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 10),
        _field("gross_weight_lbs", Usage.MUST_USE, ExchangeBehavior.PASS, 5, 5),
        _field(
            "actual_height_inches",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            3,
            3,
        ),
        _field(
            "actual_width_inches",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            3,
            3,
        ),
        _field(
            "actual_length_inches",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            3,
            3,
        ),
        _field("hmic", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 1),
        _field("handling", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 1),
        _field(
            "earliest_departure_date",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            9,
            9,
        ),
        _field(
            "latest_departure_date",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            9,
            9,
        ),
        _field(
            "departure_location",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            14,
            14,
        ),
        _field(
            "destination_location",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            14,
            14,
        ),
    ),
)


SUPPLY_REQUEST = SegmentSpec(
    name="supply_request",
    request_type="SR",
    source_reference="ADRIAN 3.5.3 / Table 17",
    fields=(
        _field("segment_number", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 1),
        _field("request_type", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field("request_priority", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field("item_by_niin", Usage.OPTIONAL, ExchangeBehavior.PASS, 9, 9),
        _field("item_quantity", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 5),
        _field(
            "required_date_local",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            9,
            9,
        ),
        _field(
            "delivery_location",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            14,
            14,
        ),
        _field("attachment", Usage.OPTIONAL, ExchangeBehavior.PASS, 1, 1),
        _field("narrative", Usage.OPTIONAL, ExchangeBehavior.PASS, 5, 100),
    ),
)


MAINTENANCE_REQUEST = SegmentSpec(
    name="maintenance_request",
    request_type=None,
    source_reference="ADRIAN 3.6.3 / Table 18",
    fields=(
        _field("segment_number", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 1),
        _field("request_type", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field("request_priority", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field("serial_number", Usage.OPTIONAL, ExchangeBehavior.PASS, 1, 15),
        _field("niin", Usage.MUST_USE, ExchangeBehavior.PASS, 9, 9),
        _field(
            "model_of_equipment",
            Usage.OPTIONAL,
            ExchangeBehavior.PASS,
            1,
            20,
        ),
        _field(
            "item_nomenclature",
            Usage.OPTIONAL,
            ExchangeBehavior.PASS,
            1,
            20,
        ),
        _field("number_of_pieces", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 5),
        _field(
            "equipment_operational_condition",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            1,
            1,
        ),
        _field(
            "date_maintenance_support_required",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            9,
            9,
        ),
        _field(
            "location_of_equipment",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            14,
            14,
        ),
        _field(
            "type_of_maintenance_support",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            2,
            2,
        ),
        _field("type_of_repair", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field(
            "repair_major_defect",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            4,
            4,
        ),
        _field("attachment", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 1),
        _field("narrative", Usage.OPTIONAL, ExchangeBehavior.PASS, 5, 100),
    ),
)


HEALTH_COLLECTION = SegmentSpec(
    name="health_collection",
    request_type="CR",
    source_reference="ADRIAN 3.8.3.1 / Table 33",
    fields=(
        _field("segment_number", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 1),
        _field("request_type", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field("request_priority", Usage.MUST_USE, ExchangeBehavior.PASS, 2, 2),
        _field("zap_or_edipi", Usage.MUST_USE, ExchangeBehavior.PASS, 8, 10),
        _field("last_name", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 15),
        _field("first_name", Usage.MUST_USE, ExchangeBehavior.PASS, 1, 15),
        _field("allergies", Usage.MUST_USE, ExchangeBehavior.PASS, 4, 10),
        _field("date_of_injury", Usage.MUST_USE, ExchangeBehavior.PASS, 9, 9),
        _field("time_of_injury", Usage.MUST_USE, ExchangeBehavior.PASS, 8, 8),
        _field(
            "location_injury_occurred",
            Usage.MUST_USE,
            ExchangeBehavior.PASS,
            14,
            14,
        ),
    ),
)


SEGMENT_REGISTRY: dict[str, SegmentSpec] = {
    segment.name: segment
    for segment in (
        PAX_MOVEMENT,
        CARGO_MOVEMENT,
        SUPPLY_REQUEST,
        MAINTENANCE_REQUEST,
        HEALTH_COLLECTION,
    )
}

REQUEST_TYPE_REGISTRY: dict[str, SegmentSpec] = {
    segment.request_type: segment
    for segment in SEGMENT_REGISTRY.values()
    if segment.request_type is not None
}


def segment_from_dict(data: dict[str, object]) -> LXDRSegment:
    """Rebuild a concrete segment from its structured form."""

    segment_name = str(data["segment_name"])
    spec = SEGMENT_REGISTRY[segment_name]
    values = data.get("values", {})
    if not isinstance(values, dict):
        raise ValueError("Segment values must be a mapping")
    return LXDRSegment(
        spec=spec,
        values={str(key): str(value) for key, value in values.items()},
    )


def parse_segment_text(line: str) -> LXDRSegment:
    """Parse a canonical text segment for known request types."""

    parts = line.split("-")
    if len(parts) < 2:
        raise ValueError("Segment text must contain at least 2 fields")

    request_type = parts[1]
    spec = REQUEST_TYPE_REGISTRY.get(request_type)
    if spec is None:
        raise ValueError(
            "Unknown or unresolved request type for text parsing: "
            f"{request_type}"
        )

    field_names = [
        field.field_name
        for field in spec.fields
        if field.exchange is ExchangeBehavior.PASS
    ]
    if len(parts) != len(field_names):
        raise ValueError(
            "Segment text field count does not match the registered spec"
        )

    return LXDRSegment(
        spec=spec,
        values=dict(zip(field_names, parts, strict=True)),
    )
