pushd "%~dp0"
curl https://localhost:8080/api/devices?interval=-1 --cacert "..\internal\infra\server\assets\cert.pem"
popd