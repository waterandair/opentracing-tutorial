# 第一节 - Hello World

## 目标

学习如何:

* 实例化一个 Tracer
* 创建一个简单的 trace
* 对 trace 进行注释

## 练习

### 一个简单的 Hello-World 程序

在 `lesson01/hello.go` 中编写一个简单的程序， 接收一个参数 arg，打印出 "Hello, {arg}!"

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) != 2 {
        panic("ERROR: Expecting one argument")
    }
    helloTo := os.Args[1]
    helloStr := fmt.Sprintf("Hello, %s!", helloTo)
    println(helloStr)
}
```

运行它:
```
$ go run ./lesson01/hello.go Kobe
Hello, Kobe!
```

### 创建一个 Trace
一个 `Trace` 是一系列 `Span` 有向无环图。一个 `span` 代表系统中具有开始时间和执行时长的逻辑运行单元。`span` 之间通过嵌套或者顺序排列建立逻辑因果关系。 
每个 `span` 都有三个基本属性： 操作名称，开始时间和结束时间。  

创建一个仅由一个 `span` 组成的 `trace`，首先需要创建一个 `opentracing.Tracer` 实例，可以用 `opentracing.GlobalTracer()`创建。

```go
tracer := opentracing.GlobalTracer()

span := tracer.StartSpan("say-hello")
println(helloStr)
span.Finish()
```

我们正在使用 OpenTracing API 的几个特性:
  * 一个 `tracer` 实例通过调用`StartSpan`函数开始一个新的 `span`
  * 每个 `span`， 都要赋予一个操作名称，在这个例子中是 `"say-hello"`
  * 每个 `span` 都必须通过调用它的 `Finish()` 函数去表示结束
  * `span` 的开始结束时间戳都会被 `tracer`自动捕获  
  
但是，再次运行程序不会看到任何变化，`tracing UI` 中也没有记录。这是因为 `opentracing.GlobalTracer()` 函数返回的是一个默认的无操作的
`tracer`

### 初始化一个真正的 tracer
初始化一个真正的 `Tracer` 实例，这里选择用 `Jaeger` (http://github.com/uber/jaeger-client-go).

```go
import (
	"fmt"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func initJaeger(service string) (opentracing.Tracer, io.Closer) {
    cfg := &config.Configuration{
        Sampler: &config.SamplerConfig{
            Type:  "const",
            Param: 1,
        },
        Reporter: &config.ReporterConfig{
            LogSpans: true,
        },
    }
    tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
    if err != nil {
        panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
    }
    return tracer, closer
}
```

为了可以使用这个 `Tracer` 实例， 对 main 函数进行修改：

```go
tracer, closer := initJaeger("hello-world")
defer closer.Close()
```

记住我们初始化时传入了字符串 `hello-world` 作为服务名称， 它对 `tracer`下的所有`span`进行标记，表明它们都来源于 `hello-world` 服务。  

现在运行程序，可以看到 `span` 的日志：

```
$ go run ./lesson01/hello.go Kobe
2019/08/22 16:33:04 Initializing logging reporter
Hello, Kobe!
2019/08/22 16:33:04 Reporting span 2aff405771d47806:2aff405771d47806:0:1

```
如果现在 Jaeger 程序在运行，你可以在 `Traceing UI` 中看到刚刚程序运行的链路追踪记录

### 使用 Tags 和 Logs 对 Trace 进行注释说明
现在我们创建的 trace 是最基础的。如果我们调用 `hello.go Jordan` 代替 `hello.go Kobe`, 产生的 Trace 几乎完全相同。
如果我们可以捕获程序的输入参数就更好了。

一个非常愚蠢的做法是把操作名称分别设置为 `"Hello, Kobe!"` 和 `"Hello, Jordan!"`， 千万不要这么做，操作名称应该是字符串常量。
因为首先操作名称在 Jaeger 的 UI中可以通过下拉列表选择，如果把操作名称设置为变量，那下拉列表中会出现数不清的项。其次，Trace 支持各种
聚合操作，如果操作名称是不固定的，那么聚合将没有任何意义。

推荐的方式是通过`tags`和`logs`去注释说明`span`，`tag` 是为`span`提供某些元数据的键值对，`log`类似于通常的日志，它包含一些与`span`相关的时间戳和
其他信息。

什么情况下我们应该使用`tags`与`logs`？ `tags`用于描述适用于一个`span`整个运行单元内的一些属性。例如，假设一个`span`代表一个`HTTP`请求，
那么请求的`URL`应该被记录在`tag`中，另外，如果服务返回一个跳转链接，应该把它记到日志中，因为围绕这个事件有一个明确的时间戳。OpenTracing 
规范提供了一份关于`tags`和`logs`字段的最佳时间，称为 [语义约定][semantic-conventions]

#### 使用 Tags

在例子`hello.go Kobe`中，字符串`Kobe`是很适合作为`span tag`，因为它存在与整个时段，而不是某个特定的时刻，我们可以这样记录它：

```go
span := tracer.StartSpan("say-hello")
span.SetTag("hello-to", helloTo)
```

#### 使用 Logs

hello world 程序过于简单，以至于很难找到设置记录到日志的信息，但是为了便于讲解案例，就强行记录一下日志。程序中对 `helloStr` 进行了格式化，并
将其打印出来，这两个操作都需要确定的时间，所以可以在日志中记录它们的完成时间：
```go
import "github.com/opentracing/opentracing-go/log"

...

helloStr := fmt.Sprintf("Hello, %s!", helloTo)
span.LogFields(
    log.String("event", "string-format"),
    log.String("value", helloStr),
)

println(helloStr)
span.LogKV("event", "println")
```

如果你之前没有使用过结构化的日志，这里的日志声明对你来说可能会有点奇怪。不同于将日志复制给一个单独的字符串，结构化的日志接口鼓励你将日志信息
分割为一些key-value键值对，这可以方便日志聚合系统自动处理。这样的做法是因为如今处理日志更多的是通过机器，而不是人用眼看。关于结构化日志的
更多讨论，可以参考这里[google "structured-logging"][google-logging]

GO 语言的`OpenTracing API`暴露出两个用于结构化日志的接口：  
  * `LogFields` 函数采用强类型的键值对，它转为零分配设计
  * `LogKV` 函数交替的传入`key-value`键值对，比如 `key1,value1,key2,value2` ， 方便使用.

`OpenTracing` 规范建议所有的日志声明都包含`event`字段，用于描述记录整个事件，事件提供的其他属性可以作为额外的字段记录。

做了上面的修改后再运行程序，就可以在 UI 中找到 `trace`， 展开它的 `span` 可以看到 `tags` 和 `logs`

## 结语
完整的示例代码可以在[solution](solution)中找到，将`initJaeger`函数作为帮助函数移动到，这样方便在其他示例中复用。  

下一节：[上下文和追踪函数](../lesson02).

[semantic-conventions]: https://github.com/opentracing/specification/blob/master/semantic_conventions.md
[google-logging]: https://www.google.com/search?q=structured-logging
