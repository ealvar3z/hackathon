package lxdr

const (
	FunctionFamilyUnspecified        = FunctionFamily_FUNCTION_FAMILY_UNSPECIFIED
	FunctionFamilyMobility           = FunctionFamily_FUNCTION_FAMILY_MOBILITY
	FunctionFamilySupply             = FunctionFamily_FUNCTION_FAMILY_SUPPLY
	FunctionFamilyMaintenance        = FunctionFamily_FUNCTION_FAMILY_MAINTENANCE
	FunctionFamilyGeneralEngineering = FunctionFamily_FUNCTION_FAMILY_GENERAL_ENGINEERING
	FunctionFamilyHealthServices     = FunctionFamily_FUNCTION_FAMILY_HEALTH_SERVICES
)

const (
	RequestPriorityCodeUnspecified = RequestPriorityCode_REQUEST_PRIORITY_CODE_UNSPECIFIED
	RequestPriorityCode01          = RequestPriorityCode_REQUEST_PRIORITY_CODE_01
	RequestPriorityCode02          = RequestPriorityCode_REQUEST_PRIORITY_CODE_02
	RequestPriorityCode03          = RequestPriorityCode_REQUEST_PRIORITY_CODE_03
	RequestPriorityCode04          = RequestPriorityCode_REQUEST_PRIORITY_CODE_04
	RequestPriorityCode05          = RequestPriorityCode_REQUEST_PRIORITY_CODE_05
	RequestPriorityCode06          = RequestPriorityCode_REQUEST_PRIORITY_CODE_06
	RequestPriorityCode07          = RequestPriorityCode_REQUEST_PRIORITY_CODE_07
	RequestPriorityCode08          = RequestPriorityCode_REQUEST_PRIORITY_CODE_08
	RequestPriorityCode09          = RequestPriorityCode_REQUEST_PRIORITY_CODE_09
	RequestPriorityCode10          = RequestPriorityCode_REQUEST_PRIORITY_CODE_10
	RequestPriorityCode11          = RequestPriorityCode_REQUEST_PRIORITY_CODE_11
	RequestPriorityCode12          = RequestPriorityCode_REQUEST_PRIORITY_CODE_12
	RequestPriorityCode13          = RequestPriorityCode_REQUEST_PRIORITY_CODE_13
	RequestPriorityCode14          = RequestPriorityCode_REQUEST_PRIORITY_CODE_14
	RequestPriorityCode15          = RequestPriorityCode_REQUEST_PRIORITY_CODE_15
)

const (
	MobilityPaxRequestTypeCodeUnspecified      = MobilityPaxRequestTypeCode_MOBILITY_PAX_REQUEST_TYPE_CODE_UNSPECIFIED
	MobilityPaxRequestTypeCodePM               = MobilityPaxRequestTypeCode_MOBILITY_PAX_REQUEST_TYPE_CODE_PM
	MobilityCargoRequestTypeCodeUnspecified    = MobilityCargoRequestTypeCode_MOBILITY_CARGO_REQUEST_TYPE_CODE_UNSPECIFIED
	MobilityCargoRequestTypeCodeCM             = MobilityCargoRequestTypeCode_MOBILITY_CARGO_REQUEST_TYPE_CODE_CM
	SupplyRequestTypeCodeUnspecified           = SupplyRequestTypeCode_SUPPLY_REQUEST_TYPE_CODE_UNSPECIFIED
	SupplyRequestTypeCodeSR                    = SupplyRequestTypeCode_SUPPLY_REQUEST_TYPE_CODE_SR
	MaintenanceRequestTypeCodeUnspecified      = MaintenanceRequestTypeCode_MAINTENANCE_REQUEST_TYPE_CODE_UNSPECIFIED
	MaintenanceRequestTypeCodeCM               = MaintenanceRequestTypeCode_MAINTENANCE_REQUEST_TYPE_CODE_CM
	HealthCollectionRequestTypeCodeUnspecified = HealthCollectionRequestTypeCode_HEALTH_COLLECTION_REQUEST_TYPE_CODE_UNSPECIFIED
	HealthCollectionRequestTypeCodeCR          = HealthCollectionRequestTypeCode_HEALTH_COLLECTION_REQUEST_TYPE_CODE_CR
)

