terraform {
  required_providers {
    firecrest = {
      source = "Hwid2121/firecrest"
    # source = "registry.terraform.io/hashicorp/firecrest"

    #   version = "0.2.5"
    }
  }
}


locals {
  job_env = <<-EOT
module load daint-gpu
module load sarus

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
  minutes        = 2
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


