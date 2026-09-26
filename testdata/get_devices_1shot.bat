pushd "%~dp0"
curl https://localhost:8080/api/devices --cacert "..\internal\infra\server\assets\cert.pem"
popd