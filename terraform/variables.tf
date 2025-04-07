variable "aws_region" {
  default = "us-east-1"
}

variable "public_key_path" {
  description = "Path to your public SSH key"
  default     = "~/.ssh/id_ed25519.pub"
}

variable "ami_id" {
  description = "Amazon Linux 2 AMI"
  default     = "ami-0a9a48ce4458e384e" # Amazon Linux 2 x86_64 (us-east-1)
}
