# ForgeNet CoT Schema

This document defines the initial Cursor-on-Target mapping used by
ForgeNet.

ForgeNet uses CoT as the wire format for exchanging operational events.
ForgeNet itself remains the application protocol and system of record.

## Principles

- Keep CoT payloads terse.
- Put W3 fields in the CoT base schema.
- Put ForgeNet workflow data under `detail/forgenet`.
- Use stable `uid` values for incident and job updates.
- Persist the full workflow state in SQLite, not in CoT.

## Base CoT Fields

Every ForgeNet CoT event must include:

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

## ForgeNet Detail Extension

ForgeNet-specific data lives under:

```xml
<detail>
  <forgenet ... />
</detail>
```

The `forgenet` element is the application extension point.

## Event Categories

ForgeNet currently uses these categories:

- `incident`
- `job`
- `capability`

Each category is identified in `detail/forgenet` via the `object`
attribute.

## Incident Events

Current CoT type:

- `b-m-r`

Current ForgeNet detail attributes:

- `object="incident"`
- `incident_id`
- `status`
- `part_number`
- `failed_component`
- `recommended_coa`

UID strategy:

- external UID if one exists
- otherwise `forgenet-incident-<incident_id>`

## Job Events

Current CoT type:

- `b-m-t`

Current ForgeNet detail attributes:

- `object="job"`
- `job_id`
- `incident_id`
- `status`
- `job_type`
- `course_of_action`

UID strategy:

- `forgenet-job-<job_id>`

## Capability Events

Capability advertisement is the next transport item to add.

Current CoT type:

- `c-f-forgenet`

Current ForgeNet detail attributes:

- `object="capability"`
- `capability_id`
- `node_id`
- `capability_type`
- `title`
- `availability_status`
- `throughput_per_day`
- `lead_time_minutes`

## Supporting Detail

When available, ForgeNet should also include:

- `detail/contact@callsign`
- `detail/remarks`

Use these for operator readability, while keeping workflow semantics in
`detail/forgenet`.

## Persistence Rules

ForgeNet stores:

- incidents
- jobs
- capabilities
- artifacts metadata
- audit events

CoT messages are projections of that state, not the source of truth.
