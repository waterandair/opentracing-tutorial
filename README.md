# go 语言的 openTracing 实践指南
`Tracing 统一翻译为链路追踪`

可以从这本书中查看一些示例 [Mastering Distributed Tracing](https://www.shkuro.com/books/2019-mastering-distributed-tracing/):
* [第4章: OpenTracing 仪表基础](https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter04)
* [第5章: 异步程序的仪表](https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter05)
* [第7章: 服务网格的链路追踪](https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter07)
* [第11章: 集成指标和日志](https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter11)
* [第12章: 通过数据挖掘收集关键信息](https://github.com/PacktPublishing/Mastering-Distributed-Tracing/tree/master/Chapter12)


## 前提要求
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


