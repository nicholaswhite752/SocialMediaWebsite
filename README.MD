# Social Media Website

[Overview](#Overview)  
[Code Structure of this Repository](#CodeStructure)   
[Features](#Features)  
[Implementation](#Implementation)  
[AWS Architecture](#AWSArch)


## Overview <a name="Overview"></a>

The purpose of this project was to create a full stack web application using the knowledge that I have been learning over the past couple of years.

The end product is a social media website that supports users interactivity. 

Note: This README gives a general overview of the architecture and setup of this repository. Each directory also has a README detailing the specifics for that directory. Each GoLang file is also fully commented with explanations for some of the more complex things.

Another Note: This program will only run on Mac and Linux Debian OS. This is because it uses Confluent with Apache Kafka, which is only supported on Linux and Mac OS as of time of writing this. I used Linux Mint for the creation and testing. 

## Code Structure of this Repository <a name="CodeStructure"></a>

### config
Directory has files for web server configuration.
- HTML Template Setup
- Database Setup
- Log File Setup
- Confluent Kafka Setup

### controllers
Directory has files for web server request completion.
- Sign-Up/Login Functionality
- Posting/Deleting Functionality
- Like/Unlike Functionality
- FindingPosts Functionality
- Default Functionality

### models
Directory has files for custom data structures.

### templates
 Directory has files for HTML page templates, CSS, and JS files. Also shows mock up pictures of the website.
 
### main.go
3 main sections for the main.go file  

Kafka Consumer go func()
- for loop that prints out Kafka Delivery reports to a text file called kafkaDelivery.txt
- Used to check that the Kafka Service is working correctly
- In production, comment this section out, and remove the import from the top of the file

Handler Functions
- These functions route a http path to a function in the controller directory
- It is the main way the web server function
- For example, when a user requests www.examplewebsitename.com/login the web server knows the execute the login function (controllers.Login)

ListenAndServe Function
- This listens for incoming requests on port 80
- Port 80 is used for HTTP requests 

## Features: <a name="Features"></a>
- User is able to sign-up for website, and login
- User stays logged in for ten minutes after lasts user activity, using web cookies
- User is able to upload a picture with a caption
- User is able to like posts
- User is able to delete own posts
- User is able to search for other users, and see that searched user's posts
- User is able to posts that they have liked on one page
- User is able to see own posts on one page
- Homepage shows the top 9 posts from all users
- User can logout
- User activity is sent to an event stream system for analysis
- Event stream system tracks page visits, and amount of user requests to server

# Implementation <a name="Implementations"></a>
### Front End
- Web pages designed using HTML, CSS, and Javascript
- Bootstrap 4 used for major web page parts
- Javascript used to implement some of the page functionality (likes, dislikes, post deletions)

### Back End
- All server requests are handled by a GoLang web server
- MongoDB used as main data store
    - Stored User Information
        - Object ID of user
        - Username
        - Hashed Password
        - Object IDs of posts
        - Object IDs of liked posts
    - Stored Post Information
        - Object ID of post
        - Username of poster
        - Caption of post
        - Filename for post photo
        - Likes Number
        - Date Posted
        - Object ID of users who liked that post
    - Stored File using MongoDB GridFs
        - File Information
            - Object ID
            - Length of File
            - Filename
            - Chunk Size
            - Upload Date
        - File Chunk
            - Actual File Binary

### Event System
- Based on Confluent and Apache Kafka
- When users send request to the server
    - The server fulfills user request
    - The server then sends a message to Apache Kafka with the request path that was filled, and the user that requested it
- Confluent and Apache Kafka is used to track how many times each request type has been requested, and how many requests to the server each user has made
    - More information on this in controllers folder README
    
# AWS Architecture <a name="AWSArch"></a>
The end product was a running basic social media website that utilized 2 AWS servers 

One AWS web server ran the Confluent Service. I was able to get what I needed running on a EC2 t2.micro instance. It was a pain, but here are the steps to the setup.
- After creating an AWS Account
- Create EC2 Instance
    - For Configure Security Group
        - Add TCP rule for 9092 (for sending data from web server)
- SSH into EC2 Instance
- Install Java
    - sudo apt-get update 
    - sudo apt upgrade
    - sudo add-apt-repository ppa:linuxuprising/java 
    - sudo apt install openjdk-11-jre-headless
- Install Confluent Kafka
    - wget -qO - https://packages.confluent.io/deb/5.2/archive.key | sudo apt-key add -
    - sudo add-apt-repository "deb [arch=amd64] https://packages.confluent.io/deb/5.2 stable main"
    - sudo apt-get update && sudo apt-get install confluent-platform-2.12
- Set up advertised listener for Kafka
    - This lets other servers write to this Kafka server
    - sudo nano /etc/kafka/server.properties
    - Uncomment
        - advertised.listeners=PLAINTEXT://"kafkahostname":9092
    - Replace "kafkahostname" with IP of the AWS Server that Kafka is running on
- EC2 t2.micro instances have 1 GB of RAM which is too small for the Confluent Platform to run completely
    - Solution is to run only some of the services, and to run them at reduced usage
- Reduce Zookeeper Service RAM Usage
    - Zookeeper is a service that keeps track of Kafka Topics and Clusters
    - sudo nano /usr/bin/zookeeper-server-start
    - Find code that looks like
    - ``` 
      if [ "x$KAFKA_HEAP_OPTS" = "x" ]; then
          export KAFKA_HEAP_OPTS="-Xmx512M -Xms512M"
      fi
      ``` 
    - Change to 
    - ``` 
      if [ "x$KAFKA_HEAP_OPTS" = "x" ]; then
          export KAFKA_HEAP_OPTS="-Xmx128M -Xms64M"
      fi
      ``` 
    - This reduces Zookeeper usage to between 64MB and 128MB of RAM, which is plenty for a one web server use case
- Reduce Kafka Service RAM Usage 
    - Kafka is a service that the messages from the web server are written to
      
    - sudo nano /usr/bin/kafka-server-start
    - Find code that looks like (may not be exact)
    - ``` 
      if [ "x$KAFKA_HEAP_OPTS" = "x" ]; then
          export KAFKA_HEAP_OPTS="-Xmx1G -Xms1G"
      fi
      ``` 
    - Change to 
    - ``` 
      if [ "x$KAFKA_HEAP_OPTS" = "x" ]; then
          export KAFKA_HEAP_OPTS="-Xmx256M -Xms128M"
      fi
      ``` 
    - This reduces Kafka usage to between 128MB and 256MB of RAM, which is plenty for a one web server use case
- Run Zookeeper Service
    - sudo systemctl start confluent-zookeeper 
- Run Kafka Service
    - sudo systemctl start confluent-kafka
- Run KSQL Service 
    - KSQL is used to do SQL like statements on Kafka Message Streams
    - sudo systemctl start confluent-ksql
    
- Go into KSQL CLI to set up SQL Streams
    - cd /usr/bin
    - ./ksql
        - starts ksql
    - Inside CLI
        - CREATE STREAM USERPAGEVISITS (username VARCHAR, pageVisited VARCHAR) WITH (kafka_topic='userPageVisits', value_format='JSON', KEY='username');
            - Creates a stream that takes input data from a Kafka Topic named userPageVisits
                - This is the topic the web server writes to
            - This new stream has a username column and pageVisited column
        - CREATE TABLE userTotal AS SELECT username, count(*) as visits FROM USERPAGEVISITS GROUP BY username;
            - Creates a table from the USERPAGEVISITS stream that was just created
            - Tracks total amount of users request to server by user
        - CREATE TABLE pageTotal AS SELECT pagevisited, count(*) as visits FROM USERPAGEVISITS GROUP BY pagevisited;
            - Creates a table from the USERPAGEVISITS stream that was just created
            - Tracks total amount of times a certain path was requested
            
- As users use the site you could use the KSQL client to run commands
    - PRINT "topic name here"
        - Print all messages the web server writes to that topic
    - SELECT * FROM "TABLE NAME HERE"
        - Print all information in the table specified

That is it for the Confluent Server. 

The second AWS server that is running is the server running the main web server. To set this up do the following
- Create EC2 t2.micro Instance
    - For Configure Security Group
        - Add TCP rule for 80 (HTTP traffic)
- transfer code folder over to EC2 instance
    - scp -r -i "PEM FILE" "FolderName" "Public Instance ID":
    - without "", and that colon at the end is needed
- SSH into EC2 instance
   - ssh -i "PEM FILE" "Public Instance IP"
- Install Go to EC2 Instance
    - sudo apt update
    - sudo apt upgrade
    - sudo apt install golang-go
- Check go env
    - go env
    - My GOPATH was /home/"user"/go
- Commands to Create Go Environment
    - cd
    - mkdir go 
    - cd go
    - mkdir bin
    - mkdir pkg
    - mkdir src
- Copy folder with code, to go src directory
    - cp -r "path where folder is" /home/"user"/go/
- Start Server
    - 1 way
        - cd go/src/"whatever folder you copied in"
        - go run main.go
        - However you can't exit that terminal or the server will shut down
    - Better Way
        - Make the program a service
        - Follow this tutorial
        - https://github.com/GoesToEleven/golang-web-dev/tree/master/031_aws/02_hands-on/02_solution#persisting-your-application
- After Server Starts Successfully
    - Go to a browser and go to the Public IP of the AWS Instance that is running the Web Server

 