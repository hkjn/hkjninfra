provider "scaleway" {
  organization = "${chomp(file(var.scaleway_organization_file))}"
  token        = "${chomp(file(var.scaleway_key_file))}"
  region = "${var.scaleway_region}"
}

resource "scaleway_ip" "hkjnprod" {
  count = "${var.hkjnprod_enabled ? 1 : 0}"
}

resource "scaleway_server" "hkjnprod" {
  count = "${var.hkjnprod_enabled ? 1 : 0}"
  name           = "prod.hkjn.me"
  image          = "${var.scaleway_image}"
  type           = "C1"
  public_ip      = "${element(scaleway_ip.hkjnprod.*.ip, count.index)}"
}

resource "scaleway_volume" "prod1" {
  name       = "proddisk1"
  count = "${var.hkjnprod_enabled ? 1 : 0}"
  size_in_gb = 150
  type       = "l_ssd"
}

resource "scaleway_volume_attachment" "attachment1" {
  count = "${var.hkjnprod_enabled ? 1 : 0}"
  server = "${scaleway_server.hkjnprod.id}"
  volume = "${scaleway_volume.prod1.id}"
}

resource "scaleway_volume" "prod2" {
  name       = "proddisk2"
  count = "${var.hkjnprod_enabled ? 1 : 0}"
  size_in_gb = 150
  type       = "l_ssd"
}

resource "scaleway_volume_attachment" "attachment2" {
  count = "${var.hkjnprod_enabled ? 1 : 0}"
  server = "${scaleway_server.hkjnprod.id}"
  volume = "${scaleway_volume.prod2.id}"
}
