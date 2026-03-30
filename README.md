# ForgeNet

ForgeNet is a Python application for contested logistics workflows built around a local ALOC node,
ATAK end-user devices, and TAK/CoT event exchange via `PyTAK`.

The current MVP topology is:

- `ALOC node`: Laptop running ForgeNet, local persistence, the ALOC TUI, and a `PyTAK`
  transport adapter
- `EUD 1`: Android device running ATAK as the forward requester
- `EUD 2`: Android device running ATAK as the maintenance or fabrication-capable responder

ForgeNet owns the business logic and state. `PyTAK` is used as a Python library to publish and
consume CoT events. `TAK Server` is optional infrastructure and is not required for the first demo.

## Why This Stack

ForgeNet follows the seam lines already proven in the reference source under
[`references/`](./references):

- `PyTAK` provides Python-side TAK/CoT send and receive primitives
- `sample-serverless-tak` shows how TAK Server fits as optional cloud infrastructure later

The ALOC application itself is intentionally standalone. It is not a TAK Server plugin and it is
not an ATAK replacement. It is the logistics workflow system that integrates with the TAK
ecosystem.

## Planned Capabilities

The first ForgeNet workflow is:

1. A forward node reports a failed component and low local stock
2. ForgeNet ingests the report and stores it locally
3. ForgeNet ranks candidate courses of action:
   - local repair
   - additive fabrication
   - reroute from another node
4. A task is assigned to a capable responder node
5. ForgeNet publishes operational state to ATAK clients via CoT
6. The ALOC UI shows readiness impact, ETA, task state, and audit history

## Project Status

This repository is being built in stages:

1. Project metadata and architecture documentation
2. Package skeleton and repository structure
3. Database schema and persistence layer
4. `PyTAK` runtime integration
5. ALOC TUI and workflow implementation

## Development

This project uses [`uv`](https://docs.astral.sh/uv/) for environment and dependency management.

### Install dependencies

```bash
uv sync
```

### Bootstrap the local database

```bash
uv run python -m forgenet.app bootstrap
```

### Run the app

The interactive TUI entrypoint will be added in the next stage. The current command shape is:

```bash
uv run python -m forgenet.app init-db
uv run python -m forgenet.app seed-demo
```

### Lint

```bash
uv run ruff check .
uv run ruff format .
```

## Initial Technical Decisions

- Language: Python 3.12
- Dependency management: `uv`
- Linting and formatting: `ruff`
- Web stack: FastAPI + Jinja2
- Persistence: SQLite + SQLAlchemy
- TAK integration: `PyTAK`
- Field client: ATAK
- Optional infrastructure: TAK Server later, not required for MVP

## Non-Goals For MVP

- No TAK Server dependency for the first demo
- No custom Android app
- No cloud-first deployment requirement
- No heavyweight SPA frontend unless the workflow proves it is necessary
