const root = document.documentElement;

const storageGet = (key) => {
  try { return localStorage.getItem(key); } catch { return null; }
};
const storageSet = (key, value) => {
  try { localStorage.setItem(key, value); } catch { /* file:// may disable storage */ }
};

const themeButton = document.querySelector('#theme-toggle');
const themeOrder = ['system', 'light', 'dark'];
const getTheme = () => storageGet('pow-theme') || 'system';
const applyTheme = (theme) => {
  root.dataset.theme = theme;
  storageSet('pow-theme', theme);
  themeButton?.setAttribute('aria-label', `Theme: ${theme}`);
};
themeButton?.addEventListener('click', () => applyTheme(themeOrder[(themeOrder.indexOf(getTheme()) + 1) % themeOrder.length]));

const attackData = [
  { rate: 0.1, value: 0.0003 }, { rate: 0.2, value: 0.011 },
  { rate: 0.3, value: 0.0823 }, { rate: 0.4, value: 0.2983 },
];

const translations = {
  en: {
    'nav.models': 'models', 'nav.experiments': 'experiments', 'nav.data': 'data',
    'hero.eyebrow': 'open experiment / deterministic by design', 'hero.title': '<em>Gears and starfields.</em><br /><span class="hero-title-line">Proof of work,</span><br /><span class="hero-title-line">two clocks.</span>',
    'hero.lede': 'The Tick model advances every node on one synchronized clock; the Routine model lets goroutines mine concurrently, with a 500 μs sleep after a failed attempt.',
    'hero.results': 'Explore the results', 'hero.compare': 'Compare the models',
    'visual.propagation': 'chain state', 'visual.live': 'LIVE', 'visual.sample': 'sample', 'visual.extends': 'block race in progress', 'visual.mainChain': 'canonical chain', 'visual.forkChain': 'competing branch', 'visual.forkNote': 'fork reconstructed from log',
    'signals.control': 'CONTROL', 'signals.controlText': 'reproducible experiments', 'signals.observe': 'OBSERVE', 'signals.observeText': 'forks, tips, miner share',
    'signals.compare': 'COMPARE', 'signals.compareText': 'same PoW, different clocks', 'signals.export': 'EXPORT', 'signals.exportText': 'human logs / machine JSON',
    'models.kicker': 'the setup', 'models.title': 'One phenomenon,<br /><span>two observatories.</span>', 'models.intro': 'The project keeps the chain model intentionally simple, then changes the clock and network around it. That makes the effect of latency visible instead of burying it in infrastructure.',
    'tick.tag': 'DISCRETE / CENTRALIZED', 'tick.title': 'Tick simulator', 'tick.description': 'Every node tries, packets wait in time buckets, then tips reconcile. A controlled lens for comparing difficulty, population and attack rate.', 'tick.rounds': '10000 rounds',
    'routine.tag': 'CONTINUOUS / GOSSIP', 'routine.title': 'Routine simulator', 'routine.description': 'Independent goroutines mine and gossip over real TCP connections. Bootnode observes the DAG while peers discover their own rhythm.', 'routine.peers': '10 peers / 2° density',
    'common.source': 'view source ↗', 'experiments.kicker': 'experiment families', 'experiments.title': 'Three lenses.<br /><span>one proof rule.</span>', 'experiments.overview': 'overview', 'experiments.attack': 'attack curve', 'experiments.sample': 'representative runs / seed 42',
    'family.blocks.kicker': 'BLOCKS / FORKS', 'family.blocks.title': 'From a clean tip to a crowded DAG', 'family.blocks.copy': 'Tick advances in synchronized rounds; Routine mines continuously and sleeps 500 μs after a failed attempt.', 'family.selfish.kicker': 'SELFISH MINING', 'family.selfish.title': 'Private blocks, public revenue', 'family.selfish.copy': 'A miner withholds work, then releases it when the private lead is valuable.', 'family.attack.kicker': 'FORK ATTACK', 'family.attack.title': 'Can a private branch catch up?', 'family.attack.copy': 'The attacker mines privately while honest miners extend the public branch.', 'table.setup': 'setup', 'table.result': 'result',
    'metrics.height': 'LONGEST CHAIN', 'metrics.heightText': 'blocks / routine sample', 'metrics.rate': 'BLOCK RATE', 'metrics.rateText': 'blocks per second', 'metrics.tips': 'ACTIVE TIPS', 'metrics.tipsText': 'visible chain endings', 'metrics.forks': 'FORK EVENTS', 'metrics.forksText': 'same-height collisions',
    'chart.kicker': 'FORK ATTACK / SUCCESS RATE', 'chart.title': 'When private work outruns public time', 'chart.unit': 'probability', 'insight.kicker': 'READOUT / 01', 'insight.title': 'Propagation is the variable.', 'insight.copy': 'With the same PoW rule, a delayed view of the chain creates room for competing tips. The Tick model makes that room measurable.', 'insight.valueText': 'success at 40% malicious rate',
    'architecture.kicker': 'under the hood', 'architecture.title': 'Simple rules.<br /><span>Rich behavior.</span>', 'architecture.intro': 'No real hashing, no transaction layer, no magic. Just enough machinery to make propagation delay, private chains and miner strategy legible.', 'architecture.mine': 'mine', 'architecture.mineText': 'probability × hash power', 'architecture.propagate': 'propagate', 'architecture.propagateText': 'delay or TCP gossip', 'architecture.reconcile': 'reconcile', 'architecture.reconcileText': 'highest visible tip wins', 'architecture.note': 'Both branches share the same probability rule. Only the clock, transport and observability layer change.',
    'footer.kicker': 'run the experiment', 'footer.title': 'Make the network<br /><em>show its work.</em>', 'footer.repository': 'Open repository <span>↗</span>', 'footer.download': 'Download Routine example log <span>↓</span>', 'footer.tagline': 'DISCRETE + CONTINUOUS / OBSERVE THE DIFFERENCE', 'footer.built': 'BUILT FOR EXPERIMENTS',
    'insight.attackTitle': 'Hash power compounds.', 'insight.attackCopy': 'At 40% malicious hash power, the sampled fork attack succeeds in 29.83% of 10,000 runs. The curve is deliberately non-linear.', 'insight.attackValue': '29.83%', 'insight.overviewValue': '3 tips',
    'download.note': 'sample visualization data',
    'routineLab.kicker': 'ROUTINE / BLOCK RACE', 'routineLab.title': 'Watch blocks arrive and split', 'routineLab.pause': 'pause', 'routineLab.play': 'play',
  },
  zh: {
    'nav.models': '模型', 'nav.experiments': '实验', 'nav.data': '数据',
    'hero.eyebrow': '开放实验 / 结果可复现', 'hero.title': '<em>齿轮与星群。</em><br /><span class="hero-title-line">两种时钟下的</span><br /><span class="hero-title-line">工作量证明。</span>',
    'hero.lede': 'Tick 模型让所有节点在同一同步时钟下推进；Routine 模型让 goroutine 并发挖矿，挖矿失败后休眠 500 微秒。',
    'hero.results': '查看实验结果', 'hero.compare': '对比两种模型',
    'visual.propagation': '链状态', 'visual.live': '运行中', 'visual.sample': '样例', 'visual.extends': '出块竞速进行中', 'visual.mainChain': '主链', 'visual.forkChain': '竞争分支', 'visual.forkNote': '根据日志还原分叉',
    'signals.control': '控制', 'signals.controlText': '可复现实验', 'signals.observe': '观察', 'signals.observeText': '分叉、Tips、矿工占比',
    'signals.compare': '对比', 'signals.compareText': '同一 PoW，不同时钟', 'signals.export': '导出', 'signals.exportText': '可读日志 / 机器 JSON',
    'models.kicker': '实验设置', 'models.title': '同一个现象，<br /><span>两种观测方式。</span>', 'models.intro': '项目刻意保持区块链模型简单，只改变周围的时钟与网络，让传播延迟的影响可以被看见，而不是被基础设施隐藏。',
    'tick.tag': '离散 / 中心化', 'tick.title': 'Tick 仿真', 'tick.description': '每个节点依次尝试出块，数据包进入时间桶，最后统一协调链头。适合对比难度、节点规模和攻击比例。', 'tick.rounds': '10000 次实验',
    'routine.tag': '连续 / Gossip', 'routine.title': 'Routine 仿真', 'routine.description': '独立 goroutine 通过真实 TCP 连接挖矿并传播区块。Bootnode 观察 DAG，节点在网络中形成自己的节奏。', 'routine.peers': '10 节点 / 2° 密度',
    'common.source': '查看源码 ↗', 'experiments.kicker': '实验类型', 'experiments.title': '三种视角。<br /><span>同一套规则。</span>', 'experiments.overview': '概览', 'experiments.attack': '攻击曲线', 'experiments.sample': '代表性结果 / 随机种子 42',
    'family.blocks.kicker': '出块 / 分叉', 'family.blocks.title': '从干净链头到拥挤 DAG', 'family.blocks.copy': 'Tick 按同步轮次推进；Routine 持续并发挖矿，失败后休眠 500 微秒。', 'family.selfish.kicker': '自私挖矿', 'family.selfish.title': '私有区块与公开收益', 'family.selfish.copy': '矿工隐藏已挖出的区块，在私有链领先时择机发布。', 'family.attack.kicker': '分叉攻击', 'family.attack.title': '私有分支能否追上？', 'family.attack.copy': '攻击者私下挖矿，同时诚实矿工继续延伸公开分支。', 'table.setup': '参数组合', 'table.result': '结果',
    'metrics.height': '最长链', 'metrics.heightText': 'Routine 样例区块', 'metrics.rate': '出块速率', 'metrics.rateText': '区块 / 秒', 'metrics.tips': '活跃 Tips', 'metrics.tipsText': '可见的链末端', 'metrics.forks': '分叉事件', 'metrics.forksText': '同高度竞争',
    'chart.kicker': '分叉攻击 / 成功率', 'chart.title': '当私有工作超过公共时间', 'chart.unit': '概率', 'insight.kicker': '读数 / 01', 'insight.title': '传播是关键变量。', 'insight.copy': 'PoW 规则不变时，对链状态的延迟观察会为竞争中的 Tips 留出空间。Tick 模型让这段空间可以被测量。', 'insight.valueText': '恶意比例 40% 时的成功率',
    'architecture.kicker': '内部机制', 'architecture.title': '简单规则。<br /><span>丰富行为。</span>', 'architecture.intro': '不做真实哈希，不加入交易层，也没有魔法。只保留足够解释传播延迟、私有链和矿工策略的机制。', 'architecture.mine': '挖矿', 'architecture.mineText': '概率 × 算力', 'architecture.propagate': '传播', 'architecture.propagateText': '延迟或 TCP Gossip', 'architecture.reconcile': '协调', 'architecture.reconcileText': '选择可见最高链头', 'architecture.note': '两个分支共享同一套出块概率规则，变化只有时钟、传输和观测层。',
    'footer.kicker': '运行实验', 'footer.title': '让网络<br /><em>展示它的工作。</em>', 'footer.repository': '打开仓库 <span>↗</span>', 'footer.download': '下载 Routine 示例日志 <span>↓</span>', 'footer.tagline': '离散 + 连续 / 观察差异', 'footer.built': '为实验而生',
    'insight.attackTitle': '算力优势会叠加。', 'insight.attackCopy': '当恶意算力达到 40% 时，样例中的分叉攻击在 10000 次实验中成功率为 29.83%，曲线呈现明显的非线性。', 'insight.attackValue': '29.83%', 'insight.overviewValue': '3 个 Tips',
    'download.note': '样例可视化数据',
    'routineLab.kicker': 'ROUTINE / 出块竞速', 'routineLab.title': '观察区块到达与分叉', 'routineLab.pause': '暂停', 'routineLab.play': '播放',
  },
};

