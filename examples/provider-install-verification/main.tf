terraform {
  required_providers {
    firecrest = {
      source  = "registry.terraform.io/hashicorp/firecrest"
      version = "6.0"
    }
  }
}

provider "firecrest" {
  client_id     = "firecrest-ntafta-coder"
  client_secret = "D1wLfcA3BfVzxYA7eJ7AivIEklWNTH3C"
}


resource "firecrest_job" "job" {
  # job_script = ""
  job_name       = "job-test"
  account        = "csstaff"
  email          = "nicolotafta@gmail.com"
  hours          = 0
  minutes        = 1
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
  value = "test"
}
