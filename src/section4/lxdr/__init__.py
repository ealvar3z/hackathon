"""LXDR protocol package for section4."""

from section4.lxdr.header import LXDRHeader
from section4.lxdr.link import (
    LXDRLinkFrame,
    embed_bundle,
    embed_inline_structured,
    embed_inline_text,
)
from section4.lxdr.registry import APPENDIX_F_FIELDS, CanonicalFieldRegistry
from section4.lxdr.request import LXDRRequestContainer
from section4.lxdr.router import LXDRRouter, OutboxEntry
from section4.lxdr.router_store import PersistentLXDRRouter

__all__ = [
    "APPENDIX_F_FIELDS",
    "CanonicalFieldRegistry",
    "LXDRHeader",
    "LXDRLinkFrame",
    "LXDRRequestContainer",
    "PersistentLXDRRouter",
    "LXDRRouter",
    "OutboxEntry",
    "embed_bundle",
    "embed_inline_structured",
    "embed_inline_text",
]
