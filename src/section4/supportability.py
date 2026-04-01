"""Supportability scoring for section4 logistics requests."""

from __future__ import annotations

from dataclasses import dataclass

from sqlalchemy import select
from sqlalchemy.orm import Session

from section4.storage.tables import (
    Capability,
    LogisticsStatusReport,
    LXDRRequestRecord,
)


@dataclass(frozen=True, slots=True)
class SupportabilityAssessment:
    """Compact supportability result for one persisted ADRIAN request."""

    supportable: bool
    recommended_course_of_action: str
    supporting_node_id: str | None
    rationale: str


def assess_request_supportability(
    session: Session,
    request: LXDRRequestRecord,
) -> SupportabilityAssessment:
    """Assess supportability using stock and capability data."""

    stock_rows = session.scalars(
        select(LogisticsStatusReport).order_by(
            LogisticsStatusReport.reported_at.desc()
        )
    ).all()
    best_stock = next(
        (
            row
            for row in stock_rows
            if (row.on_hand_quantity or 0) >= max(row.required_quantity or 0, 1)
        ),
        None,
    )
    if best_stock is not None:
        return SupportabilityAssessment(
            supportable=True,
            recommended_course_of_action="issue_from_stock",
            supporting_node_id=best_stock.node_id,
            rationale=(
                "On-hand stock is sufficient to satisfy the request and "
                "supports direct issue."
            ),
        )

    fabrication = session.scalar(
        select(Capability)
        .where(
            Capability.capability_type == "fabrication",
            Capability.availability_status == "available",
            Capability.active.is_(True),
        )
        .order_by(Capability.lead_time_minutes.asc())
        .limit(1)
    )
    if fabrication is not None:
        return SupportabilityAssessment(
            supportable=True,
            recommended_course_of_action="fabricate",
            supporting_node_id=fabrication.node_id,
            rationale=(
                "No sufficient stock is on hand, but a fabrication "
                "capability is available."
            ),
        )

    return SupportabilityAssessment(
        supportable=False,
        recommended_course_of_action="blocked",
        supporting_node_id=None,
        rationale=(
            "No sufficient stock and no available supporting capability "
            "were found."
        ),
    )
