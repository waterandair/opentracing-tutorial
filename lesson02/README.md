# 第2节 - 上下文和追踪函数

## 目标

学习:

* 追踪单独的函数
* 结合多个 `span` 到一个 `trace` 中
* 传递进程内的上下文

## 练习
首选，将[第1节](../lesson01) 目录下的 `exercise/hello.go` 复制到 `lesson02/exercise/hello.go`

### 追踪独立的函数
在 [第1节](../lesson01) 中我们编写了一个只包含一个 `span` 的追踪记录的`Hello World` 程序。 那个程序中结合了两个操作，格式化和输出打印，
现在我们将这两步操作分别封装到单独的函数：

```go
span := tracer.StartSpan("say-hello")
span.SetTag("hello-to", helloTo)
defer span.Finish()

helloStr := formatString(span, helloTo)
printHello(span, helloStr)
```

添加函数:

```go
func formatString(span opentracing.Span, helloTo string) string {
    helloStr := fmt.Sprintf("Hello, %s!", helloTo)
    span.LogFields(
        log.String("event", "string-format"),
        log.String("value", helloStr),
    )

    return helloStr
}

func printHello(span opentracing.Span, helloStr string) {
    println(helloStr)
    span.LogKV("event", "println")
}
```
当然，都这里不会对执行结果产生任何变化，我们最重要的是要将每个函数包裹到自己的 `span` 中。

```go
func formatString(rootSpan opentracing.Span, helloTo string) string {
    span := rootSpan.Tracer().StartSpan("formatString")
    defer span.Finish()

    helloStr := fmt.Sprintf("Hello, %s!", helloTo)
    span.LogFields(
        log.String("event", "string-format"),
        log.String("value", helloStr),
    )

    return helloStr
}

func printHello(rootSpan opentracing.Span, helloStr string) {
    span := rootSpan.Tracer().StartSpan("printHello")
    defer span.Finish()

    println(helloStr)
    span.LogKV("event", "println")
}
```

现在运行它：

```
$ go run ./lesson02/exercise/hello.go Kobe
2019/08/23 11:35:04 Initializing logging reporter
2019/08/23 11:35:04 Reporting span 4e5bb2cb98180833:4e5bb2cb98180833:0:1
Hello, Kobe!
2019/08/23 11:35:04 Reporting span 71375987444d387a:71375987444d387a:0:1
2019/08/23 11:35:04 Reporting span 2341eb5b187ea34f:2341eb5b187ea34f:0:1
```

我们得到了三个 `span`， 但是这里有一个问题。输出的第一个十六进制的标记表示 `Jaeger Trace ID`， 但三个`span`的`Jaeger Trace ID`
都不相同。如果我们在 `Jaeger UI` 中搜索这些 `Trace ID`， 会得到三个仅包含一个`span`的 `trace`，这并不是我们想要的。  

我们希望在两个新的`span`和`main()`函数的 `rootSpan` 之间建立因果关系。我们可以通过为 `StartSpan()` 传入一个额外的参数来完成关系的建立。

```go
    span := rootSpan.Tracer().StartSpan(
        "formatString",
        opentracing.ChildOf(rootSpan.Context()),
    )
```
我们可以认为 `trace` 是一个有向无环图，图中的节点就是`span`，节点之间的边就是`span` 之间的关系。大体上可以抽象出两种关系：`Childof` 和 
`FollowsFrom`。  
`ChildOf` 用于在类似于 `sapn` 和 `rootSpan`之间建立一种关系。在`API`定义中，`span` 间的关系用 `SpanReference` 表示，
它包含一个`SpanContext`类型和`ReferenceType`类型的字段，`SpanContext`表示一个`span`中不可变并且线程安全的部分，它可以用来建立引用或者
在整个 `trace`周期中传播。`ReferenceType`用于描述`span`间的关系。`ChildOf`关系表示`rootSpan`的逻辑依赖子`span`，`rootSpan`操作完成
前需要先完成`span`。  
`OpenTracing API`中另一个标准的`ReferenceType`是`FollowsFrom`，它意味着`rootSpan`是有向无环图的祖先，但是它的执行完成并不依赖于子`span`
的完成，例如：假设子`span`表示一个消息的缓存写入。  

