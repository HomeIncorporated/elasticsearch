/*
Copyright The KubeDB Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PgBouncerOpsRequestLister helps list PgBouncerOpsRequests.
type PgBouncerOpsRequestLister interface {
	// List lists all PgBouncerOpsRequests in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.PgBouncerOpsRequest, err error)
	// PgBouncerOpsRequests returns an object that can list and get PgBouncerOpsRequests.
	PgBouncerOpsRequests(namespace string) PgBouncerOpsRequestNamespaceLister
	PgBouncerOpsRequestListerExpansion
}

// pgBouncerOpsRequestLister implements the PgBouncerOpsRequestLister interface.
type pgBouncerOpsRequestLister struct {
	indexer cache.Indexer
}

// NewPgBouncerOpsRequestLister returns a new PgBouncerOpsRequestLister.
func NewPgBouncerOpsRequestLister(indexer cache.Indexer) PgBouncerOpsRequestLister {
	return &pgBouncerOpsRequestLister{indexer: indexer}
}

// List lists all PgBouncerOpsRequests in the indexer.
func (s *pgBouncerOpsRequestLister) List(selector labels.Selector) (ret []*v1alpha1.PgBouncerOpsRequest, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.PgBouncerOpsRequest))
	})
	return ret, err
}

// PgBouncerOpsRequests returns an object that can list and get PgBouncerOpsRequests.
func (s *pgBouncerOpsRequestLister) PgBouncerOpsRequests(namespace string) PgBouncerOpsRequestNamespaceLister {
	return pgBouncerOpsRequestNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PgBouncerOpsRequestNamespaceLister helps list and get PgBouncerOpsRequests.
type PgBouncerOpsRequestNamespaceLister interface {
	// List lists all PgBouncerOpsRequests in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.PgBouncerOpsRequest, err error)
	// Get retrieves the PgBouncerOpsRequest from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.PgBouncerOpsRequest, error)
	PgBouncerOpsRequestNamespaceListerExpansion
}

// pgBouncerOpsRequestNamespaceLister implements the PgBouncerOpsRequestNamespaceLister
// interface.
type pgBouncerOpsRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all PgBouncerOpsRequests in the indexer for a given namespace.
func (s pgBouncerOpsRequestNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.PgBouncerOpsRequest, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.PgBouncerOpsRequest))
	})
	return ret, err
}

// Get retrieves the PgBouncerOpsRequest from the indexer for a given namespace and name.
func (s pgBouncerOpsRequestNamespaceLister) Get(name string) (*v1alpha1.PgBouncerOpsRequest, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("pgbounceropsrequest"), name)
	}
	return obj.(*v1alpha1.PgBouncerOpsRequest), nil
}