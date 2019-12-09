```bash
sh -c "$(curl -fsSL https://raw.github.com/techdecaf/{{.CI_PROJECT_NAME}}/master/install.sh)"
```

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.github.com/techdecaf/{{.CI_PROJECT_NAME}}/master/install.ps1'))
```

Download Links

- [windows]({{.DOWNLOAD_URI}}/windows/{{.CI_PROJECT_NAME}}.exe)
- [mac]({{.DOWNLOAD_URI}}/latest/darwin/{{.CI_PROJECT_NAME}})
- [linux]({{.DOWNLOAD_URI}}/latest/linux/{{.CI_PROJECT_NAME}})

To install {{.CI_PROJECT_NAME}}, use the provided script, simlink it or place it in any directory that is part of your path.
i.e. `/usr/local/bin` or `c:\windows`
