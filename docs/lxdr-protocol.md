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

This draft currently defines the following Chapter 3 request segment
families in schema and validation form:

- mobility PAX request
- mobility cargo request
- supply request
- maintenance request
- engineer reconnaissance area report
- engineer reconnaissance zone report
- engineer reconnaissance route report
- engineer reconnaissance road report
- engineer reconnaissance landing zone report
- general engineering obstacle removal
- explosive ordnance disposal clearing/rendering safe
- general engineering bulk liquid support
- general engineering demolition
- health services collection
- health services triage
- health services intervention
- health services hold
- health services evacuate (CASEVAC)

The remaining Chapter 3 reports listed below remain protocol scope, but
are intentionally deferred from the current v1 implementation baseline.

The following Chapter 3 reports are intentionally deferred for a later
protocol increment:

- engineer reconnaissance tunnel report
- engineer reconnaissance bridge report
- engineer reconnaissance ford report
- engineer reconnaissance ferry report
- engineer reconnaissance river report

### 10.1 Engineer Reconnaissance Area Report

Section 3.7.3.1 defines the area report as a general engineering report
for threat, terrain, society, and infrastructure within an area a unit
must move through or occupy.

Its notable protocol properties are:

- date of evaluation or reconnaissance
- element leader identifier
- area location
- optional water source expressed as flow-rate and quantity pair
- repeatable item reports keyed by location plus free-text label
- optional attachment flag
- optional narrative

The canonical segment schema is given in Table 19.

### 10.2 Engineer Reconnaissance Zone Report

Section 3.7.3.2 defines the zone report for vague threat situations or
when cross-country trafficability information is desired.

Its notable protocol properties are:

- date of evaluation or reconnaissance
- element leader identifier
- optional amphibious crossing classification
- repeatable enemy reports keyed by location plus free-text label
- zone location
- optional water source expressed as flow-rate and quantity pair
- repeatable item reports keyed by location plus free-text label
- optional attachment flag
- optional narrative

The canonical segment schema is given in Table 20.

### 10.3 Engineer Reconnaissance Route Report

Section 3.7.3.3 defines the route report for reconnaissance along a
specified route and adjacent terrain, including lateral routes where
movement may be restricted.

Its notable protocol properties are:

- date of evaluation or reconnaissance
- element leader identifier
- enemy report references
- repeatable route locations or waypoints
- optional water source expressed as flow-rate and quantity pair
- repeatable item reports keyed by location plus free-text label
- optional attachment flag
- optional narrative

The canonical segment schema is given in Table 21.

### 10.4 Engineer Reconnaissance Road Report

Section 3.7.3.4 defines the road report as an engineer reconnaissance
segment for evaluating established-road condition and trafficability
between a predetermined start and end point.

Its notable protocol properties are:

- date of evaluation/reconnaissance
- start point location
- end point location
- closed one-character code domains for:
  - road classification
  - drainage
  - foundation
  - surface type
- required obstructions field
- optional attachment flag
- optional narrative

The canonical segment schema is given in Table 22.

#### 10.4.1 Road Classification

- `A`
  flat gradients and easy curves
- `B`
  steep gradient greater than 7 percent
- `C`
  sharp curves less than 30 meters
- `D`
  steep gradients and sharp curves

#### 10.4.2 Drainage

- `A`
  adequate crown/camber with adequate culverts in good condition
- `B`
  inadequate; poor crown/camber or blocked/poor culverts or ditches

#### 10.4.3 Foundation

- `A`
  stabilized and compacted
- `B`
  unstable or loose

#### 10.4.4 Surface Type

- `A`
  concrete
- `B`
  bituminous
- `C`
  brick
- `D`
  stone
- `E`
  crushed rock
- `F`
  water bound macadum
- `G`
  gravel
- `H`
  lightly metaled
- `I`
  natural or stabilized soil
- `J`
  other

### 10.5 Engineer Reconnaissance Landing Zone Report

Section 3.7.3.10 defines the landing zone report as an engineer
reconnaissance segment for evaluating usable terrain for aircraft landing
and takeoff. It is not an airfield or staging-facility report.

Its notable protocol properties are:

