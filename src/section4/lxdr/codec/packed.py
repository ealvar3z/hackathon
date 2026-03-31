"""Packed codecs for constrained LXDR representations."""

from __future__ import annotations

from section4.lxdr.header import LXDRHeader
from section4.lxdr.segments import PAX_MOVEMENT, LXDRSegment

PACKED_HEADER_SIZE_BYTES = 28
PACKED_PAX_SIZE_BYTES = 40

_MONTHS = {
    "JAN": 1,
    "FEB": 2,
    "MAR": 3,
    "APR": 4,
    "MAY": 5,
    "JUN": 6,
    "JUL": 7,
    "AUG": 8,
    "SEP": 9,
    "OCT": 10,
    "NOV": 11,
    "DEC": 12,
}
_MONTHS_REVERSE = {value: key for key, value in _MONTHS.items()}

# Implementation-specific constrained alphabet for ADRIAN-coded fields.
# The ADRIAN text examples and constraints rely on uppercase letters and
# digits. We exclude I and O to match the callsign/key guidance.
_ALPHABET = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ+"
_ALPHABET_INDEX = {char: index for index, char in enumerate(_ALPHABET)}


class _BitWriter:
    """Minimal bit writer for fixed-width field packing."""

    def __init__(self) -> None:
        self._value = 0
        self._bit_count = 0

    def write(self, value: int, bit_count: int) -> None:
        """Append a fixed-width unsigned integer."""

        if bit_count <= 0:
            raise ValueError("bit_count must be positive")
        if value < 0 or value >= (1 << bit_count):
            raise ValueError(
                f"value {value} does not fit in {bit_count} bits"
            )
        self._value = (self._value << bit_count) | value
        self._bit_count += bit_count

    def to_bytes(self, size_bytes: int) -> bytes:
        """Render the accumulated bits into a fixed-size byte string."""

        total_bits = size_bytes * 8
        if self._bit_count > total_bits:
            raise ValueError("packed value exceeds target byte size")
        padding = total_bits - self._bit_count
        return (self._value << padding).to_bytes(size_bytes, "big")


class _BitReader:
    """Minimal bit reader for fixed-width field unpacking."""

    def __init__(self, data: bytes) -> None:
        self._value = int.from_bytes(data, "big")
        self._remaining = len(data) * 8

    def read(self, bit_count: int) -> int:
        """Read a fixed-width unsigned integer."""

        if bit_count <= 0:
            raise ValueError("bit_count must be positive")
        if bit_count > self._remaining:
            raise ValueError("not enough bits remaining")
        shift = self._remaining - bit_count
        mask = (1 << bit_count) - 1
        value = (self._value >> shift) & mask
        self._remaining -= bit_count
        self._value &= (1 << self._remaining) - 1 if self._remaining else 0
        return value


def encode_header(header: LXDRHeader) -> bytes:
    """Encode an ADRIAN header into a fixed packed representation."""

    header.validate()
    writer = _BitWriter()

    _encode_date(writer, header.date_request_created_local)
    _encode_time(writer, header.time_request_created_local)
    _encode_alphabetic(writer, header.physical_location_of_requestor, 14)
    _encode_alphabetic(
        writer,
        header.request_unique_identification_local,
        10,
    )
    writer.write(int(header.request_priority), 4)
    _encode_alphabetic(
        writer,
        header.element_unit_identification_callsign,
        4,
    )
    writer.write(header.request_segments, 4)

    return writer.to_bytes(PACKED_HEADER_SIZE_BYTES)


def decode_header(data: bytes) -> LXDRHeader:
    """Decode a fixed packed representation into an ADRIAN header."""

    if len(data) != PACKED_HEADER_SIZE_BYTES:
        raise ValueError(
            "Packed header must be exactly "
            f"{PACKED_HEADER_SIZE_BYTES} bytes"
        )

    reader = _BitReader(data)
    date_request_created_local = _decode_date(reader)
    time_request_created_local = _decode_time(reader)
    physical_location_of_requestor = _decode_alphabetic(reader, 14)
    request_unique_identification_local = _decode_alphabetic(reader, 10)
    request_priority = f"{reader.read(4):02d}"
    element_unit_identification_callsign = _decode_alphabetic(reader, 4)
    request_segments = reader.read(4)

    header = LXDRHeader(
        date_request_created_local=date_request_created_local,
        time_request_created_local=time_request_created_local,
        physical_location_of_requestor=physical_location_of_requestor,
        request_unique_identification_local=(
            request_unique_identification_local
        ),
        request_priority=request_priority,
        element_unit_identification_callsign=(
            element_unit_identification_callsign
        ),
        request_segments=request_segments,
    )
    header.validate()
    return header


