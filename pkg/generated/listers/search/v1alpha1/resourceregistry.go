// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/karmada-io/karmada/pkg/apis/search/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ResourceRegistryLister helps list ResourceRegistries.
// All objects returned here must be treated as read-only.
type ResourceRegistryLister interface {
	// List lists all ResourceRegistries in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ResourceRegistry, err error)
	// ResourceRegistries returns an object that can list and get ResourceRegistries.
	ResourceRegistries(namespace string) ResourceRegistryNamespaceLister
	ResourceRegistryListerExpansion
}

// resourceRegistryLister implements the ResourceRegistryLister interface.
type resourceRegistryLister struct {
	indexer cache.Indexer
}

// NewResourceRegistryLister returns a new ResourceRegistryLister.
func NewResourceRegistryLister(indexer cache.Indexer) ResourceRegistryLister {
	return &resourceRegistryLister{indexer: indexer}
}

// List lists all ResourceRegistries in the indexer.
func (s *resourceRegistryLister) List(selector labels.Selector) (ret []*v1alpha1.ResourceRegistry, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ResourceRegistry))
	})
	return ret, err
}

// ResourceRegistries returns an object that can list and get ResourceRegistries.
func (s *resourceRegistryLister) ResourceRegistries(namespace string) ResourceRegistryNamespaceLister {
	return resourceRegistryNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ResourceRegistryNamespaceLister helps list and get ResourceRegistries.
// All objects returned here must be treated as read-only.
type ResourceRegistryNamespaceLister interface {
	// List lists all ResourceRegistries in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ResourceRegistry, err error)
	// Get retrieves the ResourceRegistry from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ResourceRegistry, error)
	ResourceRegistryNamespaceListerExpansion
}

// resourceRegistryNamespaceLister implements the ResourceRegistryNamespaceLister
// interface.
type resourceRegistryNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ResourceRegistries in the indexer for a given namespace.
func (s resourceRegistryNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ResourceRegistry, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ResourceRegistry))
	})
	return ret, err
}

// Get retrieves the ResourceRegistry from the indexer for a given namespace and name.
func (s resourceRegistryNamespaceLister) Get(name string) (*v1alpha1.ResourceRegistry, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("resourceregistry"), name)
	}
	return obj.(*v1alpha1.ResourceRegistry), nil
}
