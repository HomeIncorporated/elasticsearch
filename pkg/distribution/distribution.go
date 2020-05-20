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
package distribution

import (
	"fmt"

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"
	distapi "kubedb.dev/elasticsearch/pkg/distribution/api"
	"kubedb.dev/elasticsearch/pkg/distribution/elastic_stack"
	"kubedb.dev/elasticsearch/pkg/distribution/open_distro"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewElasticsearch(kc kubernetes.Interface, extClient cs.Interface, es *api.Elasticsearch) (distapi.ElasticsearchInterface, error) {
	if kc == nil {
		return nil, errors.New("Kubernetes client is empty")
	}
	if extClient == nil {
		return nil, errors.New("KubeDB client is empty")
	}
	if es == nil {
		return nil, errors.New("Elasticsearch object is empty")
	}

	v := es.Spec.Version
	esVersion, err := extClient.CatalogV1alpha1().ElasticsearchVersions().Get(v, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get elasticsearchVersion: %s", v))
	}

	if esVersion.Spec.AuthPlugin == catalog.ElasticsearchAuthPluginXpack {
		return elastic_stack.New(kc, extClient, es, esVersion), nil
	} else if esVersion.Spec.AuthPlugin == catalog.ElasticsearchAuthPluginOpenDistro {
		return open_distro.New(kc, extClient, es, esVersion), nil
	} else {
		return nil, errors.New("Unknown elasticsearch auth plugin")
	}
}