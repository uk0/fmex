package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/uk0/fmex/cmd"
	"os"
	"runtime"
	"time"
)

const VERSION = "v1"

func main() {
	app := cli.NewApp()
	app.Name = "Femx client"
	app.Usage = "a tool for driving femx cli"
	app.Author = "Zhangjianxin"
	app.Email = "zhangjianxinnet@gmail.com"
	app.Version = fmt.Sprintf("%s %s/%s %s", VERSION,
		runtime.GOOS, runtime.GOARCH, runtime.Version())
	app.EnableBashCompletion = true
	app.Compiled = time.Now()

	app.Commands = []cli.Command{
		cmd.NewBalanceCommand(),
		cmd.NewFmexCommand(),
		cmd.NewServerCommand(),
	}

	app.Run(os.Args)
}