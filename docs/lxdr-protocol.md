# LXDR Protocol Draft

## Status

This document is a draft protocol extraction for `section4`.

It is derived directly from
[project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt).
The intent is to treat `project-adrian.txt` as the normative source in
the same way an implementation would treat an RFC.

This draft does not invent fields, segment order, or exchange behavior
that are not present in the source document.

## Scope

LXDR stands for `Logistics eXchange Data Requirements`.

The protocol objective is to encode logistics requests and related
updates for disconnected, degraded, intermittent, and low-bandwidth
environments. The source document repeatedly centers the forward
tactical user, limited communications pathways, and synchronization
after periods of disconnection. See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1768
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1770
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2048

The lowest transport denominator assumed by this draft is short,
transport-agnostic data burst transmission. The transport is explicitly
out of scope for the core protocol. LXDR defines payload structure and
sync semantics, not the bearer.

This does not mean delivery is out of scope for LXDR as a whole. It
means LXDR is layered:

- `LXDR-Core`
- `LXDR-Link`
- `LXDR-Transport`

Only `LXDR-Core` is fully specified by the first ADRIAN extraction. The
other two layers exist so the ADRIAN-derived data model does not become
coupled to one bearer or one network path.

## Normative Source Mapping

The initial LXDR draft is grounded in these ADRIAN sections and tables:

