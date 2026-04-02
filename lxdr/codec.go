package lxdr

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func isOneOf[T comparable](value T, allowed ...T) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}

	return false
}

func ParseRequestPriorityCode(input string) (RequestPriorityCode, error) {
	code := RequestPriorityCode(input)
	if !code.IsValid() {
		return RequestPriorityCodeUnspecified, fmt.Errorf(
			"invalid request priority code: %q",
			input,
		)
	}

	return code, nil
}

func (c RequestPriorityCode) IsValid() bool {
	return isOneOf(
		c,
		RequestPriorityCode01,
		RequestPriorityCode02,
		RequestPriorityCode03,
		RequestPriorityCode04,
		RequestPriorityCode05,
		RequestPriorityCode06,
		RequestPriorityCode07,
		RequestPriorityCode08,
		RequestPriorityCode09,
		RequestPriorityCode10,
		RequestPriorityCode11,
		RequestPriorityCode12,
		RequestPriorityCode13,
		RequestPriorityCode14,
		RequestPriorityCode15,
	)
}

func (c CargoHMICCode) IsValid() bool {
	return isOneOf(
		c,
		CargoHMICCodeD,
		CargoHMICCodeN,
		CargoHMICCodeP,
		CargoHMICCodeY,
	)
}

func (c CargoHandlingCode) IsValid() bool {
	return isOneOf(
		c,
		CargoHandlingCodeC,
		CargoHandlingCodeM,
		CargoHandlingCodeT,
		CargoHandlingCodeR,
		CargoHandlingCodeX,
	)
}

func (c AttachmentIndicatorCode) IsValid() bool {
	return isOneOf(c, AttachmentIndicatorCodeNo, AttachmentIndicatorCodeYes)
}

func (c MaintenanceOperationalConditionCode) IsValid() bool {
	return isOneOf(
		c,
		MaintenanceOperationalConditionCodeA,
		MaintenanceOperationalConditionCodeB,
		MaintenanceOperationalConditionCodeC,
	)
}

func (c MaintenanceSupportTypeCode) IsValid() bool {
	return isOneOf(
		c,
		MaintenanceSupportTypeCodeXX,
		MaintenanceSupportTypeCodeR1,
		MaintenanceSupportTypeCodeR2,
		MaintenanceSupportTypeCodeR3,
		MaintenanceSupportTypeCodeR4,
	)
}

func (c MaintenanceRepairTypeCode) IsValid() bool {
	return isOneOf(
		c,
		MaintenanceRepairTypeCodeM1,
		MaintenanceRepairTypeCodeS1,
		MaintenanceRepairTypeCodeS2,
		MaintenanceRepairTypeCodeC1,
		MaintenanceRepairTypeCodeD1,
	)
}

func (c MaintenanceMajorDefectCode) IsValid() bool {
	return isOneOf(
		c,
		MaintenanceMajorDefectCodeMD01,
		MaintenanceMajorDefectCodeMD02,
		MaintenanceMajorDefectCodeMD03,
		MaintenanceMajorDefectCodeMD04,
		MaintenanceMajorDefectCodeMD05,
		MaintenanceMajorDefectCodeMD06,
		MaintenanceMajorDefectCodeMD07,
		MaintenanceMajorDefectCodeMD08,
		MaintenanceMajorDefectCodeMD09,
		MaintenanceMajorDefectCodeMD10,
		MaintenanceMajorDefectCodeMD11,
		MaintenanceMajorDefectCodeMD12,
		MaintenanceMajorDefectCodeMD13,
		MaintenanceMajorDefectCodeMD14,
		MaintenanceMajorDefectCodeMD15,
		MaintenanceMajorDefectCodeMD16,
		MaintenanceMajorDefectCodeNMAJ,
	)
}

func (c RoadClassificationCode) IsValid() bool {
	return isOneOf(
		c,
		RoadClassificationCodeA,
		RoadClassificationCodeB,
		RoadClassificationCodeC,
		RoadClassificationCodeD,
	)
}

func (c RoadDrainageCode) IsValid() bool {
	return isOneOf(c, RoadDrainageCodeA, RoadDrainageCodeB)
}

