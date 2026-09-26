pushd "%~dp0"
curl https://localhost:8080/api/adb/shell?args=ping,-c,4,8.8.8.8 -X POST --cacert "..\internal\infra\server\assets\cert.pem" -H "Content-Type: application/json" -d@serials.json
popd