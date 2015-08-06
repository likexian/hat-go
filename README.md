# HAT

A CLI tool for HTTP API Testing

[![Build Status](https://secure.travis-ci.org/likexian/hat-go.png)](https://secure.travis-ci.org/likexian/hat-go)

## Overview

HAT (HTTP API Testing) is a command line HTTP client. Its goal is to make HTTP API testing as easy as possible. It providers powerful ablitity to make HTTP request but with very simple arguments. HAT can be used for HTTP API testing, debugging and you can as well use it as CURL.

## Install

### Download

- OSX: https://github.com/likexian/hat-go/releases/download/v0.8.4/hat_osx_0.8.4.tar.gz
- Linux 32Bit: https://github.com/likexian/hat-go/releases/download/v0.8.4/hat_linux_x86_0.8.4.tar.gz
- Linux 64Bit: https://github.com/likexian/hat-go/releases/download/v0.8.4/hat_linux_x86_64_0.8.4.tar.gz

### Untar and move

    tar zxvf hat_*_0.8.4.tar.gz
    sudo mv hat /usr/local/bin

### Test it

    hat --version

## Documentation

hat [FLAGS] [METHOD] [URL] [OPTIONS]

### Base usage

HAT is so powerful that it works without any arguments.

    hat

*=> This will make a GET request to http://127.0.0.1/.*

### FLAGS

FLAGS specify data type of POST and PUT

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

### METHOD

METHOD specify http request method

    GET         HTTP GET        GET / HTTP/1.1 (default)
    POST        HTTP POST       POST / HTTP/1.1
    PUT         HTTP PUT        PUT / HTTP/1.1
    DELETE      HTTP DELETE     DELETE / HTTP/1.1

### URL

URL is the HTTP URL for request, support http and https

    <empty>     for http://127.0.0.1/ (default)
    :8080       for http://127.0.0.1:8080/
    :8080/api/  for http://127.0.0.1:8080/api/
    /api/       for http://127.0.0.1/api/

### OPTIONS

OPTIONS can specify the HTTP headers, HTTP body and HTTP URL query, add as many as you want

    key:value   HTTP headers        for example User-Agent:HAT/0.1.0
    key=value   HTTP body           for example name=likexian
    key?=value  HTTP query          for example name?=likexian set URL to /?name=likexian
    @file       HTTP POST/PUT Data  for example @post.json will use post.json content as post/put body data

## EXAMPLE

Just get a url

    hat http://127.0.0.1/

Get a url and specify the query

    hat http://127.0.0.1/ name?=likexian pass?=xxxxxxxx

*=> This will make a GET request to http://127.0.0.1/?name=likexian&pass=xxxxxxxx*

Get a url and specify the headers

    hat http://127.0.0.1/ User-Agent:HAT/0.1.0 X-Forward-For:192.168.1.1

POST data to url (json)

    hat post http://127.0.0.1/api/user name=likexian pass=xxxxxxxx

POST data to url (form)

    hat -f post http://127.0.0.1/api/user name=likexian pass=xxxxxxxx

POST file content as raw body data

    hat @post.json http://127.0.0.1/api/user

PUT file content as raw body data

    hat @post.json put http://127.0.0.1/api/user/1

*=> If you are using the @file option as above, the content-type will be set depending on the format of the file,
content-type will be json if the content is json (same as -f) formating, otherwise form (as -f).*

## LICENSE

Copyright 2014-2015, Li Kexian

Apache License, Version 2.0

## About

- [Li Kexian](https://www.likexian.com/)
