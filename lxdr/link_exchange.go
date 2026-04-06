package lxdr

import "fmt"

func BuildSyncResponseForRequestFrame(
	requestFrame *LinkFrame,
	synchronizedRequestID string,
	method LinkDeliveryMethod,
) (*LinkFrame, error) {
	if requestFrame == nil {
		return nil, fmt.Errorf("request frame must not be nil")
	}
	if err := requestFrame.Validate(); err != nil {
		return nil, err
	}

	request := requestFrame.GetRequestContainer()
	if request == nil {
		return nil, fmt.Errorf("request frame does not carry a request container")
	}

	return NewReferencedSynchronizedResponseLinkFrameForRequest(
		request,
		requestFrame,
		synchronizedRequestID,
		method,
	)
}

func ApplySyncResponseFrameToRequest(
	container *RequestContainer,
	syncFrame *LinkFrame,
) error {
	if container == nil {
		return fmt.Errorf("request container must not be nil")
	}
	if syncFrame == nil {
		return fmt.Errorf("sync frame must not be nil")
	}
	if err := syncFrame.Validate(); err != nil {
		return err
	}

	resp := syncFrame.GetSynchronizedResponse()
	if resp == nil {
		return fmt.Errorf("sync frame does not carry a synchronized response")
	}

	return ApplySynchronizedResponse(container, resp)
}

func ValidateLinkedSyncExchange(
	requestFrame, syncFrame *LinkFrame,
) error {
	answers, err := SyncFrameAnswersRequestFrame(syncFrame, requestFrame)
	if err != nil {
		return err
	}
	if !answers {
		return fmt.Errorf("sync frame does not answer request frame")
	}
	return nil
}
