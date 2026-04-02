package lxdr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func stringPtr(value string) *string {
	return &value
}

func ptr[T any](value T) *T {
	return &value
}

func uint32Ptr(value uint32) *uint32 {
	return &value
}

func attachmentPtr(value AttachmentIndicatorCode) *AttachmentIndicatorCode {
	return &value
}

func servicePtr(value ServiceCode) *ServiceCode {
	return &value
}

func cbrnPhysicalPropertyPtr(
	value CBRNPhysicalPropertyCode,
) *CBRNPhysicalPropertyCode {
	return &value
}

func cbrnContaminationPtr(
	value CBRNContaminationValueCode,
) *CBRNContaminationValueCode {
	return &value
}

func testHeader() *RequestHeader {
	return &RequestHeader{
		LocalSystemDate:                 "2027OCT13",
		LocalSystemTime:                 "15470352",
		SynchronizedGeospatialReference: "4QFJ123456",
		LocalRequestId:                  "3838JBNM5",
		RequestPriority:                 RequestPriorityCode02,
		ElementUnitIdOrCallsign:         "KL9K",
		RequestSegmentCount:             1,
	}
}

func wrapMobilityPax(seg *MobilityPaxRequestSegment) *RequestSegment {
	return &RequestSegment{
		FunctionFamily: FunctionFamilyMobility,
		Segment: &RequestSegment_MobilityPax{
			MobilityPax: seg,
		},
	}
}

func wrapEngineerRoad(seg *EngineerReconRoadReportSegment) *RequestSegment {
	return &RequestSegment{
		FunctionFamily: FunctionFamilyGeneralEngineering,
		Segment: &RequestSegment_EngineerReconRoad{
			EngineerReconRoad: seg,
		},
	}
}

func wrapEngineerLZ(seg *EngineerReconLandingZoneReportSegment) *RequestSegment {
	return &RequestSegment{
		FunctionFamily: FunctionFamilyGeneralEngineering,
		Segment: &RequestSegment_EngineerReconLandingZone{
			EngineerReconLandingZone: seg,
		},
	}
}

func wrapObstacleRemoval(
	seg *GeneralEngineeringObstacleRemovalSegment,
) *RequestSegment {
	return &RequestSegment{
		FunctionFamily: FunctionFamilyGeneralEngineering,
		Segment: &RequestSegment_ObstacleRemoval{
			ObstacleRemoval: seg,
		},
	}
}

func wrapEOD(seg *ExplosiveOrdnanceDisposalSegment) *RequestSegment {
	return &RequestSegment{
		FunctionFamily: FunctionFamilyGeneralEngineering,
		Segment: &RequestSegment_Eod{
			Eod: seg,
		},
	}
}

func wrapHealthCollection(seg *HealthCollectionSegment) *RequestSegment {
	return &RequestSegment{
		FunctionFamily: FunctionFamilyHealthServices,
		Segment: &RequestSegment_HealthCollection{
			HealthCollection: seg,
		},
	}
}

func wrapHealthHold(seg *HealthHoldSegment) *RequestSegment {
	return &RequestSegment{
		FunctionFamily: FunctionFamilyHealthServices,
		Segment: &RequestSegment_HealthHold{
			HealthHold: seg,
		},
	}
}

func readFixture(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %q: %v", name, err)
	}

	return strings.TrimSpace(string(data))
}

func TestParseHeaderExample(t *testing.T) {
	header, err := ParseCanonicalHeader(
		readFixture(t, "header_example.txt"),
	)
	if err != nil {
		t.Fatalf("parse header: %v", err)
	}

	if header.LocalSystemDate != "2027OCT13" {
		t.Fatalf("unexpected date: %q", header.LocalSystemDate)
	}
	if header.LocalRequestId != "3838JBNM5" {
		t.Fatalf("unexpected local request id: %q", header.LocalRequestId)
	}
	if header.RequestSegmentCount != 1 {
		t.Fatalf("unexpected segment count: %d", header.RequestSegmentCount)
	}
}

func TestRenderHeaderExample(t *testing.T) {
	header := RequestHeader{
		LocalSystemDate:                 "2027OCT13",
		LocalSystemTime:                 "15470352",
		SynchronizedGeospatialReference: "4QFJ123456",
		LocalRequestId:                  "3838JBNM5",
		RequestPriority:                 RequestPriorityCode02,
		ElementUnitIdOrCallsign:         "KL9K",
		RequestSegmentCount:             1,
	}

	got, err := header.RenderCanonical()
	if err != nil {
		t.Fatalf("render header: %v", err)
	}

	want := readFixture(t, "header_example.txt")
	if got != want {
		t.Fatalf("rendered header = %q, want %q", got, want)
	}
}