- Request header concept:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1768
- Header exchange requirements:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2014
- Header schema and synchronized response:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2023
- Mobility PAX:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2136
- Mobility cargo:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2236
- Supply request:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2692
- Maintenance request:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2917
- Engineer reconnaissance area data exchange:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):3377
- Health collection:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):5675
- Health treatment:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):5838
- Health evacuation:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):6290
- Appendix F data blocks:
  [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17830

## Protocol Model

### Layered Protocol Stack

LXDR should be treated as a layered protocol suite.

This structure is justified by ADRIAN in three ways:

1. the document explicitly requires connected and disconnected
   synchronization, which implies a protocol identity beyond any one
   momentary bearer
2. the document frames logistics exchange in kill-chain terms, meaning
   the right data must get from the right sensor to the right
   logistician even in denied or degraded communications
3. the document distinguishes data field, data block, and data file
   layers, which supports separating payload semantics from carriage

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1168
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1245
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1291

The layered stack is:

1. `LXDR-Core`
2. `LXDR-Link`
3. `LXDR-Transport`

### LXDR-Core

`LXDR-Core` is the ADRIAN-derived application payload model.

It defines:

- request header
- segment registry
- canonical field registry
- exchange metadata such as `Pass`, `Not Passed`, and `Sync Response`
- synchronization identity semantics
- canonical text-burst serialization

`LXDR-Core` is where the ADRIAN request structure lives.

### LXDR-Link

`LXDR-Link` is the transport-agnostic delivery and synchronization
envelope that carries `LXDR-Core` objects between peers or systems.

This layer is required even though ADRIAN does not name it directly,
because ADRIAN repeatedly requires:

- connected/disconnected synchronization
- acknowledgement and enterprise replacement of local identifiers
- support for constrained communication pathways
- continuity of a record across exchanges

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1168
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1768
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1810

The first draft assigns the following responsibilities to `LXDR-Link`:

- peer envelope around one or more core objects
- message identity separate from request identity
- sender and recipient addressing fields
- correlation of sync responses to local placeholder IDs
- fragmentation and reassembly for constrained bearers
- deduplication
- acknowledgement and retry semantics
- integrity and confidentiality metadata
- bundle packaging for store-and-forward exchange

This layer is conceptually similar to the delivery/session problem that
transport-agnostic messaging systems solve, but it remains a clean-room
design for `section4`.

### LXDR-Link Envelope

The first draft defines `LXDR-Link` as the minimum envelope required to
move one or more `LXDR-Core` objects through DDIL pathways without
changing their meaning.

An `LXDR-Link` frame should contain:

1. `link_message_id`
2. `sender_id`
3. `recipient_id`
4. `created_at_local`
5. `delivery_method`
6. `representation`
7. `payload_count`
8. `payload_refs` or embedded payloads
9. `fragment_index`
10. `fragment_count`
11. `correlation_id`
12. `ack_of`
13. `sync_of`
14. `attempt_count`
15. `integrity_metadata`
16. `confidentiality_metadata`

The purpose of each field is:

- `link_message_id`
  - identity of the delivery object
  - distinct from the ADRIAN request identifier
- `sender_id`
  - identity of the transmitting system, node, or operator endpoint
- `recipient_id`
  - intended peer, group, or service endpoint
- `created_at_local`
  - local time of link frame creation
- `delivery_method`
  - selected carriage strategy
- `representation`
  - how the core payload is encoded for transfer
- `payload_count`
  - count of carried core objects
- `payload_refs`
  - identifies embedded or externally referenced core objects
- `fragment_index` and `fragment_count`
  - support for constrained-bearer fragmentation
- `correlation_id`
  - allows related frames to be grouped
- `ack_of`
  - identifies the prior link frame being acknowledged
- `sync_of`
  - identifies the request or bundle being synchronized
- `attempt_count`
  - tracks retransmission attempts
- `integrity_metadata`
  - indicates integrity mechanism in use
- `confidentiality_metadata`
  - indicates confidentiality mechanism in use

This draft names these fields because a DDIL protocol needs them, not
because ADRIAN lists them directly. ADRIAN gives the requirement for
connected/disconnected synchronization, constrained exchange, and
request identity continuity; `LXDR-Link` is the clean-room mechanism for
meeting that requirement.

### LXDR-Link Delivery States

The first draft defines the following delivery lifecycle states for a
link frame:

- `CREATED`
- `QUEUED`
- `SENDING`
- `SENT`
- `ACKNOWLEDGED`
- `DELIVERED`
- `SYNCED`
- `REJECTED`
- `FAILED`
- `EXPIRED`
- `CANCELLED`

These states exist at the link layer only.

They do not alter the semantic state of the carried `LXDR-Core`
request. For example:

- a request may be semantically open while its latest link frame is
  `FAILED`
- a request may be semantically active while a later sync frame is
  `SYNCED`

This separation is necessary so delivery outcomes do not overwrite the
ADRIAN-derived request lifecycle.

### LXDR-Link Representations

The first draft defines three conceptual representations:

1. `INLINE_TEXT`
2. `INLINE_STRUCTURED`
3. `INLINE_PACKED`
4. `BUNDLE`

`INLINE_TEXT`

- carries the canonical ADRIAN-style text burst directly
- intended for the most constrained exchanges
- must preserve exact field ordering for all `Pass` fields

`INLINE_STRUCTURED`

- carries the same semantics in a structured format
- suitable for local queues, richer bearers, or internal relay
- must round-trip to the canonical text burst without semantic loss

`INLINE_PACKED`

- carries a compact fixed-allocation packed representation
- intended for constrained DDIL transports where the canonical text
  burst is still too costly
- must remain fully reversible back to the same `LXDR-Core` object
- must be specified field-by-field from ADRIAN lengths and code spaces,
  not by ad hoc compression

`BUNDLE`

- carries one or more core objects plus synchronization metadata
- intended for deferred delivery, replay, or bulk synchronization after
  disconnection
- in the first `section4` implementation, a bundle payload is a
  structured object containing:
  - `bundle_type`
  - `request_count`
  - `requests`

The required invariant is:

- every representation must preserve the meaning of the same
  `LXDR-Core` object

### LXDR-Link Delivery Methods

The first draft defines conceptual delivery methods:

- `OPPORTUNISTIC`
- `DIRECT`
- `RELAYED`
- `SYNCHRONIZATION`
- `EXPORT`

`OPPORTUNISTIC`

- send immediately over whatever pathway is presently available
- best for short one-shot burst delivery

`DIRECT`

- sender has a current path to the intended recipient

`RELAYED`

- frame is forwarded or held by an intermediate node or service

`SYNCHRONIZATION`

- frame is part of a reconciliation or catch-up exchange rather than an
  initial report

`EXPORT`

- frame is prepared for out-of-band transfer
- examples include file export, removable media, or later manual import

These method names are draft terminology for `section4`. They are not
claimed as ADRIAN terms.

### LXDR-Link Fragmentation

Because ADRIAN is optimized for constrained and intermittent exchange,
the link layer must support fragmentation without changing the carried
core object.

The first draft requires:

- fragments carry a shared `link_message_id`
- each fragment carries `fragment_index`
- each fragment carries `fragment_count`
- reassembly occurs before core decoding
- duplicate fragments are ignored after successful reassembly

Fragmentation rules for any specific bearer belong to
`LXDR-Transport`. The logical fragment model belongs to `LXDR-Link`.

### LXDR-Link Acknowledgement and Synchronization

The first draft distinguishes:

- `delivery acknowledgement`
- `synchronization acknowledgement`

Delivery acknowledgement means:

- a link frame or all fragments of a link frame were received intact

Synchronization acknowledgement means:

- the receiving side accepted the carried request container into its
  synchronized record space
- and, when relevant, returned the enterprise replacement identifier

This distinction is required by ADRIAN because simple receipt is not the
same as request synchronization.

### First Bundle Form

The first implemented `BUNDLE` form in `section4` is a structured sync
container.

It carries:

- one `LXDR-Link` frame with representation `BUNDLE`
- one embedded payload object
- that payload object contains:
  - `bundle_type`
  - `request_count`
  - `requests`

The `requests` value is an ordered list of structured `LXDR-Core`
request containers.

This is intentionally conservative:

- it gives `section4` a real store-and-forward unit now
- it avoids inventing binary bundle framing before the sync model is
  more mature
- it keeps bundle semantics distinct from `INLINE_STRUCTURED`

The current implementation rule is:

- `INLINE_STRUCTURED` carries individual structured requests
- `BUNDLE` carries a structured container of multiple requests intended
  for synchronization or deferred transfer

### LXDR Router

The first draft defines an `LXDR Router` as the runtime component that
manages `LXDR-Link` frames.

The router should own:

- outbound queueing
- inbound processing
- dedupe cache
- retry scheduling
- transport selection
- fragmentation and reassembly
- synchronization bundle generation
- sync response matching
- message expiry and cleanup

The router should not own:

- ADRIAN field definitions
- segment schemas
- request business logic

Those remain in `LXDR-Core`.

### LXDR Router Minimum Behaviors

A conforming first-draft router should be able to:

1. enqueue a new core request for delivery
2. wrap it in a link envelope
3. select a transport adapter
4. fragment if required by the chosen transport
5. retry on temporary delivery failure
6. suppress duplicates on receive
7. reassemble fragments into a complete link frame
8. decode the carried core payload
9. correlate sync responses to local placeholder IDs
10. produce a bundle of unsynchronized requests for later transfer

### LXDR-Link Security Hooks

This draft does not yet define the actual integrity or confidentiality
algorithms, but it reserves them at the link layer.

The design rule is:

- `LXDR-Core` should remain readable by the application
- `LXDR-Link` carries the metadata and processing hooks for securing the
  transfer object
- `LXDR-Transport` must not redefine security semantics independently

The actual cryptographic format requires a later draft.

### LXDR-Transport

`LXDR-Transport` is the bearer adapter layer.

It is responsible only for moving opaque `LXDR-Link` frames over a
specific pathway. Examples include:

- HF data burst
- unicast or multicast IP
- local file transfer
- TAK bridge carriage

The first rule of `LXDR-Transport` is that it must not change the
meaning of `LXDR-Core`. A transport adapter may:

- fragment
- encode for channel constraints
- schedule retransmission
- prioritize traffic

but it must preserve the semantics of the carried link frame.

### Request Container

The fundamental LXDR object is a `request container`.

Per ADRIAN, the request header:

- is generated at request creation time
- has little to no user interaction
- precedes the actual request information
- initiates the demand process
- provides common context for follow-on information gathering, status,
  and response
- acts as a container for follow-on child or sub-requests

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1768
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1771
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1935

An LXDR request container therefore consists of:

1. one `header`
2. one to nine `segments`

The segment count is part of the header and is used to verify that the
full data file has been passed and corresponds to child requests within
the request container. See
[project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1927
and
[project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1930.

In layered terms:

- the request container is a `LXDR-Core` object
- that core object is wrapped by `LXDR-Link`
- the resulting link frame is carried by `LXDR-Transport`

### Segment Semantics

Each segment is function-specific and attaches to the request header.
Header data is not duplicated inside segments. See
[project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2048
and
[project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2049.

The source document defines segments using:

- segment name
- maximum use
- repeat min/max
- usage
- exchange behavior
- schema examples in some cases

LXDR must preserve these metadata as first-class protocol properties.

## Header

### Header Data Elements

The LXDR header is derived from ADRIAN Section `3.3.3.1` and Table 7.

Header fields in order:

1. `date_request_created_local`
2. `time_request_created_local`
3. `physical_location_of_requestor`
4. `request_unique_identification_local`
5. `request_priority`
6. `element_unit_identification_callsign`
7. `request_segments`

Additionally defined by ADRIAN but not passed in the initial exchange:

- `time_request_created_utc`
- `date_time_group_utc`
- `request_unique_identification_sync`

The exchange rules are explicit in Table 7:

- `Pass`:
  - local date
  - local time
  - physical location
  - local unique identifier
  - request priority
  - element/unit/callsign
  - request segments
- `Not Passed`:
  - UTC time
  - DTG
- `Sync Response`:
  - synchronized unique identifier

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2014
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2023

### Header Field Constraints

From ADRIAN:

- local date format: `CCYYMMMDD`
- local time format for header requires decimal seconds to prevent
  duplicate timestamps
- physical location is the last synchronized MGRS reference and is 14
  characters
- local placeholder key is 10 alphanumeric characters
- synchronized enterprise request tracking ID is 12 alphanumeric
  characters
- request priority is 2 digits
- element/unit/callsign is a 4-character alphanumeric reference key
- request segments is a single digit from `1` to `9`

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1778
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1801
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1810
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1865
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2014

### Header Text Form

Initial burst format:

```text
CCYYMMMDD-HHMMSSDD-XXXX123456-XXXXXXXXXX-00-XXXX-00
```

ADRIAN example:

```text
2027OCT13-15470352-4QFJ123456-3838JBNM5-02-KL9K-01
```

Synchronized response format:

```text
3838JBNM5-KL9K15474QFJ
```

See Table 8 and Table 9:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2023
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2033

### LXDR Header Rules

From ADRIAN, the following rules are normative for LXDR:

- header fields are system-generated, not user-entered
- segment priority rolls up into header priority, with highest priority
  taking precedence
- synchronized request ID replaces the local placeholder after
  synchronization
- header text is compact and suitable for limited-bandwidth exchange

## Segment Registry

The initial LXDR segment registry is derived from the ADRIAN
function-specific tables.

### Mobility: PAX Movement

Source:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2136
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2319

Segment code:

- `PM`

Segment fields in order:

1. segment number
2. request type
3. request priority
4. ZAP/EDI-PI
5. earliest departure date
6. latest departure date
7. departure location
8. destination location
9. total estimated baggage weight
10. hazardous material type

Calculated but not passed:

11. person count
12. total baggage weight

Text schema:

```text
0-XX-00-XXXXXXXXXX-CCYYMMMDD-CCYYMMMDD-XXXX123456-XXXX123456-000-X
```

Example:

```text
1-PM-02-1010919789-2027OCT15-2027OCT20-4QFJ123456-4QFJ456789-075-X
```

### Mobility: Cargo Movement

Source:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2236
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2523

Segment code:

- `CM`

Segment fields in order:

1. segment number
2. request type
3. request priority
4. item by NIIN
5. item quantity
6. serial number
7. gross weight
8. actual height
9. actual width
10. actual length
11. HMIC
12. handling
13. earliest departure date
14. latest departure date
15. departure location
16. destination location

Calculated but not passed:

17. item count
18. total weight

Text schema:

```text
0-XX-00-000000000-00000-XXXXXXXXXX-00000-000-000-000-X-X-CCYYMMMDD-CCYYMMMDD-XXXX123456-XXXX123456
```

Note:

The ADRIAN prose says `"CM is the code for a PAX movement"` in the
cargo section. This appears to be a document wording error because the
surrounding section is explicitly `Cargo Request Data`. LXDR should
preserve the code `CM` but record this wording inconsistency as source
ambiguity rather than reinterpret the table silently.

### Supply Request

Source:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2692
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2811

Segment code:

- `SR`

Segment fields in order:

1. segment number
2. request type
3. request priority
4. item by NIIN
5. item quantity
6. required delivery date
7. delivery location
8. attachment
9. narrative

Design constraints taken directly from ADRIAN:

- the supply request is class-agnostic
- one form is used across supply classes and subclasses
- NIIN may be absent initially if the item is not known
- supporting activity may enrich the request later

### Maintenance Request

Source:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2917
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):3164

Segment fields in order:

1. segment number
2. request type
3. request priority
4. serial number
5. NIIN
6. model of equipment
7. item nomenclature
8. number of pieces
9. equipment operational condition
10. date maintenance support is required
11. location of equipment
12. type of maintenance support
13. type of repair
14. repair major defect
15. attachment
16. narrative

Maintenance support codes explicitly listed:

- `R1` repair
- `R2` recovery retrieve
- `R3` recovery free from immobility
- `R4` contact team
- `XX` none/default

Repair type codes explicitly listed:

- `M1` modification
- `S1` servicing preventive maintenance
- `S2` servicing tune/adjust
- `C1` calibration
- `D1` repair defect

Major defect codes explicitly listed:

- `MD01` through `MD16`
- `NMAJ`

### Maintenance Ambiguity

ADRIAN states:

> In this case CM is the code for data exchange.

inside the maintenance request section. See
[project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2933.

This draft treats that as an unresolved source ambiguity because:

- `CM` is already used in the cargo movement section
- the maintenance section title and field set are otherwise distinct
- Table 18 does not provide a separate schema example in the extracted
  text establishing a different code

Therefore the first LXDR draft records:

- maintenance has a normative segment structure from Table 18
- the maintenance request type code remains unresolved in this draft
  pending a stricter source reconciliation

### General Engineering: Area Reconnaissance

Source:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):3377

The extracted source confirms the existence of the engineer
reconnaissance area report data exchange and its position as Table 19.
The nearby section also shows the same pattern continuing into zone
reconnaissance. The first LXDR draft records general engineering as a
supported function family, but does not yet define all engineering
sub-segment codecs in this draft because the full field tables have not
all been transcribed here.

### Health Services: Collection

Source:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):5675
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):5828

Segment code:

- `CR`

Segment fields in order:

1. segment number
2. request type
3. request priority
4. ZAP/EDI-PI
5. last name
6. first name
7. service
8. element/unit identification/callsign
9. allergies
10. date of injury
11. time of injury
12. location injury occurred

### Health Services: Treatment / Intervention

Source:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):5838
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):6244

