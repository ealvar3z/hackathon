package lxdr

type FunctionFamily string

const (
	FunctionFamilyUnspecified        FunctionFamily = ""
	FunctionFamilyMobility           FunctionFamily = "mobility"
	FunctionFamilySupply             FunctionFamily = "supply"
	FunctionFamilyMaintenance        FunctionFamily = "maintenance"
	FunctionFamilyGeneralEngineering FunctionFamily = "general_engineering"
	FunctionFamilyHealthServices     FunctionFamily = "health_services"
)

type RequestPriorityCode string

const (
	RequestPriorityCodeUnspecified RequestPriorityCode = ""
	RequestPriorityCode01          RequestPriorityCode = "01"
	RequestPriorityCode02          RequestPriorityCode = "02"
	RequestPriorityCode03          RequestPriorityCode = "03"
	RequestPriorityCode04          RequestPriorityCode = "04"
	RequestPriorityCode05          RequestPriorityCode = "05"
	RequestPriorityCode06          RequestPriorityCode = "06"
	RequestPriorityCode07          RequestPriorityCode = "07"
	RequestPriorityCode08          RequestPriorityCode = "08"
	RequestPriorityCode09          RequestPriorityCode = "09"
	RequestPriorityCode10          RequestPriorityCode = "10"
	RequestPriorityCode11          RequestPriorityCode = "11"
	RequestPriorityCode12          RequestPriorityCode = "12"
	RequestPriorityCode13          RequestPriorityCode = "13"
	RequestPriorityCode14          RequestPriorityCode = "14"
	RequestPriorityCode15          RequestPriorityCode = "15"
)

type MobilityPaxRequestTypeCode string

const (
	MobilityPaxRequestTypeCodeUnspecified MobilityPaxRequestTypeCode = ""
	MobilityPaxRequestTypeCodePM          MobilityPaxRequestTypeCode = "PM"
)

type MobilityCargoRequestTypeCode string

const (
	MobilityCargoRequestTypeCodeUnspecified MobilityCargoRequestTypeCode = ""
	MobilityCargoRequestTypeCodeCM          MobilityCargoRequestTypeCode = "CM"
)

type CargoHMICCode string

const (
	CargoHMICCodeUnspecified CargoHMICCode = ""
	CargoHMICCodeD           CargoHMICCode = "D"
	CargoHMICCodeN           CargoHMICCode = "N"
	CargoHMICCodeP           CargoHMICCode = "P"
	CargoHMICCodeY           CargoHMICCode = "Y"
)

type CargoHandlingCode string

const (
	CargoHandlingCodeUnspecified CargoHandlingCode = ""
	CargoHandlingCodeC           CargoHandlingCode = "C"
	CargoHandlingCodeM           CargoHandlingCode = "M"
	CargoHandlingCodeT           CargoHandlingCode = "T"
	CargoHandlingCodeR           CargoHandlingCode = "R"
	CargoHandlingCodeX           CargoHandlingCode = "X"
)

type SupplyRequestTypeCode string

const (
	SupplyRequestTypeCodeUnspecified SupplyRequestTypeCode = ""
	SupplyRequestTypeCodeSR          SupplyRequestTypeCode = "SR"
)

type MaintenanceRequestTypeCode string

const (
	MaintenanceRequestTypeCodeUnspecified MaintenanceRequestTypeCode = ""
	MaintenanceRequestTypeCodeCM          MaintenanceRequestTypeCode = "CM"
)

type AttachmentIndicatorCode string

const (
	AttachmentIndicatorCodeUnspecified AttachmentIndicatorCode = ""
	AttachmentIndicatorCodeNo          AttachmentIndicatorCode = "0"
	AttachmentIndicatorCodeYes         AttachmentIndicatorCode = "1"
)

type MaintenanceOperationalConditionCode string

const (
	MaintenanceOperationalConditionCodeUnspecified MaintenanceOperationalConditionCode = ""
	MaintenanceOperationalConditionCodeA           MaintenanceOperationalConditionCode = "A"
	MaintenanceOperationalConditionCodeB           MaintenanceOperationalConditionCode = "B"
	MaintenanceOperationalConditionCodeC           MaintenanceOperationalConditionCode = "C"
)

type MaintenanceSupportTypeCode string

