package cluster

import (
	"context"
	"fmt"

	"github.com/njhale/kink/api/v1alpha1"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

// NewStrategy creates and returns a clusterStrategy instance
func NewStrategy(typer runtime.ObjectTyper) clusterStrategy {
	return clusterStrategy{typer, names.SimpleNameGenerator}
}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not a Cluster
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiserver, ok := obj.(*v1alpha1.Cluster)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Cluster")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// MatchCluster is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchCluster(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *v1alpha1.Cluster) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type clusterStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (clusterStrategy) NamespaceScoped() bool {
	return false
}

func (clusterStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (clusterStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (clusterStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (clusterStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (clusterStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (clusterStrategy) Canonicalize(obj runtime.Object) {
}

func (clusterStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
