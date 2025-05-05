---
title: Network Boot
type: docs
---

{{% alert title="Note" color="info" %}}
Network booting only applies to following machines in the homelab:
- [hp-dl360]({{% ref "/physical_infrastructure/#servers" %}})
{{% /alert %}}

## Architecture

```mermaid
flowchart LR
    subgraph homelab [Homelab]
    server@{ label: "Server", shape: rect }
    switch@{ label: "Switch", shape: circle }
    router@{ label: "Router", shape: circle }
    end

    subgraph cloud [Public Cloud]
    wiregaurd@{ label: "Wiregaurd", shape: rect }
    matchbox@{ label: "Matchbox", shape: rect }
    end

    server --- switch
    switch --- router
    router --- wiregaurd
    wiregaurd --- matchbox
```

## Homelab

### Setting up Network Boot VLAN

We begin by creating a VLAN on the router and assign by [MAC address](https://en.wikipedia.org/wiki/MAC_address)
the first port of each server to it.

### Initializing Wiregaurd node

One feature of the [Router]({{% ref "/physical_infrastructure/#router" %}}) is that
it can act as either a [Wiregaurd](https://www.wireguard.com/) server or client.
In this instance, the client implementation is initialized and all traffic on
the Network Boot VLAN will be assigned to go through the [Wiregaurd](https://www.wireguard.com/)
connection.

## Public Cloud

```mermaid
flowchart LR
    gateway@{ label: "External Passthrough Network Load Balancer", shape: rect }

    subgraph k8s [Kubernetes]
    wiregaurd@{ label: "Wiregaurd", shape: rect }
    matchbox@{ label: "Matchbox", shape: rect }
    end

    gateway --- wiregaurd
    wiregaurd --- matchbox
```

Within a public cloud (e.g. [AWS](https://aws.amazon.com/), [Azure](https://azure.microsoft.com/), [GCP](https://cloud.google.com/)),
a [Kubernetes](https://kubernetes.io/) cluster will be instantiated.

### Kubernetes Architecture

```mermaid
flowchart LR
    subgraph ns [Network Boot Namespace]
    direction LR
    wiregaurd_service@{ label: "Wiregaurd Service (LoadBalancer)", shape: rect }
    wiregaurd@{ label: "Wiregaurd Daemon Set", shape: rect }
    wiregaurd_config@{ label: "Wiregaurd Config Secret", shape: cyl }
    matchbox@{ label: "Matchbox Deployment", shape: rect }
    matchbox_config@{ label: "Matchbox Config Volume", shape: cyl }
    end

    wiregaurd_service --- wiregaurd
    wiregaurd --- wiregaurd_config
    wiregaurd --- matchbox
    matchbox --- matchbox_config
```

## Boot Sequence

```mermaid
sequenceDiagram
    participant server as Server
    participant router as Router
    participant wiregaurd as Wiregaurd
    participant matchbox as Matchbox

    server ->> router: UEFI HTTP boot request
    router ->> wiregaurd: Encrypted UEFI HTTP boot request
    wiregaurd ->> matchbox: UEFI HTTP boot request
    matchbox ->> wiregaurd: UEFI Executable Image
    wiregaurd ->> router: UEFI Executable Image
    router ->> server: UEFI Executable Image
```