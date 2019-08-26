# openTracing - Strong protection of your hair

## Why?  

### 引入场景
场景：  
老板 `A总` 安排 `小B` 开发游戏币下单业务，按照接口设计 `小B`只需要对前端请求做验证、转换，然后调用一个 `小C` 提供的下单API即可，无需关心后续具体的流程。但是在调试测试过程中，
`小C` 的下单接口频频返回`internal error 500`, `小C` 不得不查看日志，然后说是调用在 `小D` 写的一个服务的时候发生了 500 错误。。。就这样，
一下午过去了，`小B` 联系了 24 个人，才发现是 `小Z` 的写的一个 slice 切片操作没有进行长度判断，返回 500.    

一天，玩家 `Pony` 血槽马上要空了，急需购入药箱，但游戏频频提示`药箱正在路上，请您耐心等待`。  
这晚无数个像`Pony`一样的玩家，打爆了`A总`的电话，而`A总`不得不紧急拨打25个电话，叫醒`小B-Z`，后面的事情可想而知，`小B-Z`失业了。   

再后来，`Pony` 立志让天下没有难买的药箱，于是他收购了 `A总`的公司，召回`小B-Z`，发现之前的错误是因为`小X`忘记做余额不足的判断，而小Y将账户
余额的类型设置为了 uint64， 引发了 500 错误。  

`Pony` 调研后，决定引入`OpenTrace`和监控系统，当类似的问题发生时，系统自动发现 500 错误是`小Z` 服务引起的，并自动拨通了 `小Z` 的电话。  

游戏重新上线一段时间后，`Pony`喜忧参半，喜的是系统稳定，买药箱的玩家越来越多。但是系统的响应速度却越来越慢，想优化却无从入手。虽说日志记录的非常详细，
但确是分散的，有一优秀的服务使用了 RequestID 对请求进行标记，只要在各个服务的日志中，找到 RequestID 对应的日志，最后将他们整合到一起，如果日志中恰好
记录有时间，就可以进行运行时长问题的研究了。但，RequestID 规范很弱，随着服务和机器的增多，可操作行越来越差。  

其实，这个问题利用 OpenTracing 就可以解决，详细解决方案后面介绍。
### 总结
系统微服务化后多应用，多实例需要面临的一些问题：
- 错误原因快速定位
- 用户体验优化(响应时长)
- 架构(调用链路)优化


## What?

- 平台无关
- 厂商无关
- 方便的添加（或更换）追踪系统的实现
- openTracing 组织提供了大量的辅助类库(https://github.com/opentracing https://github.com/opentracing-contrib https://opentracing.io/registry/?s=go)

### openTracing API

#### Traces (跟踪，追踪)

Dependencies Directed Acyclic Graph  

```
todo： 这里需要画图说明
大家可能听过有人听过一个名词：分布式链路追踪，但如果咬文嚼字的话，其实“链(link)”这个字并不贴切，链的形状是 A->B->C，显然我们的写的代码并不都是这样的，
有的时候是树状的，但其实它是一个图状的，标准的说，是一个有向无环图。
```

##### 流程图 or 时序图？
流程图：
![一次调用](https://wu-sheng.gitbooks.io/opentracing-io/content/images/OTOV_1.png)  

时序图：
![Trace过程](https://wu-sheng.gitbooks.io/opentracing-io/content/images/OTOV_3.png)  

总结： 
- 流程图易于看组件间的调用关系，但不方便看调用时间，有局限性
- 时序图方便看执行时间，方便表示串行、并行关系。
- Trace 一般使用时序图展示。

#### Spans
一个 `span` 代表一个逻辑运行单元。
##### span 三要素

- 操作名称 （`get` or `get_user/999` or `get_user`)
- 开始时间
- 结束时间

##### span relationships 
###### ChildOf 
父 span 需要 子 span 的返回值

###### FollowsFrom 

父 span 

#### Logs


#### Tags


### SpanContext

#### Inject and Extract

#### Baggage

#### Baggage vs. Span Tags



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