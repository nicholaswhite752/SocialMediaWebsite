# config directory

This directory contains GoLang files that are ran on server start. They all use the init() to create connections and accomplish other preparatory tasks before the server starts running.

## alogSetup.go

This file creates a file called serverLog.txt in the same directory that main.go is located. The file also sets the log functions print output to that file. This way log messages are printed to the file instead of a terminal. 

The reason this file is name alogSetup is because GoLang's default for running files is in alphabetical order. Since the other config Go files have log statements in them, this file needs to run first. That way if a problem occurs a config/setup file, the error message will print to the log file.

## db.go

This file sets up the connection to the MongoDB database. Connections are created to each collection in the MongoDB database and the MongoDB GridFs database that was used for this website. These connections are used globally by the controller files to create, read, update, and delete (CRUD) records in the database.

In a working solution the username, password, and connection string need to be updated to your specific needs.

## kafkaSetup.go

This file sets up the connection to Apache Kafka. The controller files use this connection to write messages to an Apache Kafka topic. 

This file also creates a text file called kafkaDelivery.txt. This file can be used in connection with the go func() that is at the top at main.go. If the go func() is not commented out in main.go, the Kafka message delivery reports will print to the text file. If not, the file will remain empty.    

In the architecture describe in the root README, there is an AWS server running Confluent with Apache Kafka and another AWS server running this code. This file would need to updated with the IP address of the server that is fully setup and running Confluent. 

## tpl.go

This file parses the gohtml templates in the templates folder in the root directory. Those templates are placed into an object, and that object is executed in the controller files to serve web pages to users that request them. 
  