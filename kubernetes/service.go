package kubernetes

import (
	_ "log"
	"strings"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

type Service struct {
	KubeName            string
	KubeNameSpace       string
	KubernetesCluster   string
	LoadBalancerIngress string
	Name                string
	Labels              []Label
}

type Label struct {
	Key   string
	Value string
}

func (c *KubeClient) GetAllServices(namespaces *v1.NamespaceList) (services []Service, err error) {

	for _, n := range namespaces.Items {
		// log.Println(n.ObjectMeta.Name)
		service, err := c.client.Services(n.ObjectMeta.Name).List(meta_v1.ListOptions{})
		if err != nil {
			return nil, err
		}
		// log.Println(service)

		for _, s := range service.Items {
			if len(s.Status.LoadBalancer.Ingress) != 0 {
				labels := []Label{}
				for key, value := range s.ObjectMeta.Labels {
					labels = append(labels, Label{
						Key:   key,
						Value: value,
					})
				}
				services = append(services,
					Service{
						KubeName:            s.Name,
						KubeNameSpace:       s.Namespace,
						KubernetesCluster:   s.ClusterName,
						Name:                strings.Split(s.Status.LoadBalancer.Ingress[0].Hostname, "-")[0],
						LoadBalancerIngress: s.Status.LoadBalancer.Ingress[0].Hostname,
						Labels:              labels,
					})
			}
		}
	}

	return services, nil
}

func (c *KubeClient) UpdateLabelsKubernetesCluster(elbTags []Label, service Service) Service {
	service.Labels = append(service.Labels, Label{
		Key:   "kube_name",
		Value: service.KubeName,
	})
	service.Labels = append(service.Labels, Label{
		Key:   "kube_namespace",
		Value: service.KubeNameSpace,
	})
	for _, t := range elbTags {
		if t.Key == "KubernetesCluster" {
			service.Labels = append(service.Labels, Label{
				Key:   "kubernetescluster",
				Value: t.Value,
			})
		}
	}
	return service
}
