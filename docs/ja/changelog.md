# TokenRouter · 変更履歴

## v0.6.0-tokenrouter.1（2026-08-20）

TokenRouter 初のブランドリリース。上流 TokenHub v0.6.0（Apache-2.0）ベース。

### ブランド化

- **TokenHub → TokenRouter 炬枢** に名称変更（UI 文言 en/zh/ja、ページタイトル、サイドバー、ログインページ、パンくず）
- 新ブランドマーク: 角丸ヘキサゴン + 集束ビーム + アンバーコア（SVG/PNG/ICO）
- ブランドロックアップ再設計: 中国語「炬枢」を主役に（サイドバー 21px、ログイン 36px）、英語 TokenRouter を字間サブタイトルに
- バージョン表示をサイドバー下部へ移動、管理者のみ表示・アップグレード操作可能

### デザインと UX

- **デザインシステム**: `design-system/tokenrouter/MASTER.md`（ui-ux-pro-max）— プライマリ `#1E40AF`、背景 `#F8FAFC`、Fira Sans + Fira Code
- **SaaS スタイル全体アップグレード**: アンビエントブランドグロー、サイドバーグラデーション、グラデーションボタン、16px 角丸パネル、細スクロールバー、キーボードフォーカスリング、テーブル行ホバー、テーマ対応ステータスピル
- **React Bits モーション**: KPI カウントアップ（CountUp）、スポットライトカード（SpotlightCard）、グラデーション見出し（GradientText）、シャイマーボタン（ShinyText）、ビュー遷移、ナビ活性ピル（motion layoutId）
- **プレミアムチャート**: スムーズ Catmull-Rom トレンドライン（グラデーションストローク、ホバークロスヘア + ツールチップ、パルス終点ハロー）、SVG ドーナツ（12 時開始、ホバー連動凡例）、部門トークン棒グラフのホバー詳細、Top モデルのグラデーションプログレスバー
- **リクエストログ再設計**: 単行メタ（· セパレータ）、ブランドブルー選択状態、詳細パネルはビューポート内スクロール、1px 分割メタグリッド、ペイロード左右並列
- ナビゲーションのスタールコンテンツバグ修正（AnimatePresence mode="wait" → keyed motion.div）
- アクセシビリティ: `prefers-reduced-motion` フォールバック、可視フォーカス

### エンジニアリングとデモ

- `tools/mock-llm/server.js`: ゼロ依存 OpenAI 互換モック上流（ローカルで観測ページを照らす）
- エンドツーエンド検証: Provider → モデルインポート → ルーティング → API キー → `/v1` 呼び出し（ストリーミング含む）→ 使用量 / コスト / 監査
- リポジトリゲート全てグリーン: `next build` / `tsc` / `eslint` / 126 node テスト / 3 言語翻訳チェック / 環境変数契約 / 行数ベースライン

## 上流ベースライン

- **TokenHub v0.6.0**（astaxie/TokenHub、Apache-2.0）: OpenAI/Anthropic 互換 API、プロバイダーチャネルとテンプレートカタログ、モデルルーティングポリシー、プロジェクト単位キー、使用量/コスト帰属、ID ソースと RBAC、監査とアラート、SQLite/PostgreSQL。
