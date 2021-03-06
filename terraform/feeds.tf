provider "google" {
  project = var.project
  region  = var.region
}

terraform {
  backend "gcs" {
    bucket = "oss-feeds-tf-state"
    prefix = "terraform/state"
  }
}

locals {
  services = [
    "cloudbuild.googleapis.com",
    "run.googleapis.com",
    "cloudscheduler.googleapis.com",
  ]
}

resource "google_service_account" "run-invoker-account" {
  account_id   = "run-invoker-sa"
  display_name = "Feed Run Invoker"
}

resource "google_project_iam_member" "run-invoker-iam" {
  role   = "roles/run.invoker"
  member = "serviceAccount:${google_service_account.run-invoker-account.email}"
}

resource "google_project_service" "services" {
  for_each           = toset(local.services)
  service            = each.value
  disable_on_destroy = false
}

resource "google_pubsub_topic" "feed-topic" {
  name = "feed-topic"
}

resource "google_storage_bucket" "feed-functions-bucket" {
  name          = "${var.project}-feed-functions-bucket"
  force_destroy = true
}
