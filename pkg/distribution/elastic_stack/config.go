package elastic_stack

import (
	"fmt"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_util "kmodules.xyz/client-go/core/v1"
)

const (
	ConfigFileName          = "elasticsearch.yml"
	ConfigFileMountPath     = "/usr/share/elasticsearch/config"
	TempConfigFileMountPath = "/elasticsearch/temp-config"
	DatabaseConfigMapSuffix = "config"
)

var xpack_config = `
xpack.security.enabled: true

xpack.security.transport.ssl.enabled: true
xpack.security.transport.ssl.verification_mode: certificate
xpack.security.transport.ssl.keystore.path: /usr/share/elasticsearch/config/certs/node.jks
xpack.security.transport.ssl.keystore.password: ${KEY_PASS}
xpack.security.transport.ssl.truststore.path: /usr/share/elasticsearch/config/certs/root.jks
xpack.security.transport.ssl.truststore.password: ${KEY_PASS}

xpack.security.http.ssl.keystore.path: /usr/share/elasticsearch/config/certs/client.jks
xpack.security.http.ssl.keystore.password: ${KEY_PASS}
xpack.security.http.ssl.truststore.path: /usr/share/elasticsearch/config/certs/root.jks
xpack.security.http.ssl.truststore.password: ${KEY_PASS}
`

func (es *Elasticsearch) EnsureDefaultConfig() error {
	if !es.elasticsearch.Spec.DisableSecurity {
		if err := es.findDefaultConfig(); err != nil {
			return err
		}

		cmMeta := metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v-%v", es.elasticsearch.OffshootName(), DatabaseConfigMapSuffix),
			Namespace: es.elasticsearch.Namespace,
		}
		owner := metav1.NewControllerRef(es.elasticsearch, api.SchemeGroupVersion.WithKind(api.ResourceKindElasticsearch))

		if _, _, err := core_util.CreateOrPatchConfigMap(es.kClient, cmMeta, func(in *corev1.ConfigMap) *corev1.ConfigMap {
			in.Labels = core_util.UpsertMap(in.Labels, es.elasticsearch.OffshootLabels())
			core_util.EnsureOwnerReference(&in.ObjectMeta, owner)
			in.Data = map[string]string{
				ConfigFileName: xpack_config,
			}
			return in
		}); err != nil {
			return err
		}
	}
	return nil
}

func (es *Elasticsearch) findDefaultConfig() error {
	cmName := fmt.Sprintf("%v-%v", es.elasticsearch.OffshootName(), DatabaseConfigMapSuffix)

	configMap, err := es.kClient.CoreV1().ConfigMaps(es.elasticsearch.Namespace).Get(cmName, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	if configMap.Labels[api.LabelDatabaseKind] != api.ResourceKindElasticsearch &&
		configMap.Labels[api.LabelDatabaseName] != es.elasticsearch.Name {
		return fmt.Errorf(`intended configMap "%v/%v" already exists`, es.elasticsearch.Namespace, cmName)
	}

	return nil
}