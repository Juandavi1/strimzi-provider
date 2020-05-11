package main

import (
	"fmt"
	"github.com/Juandavi1/strimzi-provider/api/types/v1beta1"
	clientV1beta1 "github.com/Juandavi1/strimzi-provider/clientset/v1beta1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

const resource = "KafkaTopic"

func buildClient(configPath string) (client clientV1beta1.KafkaTopicInterface)  {
	var config *rest.Config
	var err error

	config, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = v1beta1.AddToScheme(scheme.Scheme)
	if err != nil {
		log.Fatal(err)
	}

	clientSet,err := clientV1beta1.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	client = clientSet.KafkaTopics("kafka")
	return
}

func resourceKafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_path_k8s": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"partitions": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"replicas": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"topic_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

	replicas := d.Get("replicas").(int)
	partitions := d.Get("partitions").(int)
	topicName := d.Get("topic_name").(string)
	configPath := d.Get("config_path_k8s").(string)
	namespace := d.Get("namespace").(string)

	client := buildClient(configPath)

	topic := &v1beta1.KafkaTopic{
		TypeMeta:   metav1.TypeMeta{
			Kind: resource,
			APIVersion: v1beta1.GroupName + "/" + v1beta1.GroupVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: topicName,
			Namespace: namespace,
			Labels: map[string]string{},
		},
		Spec: v1beta1.KafkaTopicSpec{
			TopicName: topicName,
			Partitions: partitions,
			Replicas: replicas,
		},
	}

	topicCreated,err := client.Create(topic, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("Failed to create KafkaTopic: %s", err)
	}

	d.SetId(topicCreated.ObjectMeta.Name)
	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {

	configPath := d.Get("config_path_k8s").(string)
	client := buildClient(configPath)
	_,err := client.Get(d.Get("topic_name").(string), metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Failed to read KafkaTopic: %s", err)
	}
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {

	replicas := d.Get("replicas").(int)
	partitions := d.Get("partitions").(int)
	topicName := d.Get("topic_name").(string)
	configPath := d.Get("config_path_k8s").(string)
	namespace := d.Get("namespace").(string)

	client := buildClient(configPath)
	kafkaTopic,err := client.Get(d.Get("topic_name").(string), metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Failed to read KafkaTopic: %s", err)
	}

	if d.HasChange("partitions") {
		kafkaTopic.Spec.Partitions = partitions
	}

	if d.HasChange("replicas") {
		kafkaTopic.Spec.Replicas = replicas
	}

	if d.HasChange("topic_name") {
		kafkaTopic.Spec.TopicName = topicName
	}

	if d.HasChange("namespace") {
		kafkaTopic.Namespace = namespace
	}

	updated,err := client.Update(kafkaTopic, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Failed to update KafkaTopic: %s", err)
	}

	d.SetId(updated.ObjectMeta.Name)
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	configPath := d.Get("config_path_k8s").(string)
	client := buildClient(configPath)
	err := client.Delete(d.Get("topic_name").(string), metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Failed to delete KafkaTopic: %s", err)
	}
	d.SetId("")
	return nil
}