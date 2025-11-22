//go:build tools
// +build tools

package api

//go:generate sh -c "oapi-codegen -generate types -o ../internal/handler/gen/types.gen.go -package gen ./openapi.yaml"
//go:generate sh -c "oapi-codegen -generate server -o ../internal/handler/gen/server.gen.go -package gen ./openapi.yaml"
