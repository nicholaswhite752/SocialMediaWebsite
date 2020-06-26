package controllers

import (
	"SocialWebsite/config"
	"SocialWebsite/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"strings"
)

//Like Post is called through a JavaScript XML Request
//So we can't redirect, but we can return a string and have JS Redirect
func LikePost(w http.ResponseWriter, r *http.Request) {
	//If trying to like but not logged in
	//Takes you to log in page
	if !alreadyLoggedIn(w, r) {
		//Sends Data for Kafka
		ProduceUserPageRequest("NotLoggedIn", "/likePost")


		//Responds to Request
		_, err := fmt.Fprintln(w, "notLogged")
		if err != nil {
			log.Println(err, "Error Writing Response : likeHandlers.go : LikePost")
		}
		return
	}

	//Reads the body of the post
	y, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err, "Error Reading Body : likeHandlers.go : LikePost")
		writeErrorForResponse(w, r)
		return
	}

	//The string from that response looks like
	//ObjectID("xxxxxxxxxHEXxxxxxxxxx")
	//Need to make an actual object of ObjectID type
	//Delimit by " to get the HexString
	temp := string(y)
	xTemp := strings.Split(temp, "\"")

	//Create an Object ID from a Hex String
	obj1, err := primitive.ObjectIDFromHex(xTemp[1])
	if err != nil {
		log.Println(err, "Error Creating Object ID from Hex : likeHandlers.go : LikePost")
		writeErrorForResponse(w, r)
		return
	}

	//Context for mongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Gets Post from ObjID
	var tempPost models.Post
	//Decode returns an error is nothing is found
	err = config.PostCol.FindOne(ctx,  bson.D{{"_id", obj1}}).Decode(&tempPost)
	if err != nil {
		log.Println(err, "Error Getting Post : likeHandlers.go : LikePost")
		writeErrorForResponse(w, r)
		return
	}

	//Gets Current User
	userTemp := GetUser(w, r)

	//Sends Data for Kafka
	ProduceUserPageRequest(userTemp.Username, "/likePost")

	//Updates the users liked posts
	userTemp.LikesIDS = append(userTemp.LikesIDS, tempPost.ID)

	//Updates User in Database for new like
	//Does not return an error
	_ = config.UserCol.FindOneAndUpdate(context.TODO(), bson.D{{"_id", userTemp.ID}},  bson.D{{"$set", bson.D{{"likesids", userTemp.LikesIDS}}},})

	//Updates PostLikes and WhoLikes it
	tempPost.Likes = tempPost.Likes + 1
	tempPost.UserLikesIDS = append(tempPost.UserLikesIDS, userTemp.ID)

	//Creates BSON format to pass into update
	//Because updating for multiple fields
	updateBsonD := bson.D{
				{"userlikesids", tempPost.UserLikesIDS},
				{"likes", tempPost.Likes},
				}

	//Updates Post Collection with new like amount and users who liked it
	_ = config.PostCol.FindOneAndUpdate(context.TODO(), bson.D{{"_id", tempPost.ID}},  bson.D{{"$set", updateBsonD},})

	_, err = fmt.Fprintln(w, "true")
	if err != nil{
		log.Println(err, "Error Writing Response for TRUE : likeHandlers.go : LikePost")
	}

	return
}


func UnlikePost(w http.ResponseWriter, r *http.Request) {
	//If trying to like but not logged in
	//Takes you to log in page
	if !alreadyLoggedIn(w, r) {
		//Sends Data for Kafka
		ProduceUserPageRequest("NotLoggedIn", "/unlikePost")

		//Responds to Request
		_, err := fmt.Fprintln(w, "notLogged")
		if err != nil {
			log.Println(err, "Error Writing Response : likeHandlers.go : UnlikePost")
		}
		return
	}

	//Reads the body of the post
	y, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err, "Error Creating Object ID from Hex : likeHandlers.go : UnlikePost")
		writeErrorForResponse(w, r)
		return
	}

	//The string from that response looks like
	//ObjectID("xxxxxxxxxHEXxxxxxxxxx")
	//Need to make an actual object of ObjectID type
	//Delimit by " to get the HexString
	temp := string(y)
	xTemp := strings.Split(temp, "\"")

	//Create an Object ID from a Hex String
	obj1, err := primitive.ObjectIDFromHex(xTemp[1])
	if err != nil {
		log.Println(err, "Error Creating Object ID from Hex : likeHandlers.go : UnlikePost")
		writeErrorForResponse(w, r)
		return
	}

	//Context
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Gets Post from ObjID
	var tempPost models.Post
	err = config.PostCol.FindOne(ctx,  bson.D{{"_id", obj1}}).Decode(&tempPost)
	if err != nil {
		log.Println(err, "Error Getting Post : likeHandlers.go : UnlikePost")
		writeErrorForResponse(w, r)
		return
	}

	//Gets Current User
	userTemp := GetUser(w, r)

	//Sends Data for Kafka
	ProduceUserPageRequest(userTemp.Username, "/likePost")

	var tempLikesIDS []primitive.ObjectID
	//Removes the post from the liked post slice
	for _, v := range userTemp.LikesIDS{
		if v != tempPost.ID{
			tempLikesIDS = append(tempLikesIDS, v )
			//I commented this out so if for some reason, there are 2+ records of the post being liked,
			//Both will be removed
			//break
		}
	}
	userTemp.LikesIDS = tempLikesIDS

	//Updates User in Database for new like
	_ = config.UserCol.FindOneAndUpdate(context.TODO(), bson.D{{"_id", userTemp.ID}},  bson.D{{"$set", bson.D{{"likesids", userTemp.LikesIDS}}},})


	//Updates PostLikes and WhoLikes it
	//In this iteration we will only remove one like, even if for some reason a bug occurred and one user liked it twice
	//In future iterations PostLikes would be dictated by how many userIDs it has in the UserLikesIDS slice on that post
	tempPost.Likes = tempPost.Likes - 1

	//Removes the user id from the list of users who like the post
	var tempUserLikesIDS []primitive.ObjectID
	for _, v := range tempPost.UserLikesIDS{
		if v != userTemp.ID{
			tempUserLikesIDS = append(tempUserLikesIDS, v)
			//I commented this out so if for some reason, there are 2+ records of the post being liked,
			//Both will be removed
			//break
		}
	}
	tempPost.UserLikesIDS = tempUserLikesIDS

	//Updating multiple fields, so do a BSON object outside of MongoDB statement for readability
	updateBsonD := bson.D{
		{"userlikesids", tempPost.UserLikesIDS},
		{"likes", tempPost.Likes},
	}

	//Updates Post Collection with new like amount and users who liked it
	_ = config.PostCol.FindOneAndUpdate(context.TODO(), bson.D{{"_id", tempPost.ID}},  bson.D{{"$set", updateBsonD},})

	_, err = fmt.Fprintln(w, "true")
	if err != nil{
		log.Println(err, "Error Writing Response for TRUE : likeHandlers.go : UnlikePost")
	}

	return

}

func writeErrorForResponse(w http.ResponseWriter, r *http.Request){
	//Write ERROR to response
	_, err := fmt.Fprintln(w, "ERROR")
	if err != nil {
		log.Println(err, "Error Writing Response : likeHandlers.go : writeError Function")
	}
}