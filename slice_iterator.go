package streamer

type SliceIterator struct {
	input   []interface{}
	current int
}

func NewSliceIterator(input []interface{}) *SliceIterator {
	res := &SliceIterator{
		input:   input,
		current: -1,
	}

	return res
}

func (it *SliceIterator) Next() (interface{}, bool) {
	if it.current+1 < len(it.input) {
		it.current++
		return it.input[it.current], true
	}
	return nil, false
}