func (c RoadFoundationCode) IsValid() bool {
	return isOneOf(c, RoadFoundationCodeA, RoadFoundationCodeB)
}

func (c RoadSurfaceTypeCode) IsValid() bool {
	return isOneOf(
		c,
		RoadSurfaceTypeCodeA,
		RoadSurfaceTypeCodeB,
		RoadSurfaceTypeCodeC,
		RoadSurfaceTypeCodeD,
		RoadSurfaceTypeCodeE,
		RoadSurfaceTypeCodeF,
		RoadSurfaceTypeCodeG,
		RoadSurfaceTypeCodeH,
		RoadSurfaceTypeCodeI,
		RoadSurfaceTypeCodeJ,
	)
}

func (c EstimateCode) IsValid() bool {
	return isOneOf(c, EstimateCodeNo, EstimateCodeYes)
}

func (c LandingZoneLayoutDesignationCode) IsValid() bool {
	return isOneOf(
		c,
		LandingZoneLayoutDesignationCodeLZ,
		LandingZoneLayoutDesignationCodeLS,
		LandingZoneLayoutDesignationCodeLP,
	)
}

func (c AircraftSupportabilityCode) IsValid() bool {
	return isOneOf(
		c,
		AircraftSupportabilityCodeA,
		AircraftSupportabilityCodeB,
		AircraftSupportabilityCodeC,
		AircraftSupportabilityCodeD,
		AircraftSupportabilityCodeE,
		AircraftSupportabilityCodeF,
	)
}

func (c CardinalDirectionCode) IsValid() bool {
	return isOneOf(
		c,
		CardinalDirectionCodeS,
		CardinalDirectionCodeE,
		CardinalDirectionCodeN,
		CardinalDirectionCodeW,
	)
}

func (c LandingZoneSurfaceSlopeCode) IsValid() bool {
	return isOneOf(
		c,
		LandingZoneSurfaceSlopeCodeA,
		LandingZoneSurfaceSlopeCodeB,
	)
}

func (c LandingZoneObstacleCode) IsValid() bool {
	return isOneOf(
		c,
		LandingZoneObstacleCode1,
		LandingZoneObstacleCode2,
		LandingZoneObstacleCode3,
	)
}

func (c ObstacleActionCode) IsValid() bool {
	return isOneOf(
		c,
		ObstacleActionCode1,
		ObstacleActionCode2,
		ObstacleActionCode3,
	)
}

func (c ObstacleDeterminationCode) IsValid() bool {
	return isOneOf(
		c,
		ObstacleDeterminationCode1,
		ObstacleDeterminationCode2,
		ObstacleDeterminationCode3,
		ObstacleDeterminationCode4,
	)
}

func (c BypassCode) IsValid() bool {
	return isOneOf(c, BypassCodeNo, BypassCodeYes)
}

func (h RequestHeader) RenderCanonical() (string, error) {
	if err := h.Validate(); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s-%s-%s-%s-%s-%s-%02d",
		h.LocalSystemDate,
		h.LocalSystemTime,
		h.SynchronizedGeospatialReference,
		h.LocalRequestID,
		h.RequestPriority,
		h.ElementUnitIDOrCallsign,
		h.RequestSegmentCount,
	), nil
}

func ParseCanonicalHeader(input string) (RequestHeader, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 7 {
		return RequestHeader{}, fmt.Errorf(
			"invalid header part count: got %d want 7",
			len(parts),
		)
	}

	count, err := strconv.ParseUint(parts[6], 10, 8)
	if err != nil {
		return RequestHeader{}, fmt.Errorf("parse segment count: %w", err)
	}

	header := RequestHeader{
		LocalSystemDate:                 parts[0],
		LocalSystemTime:                 parts[1],
		SynchronizedGeospatialReference: parts[2],
		LocalRequestID:                  parts[3],
		ElementUnitIDOrCallsign:         parts[5],
		RequestSegmentCount:             uint8(count),
	}

	priority, err := ParseRequestPriorityCode(parts[4])
	if err != nil {
		return RequestHeader{}, err
	}
	header.RequestPriority = priority

	if err := header.Validate(); err != nil {
		return RequestHeader{}, err
	}

	return header, nil
}

