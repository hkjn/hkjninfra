variable "gcloud_credentials" {
  default = ".gcp/tf-dns-editor.json"
}

variable "hkjnprod_enabled" {
  default = false
}

variable "scaleway_region" {
  default = "par1"
}

variable "scaleway_organization_file" {
  default = ".scaleway/scaleway0_organization"
}

variable "scaleway_key_file" {
  default = ".scaleway/scaleway0_key"
}

variable "scaleway_image" {
  # Fetched with:
  # curl -H "X-Auth-Token: $(cat .scaleway/scaleway0_key)" -H 'Content-Type: application/json' 'https://cp-par1.scaleway.com/images/?page=1&per_page=100' > scaleway_images.json
  default = "ee0d3a38-1e8a-4407-bc02-d35dd588efa2"
}

variable "digitalocean_token_file" {
  default = ".digitalocean/digitalocean0_token"
}

variable "digitalocean_image" {
  # Images can be fetched with:
  # curl -X GET --silent "https://api.digitalocean.com/v2/images?per_page=999" -H "Authorization: Bearer $(cat .digitalocean/digitalocean0_token)" |jq '.'
  default = "coreos-stable"
}

variable "digitalocean1_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCtb2pbWvs+LZapn7yTtKD0NrTBmgTXoDgERZAcJ47ziGawmaBtQ6UgZLppJ10atm/SddQohYwRNnZp1XukbQf83g8XgJ+W91WIHZdggmjnVIGIan+TbNsNVPAtpRXlEIiYSgpy8jHci1q0u6jat2/FGi0x3AUqGUAaUJxwvUEDnkf/g25lq/hyOZi0yoCFjytQC3/TlgFaCW4T/8RRPY4aHnEtY/D0GE4UPBsDBK+wT2/cFrcxAVSNLHW5i44ChzXbCtwTyPv9+FayngjQyPtze84KPZa7gv7XITcAmOnfRRNCN5UlNCJEwMXg2jZtext4OuUPcZo1z5/D6GaP3n+3WALIk87h8+a4aJ0hJ3RwSNEZaKicgTsziNfmuXaztC/DZAG/fpbdG0O+VULwHwTIFJVsXB9yP2bXqmJOTj9/T0NimN8XPzBa+ixo7BPebuMCwyIS31zhdzxoidi3tt/bHCozY0Aoh/sGelku0xIsT8Io3WrMX5Cqgbz3EPnjeMYP9kbG695xwRfRdS0t+qjxezG73wxs4FPwPY4e9T19G8v91XzgbKK9c33S4DyoL4Zf4nGW+i0dqkBsLgegMLXAXsNiwA1hcHDLXm7hO+viwC50RSwHUeCvwNPDzX18QSmw+mt+OOoMClnWdM6IuOzO3q4TwkwMkpqY8ncrmg+8XQ=="
}

variable "aruna_pubkey_arunallave" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC0qu+LwXD4QG1kQv4F8wjY1bIDR8J72kX/bt8kP/qdYJ0Y0dzSjwGUG79BVBx/StkDRYr+eiDk/fa/mL5ZUCtsraRiOyPoCqvfkIBvHRAjiKynwa/gTCHw1gC/fT2HgTq996VNoqnZDLV4Icc9t7BPMecaNVquBc3+eW6N8fCSV1Oj8/+NlJyV2Efom2ciG/GURknYO5DdhcH+NqmvjhMlBy/XLkfCSGiwNXm4XRHDAe5ZKELTDHH4sFXCPUjSbv7tB75iHDY71lBU+yNb7i3uOO/ODluCxeXranWbVyE76IpA8XoAPdsq0xcfZKLt3ZEGd1NU+v8ZWMXrTUJxX147UGoXsh76pPnY3pgrDxY75K8zgetS2y/TsNcnZP4KvWZK/6hVGasyhmoQxYgr4yj+4L8J9izzGSYYmm+vXL8Q+1ejgWG34KzRJjK8kQvdAKiSu+nnnPSoizsv6WEowAkBQp4pfqjHf5LIFHrfvVoU7HFLKp0wWSYeF/3chl65rhu8y67TkWHBq3mVqCEKqnciXgbYu/qMlfWpv3uHqi1Qhodpm3PPM3O1rk/+WKKrfi1gwcr6F+tz6aq6y1qnTx4YeW0Bqd1uf7J5vEe0HKjjL3Dt9BsW9f5tS3FAic+AtUlPO2vADYYonYc9FiLlux/UCDS7oIBUjsfczsM8ML0M9w== aruna@ipad"
}