const (
	CargoHMICCodeUnspecified     = CargoHMICCode_CARGO_HMIC_CODE_UNSPECIFIED
	CargoHMICCodeD               = CargoHMICCode_CARGO_HMIC_CODE_D
	CargoHMICCodeN               = CargoHMICCode_CARGO_HMIC_CODE_N
	CargoHMICCodeP               = CargoHMICCode_CARGO_HMIC_CODE_P
	CargoHMICCodeY               = CargoHMICCode_CARGO_HMIC_CODE_Y
	CargoHandlingCodeUnspecified = CargoHandlingCode_CARGO_HANDLING_CODE_UNSPECIFIED
	CargoHandlingCodeC           = CargoHandlingCode_CARGO_HANDLING_CODE_C
	CargoHandlingCodeM           = CargoHandlingCode_CARGO_HANDLING_CODE_M
	CargoHandlingCodeT           = CargoHandlingCode_CARGO_HANDLING_CODE_T
	CargoHandlingCodeR           = CargoHandlingCode_CARGO_HANDLING_CODE_R
	CargoHandlingCodeX           = CargoHandlingCode_CARGO_HANDLING_CODE_X
)

const (
	AttachmentIndicatorCodeUnspecified = AttachmentIndicatorCode_ATTACHMENT_INDICATOR_CODE_UNSPECIFIED
	AttachmentIndicatorCodeNo          = AttachmentIndicatorCode_ATTACHMENT_INDICATOR_CODE_NO
	AttachmentIndicatorCodeYes         = AttachmentIndicatorCode_ATTACHMENT_INDICATOR_CODE_YES
)

const (
	MaintenanceOperationalConditionCodeUnspecified = MaintenanceOperationalConditionCode_MAINTENANCE_OPERATIONAL_CONDITION_CODE_UNSPECIFIED
	MaintenanceOperationalConditionCodeA           = MaintenanceOperationalConditionCode_MAINTENANCE_OPERATIONAL_CONDITION_CODE_A
	MaintenanceOperationalConditionCodeB           = MaintenanceOperationalConditionCode_MAINTENANCE_OPERATIONAL_CONDITION_CODE_B
	MaintenanceOperationalConditionCodeC           = MaintenanceOperationalConditionCode_MAINTENANCE_OPERATIONAL_CONDITION_CODE_C
	MaintenanceSupportTypeCodeUnspecified          = MaintenanceSupportTypeCode_MAINTENANCE_SUPPORT_TYPE_CODE_UNSPECIFIED
	MaintenanceSupportTypeCodeXX                   = MaintenanceSupportTypeCode_MAINTENANCE_SUPPORT_TYPE_CODE_XX
	MaintenanceSupportTypeCodeR1                   = MaintenanceSupportTypeCode_MAINTENANCE_SUPPORT_TYPE_CODE_R1
	MaintenanceSupportTypeCodeR2                   = MaintenanceSupportTypeCode_MAINTENANCE_SUPPORT_TYPE_CODE_R2
	MaintenanceSupportTypeCodeR3                   = MaintenanceSupportTypeCode_MAINTENANCE_SUPPORT_TYPE_CODE_R3
	MaintenanceSupportTypeCodeR4                   = MaintenanceSupportTypeCode_MAINTENANCE_SUPPORT_TYPE_CODE_R4
	MaintenanceRepairTypeCodeUnspecified           = MaintenanceRepairTypeCode_MAINTENANCE_REPAIR_TYPE_CODE_UNSPECIFIED
	MaintenanceRepairTypeCodeM1                    = MaintenanceRepairTypeCode_MAINTENANCE_REPAIR_TYPE_CODE_M1
	MaintenanceRepairTypeCodeS1                    = MaintenanceRepairTypeCode_MAINTENANCE_REPAIR_TYPE_CODE_S1
	MaintenanceRepairTypeCodeS2                    = MaintenanceRepairTypeCode_MAINTENANCE_REPAIR_TYPE_CODE_S2
	MaintenanceRepairTypeCodeC1                    = MaintenanceRepairTypeCode_MAINTENANCE_REPAIR_TYPE_CODE_C1
	MaintenanceRepairTypeCodeD1                    = MaintenanceRepairTypeCode_MAINTENANCE_REPAIR_TYPE_CODE_D1
	MaintenanceMajorDefectCodeUnspecified          = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_UNSPECIFIED
	MaintenanceMajorDefectCodeMD01                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD01
	MaintenanceMajorDefectCodeMD02                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD02
	MaintenanceMajorDefectCodeMD03                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD03
	MaintenanceMajorDefectCodeMD04                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD04
	MaintenanceMajorDefectCodeMD05                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD05
	MaintenanceMajorDefectCodeMD06                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD06
	MaintenanceMajorDefectCodeMD07                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD07
	MaintenanceMajorDefectCodeMD08                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD08
	MaintenanceMajorDefectCodeMD09                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD09
	MaintenanceMajorDefectCodeMD10                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD10
	MaintenanceMajorDefectCodeMD11                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD11
	MaintenanceMajorDefectCodeMD12                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD12
	MaintenanceMajorDefectCodeMD13                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD13
	MaintenanceMajorDefectCodeMD14                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD14
	MaintenanceMajorDefectCodeMD15                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD15
	MaintenanceMajorDefectCodeMD16                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_MD16
	MaintenanceMajorDefectCodeNMAJ                 = MaintenanceMajorDefectCode_MAINTENANCE_MAJOR_DEFECT_CODE_NMAJ
)

