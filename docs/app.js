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

const routineSnapshot = {
  overview: { longest_height: 92, active_tips: 3, total_forks: 2, blocks_per_second: 9.1988263589 },
};
const attackData = [
  { rate: 0.1, value: 0.0003 }, { rate: 0.2, value: 0.011 },
  { rate: 0.3, value: 0.0823 }, { rate: 0.4, value: 0.2983 },
];

const translations = {
  en: {
    'nav.models': 'models', 'nav.experiments': 'experiments', 'nav.data': 'data',
    'hero.eyebrow': 'open experiment / deterministic by design', 'hero.title': 'Two clocks.<br /><em>One chain.</em>',
    'hero.lede': 'A visual field guide to Proof-of-Work behavior across discrete time and real concurrent network routines.',
    'hero.results': 'Explore the results', 'hero.compare': 'Compare the models',
    'visual.propagation': 'propagation field', 'visual.live': 'LIVE', 'visual.height': 'height', 'visual.extends': 'chain extends',
    'signals.control': 'CONTROL', 'signals.controlText': 'reproducible experiments', 'signals.observe': 'OBSERVE', 'signals.observeText': 'forks, tips, miner share',
    'signals.compare': 'COMPARE', 'signals.compareText': 'same PoW, different clocks', 'signals.export': 'EXPORT', 'signals.exportText': 'human logs / machine JSON',
    'models.kicker': 'the setup', 'models.title': 'One phenomenon,<br /><span>two observatories.</span>', 'models.intro': 'The project keeps the chain model intentionally simple, then changes the clock and network around it. That makes the effect of latency visible instead of burying it in infrastructure.',
    'tick.tag': 'DISCRETE / CENTRALIZED', 'tick.title': 'Tick simulator', 'tick.description': 'Every node tries, packets wait in time buckets, then tips reconcile. A controlled lens for comparing difficulty, population and attack rate.', 'tick.rounds': '10000 rounds',
    'routine.tag': 'CONTINUOUS / GOSSIP', 'routine.title': 'Routine simulator', 'routine.description': 'Independent goroutines mine and gossip over real TCP connections. Bootnode observes the DAG while peers discover their own rhythm.', 'routine.peers': '10 peers / 2° density',
    'common.source': 'view source ↗', 'experiments.kicker': 'observed output', 'experiments.title': 'Signals from<br /><span>the sandbox.</span>', 'experiments.overview': 'overview', 'experiments.attack': 'attack curve', 'experiments.sample': 'sample / seed 42',
    'metrics.height': 'LONGEST CHAIN', 'metrics.heightText': 'blocks / routine sample', 'metrics.rate': 'BLOCK RATE', 'metrics.rateText': 'blocks per second', 'metrics.tips': 'ACTIVE TIPS', 'metrics.tipsText': 'visible chain endings', 'metrics.forks': 'FORK EVENTS', 'metrics.forksText': 'same-height collisions',
    'chart.kicker': 'FORK ATTACK / SUCCESS RATE', 'chart.title': 'When private work outruns public time', 'chart.unit': 'probability', 'insight.kicker': 'READOUT / 01', 'insight.title': 'Propagation is the variable.', 'insight.copy': 'With the same PoW rule, a delayed view of the chain creates room for competing tips. The Tick model makes that room measurable.', 'insight.valueText': 'success at 40% malicious rate',
    'architecture.kicker': 'under the hood', 'architecture.title': 'Simple rules.<br /><span>Rich behavior.</span>', 'architecture.intro': 'No real hashing, no transaction layer, no magic. Just enough machinery to make propagation delay, private chains and miner strategy legible.', 'architecture.mine': 'mine', 'architecture.mineText': 'probability × hash power', 'architecture.propagate': 'propagate', 'architecture.propagateText': 'delay or TCP gossip', 'architecture.reconcile': 'reconcile', 'architecture.reconcileText': 'highest visible tip wins', 'architecture.note': 'Both branches share the same probability rule. Only the clock, transport and observability layer change.',
    'footer.kicker': 'run the experiment', 'footer.title': 'Make the network<br /><em>show its work.</em>', 'footer.repository': 'Open repository <span>↗</span>', 'footer.download': 'Download sample JSON <span>↓</span>', 'footer.tagline': 'DISCRETE + CONTINUOUS / OBSERVE THE DIFFERENCE', 'footer.built': 'BUILT FOR EXPERIMENTS',
    'insight.attackTitle': 'Hash power compounds.', 'insight.attackCopy': 'At 40% malicious hash power, the sampled fork attack succeeds in 29.83% of 10,000 runs. The curve is deliberately non-linear.', 'insight.attackValue': '29.83%', 'insight.overviewValue': '3 tips',
    'download.note': 'sample visualization data',
  },
  zh: {
    'nav.models': '模型', 'nav.experiments': '实验', 'nav.data': '数据',
    'hero.eyebrow': '开放实验 / 结果可复现', 'hero.title': '两种时钟。<br /><em>同一条链。</em>',
    'hero.lede': '用一个可视化实验场，观察离散时间与真实并发网络中的工作量证明行为。',
    'hero.results': '查看实验结果', 'hero.compare': '对比两种模型',
    'visual.propagation': '传播场', 'visual.live': '运行中', 'visual.height': '高度', 'visual.extends': '链正在延伸',
    'signals.control': '控制', 'signals.controlText': '可复现实验', 'signals.observe': '观察', 'signals.observeText': '分叉、Tips、矿工占比',
    'signals.compare': '对比', 'signals.compareText': '同一 PoW，不同时钟', 'signals.export': '导出', 'signals.exportText': '可读日志 / 机器 JSON',
    'models.kicker': '实验设置', 'models.title': '同一个现象，<br /><span>两种观测方式。</span>', 'models.intro': '项目刻意保持区块链模型简单，只改变周围的时钟与网络，让传播延迟的影响可以被看见，而不是被基础设施隐藏。',
    'tick.tag': '离散 / 中心化', 'tick.title': 'Tick 仿真', 'tick.description': '每个节点依次尝试出块，数据包进入时间桶，最后统一协调链头。适合对比难度、节点规模和攻击比例。', 'tick.rounds': '10000 次实验',
    'routine.tag': '连续 / Gossip', 'routine.title': 'Routine 仿真', 'routine.description': '独立 goroutine 通过真实 TCP 连接挖矿并传播区块。Bootnode 观察 DAG，节点在网络中形成自己的节奏。', 'routine.peers': '10 节点 / 2° 密度',
    'common.source': '查看源码 ↗', 'experiments.kicker': '观测结果', 'experiments.title': '来自实验场的<br /><span>信号。</span>', 'experiments.overview': '概览', 'experiments.attack': '攻击曲线', 'experiments.sample': '样例 / 随机种子 42',
    'metrics.height': '最长链', 'metrics.heightText': 'Routine 样例区块', 'metrics.rate': '出块速率', 'metrics.rateText': '区块 / 秒', 'metrics.tips': '活跃 Tips', 'metrics.tipsText': '可见的链末端', 'metrics.forks': '分叉事件', 'metrics.forksText': '同高度竞争',
    'chart.kicker': '分叉攻击 / 成功率', 'chart.title': '当私有工作超过公共时间', 'chart.unit': '概率', 'insight.kicker': '读数 / 01', 'insight.title': '传播是关键变量。', 'insight.copy': 'PoW 规则不变时，对链状态的延迟观察会为竞争中的 Tips 留出空间。Tick 模型让这段空间可以被测量。', 'insight.valueText': '恶意比例 40% 时的成功率',
    'architecture.kicker': '内部机制', 'architecture.title': '简单规则。<br /><span>丰富行为。</span>', 'architecture.intro': '不做真实哈希，不加入交易层，也没有魔法。只保留足够解释传播延迟、私有链和矿工策略的机制。', 'architecture.mine': '挖矿', 'architecture.mineText': '概率 × 算力', 'architecture.propagate': '传播', 'architecture.propagateText': '延迟或 TCP Gossip', 'architecture.reconcile': '协调', 'architecture.reconcileText': '选择可见最高链头', 'architecture.note': '两个分支共享同一套出块概率规则，变化只有时钟、传输和观测层。',
    'footer.kicker': '运行实验', 'footer.title': '让网络<br /><em>展示它的工作。</em>', 'footer.repository': '打开仓库 <span>↗</span>', 'footer.download': '下载样例 JSON <span>↓</span>', 'footer.tagline': '离散 + 连续 / 观察差异', 'footer.built': '为实验而生',
    'insight.attackTitle': '算力优势会叠加。', 'insight.attackCopy': '当恶意算力达到 40% 时，样例中的分叉攻击在 10000 次实验中成功率为 29.83%，曲线呈现明显的非线性。', 'insight.attackValue': '29.83%', 'insight.overviewValue': '3 个 Tips',
    'download.note': '样例可视化数据',
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

const overview = routineSnapshot.overview;
document.querySelector('#metric-height').textContent = overview.longest_height;
document.querySelector('#metric-rate').textContent = Number(overview.blocks_per_second).toFixed(2);
document.querySelector('#metric-tips').textContent = String(overview.active_tips).padStart(2, '0');
document.querySelector('#metric-forks').textContent = String(overview.total_forks).padStart(2, '0');
document.querySelector('#visual-height').textContent = overview.longest_height;

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

document.querySelector('#download-data')?.addEventListener('click', (event) => {
  event.preventDefault();
  const payload = JSON.stringify({ source: 'PoW Simulator', tick: attackData, note: translations[getLanguage()]['download.note'] }, null, 2);
  const link = document.createElement('a'); link.href = URL.createObjectURL(new Blob([payload], { type: 'application/json' })); link.download = 'pow-simulator-sample.json'; link.click(); URL.revokeObjectURL(link.href);
});

if ('IntersectionObserver' in window) {
  const observer = new IntersectionObserver((entries) => entries.forEach((entry) => entry.isIntersecting && entry.target.classList.add('is-visible')), { threshold: 0.12 });
  document.querySelectorAll('.section, .signal-strip, .footer-cta').forEach((element) => observer.observe(element));
}
