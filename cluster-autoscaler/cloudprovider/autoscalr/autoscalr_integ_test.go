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

package autoscalr_test

/*
	This package for performing integration testing, simulating how would be called from main & children.
	It is put in a separate package to force only external interfaces to be used.
	It tries to mimic the setup that main does before calling in to create the provider, duplicates many
	of the functions in main
 */
import (
	"testing"
	"strconv"
	"strings"
	"fmt"
	"time"
	"flag"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/autoscalr"
	"github.com/stretchr/testify/assert"
	kube_util "k8s.io/autoscaler/cluster-autoscaler/utils/kubernetes"
	kube_client "k8s.io/client-go/kubernetes"
	//"github.com/golang/glog"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"k8s.io/autoscaler/cluster-autoscaler/core"
	"k8s.io/autoscaler/cluster-autoscaler/estimator"
	"k8s.io/autoscaler/cluster-autoscaler/expander"
	"k8s.io/autoscaler/cluster-autoscaler/simulator"
)


func TestAutoScalrCloudProviderCreated(t *testing.T) {
	autoscalr.SetEnvTestCase1()
	kubeClient := createKubeClient()
	kubeEventRecorder := kube_util.CreateEventRecorder(kubeClient)
	opts := createAutoscalerOptions()
	predicateCheckerStopChannel := make(chan struct{})
	predicateChecker, _ := simulator.NewPredicateChecker(kubeClient, predicateCheckerStopChannel)
	//if err != nil {
	//	glog.Fatalf("Failed to create predicate checker: %v", err)
	//}
	listerRegistryStopChannel := make(chan struct{})
	listerRegistry := kube_util.NewListerRegistryWithDefaultListers(kubeClient, listerRegistryStopChannel)
	autoscaler, _ := core.NewAutoscaler(opts, predicateChecker, kubeClient, kubeEventRecorder, listerRegistry)
	//autoscaler, err := core.NewAutoscaler(opts, nil, nil, nil, nil)
	//if err != nil {
	//	glog.Fatalf("Failed to create autoscaler: %v", err)
	//}
	assert.NotNil(t, autoscaler)
}

type MultiStringFlag []string

