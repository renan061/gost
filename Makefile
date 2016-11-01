#
# Renan Almeida
# gosti
#

all:
	@- go build

test:
	@- go test ./...

ex: # Example
	@- go run example/main.go
