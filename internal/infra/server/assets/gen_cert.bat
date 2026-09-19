pushd "%~dp0"
go run "C:\Program Files\Go\src\crypto\tls\generate_cert.go" --host=localhost
popd