- date of evaluation/reconnaissance
- location
- estimate flag
- layout designation
- landing point, landing zone, and landing site capacities
- landing zone width and length
- aircraft supportability code
- approach and departure direction codes
- surface slope code
- obstacle code
- optional attachment flag
- optional narrative

The canonical segment schema is given in Table 28.

#### 10.5.1 Estimate

- `1`
  yes
- `0`
  no

#### 10.5.2 Layout Designation

- `LZ`
  landing zone
- `LS`
  landing site
- `LP`
  landing point

#### 10.5.3 Aircraft Supportability

- `A`
  size 4 aircraft
- `B`
  size 4 aircraft
- `C`
  size 2 aircraft
- `D`
  size 2 aircraft
- `E`
  vertical lift requires VTOL
- `F`
  vertical lift requires VTOL

#### 10.5.4 Approach and Departure Direction

- `S`
  south
- `E`
  east
- `N`
  north
- `W`
  west

#### 10.5.5 Surface Slope

- `A`
  less than 7 degrees upslope landing
- `B`
  greater than 7 degrees side slope landing

#### 10.5.6 Obstacle

- `1`
  remove
- `2`
  reduce
- `3`
  emplace

### 10.6 General Engineering Obstacle Removal

Section 3.7.3.11 defines obstacle removal as a general engineering
segment for describing an obstacle, its dimensions, bypass conditions,
and the engineer recommendation for action.

Its notable protocol properties are:

- date of evaluation/reconnaissance
- location
- obstacle code
- min/max length, width, and depth
- route number
- determination of action code
- bypass code
- bypass grid
- optional attachment flag
- optional narrative

The canonical segment schema is given in Table 29.

#### 10.6.1 Obstacle

- `1`
  remove
- `2`
  reduce
- `3`
  emplace

#### 10.6.2 Determination of Action

- `1`
  use bypass
- `2`
  neutralize obstacle
- `3`
  breach
- `4`
  continue mission

#### 10.6.3 Bypass

- `1`
  yes
- `0`
  no

### 10.7 Explosive Ordnance Disposal Clearing/Rendering Safe

Section 3.7.3.12 defines the EOD clearing/rendering safe segment for
reporting UXO discovery, requested action timing, location, CBRN
characteristics, munition details, and optional media.

Its notable protocol properties are:

- date of UXO discovery
- requested date of EOD action
- location of UXO
- type of CBRN agent
- optional physical property of CBRN agent
- optional contamination value of CBRN agent
- munition color
- munition markings
- munition purpose
- munition type
- optional attachment flag
- optional narrative

The canonical segment schema is given in Table 30.

#### 10.7.1 Type of CBRN Agent

- `1`
  chemical
- `2`
  biological
- `3`
  radiological
- `4`
  nuclear
- `5`
  none

#### 10.7.2 Physical Property of CBRN Agent

- `1`
  gas
- `2`
  liquid
- `3`
  aerosols
- `4`
  ionizing radiation

#### 10.7.3 Contamination Value of CBRN Agent

- `E`
  external
- `I`
  internal
- `W`
  wound
- `C`
  contagious casualties

#### 10.7.4 Munition Purpose

- `AA`
  anti-armor
- `AP`
  anti-personnel
- `FL`
  flare
- `SM`
  smoke
- `IM`
  improvised

#### 10.7.5 Munition Type

- `E`
  emplace
- `D`
  drop
- `T`
  throw
- `P`
  project

### 10.8 General Engineering Bulk Liquid Support

Section 3.7.3.13 defines bulk liquid support as a report of fuels and
water sources available at local or commercial sites.

Its notable protocol properties are:

- date of evaluation or reconnaissance
- bulk-liquid location
- estimate flag
- optional fuel payload encoded as a typed quantity and unit string
- optional water payload encoded as a typed quantity, unit, and water-type string
- optional attachment flag
- optional narrative

The canonical segment schema is given in Table 31.

### 10.9 General Engineering Demolition

Section 3.7.3.14 defines demolition as a general engineering request when
an object presents an obstacle to mobility or area use.

Its notable protocol properties are:

- date of evaluation or reconnaissance
- object location
- type of demolition as source-defined alpha code text
- route number
- determination of action
- bypass flag
- bypass grid
- optional attachment flag
- optional narrative

The canonical segment schema is given in Table 32.

### 10.10 Health Services Collection

Section 3.8.3.1 defines the collection segment as the initial casualty
record created at the point of collection.

Its notable protocol properties are:

- segment number
- request type code `CR`
- request priority derived from triage precedence
- ZAP or EDI-PI
- last name
- first name
- optional service code
- optional element/unit identification or callsign
- allergies/alerts
- date of injury
- time of injury
- location injury occurred

The canonical segment schema is given in Table 33.

#### 10.10.1 Service

When service is manually entered, the source text defines these codes:

- `USA`
  United States Army
- `USSF`
  United States Space Force
- `USAF`
  United States Air Force
- `USCG`
  United States Coast Guard
- `USN`
  United States Navy
- `USMC`
  United States Marine Corps
- `US CIV`
  United States civilian
- `NON-US`
  non-U.S. coalition or NATO personnel
- `EPW`
  enemy prisoner of war

### 10.11 Health Services Treatment Data (Triage Segment)

Section 3.8.3.2 defines triage as the first treatment action following
collection. The protocol-relevant fields in Table 34 are:

- primary mechanism of injury
- CBRN related/exposure
- major signs/symptoms
- one to ten injury locations
- optional triage/treatment vitals
- optional responsiveness and pain scale
- triage precedence

The source gives closed code domains for the main triage selectors.

#### 10.11.1 Primary Mechanism of Injury

- `E1`
  artillery/indirect fire, grenade
- `E2`
  improvised explosive device
- `E3`
  landmine
- `E4`
  rocket propelled grenade
- `E5`
  other type of explosion
- `P1`
  blunt force trauma
- `P2`
  fall
- `P3`
  velocity impact
- `P4`
  gunshot wound
- `P5`
  knife/bayonet/sharp object
- `D1`
  pathogen
- `D2`
  allergic reaction
- `D3`
  chemical/biological hazard exposure
- `D4`
  pathological

#### 10.11.2 CBRN Related/Exposure

- `C`
  chemical
- `B`
  biological
- `R`
  radiological
- `N`
  nuclear
- `X`
  none

#### 10.11.3 Major Signs/Symptoms

- `B`
  bleeding
- `R`
  restricted breathing
- `X`
  burned
- `C`
  incapacitation

#### 10.11.4 Pulse Location

- `W`
  wrist
- `N`
  neck

#### 10.11.5 Responsiveness

- `A`
  alert
- `V`
  responds to voice
- `P`
  responds to pain
- `U`
  unresponsive

#### 10.11.6 Triage Precedence

- `A`
  urgent surgical
- `B`
  urgent non-surgical
- `C`
  priority
- `D`
  routine
- `E`
  convenience

### 10.12 Health Services Treatment Data (Intervention Segment)

Section 3.8.3.2 continues treatment with the intervention segment in
Table 35. The protocol-relevant fields are:

- repeated tourniquet treatment entries
- one to seven wound treatment selections
- one airway treatment selection
- one breathing treatment selection
- repeated fluid, blood, and medication treatment entries
- casualty type
- first responder ZAP or EDI-PI

#### 10.12.1 Tourniquet Placement

- `TQXX`
  no tourniquets
- `TQRA`
  right arm
- `TQLA`
  left arm
- `TQRL`
  right leg
- `TQLL`
  left leg

#### 10.12.2 Tourniquet Type

- `E`
  extremity
- `J`
  junctional
- `T`
  truncal

#### 10.12.3 Wounds Treatment

- `T1`
  hemostatic dressing
- `T2`
  pressure bandage
- `T3`
  sling or splint
- `T4`
  eye shield left
- `T5`
  eye shield right
- `T6`
  eye shield both
- `T7`
  hypothermia prevention

#### 10.12.4 Airway Treatment

- `A0`
  airway intact
- `A1`
  nasopharyngeal airway
- `A2`
  cricothyroidotomy
- `A3`
  endotracheal tube
