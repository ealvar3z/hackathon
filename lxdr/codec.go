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

type validationBranch struct {
	present  bool
	validate func() error
}

func validateExactlyOneBranch(name string, branches ...validationBranch) error {
	count := 0
	for _, branch := range branches {
		if !branch.present {
			continue
		}
		count++
		if branch.validate != nil {
			if err := branch.validate(); err != nil {
				return err
			}
		}
	}
	if count != 1 {
		return fmt.Errorf("%s must set exactly one payload", name)
	}
	return nil
}

func ParseRequestPriorityCode(input string) (RequestPriorityCode, error) {
	switch input {
	case "01":
		return RequestPriorityCode01, nil
	case "02":
		return RequestPriorityCode02, nil
	case "03":
		return RequestPriorityCode03, nil
	case "04":
		return RequestPriorityCode04, nil
	case "05":
		return RequestPriorityCode05, nil
	case "06":
		return RequestPriorityCode06, nil
	case "07":
		return RequestPriorityCode07, nil
	case "08":
		return RequestPriorityCode08, nil
	case "09":
		return RequestPriorityCode09, nil
	case "10":
		return RequestPriorityCode10, nil
	case "11":
		return RequestPriorityCode11, nil
	case "12":
		return RequestPriorityCode12, nil
	case "13":
		return RequestPriorityCode13, nil
	case "14":
		return RequestPriorityCode14, nil
	case "15":
		return RequestPriorityCode15, nil
	default:
		return RequestPriorityCodeUnspecified, fmt.Errorf(
			"invalid request priority code: %q",
			input,
		)
	}
}

func formatRequestPriorityCode(code RequestPriorityCode) string {
	switch code {
	case RequestPriorityCode01:
		return "01"
	case RequestPriorityCode02:
		return "02"
	case RequestPriorityCode03:
		return "03"
	case RequestPriorityCode04:
		return "04"
	case RequestPriorityCode05:
		return "05"
	case RequestPriorityCode06:
		return "06"
	case RequestPriorityCode07:
		return "07"
	case RequestPriorityCode08:
		return "08"
	case RequestPriorityCode09:
		return "09"
	case RequestPriorityCode10:
		return "10"
	case RequestPriorityCode11:
		return "11"
	case RequestPriorityCode12:
		return "12"
	case RequestPriorityCode13:
		return "13"
	case RequestPriorityCode14:
		return "14"
	case RequestPriorityCode15:
		return "15"
	default:
		return ""
	}
}

func formatMobilityPaxRequestTypeCode(code MobilityPaxRequestTypeCode) string {
	if code == MobilityPaxRequestTypeCodePM {
		return "PM"
	}
	return ""
}

func formatMobilityCargoRequestTypeCode(
	code MobilityCargoRequestTypeCode,
) string {
	if code == MobilityCargoRequestTypeCodeCM {
		return "CM"
	}
	return ""
}

func formatCargoHMICCode(code CargoHMICCode) string {
	switch code {
	case CargoHMICCodeD:
		return "D"
	case CargoHMICCodeN:
		return "N"
	case CargoHMICCodeP:
		return "P"
	case CargoHMICCodeY:
		return "Y"
	default:
		return ""
	}
}

