output "hkjn_gz0_ip" {
  value = "${google_dns_record_set.hkjn_gz0.rrdatas[0]}"
}

output "hkjn_zs10_ip" {
  value = "${google_dns_record_set.hkjn_zs10.rrdatas[0]}"
}

output "hkjn_zg1_ip" {
  value = "${google_compute_instance.zg1.network_interface.0.access_config.0.assigned_nat_ip}"
}

output "hkjn_zg3_ip" {
  value = "${google_compute_instance.zg3.network_interface.0.access_config.0.assigned_nat_ip}"
}

output "hkjn_cities_addr" {
  value = "${google_dns_record_set.hkjn_cities.name}"
}

output "hkjn_web_ip" {
  value = "${var.hkjnweb_ip}"
}

output "hkjn_admin_addr" {
  value = "${google_dns_record_set.hkjn_admin.name}"
}

output "hkjn_vpn_addr" {
  value = "${google_dns_record_set.hkjn_vpn.name}"
}

output "hkjn_mon_addr" {
  value = "${google_dns_record_set.hkjn_mon.name}"
}

output "hkjn_iosdev_addr" {
  value = "${google_dns_record_set.hkjn_iosdev.name}"
}
