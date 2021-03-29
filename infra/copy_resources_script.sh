#!/bin/sh

jq -r '[.resources[] | select(.type=="aws_eip") | .instances[] | select(.attributes.public_dns != null) | .attributes.public_dns] | .[]' terraform.tfstate > dns
jq -r '[.resources[] | select(.type=="aws_eip") | .instances[] | select(.attributes.public_ip != null) | .attributes.public_ip] | .[]' terraform.tfstate > ../build/ip

zip -r build.zip ../build/

COUNT=0
while read p; do
  COUNT=$(expr $COUNT + 1)
  scp -o "StrictHostKeyChecking no" -o "UserKnownHostsFile /dev/null" -i key build.zip ec2-user@"$p":~/
done <dns

cd ../generator && go run main.go $COUNT && cd ../infra

INDEX=0
while read p; do
  ssh -o "StrictHostKeyChecking no" -o "UserKnownHostsFile /dev/null" -i key ec2-user@"$p" "echo $INDEX && exit" < /dev/null
  INDEX=$(expr $INDEX + 1)
done <dns