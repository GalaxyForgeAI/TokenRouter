# TokenRouter 炬枢 · 快速开始

5 分钟跑通：启动网关 → 登录控制台 → 接入第一个 Provider → 用 API Key 调通模型。

## 方式一：Docker Compose（推荐）

```bash
git clone <your-tokenrouter-repo-url> && cd tokenrouter
docker compose --env-file deploy/.env.example up -d
```

服务就绪后：

| 服务 | 地址 |
|---|---|
| 管理控制台 | `http://localhost:3000` |
| 网关 API | `http://localhost:8080` |
| 健康检查 | `http://localhost:8080/healthz` |

## 方式二：本地开发

要求：Go 1.26+、Node.js 20+。

```bash
# 1. 启动后端（SQLite，默认端口 8080）
cd backend
go build -o tokenrouter ./cmd/tokenhub
TOKENHUB_ENV=dev ./tokenrouter

# 2. 启动前端（新终端，端口 3000）
cd frontend
npm install
TOKENHUB_API_BASE_URL=http://localhost:8080 npm run dev
```

> 前端默认连 `http://localhost:8080`；后端换过端口时用 `TOKENHUB_API_BASE_URL` 指过去。

## 初始化与登录

1. 打开 `http://localhost:3000`
2. 使用默认管理员登录：**admin / admin123456**（首次登录后请立即修改密码）
3. 登录后进入「平台总览」

> 默认管理员身份可在 `TOKENHUB_BOOTSTRAP_ADMIN_PASSWORD` 环境变量中预置。

## 接入第一个 Provider

1. 左侧导航 → **AI 资源 → Provider 渠道** → **新增 Provider**
2. 从 150+ 模板目录选择或选「自定义 OpenAI 兼容」：
   - 名称：如 `DeepSeek 生产渠道`
   - Base URL：`https://api.deepseek.com/v1`
   - API Key：填写你的上游密钥
3. 保存后**引入模型**（自动拉取上游模型列表，可勾选）
4. 到**路由策略**为模型创建路由（选择渠道 + 权重/优先级），保存

## 创建 API Key 并调用

1. **组织治理 → 项目空间** → 打开项目 → **Key 管理 → 新增 Key**
2. 复制生成的 API Key（只显示一次）

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer <你的API Key>" \
  -H "content-type: application/json" \
  -d '{
    "model": "<路由中启用的模型名>",
    "messages": [{"role": "user", "content": "你好"}]
  }'
```

成功即返回 OpenAI 格式响应。此后：

- **请求日志**中能看到这次调用的完整链路（路由、上游、Token、成本）
- **用量统计**按项目/模型/成本中心归因
- **成本看板**实时汇总消耗

## 进阶

| 想做什么 | 文档 |
|---|---|
| 完整部署（PostgreSQL、反向代理、环境变量） | [部署指南](https://github.com/GalaxyForgeAI/TokenRouter/blob/main/docs/zh-CN/deployment.md) |
| 三种角色的完整操作手册 | [用户指南](https://github.com/GalaxyForgeAI/TokenRouter/blob/main/docs/zh-CN/user-guide.md) / [团队负责人指南](https://github.com/GalaxyForgeAI/TokenRouter/blob/main/docs/zh-CN/team-leader-guide.md) / [管理员指南](https://github.com/GalaxyForgeAI/TokenRouter/blob/main/docs/zh-CN/administrator-guide.md) |
| 接入 Codex 订阅账号 | [Codex 接入指南](https://github.com/GalaxyForgeAI/TokenRouter/blob/main/docs/zh-CN/codex-tokenhub-profile-quick-start.md) |
| 架构与数据流 | [架构说明](https://github.com/GalaxyForgeAI/TokenRouter/blob/main/docs/zh-CN/architecture.md) |
