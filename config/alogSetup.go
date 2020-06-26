package config

import (
	"log"
	"os"
)

//If you have multiple config files, name the one that sets the log output alphabetically first
//When go runs init functions it does it by filename so closer to the beginning of the alphabet go first
//So if you have log functions in your other func init, they will print to the file instead of stderr

func init() {
	//Creates or Opens a text file called serverLog.txt
	f, err := os.OpenFile("serverLog.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err, "File for Log did not open")
	}

	//DO NOT CLOSE
	//It will close at end of init function and be worthless
	//defer f.Close()

	//set output of logs to f
	log.SetOutput(f)
}
