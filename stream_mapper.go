package streamer

type mapperStream struct {
	input Iterator
	mapFn func(x interface{}) interface{}
}

func newMapperStream(input Iterator, mapFn func(x interface{}) interface{}) (res *mapperStream) {
	res = &mapperStream{
		input: input,
		mapFn: mapFn,
	}
	return
}

func (ms *mapperStream) Next() (interface{}, bool) {
	next, ok := ms.input.Next()
	if !ok {
		return nil, false
	}
	return ms.mapFn(next), true
}
