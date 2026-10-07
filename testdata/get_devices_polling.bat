pushd "%~dp0"
curl http://localhost:8080/api/devices?interval=5
popd