const (
	MaintenanceSupportTypeCodeUnspecified MaintenanceSupportTypeCode = ""
	MaintenanceSupportTypeCodeXX          MaintenanceSupportTypeCode = "XX"
	MaintenanceSupportTypeCodeR1          MaintenanceSupportTypeCode = "R1"
	MaintenanceSupportTypeCodeR2          MaintenanceSupportTypeCode = "R2"
	MaintenanceSupportTypeCodeR3          MaintenanceSupportTypeCode = "R3"
	MaintenanceSupportTypeCodeR4          MaintenanceSupportTypeCode = "R4"
)

type MaintenanceRepairTypeCode string

const (
	MaintenanceRepairTypeCodeUnspecified MaintenanceRepairTypeCode = ""
	MaintenanceRepairTypeCodeM1          MaintenanceRepairTypeCode = "M1"
	MaintenanceRepairTypeCodeS1          MaintenanceRepairTypeCode = "S1"
	MaintenanceRepairTypeCodeS2          MaintenanceRepairTypeCode = "S2"
	MaintenanceRepairTypeCodeC1          MaintenanceRepairTypeCode = "C1"
	MaintenanceRepairTypeCodeD1          MaintenanceRepairTypeCode = "D1"
)

type MaintenanceMajorDefectCode string

const (
	MaintenanceMajorDefectCodeUnspecified MaintenanceMajorDefectCode = ""
	MaintenanceMajorDefectCodeMD01        MaintenanceMajorDefectCode = "MD01"
	MaintenanceMajorDefectCodeMD02        MaintenanceMajorDefectCode = "MD02"
	MaintenanceMajorDefectCodeMD03        MaintenanceMajorDefectCode = "MD03"
	MaintenanceMajorDefectCodeMD04        MaintenanceMajorDefectCode = "MD04"
	MaintenanceMajorDefectCodeMD05        MaintenanceMajorDefectCode = "MD05"
	MaintenanceMajorDefectCodeMD06        MaintenanceMajorDefectCode = "MD06"
	MaintenanceMajorDefectCodeMD07        MaintenanceMajorDefectCode = "MD07"
	MaintenanceMajorDefectCodeMD08        MaintenanceMajorDefectCode = "MD08"
	MaintenanceMajorDefectCodeMD09        MaintenanceMajorDefectCode = "MD09"
	MaintenanceMajorDefectCodeMD10        MaintenanceMajorDefectCode = "MD10"
	MaintenanceMajorDefectCodeMD11        MaintenanceMajorDefectCode = "MD11"
	MaintenanceMajorDefectCodeMD12        MaintenanceMajorDefectCode = "MD12"
	MaintenanceMajorDefectCodeMD13        MaintenanceMajorDefectCode = "MD13"
	MaintenanceMajorDefectCodeMD14        MaintenanceMajorDefectCode = "MD14"
	MaintenanceMajorDefectCodeMD15        MaintenanceMajorDefectCode = "MD15"
	MaintenanceMajorDefectCodeMD16        MaintenanceMajorDefectCode = "MD16"
	MaintenanceMajorDefectCodeNMAJ        MaintenanceMajorDefectCode = "NMAJ"
)

type RoadClassificationCode string

const (
	RoadClassificationCodeUnspecified RoadClassificationCode = ""
	RoadClassificationCodeA           RoadClassificationCode = "A"
	RoadClassificationCodeB           RoadClassificationCode = "B"
	RoadClassificationCodeC           RoadClassificationCode = "C"
	RoadClassificationCodeD           RoadClassificationCode = "D"
)

type RoadDrainageCode string

const (
	RoadDrainageCodeUnspecified RoadDrainageCode = ""
	RoadDrainageCodeA           RoadDrainageCode = "A"
	RoadDrainageCodeB           RoadDrainageCode = "B"
)

type RoadFoundationCode string

const (
	RoadFoundationCodeUnspecified RoadFoundationCode = ""
	RoadFoundationCodeA           RoadFoundationCode = "A"
	RoadFoundationCodeB           RoadFoundationCode = "B"
)

type RoadSurfaceTypeCode string

const (
	RoadSurfaceTypeCodeUnspecified RoadSurfaceTypeCode = ""
	RoadSurfaceTypeCodeA           RoadSurfaceTypeCode = "A"
	RoadSurfaceTypeCodeB           RoadSurfaceTypeCode = "B"
	RoadSurfaceTypeCodeC           RoadSurfaceTypeCode = "C"
	RoadSurfaceTypeCodeD           RoadSurfaceTypeCode = "D"
	RoadSurfaceTypeCodeE           RoadSurfaceTypeCode = "E"
	RoadSurfaceTypeCodeF           RoadSurfaceTypeCode = "F"
	RoadSurfaceTypeCodeG           RoadSurfaceTypeCode = "G"
	RoadSurfaceTypeCodeH           RoadSurfaceTypeCode = "H"
	RoadSurfaceTypeCodeI           RoadSurfaceTypeCode = "I"
	RoadSurfaceTypeCodeJ           RoadSurfaceTypeCode = "J"
)

