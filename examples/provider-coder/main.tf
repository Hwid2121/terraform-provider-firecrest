terraform {
  required_providers {
    firecrest = {
      # source = "Hwid2121/firecrest"
      source = "registry.terraform.io/hashicorp/firecrest"
      # version = "0.2.6"
    }
  }
}



provider "firecrest" {
}

locals {
  job_script = <<-EOT
  chmod +x agent.sh

  module load daint-gpu
  module load sarus

  node_name=$(scontrol show hostname $SLURM_JOB_NODELIST)
  node_ip=$(getent hosts $node_name | awk '{ print $1 }')

  echo "Node name: $node_name"
  echo "Node IP: $node_ip"
  echo "Coder Token: $CODER_AGENT_TOKEN "
  echo "Coder ID: $CODER_AGENT_ID"


  sleep 180
  EOT
}

resource "firecrest_job" "job" {

}

# data "coder_provisioner" "me" {}

# data "coder_workspace" "me" {}


# resource "coder_agent" "main" {
#   os   = "linux"
#   arch = "amd64"
#   env = {
#     KC_TOKEN : data.coder_external_auth.keycloak.access_token
#   }
#   startup_script = <<-EOT
#   EOT
# }

# resource "coder_app" "code-server" {
#   agent_id     = coder_agent.main.id
#   slug         = "code-server"
#   display_name = "code-server"
#   icon         = "${data.coder_workspace.me.access_url}/icon/code.svg"
#   url          = "http://localhost:8080"

#   healthcheck {
#     url       = "http://localhost:8080/healthz"
#     interval  = 5
#     threshold = 6
#   }
# }

