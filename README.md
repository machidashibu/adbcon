# ADB Console Tool

This tool provides GUI for ADB command on WEB browser.
***adbcon*** is an acronym of ***ADB Console***.

> [!CAUTION]
> This tool accesses local resources (e.g., smartphones, local storage). You **MUST NOT** expose the TCP port used by this tool to the internet.

## Sample

![Sample screen](docs/sample.png)

## Get started

### Build

First, please create certification and private key file (PEM format) for localhost, and put to  following directory.

There are needed at build.

* `internal/infra/server/assets/cert.pem`
* `internal/infra/server/assets/key.pem`

> [!NOTE]
> No problem if It is self-signed certificate.

Readied sample bat command to create the files by `generate_cert.go` in Go sources. (for windows only)

Run `internal/infra/server/assets/gen_cert.bat`.

> [!NOTE]
> It assumes installed Go sources to `C:\Program Files\Go`)

Execute following command, executable file is generated. (`adbcon.exe` for windows)

```bash
$ go build
```

### Run

Run the generated executable file.

```bash
$ adbcon
```

Access to `https://localhost:8080/gui` by WEB browser.

You may ignore if the browser notifies means "access to not secure page".

### References

* [Open API v3.0.4](https://spec.openapis.org/oas/v3.0.4.html#encoding-object)
* [oapi-codegen](https://github.com/oapi-codegen) - Apache-2.0 license
* [Echo v4](https://echo.labstack.com/) - MIT license
