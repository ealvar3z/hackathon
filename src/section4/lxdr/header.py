"""LXDR request header model derived from ADRIAN Section 3.3.3."""

from __future__ import annotations

from dataclasses import dataclass


@dataclass(slots=True)
class LXDRHeader:
    """ADRIAN-derived request header."""

    date_request_created_local: str
    time_request_created_local: str
    physical_location_of_requestor: str
    request_unique_identification_local: str
    request_priority: str
    element_unit_identification_callsign: str
    request_segments: int
    request_unique_identification_sync: str | None = None

    def validate(self) -> None:
        """Validate obvious ADRIAN header constraints."""

        if len(self.date_request_created_local) != 9:
            raise ValueError("Header local date must be CCYYMMMDD")
        if len(self.time_request_created_local) != 8:
            raise ValueError("Header local time must be HHMMSSDD")
        if len(self.physical_location_of_requestor) != 14:
            raise ValueError("Header MGRS reference must be 14 characters")
        if len(self.request_unique_identification_local) != 10:
            raise ValueError("Header local request ID must be 10 characters")
        if len(self.request_priority) != 2:
            raise ValueError("Header priority must be 2 characters")
        if len(self.element_unit_identification_callsign) != 4:
            raise ValueError("Header callsign must be 4 characters")
        if not 1 <= self.request_segments <= 9:
            raise ValueError("Header request segment count must be 1..9")
        if (
            self.request_unique_identification_sync is not None
            and len(self.request_unique_identification_sync) != 12
        ):
            raise ValueError("Header sync request ID must be 12 characters")

    def render_text(self) -> str:
        """Render the ADRIAN header text burst."""

        self.validate()
        return "-".join(
            (
                self.date_request_created_local,
                self.time_request_created_local,
                self.physical_location_of_requestor,
                self.request_unique_identification_local,
                self.request_priority,
                self.element_unit_identification_callsign,
                f"{self.request_segments:02d}",
            )
        )

    def to_dict(self) -> dict[str, str | int | None]:
        """Return a structured header representation."""

        self.validate()
        return {
            "date_request_created_local": self.date_request_created_local,
            "time_request_created_local": self.time_request_created_local,
            "physical_location_of_requestor": (
                self.physical_location_of_requestor
            ),
            "request_unique_identification_local": (
                self.request_unique_identification_local
            ),
            "request_priority": self.request_priority,
            "element_unit_identification_callsign": (
                self.element_unit_identification_callsign
            ),
            "request_segments": self.request_segments,
            "request_unique_identification_sync": (
                self.request_unique_identification_sync
            ),
        }

    @classmethod
    def from_dict(cls, data: dict[str, str | int | None]) -> LXDRHeader:
        """Create a header from a structured representation."""

        return cls(
            date_request_created_local=str(
                data["date_request_created_local"]
            ),
            time_request_created_local=str(
                data["time_request_created_local"]
            ),
            physical_location_of_requestor=str(
                data["physical_location_of_requestor"]
            ),
            request_unique_identification_local=str(
                data["request_unique_identification_local"]
            ),
            request_priority=str(data["request_priority"]),
            element_unit_identification_callsign=str(
                data["element_unit_identification_callsign"]
            ),
            request_segments=int(data["request_segments"]),
            request_unique_identification_sync=(
                str(data["request_unique_identification_sync"])
                if data["request_unique_identification_sync"] is not None
                else None
            ),
        )

    @classmethod
    def parse_text(cls, text: str) -> LXDRHeader:
        """Parse an ADRIAN header text burst."""

        parts = text.split("-")
        if len(parts) != 7:
            raise ValueError("Header text must contain 7 fields")

        return cls(
            date_request_created_local=parts[0],
            time_request_created_local=parts[1],
            physical_location_of_requestor=parts[2],
            request_unique_identification_local=parts[3],
            request_priority=parts[4],
            element_unit_identification_callsign=parts[5],
            request_segments=int(parts[6]),
        )

    def apply_sync_identifier(self, sync_identifier: str) -> None:
        """Apply the synchronized enterprise request identifier."""

        if len(sync_identifier) != 12:
            raise ValueError("Sync identifier must be 12 characters")
        self.request_unique_identification_sync = sync_identifier
