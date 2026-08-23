import type { Model } from "../core/types";

// Model profile radar axes: each axis scores 0-100 and carries a short
// evidence string (capability tags / derived hints) shown in the tooltip.
export type ProfileAxis = {
  key: string;
  label: string;
  score: number;
  evidence: string;
};

const clampScore = (value: number) => Math.max(0, Math.min(100, Math.round(value)));

function tierBoost(tier?: string): number {
  return { basic: 0, standard: 5, flagship: 25, frontier: 40 }[tier ?? ""] ?? 5;
}

export function modelProfileAxes(model: Model): ProfileAxis[] {
  const caps = new Set(model.capabilities ?? []);
  const name = model.name.toLowerCase();
  const tier = model.tier ?? "";
  const ctx = model.context_window ?? 0;
  const price = model.input_price_usd_per_1m ?? 0;
  const boost = tierBoost(tier);

  const evidence = (tags: string[]) => tags.filter((tag) => caps.has(tag)).join(" / ") || "-";

  const language = clampScore(
    55 + (caps.has("chat") ? 15 : 0) + (caps.has("structured_outputs") ? 20 : 0) + (caps.has("tools") ? 10 : 0),
  );
  const multimodal = clampScore(
    (caps.has("vision") ? 35 : 0) +
      (caps.has("video_input") ? 25 : 0) +
      (caps.has("audio_input") ? 15 : 0) +
      (caps.has("audio_output") ? 15 : 0) +
      (caps.has("image_output") ? 10 : 0),
  );
  const reasoning = clampScore(
    (caps.has("reasoning") ? 45 : 0) +
      boost +
      (name.includes("deep") || name.includes("think") || name.includes("r1") ? 15 : 0),
  );
  const agentic = clampScore(
    (caps.has("tools") ? 40 : 0) + (caps.has("reasoning") ? 25 : 0) + (ctx >= 128000 ? 20 : 0) + (caps.has("structured_outputs") ? 15 : 0),
  );
  const coding = clampScore(
    (name.includes("codex") ? 50 : 0) +
      (name.includes("code") || name.includes("coder") || name.includes("dev") || name.includes("engineer") ? 45 : 0) +
      (caps.has("reasoning") ? 20 : 0) +
      (caps.has("tools") ? 15 : 0) +
      (boost >= 25 ? 10 : 0),
  );
  const longContext =
    ctx >= 1000000 ? 100 : ctx >= 400000 ? 90 : ctx >= 200000 ? 80 : ctx >= 128000 ? 65 : ctx >= 64000 ? 50 : ctx >= 32000 ? 35 : ctx >= 8000 ? 20 : 10;
  const costEfficiency = price <= 0 ? 45 : price < 0.1 ? 95 : price < 1 ? 80 : price < 3 ? 65 : price < 10 ? 45 : price < 30 ? 25 : 10;
  const speed = clampScore(
    ({ basic: 90, standard: 70, flagship: 45, frontier: 30 }[tier] ?? 60) +
      (name.includes("flash") || name.includes("nano") || name.includes("turbo") || name.includes("lite") ? 10 : 0),
  );

  return [
    { key: "language", label: "语言理解", score: language, evidence: evidence(["chat", "structured_outputs", "tools"]) },
    { key: "multimodal", label: "多模态", score: multimodal, evidence: evidence(["vision", "video_input", "audio_input", "audio_output", "image_output"]) },
    { key: "reasoning", label: "推理深度", score: reasoning, evidence: evidence(["reasoning"]) || name },
    { key: "agentic", label: "Agent 能力", score: agentic, evidence: evidence(["tools", "reasoning", "structured_outputs"]) },
    { key: "coding", label: "编码能力", score: coding, evidence: evidence(["reasoning", "tools"]) || name },
    { key: "long_context", label: "长上下文", score: longContext, evidence: ctx > 0 ? `${Math.round(ctx / 1000)}k` : "-" },
    { key: "cost_efficiency", label: "成本效率", score: costEfficiency, evidence: price > 0 ? `$${price}/1M in` : "-" },
    { key: "speed", label: "响应速度", score: speed, evidence: tier || "-" },
  ];
}
