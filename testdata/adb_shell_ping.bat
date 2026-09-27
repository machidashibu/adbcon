pushd "%~dp0"
curl https://localhost:8080/api/adb/shell?args=ping%%2C-c%%2C4%%2C8.8.8.8 -X POST --cacert "..\internal\infra\server\assets\cert.pem" -H "Content-Type: application/json" -d@serials.json
popd