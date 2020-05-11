package v1beta1

import (
	"github.com/Juandavi1/strimzi-provider/api/types/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type KafkaTopicV1beta1Interface interface {
	KafkaTopic(namespace string) KafkaTopicInterface
}

type KafkaTopicV1beta1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*KafkaTopicV1beta1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1beta1.GroupName, Version: v1beta1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &KafkaTopicV1beta1Client{restClient: client}, nil
}

func (c *KafkaTopicV1beta1Client) KafkaTopics(namespace string) KafkaTopicInterface {
	return &kafkaTopicClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}