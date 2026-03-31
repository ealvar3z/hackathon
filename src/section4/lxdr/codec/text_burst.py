"""Canonical ADRIAN-style text burst rendering for LXDR core objects."""

from __future__ import annotations

from section4.lxdr.request import LXDRRequestContainer


def render_segment_text(
    segment_values: dict[str, str],
    field_order: list[str],
) -> str:
    """Render one segment using the ordered pass-field set."""

    return "-".join(segment_values[name] for name in field_order)


def render_request_container(container: LXDRRequestContainer) -> str:
    """Render a header plus segment lines as a canonical text burst."""

    container.validate()
    lines = [container.header.render_text()]
    for segment in container.segments:
        field_order = segment.pass_field_names()
        lines.append(render_segment_text(segment.values, field_order))
    return "\n".join(lines)
