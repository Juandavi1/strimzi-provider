terraform {
  required_version = "0.12.24"
  required_providers {}
}

provider "strimzi" {}

resource "strimzi_kafka_topic" "mi_prueba" {
  topic_name = "product_324"
  partitions = 1
  replicas = 1
}