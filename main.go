package main

import (
  "errors"
  "fmt"
  "github.com/aerospike/aerospike-client-go"
  "github.com/gorilla/mux"
  "go-as-admin/dao"
  "go-as-admin/util"
  "log"
  "net/http"
  "sort"

  "go-as-admin/config"
  _ "go-as-admin/dao"
  "go-as-admin/view"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func init() {
  log.Println("main init")
}

func main() {
  log.Printf("Booting... listening on port %d", config.FlagListenPort)

  r := mux.NewRouter()

  r.PathPrefix("/static/").Handler(
    http.StripPrefix("/static/",
      http.FileServer(http.Dir("static"))))

  r.HandleFunc("/", Logging()(index))
  r.HandleFunc("/_n/{namespace}", Logging()(namespace))
  r.HandleFunc("/_n/{namespace}/_s/{set}", Logging()(set))

  r.HandleFunc("/stats", Logging()(stats))

  http.ListenAndServe(fmt.Sprintf(":%d", config.FlagListenPort), r)
}

func index(w http.ResponseWriter, r *http.Request) {
  namespaces := dao.GetNamespacesInfo()

  indexPage := view.IndexPage{Namespaces: namespaces}
  indexPage.Render(nil, w)
}

func namespace(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  ns := vars["namespace"]
  log.Println("namespace", ns)

  sets := dao.GetSetsInfo(ns)
  log.Println("sets", sets)

  sets = util.Map(sets, func(setInfo string) (string, error) {
    log.Println("setInfo", setInfo)

    if len(setInfo) == 0 {
      return "", errors.New("cannot be empty")
    }

    infoMap := dao.ParseAerospikeInfo(setInfo)
    return fmt.Sprintf("%s", infoMap["set"]), nil
  })

  sort.Strings(sets)

  namespacePage := view.NamespacePage{Sets: sets, Namespace: ns}
  namespacePage.Render(nil, w)
}

func set(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  set, ns := vars["set"], vars["namespace"]
  log.Println("namespace", ns, "set", set)

  sets := dao.GetSetsInfo(ns)
  log.Println("sets", sets)

  sets = util.Map(sets, func(setInfo string) (string, error) {
    log.Println("setInfo", setInfo)

    if len(setInfo) == 0 {
      return "", errors.New("cannot be empty")
    }

    infoMap := dao.ParseAerospikeInfo(setInfo)
    return fmt.Sprintf("%s", infoMap["set"]), nil
  })

  sort.Strings(sets)

  //binInfo := strings.Join(dao.GetBinsInfo(ns), "")
  //log.Println("binsInfo", binInfo)
  //
  //bins := strings.Split(binInfo, ",")
  //_, bins = bins[:2], bins[2:]
  //log.Println("bins", bins)
  //sort.Strings(bins)

  headers := new([]string)
  records := new([]*aerospike.Record)
  dao.GetRecords(ns, set, headers, records)

  log.Println(headers)
  log.Println(records)

  setPage := view.SetPage{
    Sets:      sets,
    Namespace: ns,
    Records:   records,
    Set:       set,
    Bins:      *headers}

  setPage.Render(nil, w)
}

func stats(w http.ResponseWriter, r *http.Request) {
  statsPage := view.StatsPage{}
  statsPage.Render(nil, w)
}

func Logging() Middleware {
  return func(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      log.Println(r.URL)
      next(w, r)
    }
  }
}