const browserLanguage = (navigator.languages?.[0] || navigator.language || 'en').toLowerCase();
const getLanguage = () => storageGet('pow-language') || (browserLanguage.startsWith('zh') ? 'zh' : 'en');
const languageButton = document.querySelector('#language-toggle');
const applyLanguage = (language) => {
  const dictionary = translations[language];
  root.lang = language === 'zh' ? 'zh-CN' : 'en';
  document.title = language === 'zh' ? 'PoW Simulator / 工作量证明实验场' : 'PoW Simulator / experiment lab';
  document.querySelector('meta[name="description"]')?.setAttribute('content', language === 'zh' ? 'PoW Simulator：用离散 Tick 与实时并发 Routine 对比工作量证明网络行为。' : 'PoW Simulator - discrete ticks and real concurrent routines in one experiment lab.');
  document.querySelectorAll('[data-i18n]').forEach((element) => {
    const value = dictionary[element.dataset.i18n];
    if (value) element.innerHTML = value;
  });
  languageButton?.setAttribute('aria-label', language === 'zh' ? '切换到英文' : 'Switch to Chinese');
  storageSet('pow-language', language);
};
applyLanguage(getLanguage());
languageButton?.addEventListener('click', () => applyLanguage(getLanguage() === 'zh' ? 'en' : 'zh'));