func TestParseSynchronizedResponseExample(t *testing.T) {
	resp, err := ParseSynchronizedResponse(
		readFixture(t, "sync_response_example.txt"),
	)
	if err != nil {
		t.Fatalf("parse sync response: %v", err)
	}

	if resp.LocalRequestId != "3838JBNM5" {
		t.Fatalf("unexpected local request id: %q", resp.LocalRequestId)
	}
	if resp.SynchronizedRequestId != "KL9K15474QFJ" {
		t.Fatalf(
			"unexpected synchronized request id: %q",
			resp.SynchronizedRequestId,
		)
	}
}

func TestRenderSynchronizedResponseExample(t *testing.T) {
	resp := SynchronizedResponse{
		LocalRequestId:        "3838JBNM5",
		SynchronizedRequestId: "KL9K15474QFJ",
	}

	got, err := resp.RenderCanonical()
	if err != nil {
		t.Fatalf("render sync response: %v", err)
	}

	want := readFixture(t, "sync_response_example.txt")
	if got != want {
		t.Fatalf("rendered sync response = %q, want %q", got, want)
	}
}

func TestParsePAXExample(t *testing.T) {
	seg, err := ParseCanonicalPAXSegment(
		readFixture(t, "pax_example.txt"),
	)
	if err != nil {
		t.Fatalf("parse pax: %v", err)
	}

	if seg.RequestTypeCode != MobilityPaxRequestTypeCodePM {
		t.Fatalf("unexpected request type: %q", seg.RequestTypeCode)
	}
	if seg.ZapOrEdiPi != "1010919789" {
		t.Fatalf("unexpected personnel identifier: %q", seg.ZapOrEdiPi)
	}
	if seg.DestinationLocation != "4QFJ456789" {
		t.Fatalf(
			"unexpected destination location: %q",
			seg.DestinationLocation,
		)
	}
}

func TestRenderPAXExample(t *testing.T) {
	seg := MobilityPaxRequestSegment{
		SegmentNumber:                  1,
		RequestTypeCode:                MobilityPaxRequestTypeCodePM,
		RequestPriority:                RequestPriorityCode02,
		ZapOrEdiPi:                     "1010919789",
		EarliestDepartureDateLocal:     "2027OCT15",
		LatestDepartureDateLocal:       "2027OCT20",
		DepartureLocation:              "4QFJ123456",
		DestinationLocation:            "4QFJ456789",
		TotalEstimatedBaggageWeightLbs: "075",
		HazardousMaterialType:          "X",
	}

	got, err := seg.RenderCanonical()
	if err != nil {
		t.Fatalf("render pax: %v", err)
	}

	want := readFixture(t, "pax_example.txt")
	if got != want {
		t.Fatalf("rendered pax = %q, want %q", got, want)
	}
}

func TestParseCargoExample(t *testing.T) {
	seg, err := ParseCanonicalCargoSegment(
		readFixture(t, "cargo_example.txt"),
	)
	if err != nil {
		t.Fatalf("parse cargo: %v", err)
	}

	if seg.RequestTypeCode != MobilityCargoRequestTypeCodeCM {
		t.Fatalf("unexpected request type: %q", seg.RequestTypeCode)
	}
	if seg.ItemByNiin != "015519434" {
		t.Fatalf("unexpected niin: %q", seg.ItemByNiin)
	}
	if seg.SerialNumber != "598742" {
		t.Fatalf("unexpected serial number: %q", seg.SerialNumber)
	}
}

func TestRenderCargoExample(t *testing.T) {
	seg := MobilityCargoRequestSegment{
		SegmentNumber:              1,
		RequestTypeCode:            MobilityCargoRequestTypeCodeCM,
		RequestPriority:            RequestPriorityCode02,
		ItemByNiin:                 "015519434",
		ItemQuantity:               "1",
		SerialNumber:               "598742",
		GrossWeightLbs:             "28000",
		ActualHeightInches:         "126",
		ActualWidthInches:          "100",
		ActualLengthInches:         "315",
		Hmic:                       CargoHMICCodeD,
		Handling:                   CargoHandlingCodeR,
		EarliestDepartureDateLocal: "2027OCT15",
		LatestDepartureDateLocal:   "2027OCT20",
		DepartureLocation:          "4QFJ123456",
		DestinationLocation:        "4QFJ456789",
	}

	got, err := seg.RenderCanonical()
	if err != nil {
		t.Fatalf("render cargo: %v", err)
	}

	want := readFixture(t, "cargo_example.txt")
	if got != want {
		t.Fatalf("rendered cargo = %q, want %q", got, want)
	}
}

func TestRequestContainerSegmentCountMatchesHeader(t *testing.T) {
	header := testHeader()
	pax := &MobilityPaxRequestSegment{
		SegmentNumber:                  1,
		RequestTypeCode:                MobilityPaxRequestTypeCodePM,
		RequestPriority:                RequestPriorityCode02,
		ZapOrEdiPi:                     "1010919789",
		EarliestDepartureDateLocal:     "2027OCT15",
		LatestDepartureDateLocal:       "2027OCT20",
		DepartureLocation:              "4QFJ123456",
		DestinationLocation:            "4QFJ456789",
		TotalEstimatedBaggageWeightLbs: "075",
		HazardousMaterialType:          "X",
	}
	container := RequestContainer{
		Header:   header,
		Segments: []*RequestSegment{wrapMobilityPax(pax)},
	}

	if err := container.Validate(); err != nil {
		t.Fatalf("validate container: %v", err)
	}
}

