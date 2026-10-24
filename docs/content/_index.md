---
title: Overview
type: docs
no_list: true
---

infra defines a self-hosted application platform: the place the operator's own
capabilities run by default, so that none of them depends on a vendor's hosting
product.

- **It's rebuilt, not repaired.** The platform can be stood up from nothing, from
  this repository, in at most 60 minutes.
- **It's self-hosted, not self-owned.** It spans a public cloud and the home lab.
  "Self-hosted" means the operator controls it end to end, not that every part
  runs on hardware the operator owns.
- **It's run by one operator,** in at most 2 hours a week. A designated successor
  holds sealed credentials and uses them only if the operator can't.
- **Its tenants are the operator's capabilities.** Their end users reach the
  tenants, never the platform.

The platform is being rebuilt from scratch, as the [charter](https://github.com/Zaba505/infra/discussions/970) sets out.

## Reading order

The site reads in the order its readers need it: the operator first, then a
successor who has to take over, then a capability owner bringing a tenant.

1. [Getting started]({{< ref "/getting-started/_index.md" >}}): standing the
   platform up from nothing, first as a drill and then live.
2. [Fundamentals]({{< ref "/fundamentals/_index.md" >}}): the ideas every part
   of the platform builds on.
3. [The platform]({{< ref "/platform/_index.md" >}}): one page per part, in the
   order a rebuild stands them up.
4. [Operating]({{< ref "/operating/_index.md" >}}): running the platform day to
   day, handing it to a successor, and recovering it.
5. [Hosting a tenant]({{< ref "/hosting-a-tenant/_index.md" >}}): what a
   capability owner needs to run a tenant on the platform.

Start with [Getting started]({{< ref "/getting-started/_index.md" >}}).
