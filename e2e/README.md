# E2E

This folder contains a set of e2e tests. They are not intended to be used by the end-user as examples. For that we have the [examples folder](../examples/).

Every e2e-test is a separate Nevalang module. That means you can break one test, but all the other will still compile.

Every test contains nevalang code as well as a go file with the actual test. We will continue writing tests in Go until Nevalang is mature enough.
