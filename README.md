# PoW Routine-Based P2P Simulator

这是一个使用 Go 实现的 PoW 区块链仿真系统，但**不**使用离散 Tick 驱动。
每个节点运行在独立的 goroutine 中，通过真实的 TCP 连接与邻居节点通信，
模拟更贴近真实区块链网络的并发行为。

## 与 Tick 版的核心区别

| 特性 | Tick 版 (`tick/`) | Routine 版 (`routine/`) |
|------|-------------------|-------------------------|
| 驱动方式 | 全局 Tick 循环 | 节点独立 goroutine |
| 网络通信 | 中心化时间桶延迟模拟 | 节点间真实 TCP 直连 |
| 时间模型 | 离散时间步进 | 连续实时 |
| 网络拓扑 | 全连接 + 随机延迟 | Bootnode 集中分配 + TCP |
| 区块传播 | 全局广播 | gossip 泛洪中继（洪泛全网） |
| 中心组件 | Simulator 全权控制 | Bootnode 仅做发现+观测 |
| 可扩展性 | 易于控制变量做实验 | 更贴近真实 P2P 网络 |

## 项目结构

```
routine/
├── config/           # 全局参数与随机源
│   └── config.go
├── core/             # Block 和 Blockchain 基础结构
│   ├── block.go
│   └── blockchain.go
├── node/             # 节点实现（goroutine + TCP）
│   └── node.go
├── p2p/              # P2P 消息协议与 TCP 传输工具
│   ├── message.go
│   └── transport.go
├── bootnode/         # 引导节点（发现 + 观测 + 拓扑分配）
│   └── bootnode.go
├── report/           # 报告生成（可读日志 + JSON 统计）
│   └── report.go
├── simulator/        # 仿真编排器
│   └── simulator.go
├── main.go           # 命令行入口
├── main_test.go      # 参数测试
└── README.md
```

## 架构设计

### 节点生命周期

每个节点启动后：

1. 创建 TCP 监听器（自动分配端口）
2. 连接 Bootnode（127.0.0.1:9000），注册自身地址
3. **等待 Bootnode 完成全部节点注册后，统一下发 peer 分配**
4. 根据分配列表向目标 peer 发起 TCP 连接（握手使用 `MsgHello`）
5. 启动矿工协程：循环尝试挖矿（基于 Difficulty 概率判定）
6. 挖到区块后通过 gossip 泛洪广播给所有邻居
7. 收到新区块后中继转发给所有邻居（除来源外），实现全网络传播

### 挖矿机制

- 每次尝试：生成 `[0,1)` 均匀随机数，若 `< Difficulty` 则出块成功
- 每次尝试后 `time.Sleep(500μs)` 控制速率
- **出块速度保证慢于网络传播**：本地 TCP 往返约 < 1ms，而 500μs 间隔 × 平均 1/Difficulty 次尝试 = 远大于网络传播时间

### 网络模型

#### 拓扑分配（Bootnode 集中式）

- 节点启动后先向 Bootnode 注册，**不自行选择 peer**
- Bootnode 等待全部 `NodeCount` 个节点注册完毕后，统一为每个节点随机分配邻居
- 每个节点的目标邻居数：`max(1, ceil(Density × (NodeCount - 1)))`
- 分配结果通过 `MsgPeerAssignment` 下发给各节点，节点据此建立 TCP 连接

这种设计解决了顺序注册导致的早期节点 peer 不足问题，保证每个节点都获得符合密度要求的邻居数量。

#### 区块传播（Gossip 泛洪）

- 节点挖到新区块后发送给**所有邻居**
- 节点收到**未见过**的区块后，**中继转发**给所有邻居（排除来源 peer）
- 通过 `seenBlocks` 集合去重，避免重复转发和无限循环
- 区块因此能传播到全网络所有节点，而非仅一跳范围

#### 连接握手

- 节点向 peer 发起 TCP 连接后，先发送 `MsgHello{NodeID}` 告知自身身份
- 接收方从握手消息识别对方节点 ID，正确注册到 peer 池
- 这确保了入站连接可被准确标识，支持 gossip 中继时排除来源 peer

#### 连接生命周期

- 发送区块失败时**不主动关闭连接或删除 peer**，只记录日志
- 连接清理由 `peerReadLoop` 统一管理：当读端检测到连接断开时自动清理

### 伪哈希

与 Tick 版一致，不进行真实 SHA256 计算：
- 使用 `math/rand/v2` 的 ChaCha8 流密码生成 32 字节伪哈希
- 仅用于标识父子区块关系

### 报告输出

仿真结束后默认同时生成两类结果，写入同一父目录下的不同子目录：

1. `results/log/`：面向阅读的链结构和连接拓扑报告
2. `results/json/`：面向后续分析的参数、区块 DAG、Tips、矿工统计和 peer 分配

两个目录中的时间戳结果默认不纳入 Git；示例文件为 `results/log/example.log` 和
`results/json/example.json`。

## 可调参数

所有参数定义在 `config/config.go` 中：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `NodeCount` | int | 20 | 节点数量 |
| `Difficulty` | float64 | 0.01 | 挖矿难度（每次尝试成功概率） |
| `Density` | float64 | 0.3 | 图稠密度 (0, 1] |
| `TargetHeight` | int | 50 | 仿真终止目标高度 |
| `Seed` | uint64 | 42 | 随机种子（0 = 不固定） |

## 快速开始

```bash
cd routine
go run main.go
```

运行参数扫描测试：

```bash
go test -v -run TestParameterSweep
```

## 消息类型

所有消息定义在 `p2p/message.go` 中：

| 常量 | 值 | 方向 | 说明 |
|------|------|------|------|
| `MsgRegisterPayload` | `register_payload` | Node → Bootnode | 节点注册（ID + 监听地址） |
| `MsgPeerList` | `peer_list` | Bootnode → Node | 注册确认（空列表） |
| `MsgPeerAssignment` | `peer_assignment` | Bootnode → Node | 最终 peer 分配（全量注册后统一下发） |
| `MsgHello` | `hello` | Node → Node | 连接握手（告知拨号方 ID） |
| `MsgBlock` | `block` | Node ↔ Node / Node → Bootnode | 区块传播与上报 |
| `MsgChainRequest` | `chain_req` | 预留扩展 | 向邻居请求链信息 |
| `MsgChainResponse` | `chain_resp` | 预留扩展 | 响应链信息 |

## 可扩展性

系统设计考虑了后续扩展：

1. **新消息类型**：在 `p2p/message.go` 添加消息常量，在 `node.go` 的消息分发 switch 中添加处理分支
2. **恶意节点策略**：在 `node.go` 中增加 `Strategy` 枚举和对应的挖矿/收块逻辑
3. **交易/UTXO**：在 `core/block.go` 中扩展 Block 结构体
4. **难度调整**：在 `config/config.go` 中增加动态难度逻辑
5. **P2P 发现增强**：在 `bootnode/bootnode.go` 中增加 Kademlia 等协议
6. **链同步**：在 `p2p/message.go` 中已有 `MsgChainRequest` / `MsgChainResponse` 预留

## 可复现性

设置 `Seed = 42`（或任意非零值）后，相同参数下每次运行结果一致。
设置 `Seed = 0` 则每次使用系统随机种子。
