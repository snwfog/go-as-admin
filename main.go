package main

import (
  "errors"
  "fmt"
  "github.com/aerospike/aerospike-client-go"
  "github.com/gorilla/mux"
  "go-as-admin/config"
  "go-as-admin/dao"
  _ "go-as-admin/dao"
  "go-as-admin/util"
  "go-as-admin/view"
  "log"
  "net/http"
  "sort"
  "strings"
)

var (
  QS_FilterKey       = "filter"
  QS_FilterSeparator = ":"
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

  // binInfo := strings.Join(dao.GetBinsInfo(ns), "")
  // log.Println("binsInfo", binInfo)
  //
  // bins := strings.Split(binInfo, ",")
  // _, bins = bins[:2], bins[2:]
  // log.Println("bins", bins)
  // sort.Strings(bins)

  var headers []string
  var records []*aerospike.Record
  dao.GetRecords(ns, set, &headers, &records)

  log.Println(headers)
  log.Println(records)

  if len(r.FormValue(QS_FilterKey)) > 0 {
    var err error
    records, err = filterRecords(r, records)
    if err != nil {
      log.Println("not filtered", err)
    }
  }

  setPage := view.SetPage{
    Sets:      sets,
    Namespace: ns,
    Records:   records,
    Set:       set,
    Bins:      headers}

  setPage.Render(nil, w)
}

func stats(w http.ResponseWriter, r *http.Request) {
  statsPage := view.StatsPage{}
  statsPage.Render(nil, w)
}

func filterRecords(r *http.Request, records []*aerospike.Record) ([]*aerospike.Record, error) {
  filter := r.FormValue(QS_FilterKey)
  log.Println("Filtering by", filter)

  filters := strings.Split(filter, QS_FilterSeparator)
  if len(filters) != 2 {
    return records, errors.New("error wrong filter format")
  } else {
    filterKey, filterValue := filters[0], filters[1]
    records, err := util.Filter(records, func(record *aerospike.Record) (bool, error) {
      for bin, value := range record.Bins {
        if bin == filterKey {

          recordValue := fmt.Sprintf("%v", value)
          return recordValue == filterValue, nil
        }
      }

      return false, errors.New("cannot be filtered, bin does not exist")
    })

    return records, err
  }
}

func Logging() Middleware {
  return func(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      log.Println(r.URL)
      next(w, r)
    }
  }
}
