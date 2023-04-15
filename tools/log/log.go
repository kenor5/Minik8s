package log

import (
	"fmt"
	"log"
	"os"
)

const logPath = "../../klog.log"

func init() {
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}

func LOG(msg string) {
	log.Println(msg)
}
