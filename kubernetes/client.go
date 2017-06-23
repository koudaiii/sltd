package kubernetes

import (
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	client *kubernetes.Clientset
}

func isNotExists() bool {
	file := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	_, err := os.Stat(file)
	return os.IsNotExist(err)
}

func NewKubeClient(inCluster bool) *KubeClient {
	if isNotExists() && !(inCluster) {
		log.Fatalln("kubeconfig invalid and --in-cluster is false; kubeconfig must be set to a valid file(kubeconfig default file name: $HOME/.kube/config)")
		os.Exit(1)
	}

	if inCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		tokenPresent := false
		if len(config.BearerToken) > 0 {
			tokenPresent = true
		}
		log.Printf("service account token present: %v", tokenPresent)
		log.Printf("service host: %s", config.Host)

		clientset, err := kubernetes.NewForConfig(config)
		log.Println(clientset)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		log.Println("Testing communication with server")
		_, err = clientset.Discovery().ServerVersion()

		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		log.Println("Communication with server successful")

		return &KubeClient{
			client: clientset,
		}

	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	// log.Println(clientset)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Println("Testing communication with server")
	_, err = clientset.Discovery().ServerVersion()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	log.Println("Communication with server successful")

	return &KubeClient{
		client: clientset,
	}

}
