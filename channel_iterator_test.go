package streamer_test

import (
	"testing"
	"time"

	"github.com/dc0d/streamer"

	assert "github.com/stretchr/testify/require"
)

func Test_channel_iterator_blocking(t *testing.T) {
	t.Run("channel iterator blocking", func(t *testing.T) {
		var (
			assert = assert.New(t)

			input          <-chan interface{}
			expectedOutput = []interface{}{1, 2, "3"}

			iterator streamer.Iterator
		)

		{
			ch := make(chan interface{})
			go func() {
				defer close(ch)
				for _, v := range []interface{}{1, 2, "3"} {
					ch <- v
				}
			}()

			input = ch
			iterator = streamer.NewChannelIterator(input, -1)
		}

		index := 0
		for item, ok := iterator.Next(); ok; item, ok = iterator.Next() {
			assert.Equal(expectedOutput[index], item)
			index++
		}

		assert.Equal(len(expectedOutput), index)
	})
}

func Test_channel_iterator_timeout(t *testing.T) {
	t.Run("channel iterator timeout", func(t *testing.T) {
		var (
			assert = assert.New(t)

			input          <-chan interface{}
			expectedOutput = []interface{}{1}

			iterator streamer.Iterator
		)

		{
			ch := make(chan interface{}, 1)
			ch <- 1

			input = ch
			iterator = streamer.NewChannelIterator(input, time.Millisecond*50)
		}

		index := 0
		for item, ok := iterator.Next(); ok; item, ok = iterator.Next() {
			assert.Equal(expectedOutput[index], item)
			index++
		}

		assert.Equal(len(expectedOutput), index)
	})

	t.Run("channel iterator with timeout, closing channel", func(t *testing.T) {
		var (
			assert = assert.New(t)

			input          <-chan interface{}
			expectedOutput = []interface{}{1}

			iterator streamer.Iterator
		)

		{
			ch := make(chan interface{})
			go func() {
				defer close(ch)
				ch <- 1
			}()

			input = ch
			iterator = streamer.NewChannelIterator(input, time.Millisecond*50)
		}

		index := 0
		for item, ok := iterator.Next(); ok; item, ok = iterator.Next() {
			assert.Equal(expectedOutput[index], item)
			index++
		}

		assert.Equal(len(expectedOutput), index)
	})
}

func Test_channel_iterator_nonblocking(t *testing.T) {
	t.Run("channel iterator nonblocking", func(t *testing.T) {
		var (
			assert = assert.New(t)

			input <-chan interface{}

			iterator streamer.Iterator
		)

		{
			ch := make(chan interface{})

			input = ch
			iterator = streamer.NewChannelIterator(input, 0)
		}

		index := 0
		for _, ok := iterator.Next(); ok; _, ok = iterator.Next() {
			index++
		}

		assert.Equal(0, index)
	})

	t.Run("channel iterator nonblocking with value", func(t *testing.T) {
		var (
			assert = assert.New(t)

			input          <-chan interface{}
			expectedOutput = []interface{}{1}

			iterator streamer.Iterator
		)

		{
			ch := make(chan interface{}, 1)
			ch <- 1

			input = ch
			iterator = streamer.NewChannelIterator(input, 0)
		}

		index := 0
		for item, ok := iterator.Next(); ok; item, ok = iterator.Next() {
			assert.Equal(expectedOutput[index], item)
			index++
		}

		assert.Equal(len(expectedOutput), index)
	})

	t.Run("channel iterator nonblocking, closing channel", func(t *testing.T) {
		var (
			assert = assert.New(t)

			input          <-chan interface{}
			expectedOutput = []interface{}{1}

			iterator streamer.Iterator
		)

		{
			ch := make(chan interface{}, 1)
			ch <- 1
			close(ch)

			input = ch
			iterator = streamer.NewChannelIterator(input, 0)
		}

		index := 0
		for item, ok := iterator.Next(); ok; item, ok = iterator.Next() {
			assert.Equal(expectedOutput[index], item)
			index++
		}

		assert.Equal(len(expectedOutput), index)
	})
}
