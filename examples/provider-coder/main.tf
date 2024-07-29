terraform {
  required_providers {
    firecrest = {
      # source = "Hwid2121/firecrest"
      source = "registry.terraform.io/hashicorp/firecrest"

      # version = "0.2.6"
    }
    coder = {
      source = "coder/coder"
    }
  }
}

data "coder_external_auth" "keycloak" {
  # Matches the ID of the external auth provider in Coder.
  id = "keycloak"
}

provider "firecrest" {
  client_id     = ""
  client_secret = ""

  # client_id     = ""
  # client_secret = ""
  # client_token = data.coder_external_auth.keycloak.access_token
}

locals {
  job_script = <<-EOT
    cat - > agent.sh <<< '${coder_agent.main.init_script}'
    chmod +x agent.sh

    module load daint-gpu
    module load sarus

    node_name=$(scontrol show hostname $SLURM_JOB_NODELIST)
    node_ip=$(getent hosts $node_name | awk '{ print $1 }')

    echo "Node name: $node_name"
    echo "Node IP: $node_ip"

    export CODER_AGENT_TOKEN=${coder_agent.main.token}
    export CODER_AGENT_ID=${coder_agent.main.id}

    echo "Coder Token: $CODER_AGENT_TOKEN "
    echo "Coder ID: $CODER_AGENT_ID"


    srun sarus pull nikotaft/coder-environment:latest
    srun sarus run nikotaft/coder-environment:latest /bin/bash -c "
      curl -fsSL https://code-server.dev/install.sh | sh -s -- --method=standalone --prefix=/tmp/code-server --edge &&
env SHELL=/bin/bash HOME=/home/coder PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin /tmp/code-server/bin/code-server --auth none --port 8080 --log debug    " &
     
     ./agent.sh
     echo "Agent is not blocking..."
     sleep 3600
  EOT
}

resource "firecrest_job" "job" {
  token          = "TOKEN-TESTING"
  base_url       = "https://firecrest-tds.cscs.ch"
  client         = ""
  job_name       = ""
  account        = ""
  email          = ""
  hours          = 0
  minutes        = 8
  nodes          = 1
  tasks_per_core = 1
  tasks_per_node = 1
  cpus_per_task  = 6
  partition      = "debug"
  constraint     = "gpu"
  executable     = local.job_script
  machine_name   = "daint"
}

data "coder_provisioner" "me" {}

data "coder_workspace" "me" {}


resource "coder_agent" "main" {
  os   = "linux"
  arch = "amd64"
  env = {
    KC_TOKEN : data.coder_external_auth.keycloak.access_token
  }
  startup_script = <<-EOT
EOT


}

resource "coder_app" "code-server" {
  agent_id     = coder_agent.main.id
  slug         = "code-server"
  display_name = "code-server"
  icon         = "${data.coder_workspace.me.access_url}/icon/code.svg"
  url          = "http://localhost:8080"
  # share        = "owner"
  # subdomain    = false

  healthcheck {
    url       = "http://localhost:8080/healthz"
    interval  = 5
    threshold = 6
  }
}


output "firecrest_node_ip" {
  value = firecrest_job.job.node_ip
}

output "firecrest_base_url" {
  value = firecrest_job.job.base_url
}

output "firecrest_token" {
  value = firecrest_job.job.token
}

# output "firecrest_node_ip" {
#   value = firecrest_job.job.node_ip
# }

# output "keycloak_access_token" {
#   value = data.coder_external_auth.keycloak.access_token
# }