var (
	nodeGroupsFlag         MultiStringFlag
	clusterName            = flag.String("cluster-name", "", "Autoscaled cluster name, if available")
	address                = flag.String("address", ":8085", "The address to expose prometheus metrics.")
	kubernetes             = flag.String("kubernetes", "", "Kubernetes master location. Leave blank for default")
	kubeConfigFile         = flag.String("kubeconfig", "", "Path to kubeconfig file with authorization and master location information.")
	cloudConfig            = flag.String("cloud-config", "", "The path to the cloud provider configuration file.  Empty string for no configuration file.")
	configMapName          = flag.String("configmap", "", "The name of the ConfigMap containing settings used for dynamic reconfiguration. Empty string for no ConfigMap.")
	namespace              = flag.String("namespace", "kube-system", "Namespace in which cluster-autoscaler run. If a --configmap flag is also provided, ensure that the configmap exists in this namespace before CA runs.")
	nodeGroupAutoDiscovery = flag.String("node-group-auto-discovery", "", "One or more definition(s) of node group auto-discovery. A definition is expressed `<name of discoverer per cloud provider>:[<key>[=<value>]]`. Only the `aws` cloud provider is currently supported. The only valid discoverer for it is `asg` and the valid key is `tag`. For example, specifying `--cloud-provider aws` and `--node-group-auto-discovery asg:tag=cluster-autoscaler/auto-discovery/enabled,kubernetes.io/cluster/<YOUR CLUSTER NAME>` results in ASGs tagged with `cluster-autoscaler/auto-discovery/enabled` and `kubernetes.io/cluster/<YOUR CLUSTER NAME>` to be considered as target node groups")
	scaleDownEnabled       = flag.Bool("scale-down-enabled", true, "Should CA scale down the cluster")
	scaleDownDelayAfterAdd = flag.Duration("scale-down-delay-after-add", 10*time.Minute,
		"How long after scale up that scale down evaluation resumes")
	scaleDownDelayAfterDelete = flag.Duration("scale-down-delay-after-delete", *scanInterval,
		"How long after node deletion that scale down evaluation resumes, defaults to scanInterval")
	scaleDownDelayAfterFailure = flag.Duration("scale-down-delay-after-failure", 3*time.Minute,
		"How long after scale down failure that scale down evaluation resumes")
	scaleDownUnneededTime = flag.Duration("scale-down-unneeded-time", 10*time.Minute,
		"How long a node should be unneeded before it is eligible for scale down")
	scaleDownUnreadyTime = flag.Duration("scale-down-unready-time", 20*time.Minute,
		"How long an unready node should be unneeded before it is eligible for scale down")
	scaleDownUtilizationThreshold = flag.Float64("scale-down-utilization-threshold", 0.5,
		"Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down")
	scaleDownNonEmptyCandidatesCount = flag.Int("scale-down-non-empty-candidates-count", 30,
		"Maximum number of non empty nodes considered in one iteration as candidates for scale down with drain."+
			"Lower value means better CA responsiveness but possible slower scale down latency."+
			"Higher value can affect CA performance with big clusters (hundreds of nodes)."+
			"Set to non posistive value to turn this heuristic off - CA will not limit the number of nodes it considers.")
	scaleDownCandidatesPoolRatio = flag.Float64("scale-down-candidates-pool-ratio", 0.1,
		"A ratio of nodes that are considered as additional non empty candidates for"+
			"scale down when some candidates from previous iteration are no longer valid."+
			"Lower value means better CA responsiveness but possible slower scale down latency."+
			"Higher value can affect CA performance with big clusters (hundreds of nodes)."+
			"Set to 1.0 to turn this heuristics off - CA will take all nodes as additional candidates.")
	scaleDownCandidatesPoolMinCount = flag.Int("scale-down-candidates-pool-min-count", 50,
		"Minimum number of nodes that are considered as additional non empty candidates"+
			"for scale down when some candidates from previous iteration are no longer valid."+
			"When calculating the pool size for additional candidates we take"+
			"max(#nodes * scale-down-candidates-pool-ratio, scale-down-candidates-pool-min-count).")
	scanInterval                = flag.Duration("scan-interval", 10*time.Second, "How often cluster is reevaluated for scale up or down")
	maxNodesTotal               = flag.Int("max-nodes-total", 0, "Maximum number of nodes in all node groups. Cluster autoscaler will not grow the cluster beyond this number.")
	coresTotal                  = flag.String("cores-total", minMaxFlagString(0, config.DefaultMaxClusterCores), "Minimum and maximum number of cores in cluster, in the format <min>:<max>. Cluster autoscaler will not scale the cluster beyond these numbers.")
	memoryTotal                 = flag.String("memory-total", minMaxFlagString(0, config.DefaultMaxClusterMemory), "Minimum and maximum number of gigabytes of memory in cluster, in the format <min>:<max>. Cluster autoscaler will not scale the cluster beyond these numbers.")
	cloudProviderFlag           = flag.String("cloud-provider", "gce", "Cloud provider type. Allowed values: gce, aws, kubemark")
	maxEmptyBulkDeleteFlag      = flag.Int("max-empty-bulk-delete", 10, "Maximum number of empty nodes that can be deleted at the same time.")
	maxGracefulTerminationFlag  = flag.Int("max-graceful-termination-sec", 10*60, "Maximum number of seconds CA waits for pod termination when trying to scale down a node.")
	maxTotalUnreadyPercentage   = flag.Float64("max-total-unready-percentage", 33, "Maximum percentage of unready nodes after which CA halts operations")
	okTotalUnreadyCount         = flag.Int("ok-total-unready-count", 3, "Number of allowed unready nodes, irrespective of max-total-unready-percentage")
	maxNodeProvisionTime        = flag.Duration("max-node-provision-time", 15*time.Minute, "Maximum time CA waits for node to be provisioned")
	unregisteredNodeRemovalTime = flag.Duration("unregistered-node-removal-time", 15*time.Minute, "Time that CA waits before removing nodes that are not registered in Kubernetes")

	estimatorFlag = flag.String("estimator", estimator.BinpackingEstimatorName,
		"Type of resource estimator to be used in scale up. Available values: ["+strings.Join(estimator.AvailableEstimators, ",")+"]")

	expanderFlag = flag.String("expander", expander.RandomExpanderName,
		"Type of node group expander to be used in scale up. Available values: ["+strings.Join(expander.AvailableExpanders, ",")+"]")

	writeStatusConfigMapFlag         = flag.Bool("write-status-configmap", true, "Should CA write status information to a configmap")
	maxInactivityTimeFlag            = flag.Duration("max-inactivity", 10*time.Minute, "Maximum time from last recorded autoscaler activity before automatic restart")
	maxFailingTimeFlag               = flag.Duration("max-failing-time", 15*time.Minute, "Maximum time from last recorded successful autoscaler run before automatic restart")
	balanceSimilarNodeGroupsFlag     = flag.Bool("balance-similar-node-groups", false, "Detect similar node groups and balance the number of nodes between them")
	nodeAutoprovisioningEnabled      = flag.Bool("node-autoprovisioning-enabled", false, "Should CA autoprovision node groups when needed")
	maxAutoprovisionedNodeGroupCount = flag.Int("max-autoprovisioned-node-group-count", 15, "The maximum number of autoprovisioned groups in the cluster.")

	expendablePodsPriorityCutoff = flag.Int("expendable-pods-priority_cutoff", 0, "Pods with priority below cutoff will be expendable. They can be killed without any consideration during scale down and they don't cause scale up. Pods with null priority (PodPriority disabled) are non expendable.")
)
func createKubeClient() kube_client.Interface {
	kubeConfig, err := config.GetKubeClientConfig(nil)
	if err != nil {
		//glog.Fatalf("Failed to build Kubernetes client configuration: %v", err)
	}
	return kube_client.NewForConfigOrDie(kubeConfig)
}
func parseMinMaxFlag(flag string) (int64, int64, error) {
	tokens := strings.SplitN(flag, ":", 2)
	if len(tokens) != 2 {
		return 0, 0, fmt.Errorf("wrong nodes configuration: %s", flag)
	}

	min, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to set min size: %s, expected integer, err: %v", tokens[0], err)
	}

	max, err := strconv.ParseInt(tokens[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to set max size: %s, expected integer, err: %v", tokens[1], err)
	}

	err = validateMinMaxFlag(min, max)
	if err != nil {
		return 0, 0, err
	}

	return min, max, nil
}
func validateMinMaxFlag(min, max int64) error {
	if min < 0 {
		return fmt.Errorf("min size must be greater or equal to  0")
	}
	if max < min {
		return fmt.Errorf("max size must be greater or equal to min size")
	}
	return nil
}

