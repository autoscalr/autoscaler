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
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/aws"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/kubernetes/plugin/pkg/scheduler/schedulercache"
	"github.com/golang/glog"
)

// autoScalrCloudProvider implements CloudProvider interface.
type autoScalrCloudProvider struct {
	autoScalrManager	*AutoScalrManager
	awsProvider			cloudprovider.CloudProvider
}

func BuildAutoScalrCloudProvider(autoScalrManager *AutoScalrManager, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions, resourceLimiter *cloudprovider.ResourceLimiter, awsManager *aws.AwsManager) (*autoScalrCloudProvider, error) {
	awsProv, err := aws.BuildAwsCloudProvider(awsManager, discoveryOpts, resourceLimiter)
	provider := &autoScalrCloudProvider{
		autoScalrManager: autoScalrManager,
		awsProvider: awsProv,
	}
	return provider, err
}

// Name returns name of the cloud provider.
func (asrProvider *autoScalrCloudProvider) Name() string {
	return "autoscalr"
}

// NodeGroups returns all node groups configured for this cloud provider.
func (asrProvider *autoScalrCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	awsNGs := asrProvider.awsProvider.NodeGroups()
	asrNGs := make([]cloudprovider.NodeGroup, 0, len(awsNGs))
	for _, nodeGrp := range awsNGs {
		asrNGs = append(asrNGs, BuildAutoScalrNodeGroup(nodeGrp))
	}
	return asrNGs
}

// NodeGroupForNode returns the node group for the given node.
func (asrProvider *autoScalrCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	awsNg, err := asrProvider.awsProvider.NodeGroupForNode(node)
	if err != nil {
		return awsNg, err
	} else {
		// wrap in asrNode
		return BuildAutoScalrNodeGroup(awsNg), err
	}
}

// Pricing returns pricing model for this cloud provider or error if not available.
func (asrProvider *autoScalrCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	return asrProvider.awsProvider.Pricing()
}

// GetAvailableMachineTypes get all machine types that can be requested from the cloud provider.
func (asrProvider *autoScalrCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	return asrProvider.awsProvider.GetAvailableMachineTypes()
}

// NewNodeGroup builds a theoretical node group based on the node definition provided. The node group is not automatically
// created on the cloud provider side. The node group is not returned by NodeGroups() until it is created.
func (asrProvider *autoScalrCloudProvider) NewNodeGroup(machineType string, labels map[string]string, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	awsNg, err := asrProvider.awsProvider.NewNodeGroup(machineType, labels, extraResources)
	if err != nil {
		return awsNg, err
	} else {
		// wrap in asrNode
		return BuildAutoScalrNodeGroup(awsNg), err
	}
}

// GetResourceLimiter returns struct containing limits (max, min) for resources (cores, memory etc.).
func (asrProvider *autoScalrCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	return asrProvider.awsProvider.GetResourceLimiter()
}

// Cleanup stops the go routine that is handling the current view of the ASGs in the form of a cache
func (asrProvider *autoScalrCloudProvider) Cleanup() error {
	return asrProvider.awsProvider.Cleanup()
}

// Refresh is called before every main loop and can be used to dynamically update cloud provider state.
// In particular the list of node groups returned by NodeGroups can change as a result of CloudProvider.Refresh().
func (asrProvider *autoScalrCloudProvider) Refresh() error {
	return asrProvider.awsProvider.Refresh()
}

// asrNodeGroup implements NodeGroup interface, defaulting to pass through to awsNodeGroup object
type asrNodeGroup struct {
	awsNodeGroup			cloudprovider.NodeGroup
}

func BuildAutoScalrNodeGroup(aNode cloudprovider.NodeGroup) (cloudprovider.NodeGroup) {
	asrNG := &asrNodeGroup{
		awsNodeGroup: aNode,
	}
	return asrNG
}

func (asrNG *asrNodeGroup) MaxSize() int {
	return asrNG.awsNodeGroup.MaxSize()
}

func (asrNG *asrNodeGroup) MinSize() int {
	return asrNG.awsNodeGroup.MinSize()
}

func (asrNG *asrNodeGroup) TargetSize() (int, error) {
	glog.V(0).Infof("AsrNodeGroup::TargetSize")
	return asrNG.awsNodeGroup.TargetSize()
}

func (asrNG *asrNodeGroup) IncreaseSize(delta int) error {
	glog.V(0).Infof("AsrNodeGroup::IncreaseSize")
	return asrNG.awsNodeGroup.IncreaseSize(delta)
}

func (asrNG *asrNodeGroup) DeleteNodes(nodes []*apiv1.Node) error {
	glog.V(0).Infof("AsrNodeGroup::DeleteNodes")
	return asrNG.awsNodeGroup.DeleteNodes(nodes)
}

func (asrNG *asrNodeGroup) DecreaseTargetSize(delta int) error {
	glog.V(0).Infof("AsrNodeGroup::DecreaseTargetSize")
	return asrNG.awsNodeGroup.DecreaseTargetSize(delta)
}

func (asrNG *asrNodeGroup) Id() string {
	return asrNG.awsNodeGroup.Id()
}

func (asrNG *asrNodeGroup) Debug() string {
	return asrNG.awsNodeGroup.Debug()
}

func (asrNG *asrNodeGroup) Nodes() ([]string, error) {
	return asrNG.awsNodeGroup.Nodes()
}

func (asrNG *asrNodeGroup) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
	return asrNG.awsNodeGroup.TemplateNodeInfo()
}

func (asrNG *asrNodeGroup) Exist() bool {
	return asrNG.awsNodeGroup.Exist()
}

func (asrNG *asrNodeGroup) Create() error {
	return asrNG.awsNodeGroup.Create()
}

func (asrNG *asrNodeGroup) Delete() error {
	return asrNG.awsNodeGroup.Delete()
}

func (asrNG *asrNodeGroup) Autoprovisioned() bool {
	return asrNG.awsNodeGroup.Autoprovisioned()
}
