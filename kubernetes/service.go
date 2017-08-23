package kubernetes

import (
	_ "log"
	"regexp"

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

var publicELBNameRegexp = regexp.MustCompile(`^(.*?)-[a-z0-9]*\..*.amazonaws.com$`)
var internalELBNameRegexp = regexp.MustCompile(`^internal-(.*?)$`)

func getELBName(elbHost string) (elbName string, err error) {
	elbName = string(publicELBNameRegexp.FindSubmatch([]byte(elbHost))[1])
	internal, err := regexp.MatchString(`^internal-(.*?)$`, elbName)
	if err != nil {
		return "", err
	}
	if internal {
		elbName = string(internalELBNameRegexp.FindSubmatch([]byte(elbName))[1])
	}
	return elbName, nil
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
						Key:   `kube_` + key,
						Value: value,
					})
				}

				elbName, err := getELBName(s.Status.LoadBalancer.Ingress[0].Hostname)
				if err != nil {
					return nil, err
				}
				// log.Println(elbName)

				services = append(services,
					Service{
						KubeName:            s.Name,
						KubeNameSpace:       s.Namespace,
						KubernetesCluster:   s.ClusterName,
						Name:                elbName,
						LoadBalancerIngress: s.Status.LoadBalancer.Ingress[0].Hostname,
						Labels:              labels,
					})
			}
		}
	}

	return services, nil
}

func (c *KubeClient) UpdateLabelsToDataDogFormat(elbTags []Label, service Service) Service {
	service.Labels = append(service.Labels, Label{
		Key:   "kube_service",
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
