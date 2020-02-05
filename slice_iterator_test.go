package streamer_test

import (
	"testing"

	"github.com/dc0d/streamer"

	assert "github.com/stretchr/testify/require"
)

func Test_slice_iterator(t *testing.T) {
	t.Run("iterates over a slice of objects", func(t *testing.T) {
		var (
			assert = assert.New(t)

			input          = []interface{}{1, 2, "3"}
			expectedResult = []interface{}{1, 2, "3"}

			iterator streamer.Iterator = streamer.NewSliceIterator(input)
		)

		index := 0
		for item, ok := iterator.Next(); ok; item, ok = iterator.Next() {
			assert.Equal(expectedResult[index], item)
			index++
		}

		assert.Equal(len(expectedResult), index)
	})

	t.Run("iterates over a nil slice", func(t *testing.T) {
		var (
			assert = assert.New(t)

			input = []interface{}(nil)

			iterator streamer.Iterator = streamer.NewSliceIterator(input)
		)

		index := 0
		for _, ok := iterator.Next(); ok; _, ok = iterator.Next() {
			index++
		}

		assert.Equal(0, index)
	})
}
