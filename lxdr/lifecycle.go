package lxdr

import "fmt"

type RequestLifecycleState string

const (
	RequestLifecycleStateInvalid      RequestLifecycleState = "invalid"
	RequestLifecycleStateLocalOnly    RequestLifecycleState = "local_only"
	RequestLifecycleStateCarried      RequestLifecycleState = "carried"
	RequestLifecycleStateSynchronized RequestLifecycleState = "synchronized"
)

func (s RequestLifecycleState) IsValid() bool {
	switch s {
	case RequestLifecycleStateLocalOnly,
		RequestLifecycleStateCarried,
		RequestLifecycleStateSynchronized:
		return true
	default:
		return false
	}
}

func DetermineRequestLifecycleState(
	container *RequestContainer,
	carried bool,
) (RequestLifecycleState, error) {
	if container == nil {
		return RequestLifecycleStateInvalid, fmt.Errorf("request container must not be nil")
	}
	if err := container.Validate(); err != nil {
		return RequestLifecycleStateInvalid, err
	}
	if container.Header.IsSynchronized() {
		return RequestLifecycleStateSynchronized, nil
	}
	if carried {
		return RequestLifecycleStateCarried, nil
	}
	return RequestLifecycleStateLocalOnly, nil
}

func DetermineRequestLifecycleStateFromLinkFrame(
	container *RequestContainer,
	frame *LinkFrame,
) (RequestLifecycleState, error) {
	if frame == nil {
		return RequestLifecycleStateInvalid, fmt.Errorf("link frame must not be nil")
	}
	if err := frame.Validate(); err != nil {
		return RequestLifecycleStateInvalid, err
	}

	payload := frame.GetRequestContainer()
	if payload == nil {
		return RequestLifecycleStateInvalid, fmt.Errorf("link frame does not carry a request container")
	}
	if payload.Header.LocalRequestId != container.Header.LocalRequestId {
		return RequestLifecycleStateInvalid, fmt.Errorf(
			"request container local request ID mismatch: container=%q frame=%q",
			container.Header.LocalRequestId,
			payload.Header.LocalRequestId,
		)
	}

	return DetermineRequestLifecycleState(container, true)
}

func ApplySynchronizedResponseAndDetermineState(
	container *RequestContainer,
	resp *SynchronizedResponse,
) (RequestLifecycleState, error) {
	if err := ApplySynchronizedResponse(container, resp); err != nil {
		return RequestLifecycleStateInvalid, err
	}
	return DetermineRequestLifecycleState(container, false)
}
