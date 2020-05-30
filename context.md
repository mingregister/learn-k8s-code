https://www.jianshu.com/p/f5919c78b5cd
https://studygolang.com/pkgdoc

**管理goroutine** 

在RPC或者Web服务中，当Server端接受一个request的时候，都会开启一个额外的gorountine来处理内容。有时候可能会使用多个gorountine来处理同一个request。如果对应的request被取消或者timeout，就需要所有为这个request服务的gorountine被快速回收。

Incoming requests to a server should create a Context, and outgoing calls to servers should accept a Context. The chain of function calls between them must propagate the Context, optionally replacing it with a derived Context created using WithCancel, WithDeadline, WithTimeout, or WithValue. When a Context is canceled, all Contexts derived from it are also canceled.


The same Context may be passed to functions running in different goroutines; Contexts are safe for simultaneous use by multiple goroutines.




* func Background
```go
func Background() Context
```
Background returns a non-nil, empty Context. It is never canceled, has no values, and has no deadline. It is typically used by the main function, initialization, and tests, and **as the top-level Context for incoming requests.**

* func TODO
```go
func TODO() Context
```

TODO returns a non-nil, empty Context. Code should use context.TODO when it's unclear which Context to use or it is not yet available (because the surrounding function has not yet been extended to accept a Context parameter). TODO is recognized by static analysis tools that determine whether Contexts are propagated correctly in a program.


### Snippets
context的通常开头???
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```