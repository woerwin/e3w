package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Guazi-inc/e3w/conf"
	"github.com/Guazi-inc/e3w/e3ch"
	"github.com/Guazi-inc/e3w/routers"
	"github.com/coreos/etcd/version"
	"github.com/gin-gonic/gin"
)

const (
	PROGRAM_NAME    = "e3w"
	PROGRAM_VERSION = "0.0.2"
)

var configFilepath string

func init() {
	flag.StringVar(&configFilepath, "conf", "conf/config.default.ini", "config file path")
	rev := flag.Bool("rev", false, "print rev")
	// addr := flag.String("addr", "localhost:2579", "etcd address online")
	// user := flag.String("user", "", "etcd user")
	flag.Parse()

	if *rev {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			PROGRAM_NAME, PROGRAM_VERSION,
			version.Version,
		)
		os.Exit(0)
	}
}

func main() {
	config, err := conf.NewInit(configFilepath)
	if err != nil {
		panic(err)
	}

	client, err := e3ch.InitE3chClient(config)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.UseRawPath = true
	routers.InitRouters(router, config, client)
	router.Run(":" + config.App.Port)
}
