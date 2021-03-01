build:
	cd src 
	GOOS=linux GOARCH=amd64 go build -v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o ../build/bin/app .
	cd ..

init:
	terraform init infra

plan:
	terraform plan infra

apply:
	terraform apply --auto-approve infra 
	cd infra
	./copy_resources_script.sh
	cd ..

destroy:
	terraform destroy --auto-approve infra

plot:
	python utils/plot.py

clean:
	-rm -r build
	rm infra/dns
