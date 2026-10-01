---
title: Experience
layout: timeline
entries:
  - company: Autonomous Inc
    dates: 2020 – Present
    role: Senior Software Engineer
    url: https://autonomous.ai
    bullets:
      - "Built and own much of the commerce backend — the Go service framework, the storefront and admin APIs, checkout and payments, and the MCP gateway for AI agents."
      - "Owned the pricing and promotion service (2022–2026, ~71% of commits) and turned its flat list of promotion conditions into an AND/OR rule tree with pluggable checks."
      - "Rebuilt the content service (posts, help center, SEO data, short links) from an empty repo in 2023; added an LLM step that routes contact-form messages to sales, support or spam, failing open so a real buyer is never dropped, checked against a golden set."
      - "Rewrote the notification worker around one Pub/Sub dispatcher, replacing a workflow engine it didn't need; made transactional email duplicate-safe under at-least-once delivery; built the Redis-backed WebSocket broker that reaches a user on whichever pod holds the connection."
      - "Made database reads go to replicas by default, with explicit primary reads for read-after-write paths (53 call sites); moved 22 repos onto the new MongoDB driver."
      - "Made the store-credit, loyalty-points and affiliate ledger transactional: every balance change and its history entry commit together behind an optimistic lock."
      - "Moved the admin console's sign-in onto the single sign-on service and SDK I built, with session, device and login-history views."
      - "Built and owned the platform's Elasticsearch product search (2022–2026, ~78% of its core commits): an indexing worker fed by catalog-change events, plus listing filters and attribute facets that admins add without a deploy; a teammate did most of the early relevance tuning."
      - "Built tiered flash sales in the catalog (2023–24): minute-aligned windows, automatic tier advance and next-sale chaining, with storefront, admin and import layers and the pricing hook that ties each sale line to a teammate's stock counter at checkout."
  - company: Bestarion Software Company Ltd
    dates: 2020
    role: Senior Software Engineer
  - company: WeVenture Pte Ltd
    dates: 2018
    role: Software Engineer
  - company: Gumi Vietnam
    dates: 2016
    role: Software Engineer
---
