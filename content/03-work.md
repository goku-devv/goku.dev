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

### One sign-in for every app

*Identity · OAuth2 / OIDC · 2026*

**Problem** Every storefront, admin and command-line client needed one standards-based way to sign a person in, and people needed to see, and end, the sessions their devices hold.

**Role** Designed and built the platform's single sign-on service from an empty repo (263 of 266 commits), in production since January 2026, and its TypeScript client SDK.

**Decisions**

- **One login path, pluggable grants.** Password, email code, TOTP, passkeys, Google, Apple, GitHub and refresh all plug into one pipeline; sessions, cookies and token issuance are written once. *Trade-off: every grant has to fit one request shape.*
- **Revocable per device, checked at refresh.** Each sign-in is a grant whose id rides in both tokens; ending it from a Devices & apps page, an admin tool or RFC 7009 revocation takes effect at the next refresh. *Trade-off: revocation lands within one access-token lifetime, not instantly, so the API hot path stays lookup-free.*
- **Credentials live with the service that checks them.** Passkeys are discoverable, so starting a sign-in reveals no account, and one-time-code and TOTP secrets sit in the auth service's own store, encrypted at rest. *Trade-off: passkeys are sold as faster, not stronger, while email codes remain for the same account.*
- **Ship dark, then enforce.** The client and redirect policy rolled out off → log → enforce per environment, so real traffic proved it before it refused anything.
- **Asymmetric signing without a forced logout.** Tokens moved from a shared secret to RS256; downstream services hold only the public key, and both algorithms were accepted during the switch.

**Outcome** Standards-based (PKCE, RFC 7009, RFC 9207, RFC 8252 loopback redirects), 7 grant types, 232 Go tests, and a zero-dependency SDK on npm for browser, React and CLI sign-in.

```text
storefront · admin · CLI
  → one login path
      password · code · TOTP · passkey
      Google · Apple · GitHub · refresh
  → one grant per device, id in both tokens
  → RS256 tokens; services hold the
    public key only
revoke → takes effect at the next refresh
```

### One tracking record per parcel

*Fulfillment · Integrations · 2022–2026*

**Problem** Parcels ship through parcel carriers, LTL freight carriers and vendors' own store platforms, and only some of them push updates. The shop needs one tracking record per shipment, delivery has to move the order forward, returns need tracking too, stuck parcels have to surface without anyone checking by hand, and pushed and polled updates must never overwrite each other.

**Role** Built the tracking service and the webhook ingest from their first commits in 2022. Took over the multi-carrier poller in late 2024, rebuilt its engine and wrote ~85% of its commits since. The carrier API clients and the first poller were a teammate's.

**Decisions**

- **One path per shipment.** Parcels labelled through the webhook-capable provider are updated only by webhooks; everything else is polled. *Trade-off: two code paths, and if webhooks stop arriving the poller does not fill the gap.*
- **Store first, then queue.** The ingest saves every raw event before publishing it, so the carrier gets a fast answer and there is always a replay trail; the consumer acknowledges only after processing. *Trade-off: no de-duplication on event id, so downstream writes must be idempotent.*
- **One system of record.** The tracking service writes with compare-and-swap (a stale update is rejected), keeps an audit row and emits a status event that moves the order forward, without knowing anything about orders. *Trade-off: a writer that loses the race re-reads and retries.*
- **One lane per carrier account.** A slow or rate-limited carrier holds up only its own lane. The rewrite gave each run its own state and folded forward and return tracking into one engine (+637 / −753 lines). *Trade-off: one request in flight per account, so throughput is capped on purpose.*
- **People confirm stock.** Inbound container arrivals are recorded and ops are alerted, but stock is not booked automatically: a delivery scan means the truck arrived, not that anyone counted the goods. Auto-confirm was built, then switched off the next day. *Trade-off: a person confirms every arrival.*

**Outcome** Eight carrier and store-platform adapters across ten lanes, including the two LTL freight adapters I wrote. Stuck parcels are detected and reported automatically, and returns update the shipment and notify the vendor. Built and run over four years.

```text
carrier webhook → ingest
  store raw event → queue → worker
scheduled poller (non-webhook parcels)
  one lane per carrier account
  both paths write to:
tracking service: CAS write · audit
  → status event → order moves forward
stuck parcels → ops report
```
