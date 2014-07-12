package main

import (

    "os"
    "log"
    "runtime"
)

func LogInfo(args error) {

    log.Printf("[+] Notify: %s\n", args)
}

func LogErr(args error) {

    log.Printf("[!] Notify: %s\n", args)
}

func LogErrF(args error) {

    log.Fatalf("[!] Notify: %s\n", args)
}

func GetHostName() string {

    n, e := os.Hostname()
    if e != nil {

        n = "unknown"
    }
    return n
}

func GetHostOS() (os string) {

    if (runtime.GOOS != "") {

        return runtime.GOOS
    }

    return "unknown"
}