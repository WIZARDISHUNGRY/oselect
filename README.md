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
    goos: linux
    goarch: amd64
    pkg: jonwillia.ms/oselect
    cpu: AMD Ryzen 7 2700X Eight-Core Processor         
    BenchmarkRecv4Default-16                   	 5048630	       231.2 ns/op
    BenchmarkSelect4Default_Recv-16            	 2219570	       540.8 ns/op
    BenchmarkSelect4Default_Recv_preroll-16    	 4139517	       283.0 ns/op
    Benchmark_select_4_default-16              	 5320692	       221.4 ns/op
    BenchmarkRecv4-16                          	 4550605	       257.0 ns/op
    BenchmarkSelect4_Recv-16                   	 2092376	       563.5 ns/op
    BenchmarkSelect4_RecvOK-16                 	 3795602	       306.6 ns/op
    BenchmarkSelect4_Recv_preroll-16           	 3774952	       316.1 ns/op
    Benchmark_select_4-16                      	 4616817	       252.6 ns/op
    ```

4. The generated code for the `Select` functions is ugly.

    That isn't a question!

5. Doesn't `select` just use [`runtime.Select`](https://pkg.go.dev/reflect#Select)'s `rselect` under the hood anyway?

    No idea! I should look into this.