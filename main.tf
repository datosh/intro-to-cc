terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "6.0.1"
    }
  }
}

provider "google" {
  project     = "kammel-dev"
  region      = "europe-west4"
  zone        = "europe-west4-a"
}