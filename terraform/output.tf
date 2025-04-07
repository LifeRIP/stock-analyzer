output "public_ip" {
  description = "Public IP of the application server"
  value = aws_eip.app_server_eip.public_ip
}

output "ssh_command" {
  description = "SSH command to connect to the application server"
  value = "ssh ec2-user@${aws_instance.app_server_eip.public_ip}"
}