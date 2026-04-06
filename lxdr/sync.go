package lxdr

import (
	"errors"
	"fmt"
)

func (h *RequestHeader) IsSynchronized() bool {
	if h == nil {
		return false
	}
	return h.GetSynchronizedRequestId() != ""
}

func NewSynchronizedResponseForRequest(
	container *RequestContainer,
	synchronizedRequestID string,
) (*SynchronizedResponse, error) {
	if container == nil {
		return nil, fmt.Errorf("request container must not be nil")
	}
	if err := container.Validate(); err != nil {
		return nil, err
	}
	if synchronizedRequestID == "" {
		return nil, errors.New("synchronized request ID is required")
	}

	resp := &SynchronizedResponse{
		LocalRequestId:        container.Header.LocalRequestId,
		SynchronizedRequestId: synchronizedRequestID,
	}
	if err := resp.Validate(); err != nil {
		return nil, err
	}

	return resp, nil
}

func ApplySynchronizedResponse(
	container *RequestContainer,
	resp *SynchronizedResponse,
) error {
	if container == nil {
		return fmt.Errorf("request container must not be nil")
	}
	if err := container.Validate(); err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("synchronized response must not be nil")
	}
	if err := resp.Validate(); err != nil {
		return err
	}
	if container.Header.LocalRequestId != resp.LocalRequestId {
		return fmt.Errorf(
			"synchronized response local request ID mismatch: header=%q response=%q",
			container.Header.LocalRequestId,
			resp.LocalRequestId,
		)
	}
	if container.Header.IsSynchronized() {
		current := container.Header.GetSynchronizedRequestId()
		if current == resp.SynchronizedRequestId {
			return nil
		}
		return fmt.Errorf(
			"conflicting synchronized request ID: header=%q response=%q",
			current,
			resp.SynchronizedRequestId,
		)
	}

	container.Header.SynchronizedRequestId = &resp.SynchronizedRequestId
	return nil
}

func NewSynchronizedResponseLinkFrameForRequest(
	container *RequestContainer,
	synchronizedRequestID string,
	method LinkDeliveryMethod,
) (*LinkFrame, error) {
	resp, err := NewSynchronizedResponseForRequest(
		container,
		synchronizedRequestID,
	)
	if err != nil {
		return nil, err
	}

	return NewSynchronizedResponseLinkFrame(resp, method)
}
