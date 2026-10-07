pushd "%~dp0"
curl http://localhost:8080/api/adb/logcat -X POST -H "Content-Type: application/json" -d@serials.json
popd