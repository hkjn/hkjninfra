#
# *.hkjn.me
#
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

resource "google_dns_record_set" "hkjn_admin" {
  name = "admin.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300
  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas = [
    "130.211.84.102",
  ]
}

resource "google_dns_record_set" "hkjn_cities" {
  name = "cities.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300
  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas      = ["${var.cities_ip}"]
}

resource "google_dns_record_set" "hkjn_iosdev" {
  name = "iosdev.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300
  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas = [
    "${google_compute_instance.zg1.network_interface.0.access_config.0.assigned_nat_ip}",
  ]
}

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

#
# *.tf.hkjn.me
#
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


#
# *.decenter.world
#

resource "google_dns_managed_zone" "decenter_world_zone" {
  name     = "decenter-world-zone"
  dns_name = "decenter.world."
}

resource "google_dns_record_set" "decenter_world" {
  name = "${google_dns_managed_zone.decenter_world_zone.dns_name}"
  type = "A"
  ttl  = 300
  managed_zone = "${google_dns_managed_zone.decenter_world_zone.name}"
  rrdatas = [
    "${google_compute_instance.zg3.network_interface.0.access_config.0.assigned_nat_ip}",
  ]
}

#
# *.elentari.world
#

resource "google_dns_managed_zone" "elentari_world_zone" {
  name     = "elentari-world-zone"
  dns_name = "elentari.world."
}

resource "google_dns_record_set" "elentari_world_web" {
  count = "${var.elentari_world_enabled ? 1 : 0}"
  name = "${google_dns_managed_zone.elentari_world_zone.dns_name}"
  type = "A"
  ttl  = 300
  managed_zone = "${google_dns_managed_zone.elentari_world_zone.name}"
  rrdatas      = [
    "${google_compute_instance.elentari-world.network_interface.0.access_config.0.assigned_nat_ip}",
  ]
}

resource "google_dns_record_set" "elentari_world_cname" {
  name = "www.${google_dns_managed_zone.elentari_world_zone.dns_name}"
  type = "CNAME"
  ttl  = 150
  managed_zone = "${google_dns_managed_zone.elentari_world_zone.name}"
  rrdatas      = ["sultanyoga.com."]
}

#
# *.unmiddle.men
#

resource "google_dns_managed_zone" "unmiddle_men_zone" {
  name     = "unmiddle-men-zone"
  dns_name = "ummiddle.men."
}

#
# Individual *.hkjn.me nodes
#

resource "google_dns_record_set" "hkjn_gz0" {
  name = "gz0.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300
  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas = [
    "130.211.84.102",
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

resource "google_dns_record_set" "hkjn_zg1" {
  name = "zg1.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300
  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas = [
    "${google_compute_instance.zg1.network_interface.0.access_config.0.assigned_nat_ip}",
  ]
}


resource "google_dns_record_set" "hkjn_zg3" {
  name = "zg3.${google_dns_managed_zone.hkjn_zone.dns_name}"
  type = "A"
  ttl  = 300
  managed_zone = "${google_dns_managed_zone.hkjn_zone.name}"
  rrdatas = [
    "${google_compute_instance.zg3.network_interface.0.access_config.0.assigned_nat_ip}",
  ]
}

