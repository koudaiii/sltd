package elb

import (
	"log"

	awspkg "github.com/koudaiii/sltd/aws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
)

type AwsClient struct {
	client *elb.ELB
}

type Tag struct {
	Key   string
	Value string
}

func NewELBClient() *AwsClient {
	return &AwsClient{
		client: elb.New(awspkg.Session()),
	}
}

func (c *AwsClient) DescribeTags(name string) ([]Tag, error) {
	input := &elb.DescribeTagsInput{
		LoadBalancerNames: []*string{
			aws.String(name),
		},
	}

	result, err := c.client.DescribeTags(input)
	if err != nil {
		return nil, err
	}
	tag := []Tag{}

	for _, t := range result.TagDescriptions[0].Tags {
		tag = append(tag, Tag{
			Key:   *t.Key,
			Value: *t.Value,
		})
	}

	return tag, nil
}

func (c *AwsClient) AddTag(name string, tag *Tag) error {
	input := &elb.AddTagsInput{
		LoadBalancerNames: []*string{
			aws.String(name),
		},
		Tags: []*elb.Tag{
			{
				Key:   aws.String(tag.Key),
				Value: aws.String(tag.Value),
			},
		},
	}

	result, err := c.client.AddTags(input)
	if err != nil {
		return err
	}

	log.Println(result)
	return nil
}

func (c *AwsClient) DeleteTag(name string, key string) error {
	input := &elb.RemoveTagsInput{
		LoadBalancerNames: []*string{
			aws.String(name),
		},
		Tags: []*elb.TagKeyOnly{
			{
				Key: aws.String(key),
			},
		},
	}

	result, err := c.client.RemoveTags(input)
	if err != nil {
		return err
	}

	log.Println(result)
	return nil
}