package util

import (
  "errors"
  "fmt"
  "github.com/aerospike/aerospike-client-go"
  "log"
  "reflect"
  "strconv"
  "strings"
)

func PanicOnError(msg string, err error) {
  if err != nil {
    log.Fatal(msg, err)
  }
}

func LogOnError(msg string, err error) {
  if err != nil {
    log.Println(msg, err)
  }
}

func ToString(fieldValue interface{}) string {
  valueType := reflect.ValueOf(fieldValue)
  switch valueType.Kind() {
  case reflect.String:
    return fieldValue.(string)
  case reflect.Int:
    return strconv.Itoa(fieldValue.(int))
  case reflect.Slice:
    if result, err := Join(fieldValue, ","); err != nil {
      return fmt.Sprintf("%v", err)
    } else {
      return result
    }
  default:
    return "Not string"
  }
}

func Join(slice interface{}, separator string) (string, error) {
  fmt.Printf("%T", slice)
  var results []string

  switch t := slice.(type) {
  case []string:
  case []int:
  case []uint:
  case []uint8:
    for _, value := range t {
      results = append(results, fmt.Sprintf("%v", value))
    }
  default:
    return "", errors.New("failed to transform slice into string")
  }

  return strings.Join(results, separator), nil
}

func Slice2Map(memo []string, fn func(string) (interface{}, error)) map[string]interface{} {
  mmap := make(map[string]interface{}, len(memo))
  for _, v := range memo {
    if fnValue, err := fn(v); err == nil {
      mmap[v] = fnValue
    }
  }
  return mmap
}

func Map2Slice(memo map[string]interface{}, fn func(string, interface{}) (string, error)) []string {
  sslice := make([]string, 0, len(memo))
  for key, value := range memo {
    if fnValue, err := fn(key, value); err == nil {
      sslice = append(sslice, fnValue)
    }
  }
  return sslice
}

func Map(memo []string, fn func(string) (string, error)) []string {
  slice := make([]string, 0, len(memo))

  for _, value := range memo {
    if fnValue, err := fn(value); err == nil {
      slice = append(slice, fnValue)
    }
  }

  return slice
}

func Filter(records []*aerospike.Record, fn func(record *aerospike.Record) (bool, error)) ([]*aerospike.Record, error) {
  results := make([]*aerospike.Record, 0, len(records))
  for _, r := range records {
    if ok, err := fn(r); ok && err == nil {
      results = append(results, r)
    } else if err != nil {
      return records, err
    }
  }

  return results, nil
}
