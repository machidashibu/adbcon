pushd "%~dp0"
curl http://localhost:8080/api/adb/shell?args=ping%%2C-c%%2C4%%2C8.8.8.8 -X POST -H "Content-Type: application/json"
popd