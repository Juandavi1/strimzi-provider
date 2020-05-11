package v1beta1

import (
	"context"
	"github.com/Juandavi1/strimzi-provider/api/types/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)


const Resource = "kafkatopics"

type KafkaTopicInterface interface {
	List(opts metav1.ListOptions) (*v1beta1.KafkaTopicList, error)
	Get(name string, options metav1.GetOptions) (*v1beta1.KafkaTopic, error)
	Delete(name string, options metav1.DeleteOptions) error
	Create(topic *v1beta1.KafkaTopic, setting metav1.CreateOptions) (*v1beta1.KafkaTopic, error)
	Update(topic *v1beta1.KafkaTopic, setting metav1.UpdateOptions) (*v1beta1.KafkaTopic, error)
}

type kafkaTopicClient struct {
	restClient rest.Interface
	ns         string
}

func (c *kafkaTopicClient) Get(name string, opts metav1.GetOptions) (*v1beta1.KafkaTopic, error) {
	result := v1beta1.KafkaTopic{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(Resource).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *kafkaTopicClient) List(opts metav1.ListOptions) (*v1beta1.KafkaTopicList, error) {
	result := v1beta1.KafkaTopicList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(Resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *kafkaTopicClient) Create(project *v1beta1.KafkaTopic, opts metav1.CreateOptions) (*v1beta1.KafkaTopic, error) {
	result := v1beta1.KafkaTopic{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(Resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(project).
		Do(context.TODO()).
		Into(&result)
	return &result, err
}

func (c *kafkaTopicClient) Update(kafkaTopic *v1beta1.KafkaTopic, opts metav1.UpdateOptions) (*v1beta1.KafkaTopic, error) {
	result := v1beta1.KafkaTopic{}
	err := c.restClient.
		Put().
		Namespace(c.ns).
		Resource(Resource).
		Name(kafkaTopic.Name).
		Body(kafkaTopic).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, err
}

func (c *kafkaTopicClient) Delete(name string, opts metav1.DeleteOptions) error {
	c.restClient.
		Delete().
		Namespace(c.ns).
		Resource(Resource).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO())

	return nil
}