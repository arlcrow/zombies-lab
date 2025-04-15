terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

// Реальный инстанс для production
resource "google_compute_instance" "lab_instance" {
  count        = var.debug ? 0 : 1
  name         = "lab-${var.lab_id}"
  machine_type = "e2-micro"
  zone         = "${var.region}-a"

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2004-lts"
    }
  }

  network_interface {
    network = "default"
    access_config {}
  }

  metadata = {
    ssh-keys = "ubuntu:${file(var.ssh_public_key_path)}"
  }

  tags = ["lab-instance"]
}

// Dummy ресурс для debug режима
resource "null_resource" "lab_instance_dummy" {
  count = var.debug ? 1 : 0
  
  triggers = {
    lab_id = var.lab_id
  }

  provisioner "local-exec" {
    command = "echo 'Provisioning dummy lab ${var.lab_id}'"
  }
}

output "instance_ip" {
  value = var.debug ? "192.0.2.1" : try(google_compute_instance.lab_instance[0].network_interface[0].access_config[0].nat_ip, "")
}
