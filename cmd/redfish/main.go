package main

import (
	"context"
	"os"

	"github.com/Umaaz/redfish/pkg/config/pkl/gen/appconfig"
)

func main() {
	os.Setenv("PKL_EXEC", "/home/bdonnell/repo/play/go/pickle/pkl")

	ret, err := appconfig.LoadFromPath(context.Background(), "/home/bdonnell/repo/github/Umaaz/redfish/example.pkl")
	if err != nil {
		panic(err)
	}

	println(ret)
}
