pushd "%~dp0"
del ..\sample.png
rem testing
curl https://localhost:8080/api/adb/pull?args=%%2Fsdcard%%2Fpictures%%2Fsample.png -X POST --cacert "..\internal\infra\server\assets\cert.pem" -H "Content-Type: application/json" -d@serials.json
popd