# 第3节 - 追踪远程调用请求

## 目标

学习:

* 在多个微服务中追踪调用链路
* 通过`Inject`和`Extract`在进程间传递`Context`
* 应用 `OpenTracing` 推荐的 `tags`

## 实践
这一小节的代码将不在文章中详细写出，可以在 ./lesson03/exercise/basic/ 中查看
### Hello-World 微服务程序

我们仍然以`Hello World`程序为例，通过两个下游微服务`formatter`和`publisher`调用 `formatString`和`printHello`函数。  
代码的组织结构如下:

  * `client/hello.go` 是第二节的 `hello.go`,修改为 http 调用的方式。
  * `formatter/formatter.go` 是 http 服务端，调用方式是：`GET 'http://localhost:8090/format?helloTo=Kobe'` 返回 `"Hello, Kobe!"`
  * `publisher/publisher.go` 是另一个 http 服务端，调用方式是：`GET 'http://localhost:8091/publish?helloStr=hi%20there'` 在终端打印`"hi there"`

为了方便测试，在两个终端中分别运行 `formatter`和`publisher`

```
$ go run ./lesson03/exercise/basic/formatter/formatter.go
$ go run ./lesson03/exercise/basic/publisher/publisher.go
```

向 formatter 发送 HTTP 请求

```
$ curl 'http://localhost:8081/format?helloTo=Kobe'
Hello, Kobe!%
```

向 publisher 发送 HTTP 请求

```
$ curl 'http://localhost:8082/publish?helloStr=hi%20there'
```
注意 `publisher` 程序会在终端打印出 `"hi there"`。最后运行客户端程序：

```
$ go run ./lesson03/exercise/basic/client/hello.go Kobe
2019/08/23 17:00:00 Initializing logging reporter
2019/08/23 17:00:00 Reporting span 29a86dd553a0e6ae:66fb173d2c67736a:29a86dd553a0e6ae:1
2019/08/23 17:00:00 Reporting span 29a86dd553a0e6ae:64f2fa1421475b8:29a86dd553a0e6ae:1
2019/08/23 17:00:00 Reporting span 29a86dd553a0e6ae:29a86dd553a0e6ae:0:1
```
我们可以看到 `publisher` 在终端中输出 `"Hello, Kobe!"`

### 进程间上下文传递
因为我们只是把`hello.go`中的调用改为了 HTTP 调用，所以 `trace` 记录的信息并没有什么变化。但是现在多了两个微服务，我们同样希望它们的信息也能在
`trace`中。为了在进程间和远程调用间继续记录`trace`信息，我们需要一种可以在整个请求周期中传递`context`的方式。`OpenTracing API`接口中定义
了两个函数来完成这个工作，它们是`Inject(spanContext, format, carrier)` 和 `Extract(format, carrier)`.  

`format`参数是指`OpenTracing API`定义的三种标准编码之一：  
  * `TextMap`， `spanContext` 被编码成字符串键值对的集合
  * `Binary`, `spanContext` 被编码成字节数组
  * `HTTPHeaders`， 类似于 `TextMap`， 但要求键值对可以安全的当做 HTTP headers 使用

`carrier` 参数是底层 RPC 框架的抽象。例如，如果 `format`参数为 `TextMap`，那么`carrier` 就是一个可以通过`Set(key, value)`函数设置键值对的接口(go interface);
如果 `format` 是 `Binary`， 那么 `carrier`就是一个 `io.Writer`接口。

`trace` 系统使用 `Inject` and `Extract` 函数在 RPC 调用间传递 `spanContext`

### 追踪客户端

在 `formatString` 函数中我们已经创建了一个子`span`， 为了通过 http 请求传递上下文，我们需要这样做：

#### 引入一个包
./lesson03/exercise/advance/client/hello.go
```go
import (
    "github.com/opentracing/opentracing-go/ext"
)
```

