package lxdr

import (
	"errors"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func TestRouterTrackRequestStartsGenerated(t *testing.T) {
	router := NewRouter()

	tracked, err := router.TrackRequest(testPAXRequestContainer())
	if err != nil {
		t.Fatalf("track request: %v", err)
	}

	if tracked.State != RouterRequestStateGenerated {
		t.Fatalf("state = %q, want %q", tracked.State, RouterRequestStateGenerated)
	}
}

func TestRouterQueueRequestTransitionsToOutbound(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}

	frame, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("queue request: %v", err)
	}

	if frame.PayloadKind() != LinkPayloadKindRequestContainer {
		t.Fatalf("payload kind = %q, want %q", frame.PayloadKind(), LinkPayloadKindRequestContainer)
	}

	tracked, ok := router.TrackedRequest(container.Header.LocalRequestId)
	if !ok {
		t.Fatalf("expected tracked request")
	}
	if tracked.State != RouterRequestStateOutbound {
		t.Fatalf("state = %q, want %q", tracked.State, RouterRequestStateOutbound)
	}
	if tracked.AttemptCount != 1 {
		t.Fatalf("attempt count = %d, want 1", tracked.AttemptCount)
	}
}

func TestRouterMarkRequestFrameCarriedTransitionsToCarried(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}
	if _, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect); err != nil {
		t.Fatalf("queue request: %v", err)
	}

	if err := router.MarkRequestFrameCarried(container.Header.LocalRequestId); err != nil {
		t.Fatalf("mark request frame carried: %v", err)
	}

	tracked, _ := router.TrackedRequest(container.Header.LocalRequestId)
	if tracked.State != RouterRequestStateCarried {
		t.Fatalf("state = %q, want %q", tracked.State, RouterRequestStateCarried)
	}
}

func TestRouterHandleInboundSyncFrameTransitionsToSynchronized(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}
	requestFrame, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("queue request: %v", err)
	}
	if err := router.MarkRequestFrameCarried(container.Header.LocalRequestId); err != nil {
		t.Fatalf("mark request frame carried: %v", err)
	}

	syncFrame, err := BuildSyncResponseForRequestFrame(
		requestFrame,
		"KL9K15474QFJ",
		LinkDeliveryMethodPropagated,
	)
	if err != nil {
		t.Fatalf("build sync response for request frame: %v", err)
	}

	disposition, err := router.HandleInboundFrame(syncFrame)
	if err != nil {
		t.Fatalf("handle inbound sync frame: %v", err)
	}
	if disposition != RouterFrameDispositionSynchronized {
		t.Fatalf("disposition = %q, want %q", disposition, RouterFrameDispositionSynchronized)
	}

	tracked, _ := router.TrackedRequest(container.Header.LocalRequestId)
	if tracked.State != RouterRequestStateSynchronized {
		t.Fatalf("state = %q, want %q", tracked.State, RouterRequestStateSynchronized)
	}
	if tracked.Container.Header.GetSynchronizedRequestId() != "KL9K15474QFJ" {
		t.Fatalf(
			"synchronized request id = %q, want %q",
			tracked.Container.Header.GetSynchronizedRequestId(),
			"KL9K15474QFJ",
		)
	}
}

func TestRouterSeenFrameStatus(t *testing.T) {
	router := NewRouter()
	frame, err := NewRequestContainerLinkFrame(testPAXRequestContainer(), LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("new request-container link frame: %v", err)
	}

	status, err := router.SeenFrameStatus(frame)
	if err != nil {
		t.Fatalf("seen frame status before record: %v", err)
	}
	if status != SeenFrameStatusNew {
		t.Fatalf("status = %q, want %q", status, SeenFrameStatusNew)
	}

	router.seenFrames[frame.LinkMessageId] = frame
	status, err = router.SeenFrameStatus(proto.Clone(frame).(*LinkFrame))
	if err != nil {
		t.Fatalf("seen frame status duplicate: %v", err)
	}
	if status != SeenFrameStatusDuplicate {
		t.Fatalf("status = %q, want %q", status, SeenFrameStatusDuplicate)
	}

	conflicting := proto.Clone(frame).(*LinkFrame)
	conflicting.DeliveryMethod = LinkDeliveryMethodPropagated
	status, err = router.SeenFrameStatus(conflicting)
	if err != nil {
		t.Fatalf("seen frame status conflicting: %v", err)
	}
	if status != SeenFrameStatusConflicting {
		t.Fatalf("status = %q, want %q", status, SeenFrameStatusConflicting)
	}
}

