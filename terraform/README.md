# Deploy infra and related services with Terraform

This directory contains [Terraform](https://developer.hashicorp.com/terraform) modules
for delpoying the services in this repository. It is layed out as followed:

```txt
- INFRASTRUCTURE_PROVIDER (e.g. gcp, aws, azure, etc.)
    - main.tf, variables.tf, etc. = A single module to deploy everything from scratch
    - modules = composable modules you can use to pick and choose how/what to deploy
```

## Recommended: Install using top-level infrastructure provider module

Each infrastructure provider takes all of their respective sub-modules
and combines them into a single top-level module which can be used to deploy
all services and their infrastructure dependencies in one shot.

## Advanced: Pick and choose using composable modules

There are sub-modules in each infrastructure providers which are designed to be
composable. For example, let's say you already have a Google Cloud Storage bucket.
Then you might want to skip deploying a new bucket like the top-level module does.
In this case, you can leverage the composable sub-modules provided in "terraform/gcp"
to deploy everything except a new storage bucket and instead provide your existing
storage bucket.