func formatCargoHandlingCode(code CargoHandlingCode) string {
	switch code {
	case CargoHandlingCodeC:
		return "C"
	case CargoHandlingCodeM:
		return "M"
	case CargoHandlingCodeT:
		return "T"
	case CargoHandlingCodeR:
		return "R"
	case CargoHandlingCodeX:
		return "X"
	default:
		return ""
	}
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

func (c CanonicalFileType) IsValid() bool {
	return isOneOf(
		c,
		CanonicalFileTypeRequestHeader,
		CanonicalFileTypeRequestSegment,
		CanonicalFileTypeRequestContainer,
		CanonicalFileTypeSyncResponse,
		CanonicalFileTypeAttachmentOrMediaReference,
	)
}

func (c ExchangeRole) IsValid() bool {
	return isOneOf(
		c,
		ExchangeRoleTransmitted,
		ExchangeRoleCalculated,
		ExchangeRoleSynchronized,
		ExchangeRoleLocalOnly,
	)
}

func (m LinkDeliveryMethod) IsValid() bool {
	return isOneOf(
		m,
		LinkDeliveryMethodDirect,
		LinkDeliveryMethodPropagated,
		LinkDeliveryMethodOpportunistic,
	)
}

func (r LinkRepresentation) IsValid() bool {
	return isOneOf(r, LinkRepresentationBinaryProto)
}

func (r LinkRepresentation) IsV1Supported() bool {
	return r == LinkRepresentationBinaryProto
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

func (c CBRNAgentTypeCode) IsValid() bool {
	return isOneOf(
		c,
		CBRNAgentTypeCode1,
		CBRNAgentTypeCode2,
		CBRNAgentTypeCode3,
		CBRNAgentTypeCode4,
		CBRNAgentTypeCode5,
	)
}

func (c CBRNPhysicalPropertyCode) IsValid() bool {
	return isOneOf(
		c,
		CBRNPhysicalPropertyCode1,
		CBRNPhysicalPropertyCode2,
		CBRNPhysicalPropertyCode3,
		CBRNPhysicalPropertyCode4,
	)
}

func (c CBRNContaminationValueCode) IsValid() bool {
	return isOneOf(
		c,
		CBRNContaminationValueCodeE,
		CBRNContaminationValueCodeI,
		CBRNContaminationValueCodeW,
		CBRNContaminationValueCodeC,
	)
}

func (c MunitionPurposeCode) IsValid() bool {
	return isOneOf(
		c,
		MunitionPurposeCodeAA,
		MunitionPurposeCodeAP,
		MunitionPurposeCodeFL,
		MunitionPurposeCodeSM,
		MunitionPurposeCodeIM,
	)
}

func (c MunitionTypeCode) IsValid() bool {
	return isOneOf(
		c,
		MunitionTypeCodeE,
		MunitionTypeCodeD,
		MunitionTypeCodeT,
		MunitionTypeCodeP,
	)
}

func (c HealthCollectionRequestTypeCode) IsValid() bool {
	return isOneOf(c, HealthCollectionRequestTypeCodeCR)
}

func (c HealthPrimaryMechanismCode) IsValid() bool {
	return isOneOf(
		c,
		HealthPrimaryMechanismCodeE1,
		HealthPrimaryMechanismCodeE2,
		HealthPrimaryMechanismCodeE3,
		HealthPrimaryMechanismCodeE4,
		HealthPrimaryMechanismCodeE5,
		HealthPrimaryMechanismCodeP1,
		HealthPrimaryMechanismCodeP2,
		HealthPrimaryMechanismCodeP3,
		HealthPrimaryMechanismCodeP4,
		HealthPrimaryMechanismCodeP5,
		HealthPrimaryMechanismCodeD1,
		HealthPrimaryMechanismCodeD2,
		HealthPrimaryMechanismCodeD3,
		HealthPrimaryMechanismCodeD4,
	)
}

func (c HealthCBRNExposureCode) IsValid() bool {
	return isOneOf(
		c,
		HealthCBRNExposureCodeC,
		HealthCBRNExposureCodeB,
		HealthCBRNExposureCodeR,
		HealthCBRNExposureCodeN,
		HealthCBRNExposureCodeX,
	)
}

func (c HealthMajorSignsSymptomsCode) IsValid() bool {
	return isOneOf(
		c,
		HealthMajorSignsSymptomsCodeB,
		HealthMajorSignsSymptomsCodeR,
		HealthMajorSignsSymptomsCodeX,
		HealthMajorSignsSymptomsCodeC,
	)
}

func (c ServiceCode) IsValid() bool {
	return isOneOf(
		c,
		ServiceCodeUSA,
		ServiceCodeUSSF,
		ServiceCodeUSAF,
		ServiceCodeUSCG,
		ServiceCodeUSN,
		ServiceCodeUSMC,
		ServiceCodeUSCIV,
		ServiceCodeNONUS,
		ServiceCodeEPW,
	)
}

func (c HealthPulseLocationCode) IsValid() bool {
	return isOneOf(c, HealthPulseLocationCodeW, HealthPulseLocationCodeN)
}

func (c HealthResponsivenessCode) IsValid() bool {
	return isOneOf(
		c,
		HealthResponsivenessCodeA,
		HealthResponsivenessCodeV,
		HealthResponsivenessCodeP,
		HealthResponsivenessCodeU,
	)
}

func (c HealthTriagePrecedenceCode) IsValid() bool {
	return isOneOf(
		c,
		HealthTriagePrecedenceCodeA,
		HealthTriagePrecedenceCodeB,
		HealthTriagePrecedenceCodeC,
		HealthTriagePrecedenceCodeD,
		HealthTriagePrecedenceCodeE,
	)
}

func (c TourniquetPlacementCode) IsValid() bool {
	return isOneOf(
		c,
		TourniquetPlacementCodeTQXX,
		TourniquetPlacementCodeTQRA,
		TourniquetPlacementCodeTQLA,
		TourniquetPlacementCodeTQRL,
		TourniquetPlacementCodeTQLL,
	)
}

func (c TourniquetTypeCode) IsValid() bool {
	return isOneOf(c, TourniquetTypeCodeE, TourniquetTypeCodeJ, TourniquetTypeCodeT)
}

func (c WoundTreatmentCode) IsValid() bool {
	return isOneOf(
		c,
		WoundTreatmentCodeT1,
		WoundTreatmentCodeT2,
		WoundTreatmentCodeT3,
		WoundTreatmentCodeT4,
		WoundTreatmentCodeT5,
		WoundTreatmentCodeT6,
		WoundTreatmentCodeT7,
	)
}

func (c AirwayTreatmentCode) IsValid() bool {
	return isOneOf(
		c,
		AirwayTreatmentCodeA0,
		AirwayTreatmentCodeA1,
		AirwayTreatmentCodeA2,
		AirwayTreatmentCodeA3,
		AirwayTreatmentCodeA4,
	)
}

func (c BreathingTreatmentCode) IsValid() bool {
	return isOneOf(
		c,
		BreathingTreatmentCodeB0,
		BreathingTreatmentCodeB1,
		BreathingTreatmentCodeB3,
		BreathingTreatmentCodeB4,
		BreathingTreatmentCodeB5,
	)
}

func (c FluidNameCode) IsValid() bool {
	return isOneOf(c, FluidNameCodeS, FluidNameCodeR, FluidNameCodeH)
}

func (c FluidRouteCode) IsValid() bool {
	return isOneOf(c, FluidRouteCodeIV, FluidRouteCodeIO)
}

func (c BloodProductCode) IsValid() bool {
	return isOneOf(c, BloodProductCodeWBD, BloodProductCodeRBC, BloodProductCodeFFP, BloodProductCodeFDP)
}

func (c MedicationRouteCode) IsValid() bool {
	return isOneOf(
		c,
		MedicationRouteCodeR1,
		MedicationRouteCodeR2,
		MedicationRouteCodeR3,
		MedicationRouteCodeR4,
		MedicationRouteCodeR5,
		MedicationRouteCodeR6,
		MedicationRouteCodeR7,
	)
}

func (c AnalgesicMedicationCode) IsValid() bool {
	return isOneOf(c, AnalgesicMedicationCodeK, AnalgesicMedicationCodeF, AnalgesicMedicationCodeM)
}

func (c AntibioticMedicationCode) IsValid() bool {
	return isOneOf(c, AntibioticMedicationCodeM, AntibioticMedicationCodeE, AntibioticMedicationCodeP, AntibioticMedicationCodeA)
}

func (c OtherMedicationCode) IsValid() bool {
	return isOneOf(c, OtherMedicationCodeI, OtherMedicationCodeT)
}

func (c CasualtyTypeCode) IsValid() bool {
	return isOneOf(c, CasualtyTypeCodeA, CasualtyTypeCodeB)
}

func (c HealthEvacuationRequestPriorityCode) IsValid() bool {
	return isOneOf(
		c,
		HealthEvacuationRequestPriorityCodeA,
		HealthEvacuationRequestPriorityCodeB,
		HealthEvacuationRequestPriorityCodeC,
		HealthEvacuationRequestPriorityCodeD,
		HealthEvacuationRequestPriorityCodeE,
	)
}

func (c HealthEvacuationLocationMarkingCode) IsValid() bool {
	return isOneOf(
		c,
		HealthEvacuationLocationMarkingCodeA,
		HealthEvacuationLocationMarkingCodeB,
		HealthEvacuationLocationMarkingCodeC,
		HealthEvacuationLocationMarkingCodeD,
		HealthEvacuationLocationMarkingCodeE,
	)
}

func (c HealthEvacuationContaminationCode) IsValid() bool {
	return isOneOf(
		c,
		HealthEvacuationContaminationCodeA,
		HealthEvacuationContaminationCodeB,
		HealthEvacuationContaminationCodeC,
		HealthEvacuationContaminationCodeD,
	)
}

func (c HealthEvacuationCasualtyTypeCode) IsValid() bool {
	return isOneOf(
		c,
		HealthEvacuationCasualtyTypeCodeA,
		HealthEvacuationCasualtyTypeCodeL,
	)
}

func (c HealthEvacuationRequestedEquipmentCode) IsValid() bool {
	return isOneOf(
		c,
		HealthEvacuationRequestedEquipmentCodeA,
		HealthEvacuationRequestedEquipmentCodeB,
		HealthEvacuationRequestedEquipmentCodeC,
		HealthEvacuationRequestedEquipmentCodeD,
	)
}

func (c HealthEvacuationSecurityCode) IsValid() bool {
	return isOneOf(
		c,
		HealthEvacuationSecurityCodeN,
		HealthEvacuationSecurityCodeP,
		HealthEvacuationSecurityCodeE,
		HealthEvacuationSecurityCodeX,
	)
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
		h.LocalRequestId,
		formatRequestPriorityCode(h.RequestPriority),
		h.ElementUnitIdOrCallsign,
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
		LocalRequestId:                  parts[3],
		ElementUnitIdOrCallsign:         parts[5],
		RequestSegmentCount:             uint32(count),
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
	if r.LocalRequestId == "" || r.SynchronizedRequestId == "" {
		return "", errors.New("sync response requires both request IDs")
	}

	return r.LocalRequestId + "-" + r.SynchronizedRequestId, nil
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
		LocalRequestId:        parts[0],
		SynchronizedRequestId: parts[1],
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
		formatMobilityPaxRequestTypeCode(s.RequestTypeCode),
		formatRequestPriorityCode(s.RequestPriority),
		s.ZapOrEdiPi,
		s.EarliestDepartureDateLocal,
		s.LatestDepartureDateLocal,
		s.DepartureLocation,
		s.DestinationLocation,
		s.TotalEstimatedBaggageWeightLbs,
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
		SegmentNumber:                  uint32(segNum),
		ZapOrEdiPi:                     parts[3],
		EarliestDepartureDateLocal:     parts[4],
		LatestDepartureDateLocal:       parts[5],
		DepartureLocation:              parts[6],
		DestinationLocation:            parts[7],
		TotalEstimatedBaggageWeightLbs: parts[8],
		HazardousMaterialType:          parts[9],
	}

	if parts[1] != "PM" {
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
		formatMobilityCargoRequestTypeCode(s.RequestTypeCode),
		formatRequestPriorityCode(s.RequestPriority),
		s.ItemByNiin,
		s.ItemQuantity,
		s.SerialNumber,
		s.GrossWeightLbs,
		s.ActualHeightInches,
		s.ActualWidthInches,
		s.ActualLengthInches,
		formatCargoHMICCode(s.Hmic),
		formatCargoHandlingCode(s.Handling),
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
		SegmentNumber:              uint32(segNum),
		ItemByNiin:                 parts[3],
		ItemQuantity:               parts[4],
		SerialNumber:               parts[5],
		GrossWeightLbs:             parts[6],
		ActualHeightInches:         parts[7],
		ActualWidthInches:          parts[8],
		ActualLengthInches:         parts[9],
		EarliestDepartureDateLocal: parts[12],
		LatestDepartureDateLocal:   parts[13],
		DepartureLocation:          parts[14],
		DestinationLocation:        parts[15],
	}

	if parts[1] != "CM" {
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
	switch parts[10] {
	case "D":
		segment.Hmic = CargoHMICCodeD
	case "N":
		segment.Hmic = CargoHMICCodeN
	case "P":
		segment.Hmic = CargoHMICCodeP
	case "Y":
		segment.Hmic = CargoHMICCodeY
	default:
		return MobilityCargoRequestSegment{}, fmt.Errorf(
			"invalid cargo hmic code: %q",
			parts[10],
		)
	}
	switch parts[11] {
	case "C":
		segment.Handling = CargoHandlingCodeC
	case "M":
		segment.Handling = CargoHandlingCodeM
	case "T":
		segment.Handling = CargoHandlingCodeT
	case "R":
		segment.Handling = CargoHandlingCodeR
	case "X":
		segment.Handling = CargoHandlingCodeX
	default:
		return MobilityCargoRequestSegment{}, fmt.Errorf(
			"invalid cargo handling code: %q",
			parts[11],
		)
	}

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
	if h.LocalRequestId == "" {
		return errors.New("header local request ID is required")
	}
	if !h.RequestPriority.IsValid() {
		return errors.New("header request priority is required")
	}
	if h.ElementUnitIdOrCallsign == "" {
		return errors.New("header element/unit identifier is required")
	}
	if h.RequestSegmentCount == 0 {
		return errors.New("header request segment count must be > 0")
	}
	return nil
}

func (r SynchronizedResponse) Validate() error {
	if r.LocalRequestId == "" {
		return errors.New("synchronized response local request ID is required")
	}
	if r.SynchronizedRequestId == "" {
		return errors.New("synchronized response synchronized request ID is required")
	}
	if _, err := r.RenderCanonical(); err != nil {
		return err
	}
	return nil
}

func (c RequestContainer) Validate() error {
	if c.Header == nil {
		return errors.New("request container header is required")
	}
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
	return validateExactlyOneBranch(
		"request segment",
		validationBranch{present: s.GetMobilityPax() != nil, validate: func() error { return s.GetMobilityPax().Validate() }},
		validationBranch{present: s.GetMobilityCargo() != nil, validate: func() error { return s.GetMobilityCargo().Validate() }},
		validationBranch{present: s.GetSupply() != nil, validate: func() error { return s.GetSupply().Validate() }},
		validationBranch{present: s.GetMaintenance() != nil, validate: func() error { return s.GetMaintenance().Validate() }},
		validationBranch{present: s.GetEngineerReconArea() != nil, validate: func() error { return s.GetEngineerReconArea().Validate() }},
		validationBranch{present: s.GetEngineerReconZone() != nil, validate: func() error { return s.GetEngineerReconZone().Validate() }},
		validationBranch{present: s.GetEngineerReconRoute() != nil, validate: func() error { return s.GetEngineerReconRoute().Validate() }},
		validationBranch{present: s.GetEngineerReconRoad() != nil, validate: func() error { return s.GetEngineerReconRoad().Validate() }},
		validationBranch{present: s.GetEngineerReconLandingZone() != nil, validate: func() error { return s.GetEngineerReconLandingZone().Validate() }},
		validationBranch{present: s.GetObstacleRemoval() != nil, validate: func() error { return s.GetObstacleRemoval().Validate() }},
		validationBranch{present: s.GetEod() != nil, validate: func() error { return s.GetEod().Validate() }},
		validationBranch{present: s.GetBulkLiquidSupport() != nil, validate: func() error { return s.GetBulkLiquidSupport().Validate() }},
		validationBranch{present: s.GetDemolition() != nil, validate: func() error { return s.GetDemolition().Validate() }},
		validationBranch{present: s.GetHealthCollection() != nil, validate: func() error { return s.GetHealthCollection().Validate() }},
		validationBranch{present: s.GetHealthTriage() != nil, validate: func() error { return s.GetHealthTriage().Validate() }},
		validationBranch{present: s.GetHealthIntervention() != nil, validate: func() error { return s.GetHealthIntervention().Validate() }},
		validationBranch{present: s.GetHealthHold() != nil, validate: func() error { return s.GetHealthHold().Validate() }},
		validationBranch{present: s.GetHealthEvacuation() != nil, validate: func() error { return s.GetHealthEvacuation().Validate() }},
	)
}

func (f LinkFrame) Validate() error {
	if f.LinkMessageId == "" {
		return errors.New("link frame link message ID is required")
	}
	if f.ReferenceLinkMessageId != nil && *f.ReferenceLinkMessageId == "" {
		return errors.New("link frame reference link message ID must not be empty")
	}
	if !f.DeliveryMethod.IsValid() {
		return errors.New("link frame delivery method is required")
	}
	if !f.Representation.IsValid() {
		return errors.New("link frame representation is required")
	}
	if !f.Representation.IsV1Supported() {
		return errors.New("v1 link frame representation must be BINARY_PROTO")
	}
	return validateExactlyOneBranch(
		"link frame",
		validationBranch{
			present:  f.GetRequestContainer() != nil,
			validate: func() error { return f.GetRequestContainer().Validate() },
		},
		validationBranch{
			present:  f.GetSynchronizedResponse() != nil,
			validate: func() error { return f.GetSynchronizedResponse().Validate() },
		},
		validationBranch{
			present:  f.GetCanonicalRegistry() != nil,
			validate: func() error { return f.GetCanonicalRegistry().Validate() },
		},
	)
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
		{"zap or edi-pi", s.ZapOrEdiPi},
		{"earliest departure date", s.EarliestDepartureDateLocal},
		{"latest departure date", s.LatestDepartureDateLocal},
		{"departure location", s.DepartureLocation},
		{"destination location", s.DestinationLocation},
		{"estimated baggage weight", s.TotalEstimatedBaggageWeightLbs},
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
	if !s.Hmic.IsValid() {
		return errors.New("hmic is required")
	}
	if !s.Handling.IsValid() {
		return errors.New("handling is required")
	}
	required := []struct {
		name  string
		value string
	}{
		{"item by niin", s.ItemByNiin},
		{"item quantity", s.ItemQuantity},
		{"serial number", s.SerialNumber},
		{"gross weight", s.GrossWeightLbs},
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
	if s.GetAttachmentIndicator() != AttachmentIndicatorCodeUnspecified &&
		!s.GetAttachmentIndicator().IsValid() {
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
		{"niin", s.Niin},
		{"number of pieces", s.NumberOfPiecesRequiringSupport},
		{"maintenance support date", s.DateMaintenanceSupportRequiredLocal},
		{"location of equipment", s.LocationOfEquipment},
	}
	return validateRequiredFields(required)
}

func (s EngineerReconAreaReportSegment) Validate() error {
	if s.DateOfEvaluationLocal == "" {
		return errors.New("engineer area date of evaluation is required")
	}
	if s.ElementLeaderZapOrEdiPi == "" {
		return errors.New("engineer area element leader is required")
	}
	if s.AreaLocation == "" {
		return errors.New("engineer area location is required")
	}
	if len(s.ItemsReport) == 0 {
		return errors.New("engineer area items report requires at least one entry")
	}
	for _, item := range s.ItemsReport {
		if item == nil {
			return errors.New("engineer area item report must not be nil")
		}
		if item.ItemLocation == "" {
			return errors.New("engineer area item location is required")
		}
		if item.ItemLabel == "" {
			return errors.New("engineer area item label is required")
		}
	}
	if water := s.GetWaterSource(); water != nil {
		if water.FlowRateVelocity == "" || water.QuantityKiloliters == "" {
			return errors.New("engineer area water source flow and quantity are required")
		}
	}
	if s.GetAttachmentIndicator() != AttachmentIndicatorCodeUnspecified &&
		!s.GetAttachmentIndicator().IsValid() {
		return errors.New("engineer area attachment indicator must be 0 or 1")
	}
	return nil
}

func (s EngineerReconZoneReportSegment) Validate() error {
	if s.DateOfEvaluationLocal == "" {
		return errors.New("engineer zone date of evaluation is required")
	}
	if s.ElementLeaderZapOrEdiPi == "" {
		return errors.New("engineer zone element leader is required")
	}
	if len(s.EnemyReport) == 0 {
		return errors.New("engineer zone enemy report requires at least one entry")
	}
	for _, enemy := range s.EnemyReport {
		if enemy == nil {
			return errors.New("engineer zone enemy report must not be nil")
		}
		if enemy.EnemyLocation == "" {
			return errors.New("engineer zone enemy location is required")
		}
		if enemy.EnemyLabel == "" {
			return errors.New("engineer zone enemy label is required")
		}
	}
	if s.ZoneLocation == "" {
		return errors.New("engineer zone location is required")
	}
	if water := s.GetWaterSource(); water != nil {
		if water.FlowRateVelocity == "" || water.QuantityKiloliters == "" {
			return errors.New("engineer zone water source flow and quantity are required")
		}
	}
	if len(s.ItemsReport) == 0 {
		return errors.New("engineer zone items report requires at least one entry")
	}
	for _, item := range s.ItemsReport {
		if item == nil {
			return errors.New("engineer zone item report must not be nil")
		}
		if item.ItemLocation == "" {
			return errors.New("engineer zone item location is required")
		}
		if item.ItemLabel == "" {
			return errors.New("engineer zone item label is required")
		}
	}
	if s.GetAttachmentIndicator() != AttachmentIndicatorCodeUnspecified &&
		!s.GetAttachmentIndicator().IsValid() {
		return errors.New("engineer zone attachment indicator must be 0 or 1")
	}
	return nil
}

func (s EngineerReconRouteReportSegment) Validate() error {
	if s.DateOfEvaluationLocal == "" {
		return errors.New("engineer route date of evaluation is required")
	}
	if s.ElementLeaderZapOrEdiPi == "" {
		return errors.New("engineer route element leader is required")
	}
	if len(s.EnemyReport) == 0 {
		return errors.New("engineer route enemy report requires at least one entry")
	}
	for _, enemy := range s.EnemyReport {
		if enemy == nil {
			return errors.New("engineer route enemy report must not be nil")
		}
		if enemy.EnemyLocation == "" {
			return errors.New("engineer route enemy location is required")
		}
		if enemy.EnemyLabel == "" {
			return errors.New("engineer route enemy label is required")
		}
	}
	if len(s.RouteLocations) == 0 {
		return errors.New("engineer route locations require at least one waypoint")
	}
	for _, waypoint := range s.RouteLocations {
		if waypoint == "" {
			return errors.New("engineer route waypoint must not be empty")
		}
	}
	if water := s.GetWaterSource(); water != nil {
		if water.FlowRateVelocity == "" || water.QuantityKiloliters == "" {
			return errors.New("engineer route water source flow and quantity are required")
		}
	}
	if len(s.ItemsReport) == 0 {
		return errors.New("engineer route items report requires at least one entry")
	}
	for _, item := range s.ItemsReport {
		if item == nil {
			return errors.New("engineer route item report must not be nil")
		}
		if item.ItemLocation == "" {
			return errors.New("engineer route item location is required")
		}
		if item.ItemLabel == "" {
			return errors.New("engineer route item label is required")
		}
	}
	if s.GetAttachmentIndicator() != AttachmentIndicatorCodeUnspecified &&
		!s.GetAttachmentIndicator().IsValid() {
		return errors.New("engineer route attachment indicator must be 0 or 1")
	}
	return nil
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
	if s.GetAttachmentIndicator() != AttachmentIndicatorCodeUnspecified &&
		!s.GetAttachmentIndicator().IsValid() {
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
	if s.GetAttachmentIndicator() != AttachmentIndicatorCodeUnspecified &&
		!s.GetAttachmentIndicator().IsValid() {
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
	if s.GetAttachmentIndicator() != AttachmentIndicatorCodeUnspecified &&
		!s.GetAttachmentIndicator().IsValid() {
		return errors.New(
			"obstacle attachment indicator must be 0 or 1",
		)
	}
	return nil
}

func (s ExplosiveOrdnanceDisposalSegment) Validate() error {
	if s.DateOfUxoDiscovery == "" {
		return errors.New("date of uxo discovery is required")
	}
	if s.RequestedDateOfEodAction == "" {
		return errors.New("requested date of eod action is required")
	}
	if s.LocationOfUxo == "" {
		return errors.New("location of uxo is required")
	}
	if !s.TypeOfCbrnAgent.IsValid() {
		return errors.New("type of cbrn agent is required")
	}
	if s.GetPhysicalPropertyOfCbrnAgent() != CBRNPhysicalPropertyCodeUnspecified &&
		!s.GetPhysicalPropertyOfCbrnAgent().IsValid() {
		return errors.New("physical property of cbrn agent is invalid")
	}
	if s.GetContaminationValueOfCbrnAgent() != CBRNContaminationValueCodeUnspecified &&
		!s.GetContaminationValueOfCbrnAgent().IsValid() {
		return errors.New("contamination value of cbrn agent is invalid")
	}
	if s.MunitionColor == "" {
		return errors.New("munition color is required")
	}
	if s.MunitionMarkings == "" {
		return errors.New("munition markings are required")
	}
	if !s.MunitionPurpose.IsValid() {
		return errors.New("munition purpose is required")
	}
	if !s.MunitionType.IsValid() {
		return errors.New("munition type is required")
	}
	if s.GetAttachmentIndicator() != AttachmentIndicatorCodeUnspecified &&
		!s.GetAttachmentIndicator().IsValid() {
		return errors.New("eod attachment indicator must be 0 or 1")
	}
	return nil
}

func (s GeneralEngineeringBulkLiquidSupportSegment) Validate() error {
	if s.DateOfEvaluationLocal == "" {
		return errors.New("bulk liquid date of evaluation is required")
	}
	if s.LocationOfBulkLiquid == "" {
		return errors.New("bulk liquid location is required")
	}
	if !s.Estimate.IsValid() {
		return errors.New("bulk liquid estimate is required")
	}
	if s.GetAttachmentIndicator() != AttachmentIndicatorCodeUnspecified &&
		!s.GetAttachmentIndicator().IsValid() {
		return errors.New("bulk liquid attachment indicator must be 0 or 1")
	}
	return nil
}

func (s GeneralEngineeringDemolitionSegment) Validate() error {
	if s.DateOfEvaluationLocal == "" {
		return errors.New("demolition date of evaluation is required")
	}
	if s.Location == "" {
		return errors.New("demolition location is required")
	}
	if s.TypeOfDemolition == "" {
		return errors.New("demolition type is required")
	}
	if s.RouteNumber == "" {
		return errors.New("demolition route number is required")
	}
	if !s.DeterminationOfAction.IsValid() {
		return errors.New("demolition determination of action is required")
	}
	if !s.Bypass.IsValid() {
		return errors.New("demolition bypass is required")
	}
	if s.BypassGrid == "" {
		return errors.New("demolition bypass grid is required")
	}
	if s.GetAttachmentIndicator() != AttachmentIndicatorCodeUnspecified &&
		!s.GetAttachmentIndicator().IsValid() {
		return errors.New("demolition attachment indicator must be 0 or 1")
	}
	return nil
}

func (s HealthCollectionSegment) Validate() error {
	if s.SegmentNumber == 0 {
		return errors.New("health collection segment number is required")
	}
	if !s.RequestTypeCode.IsValid() {
		return errors.New("health collection request type code must be CR")
	}
	if !s.RequestPriority.IsValid() {
		return errors.New("health collection request priority is required")
	}
	if s.ZapOrEdiPi == "" {
		return errors.New("health collection zap or edi-pi is required")
	}
	if s.LastName == "" {
		return errors.New("health collection last name is required")
	}
	if s.FirstName == "" {
		return errors.New("health collection first name is required")
	}
	if s.GetService() != ServiceCodeUnspecified && !s.GetService().IsValid() {
		return errors.New("health collection service is invalid")
	}
	if s.Allergies == "" {
		return errors.New("health collection allergies are required")
	}
	if s.DateOfInjuryLocal == "" {
		return errors.New("health collection date of injury is required")
	}
	if s.TimeOfInjuryLocal == "" {
		return errors.New("health collection time of injury is required")
	}
	if s.LocationInjuryOccurred == "" {
		return errors.New("health collection injury location is required")
	}
	return nil
}

func (s HealthTriageSegment) Validate() error {
	if !s.PrimaryMechanismOfInjury.IsValid() {
		return errors.New("health triage primary mechanism is required")
	}
	if !s.CbrnRelatedExposure.IsValid() {
		return errors.New("health triage cbrn exposure is required")
	}
	if !s.MajorSignsSymptoms.IsValid() {
		return errors.New("health triage major signs symptoms is required")
	}
	if len(s.InjuryLocations) == 0 || len(s.InjuryLocations) > 10 {
		return errors.New("health triage injury locations require 1 to 10 entries")
	}
	for _, location := range s.InjuryLocations {
		if location == "" {
			return errors.New("health triage injury locations must not be empty")
		}
	}
	if s.GetPulseLocation() != HealthPulseLocationCodeUnspecified &&
		!s.GetPulseLocation().IsValid() {
		return errors.New("health triage pulse location must be W or N")
	}
	if s.GetResponsiveness() != HealthResponsivenessCodeUnspecified &&
		!s.GetResponsiveness().IsValid() {
		return errors.New("health triage responsiveness is invalid")
	}
	if s.PainScale != nil && *s.PainScale > 10 {
		return errors.New("health triage pain scale must be 0 to 10")
	}
	if !s.TriagePrecedence.IsValid() {
		return errors.New("health triage precedence is required")
	}
	return nil
}

func (s HealthInterventionSegment) Validate() error {
	if len(s.Tourniquets) > 4 {
		return errors.New("health intervention allows at most 4 tourniquets")
	}
	for _, treatment := range s.Tourniquets {
		if treatment == nil {
			return errors.New("tourniquet treatment must not be nil")
		}
		if !treatment.Placement.IsValid() {
			return errors.New("tourniquet placement is invalid")
		}
		if !treatment.Type.IsValid() {
			return errors.New("tourniquet type is invalid")
		}
		if treatment.DateLocal == "" || treatment.TimeLocal == "" {
			return errors.New("tourniquet date and time are required")
		}
	}
	if len(s.WoundTreatments) == 0 || len(s.WoundTreatments) > 7 {
		return errors.New("health intervention wound treatments require 1 to 7 entries")
	}
	for _, treatment := range s.WoundTreatments {
		if !treatment.IsValid() {
			return errors.New("wound treatment is invalid")
		}
	}
	if !s.AirwayTreatment.IsValid() {
		return errors.New("airway treatment is required")
	}
	if !s.BreathingTreatment.IsValid() {
		return errors.New("breathing treatment is required")
	}
	for _, treatment := range s.FluidCirculationTreatments {
		if treatment == nil {
			return errors.New("fluid treatment must not be nil")
		}
		if err := validateFluidCirculationTreatment(treatment); err != nil {
			return err
		}
	}
	for _, treatment := range s.BloodCirculationTreatments {
		if treatment == nil {
			return errors.New("blood treatment must not be nil")
		}
		if err := validateBloodCirculationTreatment(treatment); err != nil {
			return err
		}
	}
	for _, treatment := range s.AnalgesicMedicationTreatments {
		if treatment == nil {
			return errors.New("analgesic treatment must not be nil")
		}
		if err := validateAnalgesicMedicationTreatment(treatment); err != nil {
			return err
		}
	}
	for _, treatment := range s.AntibioticMedicationTreatments {
		if treatment == nil {
			return errors.New("antibiotic treatment must not be nil")
		}
		if err := validateAntibioticMedicationTreatment(treatment); err != nil {
			return err
		}
	}
	for _, treatment := range s.OtherMedicationTreatments {
		if treatment == nil {
			return errors.New("other medication treatment must not be nil")
		}
		if err := validateOtherMedicationTreatment(treatment); err != nil {
			return err
		}
	}
	if !s.CasualtyType.IsValid() {
		return errors.New("casualty type is required")
	}
	if s.FirstResponderZapOrEdiPi == "" {
		return errors.New("first responder zap or edi-pi is required")
	}
	return nil
}

func (s HealthHoldSegment) Validate() error {
	if len(s.TriageEntries) == 0 && len(s.InterventionEntries) == 0 {
		return errors.New("health hold requires triage or intervention entries")
	}
	for _, entry := range s.TriageEntries {
		if entry == nil {
			return errors.New("health hold triage entry must not be nil")
		}
		if err := entry.Validate(); err != nil {
			return fmt.Errorf("health hold triage entry: %w", err)
		}
	}
	for _, entry := range s.InterventionEntries {
		if entry == nil {
			return errors.New("health hold intervention entry must not be nil")
		}
		if err := entry.Validate(); err != nil {
			return fmt.Errorf("health hold intervention entry: %w", err)
		}
	}
	return nil
}

func (s HealthEvacuationSegment) Validate() error {
	if !s.RequestPriority.IsValid() {
		return errors.New("health evacuation request priority is required")
	}
	if s.LocationOfPickup == "" {
		return errors.New("health evacuation location of pickup is required")
	}
	if s.GetLocationMarking() !=
		HealthEvacuationLocationMarkingCodeUnspecified &&
		!s.GetLocationMarking().IsValid() {
		return errors.New("health evacuation location marking is invalid")
	}
	if s.GetLocationContamination() !=
		HealthEvacuationContaminationCodeUnspecified &&
		!s.GetLocationContamination().IsValid() {
		return errors.New("health evacuation location contamination is invalid")
	}
	if s.ContactSettings == "" {
		return errors.New("health evacuation contact settings are required")
	}
	if len(s.CountOfCasualtiesPrecedence) == 0 ||
		len(s.CountOfCasualtiesPrecedence) > 5 {
		return errors.New("health evacuation precedence counts require 1 to 5 entries")
	}
	for _, entry := range s.CountOfCasualtiesPrecedence {
		if entry == nil {
			return errors.New("health evacuation precedence count must not be nil")
		}
		if !entry.Precedence.IsValid() {
			return errors.New("health evacuation precedence code is invalid")
		}
		if entry.CasualtyCount == "" {
			return errors.New("health evacuation precedence casualty count is required")
		}
	}
	if len(s.CountOfCasualtyTypes) == 0 || len(s.CountOfCasualtyTypes) > 2 {
		return errors.New("health evacuation casualty type counts require 1 to 2 entries")
	}
	for _, entry := range s.CountOfCasualtyTypes {
		if entry == nil {
			return errors.New("health evacuation casualty type count must not be nil")
		}
		if !entry.CasualtyType.IsValid() {
			return errors.New("health evacuation casualty type is invalid")
		}
		if entry.CasualtyCount == "" {
			return errors.New("health evacuation casualty type count is required")
		}
	}
	if s.GetRequestedEquipment() !=
		HealthEvacuationRequestedEquipmentCodeUnspecified &&
		!s.GetRequestedEquipment().IsValid() {
		return errors.New("health evacuation requested equipment is invalid")
	}
	if s.GetSecurity() != HealthEvacuationSecurityCodeUnspecified &&
		!s.GetSecurity().IsValid() {
		return errors.New("health evacuation security is invalid")
	}
	if len(s.Casualties) == 0 {
		return errors.New("health evacuation requires at least one casualty record")
	}
	for _, casualty := range s.Casualties {
		if casualty == nil {
			return errors.New("health evacuation casualty must not be nil")
		}
		if err := validateHealthEvacuationCasualtyRecord(casualty); err != nil {
			return err
		}
	}
	return nil
}

func validateFluidCirculationTreatment(treatment *FluidCirculationTreatment) error {
	if treatment.GetFluidNameCode() == FluidNameCodeUnspecified &&
		treatment.GetOtherFluidName() == "" {
		return errors.New("fluid treatment requires a fluid name")
	}
	if treatment.GetFluidNameCode() != FluidNameCodeUnspecified &&
		!treatment.GetFluidNameCode().IsValid() {
		return errors.New("fluid treatment name code is invalid")
	}
	if treatment.GetVolumeDose() == "" || treatment.GetDateLocal() == "" ||
		treatment.GetTimeLocal() == "" {
		return errors.New("fluid treatment volume, date, and time are required")
	}
	if !treatment.GetRoute().IsValid() {
		return errors.New("fluid treatment route is required")
	}
	return nil
}

func validateBloodCirculationTreatment(treatment *BloodCirculationTreatment) error {
	if !treatment.GetBloodProductName().IsValid() {
		return errors.New("blood treatment product is required")
	}
	if treatment.GetVolumeDose() == "" || treatment.GetDateLocal() == "" ||
		treatment.GetTimeLocal() == "" {
		return errors.New("blood treatment volume, date, and time are required")
	}
	if !treatment.GetRoute().IsValid() {
		return errors.New("blood treatment route is required")
	}
	return nil
}

func validateAnalgesicMedicationTreatment(treatment *AnalgesicMedicationTreatment) error {
	if treatment.GetMedicationCode() == AnalgesicMedicationCodeUnspecified &&
		treatment.GetOtherMedicationName() == "" {
		return errors.New("analgesic treatment requires a medication")
	}
	if treatment.GetMedicationCode() != AnalgesicMedicationCodeUnspecified &&
		!treatment.GetMedicationCode().IsValid() {
		return errors.New("analgesic medication code is invalid")
	}
	if treatment.GetVolumeDose() == "" || treatment.GetDateLocal() == "" ||
		treatment.GetTimeLocal() == "" {
		return errors.New("analgesic treatment volume, date, and time are required")
	}
	if !treatment.GetRoute().IsValid() {
		return errors.New("analgesic treatment route is required")
	}
	return nil
}

func validateAntibioticMedicationTreatment(treatment *AntibioticMedicationTreatment) error {
	if treatment.GetMedicationCode() == AntibioticMedicationCodeUnspecified &&
		treatment.GetOtherMedicationName() == "" {
		return errors.New("antibiotic treatment requires a medication")
	}
	if treatment.GetMedicationCode() != AntibioticMedicationCodeUnspecified &&
		!treatment.GetMedicationCode().IsValid() {
		return errors.New("antibiotic medication code is invalid")
	}
	if treatment.GetVolumeDose() == "" || treatment.GetDateLocal() == "" ||
		treatment.GetTimeLocal() == "" {
		return errors.New("antibiotic treatment volume, date, and time are required")
	}
	if !treatment.GetRoute().IsValid() {
		return errors.New("antibiotic treatment route is required")
	}
	return nil
}

func validateOtherMedicationTreatment(treatment *OtherMedicationTreatment) error {
	if treatment.GetMedicationCode() == OtherMedicationCodeUnspecified &&
		treatment.GetOtherMedicationName() == "" {
		return errors.New("other medication treatment requires a medication")
	}
	if treatment.GetMedicationCode() != OtherMedicationCodeUnspecified &&
		!treatment.GetMedicationCode().IsValid() {
		return errors.New("other medication code is invalid")
	}
	if treatment.GetVolumeDose() == "" || treatment.GetDateLocal() == "" ||
		treatment.GetTimeLocal() == "" {
		return errors.New("other medication treatment volume, date, and time are required")
	}
	if !treatment.GetRoute().IsValid() {
		return errors.New("other medication treatment route is required")
	}
	return nil
}

func validateHealthEvacuationCasualtyRecord(
	record *HealthEvacuationCasualtyRecord,
) error {
	if record.GetZapOrEdiPi() == "" {
		return errors.New("health evacuation casualty zap or edi-pi is required")
	}
	if record.GetPrimaryMechanismOfInjury() !=
		HealthPrimaryMechanismCodeUnspecified &&
		!record.GetPrimaryMechanismOfInjury().IsValid() {
		return errors.New("health evacuation casualty primary mechanism is invalid")
	}
	for _, treatment := range record.GetTourniquets() {
		if treatment == nil {
			return errors.New("health evacuation casualty tourniquet must not be nil")
		}
		if !treatment.Placement.IsValid() {
			return errors.New("health evacuation casualty tourniquet placement is invalid")
		}
		if !treatment.Type.IsValid() {
			return errors.New("health evacuation casualty tourniquet type is invalid")
		}
		if treatment.DateLocal == "" || treatment.TimeLocal == "" {
			return errors.New("health evacuation casualty tourniquet date and time are required")
		}
	}
	for _, treatment := range record.GetWoundTreatments() {
		if !treatment.IsValid() {
			return errors.New("health evacuation casualty wound treatment is invalid")
		}
	}
	if record.GetAirwayTreatment() != AirwayTreatmentCodeUnspecified &&
		!record.GetAirwayTreatment().IsValid() {
		return errors.New("health evacuation casualty airway treatment is invalid")
	}
	if record.GetBreathingTreatment() != BreathingTreatmentCodeUnspecified &&
		!record.GetBreathingTreatment().IsValid() {
		return errors.New("health evacuation casualty breathing treatment is invalid")
	}
	if treatment := record.GetFluidCirculationTreatment(); treatment != nil {
		if err := validateFluidCirculationTreatment(treatment); err != nil {
			return fmt.Errorf("health evacuation casualty fluid treatment: %w", err)
		}
	}
	if treatment := record.GetBloodCirculationTreatment(); treatment != nil {
		if err := validateBloodCirculationTreatment(treatment); err != nil {
			return fmt.Errorf("health evacuation casualty blood treatment: %w", err)
		}
	}
	if treatment := record.GetAnalgesicMedicationTreatment(); treatment != nil {
		if err := validateAnalgesicMedicationTreatment(treatment); err != nil {
			return fmt.Errorf("health evacuation casualty analgesic treatment: %w", err)
		}
	}
	if treatment := record.GetAntibioticMedicationTreatment(); treatment != nil {
		if err := validateAntibioticMedicationTreatment(treatment); err != nil {
			return fmt.Errorf("health evacuation casualty antibiotic treatment: %w", err)
		}
	}
	if treatment := record.GetOtherMedicationTreatment(); treatment != nil {
		if err := validateOtherMedicationTreatment(treatment); err != nil {
			return fmt.Errorf("health evacuation casualty other medication treatment: %w", err)
		}
	}
	if record.GetCbrnRelatedExposure() != HealthCBRNExposureCodeUnspecified &&
		!record.GetCbrnRelatedExposure().IsValid() {
		return errors.New("health evacuation casualty cbrn exposure is invalid")
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

func (e CanonicalRegistryFieldEntry) Validate() error {
	required := []struct {
		name  string
		value string
	}{
		{"activity", e.Activity},
		{"data element", e.DataElement},
		{"data field", e.DataField},
		{"canonical block", e.CanonicalBlock},
		{"canonical field", e.CanonicalField},
		{"source reference", e.SourceReference},
	}
	if err := validateRequiredFields(required); err != nil {
		return err
	}
	if !e.CanonicalFile.IsValid() {
		return errors.New("canonical registry file type is required")
	}
	if !e.ExchangeRole.IsValid() {
		return errors.New("canonical registry exchange role is required")
	}
	return nil
}

func (r CanonicalRegistry) Validate() error {
	if len(r.Entries) == 0 {
		return errors.New("canonical registry requires at least one entry")
	}
	for i, entry := range r.Entries {
		if entry == nil {
			return fmt.Errorf("canonical registry entry %d must not be nil", i)
		}
		if err := entry.Validate(); err != nil {
			return fmt.Errorf("canonical registry entry %d: %w", i, err)
		}
	}
	return nil
}