func TestRouterHandleInboundFrameReturnsDuplicateForRepeatedSyncFrame(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}
	requestFrame, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("queue request: %v", err)
	}

	syncFrame, err := BuildSyncResponseForRequestFrame(
		requestFrame,
		"KL9K15474QFJ",
		LinkDeliveryMethodPropagated,
	)
	if err != nil {
		t.Fatalf("build sync response for request frame: %v", err)
	}

	if _, err := router.HandleInboundFrame(syncFrame); err != nil {
		t.Fatalf("initial handle inbound sync frame: %v", err)
	}
	disposition, err := router.HandleInboundFrame(proto.Clone(syncFrame).(*LinkFrame))
	if err != nil {
		t.Fatalf("repeat handle inbound sync frame: %v", err)
	}
	if disposition != RouterFrameDispositionDuplicate {
		t.Fatalf("disposition = %q, want %q", disposition, RouterFrameDispositionDuplicate)
	}
}

func TestRouterRequestFrameForSyncResponse(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}
	requestFrame, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("queue request: %v", err)
	}
	syncFrame, err := BuildSyncResponseForRequestFrame(
		requestFrame,
		"KL9K15474QFJ",
		LinkDeliveryMethodPropagated,
	)
	if err != nil {
		t.Fatalf("build sync response for request frame: %v", err)
	}

	gotFrame, tracked, err := router.RequestFrameForSyncResponse(syncFrame)
	if err != nil {
		t.Fatalf("request frame for sync response: %v", err)
	}
	if gotFrame.LinkMessageId != requestFrame.LinkMessageId {
		t.Fatalf("request frame id = %q, want %q", gotFrame.LinkMessageId, requestFrame.LinkMessageId)
	}
	if tracked.Container.Header.LocalRequestId != container.Header.LocalRequestId {
		t.Fatalf(
			"tracked local request id = %q, want %q",
			tracked.Container.Header.LocalRequestId,
			container.Header.LocalRequestId,
		)
	}
}

func TestRouterSyncResponseStatus(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}
	requestFrame, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("queue request: %v", err)
	}
	syncFrame, err := BuildSyncResponseForRequestFrame(
		requestFrame,
		"KL9K15474QFJ",
		LinkDeliveryMethodPropagated,
	)
	if err != nil {
		t.Fatalf("build sync response for request frame: %v", err)
	}

	status, err := router.SyncResponseStatus(syncFrame)
	if err != nil {
		t.Fatalf("sync response status new: %v", err)
	}
	if status != SyncResponseStatusNew {
		t.Fatalf("status = %q, want %q", status, SyncResponseStatusNew)
	}

	if _, err := router.HandleInboundFrame(syncFrame); err != nil {
		t.Fatalf("handle inbound sync frame: %v", err)
	}

	status, err = router.SyncResponseStatus(proto.Clone(syncFrame).(*LinkFrame))
	if err != nil {
		t.Fatalf("sync response status repeated: %v", err)
	}
	if status != SyncResponseStatusRepeated {
		t.Fatalf("status = %q, want %q", status, SyncResponseStatusRepeated)
	}

	conflicting, err := BuildSyncResponseForRequestFrame(
		requestFrame,
		"OTHER0000000",
		LinkDeliveryMethodPropagated,
	)
	if err != nil {
		t.Fatalf("build conflicting sync response for request frame: %v", err)
	}
	status, err = router.SyncResponseStatus(conflicting)
	if err != nil {
		t.Fatalf("sync response status conflicting: %v", err)
	}
	if status != SyncResponseStatusConflicting {
		t.Fatalf("status = %q, want %q", status, SyncResponseStatusConflicting)
	}
}

func TestRouterHandleInboundFrameRejectsConflictingDuplicate(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}
	requestFrame, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("queue request: %v", err)
	}

	syncFrame, err := BuildSyncResponseForRequestFrame(
		requestFrame,
		"KL9K15474QFJ",
		LinkDeliveryMethodPropagated,
	)
	if err != nil {
		t.Fatalf("build sync response for request frame: %v", err)
	}
	if _, err := router.HandleInboundFrame(syncFrame); err != nil {
		t.Fatalf("initial handle inbound sync frame: %v", err)
	}

	conflicting := proto.Clone(syncFrame).(*LinkFrame)
	conflicting.ReferenceLinkMessageId = nil

	if _, err := router.HandleInboundFrame(conflicting); err == nil {
		t.Fatalf("expected conflicting duplicate error")
	}
}