The treatment segment continues the record initiated at collection and
is append-oriented. The extracted table confirms at least the
intervention portion fields `27..37` as transmitted data elements:

- extremity bleeding treatment
- wounds treatment
- airway treatment
- breathing treatment
- fluid circulation treatments
- blood circulation treatments
- analgesic medication treatments
- antibiotic medication treatments
- other medication treatments
- casualty type
- first responder ZAP/EDI-PI

ADRIAN also explicitly describes the earlier treatment fields in prose,
including:

- primary mechanism of injury
- CBRN related/exposure
- major signs/symptoms
- injury locations
- vitals date/time and vital signs

LXDR therefore treats health treatment as a continuation segment whose
elements are appended to the existing casualty record.

### Health Services: Hold

Source:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):6248

The hold phase is append-only:

- no data is overwritten
- no data is deleted
- additional date/time-driven entries are appended

This is a normative lifecycle rule for LXDR health records.

### Health Services: Evacuation

Source:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):6290
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):6377

The evacuation request is described as a specialized PAX movement
request. It prioritizes:

1. evacuation request
2. priority medical information
3. medical record / casualty report

Segment fields in order:

38. request priority
39. location of pickup
40. location marking
41. location contamination
42. contact settings
43. count of casualties by precedence
44. count of casualties by type
45. requested equipment
46. security
47. ZAP/EDI-PI of casualties/treatments

