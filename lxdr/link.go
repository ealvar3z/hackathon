package lxdr

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"google.golang.org/protobuf/proto"
)

type LinkPayloadKind string

const (
	LinkPayloadKindRequestContainer     LinkPayloadKind = "request_container"
	LinkPayloadKindSynchronizedResponse LinkPayloadKind = "synchronized_response"
	LinkPayloadKindCanonicalRegistry    LinkPayloadKind = "canonical_registry"
)

func canonicalLinkMessageID(kind LinkPayloadKind, payload []byte) string {
	sum := sha256.Sum256(payload)
	return "lf1-" + hex.EncodeToString(sum[:])
}

func newBinaryProtoLinkFrame(
	method LinkDeliveryMethod,
	kind LinkPayloadKind,
	payload proto.Message,
	wrap func(string) *LinkFrame,
) (*LinkFrame, error) {
	if !method.IsValid() {
		return nil, fmt.Errorf("invalid link delivery method: %v", method)
	}

	data, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}

	frame := wrap("")
	frame.DeliveryMethod = method
	frame.Representation = LinkRepresentationBinaryProto
	idMaterial := append(
		[]byte(fmt.Sprintf("%s:%d:%d:", kind, method, frame.Representation)),
		data...,
	)
	frame.LinkMessageId = canonicalLinkMessageID(kind, idMaterial)
	if err := frame.Validate(); err != nil {
		return nil, err
	}

	return frame, nil
}

func NewRequestContainerLinkFrame(
	container *RequestContainer,
	method LinkDeliveryMethod,
) (*LinkFrame, error) {
	if container == nil {
		return nil, fmt.Errorf("request container must not be nil")
	}
	if err := container.Validate(); err != nil {
		return nil, err
	}

	return newBinaryProtoLinkFrame(
		method,
		LinkPayloadKindRequestContainer,
		container,
		func(messageID string) *LinkFrame {
			return &LinkFrame{
				LinkMessageId: messageID,
				Payload: &LinkFrame_RequestContainer{
					RequestContainer: container,
				},
			}
		},
	)
}

func NewSynchronizedResponseLinkFrame(
	resp *SynchronizedResponse,
	method LinkDeliveryMethod,
) (*LinkFrame, error) {
	if resp == nil {
		return nil, fmt.Errorf("synchronized response must not be nil")
	}
	if err := resp.Validate(); err != nil {
		return nil, err
	}

	return newBinaryProtoLinkFrame(
		method,
		LinkPayloadKindSynchronizedResponse,
		resp,
		func(messageID string) *LinkFrame {
			return &LinkFrame{
				LinkMessageId: messageID,
				Payload: &LinkFrame_SynchronizedResponse{
					SynchronizedResponse: resp,
				},
			}
		},
	)
}

func NewCanonicalRegistryLinkFrame(
	registry *CanonicalRegistry,
	method LinkDeliveryMethod,
) (*LinkFrame, error) {
	if registry == nil {
		return nil, fmt.Errorf("canonical registry must not be nil")
	}
	if err := registry.Validate(); err != nil {
		return nil, err
	}

	return newBinaryProtoLinkFrame(
		method,
		LinkPayloadKindCanonicalRegistry,
		registry,
		func(messageID string) *LinkFrame {
			return &LinkFrame{
				LinkMessageId: messageID,
				Payload: &LinkFrame_CanonicalRegistry{
					CanonicalRegistry: registry,
				},
			}
		},
	)
}

func (f *LinkFrame) PayloadKind() LinkPayloadKind {
	switch {
	case f.GetRequestContainer() != nil:
		return LinkPayloadKindRequestContainer
	case f.GetSynchronizedResponse() != nil:
		return LinkPayloadKindSynchronizedResponse
	case f.GetCanonicalRegistry() != nil:
		return LinkPayloadKindCanonicalRegistry
	default:
		return ""
	}
}
