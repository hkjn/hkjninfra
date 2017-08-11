# Configure the Google Cloud provider
provider "google" {
  credentials = "${file(var.gcloud_credentials)}"
  project     = "${var.gcloud_project}"
  region      = "${var.gcloud_region}"
}

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

resource "google_compute_disk" "zg0_disk0" {
  name  = "test-disk"
  type  = "pd-ssd"
  zone  = "europe-west3-b"
  size  = "200"
}

resource "google_compute_instance" "zg3" {
  count = 0

  name         = "zg3"
  description  = "Bootstrap test"
  machine_type = "f1-micro"
  zone         = "europe-west3-b"
  tags = ["dev"]
  disk {
    image = "coreos-alpha-1492-1-0-v20170803"
  }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "core:${file(".keys/gz3_id_rsa.pub")}"
  }
  metadata_startup_script = "${file("bootstrap.sh")}"
}


resource "google_compute_instance" "zg1" {
  name         = "zg1"
  description  = "Dev and build instance"
  machine_type = "g1-small"
  zone         = "europe-west3-b"
  tags = ["dev", "builder", "bitcoin"]
  disk {
    image = "coreos-alpha-1492-1-0-v20170803"
  }
  attached_disk {
    source = "${google_compute_disk.zg0_disk0.self_link}"
  }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "core:${var.gz1_pubkey}"
  }
  # TODO: Start discovering metadata from inside instance in gather_facts, e.g:
  # core@zg0 ~ $ curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/description; echo
  # Dev and build instance
  #
  # TODO: Set sshd port in bootstrap:
  # cat /etc/systemd/system/sshd.socket.d/10-sshd-listen-ports.conf
  # [Socket]
  # ListenStream=
  # ListenStream=6200

  metadata_startup_script = "${file("bootstrap.yml")}"
  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}