func (r SynchronizedResponse) RenderCanonical() (string, error) {
	if r.LocalRequestID == "" || r.SynchronizedRequestID == "" {
		return "", errors.New("sync response requires both request IDs")
	}

	return r.LocalRequestID + "-" + r.SynchronizedRequestID, nil
}

func ParseSynchronizedResponse(input string) (SynchronizedResponse, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return SynchronizedResponse{}, fmt.Errorf(
			"invalid sync response part count: got %d want 2",
			len(parts),
		)
	}

	resp := SynchronizedResponse{
		LocalRequestID:        parts[0],
		SynchronizedRequestID: parts[1],
	}

	_, err := resp.RenderCanonical()
	if err != nil {
		return SynchronizedResponse{}, err
	}

	return resp, nil
}

func (s MobilityPaxRequestSegment) RenderCanonical() (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%d-%s-%s-%s-%s-%s-%s-%s-%s-%s",
		s.SegmentNumber,
		s.RequestTypeCode,
		s.RequestPriority,
		s.ZAPOrEDIPI,
		s.EarliestDepartureDateLocal,
		s.LatestDepartureDateLocal,
		s.DepartureLocation,
		s.DestinationLocation,
		s.TotalEstimatedBaggageLBS,
		s.HazardousMaterialType,
	), nil
}

func ParseCanonicalPAXSegment(
	input string,
) (MobilityPaxRequestSegment, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 10 {
		return MobilityPaxRequestSegment{}, fmt.Errorf(
			"invalid pax part count: got %d want 10",
			len(parts),
		)
	}

	segNum, err := strconv.ParseUint(parts[0], 10, 8)
	if err != nil {
		return MobilityPaxRequestSegment{}, fmt.Errorf(
			"parse pax segment number: %w",
			err,
		)
	}

	segment := MobilityPaxRequestSegment{
		SegmentNumber:              uint8(segNum),
		ZAPOrEDIPI:                 parts[3],
		EarliestDepartureDateLocal: parts[4],
		LatestDepartureDateLocal:   parts[5],
		DepartureLocation:          parts[6],
		DestinationLocation:        parts[7],
		TotalEstimatedBaggageLBS:   parts[8],
		HazardousMaterialType:      parts[9],
	}

	if parts[1] != string(MobilityPaxRequestTypeCodePM) {
		return MobilityPaxRequestSegment{}, fmt.Errorf(
			"invalid pax request type code: %q",
			parts[1],
		)
	}
	segment.RequestTypeCode = MobilityPaxRequestTypeCodePM

	priority, err := ParseRequestPriorityCode(parts[2])
	if err != nil {
		return MobilityPaxRequestSegment{}, err
	}
	segment.RequestPriority = priority

	if err := segment.Validate(); err != nil {
		return MobilityPaxRequestSegment{}, err
	}

	return segment, nil
}

func (s MobilityCargoRequestSegment) RenderCanonical() (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%d-%s-%s-%s-%s-%s-%s-%s-%s-%s-%s-%s-%s-%s-%s-%s",
		s.SegmentNumber,
		s.RequestTypeCode,
		s.RequestPriority,
		s.ItemByNIIN,
		s.ItemQuantity,
		s.SerialNumber,
		s.GrossWeightLBS,
		s.ActualHeightInches,
		s.ActualWidthInches,
		s.ActualLengthInches,
		s.HMIC,
		s.Handling,
		s.EarliestDepartureDateLocal,
		s.LatestDepartureDateLocal,
		s.DepartureLocation,
		s.DestinationLocation,
	), nil
}

