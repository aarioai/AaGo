package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func Println(args ...any) {
	ns := time.Now().Format("2006-01-02 15:04:05")
	for _, arg := range args {
		msg, ok := arg.(string)
		if !ok {
			s, err := json.Marshal(arg)
			if err != nil {
				msg = err.Error()
			} else {
				msg = string(s)
			}
		}
		log.Println(msg)
		fmt.Println(ns + " " + msg)
	}
}
