output zg0_ip {
  value = "${google_compute_instance.zg0.network_interface.0.access_config.0.assigned_nat_ip}"
}

output admin_addr {
  value = "${google_dns_record_set.hkjn_admin.name}"
}

output iosdev_addr {
  value = "${google_dns_record_set.hkjn_iosdev.name}"
}
