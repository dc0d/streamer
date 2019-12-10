package streamer_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/dc0d/streamer"

	assert "github.com/stretchr/testify/require"
)

func Test_slice_iterator(t *testing.T) {
	type (
		expectation struct {
			input          []interface{}
			expectedOutput []interface{}
		}
	)

	var (
		expectations = []expectation{
			{nil, nil},
			{[]interface{}{1}, []interface{}{1}},
			{[]interface{}{1, 2, "3"}, []interface{}{1, 2, "3"}},
		}
	)

	for i, exp := range expectations {
		var (
			input          = exp.input
			expectedOutput = exp.expectedOutput
		)

		t.Run(fmt.Sprintf("slice iterator test case %v", i+1), func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
			)

			index := 0
			for iterator.Next() {
				assert.Equal(expectedOutput[index], iterator.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}

	{
		var (
			input          = []interface{}{1, 2, "3"}
			expectedOutput = []interface{}{1, 2, "3"}
		)

		t.Run("slice iterator does not proceed on Next() until the last value is read by Value()", func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
			)

			index := 0
			for iterator.Next() {
				iterator.Next()
				assert.Equal(expectedOutput[index], iterator.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_map(t *testing.T) {
	type (
		expectation struct {
			input          []interface{}
			expectedOutput []interface{}
			mapFn          func(interface{}) interface{}
		}
	)

	var (
		expectations = []expectation{
			{nil, nil, func(x interface{}) interface{} { return x }},
			{[]interface{}{}, []interface{}{}, func(x interface{}) interface{} { return x }},
			{[]interface{}{1, 2, 3}, []interface{}{1, 2, 3}, func(x interface{}) interface{} { return x }},
			{[]interface{}{1, 2, 3}, []interface{}{2, 4, 6}, func(x interface{}) interface{} { return x.(int) * 2 }},
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

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.Map(mapFn)

			index := 0
			for stream.Next() {
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}

	{
		var (
			input          = []interface{}{1, 2, 3}
			expectedOutput = []interface{}{2, 4, 6}
			mapFn          = func(x interface{}) interface{} { return x.(int) * 2 }
		)

		t.Run("stream map does not proceed until last value is read", func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.Map(mapFn)

			index := 0
			for stream.Next() {
				stream.Next()
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_chunk_by(t *testing.T) {
	type (
		expectation struct {
			input          []interface{}
			expectedOutput []interface{}
			chunkFn        func(interface{}) interface{}
		}
	)

	var (
		expectations = []expectation{
			{
				nil,
				nil,
				func(x interface{}) interface{} { return x },
			},
			{
				[]interface{}{1, 2, 3},
				[]interface{}{
					[]interface{}{1},
					[]interface{}{2},
					[]interface{}{3},
				},
				func(x interface{}) interface{} { return x },
			},
			{
				[]interface{}{1, 3, 4, 6, 12, 7, 13, 18},
				[]interface{}{
					[]interface{}{1, 3},
					[]interface{}{4, 6, 12},
					[]interface{}{7, 13},
					[]interface{}{18},
				},
				func(x interface{}) interface{} { return x.(int) % 2 },
			},
			{
				[]interface{}{1, 2, 2, 3, 4, 4, 6, 7, 7},
				[]interface{}{
					[]interface{}{1},
					[]interface{}{2, 2},
					[]interface{}{3},
					[]interface{}{4, 4},
					[]interface{}{6},
					[]interface{}{7, 7},
				},
				func(x interface{}) interface{} { return x.(int) % 3 },
			},
			{
				[]interface{}{"str1", "123", "str2", "str3", "10", "20"},
				[]interface{}{
					[]interface{}{"str1"},
					[]interface{}{"123"},
					[]interface{}{"str2", "str3"},
					[]interface{}{"10", "20"},
				},
				func(x interface{}) interface{} {
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

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.ChunkBy(chunkFn)

			index := 0
			for stream.Next() {
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}

	{
		var (
			input          = []interface{}{1, 2, 2, 3, 4, 4, 6, 7, 7}
			expectedOutput = []interface{}{
				[]interface{}{1},
				[]interface{}{2, 2},
				[]interface{}{3},
				[]interface{}{4, 4},
				[]interface{}{6},
				[]interface{}{7, 7},
			}
			chunkFn = func(x interface{}) interface{} { return x.(int) % 3 }
		)

		t.Run("stream chunk by does not proceed until last value is read", func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.ChunkBy(chunkFn)

			index := 0
			for stream.Next() {
				stream.Next()
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_chunk_every(t *testing.T) {
	type (
		expectation struct {
			input          []interface{}
			expectedOutput []interface{}
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
				[]interface{}{1, 2, 3},
				[]interface{}{
					[]interface{}{1},
					[]interface{}{2},
					[]interface{}{3},
				},
				1,
			},
			{
				[]interface{}{1, 2, 3, 4, 5},
				[]interface{}{
					[]interface{}{1, 2},
					[]interface{}{3, 4},
					[]interface{}{5},
				},
				2,
			},
			{
				[]interface{}{1, 2, 3, 4, 5, 6, 7},
				[]interface{}{
					[]interface{}{1, 2, 3},
					[]interface{}{4, 5, 6},
					[]interface{}{7},
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

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.ChunkEvery(chunkSize)

			index := 0
			for stream.Next() {
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}

	{
		var (
			input          = []interface{}{1, 2, 3, 4, 5, 6, 7}
			expectedOutput = []interface{}{
				[]interface{}{1, 2, 3},
				[]interface{}{4, 5, 6},
				[]interface{}{7},
			}
			chunkSize = 3
		)

		t.Run("stream chunk every does not proceed until last value is read", func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.ChunkEvery(chunkSize)

			index := 0
			for stream.Next() {
				stream.Next()
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_skip(t *testing.T) {
	type (
		expectation struct {
			input          []interface{}
			expectedOutput []interface{}
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
				[]interface{}{1, 2, 3},
				[]interface{}{1, 2, 3},
				0,
			},
			{
				[]interface{}{1, 2, 3},
				[]interface{}{2, 3},
				1,
			},
			{
				[]interface{}{1, 2, 3},
				[]interface{}{3},
				2,
			},
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{3, 4, 4},
				3,
			},
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{4, 4},
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

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.Skip(skipCount)

			index := 0
			for stream.Next() {
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_skip_while(t *testing.T) {
	type (
		expectation struct {
			input          []interface{}
			expectedOutput []interface{}
			skipFn         func(interface{}) bool
		}
	)

	var (
		expectations = []expectation{
			{
				nil,
				nil,
				func(elem interface{}) bool { return false },
			},
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{1, 2, 3, 3, 4, 4},
				func(elem interface{}) bool { return false },
			},
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{},
				func(elem interface{}) bool { return true },
			},
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{3, 3, 4, 4},
				func(elem interface{}) bool { return elem.(int) < 3 },
			},
			{
				[]interface{}{2, 2, 2, 4, 4, 3, 9},
				[]interface{}{3, 9},
				func(elem interface{}) bool { return elem.(int)%2 == 0 },
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

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.SkipWhile(skipFn)

			index := 0
			for stream.Next() {
				if len(expectedOutput) > index {
					assert.Equal(expectedOutput[index], stream.Value())
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
			input          []interface{}
			expectedOutput []interface{}
			filterFn       func(interface{}) bool
		}
	)

	var (
		expectations = []expectation{
			{
				nil,
				nil,
				func(elem interface{}) bool { return true },
			},
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{1, 2, 3, 3, 4, 4},
				func(elem interface{}) bool { return true },
			},
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{},
				func(elem interface{}) bool { return false },
			},
			{
				[]interface{}{1, 3, 11, 9, 21, 17},
				[]interface{}{3, 9, 21},
				func(elem interface{}) bool { return elem.(int)%3 == 0 },
			},
			{
				[]interface{}{2, 2, 2, 4, 4, 3, 9, 8},
				[]interface{}{4, 4, 8},
				func(elem interface{}) bool { return elem.(int)%4 == 0 },
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

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.Filter(filterFn)

			index := 0
			for stream.Next() {
				if len(expectedOutput) > index {
					assert.Equal(expectedOutput[index], stream.Value())
					index++
				}
			}

			assert.Equal(len(expectedOutput), index)
		})
	}

	{
		var (
			input          = []interface{}{2, 2, 2, 4, 4, 3, 9, 8}
			expectedOutput = []interface{}{4, 4, 8}
			filterFn       = func(elem interface{}) bool { return elem.(int)%4 == 0 }
		)

		t.Run("stream filter, Next() returns true as long as Value() is not called", func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.Filter(filterFn)

			stream.Next()

			assert.True(stream.Next())

			index := 0
			for stream.Next() {
				if len(expectedOutput) > index {
					assert.Equal(expectedOutput[index], stream.Value())
					index++
				}
			}

			assert.False(stream.Next())
			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_take(t *testing.T) {
	type (
		expectation struct {
			input          []interface{}
			expectedOutput []interface{}
			takeCount      int
		}
	)

	var (
		expectations = []expectation{
			{
				[]interface{}{1, 2, 3},
				[]interface{}{1},
				1,
			},
			{
				[]interface{}{1, 2, 3},
				[]interface{}{1, 2},
				2,
			},
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{1, 2, 3},
				3,
			},
			{
				[]interface{}{},
				[]interface{}{},
				4,
			},
			{
				[]interface{}{1, 2},
				[]interface{}{1, 2},
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

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.Take(takeCount)

			index := 0
			for stream.Next() {
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}

	{
		var (
			input          = []interface{}{1, 2, 3}
			expectedOutput = []interface{}{1}
			takeCount      = 1
		)

		t.Run("stream take, Next() returns true as long as Value() is not called", func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.Take(takeCount)

			stream.Next()

			assert.True(stream.Next())

			index := 0
			for stream.Next() {
				if len(expectedOutput) > index {
					assert.Equal(expectedOutput[index], stream.Value())
					index++
				}
			}

			assert.False(stream.Next())
			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_stream_take_while(t *testing.T) {
	type (
		expectation struct {
			input          []interface{}
			expectedOutput []interface{}
			takeFn         func(interface{}) bool
		}
	)

	var (
		expectations = []expectation{
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{1, 2, 3, 3, 4, 4},
				func(elem interface{}) bool { return true },
			},
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{},
				func(elem interface{}) bool { return false },
			},
			{
				[]interface{}{1, 2, 3, 3, 4, 4},
				[]interface{}{1, 2, 3, 3},
				func(elem interface{}) bool { return elem.(int) < 4 },
			},
			{
				[]interface{}{2, 2, 2, 4, 4, 3, 9},
				[]interface{}{2, 2, 2, 4, 4},
				func(elem interface{}) bool { return elem.(int)%2 == 0 },
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

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.TakeWhile(takeFn)

			index := 0
			for stream.Next() {
				if len(expectedOutput) > index {
					assert.Equal(expectedOutput[index], stream.Value())
					index++
				}
			}

			assert.False(stream.Next())
			assert.Equal(len(expectedOutput), index)
		})
	}

	{
		var (
			input          = []interface{}{2, 2, 2, 4, 4, 3, 9}
			expectedOutput = []interface{}{2, 2, 2, 4, 4}
			takeFn         = func(elem interface{}) bool { return elem.(int)%2 == 0 }
		)

		t.Run("stream take while, Next() returns true as long as Value() is not called", func(t *testing.T) {
			var (
				assert = assert.New(t)

				iterator streamer.Iterator = streamer.NewSliceIterator(input...)
				stream                     = streamer.NewStream(iterator)
			)

			stream = stream.TakeWhile(takeFn)

			stream.Next()

			assert.True(stream.Next())

			index := 0
			for stream.Next() {
				if len(expectedOutput) > index {
					assert.Equal(expectedOutput[index], stream.Value())
					index++
				}
			}

			assert.False(stream.Next())
			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_chain_calls(t *testing.T) {
	var (
		assert = assert.New(t)

		input          []interface{}
		expectedOutput []interface{}
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

	iterator := streamer.NewSliceIterator(input...)
	stream := streamer.NewStream(iterator)

	stream = stream.
		Map(func(v interface{}) interface{} { return v.(int) + 100 }).
		Take(99).
		Filter(func(v interface{}) bool { return v.(int)%11 == 0 }).
		Skip(1).
		SkipWhile(func(v interface{}) bool { return v.(int) <= 121 })

	index := 0
	for stream.Next() {
		item := stream.Value()
		if len(expectedOutput) > index {
			assert.Equal(expectedOutput[index], item)
			index++
		}
	}

	assert.Equal(len(expectedOutput), index)
}
