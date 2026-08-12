# PoW Simulator

这是一个使用 Go 实现的 PoW（Proof of Work）区块链仿真程序。项目通过“离散时间 + 随机传播延迟 + 多策略节点”的方式，模拟区块链在多节点环境中的增长过程，并研究恶意节点在分叉攻击与自私挖矿场景下的行为表现。

项目对应课程实验目标：

- 调整节点数量与出块概率，测量区块链增长速度。
- 设置恶意节点比例，统计分叉攻击成功率。
- 设置恶意节点比例，统计自私挖矿收益比例。

完整实验分析请参阅 [REPORT.md](./REPORT.md)。

## 功能概览

- 基于 Tick 的离散事件仿真。
- 支持多节点 PoW 挖矿与链头竞争。
- 支持随机网络传播延迟。
- 支持诚实节点、分叉攻击、自私挖矿三类策略。
- 支持固定随机种子复现实验结果。
- 所有实验均可通过 `go test -v ./...` 运行。

## 技术说明

### 仿真模型

系统把时间离散化为一个个 Tick。在每个 Tick 内，仿真器执行以下流程：

1. 所有节点依次尝试挖矿。
2. 新产生的区块被写入全局链状态。
3. 网络层为新区块分配随机传播延迟。
4. 当前 Tick 到期的数据包被送达对应节点。
5. 节点从收到的最高块集合中随机选择新的本地链头。

这种模型并不追求完全复刻真实公链实现，而是强调：

- 可以稳定控制变量。
- 可以复现实验结果。
- 可以清楚观察传播延迟导致的分叉现象。
- 可以方便注入恶意策略。

### 出块概率

每个节点在单个 Tick 内出块的概率为：

`Difficulty * HashPower`

其中：

- `Difficulty` 是全局难度参数。
- `HashPower` 是单节点算力，诚实节点默认为 `1`。

当恶意节点需要达到指定算力占比时，项目通过公式自动反推其算力值，而不需要修改其他节点参数。

### 网络传播模型

网络层不直接模拟真实 P2P 拓扑，而是采用“时间桶”机制：

- 一个新区块广播时，会对每个目标节点独立生成一个传播延迟。
- 该消息会被放入未来某个 Tick 的消息桶中。
- 到达对应 Tick 时，消息再被统一投递。

这个设计非常适合课程实验，因为它能用较低复杂度模拟“节点在同一时刻看到不同链状态”的现实效果。

### 链状态维护方式

项目中只有仿真器维护完整区块集合，节点自身只维护：

- 当前本地链头 `Tip`
- 当前 Tick 内接收到的最高块集合 `Inbox`

这种设计减少了每个节点重复维护完整链带来的开销，同时足以支持本实验中的主链增长、链头竞争和攻击策略。

### 攻击策略实现

#### 分叉攻击

- 恶意节点会在私有链上持续挖矿。
- 在达到目标深度前不会主动发布。
- 当私有链达到预定长度时，仿真器判断公共链是否仍未超过对应深度，从而统计攻击成功率。

#### 自私挖矿

- 恶意节点会先私藏自己挖到的区块。
- 当公共链追上其私有链前部时，恶意节点再释放私有块。
- 最终统计主链中属于恶意节点的区块比例，作为收益指标。

## 项目结构

```text
.
├── config/             # 全局参数与随机源
├── core/               # Block 和 Blockchain 等基础结构
├── exp/                # 攻击实验实现与测试入口
├── network/            # 网络延迟与数据包投递逻辑
├── node/               # 节点挖矿、收块、换链与策略状态
├── simulator/          # 仿真主循环与全局状态管理
├── main.go             # 命令行试验入口
├── main_test.go        # 参数实验：节点数/难度/出块速度
├── results/log/        # 面向阅读的实验日志
├── results/json/       # 面向分析的结构化实验数据
├── REPORT.md           # 详细实验报告
└── LICENSE             # MIT 许可证
```

## 核心模块说明

- [config.go](./config/config.go)
  - 定义 `NodesCount`、`Difficulty`、`Target`、`DelayRange` 和随机源。
  - 是所有实验参数调度的入口。

- [block.go](./core/block.go)
  - 定义区块结构。
  - 用伪哈希标识父子关系，不进行真实哈希运算。

- [blockchain.go](./core/blockchain.go)
  - 维护创世块和全局区块映射。
  - 供仿真器和实验统计统一使用。

- [node.go](./node/node.go)
  - 实现节点出块、收块、换链和攻击策略状态。
  - `Strategy` 支持 `Honest`、`Fork`、`Selfish`。

- [network.go](./network/network.go)
  - 实现基于 Tick 的传播延迟模型。
  - 每个消息对每个目标节点拥有独立随机延迟。

- [simulator.go](./simulator/simulator.go)
  - 负责驱动全局 Tick。
  - 维护全局最高链头集合 `Tips`。
  - 在主链高度达到 `Target` 时结束仿真。