const (
	RoadClassificationCodeUnspecified = RoadClassificationCode_ROAD_CLASSIFICATION_CODE_UNSPECIFIED
	RoadClassificationCodeA           = RoadClassificationCode_ROAD_CLASSIFICATION_CODE_A
	RoadClassificationCodeB           = RoadClassificationCode_ROAD_CLASSIFICATION_CODE_B
	RoadClassificationCodeC           = RoadClassificationCode_ROAD_CLASSIFICATION_CODE_C
	RoadClassificationCodeD           = RoadClassificationCode_ROAD_CLASSIFICATION_CODE_D
	RoadDrainageCodeUnspecified       = RoadDrainageCode_ROAD_DRAINAGE_CODE_UNSPECIFIED
	RoadDrainageCodeA                 = RoadDrainageCode_ROAD_DRAINAGE_CODE_A
	RoadDrainageCodeB                 = RoadDrainageCode_ROAD_DRAINAGE_CODE_B
	RoadFoundationCodeUnspecified     = RoadFoundationCode_ROAD_FOUNDATION_CODE_UNSPECIFIED
	RoadFoundationCodeA               = RoadFoundationCode_ROAD_FOUNDATION_CODE_A
	RoadFoundationCodeB               = RoadFoundationCode_ROAD_FOUNDATION_CODE_B
	RoadSurfaceTypeCodeUnspecified    = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_UNSPECIFIED
	RoadSurfaceTypeCodeA              = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_A
	RoadSurfaceTypeCodeB              = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_B
	RoadSurfaceTypeCodeC              = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_C
	RoadSurfaceTypeCodeD              = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_D
	RoadSurfaceTypeCodeE              = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_E
	RoadSurfaceTypeCodeF              = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_F
	RoadSurfaceTypeCodeG              = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_G
	RoadSurfaceTypeCodeH              = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_H
	RoadSurfaceTypeCodeI              = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_I
	RoadSurfaceTypeCodeJ              = RoadSurfaceTypeCode_ROAD_SURFACE_TYPE_CODE_J
)

