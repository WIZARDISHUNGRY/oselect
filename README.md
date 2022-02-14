# Ordered Select

This implements [deterministic `select`](https://www.sethvargo.com/what-id-like-to-see-in-go-2/#deterministic-select) for go1.18+.
Two to 9 channels, plus a default case, are suppported. It utilizes the `ast` package for code generation & 1.18 generics to allow reading from
arbitrary channel types.

[![Go Reference](https://pkg.go.dev/badge/jonwillia.ms/oselect.svg)](https://pkg.go.dev/jonwillia.ms/oselect)

To get started:

`go get jonwillia.ms/oselect`


## FAQ

1. Can I mix and match sends and receives?

    Use the `Select` family of functions.

2. Would [variadic templates](https://www.ibm.com/docs/en/zos/2.1.0?topic=only-variadic-templates-c11)
remove the need for generating a function for every N-terms?

    No, because there's no way to generate a `select` block for an arbitrary number of of channels at compile time.

3. Which functions perform best?

    The `Recv`/`Send` families appear to be faster that the general purpose `Select` family.

    ```
    go test -benchmem -bench '.'                                                                 main  ✭
    goos: linux
    goarch: amd64
    pkg: jonwillia.ms/oselect
    cpu: AMD Ryzen 7 2700X Eight-Core Processor         
    BenchmarkRecv4Default-16                           	 4983958	       236.3 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSelect4Default_Recv-16                    	 2180769	       540.7 ns/op	      96 B/op	       4 allocs/op
    BenchmarkSelect4Default_Recv_preroll-16            	 4211283	       280.6 ns/op	       0 B/op	       0 allocs/op
    Benchmark_select_4_default-16                      	 5084107	       228.1 ns/op	       0 B/op	       0 allocs/op
    Benchmark_reflectDotSelect_4_default_preroll-16    	  592192	      2014 ns/op	     632 B/op	      23 allocs/op
    BenchmarkRecv4-16                                  	 4384048	       256.5 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSelect4_Recv-16                           	 1925407	       585.1 ns/op	      96 B/op	       4 allocs/op
    BenchmarkSelect4_RecvOK-16                         	 3690008	       323.5 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSelect4_Recv_preroll-16                   	 3730765	       316.2 ns/op	       0 B/op	       0 allocs/op
    Benchmark_select_4-16                              	 4705082	       255.5 ns/op	       0 B/op	       0 allocs/op
    PASS
    ok  	jonwillia.ms/oselect	14.923s
    ```

4. The generated code for the `Select` functions is ugly.

    That isn't a question!

5. Doesn't `select` just use [`selectgo`](https://go.dev/src/runtime/select.go#L121) under the hood anyway?
Why not use [`reflect.Select`](https://pkg.go.dev/reflect#Select)?

    I believe it does. It looks like the compiler uses looping constructs under the hood for non-trivial cases.

    ```go
	// The compiler rewrites selects that statically have
	// only 0 or 1 cases plus default into simpler constructs.
    ```

    Check out the benchmark for `Benchmark_reflectDotSelect_4_default_preroll`. Oof!