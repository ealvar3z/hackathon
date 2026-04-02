# LXDR Protocol Draft

## 1. Objective

The analysis objective in Project ADRIAN is to identify the minimum data
exchange requirements at the lowest tactical level so that logistics
decision and action can begin.

For protocol purposes, this means LXDR MUST:

- represent the minimum data required to initiate logistics action
- preserve request identity across disconnected and synchronized states
- support the five doctrinal functions analyzed in Project ADRIAN:
    - mobility
    - supply
    - maintenance
    - general engineering
    - health services
- allow the tactical user to originate data once and move it upward
  through the support chain

## 2. Scope

This draft is limited to the scope explicitly stated in Section 2.3 and
Section 3:

- lowest tactical level to higher supporting echelons
- five of the six doctrinal functions of Marine Corps logistics
- minimum essential data elements for request initiation and exchange

This draft is not yet a full logistics application architecture. It is a
request exchange protocol definition.

This draft does not define:

- aviation logistics
- services as the sixth doctrinal function
- higher-order optimization or forecasting logic
- complete enterprise logistics workflows after initial demand exchange

## 3. Problem Statement

Project ADRIAN identifies a current logistics IT environment that is:

- fragmented across many stand-alone systems
- heavily dependent on human transcription
- largely analog at the point of origin
- not federated
- not holistically integrated
- not suited to contested or degraded communications

This leads directly to the protocol requirement that LXDR be:

- compact
- structured
- hierarchical
- synchronization-aware
- suitable for digital creation at the point of demand

## 4. Foundational Concepts

### 4.1 Logistics as a Super System

Project ADRIAN defines logistics as a super system composed of functions,
sub-functions, processes, procedures, and steps.

For protocol design, this means:

- the protocol MUST model logistics by function and sub-function
- a request is not a free-form message
- request meaning is bound to doctrinal function and process context

### 4.2 Capability

Project ADRIAN defines a capability as something aligned to:

- a logistics function
- an enabling factor
- an echelon

For protocol design, this means each request type SHOULD be traceable to:

- a function
- a sub-function
- the level at which action is expected

### 4.3 Kill Chain Application to Logistics

Project ADRIAN explicitly applies the kill-chain concept to logistics.

For protocol design, this means:

- the right data must reach the right logistician at the right time
- unnecessary traffic is harmful
- data selection and exchange size are not implementation details; they are
  core protocol constraints

### 4.4 Data Is Not Information

Project ADRIAN distinguishes:

- data field
- data block
- data file

This hierarchy is protocol-significant.

#### Data Field

A single recorded fact, character, or symbol.

#### Data Block

A logical grouping of related data fields.

#### Data File

A collection of data blocks on which higher-level operations are
performed.

For protocol design, this means:

- LXDR MUST preserve a field-to-block-to-file hierarchy
- individual request segments SHOULD be understood as data blocks
- a request container is a data file composed of header and request
  segments

## 5. Knowledge Management and Protocol Implications

Sections 3.2.7.x state that useful logistics decision-making depends on:

- information accuracy
- information accessibility
- information type
- information sharing

For protocol design, this means:

- values MUST be precise enough for action
- data MUST be available at the point of need and exchangeable upward
- only relevant data should be transmitted
- the same request representation SHOULD be usable across organizations

The protocol is therefore a knowledge-enabling mechanism, not merely a
record format.

## 6. Communications Assumptions

### 6.1 Voice Is the Last Resort

Project ADRIAN states that legacy voice-over-radio reporting:

- increases EM exposure
- takes longer to transmit
- increases human transcription error
- relies on paper capture

Therefore, a digital protocol representation MUST replace voice-first
request transmission wherever possible.

### 6.2 HF Data Burst Is the Lower Bound

Project ADRIAN identifies HF radio data transmission as the relevant lower
bound:

- bandwidth may be as low as tens to thousands of bits per second
- worldwide reach may exist, but speed is highly constrained
- transmissions should be small and compact

Therefore, LXDR MUST assume:

- very small transmission budgets
- intermittent connectivity
- a need for compact data burst transmission

### 6.3 Data Burst as Primary Method

Project ADRIAN states that data processing with a data burst capability
must be the primary method.

Therefore:

- the protocol MUST support compact burst representations
- the protocol MUST minimize the transmitted set to absolute critical
  elements
- exchange structure MUST support trustworthy connectedness and
  synchronization with minimal exposure time

## 7. User Focus

Project ADRIAN centers the protocol on the forward tactical user.

For protocol design, this means:

- the tactical user is the point of request origination
- the user SHOULD not manually populate background metadata that can be
  generated automatically
- the protocol SHOULD assume limited local context and intermittent
  synchronization

## 8. Core Protocol Object Model

Based on Section 3.3, the minimum protocol object model is:

1. request header
2. one or more request segments
3. synchronized response that replaces local identity with enterprise
   identity

### 8.1 Request Header

The request header:

- precedes request information
- is generated with little to no user interaction
- provides the context for follow-on information gathering, status, and
  response
- exists to support disconnected resynchronization and limited bandwidth

The user has no direct interaction or control over the header data
elements.

### 8.2 Request Segments

The request segment holds the function- and sub-function-specific request
data.

The request header acts as a container for follow-on child or sub-requests.

The header indicates the number of request segments present.

### 8.3 Synchronized Response

The synchronized response replaces the local temporary request identifier
with the enterprise synchronized identifier.

This synchronized replacement is core to the protocol. It is not optional
metadata.

## 9. Request Header Requirements

Section 3.3.3 defines the request header as normative. The header contains
the following elements:

1. local system date
2. local system time
3. UTC time, calculated and not passed
4. military DTG, calculated and not passed
5. synchronized geospatial reference
6. local unique placeholder request identifier
7. synchronized enterprise request identifier
8. request priority
9. element or unit identification or callsign
10. request segment count

### 9.1 Header Exchange Rules

Section 3.3.3.2 states:

- some header fields are `Must Use`
- some are `Optional`
- some are `Pass`
- some are `Not Passed`
- the synchronized enterprise request identifier is returned as a
  `Sync Response`

Accordingly, LXDR MUST distinguish between:

- generated local header values
- transmitted header values
- calculated local-only values
- synchronized replacement values

### 9.2 Header Schema

Section 3.3.3.3 gives the canonical transmitted header schema as:

`CCYYMMMDD-HHMMSSDD-XXXX123456-XXXXXXXXXX-00-XXXX-00`

and provides the example:

`2027OCT13-15470352-4QFJ123456-3838JBNM5-02-KL9K-01`

This draft treats that schema as normative for the canonical header text
representation.

### 9.3 Synchronized Response Schema

Section 3.3.3.3 gives the synchronized response example as:

`3838JBNM5-KL9K15474QFJ`

This draft treats synchronized response as a first-class protocol action:

- identify the local placeholder ID
- provide the enterprise synchronized ID
- update local state accordingly

## 10. Function Families

Chapter 3 organizes request exchange requirements by logistics function.

From Sections 3.4 through 3.8, LXDR MUST be organized by function family:

- mobility
- supply
- maintenance
- general engineering
- health services

This draft currently records the following early request segment families as
normatively introduced in Sections 3.4 through 3.6:

- mobility PAX request
- mobility cargo request
- supply request
- maintenance request

Later families in Chapter 3 remain part of the protocol scope, but this
draft does not yet define their schemas because this document is limited to
the preamble and first extracted protocol baseline.

## 11. Mobility Request Requirements

Section 3.4 defines mobility as the doctrinal context for transportation
data requirements.

The protocol implications are:

- mobility request segments are separate from the request header
- embark data requirements apply across scale, from small tactical
  movements upward
- PAX and cargo movement are mode-agnostic at request time
- transportation mode selection occurs later in processing, not in the
  originating request

### 11.1 PAX Request Segment

Section 3.4.3.1 defines a PAX request segment. Its notable protocol
properties are:

- request type code `PM`
- one segment under a header
- one or more personnel identifiers
- earliest and latest departure dates
- departure and destination locations
- baggage weight
- hazardous material indication
- person count and total baggage weight as calculated, not passed

The canonical segment schema is given in Table 14.

### 11.2 Cargo Request Segment

Section 3.4.3.2 defines a cargo request segment. Its notable protocol
properties are:

- request type code `CM`
- cargo is movement of military materiel
- NIIN, quantity, serial number, weight, dimensions, HMIC, handling, dates,
  and locations are relevant
- item count and total weight are calculated, not passed

The canonical segment schema is given in Table 16.

## 12. Supply Request Requirements

Section 3.5 defines supply as a function with strong interdependence on
mobility and maintenance.

The protocol implications are:

- the tactical requester SHOULD not be burdened with commodity-chain
  complexity
- supply requests are intended to normalize requisition initiation across
  classes of supply
- a single request format may carry demands for many supply classes
- enterprise-side parsing and routing occur after initial exchange

### 12.1 Supply Request Segment