如果相应的修改 `printHello` 函数并运行程序，可以看到所有的`span` 都在同一个`trace`中。

```
$ go run ./lesson02/exercise/hello.go Kobe
2019/08/23 14:02:18 Initializing logging reporter
2019/08/23 14:02:18 Reporting span 88afe4ac5309af3:2ad73b5153a0a170:88afe4ac5309af3:1
Hello, Kobe!
2019/08/23 14:02:18 Reporting span 88afe4ac5309af3:5c9e2ddde2703e43:88afe4ac5309af3:1
2019/08/23 14:02:18 Reporting span 88afe4ac5309af3:88afe4ac5309af3:0:1

```
我们也可以看到与最后一条第三个位置 `0` 不同，前两个显示的是`88afe4ac5309af3`， 这个代表的是`rootSpan`的ID。`rootSpan`在最后打印出来是因为`rootSpan`
是最后执行完的。  

现在在`UI`中找到刚刚运行的 `trace`， 可以看出 `spans` 之间的依赖关系。

### 传递进程内的上下文

你可能已经意识到刚刚的改变引起一个令人不爽的烦恼。我们不得不在每个函数的第一个参数传入`span`对象。GO 语言不支持线程局部变量的概念，所以为了把
独立的`span`连接在一起，我们需要传递一些东西。我们不希望传递的是`span`对象，因为传递`span`对象会过多的侵入业务代码。  

还好 Go 标准库中有一个专门用于在应用中传递请求上下文的类型——`context.Context`,它除了处理超时和取消之类操作，还支持存储任意的 `key-value`对，
所以我们可以用它来存储当前有效的`span`. `OpenTracing API`整合了`context.Context`,并提供了方便的帮助函数。  

首先我们需要在`main()`函数中创建`context`并将`span`存储在里面。

```go
ctx := context.Background()
ctx = opentracing.ContextWithSpan(ctx, span)
```

然后将 `ctx` 作为第一个参数传递给另外两个函数

```go
helloStr := formatString(ctx, helloTo)
printHello(ctx, helloStr)
```

同时修改两个函数，使用函数 `StartSpanFromContext()` 取出 `span`

```go
func formatString(ctx context.Context, helloTo string) string {
    span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
    defer span.Finish()
    ...

func printHello(ctx context.Context, helloStr string) {
    span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
    defer span.Finish()
    ...
```

注意在这里我们忽略了第二个返回值，它是另外一个存储了新的 `span` 的`Context`，如果函数中需要调用更多的函数，我们应该保留第二个返回值，并在之后的
传递中使用它，而不是上层传入的 `context` 

最后，因为 `StartSpanFromContext` 函数使用 opentracing.GlobalTracer()` 去创建一个新的 `span`， 所以我们需要在`main()`函数中将它初始化为
`Jaeger Trace`

```go
tracer, closer := tracing.Init("hello-world")
defer closer.Close()
opentracing.SetGlobalTracer(tracer)
```

运行代码，看到和之前的效果一样：
```
$ go run ./lesson02/exercise/hello.go Kobe
2019/08/23 14:46:30 Initializing logging reporter
2019/08/23 14:46:30 Reporting span 728b00b596d777fc:7a2c1a8baa1662ff:728b00b596d777fc:1
Hello, Kobe!
2019/08/23 14:46:30 Reporting span 728b00b596d777fc:6639cc89c8602228:728b00b596d777fc:1
2019/08/23 14:46:30 Reporting span 728b00b596d777fc:728b00b596d777fc:0:1
```

注意， 通过帮助函数`opentracing.StartSpanFromContext()`创建的`span`是 `childOf`关系的，如果想创建`FollowsFrom`关系的`span`，可以
这样做：
```go
rootSpan := opentracing.SpanFromContext(ctx)

span := rootSpan.Tracer().StartSpan(
    "example",
    opentracing.FollowsFrom(rootSpan.Context()),
)
```

## 结语

完整的项目可以在 [solution](./solution) 中找到。

下一节: [追踪RPC请求](../lesson03).
