package util

import "log"

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