const line = document.querySelector('#chart-line');
const area = document.querySelector('#chart-area');
const points = document.querySelector('#chart-points');
const x = (index) => 18 + index * 228;
const y = (value) => 245 - value * 560;
const linePath = attackData.map((point, index) => `${index ? 'L' : 'M'} ${x(index)} ${y(point.value)}`).join(' ');
line?.setAttribute('d', linePath);
area?.setAttribute('d', `${linePath} L ${x(3)} 245 L ${x(0)} 245 Z`);
attackData.forEach((point, index) => {
  const circle = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
  circle.setAttribute('cx', x(index)); circle.setAttribute('cy', y(point.value)); circle.setAttribute('r', '5');
  points?.append(circle);
});

document.querySelectorAll('[data-view]').forEach((button) => button.addEventListener('click', () => {
  document.querySelectorAll('[data-view]').forEach((item) => item.classList.remove('active'));
  button.classList.add('active');
  const attack = button.dataset.view === 'attack';
  const dictionary = translations[getLanguage()];
  document.querySelector('#insight-title').innerHTML = dictionary[attack ? 'insight.attackTitle' : 'insight.title'];
  document.querySelector('#insight-copy').innerHTML = dictionary[attack ? 'insight.attackCopy' : 'insight.copy'];
  document.querySelector('#insight-value').textContent = dictionary[attack ? 'insight.attackValue' : 'insight.overviewValue'];
  document.querySelector('#insight-value').nextElementSibling.innerHTML = dictionary[attack ? 'insight.valueText' : 'insight.valueText'];
}));

