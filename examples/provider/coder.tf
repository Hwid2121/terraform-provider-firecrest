terraform {
  required_providers {
    coder = {
      source = "coder/coder"
    }
    firecrest = {
      source = "Hwid2121/firecrest"
      # version = "0.2.5"
    }
  }
}

data "coder_provisioner" "me" {}

data "coder_workspace" "me" {}

resource "coder_agent" "main" {
  arch           = data.coder_provisioner.me.arch
  os             = "linux"
  startup_script = <<-EOT
    set -e

    # Prepare user home with default files on first start.
    if [ ! -f ~/.init_done ]; then
      cp -rT /etc/skel ~
      touch ~/.init_done
    fi

    # Install and start code-server
    curl -fsSL https://code-server.dev/install.sh | sh -s -- --method=standalone --prefix=/tmp/code-server --version 4.19.1
    export PATH="/tmp/code-server/bin:$PATH"
    /tmp/code-server/bin/code-server --auth none --port 13337 >/tmp/code-server.log 2>&1 &
  EOT

  metadata {
    display_name = "CPU Usage"
    key          = "0_cpu_usage"
    script       = "coder stat cpu"
    interval     = 10
    timeout      = 1
  }

  metadata {
    display_name = "RAM Usage"
    key          = "1_ram_usage"
    script       = "coder stat mem"
    interval     = 10
    timeout      = 1
  }

}

# Runs a script at workspace start/stop or on a cron schedule
# details: https://registry.terraform.io/providers/coder/coder/latest/docs/resources/script
resource "coder_script" "startup_script" {
  agent_id           = coder_agent.main.id
  display_name       = "Startup Script"
  script             = <<-EOF
    #!/bin/sh
    set -e
    # Run programs at workspace startup
  EOF
  run_on_start       = true
  start_blocks_login = true
}


resource "coder_app" "code-server" {
  agent_id     = coder_agent.main.id
  slug         = "code-server"
  display_name = "code-server"
  url          = ""
  icon         = "/icon/code.svg"
  subdomain    = false
  share        = "owner"

  healthcheck {
    url       = ""
    interval  = 5
    threshold = 6
  }
}


module "vscode-web" {
  source         = "registry.coder.com/modules/vscode-web/coder"
  version        = "1.0.14"
  agent_id       = coder_agent.main.id
  accept_license = true
}




locals {
  job_env = <<-EOT
module load daint-gpu
module load sarus
node_ip=$(getent hosts $node_name | awk '{ print $1 }')


cd $SCRATCH/coder-testing/
srun sarus pull nikotaft/coder-environment:latest
srun sarus run nikotaft/coder-environment:latest /bin/bash -c "
  curl -fsSL https://code-server.dev/install.sh | sh -s -- --method=standalone --prefix=/tmp/code-server --version 4.19.1 &&
  /tmp/code-server/bin/code-server --auth none --bind-addr $node_ip:8080
"
EOT
}


provider "firecrest" {
  client_id     = "firecrest-ntafta-coder"
  client_secret = "D1wLfcA3BfVzxYA7eJ7AivIEklWNTH3C"
}


resource "firecrest_job" "job" {
  # job_script = ""
  client = "ntafta"
  job_name       = "coder-job"
  account        = "csstaff"
  email          = "nicolotafta@gmail.com"
  hours          = 0
  minutes        = 4
  nodes          = 1
  tasks_per_core = 1
  tasks_per_node = 1
  cpus_per_task  = 12
  partition      = "debug"
  constraint     = "gpu"
  executable     = local.job_env
  machine_name   = "daint"
  # env            = local.job_env
}




output "firecrest_job_id" {
  value = firecrest_job.job.id
}

output "firecrest_task_id" {
  value = firecrest_job.job.task_id
}

output "firecrest_node_ip" {
  value = firecrest_job.job.node_ip
}

