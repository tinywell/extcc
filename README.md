# 说明

extcc 实现了一种用于 external service 模式的 chaincode 的测试仿真端，包含模拟账本节点的服务端和调用 chaincode 的客户端。

chaincode 是超级账本项目下区块链平台 fabric 的智能合约，其生命周期由账本节点 peer 管理，chaincode 的编译、部署、运行都只能通过 peer 节点运行，无法独立运行进行功能测试。在传统模式下，测试 chaincode 需要部署一个用于测试的 fabric 网络，然后通过复杂的安装过程将其部署到网络中，然后通过调用 fabric 的交易相关接口实现 chaincode 的调用测试。

后续 chaincode 支持 external service 模式运行，chaincode 可以独立启动，但是依然无法独立进行功能测试，主要是因为 chaincode 的绝大部分底层 api 需要和 peer 交互才能得到结果，比如账本的读写等，chaincode 本身不具备数据存储的能力。

本项目的目的就是减少 chaincode 的功能测试的前置依赖，提高 chaincode 的早期开发阶段的共工作效率。

# 功能

通过分析 external service 模式的 chaincode 实现原理后发现，chaincode 运行一个`ChaincodeServer`  接口的服务:

```go
// ChaincodeServer is the server API for Chaincode service.
type ChaincodeServer interface {
	Connect(Chaincode_ConnectServer) error
}
```

peer 则作为客户端通过这个服务与 chaincode server 建立连接，后续所有 chaincode 与 peer 的交互通过这个连接完成。

以此为基础，extcc 实现了 peer 账本节点仿真服务、chaincode 注册管理工具和chaincode 功能测试调用工具。

# 使用

下载并编译

```Bash
git clone https://github.com/tinywell/extcc.git
cd extcc
go build -o extcc
```

使用 `server` 子命令启动仿真服务。

```Bash
./extcc server -h
模拟 peer 功能与链码通讯，完成链码调用

Usage:
  extcc server [flags]

Flags:
  -h, --help   help for server
```

使用 `register` 子命令注册链码(提前启动 chaincode 服务)

```Bash
./extcc register -h                                          
链码注册

Usage:
  extcc register [flags]

Flags:
      --addr string     仿真服务地址 (default "localhost:9999")
      --ccaddr string   外部链码服务地址 (default "localhost:9998")
      --ccname string   外部链码名称 (default "mycc")
  -h, --help            help for register
```

使用 `invoke` 调用 chaincode ，开始功能测试

```Bash
./extcc invoke -h                                                                                        
调用链码

Usage:
  extcc invoke [flags]

Flags:
      --addr string        仿真服务地址 (default "localhost:9999")
      --args string        链码调用参数(exp:'["get","A"]') (default "[]")
      --cc string          链码名称 (default "mycc")
  -h, --help               help for invoke
  -i, --init               是否调用初始化函数
      --listen string      服务监听地址 (default ":9999")
      --mspid string       交易用户 MSPID (default "Org1MSP")
  -t, --transient string   秘密参数 (default "{}")
```

chaincode 更新之后，只需要重启chaincode 服务，然后执行注册命令，就可以继续进行功能测试了。

