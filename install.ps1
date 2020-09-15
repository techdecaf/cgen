$APP_NAME = "cgen"
$BUCKET_URL = "https://s3-us-west-2.amazonaws.com/github.techdecaf.io"
$INSTALL_PATH = "c:\windows"

$APP_BINARY = "$INSTALL_PATH\$APP_NAME.exe"
$LATEST_STABLE = "$BUCKET_URL/$APP_NAME/latest/windows/$APP_NAME.exe"

if (Test-Path $APP_BINARY){
  echo "[Installed] $APP_NAME version: $(&$APP_NAME --version)"
}

echo "[Installing] $LATEST_STABLE"

Invoke-WebRequest -Uri "$LATEST_STABLE" -OutFile "$APP_BINARY"

echo "[Validation] $APP_NAME version: $(&$APP_NAME --version)"

cgen --help