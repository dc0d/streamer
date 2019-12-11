package streamer

import "time"

type ChannelIterator struct {
	input   <-chan TItem
	timeout time.Duration

	currentItem TItem
}

func NewChannelIterator(input <-chan TItem, timeout time.Duration) (res *ChannelIterator) {
	res = &ChannelIterator{
		input:   input,
		timeout: timeout,
	}
	return
}

func (ci *ChannelIterator) Next() bool {
	if ci.currentItem != nil {
		return true
	}

	if ci.timeout > 0 {
		select {
		case v, ok := <-ci.input:
			if !ok {
				return false
			}
			ci.currentItem = v
			return true
		case <-time.After(ci.timeout):
			return false
		}
	}

	if ci.timeout == 0 {
		select {
		case v, ok := <-ci.input:
			if !ok {
				return false
			}
			ci.currentItem = v
			return true
		default:
			return false
		}
	}

	select {
	case v, ok := <-ci.input:
		if !ok {
			return false
		}
		ci.currentItem = v
		return true
	}
}

func (ci *ChannelIterator) Value() TItem {
	item := ci.currentItem
	ci.currentItem = nil
	return item
}
