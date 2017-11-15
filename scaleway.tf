provider "scaleway" {
  organization = "${chomp(file(var.scaleway_organization_file))}"
  token        = "${chomp(file(var.scaleway_key_file))}"
  region = "${var.scaleway_region}"
}

resource "scaleway_ip" "ip" {
  count = "${var.scalewaytest_enabled ? 1 : 0}"
}

resource "scaleway_security_group" "allow_ssh" {
  name        = "allow_ssh"
  description = "Allow HTTP/S and SSH traffic"
}

resource "scaleway_security_group_rule" "ssh_accept" {
  security_group = "${scaleway_security_group.allow_ssh.id}"

  action    = "accept"
  direction = "inbound"
  ip_range  = "0.0.0.0/0"
  protocol  = "TCP"
  port      = 22
}

resource "scaleway_server" "sc1" {
  count = "${var.scalewaytest_enabled ? 1 : 0}"
  name           = "sc1"
  image          = "${var.scaleway_image}"
  type           = "C1"
  #bootscript     = "#!/bin/bash; echo hi > /tmp/testing"
  security_group = "${scaleway_security_group.allow_ssh.id}"
  public_ip      = "${element(scaleway_ip.ip.*.ip, count.index)}"
  #volume {
  #  size_in_gb = 20
  #  type       = "l_ssd"
  #}
}
