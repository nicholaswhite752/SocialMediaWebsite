package config

import (
	"context"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// database
var DB *mongo.Database
var UserCol *mongo.Collection
var PostCol *mongo.Collection
var FsFilesCol *mongo.Collection
var DBFiles *mongo.Database

var Bucket *gridfs.Bucket

func init() {
	// get a mongo sessions
	// Set client options

	//HAD TO SET CONNECTION TIMEOUT TO NEVER
	//Would not do in production
	//Dev environment has low speeds
	clientOptions := options.Client().ApplyURI("mongodb+srv://<specificDatabaseUsername>:<specificDatabasePassword>@<DatabaseConnectionString>").SetConnectTimeout(0)


	//Create context for mongodb
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	//Needs to be log Fatal
	//If Mongo does not start server can't run
	if err != nil {
		log.Fatal(err, "Didn't connect to MongoDB")
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err, "Mongo Connection test failed")
	}

	//Connects to database
	DB = client.Database("websiteDatabase")
	//Connects to users collection
	UserCol = DB.Collection("users")
	//Connects to posts collection
	PostCol = DB.Collection("posts")

	//Connects to myfiles database
	DBFiles = client.Database("myfiles")
	//Connects to fs.files collection
	FsFilesCol = DBFiles.Collection("fs.files")

	//Creates a bucket in the DBFiles database
	//Used for Mongodb gridfs commands
	Bucket, _ = gridfs.NewBucket(
		DBFiles,
	)


	log.Println("Connected to MongoDB!")

}