const routineCarousel = document.querySelector('#routine-carousel');
const routineParameters = document.querySelector('#chain-parameters');
const routinePagination = document.querySelector('#chain-pagination');
const routineSample = document.querySelector('#visual-sample');
const routineSource = document.querySelector('#chain-source');
const routineSummary = document.querySelector('#chain-summary');
const routineDetail = document.querySelector('#chain-detail');
let routineScenes = [];
let routineSceneIndex = 0;
let routineSceneTimer;

const createSvgElement = (tag, attributes, parent, label) => {
  const element = document.createElementNS('http://www.w3.org/2000/svg', tag);
  Object.entries(attributes).forEach(([key, value]) => element.setAttribute(key, value));
  if (label) element.textContent = label;
  parent?.append(element);
  return element;
};
const shortHash = (hash) => hash ? hash.slice(0, 6) : 'genesis';
const drawRoutineScene = (scene, index) => {
  if (!routineCarousel || !scene) return;
  routineCarousel.classList.add('is-changing');
  window.setTimeout(() => routineCarousel.classList.remove('is-changing'), 650);
  routineCarousel.replaceChildren();
  const main = scene.main_chain || [];
  const left = 42; const right = 578; const mainY = 294;
  const xStep = main.length > 1 ? (right - left) / (main.length - 1) : 0;
  const positions = new Map();
  main.forEach((block, blockIndex) => {
    positions.set(block.hash, { x: left + blockIndex * xStep, y: mainY, block });
  });
  if (main.length > 1) {
    createSvgElement('path', { d: `M ${left} ${mainY} L ${right} ${mainY}`, class: 'log-main-track' }, routineCarousel);
  }
  main.forEach((block, blockIndex) => {
    const point = positions.get(block.hash);
    const node = createSvgElement('circle', { cx: point.x, cy: point.y, r: blockIndex === main.length - 1 ? 11 : 8, class: `log-main-node${blockIndex === main.length - 1 ? ' tip-node' : ''}`, tabindex: '0' }, routineCarousel);
    createSvgElement('text', { x: point.x, y: mainY + (blockIndex % 2 ? 29 : -20), class: 'log-height', 'text-anchor': 'middle' }, routineCarousel, `H${block.height}`);
    createSvgElement('text', { x: point.x, y: mainY + (blockIndex % 2 ? 42 : -7), class: 'log-miner', 'text-anchor': 'middle' }, routineCarousel, `M${block.miner}`);
    createSvgElement('title', {}, node, `height ${block.height} · miner ${block.miner} · ${block.hash}`);
  });
  (scene.forks || []).forEach((fork, forkIndex) => {
    const parent = positions.get(fork.parent_hash);
    if (!parent) return;
    fork.branches.forEach((branch, branchIndex) => {
      const blocks = branch.blocks || [];
      if (!blocks.length) return;
      const branchY = 180 - ((forkIndex + branchIndex) % 3) * 38;
      const branchStartX = Math.min(parent.x + 8, right - 105);
      const branchEndX = Math.min(branchStartX + Math.max(54, blocks.length * 38), right - 4);
      const branchPath = `M ${parent.x} ${parent.y} C ${parent.x + 14} ${parent.y}, ${branchStartX - 12} ${branchY}, ${branchStartX} ${branchY}`;
      createSvgElement('path', { d: branchPath, class: 'log-fork-track' }, routineCarousel);
      blocks.forEach((block, blockIndex) => {
        const x = Math.min(branchStartX + blockIndex * ((branchEndX - branchStartX) / Math.max(1, blocks.length - 1)), right - 4);
        const y = branchY - blockIndex * 16;
        const forkNode = createSvgElement('circle', { cx: x, cy: y, r: 7, class: 'log-fork-node', tabindex: '0' }, routineCarousel);
        createSvgElement('text', { x, y: y - 12, class: 'log-height fork-height', 'text-anchor': 'middle' }, routineCarousel, `H${block.height}`);
        createSvgElement('title', {}, forkNode, `fork height ${block.height} · miner ${block.miner} · ${block.hash}`);
      });
    });
  });
  const p = scene.parameters;
  const o = scene.overview;
  routineParameters.textContent = `N${p.nodes} · d ${Number(p.difficulty).toFixed(3)} · ρ ${Number(p.density).toFixed(2)} · sleep 500 μs · ${scene.source_log}`;
  routineSource.textContent = `ROUTINE / ${scene.id.toUpperCase()}`;
  routineSummary.textContent = `${o.total_blocks_excluding_genesis} blocks · ${o.total_forks} forks`;
  routineDetail.textContent = `H${o.longest_chain_height} · ${o.active_tips} tips · ${Number(o.average_block_rate).toFixed(2)} blocks/sec`;
  if (routineSample) routineSample.textContent = `${String(index + 1).padStart(2, '0')} / ${String(routineScenes.length).padStart(2, '0')}`;
  routinePagination?.querySelectorAll('button').forEach((button, buttonIndex) => button.classList.toggle('active', buttonIndex === index));
};
const setupRoutineCarousel = (payload) => {
  routineScenes = payload.scenarios || [];
  if (!routineScenes.length) return;
  routinePagination.innerHTML = '';
  routineScenes.forEach((scene, index) => {
    const button = document.createElement('button');
    button.type = 'button'; button.textContent = String(index + 1).padStart(2, '0');
    button.setAttribute('aria-label', `Show ${scene.label}`);
    button.addEventListener('click', () => { routineSceneIndex = index; drawRoutineScene(routineScenes[index], index); });
    routinePagination.append(button);
  });
  drawRoutineScene(routineScenes[0], 0);
  routineSceneTimer = window.setInterval(() => { routineSceneIndex = (routineSceneIndex + 1) % routineScenes.length; drawRoutineScene(routineScenes[routineSceneIndex], routineSceneIndex); }, 6500);
};
fetch('./data/routine-visuals.json').then((response) => response.ok ? response.json() : Promise.reject(new Error('routine visual data unavailable'))).then(setupRoutineCarousel).catch(() => {
  if (routineSummary) routineSummary.textContent = 'Routine log unavailable';
  if (routineDetail) routineDetail.textContent = 'Open from the docs site or a local web server';
});

