package config

import (
  "flag"
  "log"
  "os"
  "strconv"
  "strings"
)

var FlagListenPort int
var FlagShowHelp bool
var FlagAsURL string

var AS struct {
  Host string
  Port int
}

func init() {
  flag.BoolVar(&FlagShowHelp, "help", false, "show help")
  flag.IntVar(&FlagListenPort, "listen", 9090, "listening port")
  flag.StringVar(&FlagAsURL, "host", "127.0.0.1:3000", "aerospike endpoint")

  flag.Parse()
  if FlagShowHelp {
    flag.Usage()
    os.Exit(0)
  }

  hostAndPort := strings.Split(FlagAsURL, ":")
  if len(hostAndPort) != 2 {
    log.Fatal("could not find proper config for aerospike")
  }

  AS.Host = hostAndPort[0]
  port, err := strconv.Atoi(hostAndPort[1])

  if err != nil {
    log.Fatal("failed parsing aerospike endpoint")
  }

  AS.Port = port
}
