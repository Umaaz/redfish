package main

import (
	"context"
	"github.com/Umaaz/redfish/pkg/config"
	"os"
)

func main() {
	os.Setenv("PKL_EXEC", "/home/bdonnell/repo/play/go/pickle/pkl")

	ret, err := config.LoadFromPath(context.Background(), "/home/bdonnell/repo/github/Umaaz/redfish/example.pkl")
	if err != nil {
		panic(err)
	}

	println(ret)
}
