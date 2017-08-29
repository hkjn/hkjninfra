variable "gcloud_credentials" {
  default = ".gcp/tf-dns-editor.json"
}

variable "aruna_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDOz60CdAF7cYXTlKBh+nKhyFa8PhSlSsMODOd4MtSM3Vepx1i1f716BysgfCP822vX14dPJWTEhin/DzrTKsfVDz3EKyj8i4H/eAQ5gs9l99R2DOG6qhZ0SKFMQM1gNTYYvX3qWg+GG/55xrgm6Ol3o0fzSpi0qxB+tvB2QbX7gztaaouerKAWYpe55Oe3mGHG+AYC1FJ4WZo49hQJsUSTAa0rfI2AYJXY5iOKjzjjvJ4zgw8bZax8hUXQG9CsxZRh6P15DOzrDAtBYMwYBIYJ1JolQ25h7e/9QLp18D5FHBN9j0hd7zVTs8CIxz33sBupZNpZ+k7jWOT0BFOkE+2aO4UjnZw8U1LQZpwQk7/dGA/AnlwjJ8W0GwgpTfovuLJvDoOanXoAP6ccsl8upU0R9yGY/JLRultIHAY3YBi9hb7oNNMGtVfyVntokAdiTBPLDDlvIzDU3He8MiOQwRIE6w23A7Z+U4I2dc6Pm584Ef94jUY2H8ZwDy0wYlVX0TPnyamdvqoRIb428SYWPCaSTdLx5WwX15EuR1ZFcmNpxZCM8K246KnkPRAN+qddanqm5DM+xuRjIvH45p1MeDwKGi7YhVtZMduz0KdQMwRPeqA1qnV7h6ERCEBsgDihySC9862dU6b8NNXv5+9C8qw0sINro4Q2xhV2px+bg8xsow== Arunaskilt"
}

variable "gcloud_project" {
  default = "henrik-jonsson"
}

variable "gcloud_region" {
  default = "europe-west1"
}

variable "coreos_alpha_image" {
  default = "coreos-alpha-1492-3-0-v20170810"
}

variable "coreos_beta_image" {
  default = "coreos-beta-1465-4-0-v20170810"
}

variable "gz1_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDhkGPTKaQxOkrRO4gyJFr5ewcIHVcQ+YyEeTWfK1GpwKwjMZNhbXCL0DeTxjYx+JHMqStxlKi6H/+dDJSiohbIi1MqPcpPvaY2whLPSHrV4b17tYXjnkZbowDp54QCTLgcXmfckpKosulV56t6GGXZ2PmlP85QkpyeF+8/EyS+8XFQi1GeWOwgnBvmyRRxSBwsJnMt+aJdCcIRv0sO+pKuzCJFLJ2D82sRbfLLIWOy0BkzaUMevO+sUYmmRgJ2Kis7MFNYC8c6i2GHTuRex5t0fHyCzp7qHLISAhjS4d5mmpfUeMJot7WlOiwlVTMcoQSFY0boTuDSEoVxzwu3/R2zc6Hu+RGBKZbIR6ksKZFPe64dPGJXzfX/3oyqhFQKF3A9kl9wNGq2RQweEaRtHmvToAf49iN+8q5vMdSD0XbCxJwakCj+nKLprp7yexydj6JX+lDbXtxJyF5ECuN4mE0Mi09dFJSS4e09c9bCSrIM+6AXOLpU10Kq2KqAiRxRpdV5UxheHeK89E6Z4Y8d1DggmrnQxSfvmfCctqsJVWuXC2yl+py8PmP1I9HclyBwCReTVUpOOJknQILh1lubAtboBbBK9zAIh8k13224ricxSDhiP4T+Jl6tguzxK7hQ1zQJlRipyTtaktSSikA2aMobmHo2IoifRkblJ2Wc6BiOOQ== gz1"
}

variable "hkjnweb_ip" {
  default = "163.172.173.208"
}

variable "mon_ip" {
  default = "163.172.184.153"
}

variable "vpn_ip" {
  default = "163.172.184.153"
}

variable "cities_ip" {
  default = "163.172.184.153"
}

variable "elentari_world_enabled" {
  default = false
}

variable "version" {}
