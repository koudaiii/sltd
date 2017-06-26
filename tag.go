package main

import (
	"log"

	"github.com/koudaiii/sltd/aws/elb"
	"github.com/koudaiii/sltd/kubernetes"
)

type Client struct {
	kubeclient *kubernetes.KubeClient
	awsclient  *elb.AwsClient
}

func NewClient(inCluster bool) *Client {
	return &Client{
		kubeclient: kubernetes.NewKubeClient(inCluster),
		awsclient:  elb.NewELBClient(),
	}
}

func (c *Client) process() {

	namespaces, err := c.kubeclient.GetAllNamespaces()
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println(namespaces)

	svc, err := c.kubeclient.GetAllServices(namespaces)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(svc)

	for _, s := range svc {
		tags, err := c.awsclient.DescribeTags(s.Name)
		if err != nil {
			log.Println(err)
			return
		}
		c.attachELBTags(tags, c.kubeclient.UpdateLabelsToDataDogFormat(exchangeTypeFromTagsToLabels(tags), s))
	}
}

func (c *Client) attachELBTags(tags []elb.Tag, service kubernetes.Service) error {
	log.Println(tags)
	log.Println(service)
	for _, s := range service.Labels {
		alreadyTag := false

		labelToTag := &elb.Tag{
			Key:   s.Key,
			Value: s.Value,
		}

		for _, t := range tags {
			if t.Key == s.Key && t.Value == s.Value {
				alreadyTag = true
				break
			}
			if t.Key == s.Key {
				log.Println("Replace Tag")
				log.Printf("Before: %v\n",t)
				log.Printf("After: %v\n",s)
				c.awsclient.DeleteTag(service.Name, labelToTag.Key)
				c.awsclient.AddTag(service.Name, labelToTag)
			}
		}
		if alreadyTag {
			log.Println("skip. ELB Already tagged.")
			log.Println(labelToTag)
		} else {
			log.Println("Add Tag")
			log.Println(s)
			c.awsclient.AddTag(service.Name, labelToTag)
		}
	}
	return nil
}

func exchangeTypeFromTagsToLabels(tags []elb.Tag) (labels []kubernetes.Label){
	for _, t := range tags {
		labels = append(labels, kubernetes.Label{
			Key: t.Key,
			Value: t.Value,
		})
	}
	return 	labels
}