#### 调用 `Inject` 函数
./lesson03/exercise/advance/client/hello.go
```go
ext.SpanKindRPCClient.Set(span)
ext.HTTPUrl.Set(span, url)
ext.HTTPMethod.Set(span, "GET")
span.Tracer().Inject(
    span.Context(),
    opentracing.HTTPHeaders,
    opentracing.HTTPHeadersCarrier(req.Header),
)
```

在这个例子中，`carrier` 是使用 `opentracing.HTTPHeadersCarrier()`包裹的 http 请求头部对象。注意我们也给`span`添加了两个额外的标签存储了
一些关于请求的元数据，并且将 `span` 做了 `span.kind=client` 标记，这是`openTracing API` 推荐的做法[语义约定][semantic-conventions]。  
继续在 `printHello` 函数总做类似的修改

### 追踪服务端
修改的文件：  
./lesson03/exercise/advance/formatter/formatter.go
#### 引入一些包

```go
import (
    opentracing "github.com/opentracing/opentracing-go"
    "github.com/opentracing/opentracing-go/ext"
    otlog "github.com/opentracing/opentracing-go/log"
    "github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)
```

#### 创建 `Tracer` 实例，类似于 `hello.go` 中做的

```go
tracer, closer := tracing.Init("formatter")
defer closer.Close()
```

#### 使用 `tracer.Extract` 从请求中提取出 `spanContext`

```go
spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
```

#### 创建一个子`span`
使用一个特殊的参数`RPCServerOption`， 它创建一个 `ChildOf` 关系的 `spanCcontext`，并且在新的 `span` 上设置 `span.kind=server`的标签。 

```go
span := tracer.StartSpan("format", ext.RPCServerOption(spanCtx))
defer span.Finish()
```

#### 可选的, 为 `span` 添加 `tags` 和 `logs`

```go
span.LogFields(
    otlog.String("event", "string-format"),
    otlog.String("value", helloStr),
)
```

#### 类似的修改 `./lesson03/exercise/advance/publisher/publisher.go`

### 运行一把

```
# client
$ go run lesson03/exercise/advance/client/hello.go Kobe
2019/08/23 20:14:14 Initializing logging reporter
2019/08/23 20:14:14 Reporting span 2bc39845f8bdc13a:1e620ad894cafc7b:2bc39845f8bdc13a:1
2019/08/23 20:14:14 Reporting span 2bc39845f8bdc13a:1fe33394aa12c4f1:2bc39845f8bdc13a:1
2019/08/23 20:14:14 Reporting span 2bc39845f8bdc13a:2bc39845f8bdc13a:0:1


# formatter
$ go run lesson03/exercise/advance/formatter/formatter.go
2019/08/23 20:13:44 Initializing logging reporter
2019/08/23 20:14:14 Reporting span 2bc39845f8bdc13a:26b8004ebf209ccb:1e620ad894cafc7b:1

# publisher
$ go run lesson03/exercise/advance/publisher/publisher.go 
2019/08/23 20:13:50 Initializing logging reporter
Hello, Kobe!
2019/08/23 20:14:14 Reporting span 2bc39845f8bdc13a:4968561c94b49941:1fe33394aa12c4f1:1
```
注意所有的 `span` 拥有同一个 `trace ID： 2bc39845f8bdc13a`，这表示这是一次正确的追踪. 当有错误发生时，这也是一个非常有用的调试的途径。
一个典型的错误是爱传递的过程中，在一些地方丢失了 context，这样会产生多个 `traceID`，扰乱追踪。  

最后，打开 `Jaeger UI`看一下刚刚生成的这条跨服务的追踪链路。

## 结语

前半部分的练习可以在 [basic](./exercise/basic) 中找到，后半部分的可以在 [advance](./exercise/advance) 中找到。 
最终的代码可以在 [solution](./solution) 下找到

下一节: [Baggage(行李)](../lesson04).

[semantic-conventions]: https://github.com/opentracing/specification/blob/master/semantic_conventions.md