func TestCargoHMICRejectsUnknownValue(t *testing.T) {
	seg := MobilityCargoRequestSegment{
		SegmentNumber:              1,
		RequestTypeCode:            MobilityCargoRequestTypeCodeCM,
		RequestPriority:            RequestPriorityCode02,
		ItemByNiin:                 "015519434",
		ItemQuantity:               "1",
		SerialNumber:               "598742",
		GrossWeightLbs:             "28000",
		ActualHeightInches:         "126",
		ActualWidthInches:          "100",
		ActualLengthInches:         "315",
		Hmic:                       CargoHMICCode(99),
		Handling:                   CargoHandlingCodeR,
		EarliestDepartureDateLocal: "2027OCT15",
		LatestDepartureDateLocal:   "2027OCT20",
		DepartureLocation:          "4QFJ123456",
		DestinationLocation:        "4QFJ456789",
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid hmic error")
	}
}

func TestCargoHandlingRejectsUnknownValue(t *testing.T) {
	seg := MobilityCargoRequestSegment{
		SegmentNumber:              1,
		RequestTypeCode:            MobilityCargoRequestTypeCodeCM,
		RequestPriority:            RequestPriorityCode02,
		ItemByNiin:                 "015519434",
		ItemQuantity:               "1",
		SerialNumber:               "598742",
		GrossWeightLbs:             "28000",
		ActualHeightInches:         "126",
		ActualWidthInches:          "100",
		ActualLengthInches:         "315",
		Hmic:                       CargoHMICCodeD,
		Handling:                   CargoHandlingCode(99),
		EarliestDepartureDateLocal: "2027OCT15",
		LatestDepartureDateLocal:   "2027OCT20",
		DepartureLocation:          "4QFJ123456",
		DestinationLocation:        "4QFJ456789",
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid handling error")
	}
}

func TestSupplySegmentValidateAllowsUnknownNIINWithNarrative(t *testing.T) {
	seg := SupplyRequestSegment{
		SegmentNumber:             1,
		RequestTypeCode:           SupplyRequestTypeCodeSR,
		RequestPriority:           RequestPriorityCode02,
		ItemByNiin:                nil,
		ItemQuantity:              "25",
		RequiredDeliveryDateLocal: "2027OCT21",
		DeliveryLocation:          "4QFJ123456",
		Narrative:                 stringPtr("BULK WATER"),
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate supply segment: %v", err)
	}
}

func TestSupplySegmentValidateRejectsMissingRequiredFields(t *testing.T) {
	seg := SupplyRequestSegment{
		SegmentNumber:             1,
		RequestTypeCode:           SupplyRequestTypeCodeSR,
		RequestPriority:           RequestPriorityCodeUnspecified,
		ItemQuantity:              "25",
		RequiredDeliveryDateLocal: "2027OCT21",
		DeliveryLocation:          "4QFJ123456",
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected supply validation error")
	}
}

func TestMaintenanceSegmentValidateAllowsOptionalIdentityFields(t *testing.T) {
	seg := MaintenanceRequestSegment{
		SegmentNumber:                       1,
		RequestTypeCode:                     MaintenanceRequestTypeCodeCM,
		RequestPriority:                     RequestPriorityCode02,
		SerialNumber:                        nil,
		Niin:                                "015519434",
		ModelOfEquipment:                    nil,
		ItemNomenclature:                    nil,
		NumberOfPiecesRequiringSupport:      "1",
		EquipmentOperationalCondition:       MaintenanceOperationalConditionCodeC,
		DateMaintenanceSupportRequiredLocal: "2027OCT21",
		LocationOfEquipment:                 "4QFJ123456",
		TypeOfMaintenanceSupportRequired:    MaintenanceSupportTypeCodeR1,
		TypeOfRepair:                        MaintenanceRepairTypeCodeD1,
		RepairMajorDefect:                   MaintenanceMajorDefectCodeMD07,
		AttachmentIndicator:                 AttachmentIndicatorCodeNo,
		Narrative:                           nil,
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate maintenance segment: %v", err)
	}
}

func TestMaintenanceSegmentValidateRejectsMissingRequiredFields(t *testing.T) {
	seg := MaintenanceRequestSegment{
		SegmentNumber:                       1,
		RequestTypeCode:                     MaintenanceRequestTypeCodeCM,
		RequestPriority:                     RequestPriorityCode02,
		Niin:                                "015519434",
		NumberOfPiecesRequiringSupport:      "1",
		EquipmentOperationalCondition:       MaintenanceOperationalConditionCodeUnspecified,
		DateMaintenanceSupportRequiredLocal: "2027OCT21",
		LocationOfEquipment:                 "4QFJ123456",
		TypeOfMaintenanceSupportRequired:    MaintenanceSupportTypeCodeR1,
		TypeOfRepair:                        MaintenanceRepairTypeCodeD1,
		RepairMajorDefect:                   MaintenanceMajorDefectCodeMD07,
		AttachmentIndicator:                 AttachmentIndicatorCodeNo,
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected maintenance validation error")
	}
}