type EstimateCode string

const (
	EstimateCodeUnspecified EstimateCode = ""
	EstimateCodeNo          EstimateCode = "0"
	EstimateCodeYes         EstimateCode = "1"
)

type LandingZoneLayoutDesignationCode string

const (
	LandingZoneLayoutDesignationCodeUnspecified LandingZoneLayoutDesignationCode = ""
	LandingZoneLayoutDesignationCodeLZ          LandingZoneLayoutDesignationCode = "LZ"
	LandingZoneLayoutDesignationCodeLS          LandingZoneLayoutDesignationCode = "LS"
	LandingZoneLayoutDesignationCodeLP          LandingZoneLayoutDesignationCode = "LP"
)

type AircraftSupportabilityCode string

const (
	AircraftSupportabilityCodeUnspecified AircraftSupportabilityCode = ""
	AircraftSupportabilityCodeA           AircraftSupportabilityCode = "A"
	AircraftSupportabilityCodeB           AircraftSupportabilityCode = "B"
	AircraftSupportabilityCodeC           AircraftSupportabilityCode = "C"
	AircraftSupportabilityCodeD           AircraftSupportabilityCode = "D"
	AircraftSupportabilityCodeE           AircraftSupportabilityCode = "E"
	AircraftSupportabilityCodeF           AircraftSupportabilityCode = "F"
)

type CardinalDirectionCode string

const (
	CardinalDirectionCodeUnspecified CardinalDirectionCode = ""
	CardinalDirectionCodeS           CardinalDirectionCode = "S"
	CardinalDirectionCodeE           CardinalDirectionCode = "E"
	CardinalDirectionCodeN           CardinalDirectionCode = "N"
	CardinalDirectionCodeW           CardinalDirectionCode = "W"
)

type LandingZoneSurfaceSlopeCode string

const (
	LandingZoneSurfaceSlopeCodeUnspecified LandingZoneSurfaceSlopeCode = ""
	LandingZoneSurfaceSlopeCodeA           LandingZoneSurfaceSlopeCode = "A"
	LandingZoneSurfaceSlopeCodeB           LandingZoneSurfaceSlopeCode = "B"
)

type LandingZoneObstacleCode string

const (
	LandingZoneObstacleCodeUnspecified LandingZoneObstacleCode = ""
	LandingZoneObstacleCode1           LandingZoneObstacleCode = "1"
	LandingZoneObstacleCode2           LandingZoneObstacleCode = "2"
	LandingZoneObstacleCode3           LandingZoneObstacleCode = "3"
)

type ObstacleActionCode string

const (
	ObstacleActionCodeUnspecified ObstacleActionCode = ""
	ObstacleActionCode1           ObstacleActionCode = "1"
	ObstacleActionCode2           ObstacleActionCode = "2"
	ObstacleActionCode3           ObstacleActionCode = "3"
)

type ObstacleDeterminationCode string

const (
	ObstacleDeterminationCodeUnspecified ObstacleDeterminationCode = ""
	ObstacleDeterminationCode1           ObstacleDeterminationCode = "1"
	ObstacleDeterminationCode2           ObstacleDeterminationCode = "2"
	ObstacleDeterminationCode3           ObstacleDeterminationCode = "3"
	ObstacleDeterminationCode4           ObstacleDeterminationCode = "4"
)

type BypassCode string

const (
	BypassCodeUnspecified BypassCode = ""
	BypassCodeNo          BypassCode = "0"
	BypassCodeYes         BypassCode = "1"
)

type RequestHeader struct {
	LocalSystemDate                 string
	LocalSystemTime                 string
	UTCTime                         string
	MilitaryDTG                     string
	SynchronizedGeospatialReference string
	LocalRequestID                  string
	SynchronizedRequestID           string
	RequestPriority                 RequestPriorityCode
	ElementUnitIDOrCallsign         string
	RequestSegmentCount             uint8
}

type SynchronizedResponse struct {
	LocalRequestID        string
	SynchronizedRequestID string
}

type RequestContainer struct {
	Header   RequestHeader
	Segments []RequestSegment
}

