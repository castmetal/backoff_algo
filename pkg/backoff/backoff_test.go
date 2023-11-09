package backoff

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBackoff(t *testing.T) {
	ctx := context.Background()

	type MyBackoffTest struct {
		s         string
		i         int
		numErrors int32
	}

	test1 := &MyBackoffTest{
		s: "My internal Str Var",
		i: 0,
	}

	myVar := "My External Context Var"

	fnBackoffWithoutError := func() error {
		if test1.numErrors >= 1 {
			return nil
		}

		atomic.AddInt32(&test1.numErrors, 1)

		return fmt.Errorf("test %d - runtimeContextVar %s", test1.numErrors, myVar)
	}

	test2 := &MyBackoffTest{
		s: "My internal Str Var",
		i: 0,
	}

	fnBackoffWithError := func() error {
		if test2.numErrors >= 1 {
			return fmt.Errorf("throw error - context %v", ctx)
		}

		atomic.AddInt32(&test2.numErrors, 1)

		return fmt.Errorf("test %d - runtimeContextVar %s", test2.numErrors, myVar)
	}

	testCases := []struct {
		desc   string
		input  BackoffCaller
		expect error
	}{
		{
			desc:   "Test Backoff function without error",
			input:  fnBackoffWithoutError,
			expect: nil,
		},
		{
			desc:   "Test Backoff function with 2 errors",
			input:  fnBackoffWithError,
			expect: fmt.Errorf("throw error - context %v", ctx),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			b := NewBackoff(true, 3)

			err := b.ExecuteBackoff(ctx, tC.input)
			require.Equal(t, tC.expect, err)
		})
	}
}
