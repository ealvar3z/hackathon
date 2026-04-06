package lxdr

import (
	"fmt"
	"time"
)

type RouterRequestState string

const (
	RouterRequestStateInvalid      RouterRequestState = "invalid"
	RouterRequestStateGenerated    RouterRequestState = "generated"
	RouterRequestStateOutbound     RouterRequestState = "outbound"
	RouterRequestStateCarried      RouterRequestState = "carried"
	RouterRequestStateSynchronized RouterRequestState = "synchronized"
	RouterRequestStateFailed       RouterRequestState = "failed"
)

func (s RouterRequestState) IsValid() bool {
	switch s {
	case RouterRequestStateGenerated,
		RouterRequestStateOutbound,
		RouterRequestStateCarried,
		RouterRequestStateSynchronized,
		RouterRequestStateFailed:
		return true
	default:
		return false
	}
}

type RouterFrameDisposition string

const (
	RouterFrameDispositionInvalid      RouterFrameDisposition = "invalid"
	RouterFrameDispositionAccepted     RouterFrameDisposition = "accepted"
	RouterFrameDispositionDuplicate    RouterFrameDisposition = "duplicate"
	RouterFrameDispositionSynchronized RouterFrameDisposition = "synchronized"
)

type SeenFrameStatus string

const (
	SeenFrameStatusInvalid     SeenFrameStatus = "invalid"
	SeenFrameStatusNew         SeenFrameStatus = "new"
	SeenFrameStatusDuplicate   SeenFrameStatus = "duplicate"
	SeenFrameStatusConflicting SeenFrameStatus = "conflicting"
)

type SyncResponseStatus string

const (
	SyncResponseStatusInvalid     SyncResponseStatus = "invalid"
	SyncResponseStatusNew         SyncResponseStatus = "new"
	SyncResponseStatusRepeated    SyncResponseStatus = "repeated"
	SyncResponseStatusConflicting SyncResponseStatus = "conflicting"
)

type TrackedRequest struct {
	Container     *RequestContainer
	RequestFrame  *LinkFrame
	SyncFrame     *LinkFrame
	State         RouterRequestState
	AttemptCount  uint32
	LastError     string
	NextAttemptAt *time.Time
}

type Router struct {
	requests   map[string]*TrackedRequest
	seenFrames map[string]*LinkFrame
}

func NewRouter() *Router {
	return &Router{
		requests:   map[string]*TrackedRequest{},
		seenFrames: map[string]*LinkFrame{},
	}
}

func (r *Router) TrackRequest(container *RequestContainer) (*TrackedRequest, error) {
	if r == nil {
		return nil, fmt.Errorf("router must not be nil")
	}
	if container == nil {
		return nil, fmt.Errorf("request container must not be nil")
	}
	if err := container.Validate(); err != nil {
		return nil, err
	}

	localID := container.Header.LocalRequestId
	if _, exists := r.requests[localID]; exists {
		return nil, fmt.Errorf("request %q is already tracked", localID)
	}

	tracked := &TrackedRequest{
		Container: container,
		State:     RouterRequestStateGenerated,
	}
	r.requests[localID] = tracked
	return tracked, nil
}

func (r *Router) QueueRequest(
	localRequestID string,
	method LinkDeliveryMethod,
) (*LinkFrame, error) {
	if r == nil {
		return nil, fmt.Errorf("router must not be nil")
	}
	tracked, ok := r.requests[localRequestID]
	if !ok {
		return nil, fmt.Errorf("request %q is not tracked", localRequestID)
	}
	if tracked.State == RouterRequestStateSynchronized {
		return nil, fmt.Errorf("request %q is already synchronized", localRequestID)
	}
	if tracked.State == RouterRequestStateFailed && tracked.NextAttemptAt != nil {
		if time.Now().Before(*tracked.NextAttemptAt) {
			return nil, fmt.Errorf(
				"request %q retry is not due until %s",
				localRequestID,
				tracked.NextAttemptAt.Format(time.RFC3339),
			)
		}
	}

	frame, err := NewRequestContainerLinkFrame(tracked.Container, method)
	if err != nil {
		return nil, err
	}

	tracked.RequestFrame = frame
	tracked.State = RouterRequestStateOutbound
	tracked.AttemptCount++
	tracked.LastError = ""
	tracked.NextAttemptAt = nil
	r.seenFrames[frame.LinkMessageId] = frame
	return frame, nil
}