func TestHealthTriageSegmentValidateAllowsOptionalVitals(t *testing.T) {
	painScale := uint32(7)
	seg := HealthTriageSegment{
		PrimaryMechanismOfInjury:  HealthPrimaryMechanismCodeE2,
		CbrnRelatedExposure:       HealthCBRNExposureCodeX,
		MajorSignsSymptoms:        HealthMajorSignsSymptomsCodeB,
		InjuryLocations:           []string{"FH0", "C00"},
		VitalsCheckDateLocal:      stringPtr("2027OCT13"),
		VitalsCheckTimeLocal:      stringPtr("1547"),
		BodyTemperatureFahrenheit: stringPtr("098"),
		PulseRate:                 stringPtr("084"),
		PulseLocation:             HealthPulseLocationCodeW.Enum(),
		BloodPressure:             stringPtr("120|080"),
		RespiratoryRate:           stringPtr("18"),
		PulseOximetryPercent:      stringPtr("97"),
		Responsiveness:            HealthResponsivenessCodeA.Enum(),
		PainScale:                 &painScale,
		TriagePrecedence:          HealthTriagePrecedenceCodeB,
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate health triage segment: %v", err)
	}
}

func TestHealthTriageSegmentValidateRejectsInvalidValues(t *testing.T) {
	painScale := uint32(11)
	seg := HealthTriageSegment{
		PrimaryMechanismOfInjury: HealthPrimaryMechanismCode(99),
		CbrnRelatedExposure:      HealthCBRNExposureCodeX,
		MajorSignsSymptoms:       HealthMajorSignsSymptomsCodeB,
		InjuryLocations:          []string{"FH0"},
		PainScale:                &painScale,
		TriagePrecedence:         HealthTriagePrecedenceCodeB,
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected health triage validation error")
	}
}

func TestHealthInterventionSegmentValidateAllowsStructuredTreatments(t *testing.T) {
	seg := HealthInterventionSegment{
		Tourniquets: []*TourniquetTreatment{
			{
				Placement: TourniquetPlacementCodeTQRA,
				Type:      TourniquetTypeCodeE,
				DateLocal: "2027OCT13",
				TimeLocal: "1547",
			},
		},
		WoundTreatments:    []WoundTreatmentCode{WoundTreatmentCodeT1},
		AirwayTreatment:    AirwayTreatmentCodeA1,
		BreathingTreatment: BreathingTreatmentCodeB1,
		FluidCirculationTreatments: []*FluidCirculationTreatment{
			{
				FluidNameCode: ptr(FluidNameCodeS),
				VolumeDose:    "0500",
				Route:         FluidRouteCodeIV,
				DateLocal:     "2027OCT13",
				TimeLocal:     "1548",
			},
		},
		BloodCirculationTreatments: []*BloodCirculationTreatment{
			{
				BloodProductName: BloodProductCodeWBD,
				VolumeDose:       "0500",
				Route:            FluidRouteCodeIO,
				DateLocal:        "2027OCT13",
				TimeLocal:        "1550",
			},
		},
		AnalgesicMedicationTreatments: []*AnalgesicMedicationTreatment{
			{
				MedicationCode: ptr(AnalgesicMedicationCodeK),
				VolumeDose:     "05",
				Route:          MedicationRouteCodeR3,
				DateLocal:      "2027OCT13",
				TimeLocal:      "1551",
			},
		},
		AntibioticMedicationTreatments: []*AntibioticMedicationTreatment{
			{
				MedicationCode: ptr(AntibioticMedicationCodeM),
				VolumeDose:     "10",
				Route:          MedicationRouteCodeR5,
				DateLocal:      "2027OCT13",
				TimeLocal:      "1552",
			},
		},
		OtherMedicationTreatments: []*OtherMedicationTreatment{
			{
				MedicationCode: ptr(OtherMedicationCodeT),
				VolumeDose:     "10",
				Route:          MedicationRouteCodeR4,
				DateLocal:      "2027OCT13",
				TimeLocal:      "1553",
			},
		},
		CasualtyType:             CasualtyTypeCodeA,
		FirstResponderZapOrEdiPi: "1010919789",
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate health intervention segment: %v", err)
	}
}

func TestHealthInterventionSegmentValidateRejectsMissingCoreFields(t *testing.T) {
	seg := HealthInterventionSegment{
		WoundTreatments:          []WoundTreatmentCode{WoundTreatmentCode(99)},
		AirwayTreatment:          AirwayTreatmentCodeA1,
		BreathingTreatment:       BreathingTreatmentCodeB1,
		CasualtyType:             CasualtyTypeCodeA,
		FirstResponderZapOrEdiPi: "1010919789",
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected health intervention validation error")
	}
}

