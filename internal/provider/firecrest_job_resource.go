package provider

import (
	"context"
	"fmt"
	"strings"

	// "strings"

	// "strings"

	// "strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	// "golang.org/x/mod/sumdb/tlog"
	// "github.com/hashicorp/terraform-plugin-log/tflog"
	// "github.com/hashicorp/terraform-svchost/disco"
)

var (
	_ resource.Resource              = &firecrestJobResource{}
	_ resource.ResourceWithConfigure = &firecrestJobResource{}
)

func NewFirecrestJobResource() resource.Resource {
	return &firecrestJobResource{}
}

type firecrestJobResource struct {
	client *FirecrestClient
}

type firecrestJobResourceModel struct {
	ID types.String `tfsdk:"id"`
	JobID types.String `tfsdk:"job_id"`
	State types.String `tfsdk:"state"`
	OutputFile types.String `tfsdk:"output_file"`
	JobScript types.String `tfsdk:"job_script"`
	MachineName types.String `tfsdk:"machine_name"`
	AccountName types.String `tfsdk:"account"`
	Env types.String `tfsdk:"env"`

	JobName types.String `tfsdk:"job_name"`
	Email types.String `tfsdk:"email"`
	Hours types.Int64 `tfsdk:"hours"`
	Minutes types.Int64 `tfsdk:"minutes"`
	Nodes types.Int64 `tfsdk:"nodes"`
	TasksPerCore types.Int64 `tfsdk:"tasks_per_core"`
	TasksPerNode types.Int64 `tfsdk:"tasks_per_node"`
	CpuPerTask types.Int64 `tfsdk:"cpus_per_task"`
	Partition types.String `tfsdk:"partition"`
	Constraint types.String `tfsdk:"constraint"`
	Executable types.String `tfsdk:"executable"`
}




// Schema implements resource.Resource.
func (f *firecrestJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a FirecREST Job.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the job.",
				Computed: true,
			},
			"job_id": schema.StringAttribute{
				Description: "The job ID of the submitted job.",
				Computed: true,
			},
			"state": schema.StringAttribute{
				Description: "The state of the job.",
				Computed: true,
			},
			"output_file": schema.StringAttribute{
				Description: "The path to the job's output File.",
				Computed: true,
			},
			"job_script": schema.StringAttribute{
				Description: "The sbatch script to be submitted.",
				Optional: true,
			},
			"machine_name": schema.StringAttribute{
				Description: "The name of the machine where the job will run.",
				Required: true,
			},
			"account": schema.StringAttribute{
				Description: "Account name for the job.",
				Required: true,
			},
			"env": schema.StringAttribute{
				Description: "The enviroment variables for the job.",
				Optional: true,
			},

			"job_name": schema.StringAttribute{
				Description: "The name for job.",
				Required: true,
			},
		
			"email": schema.StringAttribute{
				Description: "Specify your email address to get notified when the job changes state.",
				Optional: true,
			},
			"hours": schema.Int64Attribute{
				Description: "The hours allocated for the job.",
				Required: true,
			},
			"minutes": schema.Int64Attribute{
				Description: "The minutes allocated for the job.",
				Required: true,
			},
			"nodes": schema.Int64Attribute{
				Description: "Specify the number of nodes.",
				Required: true,
			},
			"tasks_per_core": schema.Int64Attribute{
				Description: "The number of tasks per core. Values greater than one turn hyperthreading on.",
				Required: true,
			},
			"tasks_per_node": schema.Int64Attribute{
				Description: "The number of tasks per node. Defines the number of MPI ranks per node. The maximum value depends on the number of cpus per task.",
				Required: true,
			},
			"cpus_per_task": schema.Int64Attribute{
				Description: "The number of cpus per task. Defines the number of OpenMP threads per MPI rank. The maximum value depends on the number of tasks per node.",
				Required: true,
			},
			"partition": schema.StringAttribute{
				Description: "The partition on which you want to submit your job. (normal, low, xfer, debug, prepost)",
				Required: true,
			},
			"constraint": schema.StringAttribute{
				Description: "The constraint for the job submission.",
				Optional: true,
			},
			"executable": schema.StringAttribute{
				Description: "The executable to run in the job.",
				Required: true,
			},



		},
	}
}



func (r *firecrestJobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp * resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerConfig, ok := req.ProviderData.(*firecrestProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected provider to be of type *firecrestProvider, got %T", req.ProviderData),
		)
		return
	}

	r.client = providerConfig.client
}





// Metadata implements resource.Resource.
func (f *firecrestJobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job"
}



// Create implements resource.Resource.
func (r *firecrestJobResource) Create(ctx context.Context,req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan firecrestJobResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return 
	}



	var jobScript string
	var err error


	if plan.JobScript.IsNull() {
		jobScript, err = generateJobScript(plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error generating Job Script",
				fmt.Sprintf("Could not generate job script: %s", err.Error()),
			)
			return
		}
	} else {
		jobScript = plan.JobScript.ValueString()
	}


	taskID, err := r.client.UploadJob(jobScript,  
									 plan.AccountName.ValueString(), 
									 plan.Env.ValueString(), 
									 plan.MachineName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error submitting Job",
			fmt.Sprintf("Could not submit job: %s", err.Error()),
		)
		return 
	}
	
	ctx = tflog.SetField(ctx, "Task ID: ", taskID)
	tflog.Debug(ctx, "Created Task!")

	 
	
	jobID, err := r.client.WaitForJobID(ctx, taskID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for jobID",
			fmt.Sprintf("Could not get job ID: %s", err.Error()),
		)
		return 
	}



	// ctx = tflog.SetField(ctx, "JOBID: ", types.StringValue(jobID))
	ctx = tflog.SetField(ctx, "JOBID2: ", jobID)
	// ctx = tflog.SetField(ctx, "JOBID3: ", ./)


	tflog.Debug(ctx, "CREATE status")


	// jobID = strings.TrimLeft(jobID, "=") // Ensure no leading "="
	plan.ID = types.StringValue(jobID)
	plan.JobID = types.StringValue(jobID)
	plan.State = types.StringValue("SUBMITTED")
	plan.OutputFile = types.StringValue("")
	// plan.Env = types.StringValue("")
	// plan.MachineName = types.StringValue("")
	// plan.AccountName = types.StringValue("")

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


}