func (r *Router) MarkRequestFrameCarried(localRequestID string) error {
	if r == nil {
		return fmt.Errorf("router must not be nil")
	}
	tracked, ok := r.requests[localRequestID]
	if !ok {
		return fmt.Errorf("request %q is not tracked", localRequestID)
	}
	if tracked.RequestFrame == nil {
		return fmt.Errorf("request %q has no queued request frame", localRequestID)
	}
	if tracked.State == RouterRequestStateSynchronized {
		return fmt.Errorf("request %q is already synchronized", localRequestID)
	}

	tracked.State = RouterRequestStateCarried
	return nil
}

func (r *Router) MarkRequestFailed(localRequestID string, err error) error {
	if r == nil {
		return fmt.Errorf("router must not be nil")
	}
	tracked, ok := r.requests[localRequestID]
	if !ok {
		return fmt.Errorf("request %q is not tracked", localRequestID)
	}

	tracked.State = RouterRequestStateFailed
	if err != nil {
		tracked.LastError = err.Error()
	}
	tracked.NextAttemptAt = nil
	return nil
}

func (r *Router) ScheduleRetry(
	localRequestID string,
	nextAttemptAt time.Time,
	err error,
) error {
	if r == nil {
		return fmt.Errorf("router must not be nil")
	}
	tracked, ok := r.requests[localRequestID]
	if !ok {
		return fmt.Errorf("request %q is not tracked", localRequestID)
	}

	tracked.State = RouterRequestStateFailed
	tracked.NextAttemptAt = &nextAttemptAt
	if err != nil {
		tracked.LastError = err.Error()
	}
	return nil
}

func (r *Router) TrackedRequest(localRequestID string) (*TrackedRequest, bool) {
	if r == nil {
		return nil, false
	}
	tracked, ok := r.requests[localRequestID]
	return tracked, ok
}

func (r *Router) PendingOutboundFrames() []*LinkFrame {
	if r == nil {
		return nil
	}

	frames := []*LinkFrame{}
	for _, tracked := range r.requests {
		if tracked.State == RouterRequestStateOutbound && tracked.RequestFrame != nil {
			frames = append(frames, tracked.RequestFrame)
		}
	}
	return frames
}

func (r *Router) ReadyForRetry(
	localRequestID string,
	now time.Time,
) (bool, error) {
	if r == nil {
		return false, fmt.Errorf("router must not be nil")
	}
	tracked, ok := r.requests[localRequestID]
	if !ok {
		return false, fmt.Errorf("request %q is not tracked", localRequestID)
	}
	if tracked.State != RouterRequestStateFailed {
		return false, nil
	}
	if tracked.NextAttemptAt == nil {
		return false, nil
	}
	return !now.Before(*tracked.NextAttemptAt), nil
}

func (r *Router) RetryableRequests(now time.Time) []*TrackedRequest {
	if r == nil {
		return nil
	}

	retryable := []*TrackedRequest{}
	for localRequestID, tracked := range r.requests {
		ready, err := r.ReadyForRetry(localRequestID, now)
		if err == nil && ready {
			retryable = append(retryable, tracked)
		}
	}
	return retryable
}

func (r *Router) SeenFrameStatus(frame *LinkFrame) (SeenFrameStatus, error) {
	if r == nil {
		return SeenFrameStatusInvalid, fmt.Errorf("router must not be nil")
	}
	if frame == nil {
		return SeenFrameStatusInvalid, fmt.Errorf("link frame must not be nil")
	}
	if err := frame.Validate(); err != nil {
		return SeenFrameStatusInvalid, err
	}

	existing, exists := r.seenFrames[frame.LinkMessageId]
	if !exists {
		return SeenFrameStatusNew, nil
	}

	duplicate, err := FramesAreDuplicates(existing, frame)
	if err != nil {
		return SeenFrameStatusConflicting, nil
	}
	if duplicate {
		return SeenFrameStatusDuplicate, nil
	}

	return SeenFrameStatusConflicting, nil
}

