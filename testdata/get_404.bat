pushd "%~dp0"
curl https://localhost:8080/ --cacert "..\internal\infra\server\assets\cert.pem"
popd