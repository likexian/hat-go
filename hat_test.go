/*
 * A CLI tool for HTTP API Testing
 * https://www.likexian.com/
 *
 * Copyright 2014-2015, Li Kexian
 * Released under the Apache License, Version 2.0
 *
 */

package main


import (
    "testing"
)


func TestSimplejson(t *testing.T) {
    param := Param{
        false,
        false,
        30,
        false,
        "GET",
        "https://api.github.com/",
        map[string]string{},
        map[string]string{},
        map[string]string{},
    }

    HttpRequest(param)
}
