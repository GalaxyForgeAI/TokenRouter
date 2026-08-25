#!/usr/bin/env bash
# auto-demo.sh — TokenRouter 炬枢 auto 路由一键演示
#
# 发布一个旗舰档对外模型（若未发布），然后用同一个 model=auto 发两种请求：
#   1) 短 query  → 应选择最便宜的基础档模型
#   2) 长 prompt → 应自动升档到旗舰模型（质量约束生效）
#
# 用法:
#   ./tools/demo/auto-demo.sh [API_BASE] [ADMIN_PASSWORD]
# 默认: API_BASE=http://localhost:8090  ADMIN_PASSWORD=admin123456

set -euo pipefail

API_BASE="${1:-http://localhost:8090}"
ADMIN_PASSWORD="${2:-admin123456}"
IDENTITY="${ADMIN_IDENTITY:-admin}"

FLAGSHIP_MODEL="baidu/ernie-4.5-300b-a47b-paddle"
FLAGSHIP_PROVIDER_MODEL="pmdl_2c99f3c72fe0cc21c92a"
FLAGSHIP_PRICE_IN="1.5"
FLAGSHIP_PRICE_OUT="6"

say()  { printf '\033[1;34m== %s\033[0m\n' "$*"; }
ok()   { printf '\033[1;32m✓ %s\033[0m\n' "$*"; }
warn() { printf '\033[1;33m! %s\033[0m\n' "$*"; }

json_get() { python3 -c "import json,sys; d=json.load(sys.stdin); print(d$1)" 2>/dev/null || true; }

TOKEN=""
login() {
  say "登录 "${API_BASE}"（"${IDENTITY}"）"
  TOKEN=$(curl -fsS -m 10 -X POST "$API_BASE/api/admin/auth/login" \
    -H "content-type: application/json" \
    -d "{\"identity\":\"$IDENTITY\",\"password\":\"$ADMIN_PASSWORD\"}" | json_get "['token']")
  [ -n "$TOKEN" ] || { echo "登录失败" >&2; exit 1; }
  ok "登录成功"
}

ensure_flagship_model() {
  say "检查旗舰模型 $FLAGSHIP_MODEL"
  local exists
  exists=$(curl -fsS -m 10 -H "Authorization: Bearer $TOKEN" "$API_BASE/api/admin/models" | \
    python3 -c "import json,sys; d=json.load(sys.stdin); print(any(m.get('name')=='$FLAGSHIP_MODEL' for m in d.get('data',[])))")
  if [ "$exists" = "True" ]; then
    ok "旗舰模型已发布，跳过"
    return
  fi
  say "发布旗舰模型（含自动路由）"
  curl -fsS -m 20 -X POST "$API_BASE/api/admin/models" \
    -H "Authorization: Bearer $TOKEN" -H "content-type: application/json" \
    -d "{\"name\":\"$FLAGSHIP_MODEL\",\"family\":\"ernie\",\"modality\":\"chat\",\"context_window\":131072,\"input_price_usd_per_1m\":$FLAGSHIP_PRICE_IN,\"output_price_usd_per_1m\":$FLAGSHIP_PRICE_OUT,\"capabilities\":[\"chat\",\"tools\",\"structured_outputs\",\"reasoning\"],\"tier\":\"flagship\",\"status\":\"active\",\"initial_provider_models\":[\"$FLAGSHIP_PROVIDER_MODEL\"]}" >/dev/null
  ok "已发布 "${FLAGSHIP_MODEL}"（tier=flagship, 输入 \$$FLAGSHIP_PRICE_IN/1M）"
}

ensure_demo_key() {
  say "准备演示 API Key"
  local pid key
  pid=$(curl -fsS -m 10 -H "Authorization: Bearer $TOKEN" "$API_BASE/api/admin/projects" | \
    python3 -c "import json,sys; print(json.load(sys.stdin)['data'][0]['id'])")
  key=$(curl -fsS -m 10 -X POST -H "Authorization: Bearer $TOKEN" -H "content-type: application/json" \
    -d '{"name":"auto-demo-key"}' "$API_BASE/api/admin/projects/$pid/keys" | json_get "['api_key']")
  [ -n "$key" ] || { echo "创建 API Key 失败" >&2; exit 1; }
  ok "Key 已创建"
  DEMO_KEY="$key"
}

send_auto() {
  local label="$1" body="$2"
  local model
  model=$(printf '%s' "$body" | curl -fsS -m 60 -X POST "$API_BASE/v1/chat/completions" \
    -H "Authorization: Bearer $DEMO_KEY" -H "content-type: application/json" \
    --data-binary @- | json_get "['model']")
  printf '  %-28s → %s\n' "$label" "$model"
}

login
ensure_flagship_model
ensure_demo_key

say "auto 路由对比（同一 model=auto）"
send_auto "短 query（你好）" '{"model":"auto","messages":[{"role":"user","content":"你好"}]}'
LONG_TEXT=$(python3 -c "print('这是一段需要详细分析的长文本。' * 12000)")
LONG_JSON=$(python3 -c "
import json
text = '这是一段需要详细分析的长文本。' * 12000
print(json.dumps({'model': 'auto', 'messages': [{'role': 'user', 'content': text}]}))
")
printf '  长 prompt（%s 字符 ≈ %s tokens）\n' "${#LONG_TEXT}" "$(( ${#LONG_TEXT} / 4 ))"
send_auto "长 prompt" "$LONG_JSON"

say "完成 —— 短 query 落在最便宜的基础档，长 prompt 自动升档到旗舰（质量约束生效）"