func ParseCanonicalCargoSegment(
	input string,
) (MobilityCargoRequestSegment, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 16 {
		return MobilityCargoRequestSegment{}, fmt.Errorf(
			"invalid cargo part count: got %d want 16",
			len(parts),
		)
	}

	segNum, err := strconv.ParseUint(parts[0], 10, 8)
	if err != nil {
		return MobilityCargoRequestSegment{}, fmt.Errorf(
			"parse cargo segment number: %w",
			err,
		)
	}

	segment := MobilityCargoRequestSegment{
		SegmentNumber:              uint8(segNum),
		ItemByNIIN:                 parts[3],
		ItemQuantity:               parts[4],
		SerialNumber:               parts[5],
		GrossWeightLBS:             parts[6],
		ActualHeightInches:         parts[7],
		ActualWidthInches:          parts[8],
		ActualLengthInches:         parts[9],
		EarliestDepartureDateLocal: parts[12],
		LatestDepartureDateLocal:   parts[13],
		DepartureLocation:          parts[14],
		DestinationLocation:        parts[15],
	}

	if parts[1] != string(MobilityCargoRequestTypeCodeCM) {
		return MobilityCargoRequestSegment{}, fmt.Errorf(
			"invalid cargo request type code: %q",
			parts[1],
		)
	}
	segment.RequestTypeCode = MobilityCargoRequestTypeCodeCM

	priority, err := ParseRequestPriorityCode(parts[2])
	if err != nil {
		return MobilityCargoRequestSegment{}, err
	}
	segment.RequestPriority = priority
	segment.HMIC = CargoHMICCode(parts[10])
	segment.Handling = CargoHandlingCode(parts[11])

	if err := segment.Validate(); err != nil {
		return MobilityCargoRequestSegment{}, err
	}

	return segment, nil
}

func (h RequestHeader) Validate() error {
	if h.LocalSystemDate == "" {
		return errors.New("header local system date is required")
	}
	if h.LocalSystemTime == "" {
		return errors.New("header local system time is required")
	}
	if h.SynchronizedGeospatialReference == "" {
		return errors.New("header geospatial reference is required")
	}
	if h.LocalRequestID == "" {
		return errors.New("header local request ID is required")
	}
	if !h.RequestPriority.IsValid() {
		return errors.New("header request priority is required")
	}
	if h.ElementUnitIDOrCallsign == "" {
		return errors.New("header element/unit identifier is required")
	}
	if h.RequestSegmentCount == 0 {
		return errors.New("header request segment count must be > 0")
	}
	return nil
}

func (c RequestContainer) Validate() error {
	if err := c.Header.Validate(); err != nil {
		return err
	}
	if len(c.Segments) == 0 {
		return errors.New("request container requires at least one segment")
	}
	if len(c.Segments) != int(c.Header.RequestSegmentCount) {
		return fmt.Errorf(
			"segment count mismatch: header=%d actual=%d",
			c.Header.RequestSegmentCount,
			len(c.Segments),
		)
	}
	for i, segment := range c.Segments {
		if err := segment.Validate(); err != nil {
			return fmt.Errorf("segment %d: %w", i, err)
		}
	}
	return nil
}

func (s RequestSegment) Validate() error {
	count := 0
	if s.MobilityPax != nil {
		count++
		if err := s.MobilityPax.Validate(); err != nil {
			return err
		}
	}
	if s.MobilityCargo != nil {
		count++
		if err := s.MobilityCargo.Validate(); err != nil {
			return err
		}
	}
	if s.Supply != nil {
		count++
		if err := s.Supply.Validate(); err != nil {
			return err
		}
	}
	if s.Maintenance != nil {
		count++
		if err := s.Maintenance.Validate(); err != nil {
			return err
		}
	}
	if s.EngineerRoad != nil {
		count++
		if err := s.EngineerRoad.Validate(); err != nil {
			return err
		}
	}
	if s.EngineerLZ != nil {
		count++
		if err := s.EngineerLZ.Validate(); err != nil {
			return err
		}
	}
	if s.ObstacleRemoval != nil {
		count++
		if err := s.ObstacleRemoval.Validate(); err != nil {
			return err
		}
	}
	if count != 1 {
		return errors.New("request segment must set exactly one segment body")
	}
	return nil
}