- `A4`
  supraglottic airway

#### 10.12.5 Breathing Treatment

- `B0`
  none
- `B1`
  O2
- `B3`
  needle decompression
- `B4`
  chest tube
- `B5`
  chest seal

#### 10.12.6 Fluid Circulation Treatment

- fluid names:
  - `S` saline
  - `R` ringer's lactate
  - `H` hextend
- routes:
  - `IV` intravenous
  - `IO` interosseous

#### 10.12.7 Blood Circulation Treatment

- `WBD`
  whole blood
- `RBC`
  red blood cells
- `FFP`
  fresh frozen plasma
- `FDP`
  lyophilized plasma
- routes:
  - `IV` intravenous
  - `IO` interosseous

#### 10.12.8 Analgesic Medication

- medication names:
  - `K` ketamine
  - `F` fentanyl
  - `M` morphine
- routes:
  - `R1` intrathecal
  - `R2` subcutaneous
  - `R3` intravenous
  - `R4` intramuscular
  - `R5` oral
  - `R6` inhale
  - `R7` rectal

#### 10.12.9 Antibiotic Medication

- medication names:
  - `M` moxifloxacin
  - `E` ertapenem
  - `P` penicillin
  - `A` azithromycin
- routes:
  - `R1` intrathecal
  - `R2` subcutaneous
  - `R3` intravenous
  - `R4` intramuscular
  - `R5` oral
  - `R6` inhale
  - `R7` rectal

#### 10.12.10 Other Medication

- medication names:
  - `I` ibuprofen
  - `T` tranexamic acid
- routes:
  - `R1` intrathecal
  - `R2` subcutaneous
  - `R3` intravenous
  - `R4` intramuscular
  - `R5` oral
  - `R6` inhale
  - `R7` rectal

#### 10.12.11 Casualty Type

- `A`
  litter
- `B`
  ambulatory

### 10.13 Health Services Hold

Section 3.8.3.3 does not introduce a new closed-field table. Instead,
it defines `hold` as continued monitoring and continued treatment after
initial treatment.

The protocol implication is:

- no existing data is overwritten
- additional vitals entries are appended
- additional treatment entries are appended
- all appended entries are date/time driven

For this draft, `hold` is modeled as an append-only extension of the
existing triage and intervention records. It introduces no new code
domains beyond those already defined in sections `10.11` and `10.12`.

### 10.14 Health Services Evacuate (CASEVAC)

Section 3.8.3.4 defines `evacuate (CASEVAC)` as the movement of
casualties to and/or between medical treatment facilities.

The source text makes two protocol implications explicit:

- CASEVAC is a specialized `PAX` movement request
- the evacuation segment auto-carries selected casualty data from
  collection, treatment, and hold with minimal user intervention

This draft models CASEVAC as a health-services segment with:

- evacuation request priority
- pickup location
- optional pickup-site marking
- optional contamination indicator
- contact settings
- repeated casualty-precedence counts
- repeated casualty-type counts
- optional requested equipment
- optional security status
- repeated casualty records carrying the auto-populated medical payload

#### 10.14.1 Evacuation Request Priority

- `A`
  urgent surgical
- `B`
  urgent non-surgical
- `C`
  priority
- `D`
  routine
- `E`
  convenience

#### 10.14.2 Location Marking

- `A`
  panels
- `B`
  pyrotechnic
- `C`
  smoke
- `D`
  IR flash or signal beacon
- `E`
  none

#### 10.14.3 Location Contamination

- `A`
  nuclear or radiological
- `B`
  biological
- `C`
  chemical
- `D`
  none

#### 10.14.4 Casualty Type Count

- `A`
  ambulatory
- `L`
  litter

#### 10.14.5 Requested Equipment

- `A`
  none
- `B`
  hoist
- `C`
  extraction equipment
- `D`
  ventilation

#### 10.14.6 Security

- `N`
  no enemy troops in area
- `P`
  possible enemy troops in area
- `E`
  enemy troops in area
- `X`
  enemy troops in area, armed escort required

#### 10.14.7 Casualty Record Payload

