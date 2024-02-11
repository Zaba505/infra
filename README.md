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
    imageService["Machine Management service"]
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

## Network Boot Procedure

```mermaid
sequenceDiagram
    box Home Infrastructure
        participant server
        participant router
        participant pi
    end

    box Public Infrastructure
        participant cloudflare
        participant gcp
    end

    server ->> router: DHCP
    router -->> server: PXE boot image uri
    server ->> pi: request boot image over FTP
    pi ->> cloudflare: proxy over HTTPS and mTLS
    cloudflare ->> gcp: proxy over mTLS
    gcp -->> cloudflare: bootstrap PXE to iPXE boot image
    cloudflare -->> pi: 
    pi -->> server: 
    server ->> cloudflare: chain load machine specific boot script over mTLS
    cloudflare ->> gcp: proxy over mTLS
    gcp -->> cloudflare: machine specific iPXE boot script
    cloudflare -->> server: 
```

## Services

### Boot Image Proxy service

This service runs on a Raspberry Pi in my home network and is responsible for proxying
FTP based boot image requests from the servers to HTTPS requests that are then sent
to the Machine Management service. This service also uses mTLS when reaching out to
the Machine Management service.

### Machine Management Service

This service is responsible for managing and returning iPXE scripts per machine. The scripts returned
by service are referred to as boot scripts becausethey are responsible for obtaining and loading all
the necessary files for booting an operating system on the specific machine its running on.

### Unknown Route Service

This service is only responsible for returning a 503 for any routes that the load balancer does not match
any other backend services to.
