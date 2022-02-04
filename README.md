# Ordered Select

This implements [deterministic select](https://www.sethvargo.com/what-id-like-to-see-in-go-2/#deterministic-select) for go1.18+.
Two to 9 channels, plus a default case are suppported. It utilizes the AST for code generation & 1.18 generics to allow reading from
arbitrary channel types.

[![Go Reference](https://pkg.go.dev/badge/jonwillia.ms/oselect.svg)](https://pkg.go.dev/jonwillia.ms/oselect)