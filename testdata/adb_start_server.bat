pushd "%~dp0"
curl http://localhost:8080/api/adb/start-server -X POST
popd