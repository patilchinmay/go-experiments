# Gotests

[Gotests](https://github.com/cweill/gotests) makes writing Go tests easy. It's a Golang command line tool that generates table driven tests based on its target source files' function and method signatures. Any new dependencies in the test files are automatically imported.

## Description

Write a simple function that will generate a sequence of numbers. Generate and modify the test for this function using `gotests`.

## Installation

`go install github.com/cweill/gotests/gotests@v1.6.0`

`go install` command downloads the binaries into `$GOPATH` directory.
The `$GOPATH` needs to be in the `$PATH` to be able to find this binary.
It can be done by appending `export PATH=$PATH:$(go env GOPATH)/bin` to the `.zshrc` file.

Ref: [Installation](https://github.com/cweill/gotests?tab=readme-ov-file#installation)

## Generate Test

Generate test for the `generateInts()` function.

`gotests -w -i -only "generateInts" main.go`

This generates the `main_test.go` file.

## Modifications

The generated `Test_generateInts` test will have the `want` type as `<-chan int`.

Since we want to test the values coming over a channel, we will change `want` type to `[]int`.

Then we add the table test case as:
```
{
  name: "should return 0 to 9",
  want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
},
```

Then we modify the `Test_generateInts` to:
1. Invoke the `generateInts()` and get the `resultCh`
2. Populate the `got[]` by ranging over the `resultCh`.

```
resultCh := generateInts()

got := []int{}
for v := range resultCh {
  got = append(got, v)
}
```

Finally, we can use `reflect.DeepEqual` to verify the behavior.

## Run & Test

```bash
# Run the program
❯ go run .
0
1
2
3
4
5
6
7
8
9
channel closed, exiting...

# Run the tests
❯ go test ./...
ok  	github.com/patilchinmay/go-experiments/gotests	(cached)
```

## References

- [gotests](https://github.com/cweill/gotests)


