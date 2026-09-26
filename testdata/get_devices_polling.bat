pushd "%~dp0"
curl https://localhost:8080/api/devices?interval=5 --cacert "..\internal\infra\server\assets\cert.pem"
popd