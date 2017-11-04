variable "gcloud_credentials" {
  default = ".gcp/tf-dns-editor.json"
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
  default = "coreos-alpha-1520-1-0-v20170906"
}

#
# The latest image can be found with:
# $ gcloud compute images list | grep coreos-beta
#
variable "coreos_beta_image" {
  default = "coreos-beta-1492-6-0-v20170906"
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
