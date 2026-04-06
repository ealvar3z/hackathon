package lxdr

import (
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
)

func TestMarshalUnmarshalRequestContainerBinary(t *testing.T) {
	container := &RequestContainer{
		Header: testHeader(),
		Segments: []*RequestSegment{
			wrapMobilityPax(&MobilityPaxRequestSegment{
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
			}),
		},
	}

	data, err := MarshalRequestContainerBinary(container)
	if err != nil {
		t.Fatalf("marshal request container: %v", err)
	}

	got, err := UnmarshalRequestContainerBinary(data)
	if err != nil {
		t.Fatalf("unmarshal request container: %v", err)
	}

	if got.Header.LocalRequestId != container.Header.LocalRequestId {
		t.Fatalf(
			"local request id = %q, want %q",
			got.Header.LocalRequestId,
			container.Header.LocalRequestId,
		)
	}
	if len(got.Segments) != 1 {
		t.Fatalf("segment count = %d, want 1", len(got.Segments))
	}

	pax := got.Segments[0].GetMobilityPax()
	if pax == nil {
		t.Fatalf("expected mobility pax segment")
	}
	if pax.ZapOrEdiPi != "1010919789" {
		t.Fatalf("zap/edi-pi = %q, want %q", pax.ZapOrEdiPi, "1010919789")
	}
}

func TestMarshalRequestContainerBinaryRejectsInvalid(t *testing.T) {
	container := &RequestContainer{}

	if _, err := MarshalRequestContainerBinary(container); err == nil {
		t.Fatalf("expected invalid request container error")
	}
}

func TestMarshalUnmarshalSynchronizedResponseBinary(t *testing.T) {
	resp := &SynchronizedResponse{
		LocalRequestId:        "3838JBNM5",
		SynchronizedRequestId: "KL9K15474QFJ",
	}

	data, err := MarshalSynchronizedResponseBinary(resp)
	if err != nil {
		t.Fatalf("marshal synchronized response: %v", err)
	}

	got, err := UnmarshalSynchronizedResponseBinary(data)
	if err != nil {
		t.Fatalf("unmarshal synchronized response: %v", err)
	}

	if got.LocalRequestId != resp.LocalRequestId {
		t.Fatalf(
			"local request id = %q, want %q",
			got.LocalRequestId,
			resp.LocalRequestId,
		)
	}
	if got.SynchronizedRequestId != resp.SynchronizedRequestId {
		t.Fatalf(
			"synchronized request id = %q, want %q",
			got.SynchronizedRequestId,
			resp.SynchronizedRequestId,
		)
	}
}

func TestMarshalUnmarshalCanonicalRegistryBinary(t *testing.T) {
	var registry CanonicalRegistry
	if err := prototext.Unmarshal(
		readTextProtoFixture(t, "appendix_f_header_registry.textproto"),
		&registry,
	); err != nil {
		t.Fatalf("unmarshal textproto fixture: %v", err)
	}

	data, err := MarshalCanonicalRegistryBinary(&registry)
	if err != nil {
		t.Fatalf("marshal canonical registry: %v", err)
	}

	got, err := UnmarshalCanonicalRegistryBinary(data)
	if err != nil {
		t.Fatalf("unmarshal canonical registry: %v", err)
	}

	if len(got.Entries) != len(registry.Entries) {
		t.Fatalf("entry count = %d, want %d", len(got.Entries), len(registry.Entries))
	}
	if got.Entries[0].CanonicalField != "department_of_defense_activity_address_code" {
		t.Fatalf(
			"first canonical field = %q, want %q",
			got.Entries[0].CanonicalField,
			"department_of_defense_activity_address_code",
		)
	}
}

