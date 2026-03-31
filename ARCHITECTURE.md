# section4 Architecture

## Scope

`section4` is an implementation effort for
[project-adrian.pdf](https://github.com/ealvar3z/hackathon/blob/main/docs/project-adrian.pdf).
The source document is treated as the normative design basis in the
same way an implementation would treat an RFC.

The project objective is not just to display logistics data. It is to
build the data capture, synchronization, and decision-support substrate
that ADRIAN argues is missing in contested, disconnected, degraded, and
intermittent environments.

The design center is the tactical logistics user operating at the point
of need with constrained communications, connected/disconnected
synchronization, and the need to move from data to supportability and
sustainment decisions.

## Architectural Intent

ADRIAN repeatedly implies four architectural requirements:

1. the system must remain homogeneous across the analyzed logistics
   functions
2. it must support connected and disconnected synchronization
3. it must preserve compact exchange over low-bandwidth pathways, with
   HF data-burst as the lowest denominator assumption
4. it must produce a logistics common operational picture that enables
   decision-making, not just status display

`section4` therefore separates:

1. `LXDR-Core`
2. `LXDR-Link`
3. `LXDR-Transport`
4. operator applications and decision-support tooling

This keeps the ADRIAN-derived data model stable while allowing multiple
delivery paths and operator interfaces.

## Demo Topology

The hackathon demo topology is intentionally simple:

- `ATAK` on Android EUDs
    - field-node interface
    - event and operational state visibility
- `s4net` on the laptop
    - local-first operator console
    - implemented first as a TUI
- `section4` runtime on the laptop
    - local persistence
    - protocol processing
    - synchronization logic
    - decision-support logic

This proves the key architectural point:

- the protocol core is transport and interface agnostic
- the field nodes can use one operational interface
- the laptop can use another
- the same ADRIAN-derived exchange model remains underneath both

## Protocol Stack

### LXDR-Core

`LXDR-Core` is the ADRIAN-derived payload model.

It defines:

- generated request header
- function-specific segments
- canonical field registry
- exchange metadata such as `Pass`, `Not Passed`, and `Sync Response`
- canonical text-burst representation
- packed representations for constrained transport

This layer is derived from:

- the request header sections and tables
- the function and sub-function schemas
- Appendix F canonical field/data block structure

Current implementation direction:

- semantic models for headers and requests
- segment registry
- canonical field registry
- canonical text-burst codec
- packed codecs for constrained exchange

### LXDR-Link

`LXDR-Link` is the transport-agnostic delivery and synchronization
envelope around one or more `LXDR-Core` objects.

It exists because ADRIAN requires:

- connected/disconnected synchronization
- continuity between local placeholder identity and synchronized
  enterprise identity
- compact exchange over constrained pathways
- store-and-forward behavior across interruptions

`LXDR-Link` is responsible for:

- message identity distinct from request identity
- sender and recipient addressing
- delivery method selection
- representation selection
- fragmentation and reassembly
- deduplication
- acknowledgement and retry
- synchronization bundles
- integrity and confidentiality metadata hooks

Current implemented representations:

- `INLINE_TEXT`
- `INLINE_STRUCTURED`
- `INLINE_PACKED`
- `BUNDLE`

### LXDR-Transport

`LXDR-Transport` is the bearer-adapter layer.

This layer moves opaque `LXDR-Link` frames over an available path. It
must not redefine ADRIAN request semantics.

Example transport paths include:

- HF data-burst
- local IP transport
- message relay
- export/import workflows
- bridges into existing operational ecosystems

The current demo still uses ATAK-facing integration to show state on
Android EUDs, but ATAK is not the architecture center. It is one
interface path over a transport-agnostic protocol core.

## Applications

### s4net

`s4net` is the laptop operator application.

Its job is to provide a local-first operations view over:

- requests
- capabilities
- jobs
- artifacts
- synchronization state
- event history
- decision-support outputs

The first implementation is a TUI because it is fast to iterate and
keeps the workflow close to the operator model. A web UI may follow
later, but it should mirror the same page, report, and workflow shapes.

### ATAK Demo Nodes

ATAK on Android EUDs remains part of the demo because it gives us:

- realistic field-node interaction
- a familiar operational interface
- proof that `section4` can project state into an existing ecosystem

ATAK is not the authoritative data model. `section4` owns the ADRIAN
request model and synchronization state.

## Data Model

The architecture follows ADRIAN's hierarchy:

- data field
- data block
- data file

In `section4`, that becomes:

- canonical field registry entries
- segment definitions built from those fields
- requests composed of headers plus segments
- link frames carrying one or more requests
- bundles supporting synchronized exchange

This hierarchy is critical because it lets the system remain consistent
across mobility, supply, maintenance, health, and the other logistics
functions.

## Persistence Strategy

`section4` keeps its own persistence.

The local store is the system of record for:

- requests
- capabilities
- jobs
- artifacts and references
- event and audit history
- synchronization state
- decision-support outputs

The current implementation uses:

- `SQLite`
- `SQLAlchemy`

Why:

- ADRIAN requires continuity across disconnected periods
- the operator console must work locally first
- external interfaces are transport or projection layers, not the
  authoritative workflow database

## Decision Support

The target outcome is a logistics COP that enables decisions.

That means the system should evolve from:

- data capture
- request exchange
- synchronization

into:

- supportability analysis
- candidate course-of-action comparison
- sustainment and routing decisions
- operator-facing reasoning support

One explicit hackathon objective is to build toward an ADRIAN-specialized
LLM workflow that can reason over:

- request content
- synchronization state
- supportability constraints
- response options
- logistics interdependencies

The LLM is not the protocol. It sits above the ADRIAN-derived data model
and uses it as structured input.

## Repository Shape

```text
section4/
  pyproject.toml
  README.md
  ARCHITECTURE.md
  docs/
    project-adrian.txt
    project-adrian.pdf
    lxdr-protocol.md
  assets/
  src/section4/
    app.py
    config.py
    lxdr/
    storage/
    transport/
    tui/
    services/
  tests/
  data/
    sqlite/
    artifacts/
```

## Current Build Order

1. `section4` project identity and local runtime
2. ADRIAN-derived `LXDR` protocol draft
3. `LXDR-Core` scaffolding
4. `LXDR-Link` serialization and parsing
5. packed codecs for constrained exchange
6. synchronization bundles
7. next: router, inbox/outbox, and DDIL message handling
8. later: richer operator workflows and ADRIAN-guided decision support

## Implementation Rules

- treat
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt)
  as normative
- do not invent ADRIAN payload fields where the source is silent
- separate ADRIAN payload semantics from link/session behavior
- separate link/session behavior from bearer adapters
- keep the local operator application functional while disconnected
- prefer compact, deterministic encodings for constrained paths
- use structured representations for development and interoperability,
  but do not confuse them with the long-term constrained wire format