func (s MobilityPaxRequestSegment) Validate() error {
	if s.SegmentNumber == 0 {
		return errors.New("pax segment number is required")
	}
	if s.RequestTypeCode != MobilityPaxRequestTypeCodePM {
		return errors.New("pax request type code must be PM")
	}
	if !s.RequestPriority.IsValid() {
		return errors.New("pax request priority is required")
	}
	required := []struct {
		name  string
		value string
	}{
		{"zap or edi-pi", s.ZAPOrEDIPI},
		{"earliest departure date", s.EarliestDepartureDateLocal},
		{"latest departure date", s.LatestDepartureDateLocal},
		{"departure location", s.DepartureLocation},
		{"destination location", s.DestinationLocation},
		{"estimated baggage weight", s.TotalEstimatedBaggageLBS},
		{"hazardous material type", s.HazardousMaterialType},
	}
	return validateRequiredFields(required)
}

func (s MobilityCargoRequestSegment) Validate() error {
	if s.SegmentNumber == 0 {
		return errors.New("cargo segment number is required")
	}
	if s.RequestTypeCode != MobilityCargoRequestTypeCodeCM {
		return errors.New("cargo request type code must be CM")
	}
	if !s.RequestPriority.IsValid() {
		return errors.New("cargo request priority is required")
	}
	if !s.HMIC.IsValid() {
		return errors.New("hmic is required")
	}
	if !s.Handling.IsValid() {
		return errors.New("handling is required")
	}
	required := []struct {
		name  string
		value string
	}{
		{"item by niin", s.ItemByNIIN},
		{"item quantity", s.ItemQuantity},
		{"serial number", s.SerialNumber},
		{"gross weight", s.GrossWeightLBS},
		{"actual height", s.ActualHeightInches},
		{"actual width", s.ActualWidthInches},
		{"actual length", s.ActualLengthInches},
		{"earliest departure date", s.EarliestDepartureDateLocal},
		{"latest departure date", s.LatestDepartureDateLocal},
		{"departure location", s.DepartureLocation},
		{"destination location", s.DestinationLocation},
	}
	return validateRequiredFields(required)
}

func (s SupplyRequestSegment) Validate() error {
	if s.SegmentNumber == 0 {
		return errors.New("supply segment number is required")
	}
	if s.RequestTypeCode != SupplyRequestTypeCodeSR {
		return errors.New("supply request type code must be SR")
	}
	if !s.RequestPriority.IsValid() {
		return errors.New("supply request priority is required")
	}
	if s.AttachmentIndicator != AttachmentIndicatorCodeUnspecified &&
		!s.AttachmentIndicator.IsValid() {
		return errors.New("supply attachment indicator must be 0 or 1")
	}
	required := []struct {
		name  string
		value string
	}{
		{"item quantity", s.ItemQuantity},
		{"required delivery date", s.RequiredDeliveryDateLocal},
		{"delivery location", s.DeliveryLocation},
	}
	return validateRequiredFields(required)
}

func (s MaintenanceRequestSegment) Validate() error {
	if s.SegmentNumber == 0 {
		return errors.New("maintenance segment number is required")
	}
	if s.RequestTypeCode != MaintenanceRequestTypeCodeCM {
		return errors.New("maintenance request type code must be CM")
	}
	if !s.RequestPriority.IsValid() {
		return errors.New("maintenance request priority is required")
	}
	if !s.EquipmentOperationalCondition.IsValid() {
		return errors.New("equipment operational condition is required")
	}
	if !s.TypeOfMaintenanceSupportRequired.IsValid() {
		return errors.New("type of maintenance support is required")
	}
	if !s.TypeOfRepair.IsValid() {
		return errors.New("type of repair is required")
	}
	if !s.RepairMajorDefect.IsValid() {
		return errors.New("repair major defect is required")
	}
	if !s.AttachmentIndicator.IsValid() {
		return errors.New("attachment indicator must be 0 or 1")
	}
	required := []struct {
		name  string
		value string
	}{
		{"niin", s.NIIN},
		{"number of pieces", s.NumberOfPiecesRequiringSupport},
		{"maintenance support date", s.DateMaintenanceSupportRequiredLocal},
		{"location of equipment", s.LocationOfEquipment},
	}
	return validateRequiredFields(required)
}

