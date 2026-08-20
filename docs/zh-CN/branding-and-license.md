# TokenRouter 炬枢 · 品牌与许可

## 品牌

| 项 | 值 |
|---|---|
| 产品名（英文） | **TokenRouter** |
| 产品名（中文） | **炬枢** |
| 标语 | 炬枢 · 企业 AI 网关 |
| 品牌标识 | 圆角六边形徽章 + 光束汇聚 + 琥珀光点（`frontend/public/brand/tokenrouter-logo.svg`） |

TokenRouter 与炬枢为产品专属品牌。未经授权，请勿在衍生项目中复用本标识。

## 开源许可

本项目基于 **Apache License 2.0** 发布（见仓库根目录 `LICENSE`）。

- ✅ 允许商用、修改、再分发、专利授权
- ✅ 允许闭源商用分发
- ✅ 允许基于本项目衍生自有品牌产品
- 📌 义务：保留版权声明与许可文本、声明修改、保留 NOTICE（如有）

## 上游项目致谢

本项目由 [astaxie/TokenHub](https://github.com/astaxie/TokenHub)（Apache-2.0）改造而来，上游由 Asta Xie 及社区维护。

**改造范围**：

| 层面 | 说明 |
|---|---|
| 品牌层 | 产品名、Logo、标识、文案（中/英/日）替换为 TokenRouter 炬枢 |
| 设计层 | 设计系统（`design-system/tokenrouter/MASTER.md`）、SaaS 风格界面、React Bits 动效、图表交互升级 |
| 功能层 | 保持与上游一致的网关能力与 API 契约（OpenAI 兼容 `/v1`、管理 API、环境变量体系） |

**合规要点**：

- 上游 Apache-2.0 版权声明保留于 `LICENSE` 与相关文件头部
- 环境变量、CLI 命令、存储键等内部契约沿用上游命名（`TOKENHUB_*`），确保与上游文档、工具链、迁移路径兼容
- 上游 GitHub 引用（版本检查、升级脚本）保持指向原仓库，便于跟踪上游更新

## 商标声明

TokenRouter / 炬枢 与 astaxie/TokenHub 及其作者**无隶属关系、未经其背书**；TokenHub 为其原作者的商标/项目名，归属其各自所有者。本项目的品牌仅代表本项目自身。

## 贡献与协作

- 代码风格与仓库门禁遵循仓库内 `AGENTS.md`（英文提交信息、三语文案同步、门禁检查）
- 文档三语同步：简体中文为 UI 文案基准，英文与日文翻译对齐
