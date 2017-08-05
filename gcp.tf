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

resource "google_compute_firewall" "default" {
  name    = "tf-allow-ssh-ping"
  network = "${google_compute_network.default.name}"
  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "tcp"
    ports    = ["22", "60000-60100"]
  }
  target_tags = ["dev"]
}

resource "google_compute_disk" "zg0_disk0" {
  name  = "test-disk"
  type  = "pd-ssd"
  zone  = "europe-west3-b"
}

resource "google_compute_instance" "zg0" {
  name         = "zero-dev0"
  description  = "Dev and build instance"
  machine_type = "g1-small"
  zone         = "europe-west3-b"
  tags = ["dev", "builder"]
  disk {
    image = "coreos-alpha-1492-1-0-v20170803"
  }
  # Local SSD disk
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
  metadata_startup_script = "echo hi > /tmp/test.txt"
  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_dns_managed_zone" "hkjn_zone" {
  name     = "hkjn-zone"
  dns_name = "hkjn.me."
}

resource "google_dns_record_set" "hkjn_web" {
  name = "${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 150

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas      = ["${var.hkjnweb_ip}"]
}

resource "google_dns_record_set" "hkjn_web_www" {
  name = "www.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 150

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas      = ["${var.hkjnweb_ip}"]
}

resource "google_dns_record_set" "hkjn_mail" {
  name = "${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "MX"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"

  rrdatas = [
    "1 aspmx.l.google.com.",
    "5 alt1.aspmx.l.google.com.",
    "5 alt2.aspmx.l.google.com.",
    "10 aspmx2.googlemail.com.",
    "10 aspmx3.googlemail.com.",
  ]
}

# NS and SOA records do need to be set on the top level, but seem to be auto-created
# by GCP when the managed zone is created, so if we try to create them here as well
# we'll get an HTTP 409 "already exists" error from the API.
#
#resource "google_dns_record_set" "hkjn_ns" {
#  name = "${google_dns_managed_zone.hkjn_zone.dns_name}"
#  type = "NS"
# ttl  = 21600
#
#  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
#
#	rrdatas = [
#		"ns-cloud-b1.googledomains.com.",
#		"ns-cloud-b2.googledomains.com.",
#		"ns-cloud-b3.googledomains.com.",
#		"ns-cloud-b4.googledomains.com.",
#	]
#}
#

#resource "google_dns_record_set" "hkjn_soa" {
#  name = "${google_dns_managed_zone.hkjn_zone.dns_name}"
#  type = "SOA"
#  ttl  = 21600
#
#  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
#
#	rrdatas = [
#		"ns-cloud-b1.googledomains.com. dns-admin.google.com. 0 21600 3600 1209600 300",
#	]
#}

resource "google_dns_record_set" "hkjn_mon" {
  name = "mon.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"

  rrdatas = [
    "${var.mon_ip}",
  ]
}

resource "google_dns_record_set" "hkjn_iosdev" {
  name = "iosdev.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"

  rrdatas = [
    "${var.iosdev_ip}",
  ]
}

resource "google_dns_record_set" "hkjn_gz0" {
  name = "gz0.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"

  rrdatas = [
    "130.211.84.102",
  ]
}

resource "google_dns_record_set" "hkjn_gz5" {
  name = "gz5.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas = [
    "35.189.120.224",
  ]
}

resource "google_dns_record_set" "hkjn_zs10" {
  name = "zs10.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas = [
    "163.172.184.153",
  ]
}

resource "google_dns_record_set" "hkjn_vpn" {
  name = "vpn.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas      = ["${var.vpn_ip}"]
}

resource "google_dns_record_set" "hkjn_cities" {
  name = "cities.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas      = ["${var.cities_ip}"]
}

resource "google_dns_record_set" "hkjn_tf_ns" {
  name = "tf.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "NS"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"

  rrdatas = [
    "ns-cloud-a1.googledomains.com.",
    "ns-cloud-a2.googledomains.com.",
    "ns-cloud-a3.googledomains.com.",
    "ns-cloud-a4.googledomains.com.",
  ]
}

resource "google_dns_managed_zone" "tf_zone" {
  name     = "tf-zone"
  dns_name = "tf.hkjn.me."
}

resource "google_dns_record_set" "dev" {
  name = "dev.${google_dns_managed_zone.tf_zone.dns_name}"
  type = "A"
  ttl  = 150

  managed_zone = "${google_dns_managed_zone.tf_zone.name}"
  rrdatas      = ["212.47.239.127"]                        # sz9
}

resource "google_dns_managed_zone" "elentari_world_zone" {
  name     = "elentari-world-zone"
  dns_name = "elentari.world."
}

resource "google_dns_record_set" "elentari_world_web" {
  name = "${google_dns_managed_zone.elentari_world_zone.dns_name}"
  type = "A"
  ttl  = 150

  managed_zone = "${google_dns_managed_zone.elentari_world_zone.name}"
  rrdatas      = ["${var.elentari_world_ip}"]
}

resource "google_dns_record_set" "elentari_world_cname" {
  name = "www.${google_dns_managed_zone.elentari_world_zone.dns_name}"
  type = "CNAME"
  ttl  = 150

  managed_zone = "${google_dns_managed_zone.elentari_world_zone.name}"
  rrdatas      = ["sultanyoga.com."]
}
