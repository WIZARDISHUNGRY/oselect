# Ordered Select

This implements [deterministic `select`](https://www.sethvargo.com/what-id-like-to-see-in-go-2/#deterministic-select) for go1.18+.
Two to 9 channels, plus a default case are suppported. It utilizes the AST for code generation & 1.18 generics to allow reading from
arbitrary channel types.

[![Go Reference](https://pkg.go.dev/badge/jonwillia.ms/oselect.svg)](https://pkg.go.dev/jonwillia.ms/oselect)

# FAQ

1. Why can't you mix and match sends and receives?

    While we could generate functions for every permutation of send and receives, the function naming would be even
    more unwieldy. e.g. `Send_Send_Send_RecvOK_Default`. It could be possible to make each of the args a struct:

```go
type Param[T any] struct {
    RecvChan <-chan T
    RecvFunc func(T, bool)
    SendChan chan<- T
    SendFunc func() T
}
```

    and rely on the fact that `select`ing on the nil channels will never unblock. The calling syntax for this would be even more
    labourous, even with convenience helpers:

```go
oselect.Select4(
    oselect.Recv(ctx.Done(), doneCallback),
    oselect.Recv(uiMessages, dispatchEvent),
    oselect.Recv(ircMessages, dispatchIRC),
    oselect.Recv(twitterMessages, dispatchTwitter),
    oselect.Send(metricsChan, getMetrics),
)
```
    I don't hate this as much as I originally expected. Performance implications TBD!

2. Would [variadic templates](https://www.ibm.com/docs/en/zos/2.1.0?topic=only-variadic-templates-c11)
remove the need for generating a function for every N-terms?

    No, because there's no way to generate a `select` block for an arbitrary number of of channels at compile time.