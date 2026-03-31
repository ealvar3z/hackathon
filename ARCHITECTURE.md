# section4 Architecture

## Scope

section4 is an ALOC-focused logistics workflow application that runs as
a local Python service and integrates with TAK clients through
`PyTAK`.

The source code under [`docs/`](./docs) is the implementation
guidebook for this project. In practice, that means:

- use `PyTAK` as the Python transport adapter for CoT traffic
- treat `sample-serverless-tak` as a later infrastructure reference,
  not as the MVP foundation
- borrow useful runtime and persistence patterns from prior reference
  projects without inheriting their transport assumptions
- keep section4's workflow state in its own database

## Deployment Topology

### MVP Demo Topology

- `ALOC node`
  - laptop runs section4
  - stores incidents, jobs, capabilities, artifacts, and events locally
  - hosts the operator-facing TUI
  - publishes and consumes CoT via `PyTAK`

- `EUD 1`
  - Android phone with ATAK
  - acts as forward requester

- `EUD 2`
  - Android phone with ATAK
  - acts as maintenance or fabrication-capable responder

### TAK Server Position

For the MVP:

- TAK Server is not required
- section4 remains the system of record
- `PyTAK` is used as an in-process library, not as a standalone service

For later phases:

- section4 can connect to TAK Server if we need standard TAK
  infrastructure, multi-client routing, or cloud deployment

## Architectural Layers

section4 should be implemented as three explicit layers.

### 1. Transport Layer

Responsibilities:

- initialise `PyTAK`
- manage CoT send and receive loops
- translate between CoT events and section4 domain objects
- isolate TAK-specific logic from the rest of the application

Reference guides:

- [classes.py](/Users/eax/repos/hackathon/docs/pytak/src/pytak/classes.py)
- [client_functions.py](/Users/eax/repos/hackathon/docs/pytak/src/pytak/client_functions.py)
- [send_receive.py](/Users/eax/repos/hackathon/docs/pytak/examples/send_receive.py)
- [cot-schema.md](/Users/eax/repos/hackathon/docs/cot-schema.md)

### 2. Domain Layer

Responsibilities:

- incident intake
- capability matching
- job creation and tracking
- course-of-action scoring
- readiness impact calculation
- ETA estimation
- event and audit recording

This is the part that is genuinely section4-specific.

### 3. Presentation Layer

Responsibilities:

- operator dashboard
- incident detail views
- capability views
- job board
- artifact visibility

The first operator surface is a local TUI built with `urwid`. The web
UI comes later. ATAK is not the ALOC workflow console.

## Persistence Strategy

section4 keeps its own persistence.

Store locally:

- incidents
- capabilities
- jobs
- artifacts metadata
- event and audit history
- course-of-action decisions

Use:

- `SQLite` for the MVP
- `SQLAlchemy` as the ORM and query layer

Why:

- TAK and CoT are event and interoperability layers, not section4's
  domain database
- TAK Server persistence is for TAK infrastructure, not for replacing
  application-level workflow state

## TAK Integration Strategy

Use `PyTAK` for:

- publishing incident markers and job updates
- ingesting field-originated events from ATAK
- translating section4 workflow state into TAK-visible operational
  events

Do not depend on TAK Server for:

- incident storage
- job state transitions
- artifact metadata
- course-of-action history

## FAQ

### Why section4 Is Not Built On TAK Server

The cloud solution in
[aws-serverless-tak](/Users/eax/repos/hackathon/docs/aws-serverless-tak)
is an infrastructure solution:

- ECS/Fargate deployment
- Aurora Serverless database for TAK Server
- EFS, S3, Route53, certificates, firewalling, and secrets

That is useful later, but it is not the fastest or safest path for a
hackathon MVP. section4 needs an application runtime and workflow
database first, not cloud TAK infrastructure first.

### Why section4 Still Needs Its Own UI

ATAK is useful for field awareness and operational event visibility,
but it is not sufficient as the primary ALOC workflow console for:

- reviewing incidents as records
- comparing candidate COAs
- managing assignments and queues
- tracking readiness impact
- inspecting artifact history

The ALOC node therefore needs a purpose-built TUI and later a web UI.

## Proposed Python Stack

- `pytak`
- `sqlalchemy`
- `pydantic`
- `urwid`
- `fastapi`
- `jinja2`
- `uvicorn`

Tooling:

- `uv` for dependency and task execution
- `ruff` for linting and formatting

## Proposed Repository Shape

```text
section4/
  pyproject.toml
  README.md
  ARCHITECTURE.md
  docs/
  src/section4/
    __init__.py
    app.py
    config.py
    transport/
    domain/
    storage/
    tui/
    web/
    services/
  tests/
  data/
    sqlite/
    artifacts/
```

## Initial Build Order

1. project metadata and architecture documents
2. package skeleton and repository layout
3. database schema and persistence layer
4. `PyTAK` runtime bootstrap
5. workflow services
6. section4 TUI
7. section4 web UI

## Implementation Rules

- default to patterns already proven in the reference code
- keep the runtime in one Python process unless there is a clear reason
  not to
- keep the domain model separate from transport and UI concerns
- treat `PyTAK` as a library inside section4, not as the application
  itself
- do not introduce TAK Server into the MVP without a concrete demo need
