package main

import (
	"log"

	logger "logger/src"
)

func main() {

	/*
		Simple stdout logger
	*/
	flags := 0
	stdoutLogger, err := logger.InitLogger("stdoutLogger", "debug", "", flags)
	if err != nil {
		log.Println("Error=", err)
		return
	}

	str := "formatted"

	stdoutLogger.Debug("New message with Debug severity")
	stdoutLogger.Debugf("New %s message with Debug severity", str)

	stdoutLogger.Info("New message with Info severity")
	stdoutLogger.Infof("New %s message with Info severity", str)

	stdoutLogger.Warning("New message with Warning severity")
	stdoutLogger.Warningf("New %s message with Warning severity", str)

	stdoutLogger.Error("New message with Error severity")
	stdoutLogger.Errorf("New %s message with Error severity", str)

	// stdoutLogger.Fatal("New message with Fatal severity")
	// stdoutLogger.Fatalf("New %s message with Fatal severity", str)

	// stdoutLogger.Panic("New message with Panic severity")
	// stdoutLogger.Panicf("New %s message with Panic severity", str)

	// change debug severy to error
	stdoutLogger.SetLogLevel("error")

	stdoutLogger.Info("This message will NOT be displayed")
	stdoutLogger.Error("This message will be displayed")

	/*
		JSON style logger
	*/
	flags |= logger.Ltimestamp | logger.LJSON
	jsonLogger, err := logger.InitLogger("jsonLogger", "debug", "", flags)
	if err != nil {
		log.Println("Error=", err)
		return
	}

	jsonLogger.Debug("New message with Debug severity")
	jsonLogger.Debugf("New %s message with Debug severity", str)

	jsonLogger.Info("New message with Info severity")
	jsonLogger.Infof("New %s message with Info severity", str)

	jsonLogger.Warning("New message with Warning severity")
	jsonLogger.Warningf("New %s message with Warning severity", str)

	jsonLogger.Error("New message with Error severity")
	jsonLogger.Errorf("New %s message with Error severity", str)
}
