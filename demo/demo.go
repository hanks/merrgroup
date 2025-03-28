package main

import (
	"fmt"

	"github.com/hanks/merrgroup"
	"github.com/hashicorp/go-multierror"
)

func main() {
	g := &merrgroup.Group{}

	for _, v := range []int{1, 2, 3} {
		v := v
		g.Go(func() error {
			return fmt.Errorf("error %d", v)
		})
	}

	if err := g.Wait(); err != nil {
		// Print all errors in one string
		fmt.Println(err)

		// if you want to get all the errors, you can use `Errors()`
		merr, ok := err.(*multierror.Error)
		if !ok {
			fmt.Println("Error was not a multierror.Error")
			return
		}
		// Print all errors array
		fmt.Println(merr.Errors)
	}
}

// Output:
// 3 errors occurred:
//	* error 3
//	* error 2
//	* error 1

// [error 3 error 2 error 1]