const (
	EstimateCodeUnspecified                     = EstimateCode_ESTIMATE_CODE_UNSPECIFIED
	EstimateCodeNo                              = EstimateCode_ESTIMATE_CODE_NO
	EstimateCodeYes                             = EstimateCode_ESTIMATE_CODE_YES
	LandingZoneLayoutDesignationCodeUnspecified = LandingZoneLayoutDesignationCode_LANDING_ZONE_LAYOUT_DESIGNATION_CODE_UNSPECIFIED
	LandingZoneLayoutDesignationCodeLZ          = LandingZoneLayoutDesignationCode_LANDING_ZONE_LAYOUT_DESIGNATION_CODE_LZ
	LandingZoneLayoutDesignationCodeLS          = LandingZoneLayoutDesignationCode_LANDING_ZONE_LAYOUT_DESIGNATION_CODE_LS
	LandingZoneLayoutDesignationCodeLP          = LandingZoneLayoutDesignationCode_LANDING_ZONE_LAYOUT_DESIGNATION_CODE_LP
	AircraftSupportabilityCodeUnspecified       = AircraftSupportabilityCode_AIRCRAFT_SUPPORTABILITY_CODE_UNSPECIFIED
	AircraftSupportabilityCodeA                 = AircraftSupportabilityCode_AIRCRAFT_SUPPORTABILITY_CODE_A
	AircraftSupportabilityCodeB                 = AircraftSupportabilityCode_AIRCRAFT_SUPPORTABILITY_CODE_B
	AircraftSupportabilityCodeC                 = AircraftSupportabilityCode_AIRCRAFT_SUPPORTABILITY_CODE_C
	AircraftSupportabilityCodeD                 = AircraftSupportabilityCode_AIRCRAFT_SUPPORTABILITY_CODE_D
	AircraftSupportabilityCodeE                 = AircraftSupportabilityCode_AIRCRAFT_SUPPORTABILITY_CODE_E
	AircraftSupportabilityCodeF                 = AircraftSupportabilityCode_AIRCRAFT_SUPPORTABILITY_CODE_F
	CardinalDirectionCodeUnspecified            = CardinalDirectionCode_CARDINAL_DIRECTION_CODE_UNSPECIFIED
	CardinalDirectionCodeS                      = CardinalDirectionCode_CARDINAL_DIRECTION_CODE_S
	CardinalDirectionCodeE                      = CardinalDirectionCode_CARDINAL_DIRECTION_CODE_E
	CardinalDirectionCodeN                      = CardinalDirectionCode_CARDINAL_DIRECTION_CODE_N
	CardinalDirectionCodeW                      = CardinalDirectionCode_CARDINAL_DIRECTION_CODE_W
	LandingZoneSurfaceSlopeCodeUnspecified      = LandingZoneSurfaceSlopeCode_LANDING_ZONE_SURFACE_SLOPE_CODE_UNSPECIFIED
	LandingZoneSurfaceSlopeCodeA                = LandingZoneSurfaceSlopeCode_LANDING_ZONE_SURFACE_SLOPE_CODE_A
	LandingZoneSurfaceSlopeCodeB                = LandingZoneSurfaceSlopeCode_LANDING_ZONE_SURFACE_SLOPE_CODE_B
	LandingZoneObstacleCodeUnspecified          = LandingZoneObstacleCode_LANDING_ZONE_OBSTACLE_CODE_UNSPECIFIED
	LandingZoneObstacleCode1                    = LandingZoneObstacleCode_LANDING_ZONE_OBSTACLE_CODE_1
	LandingZoneObstacleCode2                    = LandingZoneObstacleCode_LANDING_ZONE_OBSTACLE_CODE_2
	LandingZoneObstacleCode3                    = LandingZoneObstacleCode_LANDING_ZONE_OBSTACLE_CODE_3
)

const (
	ObstacleActionCodeUnspecified        = ObstacleActionCode_OBSTACLE_ACTION_CODE_UNSPECIFIED
	ObstacleActionCode1                  = ObstacleActionCode_OBSTACLE_ACTION_CODE_1
	ObstacleActionCode2                  = ObstacleActionCode_OBSTACLE_ACTION_CODE_2
	ObstacleActionCode3                  = ObstacleActionCode_OBSTACLE_ACTION_CODE_3
	ObstacleDeterminationCodeUnspecified = ObstacleDeterminationCode_OBSTACLE_DETERMINATION_CODE_UNSPECIFIED
	ObstacleDeterminationCode1           = ObstacleDeterminationCode_OBSTACLE_DETERMINATION_CODE_1
	ObstacleDeterminationCode2           = ObstacleDeterminationCode_OBSTACLE_DETERMINATION_CODE_2
	ObstacleDeterminationCode3           = ObstacleDeterminationCode_OBSTACLE_DETERMINATION_CODE_3
	ObstacleDeterminationCode4           = ObstacleDeterminationCode_OBSTACLE_DETERMINATION_CODE_4
	BypassCodeUnspecified                = BypassCode_BYPASS_CODE_UNSPECIFIED
	BypassCodeNo                         = BypassCode_BYPASS_CODE_NO
	BypassCodeYes                        = BypassCode_BYPASS_CODE_YES
)

