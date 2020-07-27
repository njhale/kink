package cluster

import (
	"github.com/njhale/kink/api/v1alpha1"
	"github.com/njhale/kink/registry"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*registry.REST, error) {
	strategy := NewStrategy(scheme)
	resource := v1alpha1.GroupVersion.WithResource("clusters").GroupResource()

	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &v1alpha1.Cluster{} },
		NewListFunc:              func() runtime.Object { return &v1alpha1.ClusterList{} },
		PredicateFunc:            MatchCluster,
		DefaultQualifiedResource: resource,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: rest.NewDefaultTableConvertor(resource),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &registry.REST{store}, nil
}
