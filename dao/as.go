package dao

import (
  "errors"
  "fmt"
  "github.com/aerospike/aerospike-client-go"
  aerospikeLog "github.com/aerospike/aerospike-client-go/logger"
  "go-as-admin/config"
  "go-as-admin/util"
  "log"
  "strings"
  "time"
)

var client = initAsClient()

func init() {
  aerospikeLog.Logger.SetLevel(aerospikeLog.DEBUG)

  log.Println("connection established with aerospike")
}

func GetInstanceStats() map[string]interface{} {
  nodeStats, err := client.Stats()

  if err != nil {
    log.Println("error getting stats", err)
  }

  return nodeStats
}

func GetNamespacesInfo() []string {
  return info("namespaces")()
}

func GetSetsInfo(namespace string) []string {
  return info(fmt.Sprintf("sets/%s", namespace))()
}

func GetBinsInfo(namespace string) []string {
  return info(fmt.Sprintf("bins/%s", namespace))()
}

func GetRecords(namespace, set string, headers *[]string, records *[]*aerospike.Record) {
  recordset, err := client.ScanAll(nil, namespace, set)

  if err != nil {
    log.Println("cannot scan", err)
  }

  binMap := make(map[string]bool)

L:
  for {
    select {
    case record := <-recordset.Records:
      if record == nil {
        break L
      }

      // compute bins of records
      for bin := range record.Bins {
        if _, ok := binMap[bin]; !ok {
          binMap[bin] = true
          *headers = append(*headers, bin)
        }
      }

      *records = append(*records, record)
    case err := <-recordset.Errors:
      // if there was an error, stop
      log.Println("error reading record set", err)
    }
  }
}

func ParseAerospikeInfo(info string) map[string]string {
  infoPair := strings.Split(info, ":")

  infoMap := make(map[string]string, len(infoPair))
  for _, infoValue := range infoPair {
    s := strings.Split(infoValue, "=")
    key, value := s[0], s[1]

    infoMap[key] = value
  }

  return infoMap
}

func initAsClient() *aerospike.Client {
  client, err := aerospike.NewClient(config.AS.Host, config.AS.Port)
  if err != nil {
    log.Fatal("failed connect to aerospike", err)
  }

  return client
}

func initAsConnection() *aerospike.Connection {
  connection, err := aerospike.NewConnection(config.FlagAsURL, 50*time.Millisecond)

  if err != nil {
    log.Println("failed connect to aerospike", err)
  }

  return connection
}

func info(infoKey string) func() []string {
  return func() []string {
    info, err := aerospike.RequestInfo(initAsConnection(), infoKey)
    util.LogOnError("failed getting asinfo", err)

    fo, ok := info[infoKey]

    if !ok {
      util.LogOnError(fmt.Sprintf("cannot find `%s` in info", infoKey), errors.New(""))
    }

    return strings.Split(fo, ";")
  }
}
