package lxdr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
	if header.LocalRequestID != "3838JBNM5" {
		t.Fatalf("unexpected local request id: %q", header.LocalRequestID)
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
		LocalRequestID:                  "3838JBNM5",
		RequestPriority:                 RequestPriorityCode02,
		ElementUnitIDOrCallsign:         "KL9K",
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

	if resp.LocalRequestID != "3838JBNM5" {
		t.Fatalf("unexpected local request id: %q", resp.LocalRequestID)
	}
	if resp.SynchronizedRequestID != "KL9K15474QFJ" {
		t.Fatalf(
			"unexpected synchronized request id: %q",
			resp.SynchronizedRequestID,
		)
	}
}

func TestRenderSynchronizedResponseExample(t *testing.T) {
	resp := SynchronizedResponse{
		LocalRequestID:        "3838JBNM5",
		SynchronizedRequestID: "KL9K15474QFJ",
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
	if seg.ZAPOrEDIPI != "1010919789" {
		t.Fatalf("unexpected personnel identifier: %q", seg.ZAPOrEDIPI)
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
		SegmentNumber:              1,
		RequestTypeCode:            MobilityPaxRequestTypeCodePM,
		RequestPriority:            RequestPriorityCode02,
		ZAPOrEDIPI:                 "1010919789",
		EarliestDepartureDateLocal: "2027OCT15",
		LatestDepartureDateLocal:   "2027OCT20",
		DepartureLocation:          "4QFJ123456",
		DestinationLocation:        "4QFJ456789",
		TotalEstimatedBaggageLBS:   "075",
		HazardousMaterialType:      "X",
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
	if seg.ItemByNIIN != "015519434" {
		t.Fatalf("unexpected niin: %q", seg.ItemByNIIN)
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
		ItemByNIIN:                 "015519434",
		ItemQuantity:               "1",
		SerialNumber:               "598742",
		GrossWeightLBS:             "28000",
		ActualHeightInches:         "126",
		ActualWidthInches:          "100",
		ActualLengthInches:         "315",
		HMIC:                       CargoHMICCodeD,
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
	header := RequestHeader{
		LocalSystemDate:                 "2027OCT13",
		LocalSystemTime:                 "15470352",
		SynchronizedGeospatialReference: "4QFJ123456",
		LocalRequestID:                  "3838JBNM5",
		RequestPriority:                 RequestPriorityCode02,
		ElementUnitIDOrCallsign:         "KL9K",
		RequestSegmentCount:             1,
	}
	pax := &MobilityPaxRequestSegment{
		SegmentNumber:              1,
		RequestTypeCode:            MobilityPaxRequestTypeCodePM,
		RequestPriority:            RequestPriorityCode02,
		ZAPOrEDIPI:                 "1010919789",
		EarliestDepartureDateLocal: "2027OCT15",
		LatestDepartureDateLocal:   "2027OCT20",
		DepartureLocation:          "4QFJ123456",
		DestinationLocation:        "4QFJ456789",
		TotalEstimatedBaggageLBS:   "075",
		HazardousMaterialType:      "X",
	}
	container := RequestContainer{
		Header: header,
		Segments: []RequestSegment{
			{
				FunctionFamily: FunctionFamilyMobility,
				MobilityPax:    pax,
			},
		},
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
		ItemByNIIN:                 "015519434",
		ItemQuantity:               "1",
		SerialNumber:               "598742",
		GrossWeightLBS:             "28000",
		ActualHeightInches:         "126",
		ActualWidthInches:          "100",
		ActualLengthInches:         "315",
		HMIC:                       CargoHMICCode("Z"),
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
		ItemByNIIN:                 "015519434",
		ItemQuantity:               "1",
		SerialNumber:               "598742",
		GrossWeightLBS:             "28000",
		ActualHeightInches:         "126",
		ActualWidthInches:          "100",
		ActualLengthInches:         "315",
		HMIC:                       CargoHMICCodeD,
		Handling:                   CargoHandlingCode("Z"),
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
		ItemByNIIN:                "",
		ItemQuantity:              "25",
		RequiredDeliveryDateLocal: "2027OCT21",
		DeliveryLocation:          "4QFJ123456",
		AttachmentIndicator:       AttachmentIndicatorCodeUnspecified,
		Narrative:                 "BULK WATER",
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
		SerialNumber:                        "",
		NIIN:                                "015519434",
		ModelOfEquipment:                    "",
		ItemNomenclature:                    "",
		NumberOfPiecesRequiringSupport:      "1",
		EquipmentOperationalCondition:       MaintenanceOperationalConditionCodeC,
		DateMaintenanceSupportRequiredLocal: "2027OCT21",
		LocationOfEquipment:                 "4QFJ123456",
		TypeOfMaintenanceSupportRequired:    MaintenanceSupportTypeCodeR1,
		TypeOfRepair:                        MaintenanceRepairTypeCodeD1,
		RepairMajorDefect:                   MaintenanceMajorDefectCodeMD07,
		AttachmentIndicator:                 AttachmentIndicatorCodeNo,
		Narrative:                           "",
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
		NIIN:                                "015519434",
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
		ItemByNIIN:                "015519434",
		ItemQuantity:              "25",
		RequiredDeliveryDateLocal: "2027OCT21",
		DeliveryLocation:          "4QFJ123456",
		AttachmentIndicator:       AttachmentIndicatorCode("2"),
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
		NIIN:                                "015519434",
		NumberOfPiecesRequiringSupport:      "1",
		EquipmentOperationalCondition:       MaintenanceOperationalConditionCodeC,
		DateMaintenanceSupportRequiredLocal: "2027OCT21",
		LocationOfEquipment:                 "4QFJ123456",
		TypeOfMaintenanceSupportRequired:    MaintenanceSupportTypeCode("R9"),
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
		NIIN:                                "015519434",
		NumberOfPiecesRequiringSupport:      "1",
		EquipmentOperationalCondition:       MaintenanceOperationalConditionCodeC,
		DateMaintenanceSupportRequiredLocal: "2027OCT21",
		LocationOfEquipment:                 "4QFJ123456",
		TypeOfMaintenanceSupportRequired:    MaintenanceSupportTypeCodeR1,
		TypeOfRepair:                        MaintenanceRepairTypeCode("Z9"),
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
		NIIN:                                "015519434",
		NumberOfPiecesRequiringSupport:      "1",
		EquipmentOperationalCondition:       MaintenanceOperationalConditionCodeC,
		DateMaintenanceSupportRequiredLocal: "2027OCT21",
		LocationOfEquipment:                 "4QFJ123456",
		TypeOfMaintenanceSupportRequired:    MaintenanceSupportTypeCodeR1,
		TypeOfRepair:                        MaintenanceRepairTypeCodeD1,
		RepairMajorDefect:                   MaintenanceMajorDefectCode("MD99"),
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
		AttachmentIndicator:   AttachmentIndicatorCodeNo,
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
		RoadClassification:    RoadClassificationCode("Z"),
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
	header := RequestHeader{
		LocalSystemDate:                 "2027OCT13",
		LocalSystemTime:                 "15470352",
		SynchronizedGeospatialReference: "4QFJ123456",
		LocalRequestID:                  "3838JBNM5",
		RequestPriority:                 RequestPriorityCode02,
		ElementUnitIDOrCallsign:         "KL9K",
		RequestSegmentCount:             1,
	}
	road := &EngineerReconRoadReportSegment{
		DateOfEvaluationLocal: "2027OCT21",
		StartPointLocation:    "4QFJ123456",
		EndPointLocation:      "4QFJ456789",
		RoadClassification:    RoadClassificationCodeA,
		Drainage:              RoadDrainageCodeA,
		Foundation:            RoadFoundationCodeA,
		SurfaceType:           RoadSurfaceTypeCodeB,
		Obstructions:          "B",
		AttachmentIndicator:   AttachmentIndicatorCodeNo,
	}
	container := RequestContainer{
		Header: header,
		Segments: []RequestSegment{
			{
				FunctionFamily: FunctionFamilyGeneralEngineering,
				EngineerRoad:   road,
			},
		},
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
		AttachmentIndicator:     AttachmentIndicatorCodeNo,
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate landing zone report: %v", err)
	}
}

func TestEngineerLandingZoneRejectsUnknownCodes(t *testing.T) {
	seg := EngineerReconLandingZoneReportSegment{
		DateOfEvaluationLocal:   "2027OCT21",
		Location:                "4QFJ123456",
		Estimate:                EstimateCode("2"),
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
	header := RequestHeader{
		LocalSystemDate:                 "2027OCT13",
		LocalSystemTime:                 "15470352",
		SynchronizedGeospatialReference: "4QFJ123456",
		LocalRequestID:                  "3838JBNM5",
		RequestPriority:                 RequestPriorityCode02,
		ElementUnitIDOrCallsign:         "KL9K",
		RequestSegmentCount:             1,
	}
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
		AttachmentIndicator:     AttachmentIndicatorCodeNo,
	}
	container := RequestContainer{
		Header: header,
		Segments: []RequestSegment{
			{
				FunctionFamily: FunctionFamilyGeneralEngineering,
				EngineerLZ:     lz,
			},
		},
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
		AttachmentIndicator:   AttachmentIndicatorCodeNo,
	}

	if err := seg.Validate(); err != nil {
		t.Fatalf("validate obstacle removal: %v", err)
	}
}

func TestObstacleRemovalRejectsUnknownCodes(t *testing.T) {
	seg := GeneralEngineeringObstacleRemovalSegment{
		DateOfEvaluationLocal: "2027OCT21",
		Location:              "4QFJ123456",
		Obstacle:              ObstacleActionCode("9"),
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
	header := RequestHeader{
		LocalSystemDate:                 "2027OCT13",
		LocalSystemTime:                 "15470352",
		SynchronizedGeospatialReference: "4QFJ123456",
		LocalRequestID:                  "3838JBNM5",
		RequestPriority:                 RequestPriorityCode02,
		ElementUnitIDOrCallsign:         "KL9K",
		RequestSegmentCount:             1,
	}
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
		AttachmentIndicator:   AttachmentIndicatorCodeNo,
	}
	container := RequestContainer{
		Header: header,
		Segments: []RequestSegment{
			{
				FunctionFamily:  FunctionFamilyGeneralEngineering,
				ObstacleRemoval: obstacle,
			},
		},
	}

	if err := container.Validate(); err != nil {
		t.Fatalf(
			"validate container with obstacle removal segment: %v",
			err,
		)
	}
}
