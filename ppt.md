## openTracing - Strong protection of your hair

### Why?
系统微服务化后需要面临的一些问题：
- 错误原因分析
- 用户体验优化(响应时长)
- 架构(调用链路)优化


### What?

- 平台无关
- 厂商无关
- 方便的添加（或更换）追踪系统的实现
- openTracing 组织提供了大量的辅助类库(https://github.com/opentracing https://github.com/opentracing-contrib https://opentracing.io/registry/?s=go)

#### openTracing API

##### Trace (跟踪，追踪)

Dependencies Directed Acyclic Graph  

```
todo： 这里需要画图说明
大家可能听过有人听过一个名词：分布式链路追踪，但如果咬文嚼字的话，其实“链(link)”这个字并不贴切，链的形状是 A->B->C，显然我们的写的代码并不都是这样的，
有的时候是树状的，但其实它是一个图状的，标准的说，是一个有向无环图。
```



### How?

### OpenTracing In Action

#### basic hello-world project 

#### nicely hello-world project

#### distributed hello-world project

- OpenTracing推荐在RPC的客户端和服务端，至少各有一个span,用于记录RPC调用的客户端和服务端信息。 

### OpenTracing In MicroPay

#### dependencies DAG

#### show code

#### show trace in Jaeger/Zipkin