Segment `47` in Table 36 is repeatable and is keyed by the casualty
`ZAP/EDI-PI`. The source text states that, when a casualty is selected
for evacuation, the following are auto-populated into the evacuation
request:

- mechanism of injury
- major injury code
- last recorded vitals readings
- last recorded treatments
- casualty CBRN code

This draft therefore models the CASEVAC casualty payload as a repeatable
record keyed by `ZAP/EDI-PI`, reusing the already-defined health code
domains where the whitepaper is explicit and preserving unresolved
subfields as strings.

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

### 13.2 Maintenance Code Domains

Section 3.6.3 defines several maintenance fields as user-facing clear text
choices that are exchanged as compact codes. For protocol purposes, these
domains are closed and normative.

#### 13.2.1 Equipment Operational Condition

The maintenance operational condition is a single alpha code:

- `A`
  operational; can perform all combat functions
- `B`
  minor; can perform most combat functions
- `C`
  degraded/deadline; cannot perform all combat functions

#### 13.2.2 Type of Maintenance Support Required

The maintenance support type is a two-character alphanumeric code:

- `XX`
  none
- `R1`
  repair
- `R2`
  recovery retrieve
- `R3`
  recovery free from immobility
- `R4`
  contact team

#### 13.2.3 Type of Repair

The repair type is a two-character alphanumeric code. Section 3.6.3
states this list activates only when the maintenance support type is
`R1`.

- `M1`
  modification
- `S1`
  servicing preventive maintenance
- `S2`
  servicing tune/adjust
- `C1`
  calibration
- `D1`
  repair defect

#### 13.2.4 Repair Major Defect

The repair major defect is a four-character alphanumeric code. Section
3.6.3 states this list activates only when the maintenance support type
is `R1` or `R4`.

- `MD01`
  engine
- `MD02`
  fuel
- `MD03`
  exhaust
- `MD04`
  cooling
- `MD05`
  lubrication
- `MD06`
  electrical/fuse
- `MD07`
  transmission/power train
- `MD08`
  chassis/body/frame/armor
- `MD09`
  suspension/axle/differential/shocks
- `MD10`
  brakes/tires/wheels/hub/track
- `MD11`
  lift/winch/boom/turret
- `MD12`
  hydraulics
- `MD13`
  receiver/input circuitry
- `MD14`
  transmitter/output circuitry
- `MD15`
  display/screen/video
- `MD16`
  weapon system
- `NMAJ`
  routine/no major defect

#### 13.2.5 Attachment Indicator

The attachment indicator is a single digit flag:

- `1`
  attachment present
- `0`
  no attachment present

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
- delivery receipts or acknowledgement semantics
- message routing policy
- transport-interface behavior
- full Appendix F registry population
- deferred engineer reconnaissance reports:
  - tunnel
  - bridge
  - ford
  - ferry
  - river

Those belong to later protocol increments beyond the current v1 draft.

## 16. Immediate Consequence for Implementation

An implementation based on this draft should begin with:

- protobuf-backed core schema as the canonical typed model
- canonical header and synchronized-response handling
- schema and validation for the implemented Chapter 3 request segments
- canonical text render and parse rules where Project ADRIAN provides
  normative serialized examples
- a minimal LXDR-Link frame carrying exactly one valid core payload
- binary marshal and unmarshal support for current core and link objects
- a strict separation between:
    - generated header data
    - transmitted fields
    - calculated but non-transmitted fields

The implementation should remain subordinate to the doctrinal structure in
Project ADRIAN rather than inventing transport or application layers first.

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

### 17.7 Cargo Request Schema Example

From Table 16, the canonical cargo request schema is:

```text
0-XX-00-XXXXXXXXX-0-XXXXXX-00000-000-000-000-X-X-CCYYMMMDD-CCYYMMMDD-XXXX123456-XXXX123456
```

The canonical cargo request example for one item is:

```text
1-CM-02-015519434-1-598742-28000-126-100-315-D-R-2027OCT15-2027OCT20-4QFJ123456-4QFJ456789
```

This example represents:

