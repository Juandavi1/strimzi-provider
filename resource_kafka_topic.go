package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schemaK8s "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)





func resourceKafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"topic_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"partitions": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"replicas": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	// uses the current context in kubeconfig
	// path-to-kubeconfig -- for example, /root/.kube/config
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/juancorrea/.kube/config")
	if err != nil {
		log.Fatal(err)
	}

	err = AddToScheme(scheme.Scheme)
	if err != nil {
		log.Fatal(err)
	}

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schemaK8s.GroupVersion{Group: GroupName, Version: GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	restClient, _ := rest.UnversionedRESTClientFor(&crdConfig)
	result := KafkaTopic{}

	topic := &KafkaTopic{
		TypeMeta:   metav1.TypeMeta{
			Kind: "KafkaTopic",
			APIVersion: GroupName + "/" + GroupVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "product_324",
			Namespace: "kafka",
			Labels: map[string]string{},
		},
		Spec: KafkaTopicSpec{
			TopicName: "JODAAA",
			Partitions: 1,
			Replicas: 2,
		},
	}

	err = restClient.
	 	Post().
		Namespace("kafka").
		Resource("kafkatopics").
		Body(&topic).
	 	Do(context.TODO()).
	 	Into(&result)

	if err != nil {
		log.Fatal(err)
	}

	d.SetId(result.Name)
	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}