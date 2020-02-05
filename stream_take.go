package streamer

type takeStream struct {
	input     Iterator
	takeCount int
}

func newTakeStream(input Iterator, takeCount int) (res *takeStream) {
	res = &takeStream{
		input:     input,
		takeCount: takeCount,
	}
	return
}

func (tc *takeStream) Next() (interface{}, bool) {
	if tc.takeCount == 0 {
		return nil, false
	}
	item, ok := tc.input.Next()
	if !ok {
		return item, false
	}
	tc.takeCount--
	return item, true
}