const (
	CBRNAgentTypeCodeUnspecified          = CBRNAgentTypeCode_CBRN_AGENT_TYPE_CODE_UNSPECIFIED
	CBRNAgentTypeCode1                    = CBRNAgentTypeCode_CBRN_AGENT_TYPE_CODE_1
	CBRNAgentTypeCode2                    = CBRNAgentTypeCode_CBRN_AGENT_TYPE_CODE_2
	CBRNAgentTypeCode3                    = CBRNAgentTypeCode_CBRN_AGENT_TYPE_CODE_3
	CBRNAgentTypeCode4                    = CBRNAgentTypeCode_CBRN_AGENT_TYPE_CODE_4
	CBRNAgentTypeCode5                    = CBRNAgentTypeCode_CBRN_AGENT_TYPE_CODE_5
	CBRNPhysicalPropertyCodeUnspecified   = CBRNPhysicalPropertyCode_CBRN_PHYSICAL_PROPERTY_CODE_UNSPECIFIED
	CBRNPhysicalPropertyCode1             = CBRNPhysicalPropertyCode_CBRN_PHYSICAL_PROPERTY_CODE_1
	CBRNPhysicalPropertyCode2             = CBRNPhysicalPropertyCode_CBRN_PHYSICAL_PROPERTY_CODE_2
	CBRNPhysicalPropertyCode3             = CBRNPhysicalPropertyCode_CBRN_PHYSICAL_PROPERTY_CODE_3
	CBRNPhysicalPropertyCode4             = CBRNPhysicalPropertyCode_CBRN_PHYSICAL_PROPERTY_CODE_4
	CBRNContaminationValueCodeUnspecified = CBRNContaminationValueCode_CBRN_CONTAMINATION_VALUE_CODE_UNSPECIFIED
	CBRNContaminationValueCodeE           = CBRNContaminationValueCode_CBRN_CONTAMINATION_VALUE_CODE_E
	CBRNContaminationValueCodeI           = CBRNContaminationValueCode_CBRN_CONTAMINATION_VALUE_CODE_I
	CBRNContaminationValueCodeW           = CBRNContaminationValueCode_CBRN_CONTAMINATION_VALUE_CODE_W
	CBRNContaminationValueCodeC           = CBRNContaminationValueCode_CBRN_CONTAMINATION_VALUE_CODE_C
	MunitionPurposeCodeUnspecified        = MunitionPurposeCode_MUNITION_PURPOSE_CODE_UNSPECIFIED
	MunitionPurposeCodeAA                 = MunitionPurposeCode_MUNITION_PURPOSE_CODE_AA
	MunitionPurposeCodeAP                 = MunitionPurposeCode_MUNITION_PURPOSE_CODE_AP
	MunitionPurposeCodeFL                 = MunitionPurposeCode_MUNITION_PURPOSE_CODE_FL
	MunitionPurposeCodeSM                 = MunitionPurposeCode_MUNITION_PURPOSE_CODE_SM
	MunitionPurposeCodeIM                 = MunitionPurposeCode_MUNITION_PURPOSE_CODE_IM
	MunitionTypeCodeUnspecified           = MunitionTypeCode_MUNITION_TYPE_CODE_UNSPECIFIED
	MunitionTypeCodeE                     = MunitionTypeCode_MUNITION_TYPE_CODE_E
	MunitionTypeCodeD                     = MunitionTypeCode_MUNITION_TYPE_CODE_D
	MunitionTypeCodeT                     = MunitionTypeCode_MUNITION_TYPE_CODE_T
	MunitionTypeCodeP                     = MunitionTypeCode_MUNITION_TYPE_CODE_P
)

const (
	ServiceCodeUnspecified = ServiceCode_SERVICE_CODE_UNSPECIFIED
	ServiceCodeUSA         = ServiceCode_SERVICE_CODE_USA
	ServiceCodeUSSF        = ServiceCode_SERVICE_CODE_USSF
	ServiceCodeUSAF        = ServiceCode_SERVICE_CODE_USAF
	ServiceCodeUSCG        = ServiceCode_SERVICE_CODE_USCG
	ServiceCodeUSN         = ServiceCode_SERVICE_CODE_USN
	ServiceCodeUSMC        = ServiceCode_SERVICE_CODE_USMC
	ServiceCodeUSCIV       = ServiceCode_SERVICE_CODE_US_CIV
	ServiceCodeNONUS       = ServiceCode_SERVICE_CODE_NON_US
	ServiceCodeEPW         = ServiceCode_SERVICE_CODE_EPW
)

