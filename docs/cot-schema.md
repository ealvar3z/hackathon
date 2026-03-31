# section4 CoT Schema

This document defines the initial Cursor-on-Target mapping used by
section4.

section4 uses CoT as the wire format for exchanging operational events.
section4 itself remains the application protocol and system of record.

## Principles

- Keep CoT payloads terse.
- Put W3 fields in the CoT base schema.
- Put section4 workflow data under `detail/section4`.
- Use stable `uid` values for incident and job updates.
- Persist the full workflow state in SQLite, not in CoT.

## Base CoT Fields

Every section4 CoT event must include:

- `event.version`
- `event.type`
- `event.uid`
- `event.time`
- `event.start`
- `event.stale`
- `event.how`
- `point.lat`
- `point.lon`
- `point.hae`
- `point.ce`
- `point.le`

These follow the CoT base schema described in
[dev-guide-to-cot.txt](/Users/eax/repos/hackathon/docs/dev-guide-to-cot.txt).

## section4 Detail Extension

section4-specific data lives under:

```xml
<detail>
  <section4 ... />
</detail>
```

The `section4` element is the application extension point.

## Event Categories

section4 currently uses these categories:

- `incident`
- `job`
- `capability`

Each category is identified in `detail/section4` via the `object`
attribute.

## Incident Events

Current CoT type:

- `b-m-r`

Current section4 detail attributes:

- `object="incident"`
- `incident_id`
- `status`
- `part_number`
- `failed_component`
- `recommended_coa`

UID strategy:

- external UID if one exists
- otherwise `section4-incident-<incident_id>`

## Job Events

Current CoT type:

- `b-m-t`

Current section4 detail attributes:

- `object="job"`
- `job_id`
- `incident_id`
- `status`
- `job_type`
- `course_of_action`

UID strategy:

- `section4-job-<job_id>`

## Capability Events

Capability advertisement is the next transport item to add.

Current CoT type:

- `c-f-section4`

Current section4 detail attributes:

- `object="capability"`
- `capability_id`
- `node_id`
- `capability_type`
- `title`
- `availability_status`
- `throughput_per_day`
- `lead_time_minutes`

## Supporting Detail

When available, section4 should also include:

- `detail/contact@callsign`
- `detail/remarks`

Use these for operator readability, while keeping workflow semantics in
`detail/section4`.

## Persistence Rules

section4 stores:

- incidents
- jobs
- capabilities
- artifacts metadata
- audit events

CoT messages are projections of that state, not the source of truth.