func minMaxFlagString(min, max int64) string {
	return fmt.Sprintf("%v:%v", min, max)
}
func createAutoscalerOptions() core.AutoscalerOptions {
	minCoresTotal, maxCoresTotal, err := parseMinMaxFlag(*coresTotal)
	if err != nil {
		//glog.Fatalf("Failed to parse flags: %v", err)
	}
	minMemoryTotal, maxMemoryTotal, err := parseMinMaxFlag(*memoryTotal)
	if err != nil {
		//glog.Fatalf("Failed to parse flags: %v", err)
	}
	// Convert memory limits to megabytes.
	minMemoryTotal = minMemoryTotal * 1024
	maxMemoryTotal = maxMemoryTotal * 1024

	autoscalingOpts := core.AutoscalingOptions{
		CloudConfig:                      *cloudConfig,
		CloudProviderName:                *cloudProviderFlag,
		NodeGroupAutoDiscovery:           *nodeGroupAutoDiscovery,
		MaxTotalUnreadyPercentage:        *maxTotalUnreadyPercentage,
		OkTotalUnreadyCount:              *okTotalUnreadyCount,
		EstimatorName:                    *estimatorFlag,
		ExpanderName:                     *expanderFlag,
		MaxEmptyBulkDelete:               *maxEmptyBulkDeleteFlag,
		MaxGracefulTerminationSec:        *maxGracefulTerminationFlag,
		MaxNodeProvisionTime:             *maxNodeProvisionTime,
		MaxNodesTotal:                    *maxNodesTotal,
		MaxCoresTotal:                    maxCoresTotal,
		MinCoresTotal:                    minCoresTotal,
		MaxMemoryTotal:                   maxMemoryTotal,
		MinMemoryTotal:                   minMemoryTotal,
		NodeGroups:                       nodeGroupsFlag,
		UnregisteredNodeRemovalTime:      *unregisteredNodeRemovalTime,
		ScaleDownDelayAfterAdd:           *scaleDownDelayAfterAdd,
		ScaleDownDelayAfterDelete:        *scaleDownDelayAfterDelete,
		ScaleDownDelayAfterFailure:       *scaleDownDelayAfterFailure,
		ScaleDownEnabled:                 *scaleDownEnabled,
		ScaleDownUnneededTime:            *scaleDownUnneededTime,
		ScaleDownUnreadyTime:             *scaleDownUnreadyTime,
		ScaleDownUtilizationThreshold:    *scaleDownUtilizationThreshold,
		ScaleDownNonEmptyCandidatesCount: *scaleDownNonEmptyCandidatesCount,
		ScaleDownCandidatesPoolRatio:     *scaleDownCandidatesPoolRatio,
		ScaleDownCandidatesPoolMinCount:  *scaleDownCandidatesPoolMinCount,
		WriteStatusConfigMap:             *writeStatusConfigMapFlag,
		BalanceSimilarNodeGroups:         *balanceSimilarNodeGroupsFlag,
		ConfigNamespace:                  *namespace,
		ClusterName:                      *clusterName,
		NodeAutoprovisioningEnabled:      *nodeAutoprovisioningEnabled,
		MaxAutoprovisionedNodeGroupCount: *maxAutoprovisionedNodeGroupCount,
		ExpendablePodsPriorityCutoff:     *expendablePodsPriorityCutoff,
	}

	configFetcherOpts := dynamic.ConfigFetcherOptions{
		ConfigMapName: *configMapName,
		Namespace:     *namespace,
	}

	return core.AutoscalerOptions{
		AutoscalingOptions:   autoscalingOpts,
		ConfigFetcherOptions: configFetcherOpts,
	}
}

