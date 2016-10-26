#
# Renan Almeida
# gosti
#

all:
	@- go build

test: all
	@- go test ./...
