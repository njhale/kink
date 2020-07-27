package cluster

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericfeatures "k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/registry/rest"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	restclient "k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/njhale/kink/api/v1alpha1"
)

type ClusterAPIProxy struct {
	client client.Client
}

func NewClusterAPIProxy(c client.Client) (*ClusterAPIProxy, error) {
	return &ClusterAPIProxy{
		client: c,
	}, nil

}

var _ rest.Connecter = &ClusterAPIProxy{}

func (c *ClusterAPIProxy) Connect(ctx context.Context, name string, options runtime.Object, responder rest.Responder) (http.Handler, error) {
	// TODO(njhale): finish this up
	// References:
	// https://github.com/kubernetes/kube-aggregator/blob/b29123120c84125198a9abc1aa5a9c4e521a628f/pkg/apiserver/handler_proxy.go#L109
	// https://github.com/kubernetes-sigs/apiserver-builder-alpha/blob/a3718943ef91ffd65400e23bd71e80b7a372f87e/example/podexec/pkg/apis/podexec/exec_pod_rest.go#L90
	cluster, err := c.getCluster(ctx, name)
	if err != nil {
		return nil, err
	}

	location, transport, err := c.location(ctx, cluster)
	if err != nil {
		return nil, err
	}

	proxyHandler := newThrottledUpgradeAwareProxyHandler(location, transport, false, true, true, responder)

	return proxyHandler, nil
}

func (c *ClusterAPIProxy) location(ctx context.Context, cluster *v1alpha1.Cluster) (*url.URL, http.RoundTripper, error) {
	// FIXME(njhale): get the rest configuration and user proxy auth right

	serviceRef := cluster.Status.ServiceRef
	if serviceRef == nil {
		return nil, nil, fmt.Errorf("missing cluster api service")
	}

	config := &restclient.Config{
		TLSClientConfig: restclient.TLSClientConfig{
			Insecure:   true,
			ServerName: fmt.Sprintf("%s.%s.svc", serviceRef.Name, serviceRef.Namespace),
			// CertData:   proxyClientCert,
			// KeyData:    proxyClientKey,
			// CAData:     apiService.Spec.CABundle,
		},
	}

	transport, err := restclient.TransportFor(config)
	if err != nil {
		return nil, nil, err
	}

	location := &url.URL{
		Scheme: "https",
		Host:   config.TLSClientConfig.ServerName,
	}

	return location, transport, nil
}

func newThrottledUpgradeAwareProxyHandler(location *url.URL, transport http.RoundTripper, wrapTransport, upgradeRequired, interceptRedirects bool, responder rest.Responder) *proxy.UpgradeAwareHandler {
	handler := proxy.NewUpgradeAwareHandler(location, transport, wrapTransport, upgradeRequired, proxy.NewErrorResponder(responder))
	handler.InterceptRedirects = interceptRedirects && utilfeature.DefaultFeatureGate.Enabled(genericfeatures.StreamingProxyRedirects)
	handler.RequireSameHostRedirects = utilfeature.DefaultFeatureGate.Enabled(genericfeatures.ValidateProxyRedirects)
	handler.MaxBytesPerSec = 0
	return handler
}

func (c *ClusterAPIProxy) getService(ctx context.Context, cluster *v1alpha1.Cluster) (*corev1.Service, error) {
	// TODO(njhale): implement
	return nil, nil
}

func (c *ClusterAPIProxy) getCluster(ctx context.Context, name string) (*v1alpha1.Cluster, error) {
	ns, found := request.NamespaceFrom(ctx)
	if !found {
		return nil, fmt.Errorf("namespace missing from request")
	}

	cluster := &v1alpha1.Cluster{}
	if err := c.client.Get(ctx, types.NamespacedName{Namespace: ns, Name: name}, cluster); err != nil {
		return nil, err
	}

	return cluster, nil
}
