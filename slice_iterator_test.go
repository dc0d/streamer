package streamer_test

import (
	"fmt"
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

				iterator streamer.Iterator = streamer.NewSliceIterator(input)
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

				iterator streamer.Iterator = streamer.NewSliceIterator(input)
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
