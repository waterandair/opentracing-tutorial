# go 语言的 openTracing 实践指南(译自: https://github.com/yurishkuro/opentracing-tutorial)
`Tracing 统一翻译为链路追踪`

可以从这本书中查看一些示例 [Mastering Distributed Tracing](https://www.shkuro.com/books/2019-mastering-distributed-tracing/):
* [第4章: OpenTracing 仪表基础](https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter04)
* [第5章: 异步程序的仪表](https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter05)
* [第7章: 服务网格的链路追踪](https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter07)
* [第11章: 集成指标和日志](https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter11)
* [第12章: 通过数据挖掘收集关键信息](https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter12)


## 安装 Jaeger
这篇指南使用 `CNCF` 基金会的 `Jaeger`(https://jaegertracing.io) 作为后端链路追踪系统，我们需要通过 `Docker` 启动一个默认内存存储的
`Jaeger`系统，只暴露必须的端口，并且将`log`模式设置为 `debug`。
docker pull jaegertracing/all-in-one:1.13
```
docker run \
  --rm \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 16686:16686 \
  jaegertracing/all-in-one:1.13 \
  --log-level=debug
```
此外， 可以从 https://jaegertracing.io/download/ 下载不同平台的名为 `all-in-one` 的` `Jaeger`二进制执行程序，`Jaeger`启动后，可以
通过 http://localhost:16686 访问 `Jaeger UI`.

## 运行项目示例
本仓库使用 `go modules`管理依赖，所以要求 GO 版本大于 1.11。  
首先将代码克隆到 `$GOPATH` 路径下
```
mkdir -p $GOPATH/src/github.com/waterandair/
cd $GOPATH/src/github.com/waterandair/
git clone https://github.com/waterandair/opentracing-tutorial.git
```

然后, 安装依赖:

```
cd $GOPATH/src/github.com/waterandair/opentracing-tutorial/
make install
```

本教程的其他命令都相对于此目录运行

## Lessons

* [第一节 - Hello World](./lesson01)
  * 实例化一个 Tracer
  * Create a simple trace
  * Annotate the trace
* [第二节 - 上下文和追踪函数](./lesson02)
  * Trace individual functions
  * Combine multiple spans into a single trace
  * Propagate the in-process context
* [第三节 - 追踪 RPC 请求](./lesson03)
  * Trace a transaction across more than one microservice
  * Pass the context between processes using `Inject` and `Extract`
  * Apply OpenTracing-recommended tags
* [第四节 - Baggage(用于传输跨进程的全局数据)](./lesson04)
  * Understand distributed context propagation
  * Use baggage to pass data through the call graph
* [其他](./extracredit)
  * 使用开源仪表程序



