# E2E

This folder contains set of e2e tests. They are not intended to be used by end-user as examples. For that we have [examples folder](../examples/).

Every e2e-test is a separate Nevalang module. That means you can break one test, but all other will still compile.

Every test contains nevalang code as well as go file with actual test. We will continue write tests in Go until Nevalang isn't mature enough.