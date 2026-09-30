---
collapse: true
---
## Selected Work

### One checkout, twelve ways to pay

*Payments · Storefront API · 2022–2026*

**Problem** Checkout is where every domain meets: cart, pricing, customer, stock and a dozen payment methods. It still has to feel like one call. A double-click must never create two orders, and a failed card may fall back to another processor, but never when that could authorise the same card twice.

**Role** Built the storefront's backend-for-frontend from its first code commit and have owned it since (~70% of commits over four years). Main contributor to the admin API (~62%), which runs on the service base I wrote. Designed the payment-gateway layer and the checkout orchestration.

**Decisions**

- **One interface for every payment method.** Each method (cards, wallets, buy-now-pay-later, wire, crypto) is a small Go adapter plugged in by id; there are 12 today. *Trade-off: a lowest-common-denominator interface, so a thin routing layer still lives in checkout.*
- **Fail over only when it is provably safe.** A second processor is tried only if the payment service answered, the decline was not final or fraud-related, no saved card was used, and at most once. *Trade-off: some recoverable declines are not retried; conversion is traded for never authorising twice.*
- **Best-effort checkout lock** per session against double-submits. *Trade-off: it fails open if the cache is down, putting availability ahead of strict de-duplication.*
- **Swapped the HTTP layer when it fought us.** The storefront API ran on a fasthttp-based framework, and a bug on its gRPC fan-out path could not be fixed fast enough, so I moved it to net/http (Gin): an 83-file migration on the live checkout path. *Trade-off: the admin API never hit the bug and stayed on the old stack, so the two gateways run different HTTP stacks.*
- **Thin order creation.** Persist the order and emit one event; emails, stock and notifications run asynchronously downstream.

**Outcome** 12 payment methods behind one four-method interface; the failover rule is covered by ~650 lines of tests.

```text
browser → storefront API
  → lock → validate → order (pending)
  → payment layer → processor A
      +- only if certainly not authorised
         → processor B
order event → async workers
  (email · stock · notifications)
```

### One chassis for ~40 Go services

*Platform · Framework · 2021–2026*

**Problem** About 40 Go services and Pub/Sub workers, run by a small team. Each needs the same lifecycle, drivers, validation, tracing and telemetry, and copying that 40 times produces 40 slightly different versions.

**Role** Original author and maintainer of the shared service framework since 2021 (~69% of 1,200+ commits). Maintainer of the gRPC contract layer.

**Decisions**

- **A service is `NewServiceApp(onInit, onClose)`.** Drivers, health checks and graceful shutdown come built in. *Trade-off: every upgrade is a fleet-wide rollout, and version skew is real.*
- **Contracts first.** ~95 gRPC services and ~1,100 RPCs live in one protobuf workspace, and ~390 declarative validation rules run by default in the framework's interceptor.
- **Tracing is a default, not a task.** I moved the fleet to the new tracer myself (38 of 42 repos, ~9 months), and context-aware logs carry the trace id.
- **Change-data-capture (2026).** Database changes stream through Pub/Sub into an append-only warehouse table, a scheduled merge keeps ~38 current-state tables, and the backfill is resumable.

**Outcome** 39 production deployables boot through it. Go 1.26 rolled out across 38 repos in ~9 days.

```text
service = chassis + domain logic
+--------------- chassis ---------------+
| lifecycle · drivers · health          |
| graceful shutdown                     |
| gRPC: recover → validate → trace      |
+---------------------------------------+
  ^ contracts: ~95 services · ~1,100 RPCs
```

### Letting AI agents run commerce operations, safely

*AI agents · MCP · 2026*

**Problem** Support bots and internal agents needed to look up orders, shipments and payments, and sometimes act on them, without anyone clicking through the admin console and without handing an LLM a master key.

**Role** Designed and built the MCP (Model Context Protocol) gateway inside the admin backend: six servers, ~99% of commits.

**Decisions**

- **Reuse business logic, don't bypass it.** Tools go through the existing use cases and gRPC services.
- **Permission tiers per key**, with a deny-by-default allowlist for support bots.
- **Preview, then commit.** An order is created from a previewed cart. There are no standalone refund tools; refunds happen only through a server-validated cancellation.
- **Warehouse queries** are dry-run first, SELECT-only (fail-closed), and capped on bytes and rows.

**Outcome** ~114 typed tools across orders, shipments, payments, catalog, content, helpdesk and analytics, with 114 unit tests.

```text
agent / support bot → MCP → gateway
  auth → tier → allowlist → redacted log
  +- orders · shipments · payments
  |    preview → commit
  +- catalog · content · helpdesk
  +- warehouse: dry-run, SELECT-only, capped
  same use cases & services as admin API
```

### Agents as teammates

*AI operations · Claude Code · 2026*

**Problem** Seven teams (growth, customer support, operations, bulk sales, a product line, finance, HR) each had an AI assistant on a third-party agent runtime, and engineers were handing real work to a coding agent across ~50 repositories. Both needed clear rules, the right tools, and a safe way to change them.

**Role** I run the assistant fleet's operations (hosting, upgrades, config management, model routing, runbooks) and I'm the main author of two of the seven assistants: internal analytics and customer-facing support. For engineering, I wrote the Claude Code setup: an architecture map, three workflow skills and a guard hook. I didn't write the agent runtime, and the other five assistants were started by their teams.

**Decisions**

- **Behaviour lives in git.** Each assistant's persona, rules, skills and approved SQL are a versioned repo; live config changes are diffed against the running copy before they land. *Trade-off: the runtime rewrites its own config, so there are two sources to reconcile.*
- **Tools come through the MCP gateway, wired per bot.** Each assistant gets only the servers its job needs. The customer assistant adds a per-tool allowlist and is gated like a support rep: order number plus email to look anything up, a one-time code and explicit confirmation for changes, and no refund tools in chat.
- **Measure before answering.** The analytics assistant has one canonical metric-definitions file, a trusted-table check and a regression eval set; every wrong answer becomes a documented anti-pattern.
- **Hard rules in a hook, judgment in skills.** A guard checks every shell command the coding agent runs (7 git and formatting rules; it denies and says what to do instead). Staging deploys, production releases and contract changes are skills, and "deployed" means the GitOps commit names your tag, not a green pipeline.

**Outcome** Seven assistants serving six internal teams plus shoppers, on one host, operated from one inventory and a ~540-line runbook. The first state-migrating upgrade took ~12 s of downtime plus ~50 s of migration, with a full rollback archive taken first. The analytics assistant's eval went from 2/10 to 7/10 in one evening of fixes. Coding-agent sessions start with the team's rules, 3 skills and a 7-rule guard.

```text
team chats · storefront widget
  → 7 assistants, one container each
      persona · rules · skills · SQL (git)
  → MCP gateway: tools wired per bot
  → in-house model relay

engineer → Claude Code
  reads the map · loads skills
  every command → guard (7 rules)
  production tag → human confirms
```
