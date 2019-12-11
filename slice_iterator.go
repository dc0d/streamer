package streamer

type SliceIterator struct {
	input       []TItem
	current     int
	currentItem TItem
}

func NewSliceIterator(input []TItem) *SliceIterator {
	res := &SliceIterator{
		input:   input,
		current: -1,
	}

	return res
}

func (it *SliceIterator) Next() bool {
	if it.currentItem != nil {
		return true
	}

	if it.current+1 < len(it.input) {
		it.current++
		it.currentItem = it.input[it.current]
		return true
	}
	return false
}

func (it *SliceIterator) Value() TItem {
	item := it.currentItem
	it.currentItem = nil
	return item
}