func TestParseRequestPriorityCodeRejectsUnknownValue(t *testing.T) {
	_, err := ParseRequestPriorityCode("16")
	if err == nil {
		t.Fatalf("expected invalid priority code error")
	}
}

func TestSupplyAttachmentIndicatorRejectsUnknownValue(t *testing.T) {
	seg := SupplyRequestSegment{
		SegmentNumber:             1,
		RequestTypeCode:           SupplyRequestTypeCodeSR,
		RequestPriority:           RequestPriorityCode02,
		ItemByNiin:                stringPtr("015519434"),
		ItemQuantity:              "25",
		RequiredDeliveryDateLocal: "2027OCT21",
		DeliveryLocation:          "4QFJ123456",
		AttachmentIndicator:       attachmentPtr(AttachmentIndicatorCode(99)),
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid attachment indicator error")
	}
}

func TestMaintenanceSupportTypeRejectsUnknownValue(t *testing.T) {
	seg := MaintenanceRequestSegment{
		SegmentNumber:                       1,
		RequestTypeCode:                     MaintenanceRequestTypeCodeCM,
		RequestPriority:                     RequestPriorityCode02,
		Niin:                                "015519434",
		NumberOfPiecesRequiringSupport:      "1",
		EquipmentOperationalCondition:       MaintenanceOperationalConditionCodeC,
		DateMaintenanceSupportRequiredLocal: "2027OCT21",
		LocationOfEquipment:                 "4QFJ123456",
		TypeOfMaintenanceSupportRequired:    MaintenanceSupportTypeCode(99),
		TypeOfRepair:                        MaintenanceRepairTypeCodeD1,
		RepairMajorDefect:                   MaintenanceMajorDefectCodeMD07,
		AttachmentIndicator:                 AttachmentIndicatorCodeNo,
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid maintenance support type error")
	}
}

func TestMaintenanceRepairTypeRejectsUnknownValue(t *testing.T) {
	seg := MaintenanceRequestSegment{
		SegmentNumber:                       1,
		RequestTypeCode:                     MaintenanceRequestTypeCodeCM,
		RequestPriority:                     RequestPriorityCode02,
		Niin:                                "015519434",
		NumberOfPiecesRequiringSupport:      "1",
		EquipmentOperationalCondition:       MaintenanceOperationalConditionCodeC,
		DateMaintenanceSupportRequiredLocal: "2027OCT21",
		LocationOfEquipment:                 "4QFJ123456",
		TypeOfMaintenanceSupportRequired:    MaintenanceSupportTypeCodeR1,
		TypeOfRepair:                        MaintenanceRepairTypeCode(99),
		RepairMajorDefect:                   MaintenanceMajorDefectCodeMD07,
		AttachmentIndicator:                 AttachmentIndicatorCodeNo,
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid maintenance repair type error")
	}
}

func TestMaintenanceMajorDefectRejectsUnknownValue(t *testing.T) {
	seg := MaintenanceRequestSegment{
		SegmentNumber:                       1,
		RequestTypeCode:                     MaintenanceRequestTypeCodeCM,
		RequestPriority:                     RequestPriorityCode02,
		Niin:                                "015519434",
		NumberOfPiecesRequiringSupport:      "1",
		EquipmentOperationalCondition:       MaintenanceOperationalConditionCodeC,
		DateMaintenanceSupportRequiredLocal: "2027OCT21",
		LocationOfEquipment:                 "4QFJ123456",
		TypeOfMaintenanceSupportRequired:    MaintenanceSupportTypeCodeR1,
		TypeOfRepair:                        MaintenanceRepairTypeCodeD1,
		RepairMajorDefect:                   MaintenanceMajorDefectCode(99),
		AttachmentIndicator:                 AttachmentIndicatorCodeNo,
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid maintenance major defect error")
	}
}

func TestEngineerRoadReportValidate(t *testing.T) {
	seg := EngineerReconRoadReportSegment{
		DateOfEvaluationLocal: "2027OCT21",
		StartPointLocation:    "4QFJ123456",
		EndPointLocation:      "4QFJ456789",
		RoadClassification:    RoadClassificationCodeA,
		Drainage:              RoadDrainageCodeA,
		Foundation:            RoadFoundationCodeA,
		SurfaceType:           RoadSurfaceTypeCodeB,
		Obstructions:          "B",
		AttachmentIndicator:   attachmentPtr(AttachmentIndicatorCodeNo),
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate road report: %v", err)
	}
}