Section 3.5.3 defines a supply request segment. Its notable protocol
properties are:

- request type code `SR`
- NIIN may be present, but narrative description is permitted when unknown
- quantity is required
- required delivery date is required
- delivery location is required
- attachment is optional
- narrative is optional

The canonical segment schema is given in Table 17.

## 13. Maintenance Request Requirements

Section 3.6 defines maintenance as an active restoration function and
assumes the owner or operator is often the first maintenance sensor.

The protocol implications are:

- the request is for maintenance support as a service
- the request is not merely a parts order
- the originating user is not assumed to be a mechanic or technician
- the request begins a maintenance conversation that continues after
  enterprise synchronization

### 13.1 Maintenance Request Segment

Section 3.6.3 defines a maintenance request segment. Its notable protocol
properties are:

- request type code `CM` in that maintenance context
- serial number is optional if the item is not serialized
- NIIN, model, and nomenclature may be derived or entered
- piece count, operational condition, required date, and location are
  required
- maintenance support type, repair type, and major defect are coded
  selections
- attachment is supported
- narrative is optional

The canonical segment schema is given in Table 18.

This draft therefore does not assume that every two-character request type
is globally unique across all logistics functions. Function and segment
context are part of interpretation.

## 14. Protocol Design Constraints Derived from Sections 1-3

From the preceding sections, LXDR MUST satisfy the following:

- support digital creation at the tactical edge
- minimize user-entered metadata in favor of generated header data
- preserve local request identity before synchronization
- support synchronized replacement of local IDs by enterprise IDs
- encode only the minimum data required to initiate logistics action
- support multiple logistics function families
- support compact exchange suitable for data burst transmission
- preserve traceability through the lifecycle of the request
- allow function-specific request segments under a common request header

LXDR SHOULD satisfy the following:

- support both vertical and horizontal exchange pathways
- support attached supporting evidence where the request family permits it
- preserve context for follow-on status and response traffic
- allow local operation in disconnected mode with later synchronization

## 15. Non-Goals for This Draft

This draft does not yet define:

- cryptographic protections
- bearer-specific transport behavior
- fragmentation rules
- binary encoding rules
- protobuf or other structured schema mappings
- appendix-derived canonical field registries
- full general engineering and health segment schemas

Those belong to later protocol sections once this preamble is accepted.

## 16. Immediate Consequence for Implementation

An implementation based on this draft should begin with:

- a canonical header model
- canonical request segment models for the families already extracted
- canonical text render and parse rules for header and segments
- synchronized response handling
- a strict separation between:
    - generated header data
    - transmitted fields
    - calculated but non-transmitted fields

The implementation should remain subordinate to the doctrinal structure in
Project ADRIAN rather than inventing a generalized messaging system first.

## 17. Canonical Examples from Project ADRIAN Tables

This section consolidates the concrete protocol examples given in the
tables of Sections 3.3 through 3.5. These examples are source-derived and
are preserved here as canonical reference material for later codecs and
conformance tests.

### 17.1 Header Schema Example

From Table 8, the canonical header text schema is:

```text
CCYYMMMDD-HHMMSSDD-XXXX123456-XXXXXXXXXX-00-XXXX-00
```

The canonical header example is:

```text
2027OCT13-15470352-4QFJ123456-3838JBNM5-02-KL9K-01
```

The example represents:

- request date: `2027OCT13`
- request time: `15470352`
- synchronized physical location: `4QFJ123456`
- local placeholder request ID: `3838JBNM5`
- priority: `02`
- element or unit identifier: `KL9K`
- request segment count: `01`

### 17.2 Synchronized Response Example

From Table 9, the canonical synchronized response example is:

```text
3838JBNM5-KL9K15474QFJ
```

The synchronized response replaces:

- local placeholder ID: `3838JBNM5`

with:

- enterprise synchronized ID: `KL9K15474QFJ`

### 17.3 UMMIPS Priority Mapping

From Table 5, the canonical priority mapping used by the request header is:

```text
F/AD I:   A=01 B=04 C=11
F/AD II:  A=02 B=05 C=12
F/AD III: A=03 B=06 C=13
F/AD IV:  A=07 B=09 C=14
F/AD V:   A=08 B=10 C=15
```

This table is not a single serialized example, but it is canonical input
for deriving the two-digit request priority field.

### 17.4 ZAP Schema Example

From Table 11, the ZAP number is constructed as:

```text
Position 1:   Last Name Initial
Position 2:   First Name Initial
Position 3-6: Last 4 digits of EDI-PI
Position 7-8: Blood Group
Position 9:   Blood RhD
```

