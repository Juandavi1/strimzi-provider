package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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