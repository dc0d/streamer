package streamer

type takeStream struct {
	input     Iterator
	takeCount int

	currentItem TItem
}

func newTakeStream(input Iterator, takeCount int) (res *takeStream) {
	res = &takeStream{
		input:     input,
		takeCount: takeCount,
	}
	return
}

func (tc *takeStream) Next() bool {
	if tc.currentItem != nil {
		return true
	}
	if tc.takeCount == 0 {
		return false
	}
	next := tc.input.Next()
	if !next {
		return false
	}
	tc.currentItem = tc.input.Value()
	tc.takeCount--
	return true
}

func (tc *takeStream) Value() TItem {
	value := tc.currentItem
	tc.currentItem = nil
	return value
}
