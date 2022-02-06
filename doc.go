//go:generate go run ./cmd/generate/... select_gen.go

/*
	Package oselect allows deterministic ordering for case labels in select statements.
	For example this contrived example of trying to prioritize ui events in a loop:

		func busyLoop(ctx context.Context, uiMessages, ircMessages, twitterMessages chan <-any) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				select {
				case <-ctx.Done():
					return
				case msg := <-uiMessages:
					dispatchEvent(msg)
					continue
				default:
				}
				select {
				case <-ctx.Done():
					return
				case msg := <-uiMessages:
					dispatchEvent(msg)
					continue
				case msg := <-ircMessages:
					dispatchIRC(msg)
					continue
				default:
				}
				select {
				case <-ctx.Done():
					return
				case msg := <-uiMessages:
					dispatchEvent(msg)
				case msg := <-ircMessages:
					dispatchIRC(msg)
				case msg := <-twitterMessages:
					dispatchTwitter(msg)
				}
			}
		}

	becomes:

		func busyLoop(ctx context.Context, uiMessages, ircMessages, twitterMessages chan <-any) {
			for {
				var done bool
				Recv4(
					ctx.Done(), func() { done = true },
					uiMessages, dispatchEvent,
					ircMessages, dispatchIRC,
					twitterMessages, dispatchTwitter,
				)
				if done {
					return
				}
			}
		}

	Whether or not this is useful is left as an exercise to the reader.

*/
package oselect
