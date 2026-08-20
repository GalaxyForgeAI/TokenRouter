// TokenRouter mock OpenAI-compatible upstream — zero dependency.
// 关键：流式末 chunk 必须带 usage 字段，网关才会写 UsageLog。
const http = require("http");
const port = process.env.PORT || 9999;

http
  .createServer((req, res) => {
    if (req.method === "POST" && req.url.startsWith("/v1/chat/completions")) {
      let body = "";
      req.on("data", (c) => (body += c));
      req.on("end", () => {
        let parsed = {};
        try {
          parsed = JSON.parse(body);
        } catch {}
        const last = (parsed.messages || []).slice(-1)[0];
        const prompt = (last && last.content) || "hi";
        const promptTokens = Math.max(4, Math.ceil(String(prompt).length / 3));
        const completionTokens = 8 + Math.floor(Math.random() * 40);
        const model = parsed.model || "tr-mock-1";
        const created = Math.floor(Date.now() / 1000);

        if (parsed.stream === true) {
          res.writeHead(200, {
            "content-type": "text/event-stream",
            "cache-control": "no-cache",
            connection: "keep-alive",
          });
          const words = ["TokenRouter", "炬枢", "mock", "upstream", "gateway", "routed", "✅"].slice(
            0,
            Math.max(1, Math.ceil(completionTokens / 12)),
          );
          words.forEach((w) => {
            res.write(
              `data: ${JSON.stringify({
                id: "chatcmpl-tr-mock",
                object: "chat.completion.chunk",
                created,
                model,
                choices: [{ index: 0, delta: { content: w + " " }, finish_reason: null }],
              })}\n\n`,
            );
          });
          res.write(
            `data: ${JSON.stringify({
              id: "chatcmpl-tr-mock",
              object: "chat.completion.chunk",
              created,
              model,
              choices: [{ index: 0, delta: {}, finish_reason: "stop" }],
              usage: {
                prompt_tokens: promptTokens,
                completion_tokens: completionTokens,
                total_tokens: promptTokens + completionTokens,
              },
            })}\n\n`,
          );
          res.write("data: [DONE]\n\n");
          res.end();
          return;
        }

        res.writeHead(200, { "content-type": "application/json" });
        res.end(
          JSON.stringify({
            id: "chatcmpl-tr-mock",
            object: "chat.completion",
            created,
            model,
            choices: [
              {
                index: 0,
                message: { role: "assistant", content: "TokenRouter 炬枢 mock upstream responding." },
                finish_reason: "stop",
              },
            ],
            usage: {
              prompt_tokens: promptTokens,
              completion_tokens: completionTokens,
              total_tokens: promptTokens + completionTokens,
            },
          }),
        );
      });
      return;
    }

    if (req.method === "GET" && req.url.startsWith("/v1/models")) {
      res.writeHead(200, { "content-type": "application/json" });
      res.end(
        JSON.stringify({
          object: "list",
          data: [{ id: "tr-mock-1", object: "model", owned_by: "tokenrouter" }],
        }),
      );
      return;
    }

    res.writeHead(404, { "content-type": "text/plain" });
    res.end("not found");
  })
  .listen(port, () => console.log(`mock LLM upstream listening on :${port}`));
