pushd "%~dp0"
curl https://localhost:8080/ -X POST --cacert "..\internal\infra\server\assets\cert.pem"
popd