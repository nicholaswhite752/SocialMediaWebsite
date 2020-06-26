package controllers

import (
	"SocialWebsite/config"
	"SocialWebsite/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func DeletePost(w http.ResponseWriter, r *http.Request) {
	//Go to login if not logged in
	if !alreadyLoggedIn(w, r) {
		//Sends Data for Kafka
		ProduceUserPageRequest("NotLoggedIn", "/deletePost")

		_, err := fmt.Fprintln(w, "notLogged")
		if err != nil {
			log.Println(err, "Error Writing Response : deletePostHandlers.go : DeletePost")
		}
		return
	}

	//Reads the body of the post
	y, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err, "Error Reading Body : deletePostHandlers.go : DeletePost")
		writeErrorForResponse(w, r)
		return
	}

	//Need to make an actual object of ObjectID type
	//Delimit by " to get the HexString
	temp := string(y)
	xTemp := strings.Split(temp, "\"")

	//Create an Object ID from a Hex String
	obj1, err := primitive.ObjectIDFromHex(xTemp[1])
	if err != nil {
		log.Println(err, "Error Creating Object ID from Hex : deletePostHandlers.go : DeletePost")
		writeErrorForDelete(w, r)
		return
	}

	//Context
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Gets Post from ObjID
	var deletingPost models.Post
	err = config.PostCol.FindOne(ctx,  bson.D{{"_id", obj1}}).Decode(&deletingPost)
	if err != nil {
		log.Println(err, "Error Getting Post : deletePostHandlers.go : DeletePost")
		writeErrorForDelete(w, r)
		return
	}

	//Gets Current User
	userTemp := GetUser(w, r)

	//Sends Data for Kafka
	ProduceUserPageRequest(userTemp.Username, "/deletePost")

	//Number of posts that reference that file
	numPostWithSameFilename, err := config.PostCol.CountDocuments(ctx,  bson.D{{"filename", deletingPost.Filename}})
	if err != nil {
		log.Println(err, "Error Getting Amount of Posts with Same Filename : deletePostHandlers.go : DeletePost")
		writeErrorForDelete(w, r)
		return
	}


	//Update all the users who liked that posts "LikesIDS"
	for _,v := range deletingPost.UserLikesIDS{
		var removeLikeUser models.User
		//Gets the user with that ID
		err = config.UserCol.FindOne(ctx,  bson.D{{"_id", v}}).Decode(&removeLikeUser)
		if err != nil {
			log.Println(err, "Error Getting User of for Unliking Post : deletePostHandlers.go : DeletePost")
			writeErrorForDelete(w, r)
			return
		}

		//Remove the post id from the user likeIds
		var tempUserLikesIDS []primitive.ObjectID
		for _, v2 := range removeLikeUser.LikesIDS{
			if v2 != deletingPost.ID{
				tempUserLikesIDS = append(tempUserLikesIDS, v2 )
				//I commented this out so if for some reason, there are 2+ records of the post,
				//Both will be removed
				//break
			}
		}

		removeLikeUser.LikesIDS = tempUserLikesIDS

		//Find user and update likes
		_ = config.UserCol.FindOneAndUpdate(context.TODO(), bson.D{{"_id", removeLikeUser.ID}},  bson.D{{"$set", bson.D{{"likesids", removeLikeUser.LikesIDS}}},})

	}


	//Removing the file from mongodb if needed
	if numPostWithSameFilename == 1{
		//Need to get fileID from mongodb
		//Have to use Cursor because Bucket does not have a FindOne
		bucketCursor, err := config.Bucket.Find(bson.D{{"filename", deletingPost.Filename}} )
		if err != nil {
			log.Println(err, "Error Finding File to Delete : deletePostHandlers.go : DeletePost")
			writeErrorForDelete(w, r)
			return
		}

		var foundFiles []models.MongoFile

		for bucketCursor.Next(context.TODO()) {
			// create a value into which the single document can be decoded
			var fileTemp models.MongoFile
			err := bucketCursor.Decode(&fileTemp)
			if err != nil {
				log.Println(err, "Error Decoding File to Delete : deletePostHandlers.go : DeletePost")
				writeErrorForDelete(w, r)
				return
			}
			foundFiles = append(foundFiles, fileTemp)
		}

		//Will delete file if one with filename is found
		if len(foundFiles) == 1{
			err = config.Bucket.Delete(foundFiles[0].ID)
			if err != nil {
				log.Println(err, "Error Deleting File to Delete : deletePostHandlers.go : DeletePost")
				writeErrorForDelete(w, r)
				return
			}
		}
	}

	//Update Users Post IDs
	var tempPostIDS []primitive.ObjectID
	for _, v := range userTemp.PostsIDS{
		if v != deletingPost.ID{
			tempPostIDS = append(tempPostIDS, v)
			//I commented this out so if for some reason, there are 2+ records of the post,
			//Both will be removed
			//break
		}
	}
	userTemp.PostsIDS = tempPostIDS

	//Updates User in Database for removing a post
	_ = config.UserCol.FindOneAndUpdate(context.TODO(), bson.D{{"_id", userTemp.ID}},  bson.D{{"$set", bson.D{{"postids", userTemp.PostsIDS}}},})


	//Lastly removing the post record
	_, err = config.PostCol.DeleteOne(ctx,  bson.D{{"_id", deletingPost.ID}})
	if err != nil {
		log.Println(err, "Error Deleting Post : deletePostHandlers.go : DeletePost")
		writeErrorForDelete(w, r)
		return
	}

	//Writes true to Javascript function to show success
	_, err = fmt.Fprintln(w, "true")
	if err != nil{
		log.Println(err, "Error Writing Response for TRUE : deletePostHandlers.go : DeletePost")
	}


	return

}

func writeErrorForDelete(w http.ResponseWriter, r *http.Request){
	//Write ERROR to response
	_, err := fmt.Fprintln(w, "ERROR")
	if err != nil {
		log.Println(err, "Error Writing Response : deletePostHandlers.go : deleteError Function")
	}
}