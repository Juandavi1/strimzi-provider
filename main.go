package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	schemaK8s "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *KafkaTopic) DeepCopyInto(out *KafkaTopic) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = KafkaTopicSpec{
		Replicas: in.Spec.Replicas,
		Partitions: in.Spec.Partitions,
		TopicName: in.Spec.TopicName,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *KafkaTopic) DeepCopyObject() runtime.Object {
	out := KafkaTopic{}
	in.DeepCopyInto(&out)
	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *KafkaTopicList) DeepCopyObject() runtime.Object {
	out := KafkaTopicList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]KafkaTopic, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return &out
}

type KafkaTopicSpec struct {
	TopicName string `json:"topicName"`
	Partitions int `json:"partitions"`
	Replicas int `json:"replicas"`
}

type KafkaTopic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KafkaTopicSpec `json:"spec"`
}

type KafkaTopicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items []KafkaTopic `json:"items"`
}

const GroupName = "kafka.strimzi.io"
const GroupVersion = "v1beta1"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&KafkaTopic{},
		&KafkaTopicList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func main() {
	/*plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})*/
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
			Name: "product324",
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
		Body(topic).
		Do(context.TODO()).
		Into(&result)

	if err != nil {
		log.Fatal(err)
	}

	log.Print(result.Name)
}