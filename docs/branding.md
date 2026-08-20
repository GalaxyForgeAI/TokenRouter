# TokenRouter 炬枢 —— 品牌化改造说明

> 本仓库由 [astaxie/TokenHub](https://github.com/astaxie/TokenHub)（Apache-2.0）改造而来。
> 改造原则：**只改品牌显示，不动功能契约**（沿用 AxonHub POC 的既定方针）。

## 品牌

| 项 | 值 |
|---|---|
| 英文名 | TokenRouter |
| 中文名 | 炬枢 |
| 标语 | 炬枢 · 企业 AI 网关 |
| 设计系统 | `design-system/tokenrouter/MASTER.md`（ui-ux-pro-max 生成） |
| Logo 源文件 | `frontend/public/brand/tokenrouter-logo.svg`（六边形路由枢纽 × 火炬 × 琥珀 token 核心） |
| Logo PNG | `frontend/public/brand/tokenrouter-logo.png`（512px）/ `tokenrouter-logo-256.png` |

## 改了哪些（品牌显示层）

- 前端控制台全部 UI 可见文案：`TokenHub` → `TokenRouter`（含侧边栏、面包屑、登录页、帮助文案、i18n 三语字典键值）
- 登录页品牌锁字：`Token<span>Router</span>` + 副标「炬枢 · 企业 AI 网关」
- 侧边栏品牌区：新 Logo + `TokenRouter` + 「炬枢」徽标
- 页面 title：`TokenRouter 炬枢 · 管理控制台`
- `frontend/package.json` name：`tokenrouter-admin`
- Logo / favicon：全新炬枢标识，替换 `tokenhub-logo.png` 与 `app/favicon.ico`
- 顶层 README（en/zh-CN/ja）与 `docs/` 全量文案

## 没改哪些（功能契约，保持上游同步）

以下保持 `tokenhub` 原样，改名会破坏功能或与上游产生 merge 冲突：

- **环境变量**：`TOKENHUB_*`（`start.sh` / `deploy/` / `.env.example` 全链路契约，`tools/check-env-contract.mjs` 校验）
- **后端 Go 内部**：module 名 `tokenhub/backend`、CLI 命令 `tokenhub db migrate`、`TOKENHUB_RELEASE_REPOSITORY`
- **前端内部契约**：localStorage 键（`tokenhub.admin.session` 等）、自定义事件名（`tokenhub-issued-key`）、查询参数（`tokenhub_provider_account`）、下载文件名、默认 client_id / webhook URL
- **上游 GitHub 引用**：`version-status.tsx` 的 `releases_url`、部署文档中的 `astaxie/TokenHub` 安装/升级地址
- **开发文档**：`AGENTS.md` / `CLAUDE.md` / `.github/`（保持上游术语，便于同步上游更新）

> 若后续要彻底私有化（连内部契约一起改），需按「重命名影响面」逐项核对：
> env 契约（start.sh/deploy/.env.example/后端）→ 存储键（一次性的，改后旧会话登出）→ 事件名（前后端同改）。

## UI 改造（React Bits + 设计系统）

| 组件 | 位置 | 用途 |
|---|---|---|
| CountUp | `frontend/components/reactbits/CountUp/` | 概览 KPI 数字滚动（请求量/成功率/延迟/成本/Token） |
| SpotlightCard | `frontend/components/reactbits/SpotlightCard/` | KPI 卡鼠标跟随聚光（浅色主题已适配） |
| GradientText | `frontend/components/reactbits/GradientText/` | 渐变标题（浅色主题已适配） |
| motion (AnimatePresence) | 控制台壳 | 视图切换淡入过渡（`admin-console.tsx`） |
| motion (layoutId) | 侧边栏 | 导航激活指示条弹簧滑动（`navigation-ui.tsx`） |

设计令牌：`frontend/app/styles/legacy/tokens.css` 已按设计系统重写
（主色 `#1E40AF`、背景 `#F8FAFC`、白卡、边框 `#E2E8F0`、Fira Sans 正文 + Fira Code 数字）。

## 验证

- `cd frontend && npm ci && npm run build` 通过
- 登录 → 概览 → KPI 数字滚动 / 卡聚光 / 导航滑条 / 视图过渡正常
- `node tools/check-ui-translations.mjs`（i18n 三语键值完整）
