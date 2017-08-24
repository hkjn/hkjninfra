variable "gcloud_credentials" {
  default = ".gcp/tf-dns-editor.json"
}

variable "elentari_world_enabled" {
  default = false
}

variable "aruna_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9AIuUhqsepqfzgNj+hz+GoGIv5BUqMThpAmTaj4qXnI7ahq7916xvE4OB61CyN0fu69N7cqHN0n0DhzrIKjRJ+mBztKZ1uNDiicoxc2PGpTff42kMA63J36YhS0wPu8hADqm7J9NSt5Ach+IKUtEQNuwnq+nZ3MXOzlRy+B/RsCzfkEatejyEyMPNZ1+TIzzh4xMnjf9ktcR0joxCGASxJPCzI0bV/4NZDWZ8gelTeNobdhfi1qB928VHHG2SrQ+VMNUsnG1D44+fnE0/SnGd8OXA1vFyzn9VLNQE5mhhP7dCaO3O/ODMioEnT2Y5Zn7LfQl4Sm8x+6f2bHh1AzuOL7ActYuur4QVzduPZn1nO6DsHYIKy3OViIQUIk1foN+OeD08QAiL7dgHwzuGYsLs1sK9bD22knb59qlAgfumLjxcxs9dNdVwD5UYJW5/FAQieAql2a5CBhGPIlmk863iVG/gyGgGT7R4dKWGyTraRh2AFEMxVr4Ua/umR1Zx1AN5JZoVv9U4RNQfqO+1oYOKCblx8OMPzGkjYkG2ryPM9bjJVEQg9XfJpmBwI0iyUSzgLDvv38URvpxqX3SeYVMS2iktUde5y+dNTS9QRrlBu7ozMWQN7oqp2Bc+9n1X220yTbvUBPWTs+9X0eYpshfQgcwaI67uGfpytueoFSSSFw=="
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

variable "elentari_world_ip" {
  default = "163.172.184.153"
}
