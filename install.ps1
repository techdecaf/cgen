$APP_NAME = "cgen"
$BUCKET_URL = "http://github.techdecaf.io"
# $INSTALL_PATH = "$env:USERPROFILE\bin"
$INSTALL_PATH = "c:\windows\system32"

$APP_BINARY = "$INSTALL_PATH\cgen.exe"
$LATEST_STABLE = "$BUCKET_URL/$APP_NAME/latest/windows/$APP_NAME.exe"

if (Test-Path C:\Windows\System32\cgen.exe){
  echo "[Installed] $APP_NAME version: $(&$APP_NAME --version)"
}

echo "[Installing] $LATEST_STABLE"

Invoke-WebRequest -Uri "$LATEST_STABLE" -OutFile "$APP_BINARY"

echo "[Validation] $APP_NAME version: $(&$APP_NAME --version)"

cgen --help