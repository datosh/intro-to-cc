variable "ssh_keys" {
  type = list(object({
    user = string
    key = string
  }))

  description = "list of public ssh keys that have access to the VM"

  default = [
      {
        user = "datosh"
        key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIPR9ON5+xDrVWCogP1c0Wh3Ne+6rNHh7qFFxIN3Up4o4 fabiankammel@Fabians-MBP.local"
      }
  ]
}
