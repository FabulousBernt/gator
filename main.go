package main

import (
    "fmt"
    "log"

    "github.com/FabulousBernt/gator/internal/config"
)

func main() {
    // 1. Read the config file
    cfg, err := config.Read()
    if err != nil {
        log.Fatal(err)
    }

    // 2. Set the current user and write to disk
    err = cfg.SetUser("Johnny")
    if err != nil {
        log.Fatal(err)
    }

    // 3. Read it again and print
    cfg, err = config.Read()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(cfg)
}