- segment number: `1`
- request type: `CM`
- request priority: `02`
- NIIN: `015519434`
- quantity: `1`
- serial number: `598742`
- gross weight pounds: `28000`
- height inches: `126`
- width inches: `100`
- length inches: `315`
- HMIC: `D`
- handling: `R`
- earliest departure date: `2027OCT15`
- latest departure date: `2027OCT20`
- departure location: `4QFJ123456`
- destination location: `4QFJ456789`

Section 3.4.3.2 also defines two closed one-character code domains for
cargo movement.

#### 17.7.1 HMIC

The hazardous material indicator code is a single alphabetic code:

- `D`
  no information in HMIRS; the NSN is in a Federal Stock Class where an
  MSDS should be available
- `N`
  no HMIRS data and the NSN is in a class not generally suspected of
  hazardous materials
- `P`
  no HMIRS data; an MSDS may be required depending on hazard
  determination or end use
- `Y`
  information is found in HMIRS

#### 17.7.2 Handling

The cargo handling code is a single alphabetic code:

- `C`
  crane or hoist
- `M`
  material handling equipment
- `T`
  tow
- `R`
  ramp
- `X`
  none

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

The early Appendix F `Header` rows also show an important modeling
constraint: Appendix F header-class metadata is broader than the Section 9
LXDR transmitted request header. Registry entries for Appendix F `Header`
must therefore preserve source provenance even when they do not map
one-for-one to the tactical request header fields.

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

## 21. Section 4 Implications

Section 4 of Project ADRIAN does not replace the protocol core defined by
chapters 1 through 3, but it does clarify the intended operational system
that LXDR is meant to enable.

This section records only the protocol-relevant implications of Section 4.

### 21.1 Requests Are Broader Than Materiel

Section 4 explicitly broadens the meaning of a request into three classes:

- request for information
- request for service
- request for materiel

Therefore, LXDR must not be interpreted as a materiel-only requisition
format. The protocol core should be capable of carrying requests across
those three classes even when a given function or segment family is more
heavily associated with one class than the others.

### 21.2 Supply Chain Scope Is Broader Than Materiel

Section 4 defines the supply chain in broader terms than traditional
material flow. In protocol terms, the exchanged object may represent a
demand for:

- information
- service
- materiel

This reinforces the requirement that LXDR request objects remain function-
 and demand-oriented rather than hard-coded to inventory issuance alone.

### 21.3 Mobility Is the Intended Function Frame

Section 4 recommends a doctrinal shift from `transportation` to
`mobility`.

For protocol purposes, this means:

- mobility is the intended top-level function frame for movement-related
  requests
- transportation mode selection remains a later processing activity
- future mobility request families may need to accommodate:
  - intermodal planning
  - littoral connectors
  - unmanned systems
  - small aerial craft

This draft does not yet add those future mobility segment families, but
their possibility should be assumed in future schema growth.

### 21.4 Planning and Forecasting Are Protocol Consumers

Section 4 identifies a major gap between operations, logistics, and
planning tools. The desired end state is a common operational picture and
planner workbench that integrates logistics digitally.

For LXDR, the immediate implication is not that planning objects become
wire-format messages now. The implication is:

- planning
- forecasting
- readiness
- execution tracking

are future consumers and producers of LXDR-aligned request data.

Therefore, LXDR should continue to emphasize:

- stable request identity
- synchronization
- function-family structure
- minimum critical data at the point of origin

so that later planning and COP layers can consume the same protocol
objects without redefining them.

### 21.5 Section 4 Is Guidance for Future Extension

Section 4 should currently be treated as:

- operational design guidance
- future protocol extension guidance
- future application and COP guidance

It should not yet be treated as a source of new wire fields unless the
whitepaper provides explicit exchange structures comparable to those in
Section 3.

## 22. Initial LXDR-Link Frame

LXDR v1 now distinguishes between:

- `LXDR-Core`
  the ADRIAN-derived typed objects such as:
  - request container
  - synchronized response
  - canonical registry
- `LXDR-Link`
  the carried wire object that transports exactly one core payload

The initial `LXDR-Link` frame is intentionally minimal. It exists only to
distinguish what kind of LXDR core object is being exchanged, without yet
introducing transport, routing, fragmentation, or bearer-specific
metadata.