func TestMarshalUnmarshalLinkFrameBinaryWithRequestContainer(t *testing.T) {
	frame := &LinkFrame{
		LinkMessageId:  "LF-0001",
		DeliveryMethod: LinkDeliveryMethodDirect,
		Representation: LinkRepresentationBinaryProto,
		Payload: &LinkFrame_RequestContainer{
			RequestContainer: &RequestContainer{
				Header: testHeader(),
				Segments: []*RequestSegment{
					wrapMobilityPax(&MobilityPaxRequestSegment{
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
					}),
				},
			},
		},
	}

	data, err := MarshalLinkFrameBinary(frame)
	if err != nil {
		t.Fatalf("marshal link frame: %v", err)
	}

	got, err := UnmarshalLinkFrameBinary(data)
	if err != nil {
		t.Fatalf("unmarshal link frame: %v", err)
	}

	if got.GetRequestContainer() == nil {
		t.Fatalf("expected request container payload")
	}
	if got.LinkMessageId != "LF-0001" {
		t.Fatalf("link message id = %q, want %q", got.LinkMessageId, "LF-0001")
	}
	if got.GetRequestContainer().Header.LocalRequestId != "3838JBNM5" {
		t.Fatalf(
			"local request id = %q, want %q",
			got.GetRequestContainer().Header.LocalRequestId,
			"3838JBNM5",
		)
	}
}

func TestMarshalLinkFrameBinaryRejectsEmptyFrame(t *testing.T) {
	frame := &LinkFrame{}

	if _, err := MarshalLinkFrameBinary(frame); err == nil {
		t.Fatalf("expected invalid link frame error")
	}
}

func TestMarshalLinkFrameBinaryRejectsMissingLinkMetadata(t *testing.T) {
	frame := &LinkFrame{
		Payload: &LinkFrame_SynchronizedResponse{
			SynchronizedResponse: &SynchronizedResponse{
				LocalRequestId:        "3838JBNM5",
				SynchronizedRequestId: "KL9K15474QFJ",
			},
		},
	}

	if _, err := MarshalLinkFrameBinary(frame); err == nil {
		t.Fatalf("expected invalid link frame metadata error")
	}
}

func TestNewRequestContainerLinkFrame(t *testing.T) {
	container := &RequestContainer{
		Header: testHeader(),
		Segments: []*RequestSegment{
			wrapMobilityPax(&MobilityPaxRequestSegment{
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
			}),
		},
	}

	frame, err := NewRequestContainerLinkFrame(container, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("new request-container link frame: %v", err)
	}

	if frame.LinkMessageId == "" {
		t.Fatalf("expected link message id")
	}
	if frame.DeliveryMethod != LinkDeliveryMethodDirect {
		t.Fatalf("delivery method = %v, want %v", frame.DeliveryMethod, LinkDeliveryMethodDirect)
	}
	if frame.Representation != LinkRepresentationBinaryProto {
		t.Fatalf("representation = %v, want %v", frame.Representation, LinkRepresentationBinaryProto)
	}
	if frame.PayloadKind() != LinkPayloadKindRequestContainer {
		t.Fatalf("payload kind = %q, want %q", frame.PayloadKind(), LinkPayloadKindRequestContainer)
	}
}

func TestNewRequestContainerLinkFrameIsDeterministic(t *testing.T) {
	container := &RequestContainer{
		Header: testHeader(),
		Segments: []*RequestSegment{
			wrapMobilityPax(&MobilityPaxRequestSegment{
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
			}),
		},
	}

	frameA, err := NewRequestContainerLinkFrame(container, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("new request-container link frame A: %v", err)
	}
	frameB, err := NewRequestContainerLinkFrame(container, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("new request-container link frame B: %v", err)
	}

	if frameA.LinkMessageId != frameB.LinkMessageId {
		t.Fatalf("link message ids differ: %q != %q", frameA.LinkMessageId, frameB.LinkMessageId)
	}
}

func TestNewSynchronizedResponseLinkFrame(t *testing.T) {
	resp := &SynchronizedResponse{
		LocalRequestId:        "3838JBNM5",
		SynchronizedRequestId: "KL9K15474QFJ",
	}

	frame, err := NewSynchronizedResponseLinkFrame(resp, LinkDeliveryMethodPropagated)
	if err != nil {
		t.Fatalf("new synchronized-response link frame: %v", err)
	}

	if frame.PayloadKind() != LinkPayloadKindSynchronizedResponse {
		t.Fatalf("payload kind = %q, want %q", frame.PayloadKind(), LinkPayloadKindSynchronizedResponse)
	}
	if frame.DeliveryMethod != LinkDeliveryMethodPropagated {
		t.Fatalf("delivery method = %v, want %v", frame.DeliveryMethod, LinkDeliveryMethodPropagated)
	}
}

