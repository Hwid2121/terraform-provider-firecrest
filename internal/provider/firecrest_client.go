package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	// "os"
	// "path/filepath"

	// "time"

	// "log"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	// "golang.org/x/tools/go/analysis/passes/defers"
	// "google.golang.org/grpc/internal/status"
	// "google.golang.org/grpc/status"
)

type FirecrestClient struct {
	clientID string
	clientSecret string
	baseURL string
	apiToken string
	httpClient *http.Client
}

type JobStatus struct {
	Success string `json:"success"`
	TaskID string `json:"task_id"`
	TaskURL string `json:"task_url"`
}

type TaskStatus struct {
	CreatedAt   string                 `json:"created_at"`
    Data        interface{} 		   `json:"data"`
	Description string                 `json:"description"`
	HashID      string                 `json:"hash_id"`
	LastModify  string                 `json:"last_modify"`
	Service     string                 `json:"service"`
	Status      string                 `json:"status"`
	System      string                 `json:"system"`
	TaskID      string                 `json:"task_id"`
	TaskURL     string                 `json:"task_url"`
	UpdatedAt   string                 `json:"updated_at"`
	User        string                 `json:"user"`
}


type TaskData struct {
    JobDataErr   string `json:"job_data_err"`
    JobDataOut   string `json:"job_data_out"`
    JobFile      string `json:"job_file"`
    JobFileErr   string `json:"job_file_err"`
    JobFileOut   string `json:"job_file_out"`
    JobInfoExtra string `json:"job_info_extra"`
    JobID        int    `json:"jobid"`
    Result       string `json:"result"`
}





func NewFireCrestClient(clientID, clientSecret string) *FirecrestClient {
	return &FirecrestClient{
		clientID: clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{},
		baseURL: "https://firecrest.cscs.ch",
		// apiToken: apiToken,
	}
}

func (c *FirecrestClient) GetToken(clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	log.Println("baseURL: ", c.baseURL)	

	// log.Println("Request Body:", data.Encode())



	req, err := http.NewRequest("POST", "https://auth.cscs.ch//auth/realms/firecrest-clients/protocol/openid-connect/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve token, status code: %s", resp.Status)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("failed to parse access token")
	}

	return token, nil
}



func (c *FirecrestClient) SetToken(token string) {
	c.apiToken = token
}

func (c *FirecrestClient) UploadFile(sourcePath, targetPath string) error {
	return nil
}

func (c *FirecrestClient) DeleteFile(sourcePath, targetPath string) error {
	return nil
}



/*
	REST management for tasks
*/
func (c *FirecrestClient) GetTaskStatus(ctx context.Context, taskID  string) (*TaskStatus, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tasks/%s", c.baseURL, taskID), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer " + c.apiToken)
	req.Header.Set("accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get task status, status code: %s", resp.Status)
	}


	var taskStatus struct {
		Task TaskStatus `json:"task"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&taskStatus); err != nil {
		return nil, err
	}

	ctx = tflog.SetField(ctx, "Task Status  ", taskStatus)
	tflog.Debug(ctx, "task status json")

	return &taskStatus.Task, nil
}


func (c *FirecrestClient) WaitForJobID(ctx context.Context, taskID string) (string, error) {


	for {
		taskStatus, err := c.GetTaskStatus(ctx, taskID)
		if err != nil {
			return "", err
		}

        if taskStatus.Status == "200" {
            if dataMap, ok := taskStatus.Data.(map[string]interface{}); ok {
                jobID, ok := dataMap["jobid"].(float64) // JSON numbers are decoded as float64
                if !ok {
                    return "", fmt.Errorf("job ID not found in task data")
                }
                return fmt.Sprintf("%.0f", jobID), nil
            } else {
                return "", fmt.Errorf("unexpected data format: %v", taskStatus.Data)
            }
        }
		

        if taskStatus.Status == "400" {
            if dataString, ok := taskStatus.Data.(string); ok {
                return "", fmt.Errorf("task %s failed. Description: %s, Data: %s", taskID, taskStatus.Description, dataString)
            } else {
                return "", fmt.Errorf("task %s failed. Description: %s", taskID, taskStatus.Description)
            }
        }
		time.Sleep(5 * time.Second)

		// return "", fmt.Errorf("task %s boh. \n Description: \t %s", taskID, taskStatus.Description)
	}

}

/*
	REST management for JOB
*/


func (c *FirecrestClient) DeleteJob(jobID, machineName string)  error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/compute/jobs/%s", c.baseURL,jobID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-Machine-Name", machineName)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if  resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete job, status code: %s", resp.Status)
	}

	return nil
}





func (c *FirecrestClient) GetJobStatus(ctx context.Context, jobID  string, machineName string) (*JobStatus, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/compute/jobs/%s", c.baseURL, jobID), nil)


	ctx = tflog.SetField(ctx, "Get Job Status req: ", req)
	tflog.Debug(ctx, "Get Job Status ")

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer " + c.apiToken)
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-Machine-Name", machineName)


	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get job status, status code: %s", resp.Status)
	}


	var jobStatus JobStatus
	if err := json.NewDecoder(resp.Body).Decode(&jobStatus); err != nil {
		return nil, err
	}

	ctx = tflog.SetField(ctx, "Job Status  ", jobStatus)
	tflog.Debug(ctx, "job status json")
	return &jobStatus, nil

}


func (c *FirecrestClient) UploadJob(JobScript, Account, Env, MachineName string ) (string, error) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", "job_script.sh")
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(fw, bytes.NewReader([]byte(JobScript))); err != nil {
		return "", err
	}
	// Add other form fields
	if err := w.WriteField("type", "application/x-shellscript"); err != nil {
		return "", err
	}
	if err := w.WriteField("account", Account); err != nil {
		return "", err
	}
	if err := w.WriteField("env", Env); err != nil {
		return "", err
	}

	w.Close()

	req, err := http.NewRequest("POST", c.baseURL+"/compute/jobs/upload", &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-Machine-Name", MachineName)


	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("failed to submit job, status code: %s", resp.Status)
	}


	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	
	taskID, ok := result["task_id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to parse task id")
	}

	return taskID, nil
}	


