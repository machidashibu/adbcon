pushd "%~dp0"
curl http://localhost:8080/api/adb/shell?args=ls%%2C-la -X POST -H "Content-Type: application/json" -d@serials.json
popd