## Canonical Field Registry

Appendix F is treated as the basis for a canonical LXDR field registry.
The source document organizes data as:

- category
- data block
- data element
- data field

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1313
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17830

The initial LXDR rule is:

- segment definitions reference canonical field definitions
- canonical field definitions are grouped into data blocks
- a future implementation should encode field metadata such as:
  - display label
  - fixed length or range
  - exchange behavior
  - code list, if any

This is required to avoid ad hoc per-segment redefinition of the same
underlying logistics fields.

### Canonical Registry Format

The initial canonical field registry format is derived directly from the
ADRIAN data taxonomy and Appendix F structure.

ADRIAN defines:

- `Data Field (Element)` as the smallest container for a recorded fact,
  character, or symbol
- `Data Block (Object)` as a logical, topic-based grouping of data
  fields
- `Data File` as a collection of data blocks

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1307
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1313

Appendix F then presents the reusable registry in four columns:

1. `Category`
2. `Data Block`
3. `Data Element`
4. `Data Field`

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17830
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17832
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17929
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17932

The first-draft LXDR canonical registry should therefore use this
logical record shape:

```text
category
data_block
data_element
data_field
```

This is the minimum normative shape. A registry entry in LXDR should be
traceable back to an Appendix F row or to another ADRIAN-defined field
source.

