# TokenRouter 炬枢 · 变更记录

## v0.6.0-tokenrouter.1（2026-08-20）

TokenRouter 首个品牌化发布。基于上游 TokenHub v0.6.0（Apache-2.0）。

### 品牌化

- 产品更名 **TokenHub → TokenRouter 炬枢**：全站 UI 文案（中/英/日三语）、页面标题、侧边栏、登录页、面包屑
- 全新品牌标识：圆角六边形徽章 + 光束汇聚 + 琥珀光点（SVG/PNG/ICO），替换旧 logo 与 favicon
- 品牌锁字重构：中文「炬枢」为主角（侧边栏 21px 主位、登录页 36px 大号锁字），英文 TokenRouter 为宽字距副标
- 版本状态移入侧边栏底部，仅管理员可见、可操作升级

### 设计与体验

- **设计系统**：`design-system/tokenrouter/MASTER.md`（ui-ux-pro-max 生成）——主色 `#1E40AF`、背景 `#F8FAFC`、Fira Sans + Fira Code、数据密集企业仪表盘规范
- **SaaS 风格全局升级**：内容区品牌氛围光、侧边栏渐变、渐变主按钮、16px 圆角面板、细滚动条、键盘焦点光环、表格 hover、状态胶囊主题化
- **React Bits 动效**：KPI 数字滚动（CountUp）、卡片聚光（SpotlightCard）、渐变标题（GradientText）、按钮光泽（ShinyText）、视图切换淡入、导航激活滑条（motion layoutId）
- **图表高级化**：趋势图平滑曲线 + 渐变描边 + hover 竖线 tooltip + 末点光晕；环形占比图 SVG 化（12 点起始、悬停联动图例）；部门 Token 消耗对比柱状图 hover 明细；Top 模型渐变进度条
- **请求日志重排**：列表单行 meta（· 分隔）、品牌蓝选中态；详情面板视口内滚动、元信息分割线网格、payload 双列并排
- 修复导航动画陈旧内容 bug（AnimatePresence mode="wait" → keyed motion.div）
- 无障碍：`prefers-reduced-motion` 降级、可见焦点

### 工程与演示

- `tools/mock-llm/server.js`：零依赖 OpenAI 兼容 mock 上游，用于本地点亮观测页
- 端到端数据链路验证：Provider → 模型导入 → 路由 → API Key → `/v1` 调用（含流式）→ 用量/成本/审计全链路
- 仓库门禁全绿：`next build` / `tsc` / `eslint` / 126 项 node 测试 / 三语翻译一致性 / 环境变量契约 / 行数基线

## 上游基线

- **TokenHub v0.6.0**（astaxie/TokenHub，Apache-2.0）：企业 AI 网关核心能力——OpenAI/Anthropic 兼容 API、Provider 渠道与模板目录、模型路由策略、项目级 Key、用量成本归因、身份源与 RBAC、审计与告警、SQLite/PostgreSQL。
