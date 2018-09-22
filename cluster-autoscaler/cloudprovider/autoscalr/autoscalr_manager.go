/*
Copyright 2017 AutoScalr

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package autoscalr

import (
	"io"
	"bytes"
	"fmt"
	"encoding/json"
	"net/http"
	"time"
	"errors"
	"os"
	"strings"
	"strconv"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	apiappsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
)

// AutoScalrManager is handles communication and data caching.
type AutoScalrManager struct {
	random   string
}

func createAutoScalrManagerInternal(configReader io.Reader, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions) (*AutoScalrManager, error) {
	manager := &AutoScalrManager{
		random: "Test-jay",
	}
	return manager, nil
}

func CreateAutoScalrManager(configReader io.Reader, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions) (*AutoScalrManager, error) {
	return createAutoScalrManagerInternal(configReader, discoveryOpts)
}

type AppDef struct {
	AutoScalingGroupName        string   `json:"aws_autoscaling_group_name"`
	AwsRegion                   string   `json:"aws_region"`
	AppType						string	 `json:"app_type'`
	InstanceTypes               []string `json:"instance_types"`
	ScaleMode                   string   `json:"scale_mode"`
	MaxSpotPercentTotal         int      `json:"max_spot_percent_total"`
	MaxSpotPercentOneMarket     int      `json:"max_spot_percent_one_market"`
	TargetSpareCPUPercent       int      `json:"target_spare_cpu_percent"`
	ClusterName                 string   `json:"cluster_name"`
	TargetSpareMemoryPercent    int      `json:"target_spare_memory_percent"`
	QueueName                   string   `json:"queue_name"`
	TargetQueueSize             int      `json:"target_queue_size"`
	InstanceSpinUpSeconds       int      `json:"instance_spin_up_seconds"`
	MaxMinutesToTargetQueueSize int      `json:"max_minutes_to_target_queue_size"`
	DisplayName                 string   `json:"display_name"`
	DetailedMonitoringEnabled   bool     `json:"detailed_monitoring_enabled"`
	AutoscalrEnabled            bool     `json:"autoscalr_enabled"`
	OsFamily                    string   `json:"os_family"`
	MaxHoursInstanceAge         int      `json:"max_hours_instance_age"`
	TargetCapacity		        int      `json:"target_capacity"`
}

type AppDefUpdate struct {
	AutoScalingGroupName        string   `json:"aws_autoscaling_group_name"`
	AwsRegion                   string   `json:"aws_region"`
	TargetCapacity		        int      `json:"target_capacity"`
}
type AppDefNodeDelete struct {
	AutoScalingGroupName        string   `json:"aws_autoscaling_group_name"`
	AwsRegion                   string   `json:"aws_region"`
	DeltaVCpu			        int      `json:"delta_vcpu"`
	NodesToDelete               []string `json:"nodes_to_delete"`
}

type AutoScalrRequest struct {
	AsrToken    string  `json:"api_key"`
	RequestType string  `json:"request_type"`
	AsrAppDef   *AppDef `json:"autoscalr_app_def"`
}

type AutoScalrUpdateRequest struct {
	AsrToken    string  `json:"api_key"`
	RequestType string  `json:"request_type"`
	AsrAppDef   *AppDefUpdate `json:"autoscalr_app_def"`
}

type AutoScalrNodeDeleteRequest struct {
	AsrToken    string  `json:"api_key"`
	RequestType string  `json:"request_type"`
	AsrAppDef   *AppDefNodeDelete `json:"autoscalr_app_def"`
}

type AutoScalrClusterState struct {
	AsrToken    string  `json:"api_key"`
	AwsRegion   string  `json:"AwsRegion"`
	AutoScalingGroupName   string  `json:"AutoScalingGroupName"`
	Deployments []apiappsv1.Deployment `json:"deployments"`
	Nodes []apiv1.Node `json:"nodes"`
}

type AsrDeployment struct {
	Name    string  `json:"Name"`
}

type AsrApiErrorResponse struct {
	Error    *AsrApiError  `json:"error"`
}

type AsrApiError struct {
	ErrorMessage    	string  `json:"errorMessage"`
	Code 	 	string  `json:"code"`
}

func numVCpusBaseType() int {
	instanceTypesStr := os.Getenv("INSTANCE_TYPES")
	instanceTypesArr := strings.Split(instanceTypesStr, ",")
	baseType := instanceTypesArr[0]
	baseTypeStats := InstanceTypes[baseType]
	return int(baseTypeStats.VCPU)
}

func InstanceIdFromProviderId(id string) (string) {
	splitted := strings.Split(id[7:], "/")
	return splitted[1]
}

func SendClusterState(cState *AutoScalrClusterState) (int, error) {
	url := "https://api.autoscalr.com/v1/k8sClusterState"
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	postBody := new(bytes.Buffer)
	json.NewEncoder(postBody).Encode(cState)
	resp, err := client.Post(url, "application/json", postBody)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			// make 2 copies of response, one for error decoding and one for good response
			respBuf := new(bytes.Buffer)
			respBuf.ReadFrom(resp.Body)
			errBuf := bytes.NewBuffer(respBuf.Bytes())
			// Check for error response json
			jsonErr := new(AsrApiErrorResponse)
			json.NewDecoder(errBuf).Decode(jsonErr)
			if jsonErr.Error != nil && jsonErr.Error.ErrorMessage != ""  {
				// error response
				err = errors.New(fmt.Sprintf("Error response: %s", jsonErr.Error.ErrorMessage))
			} else {
				// looks like good response
				//json.NewDecoder(respBuf).Decode(app)
			}
			return resp.StatusCode, err
		} else {
			err = errors.New(fmt.Sprintf("k8sClusterStateAPI returned: %d", resp.Status))
			return resp.StatusCode, err
		}
	} else {
		//log.Println("Error: %s", err.Error())
		return 500, err
	}
}

func makeApiCall(asrReq *AutoScalrRequest) (int, *AppDef, error) {
	url := "https://app.autoscalr.com/api/autoScalrApp"
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	postBody := new(bytes.Buffer)
	json.NewEncoder(postBody).Encode(asrReq)
	app := new(AppDef)
	resp, err := client.Post(url, "application/json", postBody)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			// make 2 copies of response, one for error decoding and one for good response
			respBuf := new(bytes.Buffer)
			respBuf.ReadFrom(resp.Body)
			errBuf := bytes.NewBuffer(respBuf.Bytes())
			// Check for error response json
			jsonErr := new(AsrApiErrorResponse)
			json.NewDecoder(errBuf).Decode(jsonErr)
			if jsonErr.Error != nil && jsonErr.Error.ErrorMessage != ""  {
				// error response
				err = errors.New(fmt.Sprintf("Error response: %s", jsonErr.Error.ErrorMessage))
			} else {
				// looks like good response
				json.NewDecoder(respBuf).Decode(app)
			}
			return resp.StatusCode, app, err
		} else {
			err = errors.New(fmt.Sprintf("AutoScalr API returned: %d", resp.Status))
			return resp.StatusCode, app, err
		}
	} else {
		//log.Println("Error: %s", err.Error())
		return 500, app, err
	}
}

func makeUpdateApiCall(asrReq *AutoScalrUpdateRequest) (int, *AppDef, error) {
	url := "https://app.autoscalr.com/api/autoScalrApp"
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	postBody := new(bytes.Buffer)
	json.NewEncoder(postBody).Encode(asrReq)
	app := new(AppDef)
	resp, err := client.Post(url, "application/json", postBody)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			// make 2 copies of response, one for error decoding and one for good response
			respBuf := new(bytes.Buffer)
			respBuf.ReadFrom(resp.Body)
			errBuf := bytes.NewBuffer(respBuf.Bytes())
			// Check for error response json
			jsonErr := new(AsrApiErrorResponse)
			json.NewDecoder(errBuf).Decode(jsonErr)
			if jsonErr.Error != nil && jsonErr.Error.ErrorMessage != ""  {
				// error response
				err = errors.New(fmt.Sprintf("Error response: %s", jsonErr.Error.ErrorMessage))
			} else {
				// looks like good response
				json.NewDecoder(respBuf).Decode(app)
			}
			return resp.StatusCode, app, err
		} else {
			err = errors.New(fmt.Sprintf("AutoScalr API returned: %d", resp.Status))
			return resp.StatusCode, app, err
		}
	} else {
		//log.Println("Error: %s", err.Error())
		return 500, app, err
	}
}

func makeDeleteNodesApiCall(asrReq *AutoScalrNodeDeleteRequest) (int, *AppDef, error) {
	url := "https://app.autoscalr.com/api/autoScalrApp"
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	postBody := new(bytes.Buffer)
	json.NewEncoder(postBody).Encode(asrReq)
	app := new(AppDef)
	resp, err := client.Post(url, "application/json", postBody)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			// make 2 copies of response, one for error decoding and one for good response
			respBuf := new(bytes.Buffer)
			respBuf.ReadFrom(resp.Body)
			errBuf := bytes.NewBuffer(respBuf.Bytes())
			// Check for error response json
			jsonErr := new(AsrApiErrorResponse)
			json.NewDecoder(errBuf).Decode(jsonErr)
			if jsonErr.Error != nil && jsonErr.Error.ErrorMessage != ""  {
				// error response
				err = errors.New(fmt.Sprintf("Error response: %s", jsonErr.Error.ErrorMessage))
			} else {
				// looks like good response
				json.NewDecoder(respBuf).Decode(app)
			}
			return resp.StatusCode, app, err
		} else {
			err = errors.New(fmt.Sprintf("AutoScalr API returned: %d", resp.Status))
			return resp.StatusCode, app, err
		}
	} else {
		//log.Println("Error: %s", err.Error())
		return 500, app, err
	}
}

func appDefCreate() error {
	instanceTypesStr := os.Getenv("INSTANCE_TYPES")
	instanceTypesArr := strings.Split(instanceTypesStr, ",")
	maxSpotPercTotal, _ := strconv.Atoi(os.Getenv("MAX_SPOT_PERCENT_TOTAL"))
	maxSpotPercOneMarket, _ := strconv.Atoi(os.Getenv("MAX_SPOT_PERCENT_ONE_MARKET"))
	maxHoursInstAge, _ := strconv.Atoi(os.Getenv("MAX_HOURS_INSTANCE_AGE"))
	targVcpuCapacity, _ := strconv.Atoi(os.Getenv("TARGET_CAPACITY_VCPUS"))
	detailedMonitoring, _ := strconv.ParseBool(os.Getenv("DETAILED_MONITORING_ENABLED"))
	body := &AutoScalrRequest{
		AsrToken:    os.Getenv("AUTOSCALR_API_KEY"),
		RequestType: "Create",
		AsrAppDef: &AppDef{
			AutoScalingGroupName:        os.Getenv("AUTOSCALING_GROUP_NAME"),
			AwsRegion:                   os.Getenv("AWS_REGION"),
			AppType:					 "k8s",
			InstanceTypes:               instanceTypesArr,
			ScaleMode:                   "fixed",
			MaxSpotPercentTotal:         maxSpotPercTotal,
			MaxSpotPercentOneMarket:     maxSpotPercOneMarket,
			ClusterName:                 "",
			TargetSpareCPUPercent:       0,
			TargetSpareMemoryPercent:    0,
			QueueName:                   "",
			TargetQueueSize:             0,
			InstanceSpinUpSeconds:       15,
			MaxMinutesToTargetQueueSize: 0,
			DisplayName:                 os.Getenv("DISPLAY_NAME"),
			DetailedMonitoringEnabled:   detailedMonitoring,
			AutoscalrEnabled:            true,
			OsFamily:                    os.Getenv("OS_FAMILY"),
			MaxHoursInstanceAge:         maxHoursInstAge,
			TargetCapacity:         	 targVcpuCapacity,
		},
	}
	respCode, _, err := makeApiCall(body)
	if respCode > 400 {
		err = fmt.Errorf("AutoScalr API returned status code: %d", respCode)
	}
	return err
}

func appDefRead() (*AppDef, error) {
	body := &AutoScalrRequest{
		AsrToken:    os.Getenv("AUTOSCALR_API_KEY"),
		RequestType: "Get",
		AsrAppDef: &AppDef{
			AutoScalingGroupName:        os.Getenv("AUTOSCALING_GROUP_NAME"),
			AwsRegion:                   os.Getenv("AWS_REGION"),
		},
	}
	respCode, app, err := makeApiCall(body)
	if respCode > 400 {
		err = fmt.Errorf("AutoScalr API returned status code: %d", respCode)
	}
	return app, err
}


func appDefUpdate(target_capacity int) error {
	body := &AutoScalrUpdateRequest{
		AsrToken:    os.Getenv("AUTOSCALR_API_KEY"),
		RequestType: "Update",
		AsrAppDef: &AppDefUpdate{
			AutoScalingGroupName:        os.Getenv("AUTOSCALING_GROUP_NAME"),
			AwsRegion:                   os.Getenv("AWS_REGION"),
			TargetCapacity:         	 target_capacity,
		},
	}
	respCode, _, err := makeUpdateApiCall(body)
	if respCode > 400 {
		err = fmt.Errorf("AutoScalr API returned status code: %d", respCode)
	}
	return err
}

func appDefDeleteNodes(deltaVcpu int, nodesToDel []string) error {
	body := &AutoScalrNodeDeleteRequest{
		AsrToken:    os.Getenv("AUTOSCALR_API_KEY"),
		RequestType: "DeleteAppNodes",
		AsrAppDef: &AppDefNodeDelete{
			AutoScalingGroupName:        os.Getenv("AUTOSCALING_GROUP_NAME"),
			AwsRegion:                   os.Getenv("AWS_REGION"),
			DeltaVCpu:         	 		1,
			NodesToDelete: 				nodesToDel,
		},
	}
	respCode, _, err := makeDeleteNodesApiCall(body)
	if respCode > 400 {
		err = fmt.Errorf("AutoScalr API returned status code: %d", respCode)
	}
	return err
}

func appDefDelete() error {
	body := &AutoScalrRequest{
		AsrToken:    os.Getenv("AUTOSCALR_API_KEY"),
		RequestType: "Delete",
		AsrAppDef: &AppDef{
			AutoScalingGroupName:        os.Getenv("AUTOSCALING_GROUP_NAME"),
			AwsRegion:                   os.Getenv("AWS_REGION"),
		},
	}
	respCode, _, err := makeApiCall(body)
	if respCode > 400 {
		err = fmt.Errorf("AutoScalr API returned status code: %d", respCode)
	}
	return err
}
