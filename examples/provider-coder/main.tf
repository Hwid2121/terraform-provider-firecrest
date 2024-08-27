terraform {
  required_providers {
    firecrest = {
      source = "registry.terraform.io/hashicorp/firecrest"
      # version = "0.1.0"
    }
  }
}

provider "firecrest" {
}



locals {
  job_script = <<-EOT

    # module load daint-gpu
    module load sarus

    node_name=$(scontrol show hostname $SLURM_JOB_NODELIST)
    node_ip=$(getent hosts $node_name | awk '{ print $1 }')

    echo "Node name: $node_name"
    echo "Node IP: $node_ip"
    sleep 120
  EOT
}

resource "firecrest_job" "job" {
  client_id     = "firecrest-ntafta-coder"
  client_secret = "D1wLfcA3BfVzxYA7eJ7AivIEklWNTH3C"
  # client_id     = ""
  # client_secret = ""
  base_url       = "https://api.cscs.ch/hpc/firecrest/v1"
  token         = ""
  job_name       = "coder-job"
  account        = "csstaff"
  email          = "nicolotafta@gmail.com"
  hours          = 0
  minutes        = 3
  nodes          = 1                            
  tasks_per_core = 1
  tasks_per_node = 1
  cpus_per_task  = 6
  partition      = "normal"
  constraint     = "gpu"
  executable     = local.job_script
  machine_name   = "eiger"
}

# data "coder_provisioner" "me" {}

# data "coder_workspace" "me" {}


# resource "coder_agent" "main" {
#   os             = "linux"
#   arch           = "amd64"

#   env = {
#     KC_TOKEN : data.coder_external_auth.keycloak.access_token
#   }

#   startup_script = <<-EOT
# EOT

  
# }


# data "coder_external_auth" "keycloak" {
#   # Matches the ID of the external auth provider in Coder.
#   id = "keycloak"
# }


# resource "coder_app" "code-server" {
#   agent_id     = coder_agent.main.id
#   slug         = "code-server"
#   display_name = "code-server"
#   icon         = "${data.coder_workspace.me.access_url}/icon/code.svg"
#   url          = "http://localhost:8080"
#   # share        = "owner"
#   # subdomain    = false

#   healthcheck {
#     url       = "http://localhost:8080/healthz"
#     interval  = 5
#     threshold = 6
#   }
# }
