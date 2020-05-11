install-dev: setup build link

install: setup build move clean

clean:
	@rm -rf release/ .terraform/ terraform-provider-mongodb terraform.tfstate terraform.tfstate.backup

setup:
	@git pull origin master
	@mkdir -p ~/.terraform.d/plugins

build:
	@go get
	@go build -o terraform-provider-mongodb

link:
	@ln -sf $(shell pwd)/terraform-provider-mongodb ~/.terraform.d/plugins

move:
	@mv terraform-provider-mongodb ~/terraform.d/plugins/

uninstall: clean
	@rm ~/.terraform.d/plugins/terraform-provider-mongodb

release-%: build-release-% publish-release-%
	@

build-release-%:
	@echo "Building $*"
	@./scripts/build-release.sh $*
	@echo "Build completed!"

publish-release-%:
	@echo "Publish release $* manually to Github!"