variable "zg0_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCw76RMWPriShiqQKy06oJjlYURKdwStODyUsgwBIEdf/5Vzs8CDCVN9MtML1F+D7p3HsqUMLuMaW8516dI/KgRXVEkNJPTBGGPe0sxYDztDBDYeCxNgvOzCT/A39LPULHM1tUE84+Q33/yVSKhHl5Bzg+WU02ksN3DNQijs3OfYZ1E6ltnutmIgS0qnFFh54EUHhUg+p3T0sUM7RaB6sfJWU1JqsKNrFvK4+LjSjjvRfLB9eB+jckk92DverhOn/4fcbToeDU6fV+11X86EvqKpSD4wQrxvtT5iszYxDFPngeESDs/4c+n3qj+BXz2dvK6BUzseOeV/cpRkP3fwJy6j9Uklw2Y22JR9lLuhZWJdZR8dsY01zs7f0N3y+urmgOWDD5WHoJeqDR7NOU9NNXEgC/b6/5EhNfhLYCpuw4YQf7V3AGR6qFSj5R0Jjye8hZ6i74NshmLoaszFHCZX5XI7kmzI9BJCPvf1ZIdHCst+mh8yMRygJ3LvBuMkRrzENFpfTPXtTbdu1pJTPOa1JGpqSJwaM4MfkRO9sh/ddtT0ADtSkWwcNekUvrTHEs6W5IYUm3O0EhFU39m+KOedD+rv7I5uK7RqV4HRPXb4ewE1b3+SpJ44EujlkPTBu38AVI0UV24jEOTHNJprLSgV7tYLwzZLBetvXeuQJc5mfNYSQ== zero@gz0"
}

variable "gz1_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDhkGPTKaQxOkrRO4gyJFr5ewcIHVcQ+YyEeTWfK1GpwKwjMZNhbXCL0DeTxjYx+JHMqStxlKi6H/+dDJSiohbIi1MqPcpPvaY2whLPSHrV4b17tYXjnkZbowDp54QCTLgcXmfckpKosulV56t6GGXZ2PmlP85QkpyeF+8/EyS+8XFQi1GeWOwgnBvmyRRxSBwsJnMt+aJdCcIRv0sO+pKuzCJFLJ2D82sRbfLLIWOy0BkzaUMevO+sUYmmRgJ2Kis7MFNYC8c6i2GHTuRex5t0fHyCzp7qHLISAhjS4d5mmpfUeMJot7WlOiwlVTMcoQSFY0boTuDSEoVxzwu3/R2zc6Hu+RGBKZbIR6ksKZFPe64dPGJXzfX/3oyqhFQKF3A9kl9wNGq2RQweEaRtHmvToAf49iN+8q5vMdSD0XbCxJwakCj+nKLprp7yexydj6JX+lDbXtxJyF5ECuN4mE0Mi09dFJSS4e09c9bCSrIM+6AXOLpU10Kq2KqAiRxRpdV5UxheHeK89E6Z4Y8d1DggmrnQxSfvmfCctqsJVWuXC2yl+py8PmP1I9HclyBwCReTVUpOOJknQILh1lubAtboBbBK9zAIh8k13224ricxSDhiP4T+Jl6tguzxK7hQ1zQJlRipyTtaktSSikA2aMobmHo2IoifRkblJ2Wc6BiOOQ== gz1"
}

variable "gcloud_project" {
  default = "henrik-jonsson"
}

variable "gcloud_region" {
  default = "europe-west1"
}

#
# The latest image can be found with:
# $ gcloud compute images list | grep coreos-alpha
#
variable "coreos_alpha_image" {
  default = "coreos-alpha-1576-1-0-v20171026"
}

#
# The latest image can be found with:
# $ gcloud compute images list | grep coreos-beta
#
variable "coreos_beta_image" {
  default = "coreos-beta-1492-6-0-v20170906"
}

variable "admin_ip" {
  default = "130.211.84.102"
}

variable "exocore_ip" {
  default = "159.100.250.108"
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

variable "builder_enabled" {
  description = "Whether builder node is enabled."
  default = false
}

variable "dropcore_enabled" {
  description = "Whether dropcore node is enabled."
  default = false
}

variable "elentari_world_enabled" {
  description = "Whether elentari.world node is enabled."
  default = false
}

variable "guac_enabled" {
  description = "Whether guac.hkjn.me node is enabled."
  default = false
}

variable "blockpress_me_enabled" {
  description = "Whether blockpress.me node is enabled."
  default = true
}

variable "version" {}
