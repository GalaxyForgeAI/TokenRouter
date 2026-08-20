import { defineConfig } from "vitepress";

export default defineConfig({
  title: "TokenRouter 炬枢",
  description: "炬枢 · 企业私有化 AI 网关 —— 统一接入、策略治理、用量归因、全程审计",
  lang: "zh-CN",
  cleanUrls: true,
  outDir: "../dist",
  ignoreDeadLinks: "localhost",
  head: [
    ["link", { rel: "icon", type: "image/png", href: "/brand/tokenrouter-logo.png" }],
  ],
  themeConfig: {
    logo: "/brand/tokenrouter-logo.png",
    siteTitle: "TokenRouter 炬枢",
    nav: [
      { text: "首页", link: "/" },
      { text: "产品文档", link: "/guide/product-overview" },
      { text: "快速开始", link: "/guide/quickstart" },
      {
        text: "GitHub",
        link: "https://github.com/GalaxyForgeAI/TokenRouter",
      },
    ],
    sidebar: [
      {
        text: "项目文档",
        items: [
          { text: "产品概述", link: "/guide/product-overview" },
          { text: "快速开始", link: "/guide/quickstart" },
          { text: "品牌与许可", link: "/guide/branding-and-license" },
          { text: "变更记录", link: "/guide/changelog" },
        ],
      },
    ],
    footer: {
      message: "Apache-2.0 Licensed · 上游致谢 astaxie/TokenHub",
      copyright: "Copyright © 2026 TokenRouter 炬枢",
    },
    search: {
      provider: "local",
    },
  },
});