### 17.5 ZAP Conversion Example

From Table 12, the canonical ZAP conversion example is:

```text
Last Name:      Smith-Hernandez
First Name:     Chris
EDI-PI:         1010919789
Blood Group:    AB
Blood RhD:      Positive
ZAP:            SC9789AB+
```

### 17.6 PAX Request Schema Example

From Table 14, the canonical PAX request schema is:

```text
0-XX-00-XXXXXXXXXX-CCYYMMMDD-CCYYMMMDD-XXXX123456-XXXX123456-000-X
```

The canonical PAX request example for one passenger is:

```text
1-PM-02-1010919789-2027OCT15-2027OCT20-4QFJ123456-4QFJ456789-075-X
```

This example represents:

- segment number: `1`
- request type: `PM`
- request priority: `02`
- passenger identifier: `1010919789`
- earliest departure date: `2027OCT15`
- latest departure date: `2027OCT20`
- departure location: `4QFJ123456`
- destination location: `4QFJ456789`

## 18. Role of the Appendices

Sections 1 through 3 establish the protocol rationale, operating
constraints, request header concept, and the first explicit exchange
schemas.

The appendices serve a different but essential purpose:

- `Appendix A` through `Appendix E` enumerate the data-point universe by
  logistics function and sub-function.
- `Appendix F` organizes logistics data into canonical data blocks.
- `Appendix G` identifies authoritative data sources and managers.

Accordingly, the appendices are not merely reference material. They are
the basis for protocol completeness.

### 18.1 What the Appendices Contribute to LXDR

The appendices enable LXDR to define:

- a canonical field registry
- a canonical block registry
- a canonical file or request-container registry
- shared field reuse across functions and sub-functions
- provenance and authority metadata for later validation work

### 18.2 What the Appendices Do Not Imply

The appendices do not imply that every listed data point must be passed in
every transmission.

Project ADRIAN is explicit that the protocol is concerned with minimum
critical data at the tactical edge. Therefore:

- Appendix field inventories define the available data vocabulary
- Chapters 1 through 3 define the minimum exchange philosophy
- LXDR must select and normalize, not blindly transmit every field

### 18.3 Protocol Consequence

LXDR should therefore be built in three layers of rigor:

1. chapters 1 through 3 define the request and exchange model
2. appendices A through E define the function-specific field universe
3. appendix F defines the canonical field, block, and file organization

Appendix G should inform provenance, validation, and governance but should
not be treated as direct wire-format material.

## 19. Initial Canonical Registry Model from Appendix F

Appendix F is titled `Logistics Data Blocks` and is the strongest basis in
Project ADRIAN for a normalized protocol registry.

The visible hierarchy in Appendix F is:

- `Activity`
- `Data Element`
- `Data Field`

This aligns with the chapter 3 hierarchy of:

- data field
- data block
- data file

For LXDR purposes, the initial canonical registry model should therefore
be:

1. canonical file
2. canonical block
3. canonical field

where:

- a canonical file is a request container or other top-level exchangeable
  object
- a canonical block is a logical grouping of related fields
- a canonical field is the atomic named datum

### 19.1 Initial Registry Shape

The initial registry entries should preserve Appendix F provenance. A
single canonical field entry SHOULD include:

- `activity`
- `data_element`
- `data_field`
- `canonical_file`
- `canonical_block`
- `canonical_field`
- `source_reference`
- `exchange_role`

Where:

- `activity` is the source system or organizational activity named by
  Appendix F
- `data_element` is the Appendix F grouping label
- `data_field` is the Appendix F field label
- `canonical_file` is the LXDR top-level container classification
- `canonical_block` is the normalized block classification
- `canonical_field` is the normalized field name
- `source_reference` points back to the ADRIAN appendix entry
- `exchange_role` indicates whether the field is:
  - transmitted
  - calculated
  - synchronized
  - local-only

### 19.2 Early Observations from Appendix F

The excerpted Appendix F rows already show several recurring canonical
blocks that are strong LXDR candidates:

- `Header`
- `Activity`
- `Requisitioning`
- `Location, Physical`
- `Item Identification`
- `Item Catalog Information`
- `Item Attributes`
- `Item Characteristics`
- `Manifest`
- `Transportation`
- `Individual Member`
- `Date`
- `Communications`
- `Financials`
- `Report, General Engineering`
- `Media Object`
- `Map Overlay Object`
- `Estimate`
- `Position Information`

