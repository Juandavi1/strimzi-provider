terraform {
  required_version = "0.12.24"
  required_providers {}
}

provider "strimzi" {}

resource "strimzi_kafka_topic" "mi_prueba" {
  config_path_k8s = "~/.kube/config"
  namespace = "kafka"
  topic_name = "my_topic"
  partitions = 10
  replicas = 10
}