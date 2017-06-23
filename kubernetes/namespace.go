package kubernetes

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func (c *KubeClient) GetAllNamespaces() (namespaces *v1.NamespaceList, err error) {
	if namespaces, err = c.client.Namespaces().List(meta_v1.ListOptions{}); err != nil {
		return nil, err
	}
	return namespaces, nil
}