func TestRouterMarkRequestFailedTransitionsToFailed(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}

	if err := router.MarkRequestFailed(container.Header.LocalRequestId, errors.New("no path")); err != nil {
		t.Fatalf("mark request failed: %v", err)
	}

	tracked, _ := router.TrackedRequest(container.Header.LocalRequestId)
	if tracked.State != RouterRequestStateFailed {
		t.Fatalf("state = %q, want %q", tracked.State, RouterRequestStateFailed)
	}
	if tracked.LastError != "no path" {
		t.Fatalf("last error = %q, want %q", tracked.LastError, "no path")
	}
	if tracked.NextAttemptAt != nil {
		t.Fatalf("expected no next attempt time for plain failure")
	}
}

func TestRouterScheduleRetrySetsRetryMetadata(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}

	nextAttemptAt := time.Date(2026, 4, 5, 12, 30, 0, 0, time.UTC)
	if err := router.ScheduleRetry(container.Header.LocalRequestId, nextAttemptAt, errors.New("link busy")); err != nil {
		t.Fatalf("schedule retry: %v", err)
	}

	tracked, _ := router.TrackedRequest(container.Header.LocalRequestId)
	if tracked.State != RouterRequestStateFailed {
		t.Fatalf("state = %q, want %q", tracked.State, RouterRequestStateFailed)
	}
	if tracked.LastError != "link busy" {
		t.Fatalf("last error = %q, want %q", tracked.LastError, "link busy")
	}
	if tracked.NextAttemptAt == nil || !tracked.NextAttemptAt.Equal(nextAttemptAt) {
		t.Fatalf("next attempt at = %v, want %v", tracked.NextAttemptAt, nextAttemptAt)
	}
}

func TestRouterReadyForRetry(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}

	nextAttemptAt := time.Date(2026, 4, 5, 12, 30, 0, 0, time.UTC)
	if err := router.ScheduleRetry(container.Header.LocalRequestId, nextAttemptAt, errors.New("link busy")); err != nil {
		t.Fatalf("schedule retry: %v", err)
	}

	ready, err := router.ReadyForRetry(container.Header.LocalRequestId, nextAttemptAt.Add(-time.Second))
	if err != nil {
		t.Fatalf("ready for retry before due: %v", err)
	}
	if ready {
		t.Fatalf("expected retry not ready before due time")
	}

	ready, err = router.ReadyForRetry(container.Header.LocalRequestId, nextAttemptAt)
	if err != nil {
		t.Fatalf("ready for retry at due time: %v", err)
	}
	if !ready {
		t.Fatalf("expected retry ready at due time")
	}
}

func TestRouterRetryableRequests(t *testing.T) {
	router := NewRouter()
	containerA := testPAXRequestContainer()
	containerB := testPAXRequestContainer()
	containerB.Header.LocalRequestId = "ZXCVBNM12"

	if _, err := router.TrackRequest(containerA); err != nil {
		t.Fatalf("track request A: %v", err)
	}
	if _, err := router.TrackRequest(containerB); err != nil {
		t.Fatalf("track request B: %v", err)
	}

	now := time.Date(2026, 4, 5, 12, 30, 0, 0, time.UTC)
	if err := router.ScheduleRetry(containerA.Header.LocalRequestId, now.Add(-time.Second), errors.New("retry A")); err != nil {
		t.Fatalf("schedule retry A: %v", err)
	}
	if err := router.ScheduleRetry(containerB.Header.LocalRequestId, now.Add(time.Minute), errors.New("retry B")); err != nil {
		t.Fatalf("schedule retry B: %v", err)
	}

	retryable := router.RetryableRequests(now)
	if len(retryable) != 1 {
		t.Fatalf("retryable requests = %d, want 1", len(retryable))
	}
	if retryable[0].Container.Header.LocalRequestId != containerA.Header.LocalRequestId {
		t.Fatalf(
			"retryable local request id = %q, want %q",
			retryable[0].Container.Header.LocalRequestId,
			containerA.Header.LocalRequestId,
		)
	}
}

func TestRouterQueueRequestClearsRetryMetadata(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}

	past := time.Now().Add(-time.Minute)
	if err := router.ScheduleRetry(container.Header.LocalRequestId, past, errors.New("retry later")); err != nil {
		t.Fatalf("schedule retry: %v", err)
	}

	if _, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect); err != nil {
		t.Fatalf("queue request after retry due: %v", err)
	}

	tracked, _ := router.TrackedRequest(container.Header.LocalRequestId)
	if tracked.State != RouterRequestStateOutbound {
		t.Fatalf("state = %q, want %q", tracked.State, RouterRequestStateOutbound)
	}
	if tracked.LastError != "" {
		t.Fatalf("last error = %q, want empty", tracked.LastError)
	}
	if tracked.NextAttemptAt != nil {
		t.Fatalf("expected next attempt time to be cleared")
	}
}

