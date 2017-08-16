
output "decenter_world_addr" {
  value = "${google_dns_record_set.decenter_world.name}"
}


#
# *.hkjn.me
#

output "hkjn_addr_admin" {
  value = "${google_dns_record_set.hkjn_admin.name}"
}

output "hkjn_addr_cities" {
  value = "${google_dns_record_set.hkjn_cities.name}"
}

output "hkjn_addr_iosdev" {
  value = "${google_dns_record_set.hkjn_iosdev.name}"
}

output "hkjn_addr_mon" {
  value = "${google_dns_record_set.hkjn_mon.name}"
}

output "hkjn_addr_vpn" {
  value = "${google_dns_record_set.hkjn_vpn.name}"
}

#
# IP addresses of *.hkjn.me
#

output "hkjn_ip_gz0" {
  value = "${google_dns_record_set.hkjn_gz0.rrdatas[0]}"
}

output "hkjn_ip_zs10" {
  value = "${google_dns_record_set.hkjn_zs10.rrdatas[0]}"
}

output "hkjn_ip_zg1" {
  value = "${google_compute_instance.zg1.network_interface.0.access_config.0.assigned_nat_ip}"
}

output "hkjn_ip_zg3" {
  value = "${google_compute_instance.zg3.network_interface.0.access_config.0.assigned_nat_ip}"
}

output "hkjn_ip_web" {
  value = "${var.hkjnweb_ip}"
}
