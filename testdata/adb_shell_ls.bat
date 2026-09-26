pushd "%~dp0"
curl https://localhost:8080/api/adb/shell?args=ls,-la -X POST --cacert "..\internal\infra\server\assets\cert.pem" -H "Content-Type: application/json" -d@serials.json
popd