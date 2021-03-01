output "servers" {
  value = [aws_instance.airwave.*.public_ip]
}
