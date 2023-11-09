package main

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/castmetal/backoff_algo/pkg/backoff"
)

type MyBackoffTest struct {
	a         string
	b         int
	numErrors int32
}

func main() {
	start := time.Now()
	test := &MyBackoffTest{
		a: "asdd",
		b: 0,
	}

	myVar := "asdasdasd"
	b := backoff.NewBackoff(true, 3)

	fn := func() error {
		if test.numErrors >= 3 {
			fmt.Println("Executed")
			return nil
		}

		fmt.Println(test.a)
		fmt.Println(myVar)

		atomic.AddInt32(&test.numErrors, 1)

		return fmt.Errorf("test %d", test.numErrors)
	}

	ctx := context.Background()

	if err := b.ExecuteBackoff(ctx, fn); err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)

	fmt.Println(elapsed)
}