const (
	HealthPrimaryMechanismCodeUnspecified   = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_UNSPECIFIED
	HealthPrimaryMechanismCodeE1            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_E1
	HealthPrimaryMechanismCodeE2            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_E2
	HealthPrimaryMechanismCodeE3            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_E3
	HealthPrimaryMechanismCodeE4            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_E4
	HealthPrimaryMechanismCodeE5            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_E5
	HealthPrimaryMechanismCodeP1            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_P1
	HealthPrimaryMechanismCodeP2            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_P2
	HealthPrimaryMechanismCodeP3            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_P3
	HealthPrimaryMechanismCodeP4            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_P4
	HealthPrimaryMechanismCodeP5            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_P5
	HealthPrimaryMechanismCodeD1            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_D1
	HealthPrimaryMechanismCodeD2            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_D2
	HealthPrimaryMechanismCodeD3            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_D3
	HealthPrimaryMechanismCodeD4            = HealthPrimaryMechanismCode_HEALTH_PRIMARY_MECHANISM_CODE_D4
	HealthCBRNExposureCodeUnspecified       = HealthCBRNExposureCode_HEALTH_CBRN_EXPOSURE_CODE_UNSPECIFIED
	HealthCBRNExposureCodeC                 = HealthCBRNExposureCode_HEALTH_CBRN_EXPOSURE_CODE_C
	HealthCBRNExposureCodeB                 = HealthCBRNExposureCode_HEALTH_CBRN_EXPOSURE_CODE_B
	HealthCBRNExposureCodeR                 = HealthCBRNExposureCode_HEALTH_CBRN_EXPOSURE_CODE_R
	HealthCBRNExposureCodeN                 = HealthCBRNExposureCode_HEALTH_CBRN_EXPOSURE_CODE_N
	HealthCBRNExposureCodeX                 = HealthCBRNExposureCode_HEALTH_CBRN_EXPOSURE_CODE_X
	HealthMajorSignsSymptomsCodeUnspecified = HealthMajorSignsSymptomsCode_HEALTH_MAJOR_SIGNS_SYMPTOMS_CODE_UNSPECIFIED
	HealthMajorSignsSymptomsCodeB           = HealthMajorSignsSymptomsCode_HEALTH_MAJOR_SIGNS_SYMPTOMS_CODE_B
	HealthMajorSignsSymptomsCodeR           = HealthMajorSignsSymptomsCode_HEALTH_MAJOR_SIGNS_SYMPTOMS_CODE_R
	HealthMajorSignsSymptomsCodeX           = HealthMajorSignsSymptomsCode_HEALTH_MAJOR_SIGNS_SYMPTOMS_CODE_X
	HealthMajorSignsSymptomsCodeC           = HealthMajorSignsSymptomsCode_HEALTH_MAJOR_SIGNS_SYMPTOMS_CODE_C
	HealthPulseLocationCodeUnspecified      = HealthPulseLocationCode_HEALTH_PULSE_LOCATION_CODE_UNSPECIFIED
	HealthPulseLocationCodeW                = HealthPulseLocationCode_HEALTH_PULSE_LOCATION_CODE_W
	HealthPulseLocationCodeN                = HealthPulseLocationCode_HEALTH_PULSE_LOCATION_CODE_N
	HealthResponsivenessCodeUnspecified     = HealthResponsivenessCode_HEALTH_RESPONSIVENESS_CODE_UNSPECIFIED
	HealthResponsivenessCodeA               = HealthResponsivenessCode_HEALTH_RESPONSIVENESS_CODE_A
	HealthResponsivenessCodeV               = HealthResponsivenessCode_HEALTH_RESPONSIVENESS_CODE_V
	HealthResponsivenessCodeP               = HealthResponsivenessCode_HEALTH_RESPONSIVENESS_CODE_P
	HealthResponsivenessCodeU               = HealthResponsivenessCode_HEALTH_RESPONSIVENESS_CODE_U
	HealthTriagePrecedenceCodeUnspecified   = HealthTriagePrecedenceCode_HEALTH_TRIAGE_PRECEDENCE_CODE_UNSPECIFIED
	HealthTriagePrecedenceCodeA             = HealthTriagePrecedenceCode_HEALTH_TRIAGE_PRECEDENCE_CODE_A
	HealthTriagePrecedenceCodeB             = HealthTriagePrecedenceCode_HEALTH_TRIAGE_PRECEDENCE_CODE_B
	HealthTriagePrecedenceCodeC             = HealthTriagePrecedenceCode_HEALTH_TRIAGE_PRECEDENCE_CODE_C
	HealthTriagePrecedenceCodeD             = HealthTriagePrecedenceCode_HEALTH_TRIAGE_PRECEDENCE_CODE_D
	HealthTriagePrecedenceCodeE             = HealthTriagePrecedenceCode_HEALTH_TRIAGE_PRECEDENCE_CODE_E
)

