package lxdr

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

func FramesAreDuplicates(a, b *LinkFrame) (bool, error) {
	if a == nil || b == nil {
		return false, fmt.Errorf("link frames must not be nil")
	}
	if err := a.Validate(); err != nil {
		return false, err
	}
	if err := b.Validate(); err != nil {
		return false, err
	}

	if a.LinkMessageId != b.LinkMessageId {
		return false, nil
	}
	if !proto.Equal(a, b) {
		return false, fmt.Errorf(
			"conflicting link frames share link_message_id %q",
			a.LinkMessageId,
		)
	}

	return true, nil
}

func FrameRefersTo(frame, target *LinkFrame) (bool, error) {
	if frame == nil || target == nil {
		return false, fmt.Errorf("link frames must not be nil")
	}
	if err := frame.Validate(); err != nil {
		return false, err
	}
	if err := target.Validate(); err != nil {
		return false, err
	}

	return frame.GetReferenceLinkMessageId() == target.LinkMessageId &&
		frame.GetReferenceLinkMessageId() != "", nil
}

func SyncFrameAnswersRequestFrame(
	syncFrame, requestFrame *LinkFrame,
) (bool, error) {
	if syncFrame == nil || requestFrame == nil {
		return false, fmt.Errorf("link frames must not be nil")
	}
	if err := syncFrame.Validate(); err != nil {
		return false, err
	}
	if err := requestFrame.Validate(); err != nil {
		return false, err
	}

	resp := syncFrame.GetSynchronizedResponse()
	if resp == nil {
		return false, fmt.Errorf("sync frame does not carry a synchronized response")
	}
	request := requestFrame.GetRequestContainer()
	if request == nil {
		return false, fmt.Errorf("request frame does not carry a request container")
	}

	if resp.LocalRequestId != request.Header.LocalRequestId {
		return false, nil
	}

	reference := syncFrame.GetReferenceLinkMessageId()
	if reference != "" && reference != requestFrame.LinkMessageId {
		return false, nil
	}

	return true, nil
}