func (s EngineerReconRoadReportSegment) Validate() error {
	if s.DateOfEvaluationLocal == "" {
		return errors.New("road evaluation date is required")
	}
	if s.StartPointLocation == "" {
		return errors.New("road start point location is required")
	}
	if s.EndPointLocation == "" {
		return errors.New("road end point location is required")
	}
	if !s.RoadClassification.IsValid() {
		return errors.New("road classification is required")
	}
	if !s.Drainage.IsValid() {
		return errors.New("road drainage is required")
	}
	if !s.Foundation.IsValid() {
		return errors.New("road foundation is required")
	}
	if !s.SurfaceType.IsValid() {
		return errors.New("road surface type is required")
	}
	if s.Obstructions == "" {
		return errors.New("road obstructions are required")
	}
	if s.AttachmentIndicator != AttachmentIndicatorCodeUnspecified &&
		!s.AttachmentIndicator.IsValid() {
		return errors.New("road attachment indicator must be 0 or 1")
	}
	return nil
}

func (s EngineerReconLandingZoneReportSegment) Validate() error {
	if s.DateOfEvaluationLocal == "" {
		return errors.New("landing zone evaluation date is required")
	}
	if s.Location == "" {
		return errors.New("landing zone location is required")
	}
	if !s.Estimate.IsValid() {
		return errors.New("landing zone estimate is required")
	}
	if !s.LayoutDesignation.IsValid() {
		return errors.New("landing zone layout designation is required")
	}
	if s.LandingPointCapacity == 0 {
		return errors.New("landing point capacity is required")
	}
	if s.LandingZoneCapacity == 0 {
		return errors.New("landing zone capacity is required")
	}
	if s.LandingSiteCapacity == 0 {
		return errors.New("landing site capacity is required")
	}
	if s.LandingZoneWidthFeet == "" {
		return errors.New("landing zone width is required")
	}
	if s.LandingZoneLengthFeet == "" {
		return errors.New("landing zone length is required")
	}
	if !s.AircraftSupportability.IsValid() {
		return errors.New("aircraft supportability is required")
	}
	if !s.LandingZoneApproach.IsValid() {
		return errors.New("landing zone approach is required")
	}
	if !s.LandingZoneDeparture.IsValid() {
		return errors.New("landing zone departure is required")
	}
	if !s.LandingZoneSurfaceSlope.IsValid() {
		return errors.New("landing zone surface slope is required")
	}
	if !s.Obstacle.IsValid() {
		return errors.New("landing zone obstacle is required")
	}
	if s.AttachmentIndicator != AttachmentIndicatorCodeUnspecified &&
		!s.AttachmentIndicator.IsValid() {
		return errors.New(
			"landing zone attachment indicator must be 0 or 1",
		)
	}
	return nil
}

func (s GeneralEngineeringObstacleRemovalSegment) Validate() error {
	if s.DateOfEvaluationLocal == "" {
		return errors.New("obstacle removal evaluation date is required")
	}
	if s.Location == "" {
		return errors.New("obstacle removal location is required")
	}
	if !s.Obstacle.IsValid() {
		return errors.New("obstacle action is required")
	}
	if s.MinMaxLengthFeet == "" {
		return errors.New("obstacle min/max length is required")
	}
	if s.MinMaxWidthFeet == "" {
		return errors.New("obstacle min/max width is required")
	}
	if s.MinMaxDepthFeet == "" {
		return errors.New("obstacle min/max depth is required")
	}
	if s.RouteNumber == "" {
		return errors.New("obstacle route number is required")
	}
	if !s.DeterminationOfAction.IsValid() {
		return errors.New("obstacle determination of action is required")
	}
	if !s.Bypass.IsValid() {
		return errors.New("obstacle bypass is required")
	}
	if s.BypassGrid == "" {
		return errors.New("obstacle bypass grid is required")
	}
	if s.AttachmentIndicator != AttachmentIndicatorCodeUnspecified &&
		!s.AttachmentIndicator.IsValid() {
		return errors.New(
			"obstacle attachment indicator must be 0 or 1",
		)
	}
	return nil
}

func validateRequiredFields(
	fields []struct {
		name  string
		value string
	},
) error {
	for _, field := range fields {
		if field.value == "" {
			return fmt.Errorf("%s is required", field.name)
		}
	}
	return nil
}
