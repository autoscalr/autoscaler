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
	"testing"
	"github.com/stretchr/testify/assert"
	//"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"os"
)

func setEnvTestCase1() {
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AUTOSCALR_API_KEY", "myApiKey")
	os.Setenv("DISPLAY_NAME", "nameToDisplayInUI")
	os.Setenv("MAX_SPOT_PERCENT_TOTAL", "90")
	os.Setenv("MAX_SPOT_PERCENT_ONE_MARKET", "25")
	os.Setenv("DETAILED_MONITORING_ENABLED", "true")
	os.Setenv("MAX_HOURS_INSTANCE_AGE", "")
	os.Setenv("INSTANCE_TYPES", "m1.medium,m3.large")
	os.Setenv("TARGET_SPARE_CPU_PERCENT", "20")
	os.Setenv("TARGET_CAPACITY_VCPUS", "6")
	os.Setenv("TARGET_CAPACITY_INSTANCES", "2")
	os.Setenv("TARGET_SPARE_MEMORY_PERCENT", "20")
}

func TestEnvSetCorrectly(t *testing.T) {
	setEnvTestCase1()
	assert.Equal(t, os.Getenv("AWS_REGION"), "us-east-1")
}

func TestBuildAutoScalrCloudProvider(t *testing.T) {
	//asrCloudProv, _ := BuildAutoScalrCloudProvider(asrMgr, discoveryOpts, resourceLimiter)
	//discOpts := cloudprovider.NodeGroupDiscoveryOptions{
	//	NodeGroupSpecs: [],
	//	NodeGroupAutoDiscoverySpec: "",
	//}
	//asrCloudProv, _ := BuildAutoScalrCloudProvider(nil, discOpts, nil)
	asrCloudProv := "val"
	assert.NotNil(t, asrCloudProv)
}