### Registry Levels

The first-draft registry levels are:

#### 1. Category

The broadest organizational grouping from Appendix F.

Example extracted from Appendix F:

- `Department of Defense Activity Address Directory`

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17834

#### 2. Data Block

The topic-based object grouping within a category.

Examples extracted from Appendix F:

- `Activity`
- `Type Address Code 1 (Postal/Mail)`
- `Type Address Code 2 (Freight)`

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17834
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17932

#### 3. Data Element

The sub-grouping inside the data block.

Examples extracted from Appendix F:

- `Header`
- `COMMRI`

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17929

#### 4. Data Field

The concrete reusable field definition.

Examples extracted from Appendix F:

- `1. Department of Defense Activity Address Code`
- `2. Unit Identifier Code`
- `1. Bill`
- `2. Data Pattern`
- `1. T1_Line 1`
- `6. T1_City`

See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):17932

### Registry Entry Shape

The first-draft canonical registry entry should include:

1. `registry_id`
2. `category`
3. `data_block`
4. `data_element`
5. `data_field`
6. `field_label`
7. `field_order`
8. `source_reference`

Only items `2` through `5` are directly prescribed by ADRIAN Appendix
F. The additional fields are implementation scaffolding required to make
the registry executable while remaining source-traceable.

The meaning of the added fields is:

- `registry_id`
  - stable internal identifier for the registry entry
- `field_label`
  - normalized field display label
- `field_order`
  - preserves the numeric ordering embedded in many Appendix F field
    names
- `source_reference`
  - exact source provenance, such as section or appendix row family

### Registry Identifier Rule

The first-draft identifier rule is:

```text
<category>::<data_block>::<data_element>::<data_field>
```

This rule is draft implementation guidance, not an ADRIAN term. Its
purpose is to preserve the four-level Appendix F structure without
collapsing it into ad hoc local names.

### Registry Normalization Rule

Appendix F contains ordered field labels such as:

- `1. Department of Defense Activity Address Code`
- `2. Unit Identifier Code`
- `1. T1_Line 1`

The first-draft normalization rule is:

- preserve the original ADRIAN label exactly in `data_field`
- extract the numeric prefix into `field_order` when present
- derive a normalized implementation label separately in
  `field_label`

This prevents losing the source ordering while still allowing cleaner
internal code references.

### Registry-to-Segment Rule

Segments should not define fields from scratch when a canonical field
already exists.

Instead:

- a segment field should reference a canonical registry entry where
  possible
- a segment may add segment-specific exchange metadata such as:
  - `Must Use`
  - `Optional`
  - `Pass`
  - `Not Passed`
  - repeat min/max
- segment metadata augments a field reference; it does not replace the
  canonical field definition

This follows ADRIAN’s data model:

- reusable fields form blocks
- blocks form files
- requests are built from those structures

### Registry Families

The first draft recognizes that the canonical registry will contain at
least two families:

1. `Appendix F canonical fields`
2. `ADRIAN request/segment fields not yet visible in Appendix F`

This distinction is necessary because:

- Appendix F is clearly intended as a reusable data block catalog
- some request table fields are operational exchange fields whose full
  Appendix F mapping may require additional extraction work

The implementation rule is:

- when Appendix F has the field, use Appendix F as canonical source
- otherwise create a provisional registry entry marked with direct
  section/table provenance from the request table

### Example Registry Record

The following is a format example, not a claim of full Appendix F
coverage:

```text
registry_id:
  Department of Defense Activity Address Directory::Activity::Header::1. Department of Defense Activity Address Code
category:
  Department of Defense Activity Address Directory
data_block:
  Activity
data_element:
  Header
data_field:
  1. Department of Defense Activity Address Code
field_label:
  department_of_defense_activity_address_code
field_order:
  1
source_reference:
  Appendix F
```

