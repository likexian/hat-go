/*
 * A CLI tool for HTTP API Testing
 * http://www.likexian.com/
 *
 * Copyright 2014, Kexian Li
 * Released under the Apache License, Version 2.0
 *
 */

package main


import (
    "fmt"
    "os"
    "strings"
    "net/http"
    "io/ioutil"
    "bytes"
    "strconv"
    "net/url"
    "time"
    "github.com/likexian/simplejson-go"
)


type Param struct {
    Verbose bool                `json:"verbose"`
    Timeout int                 `json:"timeout"`
    IsJson  bool                `json:"is_json"`
    Method  string              `json:"method"`
    URL     string              `json:"url"`
    Header  map[string]string   `json:"header"`
    Data    map[string]string   `json:"data"`
}


func Version() string {
    return "0.3.0"
}


func Author() string {
    return "[Li Kexian](http://www.likexian.com/)"
}


func License() string {
    return "Apache License, Version 2.0"
}


func main() {
    param := Param{
        false,
        30,
        false,
        "GET",
        "http://127.0.0.1",
        map[string]string{},
        map[string]string{},
    }

    args := os.Args
    for i:=1; i<len(args); i++ {
        v := args[i]
        if v[0] == ':' || v[0] == '/' {
            param.URL += v
            continue
        }

        if v[0] == '-' {
            if v == "-j" || v == "--json" {
                param.IsJson = true
                continue
            } else if v == "-f" || v == "--form" {
                param.IsJson = false
                continue
            }

            if v == "-v" || v == "--verbose" {
                param.Verbose = true
                continue
            }

            if len(v) > 10 && v[:10] == "--timeout=" {
                timeout, err := strconv.Atoi(v[10:])
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
                param.Timeout = timeout
                continue
            }
        }

        if len(v) > 7 && v[:7] == "http://" {
            param.URL = v
            continue
        } else if len(v) > 8 && v[:8] == "https://" {
            param.URL = v
            continue
        }

        if len(v) < 6 {
            _v := strings.ToUpper(v)
            if _v == "GET" || _v == "POST" || _v == "PUT" || _v == "DELETE" {
                param.Method = _v
                continue
            }
        }

        if strings.Contains(v, "=") {
            vv := strings.Split(v, "=")
            if strings.Contains(vv[0], ":") {
                param.Header[vv[0]] = vv[1]
            } else {
                param.Data[vv[0]] = vv[1]
            }
            continue
        }

        if strings.Contains(v, ":") {
            vv := strings.Split(v, ":")
            _v := vv[1]
            if strings.Contains(vv[1], "/") {
                _vv := strings.Split(vv[1], "/")
                _v = _vv[0]
            }
            if len(_v) > 0 && len(_v) < 6 {
                _, err := strconv.Atoi(_v)
                if err == nil {
                    param.URL = v
                    continue
                }
            }
            param.Header[vv[0]] = vv[1]
            continue
        }

        param.URL = v
    }

    HttpRequest(param)
}


func HttpRequest(param Param) {
    if len(param.URL) <= 7 || (param.URL[:8] != "https://" && param.URL[:7] != "http://") {
        param.URL = "http://" + param.URL
    }

    body := ""
    if param.Method == "POST" || param.Method == "PUT" {
        if param.IsJson {
            data_json := simplejson.New()
            for k, v := range param.Data {
                data_json.Set(k, v)
            }
            data, err := simplejson.Dumps(data_json)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
            body = data
        } else {
            data := url.Values{}
            for k, v := range param.Data {
                data.Add(k, v)
            }
            body = data.Encode()
        }
    }

    request, err := http.NewRequest(param.Method, param.URL, bytes.NewBuffer([]byte(body)))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    request.Header.Set("Accept", "*/*")
    request.Header.Set("Accept-Encoding", "gzip")
    request.Header.Set("User-Agent", fmt.Sprintf("HAT/%s (i@likexian.com)", Version()))

    if param.Method == "POST" || param.Method == "PUT" {
        if param.IsJson {
            request.Header.Set("Accept", "application/json")
            request.Header.Set("Content-Type", "application/json")
        } else {
            request.Header.Set("Accept", "*/*")
            request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        }
    }

    for k, v := range param.Header {
        request.Header.Set(k, v)
    }

    client := &http.Client{Timeout: time.Duration(param.Timeout) * time.Second}
    response, err := client.Do(request)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer response.Body.Close()

    if param.Verbose {
        path := request.URL.Path
        if path == "" {
            path = "/"
        }

        header := fmt.Sprintf("> %s %s HTTP/1.1\r\n", param.Method, path)
        header += fmt.Sprintf("> Host: %s\r\n", request.URL.Host)
        for k, v := range request.Header {
            header += fmt.Sprintf("> %s: %s\r\n", k, v[0])
        }

        fmt.Print(header)
        fmt.Print("> \r\n")

        if body != "" {
            fmt.Print(body)
            fmt.Print("> \r\n")
        }
    }

    is_json := false
    for k, v := range response.Header {
        if k == "Content-Type" {
            vv := strings.Split(v[0], ";")
            if strings.ToLower(vv[0]) == "application/json" {
                is_json = true
                break
            }
        }
    }

    if param.Verbose {
        for k, v := range response.Header {
            fmt.Print(fmt.Sprintf("< %s: %s\r\n", k, v[0]))
        }
        fmt.Print("< \r\n")
    }

    data, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    text := string(data)
    if is_json {
        data, err := simplejson.Loads(text)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        text, _ = simplejson.PrettyDumps(data)
    } else {
        if text != "" && (text[0] == '{' || text[0] == '[') && (text[len(text) - 1] == '}' || text[len(text) - 1] == ']') {
            data, err := simplejson.Loads(text)
            if err == nil {
                text, _ = simplejson.PrettyDumps(data)
            }
        }
    }

    fmt.Println(text)
}
