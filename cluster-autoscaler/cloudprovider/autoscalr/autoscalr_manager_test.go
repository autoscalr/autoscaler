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
	"os"
)

func TestOne(t *testing.T) {
	assert.Equal(t, "us-east-1", "us-east-1")
}

func TestCreateAutoScalrManager(t *testing.T) {
	asrMgr, _ := CreateAutoScalrManager(nil)
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

func TestAppDefDelete(t *testing.T){
	os.Setenv("AUTOSCALING_GROUP_NAME", "testASG")
	os.Setenv("AWS_REGION", "us-east-1")

	err := appDefDelete()
	assert.NoError(t, err)
}