const (
	TourniquetPlacementCodeUnspecified  = TourniquetPlacementCode_TOURNIQUET_PLACEMENT_CODE_UNSPECIFIED
	TourniquetPlacementCodeTQXX         = TourniquetPlacementCode_TOURNIQUET_PLACEMENT_CODE_TQXX
	TourniquetPlacementCodeTQRA         = TourniquetPlacementCode_TOURNIQUET_PLACEMENT_CODE_TQRA
	TourniquetPlacementCodeTQLA         = TourniquetPlacementCode_TOURNIQUET_PLACEMENT_CODE_TQLA
	TourniquetPlacementCodeTQRL         = TourniquetPlacementCode_TOURNIQUET_PLACEMENT_CODE_TQRL
	TourniquetPlacementCodeTQLL         = TourniquetPlacementCode_TOURNIQUET_PLACEMENT_CODE_TQLL
	TourniquetTypeCodeUnspecified       = TourniquetTypeCode_TOURNIQUET_TYPE_CODE_UNSPECIFIED
	TourniquetTypeCodeE                 = TourniquetTypeCode_TOURNIQUET_TYPE_CODE_E
	TourniquetTypeCodeJ                 = TourniquetTypeCode_TOURNIQUET_TYPE_CODE_J
	TourniquetTypeCodeT                 = TourniquetTypeCode_TOURNIQUET_TYPE_CODE_T
	WoundTreatmentCodeUnspecified       = WoundTreatmentCode_WOUND_TREATMENT_CODE_UNSPECIFIED
	WoundTreatmentCodeT1                = WoundTreatmentCode_WOUND_TREATMENT_CODE_T1
	WoundTreatmentCodeT2                = WoundTreatmentCode_WOUND_TREATMENT_CODE_T2
	WoundTreatmentCodeT3                = WoundTreatmentCode_WOUND_TREATMENT_CODE_T3
	WoundTreatmentCodeT4                = WoundTreatmentCode_WOUND_TREATMENT_CODE_T4
	WoundTreatmentCodeT5                = WoundTreatmentCode_WOUND_TREATMENT_CODE_T5
	WoundTreatmentCodeT6                = WoundTreatmentCode_WOUND_TREATMENT_CODE_T6
	WoundTreatmentCodeT7                = WoundTreatmentCode_WOUND_TREATMENT_CODE_T7
	AirwayTreatmentCodeUnspecified      = AirwayTreatmentCode_AIRWAY_TREATMENT_CODE_UNSPECIFIED
	AirwayTreatmentCodeA0               = AirwayTreatmentCode_AIRWAY_TREATMENT_CODE_A0
	AirwayTreatmentCodeA1               = AirwayTreatmentCode_AIRWAY_TREATMENT_CODE_A1
	AirwayTreatmentCodeA2               = AirwayTreatmentCode_AIRWAY_TREATMENT_CODE_A2
	AirwayTreatmentCodeA3               = AirwayTreatmentCode_AIRWAY_TREATMENT_CODE_A3
	AirwayTreatmentCodeA4               = AirwayTreatmentCode_AIRWAY_TREATMENT_CODE_A4
	BreathingTreatmentCodeUnspecified   = BreathingTreatmentCode_BREATHING_TREATMENT_CODE_UNSPECIFIED
	BreathingTreatmentCodeB0            = BreathingTreatmentCode_BREATHING_TREATMENT_CODE_B0
	BreathingTreatmentCodeB1            = BreathingTreatmentCode_BREATHING_TREATMENT_CODE_B1
	BreathingTreatmentCodeB3            = BreathingTreatmentCode_BREATHING_TREATMENT_CODE_B3
	BreathingTreatmentCodeB4            = BreathingTreatmentCode_BREATHING_TREATMENT_CODE_B4
	BreathingTreatmentCodeB5            = BreathingTreatmentCode_BREATHING_TREATMENT_CODE_B5
	FluidNameCodeUnspecified            = FluidNameCode_FLUID_NAME_CODE_UNSPECIFIED
	FluidNameCodeS                      = FluidNameCode_FLUID_NAME_CODE_S
	FluidNameCodeR                      = FluidNameCode_FLUID_NAME_CODE_R
	FluidNameCodeH                      = FluidNameCode_FLUID_NAME_CODE_H
	FluidRouteCodeUnspecified           = FluidRouteCode_FLUID_ROUTE_CODE_UNSPECIFIED
	FluidRouteCodeIV                    = FluidRouteCode_FLUID_ROUTE_CODE_IV
	FluidRouteCodeIO                    = FluidRouteCode_FLUID_ROUTE_CODE_IO
	BloodProductCodeUnspecified         = BloodProductCode_BLOOD_PRODUCT_CODE_UNSPECIFIED
	BloodProductCodeWBD                 = BloodProductCode_BLOOD_PRODUCT_CODE_WBD
	BloodProductCodeRBC                 = BloodProductCode_BLOOD_PRODUCT_CODE_RBC
	BloodProductCodeFFP                 = BloodProductCode_BLOOD_PRODUCT_CODE_FFP
	BloodProductCodeFDP                 = BloodProductCode_BLOOD_PRODUCT_CODE_FDP
	MedicationRouteCodeUnspecified      = MedicationRouteCode_MEDICATION_ROUTE_CODE_UNSPECIFIED
	MedicationRouteCodeR1               = MedicationRouteCode_MEDICATION_ROUTE_CODE_R1
	MedicationRouteCodeR2               = MedicationRouteCode_MEDICATION_ROUTE_CODE_R2
	MedicationRouteCodeR3               = MedicationRouteCode_MEDICATION_ROUTE_CODE_R3
	MedicationRouteCodeR4               = MedicationRouteCode_MEDICATION_ROUTE_CODE_R4
	MedicationRouteCodeR5               = MedicationRouteCode_MEDICATION_ROUTE_CODE_R5
	MedicationRouteCodeR6               = MedicationRouteCode_MEDICATION_ROUTE_CODE_R6
	MedicationRouteCodeR7               = MedicationRouteCode_MEDICATION_ROUTE_CODE_R7
	AnalgesicMedicationCodeUnspecified  = AnalgesicMedicationCode_ANALGESIC_MEDICATION_CODE_UNSPECIFIED
	AnalgesicMedicationCodeK            = AnalgesicMedicationCode_ANALGESIC_MEDICATION_CODE_K
	AnalgesicMedicationCodeF            = AnalgesicMedicationCode_ANALGESIC_MEDICATION_CODE_F
	AnalgesicMedicationCodeM            = AnalgesicMedicationCode_ANALGESIC_MEDICATION_CODE_M
	AntibioticMedicationCodeUnspecified = AntibioticMedicationCode_ANTIBIOTIC_MEDICATION_CODE_UNSPECIFIED
	AntibioticMedicationCodeM           = AntibioticMedicationCode_ANTIBIOTIC_MEDICATION_CODE_M
	AntibioticMedicationCodeE           = AntibioticMedicationCode_ANTIBIOTIC_MEDICATION_CODE_E
	AntibioticMedicationCodeP           = AntibioticMedicationCode_ANTIBIOTIC_MEDICATION_CODE_P
	AntibioticMedicationCodeA           = AntibioticMedicationCode_ANTIBIOTIC_MEDICATION_CODE_A
	OtherMedicationCodeUnspecified      = OtherMedicationCode_OTHER_MEDICATION_CODE_UNSPECIFIED
	OtherMedicationCodeI                = OtherMedicationCode_OTHER_MEDICATION_CODE_I
	OtherMedicationCodeT                = OtherMedicationCode_OTHER_MEDICATION_CODE_T
	CasualtyTypeCodeUnspecified         = CasualtyTypeCode_CASUALTY_TYPE_CODE_UNSPECIFIED
	CasualtyTypeCodeA                   = CasualtyTypeCode_CASUALTY_TYPE_CODE_A
	CasualtyTypeCodeB                   = CasualtyTypeCode_CASUALTY_TYPE_CODE_B
)
