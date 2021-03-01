#!/bin/sh

jq -r '[.resources[] | select(.type=="aws_eip") | .instances[] | select(.attributes.public_dns != null) | .attributes.public_dns] | .[]' terraform.tfstate > dns
jq -r '[.resources[] | select(.type=="aws_eip") | .instances[] | select(.attributes.public_ip != null) | .attributes.public_ip] | .[]' terraform.tfstate > ../build/ip

zip -r build.zip ../build/
while read p; do
  scp -o "StrictHostKeyChecking no" -o "UserKnownHostsFile /dev/null" -i key build.zip ec2-user@"$p":~/ &
done <dns

