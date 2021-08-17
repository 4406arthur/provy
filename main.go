// Command leproxy implements https reverse proxy with automatic Letsencrypt usage for multiple
// hostnames/backends
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"provy/cmd"
	"provy/load-balancer/usecase"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var usageStr = `
Advance Reverse proxy

Server Options:
    -c, --config <file>              Configuration file path
    -h, --help                       Show this message
    -v, --version                    Show version
`

// usage will print out the flag options for the server.
func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func parseConfig(path string) *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	if path != "" {
		v.AddConfigPath(path)
	} else {
		v.AddConfigPath("./")
	}

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	return v
}

var version string

func printVersion() {

	fmt.Printf(`Advance Reverse proxy %s, Compiler: %s %s, Copyright (C) 2021 E.sun Bank, Inc.`,
		version,
		runtime.Compiler,
		runtime.Version())
	fmt.Println()
}

func main() {
	var configFile string
	var showVersion bool
	version = "0.0.1"
	flag.BoolVar(&showVersion, "v", false, "Print version information.")
	flag.StringVar(&configFile, "c", "", "Configuration file path.")
	flag.Usage = usage
	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
	}
	config := parseConfig(configFile)

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(4)
	logger := logrus.WithFields(logrus.Fields{
		"service": "reverse-proxy",
	})
	var rules []usecase.Member
	config.UnmarshalKey("load_balancer.rules", &rules)
	lb := usecase.NewLoadBalancerUsecase(logger, rules)

	directorFunc := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = req.Host
		versionPath := lb.Locate(req.Header.Get("Hash-Id"))
		logger.Infof("Get Hash-Id: %s distribute to path: %s", req.Header.Get("Hash-Id"), versionPath)
		req.URL.Path = req.URL.Path + versionPath
	}
	provy := cmd.NewReveseProxy(config.GetString("protocol"), config.GetString("backend"), directorFunc)
	srv := &http.Server{
		Handler:  provy.Handler,
		Addr:     config.GetString("address"),
		ErrorLog: &log.Logger{},
	}

	logger.Infof("Starting reverse proxy on %s", config.GetString("address"))
	srv.ListenAndServe()
}
