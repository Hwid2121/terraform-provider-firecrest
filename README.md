# Firecrest Terraform Provider

<!-- ![Firecrest Logo](link-to-logo-image) -->

The Firecrest Terraform Provider allows you to manage and interact with the Firecrest API using Terraform. Firecrest is a powerful tool for managing computational resources on HPC systems. This provider simplifies the automation and management of these resources using Terraform's infrastructure-as-code approach.

## Table of Contents

- [Firecrest Terraform Provider](#firecrest-terraform-provider)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)

## Features

- **Job Management**: Submit, monitor, and manage computational jobs.
- **Resource Management**: Manage computational resources and environments.
- **Integration with Coder**: Seamless integration with Coder for managing cloud IDEs.

## Installation

To install the Firecrest Terraform Provider, add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    firecrest = {
      source  = "Hwid2121/firecrest"
    }
  }
}
```

## Usage
Provider Configuration
First, configure the Firecrest provider with the necessary credentials:

```hcl
Copy code
provider "firecrest" {
  client_id     = var.client_id
  client_secret = var.client_secret
  client_token  = data.coder_external_auth.keycloak.access_token
}
```

Make sure to replace var.client_id, var.client_secret, and data.coder_external_auth.keycloak.access_token with your actual credentials or data sources.

# Resources
The Firecrest provider currently supports the following resources:

- firecrest_job: Manages computational jobs.

Here is an example of how to use the Firecrest provider to submit a job:

```hcl
locals {
  job_script = <<-EOT
    #!/bin/bash
    echo "Hello, Firecrest!"
  EOT
}

resource "firecrest_job" "example" {
  base_url     = "https://firecrest.cscs.ch"
  token        = var.token
  job_name     = "example-job"
  account      = "example-account"
  machine_name = "daint"
  executable   = local.job_script
  hours        = 1
  minutes      = 0
  nodes        = 1
  tasks_per_core = 1
  tasks_per_node = 1
  cpus_per_task  = 1
  partition    = "normal"
}
```



Please note: We take Terraform's security and our users' trust very seriously. If you believe you have found a security issue in the Terraform AWS Provider, please responsibly disclose it by contacting us at security@hashicorp.com.