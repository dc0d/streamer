package streamer

type takeWhileStream struct {
	input  Iterator
	takeFn func(interface{}) bool

	taken bool
}

func newTakeWhileStream(input Iterator, takeFn func(interface{}) bool) (res *takeWhileStream) {
	res = &takeWhileStream{
		input:  input,
		takeFn: takeFn,
	}
	return
}

func (tw *takeWhileStream) Next() (interface{}, bool) {
	if tw.taken {
		return nil, false
	}
	item, ok := tw.input.Next()
	if !ok {
		return nil, false
	}
	if tw.takeFn(item) {
		return item, true
	}
	tw.taken = true
	return nil, false
}
