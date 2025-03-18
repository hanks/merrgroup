package merrgroup

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	t.Run("runs Go functions concurrently", func(t *testing.T) {
		g := &Group{}
		var counter int32

		for i := 0; i < 3; i++ {
			g.Go(func() error {
				atomic.AddInt32(&counter, 1)
				return nil
			})
		}

		err := g.Wait()
		assert.NoError(t, err)
		assert.Equal(t, int32(3), counter)
	})

	t.Run("runs TryGo functions concurrently", func(t *testing.T) {
		g := &Group{}
		var counter int32

		for i := 0; i < 3; i++ {
			ok := g.TryGo(func() error {
				atomic.AddInt32(&counter, 1)
				return nil
			})
			assert.True(t, ok)
		}

		err := g.Wait()
		assert.NoError(t, err)
		assert.Equal(t, int32(3), counter)
	})

	t.Run("handles errors and combines them in Go", func(t *testing.T) {
		g := &Group{}
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		g.Go(func() error {
			return err1
		})
		g.Go(func() error {
			return err2
		})

		err := g.Wait()
		merr, ok := err.(*multierror.Error)
		assert.True(t, ok)
		assert.Contains(t, merr.Errors, err1)
		assert.Contains(t, merr.Errors, err2)
	})

	t.Run("handles errors and combines them in TryGo", func(t *testing.T) {
		g := &Group{}
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		ok := g.TryGo(func() error {
			return err1
		})
		assert.True(t, ok)

		g.TryGo(func() error {
			return err2
		})
		assert.True(t, ok)

		err := g.Wait()
		merr, ok := err.(*multierror.Error)
		assert.True(t, ok)
		assert.Contains(t, merr.Errors, err1)
		assert.Contains(t, merr.Errors, err2)
	})

	t.Run("respects semaphore limit > 0 in Go", func(t *testing.T) {
		g := &Group{}
		g.SetLimit(2)

		var running int32
		var maxRunning int32

		for i := 0; i < 5; i++ {
			g.Go(func() error {
				current := atomic.AddInt32(&running, 1)
				if current > atomic.LoadInt32(&maxRunning) {
					atomic.StoreInt32(&maxRunning, current)
				}
				time.Sleep(10 * time.Millisecond)
				atomic.AddInt32(&running, -1)
				return nil
			})
		}

		err := g.Wait()
		assert.NoError(t, err)
		assert.LessOrEqual(t, maxRunning, int32(2))
	})

	t.Run("respects semaphore limit > 0 in TryGo", func(t *testing.T) {
		g := &Group{}
		g.SetLimit(2)

		var running int32
		var maxRunning int32

		for i := 0; i < 5; i++ {
			g.TryGo(func() error {
				current := atomic.AddInt32(&running, 1)
				if current > atomic.LoadInt32(&maxRunning) {
					atomic.StoreInt32(&maxRunning, current)
				}
				time.Sleep(10 * time.Millisecond)
				atomic.AddInt32(&running, -1)
				return nil
			})
		}

		err := g.Wait()
		assert.NoError(t, err)
		assert.LessOrEqual(t, maxRunning, int32(2))
	})

	t.Run("respects semaphore limit < 0", func(t *testing.T) {
		g := &Group{}
		g.SetLimit(-1)

		var counter int32

		for i := 0; i < 3; i++ {
			g.Go(func() error {
				atomic.AddInt32(&counter, 1)
				return nil
			})
		}

		err := g.Wait()
		assert.NoError(t, err)
		assert.Equal(t, int32(3), counter)
	})

	t.Run("calls cancel function on error", func(t *testing.T) {
		g, ctx := WithContext(context.Background())

		g.Go(func() error {
			return errors.New("test error")
		})

		err := g.Wait()
		assert.Error(t, err, ctx.Err())
	})
}
