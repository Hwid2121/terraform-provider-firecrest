# Manages a Firecrest Job.
resource "firecrest_job" "example" {
  job_name       = "example-job"
  machine_name   = "machine-example"
  account        = "account-example"
  hours          = 1
  minutes        = 30
  nodes          = 1
  tasks_per_core = 1
  tasks_per_node = 1
  cpus_per_task  = 1
  partition      = "normal"
  executable     = "example_executable.sh"
}
