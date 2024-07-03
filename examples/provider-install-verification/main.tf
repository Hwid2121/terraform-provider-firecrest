terraform {
  required_providers {
    firecrest = {
      source = "registry.terraform.io/hashicorp/firecrest"
      # version = "0.1.0"
    }
  }
}

provider "firecrest" {
  client_id     = var.client_id
  client_secret = var.client_secret
}


resource "firecrest_job" "job" {
  job_name       = "job-test"
  account        = var.account_name
  email          = var.account_email
  hours          = 0
  minutes        = 30
  nodes          = 1
  tasks_per_core = 2
  tasks_per_node = 6
  cpus_per_task  = 1
  partition      = "debug"
  constraint     = "gpu"
  executable     = ""
  machine_name   = "daint"
  env            = ""
}





output "firecrest_job_id" {
  value = firecrest_job.job.id
}

output "firecrest_task_id" {
  value = firecrest_job.task.id
}