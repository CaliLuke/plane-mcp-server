"""Short UUID encoding/decoding for compact MCP responses.

UUIDs are encoded to ~22-char base57 strings in all outgoing responses.
Incoming parameters accept both short and full UUID formats transparently.
"""

import re
import uuid as _uuid_mod
from typing import Annotated

import shortuuid as _su
from pydantic import BeforeValidator

_UUID_RE = re.compile(
    r"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", re.I
)


def encode(uuid_str: str | None) -> str | None:
    """Encode a full UUID string to a short base57 string."""
    if not uuid_str:
        return uuid_str
    return _su.encode(_uuid_mod.UUID(uuid_str))


def decode(s: str) -> str:
    """Accept a short or full UUID string; always return a full UUID string."""
    if _UUID_RE.match(s):
        return s
    return str(_su.decode(s))


def encode_value(v: object) -> object:
    """Recursively encode any UUID strings found in dicts/lists/scalars."""
    if isinstance(v, dict):
        return {k: encode_value(val) for k, val in v.items()}
    if isinstance(v, list):
        return [encode_value(item) for item in v]
    if isinstance(v, str) and _UUID_RE.match(v):
        return encode(v)
    return v


# Annotated type for tool parameters — FastMCP/Pydantic decodes transparently.
ShortUUID = Annotated[str, BeforeValidator(decode)]
