package lxdr

import "testing"

func testPAXRequestContainer() *RequestContainer {
	return &RequestContainer{
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
}

func TestDetermineRequestLifecycleStateLocalOnly(t *testing.T) {
	state, err := DetermineRequestLifecycleState(testPAXRequestContainer(), false)
	if err != nil {
		t.Fatalf("determine local-only lifecycle state: %v", err)
	}

	if state != RequestLifecycleStateLocalOnly {
		t.Fatalf("state = %q, want %q", state, RequestLifecycleStateLocalOnly)
	}
}

func TestDetermineRequestLifecycleStateFromLinkFrameCarried(t *testing.T) {
	container := testPAXRequestContainer()
	frame, err := NewRequestContainerLinkFrame(container, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("new request-container link frame: %v", err)
	}

	state, err := DetermineRequestLifecycleStateFromLinkFrame(container, frame)
	if err != nil {
		t.Fatalf("determine carried lifecycle state from link frame: %v", err)
	}

	if state != RequestLifecycleStateCarried {
		t.Fatalf("state = %q, want %q", state, RequestLifecycleStateCarried)
	}
}

func TestApplySynchronizedResponseAndDetermineState(t *testing.T) {
	container := testPAXRequestContainer()
	resp := &SynchronizedResponse{
		LocalRequestId:        "3838JBNM5",
		SynchronizedRequestId: "KL9K15474QFJ",
	}

	state, err := ApplySynchronizedResponseAndDetermineState(container, resp)
	if err != nil {
		t.Fatalf("apply synchronized response and determine state: %v", err)
	}

	if state != RequestLifecycleStateSynchronized {
		t.Fatalf("state = %q, want %q", state, RequestLifecycleStateSynchronized)
	}
}

func TestDetermineRequestLifecycleStateFromLinkFrameRejectsNonRequestPayload(t *testing.T) {
	frame, err := NewSynchronizedResponseLinkFrame(
		&SynchronizedResponse{
			LocalRequestId:        "3838JBNM5",
			SynchronizedRequestId: "KL9K15474QFJ",
		},
		LinkDeliveryMethodDirect,
	)
	if err != nil {
		t.Fatalf("new synchronized-response link frame: %v", err)
	}

	if _, err := DetermineRequestLifecycleStateFromLinkFrame(testPAXRequestContainer(), frame); err == nil {
		t.Fatalf("expected non-request payload error")
	}
}
