// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package k8s

import (
	"fmt"
	"sync"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/cilium/cilium/pkg/hive/cell"
	cilium_api_v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	cilium_api_v2alpha1 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2alpha1"
	"github.com/cilium/cilium/pkg/k8s/client"
	"github.com/cilium/cilium/pkg/k8s/resource"
	slim_corev1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/core/v1"
	slim_discoveryv1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/discovery/v1"
	slim_discoveryv1beta1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/discovery/v1beta1"
	slim_networkingv1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/api/networking/v1"
	"github.com/cilium/cilium/pkg/k8s/types"
	"github.com/cilium/cilium/pkg/k8s/utils"
	"github.com/cilium/cilium/pkg/k8s/version"
	"github.com/cilium/cilium/pkg/node"
)

// Config defines the configuration options for k8s resources.
type Config struct {
	EnableK8sEndpointSlice bool

	// K8sServiceProxyName is the value of service.kubernetes.io/service-proxy-name label,
	// that identifies the service objects Cilium should handle.
	// If the provided value is an empty string, Cilium will manage service objects when
	// the label is not present. For more details -
	// https://github.com/kubernetes/enhancements/tree/master/keps/sig-network/2447-Make-kube-proxy-service-abstraction-optional
	K8sServiceProxyName string

	// K8sListTimeoutSeconds is the timeout for listing resources from k8s.
	K8sListTimeoutSeconds int64
}

// DefaultConfig represents the default k8s resources config values.
var DefaultConfig = Config{
	EnableK8sEndpointSlice: true,
}

// Flags implements the cell.Flagger interface.
func (def Config) Flags(flags *pflag.FlagSet) {
	flags.Bool("enable-k8s-endpoint-slice", def.EnableK8sEndpointSlice, "Enables k8s EndpointSlice feature in Cilium if the k8s cluster supports it")
	flags.String("k8s-service-proxy-name", def.K8sServiceProxyName, "Value of K8s service-proxy-name label for which Cilium handles the services (empty = all services without service.kubernetes.io/service-proxy-name label)")
	flags.Int64("k8s-list-timeout-seconds", def.K8sListTimeoutSeconds, "The maximum duration to wait before giving up on a list operation")
}

// ListOptionOut defines the parameters for building k8s resources.
type ListOptionOut struct {
	cell.Out
	// ListOptions is the list options to be used for all k8s resources.
	Opts func(*metav1.ListOptions) `group:"opts"`
}

// SetListTimeout sets the list timeout for all k8s resources.
func SetListTimeout(def Config) ListOptionOut {
	return ListOptionOut{
		Opts: func(opts *metav1.ListOptions) {
			if def.K8sListTimeoutSeconds > 0 {
				opts.TimeoutSeconds = &def.K8sListTimeoutSeconds
			}
		},
	}
}

// ResourceParam defines the parameters for building k8s resources.
type ResourceParam struct {
	cell.In
	Lc cell.Lifecycle

	Cs client.Clientset
	Cfg Config

	// ListOptions is the list options to be used for all k8s resources.
	Opts []func(*metav1.ListOptions) `group:"opts"`
}

// ServiceResource builds the Resource[Service] object.
func ServiceResource(params ResourceParam) (resource.Resource[*slim_corev1.Service], error) {
	if !params.Cs.IsEnabled() {
		return nil, nil
	}
	optsModifier, err := utils.GetServiceAndEndpointListOptionsModifier(params.Cfg.K8sServiceProxyName)
	if err != nil {
		return nil, err
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*slim_corev1.ServiceList](params.Cs.Slim().CoreV1().Services("")),
		append(params.Opts, optsModifier)...,
	)
	return resource.New[*slim_corev1.Service](params.Lc, lw, resource.WithMetric("Service")), nil
}

func NodeResource(lc cell.Lifecycle, cs client.Clientset, opts ...func(*metav1.ListOptions)) (resource.Resource[*slim_corev1.Node], error) {
	if !cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*slim_corev1.NodeList](cs.Slim().CoreV1().Nodes()),
		opts...,
	)
	return resource.New[*slim_corev1.Node](lc, lw, resource.WithMetric("Node")), nil
}

func CiliumNodeResource(params ResourceParam) (resource.Resource[*cilium_api_v2.CiliumNode], error) {
	if !params.Cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*cilium_api_v2.CiliumNodeList](params.Cs.CiliumV2().CiliumNodes()),
		params.Opts...,
	)
	return resource.New[*cilium_api_v2.CiliumNode](params.Lc, lw, resource.WithMetric("CiliumNode")), nil
}