/* Legacy experiment chart code is retained for compatibility with older snapshots. */
/* const routineScenarios = {
  calm: {
    params: '5 nodes · difficulty 0.005 · density 0.30', result: '1 fork · 2 tips · height 30',
    events: [
      ['block mined', 'node 02 extends height 08', 'main'], ['propagated', 'node 04 receives block 08', 'main'],
      ['block mined', 'node 01 extends height 09', 'main'], ['block mined', 'node 03 extends height 10', 'main'],
      ['fork detected', 'node 00 mines a competing height 10', 'fork'], ['reconciled', 'longest visible tip wins', 'main'],
      ['block mined', 'node 02 extends height 11', 'main'], ['propagated', 'gossip reaches every peer', 'main'],
      ['block mined', 'node 04 extends height 12', 'main'], ['block mined', 'node 01 extends height 13', 'main'],
      ['reconciled', 'stale tip becomes inactive', 'main'], ['target reached', 'height 30 · 2 active tips', 'main'],
    ],
  },
  forks: {
    params: '10 nodes · difficulty 0.02 · density 0.60', result: '11 forks · 13 tips · height 31',
    events: [
      ['block mined', 'node 07 extends height 14', 'main'], ['block mined', 'node 03 extends height 14', 'fork'],
      ['fork detected', 'two tips share height 14', 'fork'], ['propagated', 'half the peers still see the fork', 'fork'],
      ['block mined', 'node 09 extends the orange tip', 'fork'], ['block mined', 'node 01 extends the cyan tip', 'main'],
      ['fork detected', 'height 15 collision', 'fork'], ['propagated', 'gossip closes the visibility gap', 'main'],
      ['block mined', 'node 05 extends height 16', 'main'], ['block mined', 'node 08 extends a side tip', 'fork'],
      ['reconciled', 'the longer branch becomes canonical', 'main'], ['target reached', 'height 31 · 13 active tips', 'fork'],
    ],
  },
  burst: {
    params: '20 nodes · difficulty 0.05 · density 0.60', result: '65 forks · 131 tips · height 72',
    events: [
      ['block mined', 'node 12 extends height 23', 'main'], ['block mined', 'node 04 extends height 23', 'fork'],
      ['block mined', 'node 19 extends height 23', 'fork'], ['fork detected', 'three miners race at one height', 'fork'],
      ['block mined', 'node 02 extends the first tip', 'main'], ['block mined', 'node 16 extends a private tip', 'fork'],
      ['propagated', 'packets fan out across 20 peers', 'main'], ['fork detected', 'new collisions form before convergence', 'fork'],
      ['block mined', 'node 07 wins the next race', 'main'], ['block mined', 'node 11 starts another branch', 'fork'],
      ['reconciled', 'only the highest visible tip advances', 'main'], ['target reached', 'height 72 · 131 active tips', 'fork'],
    ],
  },
};
const routineSvg = document.querySelector('#routine-chain');
const routineEvent = document.querySelector('#routine-event');
const routineDetail = document.querySelector('#routine-detail');
const routineStep = document.querySelector('#routine-step');
const routineParams = document.querySelector('#routine-params');
const routineResult = document.querySelector('#routine-result');
const playButton = document.querySelector('#routine-play');
let activeScenario = 'calm'; let routineIndex = 0; let routinePlaying = true;
const drawRoutine = () => {
  const scenario = routineScenarios[activeScenario]; const event = scenario.events[routineIndex];
  routineSvg.innerHTML = '';
  const ns = 'http://www.w3.org/2000/svg';
  const branch = event[2] === 'fork';
  const baseY = 155;
  const points = Array.from({ length: 9 }, (_, i) => ({ x: 62 + i * 94, y: i > 4 && branch ? 155 + (i - 4) * 17 : baseY }));
  const path = document.createElementNS(ns, 'path'); path.setAttribute('d', points.map((p, i) => `${i ? 'L' : 'M'} ${p.x} ${p.y}`).join(' ')); path.setAttribute('class', 'routine-main-line'); routineSvg.append(path);
  if (branch || routineIndex >= 4) { const forkPath = document.createElementNS(ns, 'path'); forkPath.setAttribute('d', 'M 438 155 C 500 108, 566 104, 720 92'); forkPath.setAttribute('class', 'routine-fork-line'); routineSvg.append(forkPath); }
  points.forEach((point, i) => { const circle = document.createElementNS(ns, 'circle'); circle.setAttribute('cx', point.x); circle.setAttribute('cy', point.y); circle.setAttribute('r', i <= routineIndex % 9 + 1 ? '8' : '5'); circle.setAttribute('class', i <= routineIndex % 9 + 1 ? (i === routineIndex % 9 + 1 && event[2] === 'fork' ? 'routine-node fork-node' : 'routine-node active-node') : 'routine-node'); routineSvg.append(circle); });
  routineStep.textContent = `EVENT ${String(routineIndex + 1).padStart(2, '0')} / ${scenario.events.length}`; routineEvent.textContent = event[0]; routineDetail.textContent = event[1]; routineParams.textContent = scenario.params; routineResult.textContent = scenario.result;
};
const setScenario = (scenario) => { activeScenario = scenario; routineIndex = 0; document.querySelectorAll('[data-scenario]').forEach((button) => button.classList.toggle('active', button.dataset.scenario === scenario)); drawRoutine(); };
document.querySelectorAll('[data-scenario]').forEach((button) => button.addEventListener('click', () => setScenario(button.dataset.scenario)));
playButton?.addEventListener('click', () => { routinePlaying = !routinePlaying; playButton.querySelector('.play-icon').textContent = routinePlaying ? 'Ⅱ' : '▶'; playButton.querySelector('.play-label').textContent = translations[getLanguage()][routinePlaying ? 'routineLab.pause' : 'routineLab.play']; playButton.setAttribute('aria-label', routinePlaying ? 'Pause simulation' : 'Play simulation'); });
setInterval(() => { if (!routinePlaying) return; routineIndex = (routineIndex + 1) % routineScenarios[activeScenario].events.length; drawRoutine(); }, 1500);
if (routineSvg) drawRoutine(); */

if ('IntersectionObserver' in window) {
  const observer = new IntersectionObserver((entries) => entries.forEach((entry) => entry.isIntersecting && entry.target.classList.add('is-visible')), { threshold: 0.12 });
  document.querySelectorAll('.section, .signal-strip, .footer-cta').forEach((element) => observer.observe(element));
}