def encode_pax_segment(segment: LXDRSegment) -> bytes:
    """Encode a PAX segment into a fixed packed representation."""

    if segment.spec is not PAX_MOVEMENT:
        raise ValueError("Packed PAX codec requires the PAX segment spec")

    values = segment.values
    writer = _BitWriter()

    writer.write(int(values["segment_number"]), 4)
    _encode_alphabetic(writer, values["request_type"], 2)
    writer.write(int(values["request_priority"]), 4)
    _encode_variable_alphabetic(writer, values["zap_or_edipi"], 10, 4)
    _encode_date(writer, values["earliest_departure_date"])
    _encode_date(writer, values["latest_departure_date"])
    _encode_alphabetic(writer, values["departure_location"], 14)
    _encode_alphabetic(writer, values["destination_location"], 14)
    writer.write(int(values["estimated_baggage_weight_lbs"]), 10)
    _encode_alphabetic(writer, values["hazardous_material_type"], 1)

    return writer.to_bytes(PACKED_PAX_SIZE_BYTES)


def decode_pax_segment(data: bytes) -> LXDRSegment:
    """Decode a fixed packed representation into a PAX segment."""

    if len(data) != PACKED_PAX_SIZE_BYTES:
        raise ValueError(
            "Packed PAX segment must be exactly "
            f"{PACKED_PAX_SIZE_BYTES} bytes"
        )

    reader = _BitReader(data)
    values = {
        "segment_number": str(reader.read(4)),
        "request_type": _decode_alphabetic(reader, 2),
        "request_priority": f"{reader.read(4):02d}",
        "zap_or_edipi": _decode_variable_alphabetic(reader, 10, 4),
        "earliest_departure_date": _decode_date(reader),
        "latest_departure_date": _decode_date(reader),
        "departure_location": _decode_alphabetic(reader, 14),
        "destination_location": _decode_alphabetic(reader, 14),
        "estimated_baggage_weight_lbs": f"{reader.read(10):03d}",
        "hazardous_material_type": _decode_alphabetic(reader, 1),
    }
    return LXDRSegment(spec=PAX_MOVEMENT, values=values)


def _encode_date(writer: _BitWriter, value: str) -> None:
    """Encode CCYYMMMDD as year, month, day fields."""

    if len(value) != 9:
        raise ValueError("date must be CCYYMMMDD")
    year = int(value[:4])
    month = _MONTHS[value[4:7]]
    day = int(value[7:9])

    writer.write(year, 14)
    writer.write(month, 4)
    writer.write(day, 5)


def _decode_date(reader: _BitReader) -> str:
    """Decode year, month, day into CCYYMMMDD."""

    year = reader.read(14)
    month = reader.read(4)
    day = reader.read(5)
    return f"{year:04d}{_MONTHS_REVERSE[month]}{day:02d}"


def _encode_time(writer: _BitWriter, value: str) -> None:
    """Encode HHMMSSDD as hour, minute, second, hundredths."""

    if len(value) != 8:
        raise ValueError("time must be HHMMSSDD")
    hours = int(value[0:2])
    minutes = int(value[2:4])
    seconds = int(value[4:6])
    hundredths = int(value[6:8])

    writer.write(hours, 5)
    writer.write(minutes, 6)
    writer.write(seconds, 6)
    writer.write(hundredths, 7)


def _decode_time(reader: _BitReader) -> str:
    """Decode hour, minute, second, hundredths into HHMMSSDD."""

    hours = reader.read(5)
    minutes = reader.read(6)
    seconds = reader.read(6)
    hundredths = reader.read(7)
    return f"{hours:02d}{minutes:02d}{seconds:02d}{hundredths:02d}"


def _encode_alphabetic(
    writer: _BitWriter,
    value: str,
    expected_length: int,
) -> None:
    """Encode a fixed-length ADRIAN uppercase alphanumeric field."""

    if len(value) != expected_length:
        raise ValueError(
            f"value must be {expected_length} characters long"
        )
    for char in value:
        try:
            encoded = _ALPHABET_INDEX[char]
        except KeyError as exc:
            raise ValueError(
                f"unsupported packed character: {char!r}"
            ) from exc
        writer.write(encoded, 6)


def _decode_alphabetic(reader: _BitReader, length: int) -> str:
    """Decode a fixed-length ADRIAN uppercase alphanumeric field."""

    chars = []
    for _ in range(length):
        chars.append(_ALPHABET[reader.read(6)])
    return "".join(chars)


def _encode_variable_alphabetic(
    writer: _BitWriter,
    value: str,
    max_length: int,
    length_bits: int,
) -> None:
    """Encode a bounded-length ADRIAN field with an explicit length."""

    if len(value) > max_length:
        raise ValueError(
            f"value exceeds maximum packed length of {max_length}"
        )
    writer.write(len(value), length_bits)
    padded = value.ljust(max_length, "0")
    _encode_alphabetic(writer, padded, max_length)


def _decode_variable_alphabetic(
    reader: _BitReader,
    max_length: int,
    length_bits: int,
) -> str:
    """Decode a bounded-length ADRIAN field with explicit length."""

    length = reader.read(length_bits)
    if length < 1 or length > max_length:
        raise ValueError("decoded variable-length field is out of bounds")
    value = _decode_alphabetic(reader, max_length)
    return value[:length]