These recurring blocks strongly suggest that:

- function-specific segments should be composed from a shared canonical
  block vocabulary
- different functions will reuse many of the same atomic fields
- request segments should not be modeled as isolated one-off forms

### 19.3 Initial Canonical File Interpretation

From chapters 1 through 3 plus Appendix F, the most defensible initial
canonical file types are:

- `request_header`
- `request_segment`
- `request_container`
- `sync_response`
- `attachment_or_media_reference`

This draft does not yet define every canonical file formally, but this is
the correct initial organizing direction.

### 19.4 Immediate Use in the Draft

Appendix F should next be used to:

- normalize the header fields already defined in Section 9
- map mobility, supply, and maintenance segment fields into shared blocks
- identify which fields recur across functions and should have one
  canonical name
- prepare for later formal schema work, including protobuf or other typed
  representations

## 20. Protobuf Mapping Rules

This draft permits a protobuf schema as the canonical typed model for
LXDR-Core, provided that the protobuf schema remains subordinate to the
normative protocol text in this document.

Protobuf is therefore an implementation and interchange schema for the
typed protocol model. It is not, by itself, the full protocol.

### 20.1 Scope of Protobuf in LXDR

At the current maturity of this draft, protobuf is appropriate for:

- request header
- synchronized response
- request container
- typed request segments
- canonical field, block, and file registries once they are formally
  extracted

At the current maturity of this draft, protobuf is not yet authoritative
for:

- bearer-specific transport behavior
- fragmentation policy
- cryptographic behavior
- all future function families not yet defined in this draft

### 20.2 Mapping Principles

The protobuf schema MUST follow these rules:

1. The protobuf schema MUST NOT invent fields absent from this draft.
2. The protobuf schema MUST preserve the distinction between:
   - generated fields
   - transmitted fields
   - calculated but non-transmitted fields
   - synchronized replacement fields
3. The protobuf schema MUST preserve the request hierarchy:
   - header
   - one or more request segments
   - synchronized response as a separate object
4. The protobuf schema SHOULD preserve function-family context so that
   similarly named request types are not assumed globally unique.
5. The protobuf schema SHOULD prefer neutral field typing when the draft
   defines a canonical text code or fixed-format value but does not yet
   define its deeper numeric or enum domain.

### 20.3 Current Typed Mapping Choice

At the current stage of protocol extraction, many fields are represented
best as strings in the protobuf schema because the draft defines them by:

- canonical textual format
- fixed-width code
- structured identifier
- coded value whose full enum domain is not yet extracted into this draft

Examples include:

- local system date
- local system time
- geospatial references
- request type codes
- priority codes
- NIIN values
- serial numbers
- maintenance defect codes

Fields should only be promoted to stricter enums or numerics when the
protocol draft explicitly defines the stable domain required to do so.

### 20.4 Relationship to Canonical Text

The canonical text examples in Section 17 remain normative examples of the
exchange shape. The protobuf schema is a typed representation of the same
objects.

Therefore:

- canonical text examples remain valid conformance references
- protobuf messages represent the same request objects in a typed form
- future codecs may translate between:
  - protobuf
  - canonical text burst
  - other compact representations

### 20.5 Relationship to Future Packed Encodings

If LXDR later defines packed or burst-optimized binary forms, those packed
forms SHOULD be derived from the canonical protocol objects and field
registries, not treated as an alternate semantic model.

In that model:

- this document remains the protocol authority
- protobuf remains the canonical typed schema
- packed encodings become a constrained representation of the same objects
- baggage weight: `075`
- hazardous material type: `X`

### 17.7 Cargo Request Schema Example

From Table 16, the canonical cargo request schema is:

```text
0-XX-00-000000000-00000-XXXXXXXXXX-00000-000-000-000-X-X-CCYYMMMDD-CCYYMMMDD-XXXX123456-XXXX123456
```

The canonical cargo request example for one MTVR truck, serial number
`598742`, is:

```text
1-CM-02-015519434-1-598742-28000-126-100-315-D-R-2027OCT15-2027OCT20-4QFJ123456-4QFJ456789
```

This example represents:

- segment number: `1`
- request type: `CM`
- request priority: `02`
- NIIN: `015519434`
- item quantity: `1`
- serial number: `598742`
- gross weight: `28000`
- height: `126`
- width: `100`
- length: `315`
- HMIC: `D`
- handling: `R`
- earliest departure date: `2027OCT15`
- latest departure date: `2027OCT20`
- departure location: `4QFJ123456`
- destination location: `4QFJ456789`
