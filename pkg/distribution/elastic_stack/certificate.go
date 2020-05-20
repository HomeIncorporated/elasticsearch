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
package elastic_stack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha1"
	"kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1/util"
	certlib "kubedb.dev/elasticsearch/pkg/lib/cert"

	"github.com/appscode/go/crypto/rand"
	"gomodules.xyz/cert"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (es *Elasticsearch) EnsureCertSecret() error {
	if es.elasticsearch.Spec.DisableSecurity {
		return nil
	}

	certSecretVolumeSource := es.elasticsearch.Spec.CertificateSecret
	if certSecretVolumeSource == nil {
		var err error
		if certSecretVolumeSource, err = es.createCertSecret(); err != nil {
			return err
		}
		newES, _, err := util.PatchElasticsearch(es.extClient.KubedbV1alpha1(), es.elasticsearch, func(in *api.Elasticsearch) *api.Elasticsearch {
			in.Spec.CertificateSecret = certSecretVolumeSource
			return in
		})
		if err != nil {
			return err
		}
		es.elasticsearch = newES
	}
	return nil
}

func (es *Elasticsearch) createCertSecret() (*corev1.SecretVolumeSource, error) {
	certSecret, err := es.findCertSecret()
	if err != nil {
		return nil, err
	}

	if certSecret != nil {
		return &corev1.SecretVolumeSource{
			SecretName: certSecret.Name,
		}, nil
	}

	certPath := fmt.Sprintf("%v/%v", certlib.CertsDir, rand.Characters(3))
	if err := os.MkdirAll(certPath, os.ModePerm); err != nil {
		return nil, err
	}

	caKey, caCert, pass, err := certlib.CreateCaCertificate(certPath)
	if err != nil {
		return nil, err
	}
	err = certlib.CreateNodeCertificate(certPath, es.elasticsearch, caKey, caCert, pass)
	if err != nil {
		return nil, err
	}

	root, err := ioutil.ReadFile(filepath.Join(certPath, certlib.RootKeyStore))
	if err != nil {
		return nil, err
	}
	node, err := ioutil.ReadFile(filepath.Join(certPath, certlib.NodeKeyStore))
	if err != nil {
		return nil, err
	}

	data := map[string][]byte{
		certlib.RootKeyStore: root,
		certlib.NodeKeyStore: node,
	}

	if err := certlib.CreateClientCertificate(certPath, es.elasticsearch, caKey, caCert, pass); err != nil {
		return nil, err
	}

	client, err := ioutil.ReadFile(filepath.Join(certPath, certlib.ClientKeyStore))
	if err != nil {
		return nil, err
	}

	data[certlib.RootCert] = cert.EncodeCertPEM(caCert)
	data[certlib.ClientKeyStore] = client

	name := fmt.Sprintf("%v-cert", es.elasticsearch.OffshootName())
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: es.elasticsearch.OffshootLabels(),
		},
		Type: corev1.SecretTypeOpaque,
		Data: data,
		StringData: map[string]string{
			"key_pass": pass,
		},
	}
	if _, err := es.kClient.CoreV1().Secrets(es.elasticsearch.Namespace).Create(secret); err != nil {
		return nil, err
	}

	secretVolumeSource := &corev1.SecretVolumeSource{
		SecretName: secret.Name,
	}

	return secretVolumeSource, nil
}

func (es *Elasticsearch) findCertSecret() (*corev1.Secret, error) {
	name := fmt.Sprintf("%v-cert", es.elasticsearch.OffshootName())

	secret, err := es.kClient.CoreV1().Secrets(es.elasticsearch.Namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if secret.Labels[api.LabelDatabaseKind] != api.ResourceKindElasticsearch ||
		secret.Labels[api.LabelDatabaseName] != es.elasticsearch.Name {
		return nil, fmt.Errorf(`intended secret "%v/%v" already exists`, es.elasticsearch.Namespace, name)
	}

	return secret, nil
}