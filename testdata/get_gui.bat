pushd "%~dp0"
curl https://localhost:8080/gui --cacert "..\internal\infra\server\assets\cert.pem"
popd