// Delete implements resource.Resource.
func (f *firecrestJobResource) Delete(ctx context.Context,req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan firecrestJobResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return 
	}

	err := f.client.DeleteJob(plan.JobID.ValueString(), plan.MachineName.ValueString()) 
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Job", 
			fmt.Sprintf("Could not delete job: %s", err.Error()),
		)
		return 
	}
 }


// Read implements resource.Resource.
func (f *firecrestJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan firecrestJobResourceModel

	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Clean the job ID before using it
	// jobID := "54037022"
	// jobID := plan.JobID.ValueString()
    // if strings.HasPrefix(jobID, "=") {
        // jobID = jobID[1:]
    // }

	jobID := plan.JobID.String()
    jobID = strings.Trim(jobID, "=\"")

	ctx = tflog.SetField(ctx, "JOBID: ", jobID)
	tflog.Debug(ctx, "READ status")

	jobStatus, err := f.client.GetJobStatus(ctx, plan.JobID.ValueString(), plan.MachineName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Job",
			fmt.Sprintf("Could not read job status: %s", err.Error()),
		)
		return
	}

	plan.State = types.StringValue(jobStatus.Success)
	// state.OutputFile = types.StringValue(jobStatus.OutputFile)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return 
	}
}

	// Update implements resource.Resource.
	func (f *firecrestJobResource) Update(ctx context.Context,req resource.UpdateRequest, resp *resource.UpdateResponse) {
		var plan firecrestJobResourceModel
		
		diags := req.Plan.Get(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		var jobScript string
		var err error

		if plan.JobScript.IsNull() {
			jobScript, err = generateJobScript(plan)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error generating updated Job Script",
					fmt.Sprintf("Could not generate job script: %s", err.Error()),
				)
				return
			}
		} else {
			jobScript = plan.JobScript.ValueString()
		}

		newTaskID, err := f.client.UploadJob(
			jobScript, plan.AccountName.ValueString(),
			plan.Env.ValueString(), plan.MachineName.ValueString())

		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating Job",
				fmt.Sprintf("Could not update job: %s", err.Error()),
			)
			return
		}



		ctx = tflog.SetField(ctx, "New Task ID: ", newTaskID)
		tflog.Debug(ctx, "Created new Task for update!")

		newJobID, err := f.client.WaitForJobID(ctx, newTaskID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for new jobID",
				fmt.Sprintf("Could not get new job ID: %s", err.Error()),
			)
			return
		}


		// Optionally delete the old job
		if plan.JobID.ValueString() != "" {
			err = f.client.DeleteJob(plan.JobID.ValueString(), plan.MachineName.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error deleting old Job",
					fmt.Sprintf("Could not delete old job: %s", err.Error()),
				)
				return
			}
		}


		plan.ID = types.StringValue(newJobID)
		plan.JobID = types.StringValue(newJobID)
		plan.State = types.StringValue("UPDATED")
		plan.OutputFile = types.StringValue("")
	
		diags = resp.State.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}






// func generateOptionalParameters(plan firecrestJobResourceModel) string {
// 	params := ""
// 	if !plan.Constraint.IsNull() {
// 		params += fmt.Sprintf("#SBATCH --contraint=%s\n", plan.Constraint.ValueString())
// 	}
// 	return params
// }

func generateJobScript(plan firecrestJobResourceModel) (string, error) {
	walltime, err := ConvertHoursMinutesToWalltime(plan.Hours.ValueInt64(), plan.Minutes.ValueInt64())
	if err != nil {
		return "", err
	}

	script := fmt.Sprintf(
		`#!/bin/bash -l
#SBATCH --job-name="%s"
#SBATCH --mail-type=ALL
#SBATCH --mail-user="%s"
#SBATCH --time=%s
#SBATCH --nodes=%d
#SBATCH --ntasks-per-core=%d
#SBATCH --ntasks-per-node=%d
#SBATCH --cpus-per-task=%d
#SBATCH --partition=%s
%s

export OMP_NUM_THREADS=$SLURM_CPUS_PER_TASK
export CRAY_CUDA_MPS=1

srun %s
`,
		plan.JobName.ValueString(),
		plan.Email.ValueString(),
		walltime,
		plan.Nodes.ValueInt64(),
		plan.TasksPerCore.ValueInt64(),
		plan.TasksPerNode.ValueInt64(),
		plan.CpuPerTask.ValueInt64(),
		plan.Partition.ValueString(),
		optionalField(plan.Constraint.ValueString(), "#SBATCH --constraint=%s"),
		plan.Executable.ValueString(),
	)

	return script, nil
}


func optionalField(value, format string) string {
	if value == "" {
		return ""
	}
	return fmt.Sprintf(format, value)
}

func ConvertHoursMinutesToWalltime(hours, minutes int64) (string, error) {
	if hours < 0 || minutes < 0 || minutes >= 60 {
		return "", fmt.Errorf("invalid input: hours should be >= 0 and minutes should be between 0 and 59")
	}

	walltime := fmt.Sprintf("%02d:%02d:00", hours, minutes)
	return walltime, nil
}


