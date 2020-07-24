## Strimzi Provider for Terraform 

  provider "strimzi" {}

  resource "strimzi_kafka_topic" "mi_prueba" {
    config_path_k8s = "~/.kube/config"
    namespace = "kafka"
    topic_name = "my_topic"
    partitions = 10
    replicas = 10
  }