func TestRouterQueueRequestRejectsRetryBeforeDue(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}

	future := time.Now().Add(time.Hour)
	if err := router.ScheduleRetry(container.Header.LocalRequestId, future, errors.New("retry later")); err != nil {
		t.Fatalf("schedule retry: %v", err)
	}

	if _, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect); err == nil {
		t.Fatalf("expected retry-before-due error")
	}
}

func TestRouterPendingOutboundFrames(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()
	if _, err := router.TrackRequest(container); err != nil {
		t.Fatalf("track request: %v", err)
	}
	if _, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect); err != nil {
		t.Fatalf("queue request: %v", err)
	}

	frames := router.PendingOutboundFrames()
	if len(frames) != 1 {
		t.Fatalf("pending outbound frames = %d, want 1", len(frames))
	}

	if err := router.MarkRequestFrameCarried(container.Header.LocalRequestId); err != nil {
		t.Fatalf("mark request frame carried: %v", err)
	}
	frames = router.PendingOutboundFrames()
	if len(frames) != 0 {
		t.Fatalf("pending outbound frames = %d, want 0", len(frames))
	}
}

func TestRouterLocalExchangeWorkflow(t *testing.T) {
	router := NewRouter()
	container := testPAXRequestContainer()

	tracked, err := router.TrackRequest(container)
	if err != nil {
		t.Fatalf("track request: %v", err)
	}
	if tracked.State != RouterRequestStateGenerated {
		t.Fatalf("initial state = %q, want %q", tracked.State, RouterRequestStateGenerated)
	}

	requestFrame, err := router.QueueRequest(container.Header.LocalRequestId, LinkDeliveryMethodDirect)
	if err != nil {
		t.Fatalf("queue request: %v", err)
	}
	if requestFrame.PayloadKind() != LinkPayloadKindRequestContainer {
		t.Fatalf("request frame payload kind = %q, want %q", requestFrame.PayloadKind(), LinkPayloadKindRequestContainer)
	}

	pending := router.PendingOutboundFrames()
	if len(pending) != 1 || pending[0].LinkMessageId != requestFrame.LinkMessageId {
		t.Fatalf("pending outbound frames not updated after queue")
	}

	if err := router.MarkRequestFrameCarried(container.Header.LocalRequestId); err != nil {
		t.Fatalf("mark request frame carried: %v", err)
	}

	syncFrame, err := BuildSyncResponseForRequestFrame(
		requestFrame,
		"KL9K15474QFJ",
		LinkDeliveryMethodPropagated,
	)
	if err != nil {
		t.Fatalf("build sync response for request frame: %v", err)
	}
	if syncFrame.GetReferenceLinkMessageId() != requestFrame.LinkMessageId {
		t.Fatalf(
			"sync frame reference = %q, want %q",
			syncFrame.GetReferenceLinkMessageId(),
			requestFrame.LinkMessageId,
		)
	}

	disposition, err := router.HandleInboundFrame(syncFrame)
	if err != nil {
		t.Fatalf("handle inbound sync frame: %v", err)
	}
	if disposition != RouterFrameDispositionSynchronized {
		t.Fatalf("disposition = %q, want %q", disposition, RouterFrameDispositionSynchronized)
	}

	tracked, ok := router.TrackedRequest(container.Header.LocalRequestId)
	if !ok {
		t.Fatalf("expected tracked request after sync")
	}
	if tracked.State != RouterRequestStateSynchronized {
		t.Fatalf("state = %q, want %q", tracked.State, RouterRequestStateSynchronized)
	}
	if tracked.SyncFrame == nil || tracked.SyncFrame.LinkMessageId != syncFrame.LinkMessageId {
		t.Fatalf("expected tracked sync frame to be stored")
	}
	if tracked.Container.Header.GetSynchronizedRequestId() != "KL9K15474QFJ" {
		t.Fatalf(
			"synchronized request id = %q, want %q",
			tracked.Container.Header.GetSynchronizedRequestId(),
			"KL9K15474QFJ",
		)
	}

	disposition, err = router.HandleInboundFrame(proto.Clone(syncFrame).(*LinkFrame))
	if err != nil {
		t.Fatalf("handle duplicate sync frame: %v", err)
	}
	if disposition != RouterFrameDispositionDuplicate {
		t.Fatalf("duplicate disposition = %q, want %q", disposition, RouterFrameDispositionDuplicate)
	}

	conflicting := proto.Clone(syncFrame).(*LinkFrame)
	conflicting.ReferenceLinkMessageId = nil
	if _, err := router.HandleInboundFrame(conflicting); err == nil {
		t.Fatalf("expected conflicting replay error")
	}
}
