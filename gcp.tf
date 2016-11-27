// Configure the Google Cloud provider
provider "google" {
  credentials = "${file(var.gcloud_credentials)}"
  project     = "${var.gcloud_project}"
  region      = "${var.gcloud_region}"
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

resource "google_dns_record_set" "hkjn_vpn" {
  name = "vpn.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas      = ["${var.vpn_ip}"]
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

resource "google_dns_record_set" "sample" {
  name = "m1.${google_dns_managed_zone.tf_zone.dns_name}"
  type = "A"
  ttl  = 150

  managed_zone = "${google_dns_managed_zone.tf_zone.name}"
  rrdatas      = ["1.2.3.4"]
}