func TestEngineerRoadReportRejectsUnknownRoadCodes(t *testing.T) {
	seg := EngineerReconRoadReportSegment{
		DateOfEvaluationLocal: "2027OCT21",
		StartPointLocation:    "4QFJ123456",
		EndPointLocation:      "4QFJ456789",
		RoadClassification:    RoadClassificationCode(99),
		Drainage:              RoadDrainageCodeA,
		Foundation:            RoadFoundationCodeA,
		SurfaceType:           RoadSurfaceTypeCodeB,
		Obstructions:          "B",
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid road classification error")
	}
}

func TestRequestContainerAcceptsEngineerRoadSegment(t *testing.T) {
	header := testHeader()
	road := &EngineerReconRoadReportSegment{
		DateOfEvaluationLocal: "2027OCT21",
		StartPointLocation:    "4QFJ123456",
		EndPointLocation:      "4QFJ456789",
		RoadClassification:    RoadClassificationCodeA,
		Drainage:              RoadDrainageCodeA,
		Foundation:            RoadFoundationCodeA,
		SurfaceType:           RoadSurfaceTypeCodeB,
		Obstructions:          "B",
		AttachmentIndicator:   attachmentPtr(AttachmentIndicatorCodeNo),
	}
	container := RequestContainer{
		Header:   header,
		Segments: []*RequestSegment{wrapEngineerRoad(road)},
	}

	if err := container.Validate(); err != nil {
		t.Fatalf("validate container with road segment: %v", err)
	}
}

func TestEngineerLandingZoneReportValidate(t *testing.T) {
	seg := EngineerReconLandingZoneReportSegment{
		DateOfEvaluationLocal:   "2027OCT21",
		Location:                "4QFJ123456",
		Estimate:                EstimateCodeNo,
		LayoutDesignation:       LandingZoneLayoutDesignationCodeLZ,
		LandingPointCapacity:    1,
		LandingZoneCapacity:     2,
		LandingSiteCapacity:     3,
		LandingZoneWidthFeet:    "120",
		LandingZoneLengthFeet:   "900",
		AircraftSupportability:  AircraftSupportabilityCodeC,
		LandingZoneApproach:     CardinalDirectionCodeN,
		LandingZoneDeparture:    CardinalDirectionCodeS,
		LandingZoneSurfaceSlope: LandingZoneSurfaceSlopeCodeA,
		Obstacle:                LandingZoneObstacleCode1,
		AttachmentIndicator:     attachmentPtr(AttachmentIndicatorCodeNo),
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate landing zone report: %v", err)
	}
}

func TestEngineerLandingZoneRejectsUnknownCodes(t *testing.T) {
	seg := EngineerReconLandingZoneReportSegment{
		DateOfEvaluationLocal:   "2027OCT21",
		Location:                "4QFJ123456",
		Estimate:                EstimateCode(99),
		LayoutDesignation:       LandingZoneLayoutDesignationCodeLZ,
		LandingPointCapacity:    1,
		LandingZoneCapacity:     2,
		LandingSiteCapacity:     3,
		LandingZoneWidthFeet:    "120",
		LandingZoneLengthFeet:   "900",
		AircraftSupportability:  AircraftSupportabilityCodeC,
		LandingZoneApproach:     CardinalDirectionCodeN,
		LandingZoneDeparture:    CardinalDirectionCodeS,
		LandingZoneSurfaceSlope: LandingZoneSurfaceSlopeCodeA,
		Obstacle:                LandingZoneObstacleCode1,
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid landing zone estimate error")
	}
}

func TestRequestContainerAcceptsEngineerLandingZoneSegment(t *testing.T) {
	header := testHeader()
	lz := &EngineerReconLandingZoneReportSegment{
		DateOfEvaluationLocal:   "2027OCT21",
		Location:                "4QFJ123456",
		Estimate:                EstimateCodeNo,
		LayoutDesignation:       LandingZoneLayoutDesignationCodeLZ,
		LandingPointCapacity:    1,
		LandingZoneCapacity:     2,
		LandingSiteCapacity:     3,
		LandingZoneWidthFeet:    "120",
		LandingZoneLengthFeet:   "900",
		AircraftSupportability:  AircraftSupportabilityCodeC,
		LandingZoneApproach:     CardinalDirectionCodeN,
		LandingZoneDeparture:    CardinalDirectionCodeS,
		LandingZoneSurfaceSlope: LandingZoneSurfaceSlopeCodeA,
		Obstacle:                LandingZoneObstacleCode1,
		AttachmentIndicator:     attachmentPtr(AttachmentIndicatorCodeNo),
	}
	container := RequestContainer{
		Header:   header,
		Segments: []*RequestSegment{wrapEngineerLZ(lz)},
	}

	if err := container.Validate(); err != nil {
		t.Fatalf(
			"validate container with landing zone segment: %v",
			err,
		)
	}
}

