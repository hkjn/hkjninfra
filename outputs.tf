output zero-dev0-ip {
  value = "${google_compute_instance.zg0.network_interface.0.access_config.0.assigned_nat_ip}"
}