This example is derived from the Appendix F rows visible in the source
text and demonstrates how the ADRIAN hierarchy becomes an executable
registry without changing the underlying structure.

## Exchange Semantics

ADRIAN uses explicit exchange metadata in the segment tables. LXDR must
preserve at least these properties for each field:

- `usage`
  - `Must Use`
  - `Optional`
- `exchange`
  - `Pass`
  - `Not Passed`
  - `Sync Response`
- repeatability
  - `No`
  - `Yes`
- min/max length

Fields marked `Not Passed` are part of the local application state and
derived record, but are not part of the over-the-wire burst payload.

## Synchronization

Synchronization is built into the protocol model.

The header’s local placeholder key exists specifically until the
request is connected, acknowledged, and synchronized to the enterprise
system. The synchronized response replaces the local identifier with the
enterprise tracking identifier. See:

- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):1810
- [project-adrian.txt](/Users/eax/repos/hackathon/docs/project-adrian.txt):2030

First-draft LXDR synchronization rules:

1. a request originates with a local identifier
2. a request may remain unsynchronized for some period
3. synchronization returns an enterprise identifier
4. the enterprise identifier replaces the local placeholder in the
   canonical record
5. segment count is used to validate the received container

`LXDR-Core` defines the request identity semantics.

`LXDR-Link` is expected to define the mechanics by which:

- sync acknowledgements are correlated
- out-of-order delivery is tolerated
- duplicate link frames are suppressed
- bundles of unsynchronized requests are replayed after disconnection

## Lowest-Denominator Text Burst

The ADRIAN schema examples are the normative basis for a lowest-common-
denominator text codec.

The first LXDR draft therefore defines a required canonical text
representation:

- one header text line
- one or more segment text lines
- field ordering exactly matches ADRIAN table order
- only fields marked `Pass` appear in the transmitted text burst

This is the representation intended for the most constrained transports.
Higher-capability transports may carry richer envelopes, but they must
round-trip to the canonical text burst without semantic loss.

## Packed Core Representation

The first draft also permits a packed `LXDR-Core` representation for
constrained transports, as long as it remains reversible to the same
canonical request object.

The implementation rule is:

- `INLINE_TEXT` remains the canonical human-readable burst
- `INLINE_PACKED` is an implementation-specific compact encoding
- `INLINE_PACKED` must decode to the same header and segment values that
  would appear in `INLINE_TEXT`

### Packed Header Allocation

The first implemented packed allocation in `section4` is the ADRIAN
header.

Current allocation:

- local date
  - year: 14 bits
  - month: 4 bits
  - day: 5 bits
- local time
  - hour: 5 bits
  - minute: 6 bits
  - second: 6 bits
  - hundredths: 7 bits
- physical location of requestor
  - 14 characters
  - 6 bits per character
- local request identifier
  - 10 characters
  - 6 bits per character
- request priority
  - 4 bits
- element/unit identification/callsign
  - 4 characters
  - 6 bits per character
- request segments
  - 4 bits

Current total:

- 224 bits
- 28 bytes

This is an implementation-specific first allocation for `section4`. It
is ADRIAN-derived, but ADRIAN itself does not prescribe this exact bit
layout.

### Packed Representation Rules

The first-draft packed rules are:

1. packing must preserve ADRIAN field order and field meaning
2. fields marked `Not Passed` stay out of the packed wire form
3. variable-length coded fields must carry enough metadata to decode
   back to their original value
4. packed forms are not allowed to invent values missing from the core
   object
5. packed and text forms must round-trip to the same canonical request
   container

## Non-Goals of This Draft

This draft does not yet define:

- encryption format
- signature format
- link frame format
- fragmentation framing across specific bearers
- retransmission timing
- peer addressing model
- engineering sub-report codecs beyond the first identified report
- a resolved maintenance request type code beyond the source ambiguity

Those items require either more ADRIAN extraction or a separate design
document clearly marked as implementation-specific rather than
source-derived.

## Initial Implementation Consequences

A clean implementation based on this draft should separate:

- `header` model
- `segment registry`
- `canonical field registry`
- `link envelope`
- `text burst codec`
- `sync response handling`
- transport adapters

The implementation should never bake transport assumptions into the core
LXDR schema.
