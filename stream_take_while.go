package streamer

type takeWhileStream struct {
	input  Iterator
	takeFn func(TItem) bool

	currentItem TItem
	taken       bool
}

func newTakeWhileStream(input Iterator, takeFn func(TItem) bool) (res *takeWhileStream) {
	res = &takeWhileStream{
		input:  input,
		takeFn: takeFn,
	}
	return
}

func (tw *takeWhileStream) Next() bool {
	if tw.taken {
		return false
	}
	if tw.currentItem != nil {
		return true
	}
	next := tw.input.Next()
	if !next {
		return false
	}
	item := tw.input.Value()
	if tw.takeFn(item) {
		tw.currentItem = item
		return true
	}
	tw.taken = true
	return false
}

func (tw *takeWhileStream) Value() TItem {
	value := tw.currentItem
	tw.currentItem = nil
	return value
}
