# 第4节 - Baggage（行李）

## Objectives

* 理解分布式上下文传递
* 使用 `baggage` 在整个调用周期里传递数据

### 理解概念
Baggage是存储在SpanContext中的一个键值对(SpanContext)集合。它会在一条追踪链路上的所有span内全局传输，包含这些span对应的SpanContexts。在这种情况下，"Baggage"会随着trace一同传播，他因此得名（Baggage可理解为随着trace运行过程传送的行李）。鉴于全栈OpenTracing集成的需要，Baggage通过透明化的传输任意应用程序的数据，实现强大的功能。例如：可以在最终用户的手机端添加一个Baggage元素，并通过分布式追踪系统传递到存储层，然后再通过反向构建调用栈，定位过程中消耗很大的SQL查询语句。
Baggage拥有强大功能，也会有很大的消耗。由于Baggage的全局传输，如果包含的数量量太大，或者元素太多，它将降低系统的吞吐量或增加RPC的延迟。

### 实践
在第3节中，我们俺看到如何在整个调用链中国传递 `spanContext`。不难看出，我们不仅仅可以传递上下文。在 `Opentracing` 追踪中，我们支持在分布式
上下文中同 RPC 请求一起传递一些全局的元数据，在 `Opentracing` 中这些元数据称为 baggage， 需要强调的是，它与所有 RPC 请求一起传输，顾名思义就是行李。  

可以把 [../lesson03/solution](../lesson03/solution) 下的代码复制到 `/lesson04/exercise` 下，继续编码实践.  

`formatter` 服务拿到`helloTo`参数并返回字符串`Hello， {helloTo}`，现在修改代码，支持自定义问候语。

### 在客户端设置`Baggage`

修改 `/lesson04/exercise/client/hello.go`:

```go
if len(os.Args) != 3 {
    panic("ERROR: Expecting two arguments")
}

greeting := os.Args[2]

// 在创建 span 后
span.SetBaggageItem("greeting", greeting)
```
这里我们读取命令行第二个参数，并将其存到`Baggage`中。

### 在 `formatter` 服务中读取 `Baggage`

修改 `/lesson04/exercise/formatter/formatter.go`

```go
greeting := span.BaggageItem("greeting")
if greeting == "" {
    greeting = "Hello"
}

helloTo := r.FormValue("helloTo")
helloStr := fmt.Sprintf("%s, %s!", greeting, helloTo)
```

### 运行  

运行代码，传入两个参数:

```
# client
$ go run ./lesson04/exercise/client/hello.go 科比 你好    
2019/08/24 14:24:57 Initializing logging reporter
2019/08/24 14:24:57 Reporting span 5372cc2e08329449:135abca4e4f70410:5372cc2e08329449:1
2019/08/24 14:24:57 Reporting span 5372cc2e08329449:320a284ec72ebf97:5372cc2e08329449:1
2019/08/24 14:24:57 Reporting span 5372cc2e08329449:5372cc2e08329449:0:1

# formatter
$ go run ./lesson04/exercise/formatter/formatter.go
2019/08/24 14:16:28 Initializing logging reporter
2019/08/24 14:24:57 Reporting span 5372cc2e08329449:4948a00b953ce346:135abca4e4f70410:1

# publisher
$ go run ./lesson04/exercise/publisher/publisher.go
你好, 科比!
2019/08/24 14:24:57 Reporting span 5372cc2e08329449:11df7cc3cc3ae559:320a284ec72ebf97:1
```

### What's the Big Deal（有什么大不了的）?

你可能会疑惑 - 这有什么了不起,我们简单的通过http请求参数完成这件事。但是，这正式本文的重点。我们不再需要修改API参数了，如果编写一个更复杂的程序，
拥有更深的调用树，如果用http请求参数的方式，我们将需要修改很多微服务接口。

`Baggage` 的一些可能的应用场景：
  * 在多租户系统总传递租户
  * 传递顶层调用者的身份信息
  * 为混沌工程注入的故障注入指令
  * 传递一些监控数据的请求范围维度，比如分离生产环境和测试环境的流量指标


### 警告
`baggage`是一个强大的机制，但同时也是危险的。如果你在`baggage`中存储了`1mb`的key-value数据，那你就收拾好`baggage`回家吧，
因为在调用链路的每个请求都会携带 `1MB`的数据，要很慎重的使用 `baggage`。


## 结语

完整的代码可以在 [solution](./solution) 中找到.

加餐: [使用现成的开源工具](../extracredit).
