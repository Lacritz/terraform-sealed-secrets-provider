OS_TARGET = $(shell go env GOOS)_$(shell go env GOARCH)

build-local:
	mkdir -p ~/.terraform.d/plugins/datalbry/local/sealedsecret/0.0.1/$(OS_TARGET) \
	&& go build -o terraform-provider-sealedsecret \
	&& mv terraform-provider-sealedsecret ~/.terraform.d/plugins/datalbry/local/sealedsecret/0.0.1/$(OS_TARGET)