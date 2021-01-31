resource "google_cloud_run_service" "run-rubygems" {
  name     = "rubygems-run-srv"
  location = var.region

  template {
    spec {
      containers {
        image = "gcr.io/${var.project}/feeds-rubygems"
        env {
          name  = "OSSMALWARE_TOPIC_URL"
          value = "gcppubsub://${google_pubsub_topic.feed-topic.id}"
        }
      }
    }
  }
}
