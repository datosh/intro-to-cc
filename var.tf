variable "ssh_keys" {
  type = list(object({
    user = string
    key = string
  }))

  description = "list of public ssh keys that have access to the VM"

  default = [
      {
        user = "datosh"
        key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMN+egbaEF+2hl7TMtjEpCcnFmQbgqOC1v1ijFpuTJGW datosh@bingo"
      }
  ]
}