- [fork_attack.go](./exp/fork_attack.go)
  - 实现分叉攻击实验。
  - 默认重复 10000 次统计成功概率。

- [selfish_mining.go](./exp/selfish_mining.go)
  - 实现自私挖矿实验。
  - 默认重复 1000 次统计主链收益比例。

## 快速开始

### 环境要求

- Go 1.22 或更高版本

### 获取依赖

```bash
go mod download
```

### 运行全部实验

```bash
go test -v ./...
```

### 只运行参数实验

```bash
go test -run TestParameters -v
```

### 只运行分叉攻击实验

```bash
go test -run TestForkATK -v ./exp
```

### 只运行自私挖矿实验

```bash
go test -run TestSelfish -v ./exp
```

## GitHub Pages 展示页

项目展示页位于 [`docs/`](./docs/)；它会同时展示 Tick 版和 Routine 版，使用两版的
`results/*/example.json` 作为静态数据快照。推送到 `main` 分支后，
`.github/workflows/pages.yml` 会自动发布到 GitHub Pages。

首次启用时，在仓库设置中将 `Settings → Pages → Build and deployment → Source`
设置为 `GitHub Actions` 即可。

## 结果说明

默认运行 `go test -v ./...` 会同时生成两类结果：

- `results/log/`：按参数实验、分叉攻击和自私挖矿实验生成的可读表格日志。
- `results/json/`：与日志对应的结构化统计数据。

两个目录中的时间戳文件默认不纳入 Git，仅保留 `example.log` 和 `example.json` 作为格式样例。

实验主要对应三类数据：

1. `TestParameters`
   - 测量节点数量与难度变化对 `ticks/block` 的影响。

2. `TestForkATKVar1/2/3`
   - 测量不同恶意比例下，目标为 6 长度分叉时的攻击成功率。

3. `TestSelfishVar1/2/3`
   - 测量不同恶意比例下，自私挖矿在最终主链中的收益占比。

从现有结果可得到几个直观结论：

- `NodesCount * Difficulty` 越大，链增长越快。
- 恶意算力比例越高，分叉攻击成功率增长越快。
- 当恶意比例较高时，自私挖矿收益可能超过其算力份额。

## 可复现性

项目已经为主要实验设置固定随机种子，因此在相同代码、相同测试顺序下可以稳定复现已有结果：

- 默认全局随机源定义在 [config.go](./config/config.go)
- 分叉攻击测试种子定义在 [fork_attack_test.go](./exp/fork_attack_test.go)
- 自私挖矿测试种子定义在 [selfish_mining_test.go](./exp/selfish_mining_test.go)

如果修改种子，实验结果数值会发生变化，但整体趋势应保持一致。

## 开发说明

### 修改实验参数

如果要修改默认参数，可以优先查看 [config.go](./config/config.go) 中的：

- `NodesCount`
- `Difficulty`
- `Target`
- `DelayRange`

如果只是做实验对比，建议优先在测试文件中临时覆盖这些参数，而不是直接改默认值。

### 添加新的攻击策略

如果想扩展新的恶意行为，推荐按以下步骤进行：

1. 在 [node.go](./node/node.go) 的 `Strategy` 中增加新策略枚举。
2. 在 `Node.Mine()`、`Node.Receive()`、`Node.UpdateTip()` 中加入对应逻辑。
3. 如需额外状态，可扩展 `ExpVar`。
4. 在 `exp/` 下新增实验函数与测试用例。
5. 通过 `go test -v ./...` 验证输出。

### 添加新的统计指标

当前项目已经统计：

- 平均 `ticks/block`
- 分叉攻击成功率
- 自私挖矿收益占比

如果需要扩展，可以考虑在 [simulator.go](./simulator/simulator.go) 或 `exp/` 中增加：

- 分叉数量
- 孤块比例
- 主链切换次数
- 平均传播路径长度

### 开发建议

- 保持测试驱动方式，优先在 `*_test.go` 中组织实验。
- 尽量复用固定随机种子，便于对比改动前后的结果。
- 修改攻击逻辑后，建议同时查看 `results/log/example.log` 风格输出，确认趋势是否符合预期。
- 若引入更复杂网络模型，建议优先保持接口清晰，而不是直接把逻辑耦合到 `Simulator.Run()` 中。

## 局限性说明

该项目是课程实验仿真器，不是生产级区块链实现。它有以下刻意简化：

- 不进行真实哈希求解。
- 不实现交易、手续费、难度调整。
- 不实现真实网络拓扑与带宽瓶颈。
- 不实现完整论文级自私挖矿状态机。

这些简化不会妨碍观察 PoW 的核心现象，但也意味着结果更适合用于趋势分析，而不是直接映射到真实公链。

## 许可证

本项目采用 [MIT License](./LICENSE)。
