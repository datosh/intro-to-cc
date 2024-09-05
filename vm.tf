resource "google_service_account" "default" {
  account_id   = "my-custom-sa"
  display_name = "Custom SA for VM Instance"
}

resource "google_compute_network" "vpc_network" {
  name                    = "default"
  auto_create_subnetworks = true
  mtu                     = 1460
}

resource "google_compute_firewall" "default" {
  name    = "allow-ssh"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = ["22", "80"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_instance" "default" {
  name         = "cvm"
  # Must be N2D machine to support AMD SEV-SNP
  # https://cloud.google.com/confidential-computing/confidential-vm/docs/create-a-confidential-vm-instance#create-instance
  machine_type = "n2d-standard-2"

  # https://cloud.google.com/compute/docs/instances/specify-min-cpu-platform#availablezones
  # or dynamically discover via gcloud
  # `gcloud compute zones describe europe-west3-a`
  min_cpu_platform = "AMD Milan"

  confidential_instance_config {
    confidential_instance_type = "SEV_SNP"
  }

  scheduling {
    on_host_maintenance = "TERMINATE"
  }

  tags = ["cvm"]

  boot_disk {
    initialize_params {
      image = "ubuntu-minimal-2404-noble-amd64-v20240829"
      labels = {
        my_label = "value"
      }
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "NVME"
  }

  network_interface {
    network = google_compute_network.vpc_network.name

    access_config {
      // Ephemeral public IP
    }
  }

  metadata = {
    ssh-keys = join("\n", [for key in var.ssh_keys : "${key.user}:${key.key}"])
  }

  metadata_startup_script = file("startup.sh")

  service_account {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    email  = google_service_account.default.email
    scopes = ["cloud-platform"]
  }
}
