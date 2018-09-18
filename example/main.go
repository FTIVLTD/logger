package main

import (
	"log"

	logger "github.com/FTIVLTD/logger/src"
)

func main() {
	lg, err := logger.InitLogger("AppName", "debug", "")

	if err != nil {
		log.Println("Error=", err)
		return
	}

	str := "formatted"

	lg.Debug("New message with Debug severity")
	lg.Debugf("New %s message with Debug severity", str)

	lg.Info("New message with Info severity")
	lg.Infof("New %s message with Info severity", str)

	lg.Warning("New message with Warning severity")
	lg.Warningf("New %s message with Warning severity", str)

	lg.Error("New message with Error severity")
	lg.Errorf("New %s message with Error severity", str)

	// lg.Fatal("New message with Fatal severity")
	// lg.Fatalf("New %s message with Fatal severity", str)

	// lg.Panic("New message with Panic severity")
	// lg.Panicf("New %s message with Panic severity", str)

	// change debug severy to error
	lg.SetLogLevel("error")

	lg.Info("This message will NOT be displayed")
	lg.Error("This message will be displayed")
}
