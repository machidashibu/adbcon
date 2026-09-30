pushd "%~dp0"
curl http://localhost:8080/api/adb/kill-server -X POST
popd