func TestObstacleRemovalSegmentValidate(t *testing.T) {
	seg := GeneralEngineeringObstacleRemovalSegment{
		DateOfEvaluationLocal: "2027OCT21",
		Location:              "4QFJ123456",
		Obstacle:              ObstacleActionCode1,
		MinMaxLengthFeet:      "00010/00020",
		MinMaxWidthFeet:       "00005/00010",
		MinMaxDepthFeet:       "00002/00005",
		RouteNumber:           "101",
		DeterminationOfAction: ObstacleDeterminationCode2,
		Bypass:                BypassCodeYes,
		BypassGrid:            "4QFJ456789",
		AttachmentIndicator:   attachmentPtr(AttachmentIndicatorCodeNo),
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate obstacle removal: %v", err)
	}
}

func TestObstacleRemovalRejectsUnknownCodes(t *testing.T) {
	seg := GeneralEngineeringObstacleRemovalSegment{
		DateOfEvaluationLocal: "2027OCT21",
		Location:              "4QFJ123456",
		Obstacle:              ObstacleActionCode(99),
		MinMaxLengthFeet:      "00010/00020",
		MinMaxWidthFeet:       "00005/00010",
		MinMaxDepthFeet:       "00002/00005",
		RouteNumber:           "101",
		DeterminationOfAction: ObstacleDeterminationCode2,
		Bypass:                BypassCodeYes,
		BypassGrid:            "4QFJ456789",
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid obstacle action error")
	}
}

func TestRequestContainerAcceptsObstacleRemovalSegment(t *testing.T) {
	header := testHeader()
	obstacle := &GeneralEngineeringObstacleRemovalSegment{
		DateOfEvaluationLocal: "2027OCT21",
		Location:              "4QFJ123456",
		Obstacle:              ObstacleActionCode1,
		MinMaxLengthFeet:      "00010/00020",
		MinMaxWidthFeet:       "00005/00010",
		MinMaxDepthFeet:       "00002/00005",
		RouteNumber:           "101",
		DeterminationOfAction: ObstacleDeterminationCode2,
		Bypass:                BypassCodeYes,
		BypassGrid:            "4QFJ456789",
		AttachmentIndicator:   attachmentPtr(AttachmentIndicatorCodeNo),
	}
	container := RequestContainer{
		Header:   header,
		Segments: []*RequestSegment{wrapObstacleRemoval(obstacle)},
	}

	if err := container.Validate(); err != nil {
		t.Fatalf(
			"validate container with obstacle removal segment: %v",
			err,
		)
	}
}

func TestEODSegmentValidate(t *testing.T) {
	seg := ExplosiveOrdnanceDisposalSegment{
		DateOfUxoDiscovery:            "2027OCT21",
		RequestedDateOfEodAction:      "2027OCT22",
		LocationOfUxo:                 "4QFJ123456",
		TypeOfCbrnAgent:               CBRNAgentTypeCode1,
		PhysicalPropertyOfCbrnAgent:   cbrnPhysicalPropertyPtr(CBRNPhysicalPropertyCode2),
		ContaminationValueOfCbrnAgent: cbrnContaminationPtr(CBRNContaminationValueCodeE),
		MunitionColor:                 "R|RED",
		MunitionMarkings:              "TRAINING ROUND",
		MunitionPurpose:               MunitionPurposeCodeAA,
		MunitionType:                  MunitionTypeCodeE,
		AttachmentIndicator:           attachmentPtr(AttachmentIndicatorCodeNo),
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate eod segment: %v", err)
	}
}

func TestEODRejectsUnknownCodes(t *testing.T) {
	seg := ExplosiveOrdnanceDisposalSegment{
		DateOfUxoDiscovery:       "2027OCT21",
		RequestedDateOfEodAction: "2027OCT22",
		LocationOfUxo:            "4QFJ123456",
		TypeOfCbrnAgent:          CBRNAgentTypeCode(9),
		MunitionColor:            "R|RED",
		MunitionMarkings:         "TRAINING ROUND",
		MunitionPurpose:          MunitionPurposeCodeAA,
		MunitionType:             MunitionTypeCodeE,
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid cbrn agent code error")
	}
}

func TestRequestContainerAcceptsEODSegment(t *testing.T) {
	header := testHeader()
	eod := &ExplosiveOrdnanceDisposalSegment{
		DateOfUxoDiscovery:       "2027OCT21",
		RequestedDateOfEodAction: "2027OCT22",
		LocationOfUxo:            "4QFJ123456",
		TypeOfCbrnAgent:          CBRNAgentTypeCode1,
		MunitionColor:            "R|RED",
		MunitionMarkings:         "TRAINING ROUND",
		MunitionPurpose:          MunitionPurposeCodeAA,
		MunitionType:             MunitionTypeCodeE,
		AttachmentIndicator:      attachmentPtr(AttachmentIndicatorCodeNo),
	}
	container := RequestContainer{
		Header:   header,
		Segments: []*RequestSegment{wrapEOD(eod)},
	}

	if err := container.Validate(); err != nil {
		t.Fatalf("validate container with eod segment: %v", err)
	}
}

