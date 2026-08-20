# TokenRouter · クイックスタート

管理されたゲートウェイを 5 分で起動します: 起動 → ログイン → プロバイダー接続 → API キーでモデル呼び出し。

## 方法 1: Docker Compose（推奨）

```bash
git clone <your-tokenrouter-repo-url> && cd tokenrouter
docker compose --env-file deploy/.env.example up -d
```

| サービス | アドレス |
|---|---|
| 管理コンソール | http://localhost:3000 |
| ゲートウェイ API | http://localhost:8080 |
| ヘルスチェック | http://localhost:8080/healthz |

## 方法 2: ローカル開発

要件: Go 1.26+、Node.js 20+。

```bash
# 1. バックエンド（SQLite、デフォルトポート 8080）
cd backend
go build -o tokenrouter ./cmd/tokenhub
TOKENHUB_ENV=dev ./tokenrouter

# 2. フロントエンド（別ターミナル、ポート 3000）
cd frontend
npm install
TOKENHUB_API_BASE_URL=http://localhost:8080 npm run dev
```

> フロントエンドはデフォルトで `http://localhost:8080` を向きます。バックエンドのポートを変えた場合は `TOKENHUB_API_BASE_URL` で指定します。

## 初期化とログイン

1. http://localhost:3000 を開く
2. ブートストラップ管理者でログイン: **admin / admin123456**（初回ログイン後すぐにパスワード変更）
3. ダッシュボードが表示される

> ブートストラップ管理者パスワードは `TOKENHUB_BOOTSTRAP_ADMIN_PASSWORD` 環境変数で事前設定できます。

## 最初のプロバイダー接続

1. サイドバー → **AI リソース → Provider チャネル** → **プロバイダー追加**
2. 150+ テンプレートから選択、または「カスタム OpenAI 互換」を選択:
   - 名前: 例 `DeepSeek Production`
   - Base URL: `https://api.deepseek.com/v1`
   - API キー: 上流の認証情報
3. 保存後、**モデルをインポート**（上流モデル一覧が自動取得）
4. **ルーティングポリシー**でモデルにルートを作成（チャネル + 重み/優先度）して保存

## API キー作成と呼び出し

1. **組織 → プロジェクトスペース** → プロジェクトを開く → **キー管理 → キー追加**
2. 生成された API キーをコピー（一度だけ表示）

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer <YOUR_API_KEY>" \
  -H "content-type: application/json" \
  -d '{
    "model": "<routingで有効なモデル名>",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```

成功すると OpenAI 形式のレスポンスが返ります。以降:

- **リクエストログ**に呼び出し経路（ルート・上流・トークン・コスト）が記録されます
- **使用量分析**がプロジェクト/モデル/コストセンター別に集計されます
- **コストダッシュボード**がリアルタイム集計します

## 次のステップ

| 目的 | ドキュメント |
|---|---|
| 本番デプロイ（PostgreSQL、リバースプロキシ） | [デプロイガイド](deployment.md) |
| ロール別マニュアル | [ユーザーガイド](user-guide.md) / [チームリーダーガイド](team-leader-guide.md) / [管理者ガイド](administrator-guide.md) |
| Codex サブスクリプション接続 | [Codex クイックセットアップ](codex-tokenhub-profile-quick-start.md) |
| アーキテクチャ | [アーキテクチャ](architecture.md) |