type RequestSegment struct {
	FunctionFamily  FunctionFamily
	MobilityPax     *MobilityPaxRequestSegment
	MobilityCargo   *MobilityCargoRequestSegment
	Supply          *SupplyRequestSegment
	Maintenance     *MaintenanceRequestSegment
	EngineerRoad    *EngineerReconRoadReportSegment
	EngineerLZ      *EngineerReconLandingZoneReportSegment
	ObstacleRemoval *GeneralEngineeringObstacleRemovalSegment
}

type MobilityPaxRequestSegment struct {
	SegmentNumber              uint8
	RequestTypeCode            MobilityPaxRequestTypeCode
	RequestPriority            RequestPriorityCode
	ZAPOrEDIPI                 string
	EarliestDepartureDateLocal string
	LatestDepartureDateLocal   string
	DepartureLocation          string
	DestinationLocation        string
	TotalEstimatedBaggageLBS   string
	HazardousMaterialType      string
	PersonCount                *uint32
	TotalBaggageWeightLBS      string
}

type MobilityCargoRequestSegment struct {
	SegmentNumber              uint8
	RequestTypeCode            MobilityCargoRequestTypeCode
	RequestPriority            RequestPriorityCode
	ItemByNIIN                 string
	ItemQuantity               string
	SerialNumber               string
	GrossWeightLBS             string
	ActualHeightInches         string
	ActualWidthInches          string
	ActualLengthInches         string
	HMIC                       CargoHMICCode
	Handling                   CargoHandlingCode
	EarliestDepartureDateLocal string
	LatestDepartureDateLocal   string
	DepartureLocation          string
	DestinationLocation        string
	ItemCount                  *uint32
	TotalWeightLBS             string
}

type SupplyRequestSegment struct {
	SegmentNumber             uint8
	RequestTypeCode           SupplyRequestTypeCode
	RequestPriority           RequestPriorityCode
	ItemByNIIN                string
	ItemQuantity              string
	RequiredDeliveryDateLocal string
	DeliveryLocation          string
	AttachmentIndicator       AttachmentIndicatorCode
	Narrative                 string
}

type MaintenanceRequestSegment struct {
	SegmentNumber                       uint8
	RequestTypeCode                     MaintenanceRequestTypeCode
	RequestPriority                     RequestPriorityCode
	SerialNumber                        string
	NIIN                                string
	ModelOfEquipment                    string
	ItemNomenclature                    string
	NumberOfPiecesRequiringSupport      string
	EquipmentOperationalCondition       MaintenanceOperationalConditionCode
	DateMaintenanceSupportRequiredLocal string
	LocationOfEquipment                 string
	TypeOfMaintenanceSupportRequired    MaintenanceSupportTypeCode
	TypeOfRepair                        MaintenanceRepairTypeCode
	RepairMajorDefect                   MaintenanceMajorDefectCode
	AttachmentIndicator                 AttachmentIndicatorCode
	Narrative                           string
}

type EngineerReconRoadReportSegment struct {
	DateOfEvaluationLocal string
	StartPointLocation    string
	EndPointLocation      string
	RoadClassification    RoadClassificationCode
	Drainage              RoadDrainageCode
	Foundation            RoadFoundationCode
	SurfaceType           RoadSurfaceTypeCode
	Obstructions          string
	AttachmentIndicator   AttachmentIndicatorCode
	Narrative             string
}

type EngineerReconLandingZoneReportSegment struct {
	DateOfEvaluationLocal   string
	Location                string
	Estimate                EstimateCode
	LayoutDesignation       LandingZoneLayoutDesignationCode
	LandingPointCapacity    uint8
	LandingZoneCapacity     uint8
	LandingSiteCapacity     uint8
	LandingZoneWidthFeet    string
	LandingZoneLengthFeet   string
	AircraftSupportability  AircraftSupportabilityCode
	LandingZoneApproach     CardinalDirectionCode
	LandingZoneDeparture    CardinalDirectionCode
	LandingZoneSurfaceSlope LandingZoneSurfaceSlopeCode
	Obstacle                LandingZoneObstacleCode
	AttachmentIndicator     AttachmentIndicatorCode
	Narrative               string
}

type GeneralEngineeringObstacleRemovalSegment struct {
	DateOfEvaluationLocal string
	Location              string
	Obstacle              ObstacleActionCode
	MinMaxLengthFeet      string
	MinMaxWidthFeet       string
	MinMaxDepthFeet       string
	RouteNumber           string
	DeterminationOfAction ObstacleDeterminationCode
	Bypass                BypassCode
	BypassGrid            string
	AttachmentIndicator   AttachmentIndicatorCode
	Narrative             string
}