### 22.1 Link Frame Scope

The initial frame supports exactly one payload of:

- `RequestContainer`
- `SynchronizedResponse`
- `CanonicalRegistry`

In addition, the initial frame carries three link-level metadata fields:

- `link_message_id`
  a stable identifier for the carried frame
- `delivery_method`
  the delivery intent class for the frame:
  - `DIRECT`
  - `PROPAGATED`
  - `OPPORTUNISTIC`
- `representation`
  the carried representation of the core payload

In the current implementation, `representation` is constrained to:

- `BINARY_PROTO`

This means the frame is a typed wrapper around one core object. It does
not yet define:

- routing metadata
- acknowledgements
- fragmentation fields
- integrity fields
- correlation or reply-reference fields
- transport hints

Those remain later concerns.

### 22.2 Validation Rules

An `LXDR-Link` frame is valid only when:

- `link_message_id` is present
- `delivery_method` is a known link-level delivery intent
- `representation` is a known link-level representation
- exactly one payload is present
- the carried payload is itself valid under the corresponding `LXDR-Core`
  rules

The link layer therefore does not weaken core validation. It only adds a
typed carried-message boundary.

The current implementation also constrains the frame/payload combination:

- when a payload is carried directly as a protobuf-typed object inside the
  frame, `representation` must be `BINARY_PROTO`

Future canonical-text or packed link representations may require
different payload carriers, but those are not part of this draft yet.

## 23. Conformance Matrix

Status markers in this table mean:

- `[x]`
  implemented in the current repo
- `[ ]`
  not yet implemented

The columns are:

- `Spec`
  present in this draft
- `Proto`
  modeled in `proto/lxdr/v1/lxdr.proto`
- `Validate`
  validated in Go
- `Canon`
  canonical text parse/render implemented
- `Fixture`
  canonical example fixture present

| Object / Segment | Spec | Proto | Validate | Canon | Fixture |
| --- | --- | --- | --- | --- | --- |
| LXDR-Link frame | [x] | [x] | [x] | [ ] | [ ] |
| Request header | [x] | [x] | [x] | [x] | [x] |
| Synchronized response | [x] | [x] | [x] | [x] | [x] |
| Mobility PAX | [x] | [x] | [x] | [x] | [x] |
| Mobility cargo | [x] | [x] | [x] | [x] | [x] |
| Supply request | [x] | [x] | [x] | [ ] | [ ] |
| Maintenance request | [x] | [x] | [x] | [ ] | [ ] |
| Engineer area report | [x] | [x] | [x] | [ ] | [ ] |
| Engineer zone report | [x] | [x] | [x] | [ ] | [ ] |
| Engineer route report | [x] | [x] | [x] | [ ] | [ ] |
| Engineer road report | [x] | [x] | [x] | [ ] | [ ] |
| Engineer landing zone report | [x] | [x] | [x] | [ ] | [ ] |
| Obstacle removal | [x] | [x] | [x] | [ ] | [ ] |
| EOD clearing/rendering safe | [x] | [x] | [x] | [ ] | [ ] |
| Bulk liquid support | [x] | [x] | [x] | [ ] | [ ] |
| Demolition | [x] | [x] | [x] | [ ] | [ ] |
| Engineer tunnel report | [ ] | [ ] | [ ] | [ ] | [ ] |
| Engineer bridge report | [ ] | [ ] | [ ] | [ ] | [ ] |
| Engineer ford report | [ ] | [ ] | [ ] | [ ] | [ ] |
| Engineer ferry report | [ ] | [ ] | [ ] | [ ] | [ ] |
| Engineer river report | [ ] | [ ] | [ ] | [ ] | [ ] |
| Health collection | [x] | [x] | [x] | [ ] | [ ] |
| Health triage | [x] | [x] | [x] | [ ] | [ ] |
| Health intervention | [x] | [x] | [x] | [ ] | [ ] |
| Health hold | [x] | [x] | [x] | [ ] | [ ] |
| Health evacuate / CASEVAC | [x] | [x] | [x] | [ ] | [ ] |
| Appendix F canonical registry | [x] | [x] | [x] | [ ] | [ ] |
