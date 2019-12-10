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

			stream.Map(mapFn)

			index := 0
			for stream.Next() {
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

			stream.ChunkBy(chunkFn)

			index := 0
			for stream.Next() {
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_chunk_every(t *testing.T) {
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

			stream.ChunkEvery(chunkSize)

			index := 0
			for stream.Next() {
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_skip(t *testing.T) {
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

			stream.Skip(skipCount)

			index := 0
			for stream.Next() {
				assert.Equal(expectedOutput[index], stream.Value())
				index++
			}

			assert.Equal(len(expectedOutput), index)
		})
	}
}

func Test_skip_while(t *testing.T) {
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

			stream.SkipWhile(skipFn)

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