func TestNewCanonicalRegistryLinkFrame(t *testing.T) {
	var registry CanonicalRegistry
	if err := prototext.Unmarshal(
		readTextProtoFixture(t, "appendix_f_header_registry.textproto"),
		&registry,
	); err != nil {
		t.Fatalf("unmarshal textproto fixture: %v", err)
	}

	frame, err := NewCanonicalRegistryLinkFrame(&registry, LinkDeliveryMethodOpportunistic)
	if err != nil {
		t.Fatalf("new canonical-registry link frame: %v", err)
	}

	if frame.PayloadKind() != LinkPayloadKindCanonicalRegistry {
		t.Fatalf("payload kind = %q, want %q", frame.PayloadKind(), LinkPayloadKindCanonicalRegistry)
	}
	if frame.DeliveryMethod != LinkDeliveryMethodOpportunistic {
		t.Fatalf("delivery method = %v, want %v", frame.DeliveryMethod, LinkDeliveryMethodOpportunistic)
	}
}

func TestNewSynchronizedResponseForRequest(t *testing.T) {
	container := &RequestContainer{
		Header: testHeader(),
		Segments: []*RequestSegment{
			wrapMobilityPax(&MobilityPaxRequestSegment{
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
			}),
		},
	}

	resp, err := NewSynchronizedResponseForRequest(container, "KL9K15474QFJ")
	if err != nil {
		t.Fatalf("new synchronized response for request: %v", err)
	}

	if resp.LocalRequestId != "3838JBNM5" {
		t.Fatalf("local request id = %q, want %q", resp.LocalRequestId, "3838JBNM5")
	}
	if resp.SynchronizedRequestId != "KL9K15474QFJ" {
		t.Fatalf(
			"synchronized request id = %q, want %q",
			resp.SynchronizedRequestId,
			"KL9K15474QFJ",
		)
	}
}

func TestApplySynchronizedResponse(t *testing.T) {
	container := &RequestContainer{
		Header: testHeader(),
		Segments: []*RequestSegment{
			wrapMobilityPax(&MobilityPaxRequestSegment{
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
			}),
		},
	}
	resp := &SynchronizedResponse{
		LocalRequestId:        "3838JBNM5",
		SynchronizedRequestId: "KL9K15474QFJ",
	}

	if err := ApplySynchronizedResponse(container, resp); err != nil {
		t.Fatalf("apply synchronized response: %v", err)
	}
	if !container.Header.IsSynchronized() {
		t.Fatalf("expected synchronized header state")
	}
	if container.Header.GetSynchronizedRequestId() != "KL9K15474QFJ" {
		t.Fatalf(
			"synchronized request id = %q, want %q",
			container.Header.GetSynchronizedRequestId(),
			"KL9K15474QFJ",
		)
	}
}

func TestApplySynchronizedResponseRejectsMismatchedLocalRequestID(t *testing.T) {
	container := &RequestContainer{
		Header: testHeader(),
		Segments: []*RequestSegment{
			wrapMobilityPax(&MobilityPaxRequestSegment{
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
			}),
		},
	}
	resp := &SynchronizedResponse{
		LocalRequestId:        "ZZZZZZZZZ",
		SynchronizedRequestId: "KL9K15474QFJ",
	}

	if err := ApplySynchronizedResponse(container, resp); err == nil {
		t.Fatalf("expected mismatched synchronized response error")
	}
}

func TestNewSynchronizedResponseLinkFrameForRequest(t *testing.T) {
	container := &RequestContainer{
		Header: testHeader(),
		Segments: []*RequestSegment{
			wrapMobilityPax(&MobilityPaxRequestSegment{
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
			}),
		},
	}

	frame, err := NewSynchronizedResponseLinkFrameForRequest(
		container,
		"KL9K15474QFJ",
		LinkDeliveryMethodDirect,
	)
	if err != nil {
		t.Fatalf("new synchronized response link frame for request: %v", err)
	}

	if frame.PayloadKind() != LinkPayloadKindSynchronizedResponse {
		t.Fatalf("payload kind = %q, want %q", frame.PayloadKind(), LinkPayloadKindSynchronizedResponse)
	}
	if frame.GetSynchronizedResponse().LocalRequestId != "3838JBNM5" {
		t.Fatalf(
			"local request id = %q, want %q",
			frame.GetSynchronizedResponse().LocalRequestId,
			"3838JBNM5",
		)
	}
}
