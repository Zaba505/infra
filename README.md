# My Infrastructure as Code Monorepo

```mermaid
flowchart LR

    dns["DNS"]
    https["HTTP(S)"]
    cloudflare["Cloudflare"]
    homeRouter["Router"]
    homeLab["Home Lab Cluster"]
    k8s["Kubernetes"]
    loadBalancer["Load balancer"]
    bootImages[("Machine boot images")]
    imageService["Machine Image service"]
    sinkService["Unknown Route service"]

    subgraph a [Public Internet]
        dns --> cloudflare
        https --> cloudflare
    end

    subgraph one [Home]
        homeRouter --> homeLab
        homeLab --> k8s
    end

    subgraph two [Cloud]
        loadBalancer --> imageService
        loadBalancer --> sinkService
        imageService --> bootImages
    end

    cloudflare -->|mTLS| loadBalancer
    cloudflare -->|mTLS| homeRouter
```