func PodResource(lc cell.Lifecycle, cs client.Clientset, opts ...func(*metav1.ListOptions)) (resource.Resource[*slim_corev1.Pod], error) {
	if !cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*slim_corev1.PodList](cs.Slim().CoreV1().Pods("")),
		opts...,
	)
	return resource.New[*slim_corev1.Pod](lc, lw, resource.WithMetric("Pod")), nil
}

func NamespaceResource(lc cell.Lifecycle, cs client.Clientset, opts ...func(*metav1.ListOptions)) (resource.Resource[*slim_corev1.Namespace], error) {
	if !cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*slim_corev1.NamespaceList](cs.Slim().CoreV1().Namespaces()),
		opts...,
	)
	return resource.New[*slim_corev1.Namespace](lc, lw, resource.WithMetric("Namespace")), nil
}

func LBIPPoolsResource(lc cell.Lifecycle, cs client.Clientset, opts ...func(*metav1.ListOptions)) (resource.Resource[*cilium_api_v2alpha1.CiliumLoadBalancerIPPool], error) {
	if !cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*cilium_api_v2alpha1.CiliumLoadBalancerIPPoolList](cs.CiliumV2alpha1().CiliumLoadBalancerIPPools()),
		opts...,
	)
	return resource.New[*cilium_api_v2alpha1.CiliumLoadBalancerIPPool](lc, lw, resource.WithMetric("CiliumLoadBalancerIPPool")), nil
}

func CiliumIdentityResource(lc cell.Lifecycle, cs client.Clientset, opts ...func(*metav1.ListOptions)) (resource.Resource[*cilium_api_v2.CiliumIdentity], error) {
	if !cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*cilium_api_v2.CiliumIdentityList](cs.CiliumV2().CiliumIdentities()),
		opts...,
	)
	return resource.New[*cilium_api_v2.CiliumIdentity](lc, lw, resource.WithMetric("CiliumIdentityList")), nil
}

func NetworkPolicyResource(params ResourceParam) (resource.Resource[*slim_networkingv1.NetworkPolicy], error) {
	if !params.Cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*slim_networkingv1.NetworkPolicyList](params.Cs.Slim().NetworkingV1().NetworkPolicies("")),
		params.Opts...,
	)
	return resource.New[*slim_networkingv1.NetworkPolicy](params.Lc, lw, resource.WithMetric("NetworkPolicy")), nil
}

func CiliumNetworkPolicyResource(params ResourceParam) (resource.Resource[*cilium_api_v2.CiliumNetworkPolicy], error) {
	if !params.Cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*cilium_api_v2.CiliumNetworkPolicyList](params.Cs.CiliumV2().CiliumNetworkPolicies("")),
		params.Opts...,
	)
	return resource.New[*cilium_api_v2.CiliumNetworkPolicy](params.Lc, lw, resource.WithMetric("CiliumNetworkPolicy")), nil
}

func CiliumClusterwideNetworkPolicyResource(params ResourceParam) (resource.Resource[*cilium_api_v2.CiliumClusterwideNetworkPolicy], error) {
	if !params.Cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*cilium_api_v2.CiliumClusterwideNetworkPolicyList](params.Cs.CiliumV2().CiliumClusterwideNetworkPolicies()),
		params.Opts...,
	)
	return resource.New[*cilium_api_v2.CiliumClusterwideNetworkPolicy](params.Lc, lw, resource.WithMetric("CiliumClusterwideNetworkPolicy")), nil
}

func CiliumCIDRGroupResource(params ResourceParam) (resource.Resource[*cilium_api_v2alpha1.CiliumCIDRGroup], error) {
	if !params.Cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*cilium_api_v2alpha1.CiliumCIDRGroupList](params.Cs.CiliumV2alpha1().CiliumCIDRGroups()),
		params.Opts...,
	)
	return resource.New[*cilium_api_v2alpha1.CiliumCIDRGroup](params.Lc, lw, resource.WithMetric("CiliumCIDRGroup")), nil
}

func CiliumPodIPPoolResource(lc cell.Lifecycle, cs client.Clientset, opts ...func(*metav1.ListOptions)) (resource.Resource[*cilium_api_v2alpha1.CiliumPodIPPool], error) {
	if !cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*cilium_api_v2alpha1.CiliumPodIPPoolList](cs.CiliumV2alpha1().CiliumPodIPPools()),
		opts...,
	)
	return resource.New[*cilium_api_v2alpha1.CiliumPodIPPool](lc, lw, resource.WithMetric("CiliumPodIPPool")), nil
}

