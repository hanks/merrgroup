[![Go Tests](https://github.com/hanks/merrgroup/actions/workflows/go-test.yml/badge.svg)](https://github.com/hanks/merrgroup/actions/workflows/go-test.yml)

# merrgroup

A multi error version of [x/sync/errorgroup](https://pkg.go.dev/golang.org/x/sync/errgroup), the different is that with a `Wait` method that `errgroup` returns a single error but `merrgroup` returns all the errors information.

## Usage

The interface is the same as `errgroup`, and you can check [demo.go](./demo/demo.go) for more details.
