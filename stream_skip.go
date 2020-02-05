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

func (sc *skipStream) Next() (interface{}, bool) {
	if sc.skipCount != 0 {
		ok := true
		for i := 0; i < sc.skipCount && ok; i++ {
			_, ok = sc.input.Next()
		}
		sc.skipCount = 0
	}
	return sc.input.Next()
}
