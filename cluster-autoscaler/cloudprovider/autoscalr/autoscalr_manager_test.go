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
<<<<<<< HEAD
	"os"
=======
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"os"
	apiv1 "k8s.io/api/core/v1"
>>>>>>> cluster-autoscaler-release-1.2
)

func TestOne(t *testing.T) {
	assert.Equal(t, "us-east-1", "us-east-1")
}

func TestCreateAutoScalrManager(t *testing.T) {
<<<<<<< HEAD
	asrMgr, _ := CreateAutoScalrManager(nil)
=======
	do := cloudprovider.NodeGroupDiscoveryOptions{}
	asrMgr, _ := CreateAutoScalrManager(nil, do)
>>>>>>> cluster-autoscaler-release-1.2
	assert.NotNil(t, asrMgr)
}

func TestAppDefCreate(t *testing.T){
	os.Setenv("AUTOSCALING_GROUP_NAME", "testASG")
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("INSTANCE_TYPES", "c3.large,c3.xlarge")
	os.Setenv("TARGET_CAPACITY_VCPUS", "1")
	err := appDefCreate()

	assert.NoError(t, err)
}

func TestAppDefRead(t *testing.T){
	os.Setenv("AUTOSCALING_GROUP_NAME", "testASG")
	os.Setenv("AWS_REGION", "us-east-1")

	appDefTest, err := appDefRead()

	assert.NoError(t, err, appDefTest)
}

func TestAppDefUpdate(t *testing.T){
	os.Setenv("AUTOSCALING_GROUP_NAME", "testASG")
	os.Setenv("AWS_REGION", "us-east-1")
	target_capacity := 2

	err := appDefUpdate(target_capacity)
	assert.NoError(t, err)
}

func TestAppDefDeleteNodes(t *testing.T){
	os.Setenv("AUTOSCALING_GROUP_NAME", "testASG")
	os.Setenv("AWS_REGION", "us-east-1")
	delta_vcpu := 1
	nodes2Del := make([]string, 1)
	nodes2Del[0] = "nodeId1"

	err := appDefDeleteNodes(delta_vcpu, nodes2Del)
	assert.NoError(t, err)
}

func TestAppDefDelete(t *testing.T){
	os.Setenv("AUTOSCALING_GROUP_NAME", "testASG")
	os.Setenv("AWS_REGION", "us-east-1")

	err := appDefDelete()
	assert.NoError(t, err)
}
<<<<<<< HEAD
=======

func TestApplyLabels(t *testing.T){
	scsResp := new(SendClusterStateResponse)
	lblEntry := new(LabelUpdate)
	lblEntry.InstanceId = "id1"
	lblEntry.UID = "uid1"
	lblEntry.PayModel = "spot"
	scsResp.LabelUpdates = append(scsResp.LabelUpdates, *lblEntry)
	nodes := make([]apiv1.Node, 1)

	err := ApplyLabels(scsResp, nodes)
	assert.NoError(t, err)
}
>>>>>>> cluster-autoscaler-release-1.2
