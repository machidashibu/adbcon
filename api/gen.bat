pushd "%~dp0"
go tool oapi-codegen -config config.yaml openapi.yaml
popd
