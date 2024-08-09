terraform {
  required_providers {
    firecrest = {
      # source = "Hwid2121/firecrest"
      source = "registry.terraform.io/hashicorp/firecrest"

      # version = "0.2.6"
    }
    # coder = {
    #   source = "coder/coder"
    # }
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

  # client_id     = ""
  # client_secret = ""
  base_url       = "https://firecrest-tds.cscs.ch"
  token          = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJDSXhkdHZHeV92NEY4dXFXdlYzNFdndXRvM1BkN3h0RGk5dG9nMVV0MFIwIn0.eyJleHAiOjE3MjMyMTM4NDUsImlhdCI6MTcyMzIxMzU0NSwiYXV0aF90aW1lIjoxNzIzMjAyNTM4LCJqdGkiOiI4ZDIxZDM2ZS1kOWExLTQ3MmItOTEzYy02Y2E1ZTc5ODljOTciLCJpc3MiOiJodHRwczovL2F1dGgtdGRzLmNzY3MuY2gvYXV0aC9yZWFsbXMvY3NjcyIsImF1ZCI6WyJyZWFsbS1tYW5hZ2VtZW50IiwiYWNjb3VudCJdLCJzdWIiOiIwMDcyODRkZi0xNGUwLTRmMDktYTJmMC1jODczMjY5NTliMTAiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJjbG91ZC1pZGUiLCJzZXNzaW9uX3N0YXRlIjoiMzVlODUxYTEtMDU1Yy00MjEwLTk3NzAtOTAxM2IxZWQ4ZGU0IiwiYWNyIjoiMSIsInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJ3c28yYXBpbS1wdWJsaXNoZXIiLCJ3c28yYXBpbS1jcmVhdG9yIiwib2ZmbGluZV9hY2Nlc3MiLCJ1bWFfYXV0aG9yaXphdGlvbiIsImNzc3RhZmZfcm9sZSIsIndzbzJhcGltLXN1YnNjcmliZXIiLCJkZWZhdWx0LXJvbGVzLWNzY3MiXX0sInJlc291cmNlX2FjY2VzcyI6eyJyZWFsbS1tYW5hZ2VtZW50Ijp7InJvbGVzIjpbInZpZXctdXNlcnMiLCJxdWVyeS1ncm91cHMiLCJxdWVyeS11c2VycyJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJmaXJlY3Jlc3QgcHJvZmlsZSBlbWFpbCBmaXJlY3Jlc3QtdjIiLCJzaWQiOiIzNWU4NTFhMS0wNTVjLTQyMTAtOTc3MC05MDEzYjFlZDhkZTQiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwib3JnIjoiU3dpc3MgTmF0aW9uYWwgU3VwZXJjb21wdXRpbmcgQ2VudHJlIiwibmFtZSI6Ik5pY29sw7IgVGFmdGEiLCJncm91cHMiOlsid3NvMmFwaW0tcHVibGlzaGVyIiwid3NvMmFwaW0tY3JlYXRvciIsIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iLCJjc3N0YWZmX3JvbGUiLCJ3c28yYXBpbS1zdWJzY3JpYmVyIiwiZGVmYXVsdC1yb2xlcy1jc2NzIl0sInByZWZlcnJlZF91c2VybmFtZSI6Im50YWZ0YSIsImdpdmVuX25hbWUiOiJOaWNvbMOyIiwiZmFtaWx5X25hbWUiOiJUYWZ0YSIsImVtYWlsIjoibmljb2xvLnRhZnRhQGNzY3MuY2giLCJ1c2VybmFtZSI6Im50YWZ0YSJ9.Abwjb-LM-EsIz6cY3Hk7u-12J7Q5pAWXGNyFXHEvAFIdLPNUyRf8HIFb8lXn9znux2mWix1QpaeG2I5yeAa69vHYq_ieO_V9hP7MzaL6u9eYG58sTx76yjr39CSlHu6ZDJO20WgkCoeNhm6_Sq6w6RWB-buMOlUuVq-dxFMz068hCioG6e0jLLmaZLdUshorBULtbC_0_iwrIqgqbpvkiGI8Qm2weS-V0PnxMxgELfUB5mAvcN4u8e5aJZXJXyujFn8sjrD1ImoUR3ljBQdP6iIzFyYgOyCTP4IW6YCs8dwW1TB0hR-Tb52NXy7Y14H7G_QD7-e0O70x_LLS0gOeDQ"
  job_name       = "coder-job"
  account        = "csstaff"
  email          = "nicolotafta@gmail.com"
  hours          = 0
  minutes        = 5
  nodes          = 1
  tasks_per_core = 1
  tasks_per_node = 1
  cpus_per_task  = 6
  partition      = "normal"
  constraint     = "gpu"
  executable     = local.job_script
  machine_name   = "dom"
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