func (r *Router) RequestFrameForSyncResponse(
	syncFrame *LinkFrame,
) (*LinkFrame, *TrackedRequest, error) {
	if r == nil {
		return nil, nil, fmt.Errorf("router must not be nil")
	}
	if syncFrame == nil {
		return nil, nil, fmt.Errorf("sync frame must not be nil")
	}
	if err := syncFrame.Validate(); err != nil {
		return nil, nil, err
	}

	resp := syncFrame.GetSynchronizedResponse()
	if resp == nil {
		return nil, nil, fmt.Errorf("sync frame does not carry a synchronized response")
	}

	tracked, ok := r.requests[resp.LocalRequestId]
	if !ok {
		return nil, nil, fmt.Errorf(
			"no tracked request for synchronized response local request ID %q",
			resp.LocalRequestId,
		)
	}
	if tracked.RequestFrame == nil {
		return nil, nil, fmt.Errorf("tracked request %q has no request frame", resp.LocalRequestId)
	}

	return tracked.RequestFrame, tracked, nil
}

func (r *Router) SyncResponseStatus(syncFrame *LinkFrame) (SyncResponseStatus, error) {
	if r == nil {
		return SyncResponseStatusInvalid, fmt.Errorf("router must not be nil")
	}
	requestFrame, tracked, err := r.RequestFrameForSyncResponse(syncFrame)
	if err != nil {
		return SyncResponseStatusInvalid, err
	}

	seenStatus, err := r.SeenFrameStatus(syncFrame)
	if err != nil {
		return SyncResponseStatusInvalid, err
	}
	switch seenStatus {
	case SeenFrameStatusConflicting:
		return SyncResponseStatusConflicting, nil
	case SeenFrameStatusDuplicate:
		return SyncResponseStatusRepeated, nil
	}

	if err := ValidateLinkedSyncExchange(requestFrame, syncFrame); err != nil {
		return SyncResponseStatusConflicting, nil
	}

	resp := syncFrame.GetSynchronizedResponse()
	if tracked.Container.Header.IsSynchronized() {
		current := tracked.Container.Header.GetSynchronizedRequestId()
		switch {
		case current == resp.SynchronizedRequestId:
			return SyncResponseStatusRepeated, nil
		default:
			return SyncResponseStatusConflicting, nil
		}
	}

	return SyncResponseStatusNew, nil
}

func (r *Router) HandleInboundFrame(
	frame *LinkFrame,
) (RouterFrameDisposition, error) {
	if r == nil {
		return RouterFrameDispositionInvalid, fmt.Errorf("router must not be nil")
	}
	if frame == nil {
		return RouterFrameDispositionInvalid, fmt.Errorf("link frame must not be nil")
	}
	if err := frame.Validate(); err != nil {
		return RouterFrameDispositionInvalid, err
	}

	seenStatus, err := r.SeenFrameStatus(frame)
	if err != nil {
		return RouterFrameDispositionInvalid, err
	}
	switch seenStatus {
	case SeenFrameStatusDuplicate:
		return RouterFrameDispositionDuplicate, nil
	case SeenFrameStatusConflicting:
		return RouterFrameDispositionInvalid, fmt.Errorf(
			"conflicting link frames share link_message_id %q",
			frame.LinkMessageId,
		)
	}

	switch {
	case frame.GetSynchronizedResponse() != nil:
		requestFrame, tracked, err := r.RequestFrameForSyncResponse(frame)
		if err != nil {
			return RouterFrameDispositionInvalid, err
		}
		status, err := r.SyncResponseStatus(frame)
		if err != nil {
			return RouterFrameDispositionInvalid, err
		}
		switch status {
		case SyncResponseStatusRepeated:
			return RouterFrameDispositionDuplicate, nil
		case SyncResponseStatusConflicting:
			return RouterFrameDispositionInvalid, fmt.Errorf(
				"conflicting synchronized response for local request ID %q",
				tracked.Container.Header.LocalRequestId,
			)
		}
		if err := ValidateLinkedSyncExchange(requestFrame, frame); err != nil {
			return RouterFrameDispositionInvalid, err
		}
		if err := ApplySyncResponseFrameToRequest(tracked.Container, frame); err != nil {
			return RouterFrameDispositionInvalid, err
		}

		r.seenFrames[frame.LinkMessageId] = frame
		tracked.SyncFrame = frame
		tracked.State = RouterRequestStateSynchronized
		tracked.LastError = ""
		return RouterFrameDispositionSynchronized, nil

	case frame.GetRequestContainer() != nil, frame.GetCanonicalRegistry() != nil:
		r.seenFrames[frame.LinkMessageId] = frame
		return RouterFrameDispositionAccepted, nil

	default:
		return RouterFrameDispositionInvalid, fmt.Errorf("unsupported link frame payload")
	}
}
