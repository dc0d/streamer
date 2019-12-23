package streamer

type skipStream struct {
	input     Iterator
	skipCount int
}

func newSkipStream(input Iterator, skipCount int) (res *skipStream) {
	res = &skipStream{
		input:     input,
		skipCount: skipCount,
	}
	return
}

func (sc *skipStream) Next() bool {
	if sc.skipCount != 0 {
		for i := 0; i < sc.skipCount && sc.input.Next(); i++ {
			sc.input.Value()
		}
		sc.skipCount = 0
	}
	return sc.input.Next()
}

func (sc *skipStream) Value() interface{} { return sc.input.Value() }