func EndpointsResource(params ResourceParam) (resource.Resource[*Endpoints], error) {
	if !params.Cs.IsEnabled() {
		return nil, nil
	}
	endpointsOptsModifier, err := utils.GetServiceAndEndpointListOptionsModifier(params.Cfg.K8sServiceProxyName)
	if err != nil {
		return nil, err
	}

	endpointSliceOpsModifier, err := utils.GetEndpointSliceListOptionsModifier()
	if err != nil {
		return nil, err
	}
	lw := &endpointsListerWatcher{
		cs:                         params.Cs,
		enableK8sEndpointSlice:     params.Cfg.EnableK8sEndpointSlice,
		endpointsOptsModifier:      endpointsOptsModifier,
		endpointSlicesOptsModifier: endpointSliceOpsModifier,
	}
	return resource.New[*Endpoints](
		params.Lc,
		lw,
		resource.WithLazyTransform(lw.getSourceObj, transformEndpoint),
		resource.WithMetric("Endpoint"),
		resource.WithName("endpoints"),
	), nil
}

// endpointsListerWatcher implements List and Watch for endpoints/endpointslices. It
// performs the capability check on first call to List/Watch. This allows constructing
// the resource before the client has been started and capabilities have been probed.
type endpointsListerWatcher struct {
	cs                         client.Clientset
	enableK8sEndpointSlice     bool
	endpointsOptsModifier      func(*metav1.ListOptions)
	endpointSlicesOptsModifier func(*metav1.ListOptions)
	sourceObj                  k8sRuntime.Object

	once                sync.Once
	cachedListerWatcher cache.ListerWatcher
}

func (lw *endpointsListerWatcher) getSourceObj() k8sRuntime.Object {
	lw.getListerWatcher() // force the construction
	return lw.sourceObj
}

func (lw *endpointsListerWatcher) getListerWatcher() cache.ListerWatcher {
	lw.once.Do(func() {
		if lw.enableK8sEndpointSlice && version.Capabilities().EndpointSlice {
			if version.Capabilities().EndpointSliceV1 {
				log.Info("Using discoveryv1.EndpointSlice")
				lw.cachedListerWatcher = utils.ListerWatcherFromTyped[*slim_discoveryv1.EndpointSliceList](
					lw.cs.Slim().DiscoveryV1().EndpointSlices(""),
				)
				lw.sourceObj = &slim_discoveryv1.EndpointSlice{}
			} else {
				log.Info("Using discoveryv1beta1.EndpointSlice")
				lw.cachedListerWatcher = utils.ListerWatcherFromTyped[*slim_discoveryv1beta1.EndpointSliceList](
					lw.cs.Slim().DiscoveryV1beta1().EndpointSlices(""),
				)
				lw.sourceObj = &slim_discoveryv1beta1.EndpointSlice{}
			}
			lw.cachedListerWatcher = utils.ListerWatcherWithModifier(lw.cachedListerWatcher, lw.endpointSlicesOptsModifier)
		} else {
			log.Info("Using v1.Endpoints")
			lw.cachedListerWatcher = utils.ListerWatcherFromTyped[*slim_corev1.EndpointsList](
				lw.cs.Slim().CoreV1().Endpoints(""),
			)
			lw.sourceObj = &slim_corev1.Endpoints{}
			lw.cachedListerWatcher = utils.ListerWatcherWithModifier(lw.cachedListerWatcher, lw.endpointsOptsModifier)
		}
	})
	return lw.cachedListerWatcher
}

func (lw *endpointsListerWatcher) List(opts metav1.ListOptions) (k8sRuntime.Object, error) {
	return lw.getListerWatcher().List(opts)
}

func (lw *endpointsListerWatcher) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return lw.getListerWatcher().Watch(opts)
}

func transformEndpoint(obj any) (any, error) {
	switch obj := obj.(type) {
	case *slim_corev1.Endpoints:
		return ParseEndpoints(obj), nil
	case *slim_discoveryv1.EndpointSlice:
		return ParseEndpointSliceV1(obj), nil
	case *slim_discoveryv1beta1.EndpointSlice:
		return ParseEndpointSliceV1Beta1(obj), nil
	default:
		return nil, fmt.Errorf("%T not a known endpoint or endpoint slice object", obj)
	}
}