func TestHealthCollectionSegmentValidate(t *testing.T) {
	seg := HealthCollectionSegment{
		SegmentNumber:                       1,
		RequestTypeCode:                     HealthCollectionRequestTypeCodeCR,
		RequestPriority:                     RequestPriorityCode02,
		ZapOrEdiPi:                          "1010919789",
		LastName:                            "SMITH",
		FirstName:                           "CHRIS",
		Service:                             servicePtr(ServiceCodeUSMC),
		ElementUnitIdentificationOrCallsign: stringPtr("KL9K"),
		Allergies:                           "NKDA",
		DateOfInjuryLocal:                   "2027OCT21",
		TimeOfInjuryLocal:                   "1547",
		LocationInjuryOccurred:              "4QFJ123456",
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate health collection segment: %v", err)
	}
}

func TestHealthCollectionRejectsUnknownService(t *testing.T) {
	seg := HealthCollectionSegment{
		SegmentNumber:          1,
		RequestTypeCode:        HealthCollectionRequestTypeCodeCR,
		RequestPriority:        RequestPriorityCode02,
		ZapOrEdiPi:             "1010919789",
		LastName:               "SMITH",
		FirstName:              "CHRIS",
		Service:                servicePtr(ServiceCode(99)),
		Allergies:              "NKDA",
		DateOfInjuryLocal:      "2027OCT21",
		TimeOfInjuryLocal:      "1547",
		LocationInjuryOccurred: "4QFJ123456",
	}

	if err := seg.Validate(); err == nil {
		t.Fatalf("expected invalid health collection service error")
	}
}

func TestRequestContainerAcceptsHealthCollectionSegment(t *testing.T) {
	header := testHeader()
	health := &HealthCollectionSegment{
		SegmentNumber:          1,
		RequestTypeCode:        HealthCollectionRequestTypeCodeCR,
		RequestPriority:        RequestPriorityCode02,
		ZapOrEdiPi:             "1010919789",
		LastName:               "SMITH",
		FirstName:              "CHRIS",
		Service:                servicePtr(ServiceCodeUSMC),
		Allergies:              "NKDA",
		DateOfInjuryLocal:      "2027OCT21",
		TimeOfInjuryLocal:      "1547",
		LocationInjuryOccurred: "4QFJ123456",
	}
	container := RequestContainer{
		Header:   header,
		Segments: []*RequestSegment{wrapHealthCollection(health)},
	}

	if err := container.Validate(); err != nil {
		t.Fatalf(
			"validate container with health collection segment: %v",
			err,
		)
	}
}

func TestHealthHoldSegmentValidate(t *testing.T) {
	painScale := uint32(4)
	hold := &HealthHoldSegment{
		TriageEntries: []*HealthTriageSegment{
			{
				PrimaryMechanismOfInjury: HealthPrimaryMechanismCodeP4,
				CbrnRelatedExposure:      HealthCBRNExposureCodeX,
				MajorSignsSymptoms:       HealthMajorSignsSymptomsCodeB,
				InjuryLocations:          []string{"RA0"},
				VitalsCheckDateLocal:     stringPtr("2027OCT21"),
				VitalsCheckTimeLocal:     stringPtr("1600"),
				PulseRate:                stringPtr("090"),
				PulseLocation:            HealthPulseLocationCodeW.Enum(),
				Responsiveness:           HealthResponsivenessCodeA.Enum(),
				PainScale:                &painScale,
				TriagePrecedence:         HealthTriagePrecedenceCodeB,
			},
		},
		InterventionEntries: []*HealthInterventionSegment{
			{
				WoundTreatments:          []WoundTreatmentCode{WoundTreatmentCodeT2},
				AirwayTreatment:          AirwayTreatmentCodeA0,
				BreathingTreatment:       BreathingTreatmentCodeB0,
				CasualtyType:             CasualtyTypeCodeA,
				FirstResponderZapOrEdiPi: "1010919789",
			},
		},
	}

	if err := hold.Validate(); err != nil {
		t.Fatalf("validate health hold segment: %v", err)
	}
}

func TestHealthHoldSegmentRejectsEmptyAppendSet(t *testing.T) {
	hold := &HealthHoldSegment{}

	if err := hold.Validate(); err == nil {
		t.Fatalf("expected empty health hold validation error")
	}
}

func TestRequestContainerAcceptsHealthHoldSegment(t *testing.T) {
	header := testHeader()
	hold := &HealthHoldSegment{
		TriageEntries: []*HealthTriageSegment{
			{
				PrimaryMechanismOfInjury: HealthPrimaryMechanismCodeE1,
				CbrnRelatedExposure:      HealthCBRNExposureCodeX,
				MajorSignsSymptoms:       HealthMajorSignsSymptomsCodeB,
				InjuryLocations:          []string{"FH0"},
				TriagePrecedence:         HealthTriagePrecedenceCodeB,
			},
		},
	}
	container := RequestContainer{
		Header:   header,
		Segments: []*RequestSegment{wrapHealthHold(hold)},
	}

	if err := container.Validate(); err != nil {
		t.Fatalf("validate container with health hold segment: %v", err)
	}
}
