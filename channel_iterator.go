package streamer

import "time"

type ChannelIterator struct {
	input   <-chan interface{}
	timeout time.Duration
}

func NewChannelIterator(input <-chan interface{}, timeout time.Duration) (res *ChannelIterator) {
	res = &ChannelIterator{
		input:   input,
		timeout: timeout,
	}
	return
}

func (ci *ChannelIterator) Next() (interface{}, bool) {
	if ci.timeout > 0 {
		select {
		case v, ok := <-ci.input:
			if !ok {
				return nil, false
			}
			return v, true
		case <-time.After(ci.timeout):
			return nil, false
		}
	}

	if ci.timeout == 0 {
		select {
		case v, ok := <-ci.input:
			if !ok {
				return nil, false
			}
			return v, true
		default:
			return nil, false
		}
	}

	select {
	case v, ok := <-ci.input:
		if !ok {
			return nil, false
		}
		return v, true
	}
}
