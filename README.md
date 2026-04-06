# section4

![section4 logo](assets/section4-small.png)

section4 is implementing `LXDR`, an ADRIAN-derived logistics exchange
protocol built strictly from
[`docs/project-adrian.txt`](./docs/project-adrian.txt).

The current project direction is deliberately narrow:

- implement `LXDR` fully from the ADRIAN whitepaper
- keep the protocol spec and code aligned
- do not build a custom network stack
- use `TAK` as the operational transport and network surface
- build an `lxmfcot` bridge, in the style of `aiscot`, `djicot`,
  `adsbcot`, and related TAK adapters
- use `PyTAK` to move `LXDR`-derived exchange data over CoT

## What This Repo Is

This repository is first and foremost an `LXDR` protocol implementation.

It currently contains:

- the evolving protocol draft in
  [`docs/lxdr-protocol.md`](./docs/lxdr-protocol.md)
- the ADRIAN source material in
  [`docs/project-adrian.txt`](./docs/project-adrian.txt)
- a protobuf schema for `LXDR` core and link objects in
  [`proto/lxdr/v1/lxdr.proto`](./proto/lxdr/v1/lxdr.proto)
- a Go implementation of:
    - `LXDR-Core`
    - `LXDR-Link`
    - `LXDR-Router v1`

This repository is not currently trying to be:

- a bespoke transport stack
- a TAK server
- a dashboard-first application

## Architecture Direction

The working architecture is:

1. `LXDR`
    - the ADRIAN protocol implementation
    - request/header/segment/schema/sync/router semantics

2. `lxmfcot`
    - a planned bridge process that converts between:
        - `LXDR`
        - `Cursor on Target`
    - built in the spirit of TAK adapter tools such as:
        - `aiscot`
        - `djicot`
        - `adsbcot`

3. `PyTAK`
    - the transport/client library used to publish and receive CoT over
      the TAK ecosystem

4. `TAK`
    - the operational transport and network substrate for demo and
      integration use

The key decision is that `TAK` is the network and transport environment.
`section4` will not spend time building a new bearer or routing stack
before the protocol itself is mature.

## Why This Direction

`project-adrian.txt` is the normative design source for this work.

The ADRIAN whitepaper argues for:

- critical logistics data capture at the point of need
- connected and disconnected synchronization
- homogeneous exchange across the logistics functions
- minimum critical data under constrained communications

Those requirements point to a protocol implementation effort first.

For practical integration and demo speed, `TAK` already gives us:

- fielded user interfaces
- familiar operational workflows
- established CoT exchange patterns
- a transport and networking environment we can ride immediately

So the project strategy is:

- finish `LXDR`
- bridge it into the TAK ecosystem
- demonstrate contested-logistics workflows there

## Planned Demo Direction

The intended demo path is:

1. edge users originate logistics-relevant information in TAK
2. `lxmfcot` receives CoT via `PyTAK`
3. `lxmfcot` converts that information into valid `LXDR`
4. `LXDR` request and synchronization logic runs locally
5. resulting state or synchronized outputs are emitted back into TAK

This keeps the main thing the main thing:

- the protocol is `LXDR`
- TAK is the transport and operator ecosystem
- `lxmfcot` is the bridge

## Current Protocol Status

The repository currently has a working `LXDR v1` baseline, including:

- protobuf-backed core schema
- generated Go types
- canonical text support where ADRIAN gives explicit examples:
    - request header
    - synchronized response
    - mobility PAX
    - mobility cargo
- validated Chapter 3 segment coverage for the implemented v1 families
- minimal `LXDR-Link` frame semantics
- formal synchronization exchange helpers
- local `LXDR-Router v1` state and workflow semantics

See:

- [`docs/lxdr-protocol.md`](./docs/lxdr-protocol.md)
- [`proto/lxdr/v1/lxdr.proto`](./proto/lxdr/v1/lxdr.proto)

## Canonical ADRIAN Workflows

These extracted figures from the ADRIAN source document remain the
conceptual guiderails for the implementation.

### HF Data Transmission

![ADRIAN HF to HF data transmission](assets/adrian-hf-to-hf-data-transmission.jpg)

### Data Synchronization Pathway

![ADRIAN data synchronization pathway](assets/adrian-data-sync-pathway.jpg)

### Mobility Interdependency

![ADRIAN mobility interdependency](assets/adrian-mobility-interdependency.jpg)

### Data to Decision Support

![ADRIAN data to wisdom model](assets/adrian-data-to-wisdom.jpg)

## Development

The implementation language is Go, with protobuf as the canonical typed
schema layer.

### Regenerate protobuf types

```bash
make proto
```

### Format Go sources

```bash
make fmt
```

### Run tests

```bash
make test
```

## Near-Term Work

Near-term work is expected to focus on:

- continuing to harden `LXDR`
- documenting and testing v1 protocol behavior
- building `lxmfcot`
- integrating with TAK through `PyTAK`

## Non-Goals Right Now

- building a custom network stack
- building a new TAK replacement
- inventing protocol fields not justified by
  [`docs/project-adrian.txt`](./docs/project-adrian.txt)
- drifting into UI-first work before the protocol and bridge are ready
