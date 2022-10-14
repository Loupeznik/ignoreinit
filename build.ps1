$appName = "ignoreinit"
$archs = @("amd64", "386", "arm64")
$module = "github.com/loupeznik/$appName"

foreach ($arch in $archs)
{
    $Env:GOOS="windows"; $Env:GOARCH=$arch; go build -o bin/$appname-$arch-win.exe $module

    if ($arch -ne "386") {
        $Env:GOOS="darwin"; $Env:GOARCH=$arch; go build -o bin/$appname-$arch-darwin $module
    }
    
    $Env:GOOS="linux"; $Env:GOARCH=$arch; go build -o bin/$appname-$arch-linux $module
}
