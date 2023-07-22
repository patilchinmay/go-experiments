# Testing Time

The purpose of this module is to mock time.

## Context

Suppose, we are using the `time` package. One of our function depends on the `time.After(1* time.Minute)`.

If we test this function in it current form, our test case may need to wait 1 minute.

A better option is to substitute the usage of `time` package with another package which allows us the ability to either use real time or mock time depending on the initialization.

In this example we are using the package `github.com/benbjohnson/clock`.
