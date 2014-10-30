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


const (
    VERSION = "0.7.0"
    HELP_INFO = `Usage:
    hat [FLAGS] [METHOD] [URL] [OPTIONS]

FLAGS: specify data type of POST and PUT
    -j, --json  POST/PUT data as json encode (default)
    -f, --form  POST/PUT data as form encode

FLAGS specify verbose
    -v, --verbose

FLAGS specify show request/response time and download speed
    -t

FLAGS specify request and response total timeout
    --timeout=<int>

FLAGS specify show the version
    -V, --version

FLAGS specify show the help
    -h, --help

METHOD: specify http request method
    GET         HTTP GET        GET / HTTP/1.1 (default)
    POST        HTTP POST       POST / HTTP/1.1
    PUT         HTTP PUT        PUT / HTTP/1.1
    DELETE      HTTP DELETE     DELETE / HTTP/1.1

URL: the HTTP URL for request, support http and https
    <empty>     for http://127.0.0.1/ (default)
    :8080       for http://127.0.0.1:8080/
    :8080/api/  for http://127.0.0.1:8080/api/
    /api/       for http://127.0.0.1/api/

OPTIONS: the HTTP headers and HTTP body, add as many as you want
    key:value   HTTP headers    for example User-Agent:HAT/0.1.0
    key=value   HTTP body       for example name=likexian
    key?=value  HTTP query      for example name?=likexian set URL to /?name=likexian
`
    URL_DEFAULT = "http://127.0.0.1"
)


type Param struct {
    Verbose bool                `json:"verbose"`
    Timer   bool                `json:"timer"`
    Timeout int                 `json:"timeout"`
    IsJson  bool                `json:"is_json"`
    Method  string              `json:"method"`
    URL     string              `json:"url"`
    Header  map[string]string   `json:"header"`
    Query   map[string]string   `json:"query"`
    Data    map[string]string   `json:"data"`
}


func Version() string {
    return VERSION
}


func Author() string {
    return "[Li Kexian](http://www.likexian.com/)"
}


func Copyright() string {
    return "Copyright 2014, Kexian Li"
}


func License() string {
    return "Apache License, Version 2.0"
}


func main() {
    param := Param{
        false,
        false,
        30,
        true,
        "GET",
        URL_DEFAULT,
        map[string]string{},
        map[string]string{},
        map[string]string{},
    }

    args := os.Args
    for i:=1; i<len(args); i++ {
        v := args[i]
        if v[0] == ':' || v[0] == '/' {
            param.URL = URL_DEFAULT + v
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

            if v == "-t" {
                param.Timer = true
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

            if v == "-V" || v == "--version" {
                fmt.Println("HAT version " + Version())
                fmt.Println(Copyright())
                fmt.Println(License())
                os.Exit(0)
            }

            if v == "-h" || v == "--help" {
                fmt.Println(HELP_INFO)
                os.Exit(0)
            }

            continue
        }

        if len(v) < 7 {
            _v := strings.ToUpper(v)
            if _v == "GET" || _v == "POST" || _v == "PUT" || _v == "DELETE" {
                param.Method = _v
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

        if strings.Contains(v, "?=") {
            vv := strings.SplitN(v, "?=", 2)
            if !strings.Contains(vv[0], "=") && !strings.Contains(vv[0], ":") {
                param.Query[vv[0]] = vv[1]
                continue
            }
        }

        if strings.Contains(v, "=") {
            vv := strings.SplitN(v, "=", 2)
            if !strings.Contains(vv[0], ":") {
                param.Data[vv[0]] = vv[1]
                continue
            }
        }

        if strings.Contains(v, ":") {
            vv := strings.SplitN(v, ":", 2)
            _v := vv[1]
            if strings.Contains(vv[1], "/") {
                _vv := strings.SplitN(vv[1], "/", 2)
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

    q_data := url.Values{}
    for k, v := range param.Query {
        q_data.Add(k, v)
    }

    q_query := q_data.Encode()
    if q_query != "" {
        if !strings.Contains(param.URL, "?") {
            param.URL += "?" + q_query
        } else {
            param.URL += "&" + q_query
        }
    }

    w_body := ""
    v_w_body := ""
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
            w_body = data
            v_w_body, err = simplejson.PrettyDumps(data_json)
            if err != nil {
                v_w_body = w_body
            }
        } else {
            data := url.Values{}
            for k, v := range param.Data {
                data.Add(k, v)
            }
            w_body = data.Encode()
            v_w_body = w_body
        }
    }

    request, err := http.NewRequest(param.Method, param.URL, bytes.NewBuffer([]byte(w_body)))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    request.Header.Set("Accept", "*/*")
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
    w_start_time := time.Now().UnixNano() / 1e6
    response, err := client.Do(request)
    w_end_time := time.Now().UnixNano() / 1e6
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

        query := request.URL.RawQuery
        if query != "" {
            query = "?" + query
        }

        header := fmt.Sprintf("> %s %s%s HTTP/1.1\r\n", param.Method, path, query)
        header += fmt.Sprintf("> Host: %s\r\n", request.URL.Host)
        header += "> Accept-Encoding: gzip\r\n"
        for k, v := range request.Header {
            header += fmt.Sprintf("> %s: %s\r\n", k, v[0])
        }

        fmt.Print(header)
        fmt.Println(">")

        if w_body != "" {
            fmt.Println(v_w_body)
            fmt.Println(">")
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
        fmt.Print(fmt.Sprintf("< %s %s\r\n", response.Proto, response.Status))
        for k, v := range response.Header {
            fmt.Print(fmt.Sprintf("< %s: %s\r\n", k, v[0]))
        }
        fmt.Print("< \r\n")
    }

    r_start_time := time.Now().UnixNano() / 1e6
    r_body, err := ioutil.ReadAll(response.Body)
    r_end_time := time.Now().UnixNano() / 1e6
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    text := string(r_body)
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
    if param.Timer {
        w_time := w_end_time - w_start_time
        r_time := r_end_time - r_start_time
        fmt.Println("")
        fmt.Println(fmt.Sprintf("request:\t%.2fs", float64(w_time) / 1000.0))
        fmt.Println(fmt.Sprintf("response:\t%.2fs", float64(r_time) / 1000.0))
        if (r_time == 0) {
            r_time = 1
        }
        fmt.Println(fmt.Sprintf("download:\t%dk/s", int64(len(r_body) * 1000 / 1024) / r_time))
    }
}
