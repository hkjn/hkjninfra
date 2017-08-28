# Configure the Google Cloud provider
provider "google" {
  credentials = "${file(var.gcloud_credentials)}"
  project     = "${var.gcloud_project}"
  region      = "${var.gcloud_region}"
}

#
# Network and firewalls.
#

resource "google_compute_network" "default" {
  name                    = "tf-net0"
  auto_create_subnetworks = "true"
}

resource "google_compute_firewall" "allow_ssh" {
  name    = "tf-allow-ssh-ping"
  network = "${google_compute_network.default.name}"
  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "tcp"
    ports    = ["22", "6200"]
  }
  allow {
    protocol = "udp"
    ports    = ["60000-60100"]
  }
  source_ranges = ["0.0.0.0/0"]
  target_tags = ["dev"]
}

resource "google_compute_firewall" "http" {
  name = "tf-allow-http"
  network = "${google_compute_network.default.name}"
  allow {
    protocol = "tcp"
    ports = ["80"]
  }
  allow {
    protocol = "tcp"
    ports = ["443"]
  }
  source_ranges = ["0.0.0.0/0"]
  target_tags = ["http"]
}

resource "google_compute_firewall" "bitcoin" {
  name = "tf-allow-bitcoin"
  network = "${google_compute_network.default.name}"
  allow {
    protocol = "tcp"
    ports = ["8333"]
  }
  source_ranges = ["0.0.0.0/0"]
  target_tags = ["bitcoin"]
}

#
# Storage
#

#
# Bitcoind storage disk
#
resource "google_compute_disk" "zg0_disk0" {
  name  = "test-disk"
  type  = "pd-ssd"
  zone  = "europe-west3-b"
  size  = "200"
}

#
# decenter.world persistent disk
#
resource "google_compute_disk" "zdisk1" {
  name  = "zdisk1"
  type  = "pd-ssd"
  zone  = "europe-west3-b"
  size  = "20"
}

#
# GCP instances
#

resource "google_compute_instance" "zg1" {
  name         = "zg1"
  description  = "Dev and bitcoind"
  machine_type = "g1-small"
  zone         = "europe-west3-b"
  tags = ["dev", "builder", "bitcoin", "http"]
  disk {
    image = "${var.coreos_beta_image}"
  }
  attached_disk {
    source = "${google_compute_disk.zg0_disk0.self_link}"
  }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "core:${file(".keys/gz3_id_rsa.pub")}"
    user-data = "${file("bootstrap_zg1.json")}"
  }
  # TODO: Set sshd port in bootstrap:
  # cat /etc/systemd/system/sshd.socket.d/10-sshd-listen-ports.conf
  # [Socket]
  # ListenStream=
  # ListenStream=6200

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_instance" "zg3" {
  count = 1
  name         = "zg3"
  description  = "Bootstrap test"
  machine_type = "f1-micro"
  zone         = "europe-west3-b"
  tags = ["dev", "http"]
  disk {
    image = "${var.coreos_alpha_image}"
  }
  attached_disk {
    source = "${google_compute_disk.zdisk1.self_link}"
  }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "core:${file(".keys/gz3_id_rsa.pub")}"
    buildCommit = "${var.git_commit}"
    user-data = "${file("bootstrap_zg3.json")}"
  }
}

resource "google_compute_instance" "elentari-world" {
  count = "${var.elentari_world_enabled ? 1 : 0}"
  name         = "elentari-world"
  description  = "The server for elentari.world."
  machine_type = "f1-micro"
  zone         = "europe-west3-b"
  tags = ["dev", "http"]
  disk { image = "${var.coreos_alpha_image}" }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "core:${var.aruna_pubkey}"
  }
}
