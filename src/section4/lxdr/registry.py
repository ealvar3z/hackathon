"""Canonical LXDR field registry based on ADRIAN Appendix F."""

from __future__ import annotations

from dataclasses import dataclass


@dataclass(frozen=True, slots=True)
class CanonicalFieldEntry:
    """Executable representation of an ADRIAN canonical field."""

    registry_id: str
    category: str
    data_block: str
    data_element: str
    data_field: str
    field_label: str
    field_order: int | None
    source_reference: str


def build_registry_id(
    category: str,
    data_block: str,
    data_element: str,
    data_field: str,
) -> str:
    """Build the draft canonical registry identifier."""

    return "::".join((category, data_block, data_element, data_field))


def canonical_field(
    *,
    category: str,
    data_block: str,
    data_element: str,
    data_field: str,
    field_label: str,
    field_order: int | None,
    source_reference: str,
) -> CanonicalFieldEntry:
    """Create a canonical field entry with a stable registry ID."""

    return CanonicalFieldEntry(
        registry_id=build_registry_id(
            category=category,
            data_block=data_block,
            data_element=data_element,
            data_field=data_field,
        ),
        category=category,
        data_block=data_block,
        data_element=data_element,
        data_field=data_field,
        field_label=field_label,
        field_order=field_order,
        source_reference=source_reference,
    )


APPENDIX_F_FIELDS: tuple[CanonicalFieldEntry, ...] = (
    canonical_field(
        category="Department of Defense Activity Address Directory",
        data_block="Activity",
        data_element="Header",
        data_field="1. Department of Defense Activity Address Code",
        field_label="department_of_defense_activity_address_code",
        field_order=1,
        source_reference="Appendix F",
    ),
    canonical_field(
        category="Department of Defense Activity Address Directory",
        data_block="Activity",
        data_element="Header",
        data_field="2. Unit Identifier Code",
        field_label="unit_identifier_code",
        field_order=2,
        source_reference="Appendix F",
    ),
    canonical_field(
        category="Department of Defense Activity Address Directory",
        data_block="Activity",
        data_element="COMMRI",
        data_field="1. Bill",
        field_label="bill",
        field_order=1,
        source_reference="Appendix F",
    ),
    canonical_field(
        category="Department of Defense Activity Address Directory",
        data_block="Activity",
        data_element="COMMRI",
        data_field="2. Data Pattern",
        field_label="data_pattern",
        field_order=2,
        source_reference="Appendix F",
    ),
    canonical_field(
        category="Department of Defense Activity Address Directory",
        data_block="Type Address Code 1 (Postal/Mail)",
        data_element="Type Address Code 1 (Postal/Mail)",
        data_field="1. T1_Line 1",
        field_label="t1_line_1",
        field_order=1,
        source_reference="Appendix F",
    ),
    canonical_field(
        category="Department of Defense Activity Address Directory",
        data_block="Type Address Code 1 (Postal/Mail)",
        data_element="Type Address Code 1 (Postal/Mail)",
        data_field="6. T1_City",
        field_label="t1_city",
        field_order=6,
        source_reference="Appendix F",
    ),
)


class CanonicalFieldRegistry:
    """In-memory canonical registry for draft LXDR implementation."""

    def __init__(
        self,
        entries: tuple[CanonicalFieldEntry, ...] = APPENDIX_F_FIELDS,
    ) -> None:
        self._entries = {entry.registry_id: entry for entry in entries}

    def all(self) -> tuple[CanonicalFieldEntry, ...]:
        """Return all canonical entries."""

        return tuple(self._entries.values())

    def get(self, registry_id: str) -> CanonicalFieldEntry | None:
        """Look up a canonical entry by registry ID."""

        return self._entries.get(registry_id)
