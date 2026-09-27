pushd "%~dp0"
rem remove all files
curl https://localhost:8080/api/adb/shell?args=rm%%2C%%2Fsdcard%%2Fpictures%%2Fsample.png -X POST --cacert "..\internal\infra\server\assets\cert.pem" -H "Content-Type: application/json" -d@serials.json
rem testing
curl https://localhost:8080/api/adb/push?args=docs%%2Fsample.png%%2C%%2Fsdcard%%2Fpictures -X POST --cacert "..\internal\infra\server\assets\cert.pem" -H "Content-Type: application/json" -d@serials.json
rem check result
curl https://localhost:8080/api/adb/shell?args=ls%%2C-la%%2C%%2Fsdcard%%2Fpictures -X POST --cacert "..\internal\infra\server\assets\cert.pem" -H "Content-Type: application/json" -d@serials.json
popd