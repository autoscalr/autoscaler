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
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

func SetEnvTestCase1() {
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
//func getDiscoveryOptionsTestCase1() cloudprovider.NodeGroupDiscoveryOptions {
//	return cloudprovider.NodeGroupDiscoveryOptions{
//		NodeGroupSpecs: []string{"1:6:asgName"},
//		NodeGroupAutoDiscoverySpec: "",
//	}
//}

func TestEnvSetCorrectly(t *testing.T) {
	SetEnvTestCase1()
	assert.Equal(t, os.Getenv("AWS_REGION"), "us-east-1")
}

func TestBuildAutoScalrCloudProvider(t *testing.T) {
	SetEnvTestCase1()
	do := cloudprovider.NodeGroupDiscoveryOptions{}
	resourceLimiter := cloudprovider.NewResourceLimiter(
		map[string]int64{cloudprovider.ResourceNameCores: 1, cloudprovider.ResourceNameMemory: 10000000},
		map[string]int64{cloudprovider.ResourceNameCores: 10, cloudprovider.ResourceNameMemory: 100000000})

	asrMgr, err := CreateAutoScalrManager(nil, do)
	awsMgr, err := aws.CreateAwsManager(nil, do)
	assert.NoError(t, err)
	//rl := nil
	//resourceLimiter := cloudprovider.NewResourceLimiter(
	//	map[string]int64{cloudprovider.ResourceNameCores: 1, cloudprovider.ResourceNameMemory: 10000000},
	//	map[string]int64{cloudprovider.ResourceNameCores: 10, cloudprovider.ResourceNameMemory: 100000000})
	asrCloudProv, err := BuildAutoScalrCloudProvider(asrMgr, resourceLimiter, awsMgr)
	assert.NoError(t, err)
	assert.NotNil(t, asrCloudProv)
	assert.Equal(t, asrCloudProv.Name(), "autoscalr")
	nodeGrps := asrCloudProv.NodeGroups()
	assert.Equal(t, len(nodeGrps), 1)
	ng1 := nodeGrps[0]
	assert.Equal(t, ng1.MaxSize(), 6)
	assert.Equal(t, ng1.MinSize(), 1)
	assert.Equal(t, ng1.Id(), "asgName")
	//assert.True(t, ng1.Exist())
}

// Create a mock for awsProvider that requests will be forwarded to by default
type CloudProviderMock struct {
	mock.Mock
}

func (a *CloudProviderMock) Name() (string) {
	args := a.Called()
	return args.String(0)
}