// CiliumSlimEndpointResource uses the "localNode" IndexFunc to build the resource indexer.
// The IndexFunc accesses the local node info to get its IP, so it depends on the local node store
// to initialize it before the first access.
// To reflect this, the node.LocalNodeStore dependency is explicitly requested in the function
// signature.
func CiliumSlimEndpointResource(params ResourceParam, _ *node.LocalNodeStore) (resource.Resource[*types.CiliumEndpoint], error) {
	if !params.Cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*cilium_api_v2.CiliumEndpointList](params.Cs.CiliumV2().CiliumEndpoints(slim_corev1.NamespaceAll)),
		params.Opts...,
	)
	indexers := cache.Indexers{
		"localNode": ciliumEndpointLocalPodIndexFunc,
	}
	return resource.New[*types.CiliumEndpoint](params.Lc, lw,
		resource.WithLazyTransform(func() runtime.Object {
			return &cilium_api_v2.CiliumEndpoint{}
		}, TransformToCiliumEndpoint),
		resource.WithMetric("CiliumEndpoint"),
		resource.WithIndexers(indexers),
		resource.WithStoppableInformer(),
	), nil
}

// ciliumEndpointLocalPodIndexFunc is an IndexFunc that indexes only local
// CiliumEndpoints, by their local Node IP.
func ciliumEndpointLocalPodIndexFunc(obj any) ([]string, error) {
	cep, ok := obj.(*types.CiliumEndpoint)
	if !ok {
		return nil, fmt.Errorf("unexpected object type: %T", obj)
	}
	indices := []string{}
	if cep.Networking == nil {
		log.WithField("ciliumendpoint", cep.GetNamespace()+"/"+cep.GetName()).
			Debug("cannot index CiliumEndpoint by node without network status")
		return nil, nil
	}
	if cep.Networking.NodeIP == node.GetCiliumEndpointNodeIP() {
		indices = append(indices, cep.Networking.NodeIP)
	}
	return indices, nil
}

// CiliumEndpointSliceResource uses the "localNode" IndexFunc to build the resource indexer.
// The IndexFunc accesses the local node info to get its IP, so it depends on the local node store
// to initialize it before the first access.
// To reflect this, the node.LocalNodeStore dependency is explicitly requested in the function
// signature.
func CiliumEndpointSliceResource(params ResourceParam, _ *node.LocalNodeStore) (resource.Resource[*cilium_api_v2alpha1.CiliumEndpointSlice], error) {
	if !params.Cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*cilium_api_v2alpha1.CiliumEndpointSliceList](params.Cs.CiliumV2alpha1().CiliumEndpointSlices()),
		params.Opts...,
	)
	indexers := cache.Indexers{
		"localNode": ciliumEndpointSliceLocalPodIndexFunc,
	}
	return resource.New[*cilium_api_v2alpha1.CiliumEndpointSlice](params.Lc, lw,
		resource.WithMetric("CiliumEndpointSlice"),
		resource.WithIndexers(indexers),
		resource.WithStoppableInformer(),
	), nil
}

// ciliumEndpointSliceLocalPodIndexFunc is an IndexFunc that indexes CiliumEndpointSlices
// by their corresponding Pod, which are running locally on this Node.
func ciliumEndpointSliceLocalPodIndexFunc(obj any) ([]string, error) {
	ces, ok := obj.(*cilium_api_v2alpha1.CiliumEndpointSlice)
	if !ok {
		return nil, fmt.Errorf("unexpected object type: %T", obj)
	}
	indices := []string{}
	for _, ep := range ces.Endpoints {
		if ep.Networking.NodeIP == node.GetCiliumEndpointNodeIP() {
			indices = append(indices, ep.Networking.NodeIP)
			break
		}
	}
	return indices, nil
}

func CiliumExternalWorkloads(lc cell.Lifecycle, cs client.Clientset, opts ...func(*metav1.ListOptions)) (resource.Resource[*cilium_api_v2.CiliumExternalWorkload], error) {
	if !cs.IsEnabled() {
		return nil, nil
	}
	lw := utils.ListerWatcherWithModifiers(
		utils.ListerWatcherFromTyped[*cilium_api_v2.CiliumExternalWorkloadList](cs.CiliumV2().CiliumExternalWorkloads()),
		opts...,
	)
	return resource.New[*cilium_api_v2.CiliumExternalWorkload](lc, lw, resource.WithMetric("CiliumExternalWorkloads")), nil
}
