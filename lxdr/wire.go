package lxdr

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

type validatableMessage interface {
	proto.Message
	Validate() error
}

func marshalValidatedBinary[T validatableMessage](
	msg T,
	nilLabel string,
) ([]byte, error) {
	if any(msg) == nil {
		return nil, fmt.Errorf("%s must not be nil", nilLabel)
	}
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return proto.Marshal(msg)
}

func unmarshalValidatedBinary[T validatableMessage](
	data []byte,
	target T,
) (T, error) {
	if err := proto.Unmarshal(data, target); err != nil {
		return target, err
	}
	if err := target.Validate(); err != nil {
		return target, err
	}
	return target, nil
}

func MarshalRequestContainerBinary(container *RequestContainer) ([]byte, error) {
	return marshalValidatedBinary(container, "request container")
}

func UnmarshalRequestContainerBinary(data []byte) (*RequestContainer, error) {
	return unmarshalValidatedBinary(data, &RequestContainer{})
}

func MarshalSynchronizedResponseBinary(resp *SynchronizedResponse) ([]byte, error) {
	return marshalValidatedBinary(resp, "synchronized response")
}

func UnmarshalSynchronizedResponseBinary(data []byte) (*SynchronizedResponse, error) {
	return unmarshalValidatedBinary(data, &SynchronizedResponse{})
}

func MarshalCanonicalRegistryBinary(registry *CanonicalRegistry) ([]byte, error) {
	return marshalValidatedBinary(registry, "canonical registry")
}

func UnmarshalCanonicalRegistryBinary(data []byte) (*CanonicalRegistry, error) {
	return unmarshalValidatedBinary(data, &CanonicalRegistry{})
}

func MarshalLinkFrameBinary(frame *LinkFrame) ([]byte, error) {
	return marshalValidatedBinary(frame, "link frame")
}

func UnmarshalLinkFrameBinary(data []byte) (*LinkFrame, error) {
	return unmarshalValidatedBinary(data, &LinkFrame{})
}
