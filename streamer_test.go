package streamer_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/dc0d/streamer"

	assert "github.com/stretchr/testify/require"
)

type (
	T = interface{}
)

func Test_stream_map(t *testing.T) {
	type (
		expectation struct {
			input          []T
			expectedOutput []T
			mapFn          func(T) T
		}
	)

	var (
		expectations = []expectation{
			{nil, nil, func(x T) T { return x }},
			{[]T{}, []T{}, func(x T) T { return x }},
			{[]T{1, 2, 3}, []T{1, 2, 3}, func(x T) T { return x }},
			{[]T{1, 2, 3}, []T{2, 4, 6}, func(x T) T { return x.(int) * 2 }},
		}
	)

	for i, exp := range expectations {
		var (
			input          = exp.input
			expectedOutput = exp.expectedOutput
			mapFn          = exp.mapFn
		)

		t.Run(fmt.Sprintf("stream map test case %v", i+1), func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input)
				stream                     = streamer.NewStream(iterator)

				_ streamer.Iterator = stream
			)

			stream = stream.Map(mapFn)

			index := 0
			for item, ok := stream.Next(); ok; item, ok = stream.Next() {
				assert.Equal(expectedOutput[index], item)
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_chunk_by(t *testing.T) {
	type (
		expectation struct {
			input          []T
			expectedOutput []T
			chunkFn        func(T) T
		}
	)

	var (
		expectations = []expectation{
			{
				nil,
				nil,
				func(x T) T { return x },
			},
			{
				[]T{1, 2, 3},
				[]T{
					[]T{1},
					[]T{2},
					[]T{3},
				},
				func(x T) T { return x },
			},
			{
				[]T{1, 3, 4, 6, 12, 7, 13, 18},
				[]T{
					[]T{1, 3},
					[]T{4, 6, 12},
					[]T{7, 13},
					[]T{18},
				},
				func(x T) T { return x.(int) % 2 },
			},
			{
				[]T{1, 2, 2, 3, 4, 4, 6, 7, 7},
				[]T{
					[]T{1},
					[]T{2, 2},
					[]T{3},
					[]T{4, 4},
					[]T{6},
					[]T{7, 7},
				},
				func(x T) T { return x.(int) % 3 },
			},
			{
				[]T{"str1", "123", "str2", "str3", "10", "20"},
				[]T{
					[]T{"str1"},
					[]T{"123"},
					[]T{"str2", "str3"},
					[]T{"10", "20"},
				},
				func(x T) T {
					_, err := strconv.Atoi(x.(string))
					return err == nil
				},
			},
		}
	)

	for i, exp := range expectations {
		var (
			input          = exp.input
			expectedOutput = exp.expectedOutput
			chunkFn        = exp.chunkFn
		)

		t.Run(fmt.Sprintf("stream chunk by, test case %v", i+1), func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.ChunkBy(chunkFn)

			index := 0
			for item, ok := stream.Next(); ok; item, ok = stream.Next() {
				assert.Equal(expectedOutput[index], item)
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_chunk_every(t *testing.T) {
	type (
		expectation struct {
			input          []T
			expectedOutput []T
			chunkSize      int
		}
	)

	var (
		expectations = []expectation{
			{
				nil,
				nil,
				100,
			},
			{
				[]T{1, 2, 3},
				[]T{
					[]T{1},
					[]T{2},
					[]T{3},
				},
				1,
			},
			{
				[]T{1, 2, 3, 4, 5},
				[]T{
					[]T{1, 2},
					[]T{3, 4},
					[]T{5},
				},
				2,
			},
			{
				[]T{1, 2, 3, 4, 5, 6, 7},
				[]T{
					[]T{1, 2, 3},
					[]T{4, 5, 6},
					[]T{7},
				},
				3,
			},
		}
	)

	for i, exp := range expectations {
		var (
			input          = exp.input
			expectedOutput = exp.expectedOutput
			chunkSize      = exp.chunkSize
		)

		t.Run(fmt.Sprintf("stream chunk every, test case %v", i+1), func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.ChunkEvery(chunkSize)

			index := 0
			for item, ok := stream.Next(); ok; item, ok = stream.Next() {
				assert.Equal(expectedOutput[index], item)
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_skip(t *testing.T) {
	type (
		expectation struct {
			input          []T
			expectedOutput []T
			skipCount      int
		}
	)

	var (
		expectations = []expectation{
			{
				nil,
				nil,
				100,
			},
			{
				[]T{1, 2, 3},
				[]T{1, 2, 3},
				0,
			},
			{
				[]T{1, 2, 3},
				[]T{2, 3},
				1,
			},
			{
				[]T{1, 2, 3},
				[]T{3},
				2,
			},
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{3, 4, 4},
				3,
			},
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{4, 4},
				4,
			},
		}
	)

	for i, exp := range expectations {
		var (
			input          = exp.input
			expectedOutput = exp.expectedOutput
			skipCount      = exp.skipCount
		)

		t.Run(fmt.Sprintf("stream skip, test case %v", i+1), func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.Skip(skipCount)

			index := 0
			for item, ok := stream.Next(); ok; item, ok = stream.Next() {
				assert.Equal(expectedOutput[index], item)
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_skip_while(t *testing.T) {
	type (
		expectation struct {
			input          []T
			expectedOutput []T
			skipFn         func(T) bool
		}
	)

	var (
		expectations = []expectation{
			{
				nil,
				nil,
				func(elem T) bool { return false },
			},
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{1, 2, 3, 3, 4, 4},
				func(elem T) bool { return false },
			},
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{},
				func(elem T) bool { return true },
			},
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{3, 3, 4, 4},
				func(elem T) bool { return elem.(int) < 3 },
			},
			{
				[]T{2, 2, 2, 4, 4, 3, 9},
				[]T{3, 9},
				func(elem T) bool { return elem.(int)%2 == 0 },
			},
		}
	)

	for i, exp := range expectations {
		var (
			input          = exp.input
			expectedOutput = exp.expectedOutput
			skipFn         = exp.skipFn
		)

		t.Run(fmt.Sprintf("stream skip while, test case %v", i+1), func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.SkipWhile(skipFn)

			index := 0
			for item, ok := stream.Next(); ok; item, ok = stream.Next() {
				if len(expectedOutput) > index {
					assert.Equal(expectedOutput[index], item)
					index++
				}
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_filter(t *testing.T) {
	type (
		expectation struct {
			input          []T
			expectedOutput []T
			filterFn       func(T) bool
		}
	)

	var (
		expectations = []expectation{
			{
				nil,
				nil,
				func(elem T) bool { return true },
			},
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{1, 2, 3, 3, 4, 4},
				func(elem T) bool { return true },
			},
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{},
				func(elem T) bool { return false },
			},
			{
				[]T{1, 3, 11, 9, 21, 17},
				[]T{3, 9, 21},
				func(elem T) bool { return elem.(int)%3 == 0 },
			},
			{
				[]T{2, 2, 2, 4, 4, 3, 9, 8},
				[]T{4, 4, 8},
				func(elem T) bool { return elem.(int)%4 == 0 },
			},
		}
	)

	for i, exp := range expectations {
		var (
			input          = exp.input
			expectedOutput = exp.expectedOutput
			filterFn       = exp.filterFn
		)

		t.Run(fmt.Sprintf("stream filter, test case %v", i+1), func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.Filter(filterFn)

			index := 0
			for item, ok := stream.Next(); ok; item, ok = stream.Next() {
				if len(expectedOutput) > index {
					assert.Equal(expectedOutput[index], item)
					index++
				}
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_take(t *testing.T) {
	type (
		expectation struct {
			input          []T
			expectedOutput []T
			takeCount      int
		}
	)

	var (
		expectations = []expectation{
			{
				[]T{1, 2, 3},
				[]T{1},
				1,
			},
			{
				[]T{1, 2, 3},
				[]T{1, 2},
				2,
			},
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{1, 2, 3},
				3,
			},
			{
				[]T{},
				[]T{},
				4,
			},
			{
				[]T{1, 2},
				[]T{1, 2},
				4,
			},
		}
	)

	for i, exp := range expectations {
		var (
			input          = exp.input
			expectedOutput = exp.expectedOutput
			takeCount      = exp.takeCount
		)

		t.Run(fmt.Sprintf("stream take, test case %v", i+1), func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.Take(takeCount)

			index := 0
			for item, ok := stream.Next(); ok; item, ok = stream.Next() {
				assert.Equal(expectedOutput[index], item)
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_take_while(t *testing.T) {
	type (
		expectation struct {
			input          []T
			expectedOutput []T
			takeFn         func(T) bool
		}
	)

	var (
		expectations = []expectation{
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{1, 2, 3, 3, 4, 4},
				func(elem T) bool { return true },
			},
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{},
				func(elem T) bool { return false },
			},
			{
				[]T{1, 2, 3, 3, 4, 4},
				[]T{1, 2, 3, 3},
				func(elem T) bool { return elem.(int) < 4 },
			},
			{
				[]T{2, 2, 2, 4, 4, 3, 9},
				[]T{2, 2, 2, 4, 4},
				func(elem T) bool { return elem.(int)%2 == 0 },
			},
		}
	)

	for i, exp := range expectations {
		var (
			input          = exp.input
			expectedOutput = exp.expectedOutput
			takeFn         = exp.takeFn
		)

		t.Run(fmt.Sprintf("stream take while, test case %v", i+1), func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.TakeWhile(takeFn)

			index := 0
			for item, ok := stream.Next(); ok; item, ok = stream.Next() {
				if len(expectedOutput) > index {
					assert.Equal(expectedOutput[index], item)
					index++
				}
			}

			_, ok := stream.Next()
			assert.False(ok)
			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_chain_calls(t *testing.T) {
	var (
		assert = assert.New(t)

		input          []T
		expectedOutput []T
	)

	for i := 1; i <= 100; i++ {
		input = append(input, i)
	}

	for i := 1; i <= 99; i++ {
		v := i + 100
		if v%11 != 0 {
			continue
		}
		if v <= 121 {
			continue
		}
		expectedOutput = append(expectedOutput, v)
	}

	iterator := streamer.NewSliceIterator(input)
	stream := streamer.NewStream(iterator)

	stream = stream.
		Map(func(v T) T { return v.(int) + 100 }).
		Take(99).
		Filter(func(v T) bool { return v.(int)%11 == 0 }).
		Skip(1).
		SkipWhile(func(v T) bool { return v.(int) <= 121 })

	index := 0
	for item, ok := stream.Next(); ok; item, ok = stream.Next() {
		if len(expectedOutput) > index {
			assert.Equal(expectedOutput[index], item)
			index++
		}
	}

	assert.Equal(len(expectedOutput), index)
}
