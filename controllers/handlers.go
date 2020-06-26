package controllers

import (
	"SocialWebsite/config"
	"SocialWebsite/models"
	"bytes"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"context"
	"go.mongodb.org/mongo-driver/bson"
	//"log"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {

	//Context for MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//Sets options to sort by highest likes and limit to 9 posts
	findOptions := options.Find()
	// Sort by likes field descending
	findOptions.SetSort(bson.D{{"likes", -1}})
	//Set limit to 9 posts
	findOptions.SetLimit(9)

	//Creates structure to store the posts
	var postSlice []models.Post

	//Searches Database
	postCursor, err := config.PostCol.Find(ctx,  bson.D{{}}, findOptions)
	if err != nil {
		log.Println(err, "Error reading from PostCol : handlers.go : index")
		err := config.TPL.ExecuteTemplate(w, "index.gohtml", nil)
		if err != nil {
			log.Println(err, "Error Executing Template : handlers.go : index Finding Post Error")
		}
		return
	}

	//Puts found post into a slice of posts
	for postCursor.Next(context.TODO()) {
		// create a value to decode post
		var postTemp models.Post
		err := postCursor.Decode(&postTemp)
		if err != nil {
			//If post can't be decoded skip post and continue with the others
			log.Println(err, "Error Getting a Post : handlers.go : index")
			continue
		}
		postSlice = append(postSlice, postTemp)
	}

	//Different logic for logged in and not logged in users
	if alreadyLoggedIn(w, r) {
		//If logged in get current user
		userTemp := GetUser(w, r)

		//Sends Data for Kafka
		ProduceUserPageRequest(userTemp.Username, "/")

		var postSends []models.PostWithEncode

		for _, v := range postSlice{
			var buf bytes.Buffer
			_, err := config.Bucket.DownloadToStreamByName(v.Filename, &buf)
			if err != nil {
				log.Println(err, "Error Getting File for a Post : handlers.go : logged in")
				//A continue here will allow other posts to be sent, but this post will not be
				continue
			}

			//turns buffer into a slice of bytes
			content := buf.Bytes()

			//encodes to base64
			encoded := base64.StdEncoding.EncodeToString(content)

			//Turns date posted into something usable
			tempDatePosted := v.DatePosted.Format("01-02-2006")

			isLikedString := ""
			//Ranges through the users liked posts array
			//If the id of the post is in the user likes, then give the attribute liked
			//It is passed into the template as a custom css class
			for _, v2 := range userTemp.LikesIDS{
				if v2 == v.ID{
					isLikedString = "liked"
				}
			}

			//Tells if have ability to delete
			deleteTemp := false
			//If the username of the post, matches the username of the current user
			//Set deleteTemp to true
			//This will show a button that allows that user to delete the post
			if userTemp.Username == v.Username{
				deleteTemp = true
			}

			//Object of a single post
			tempPostEncodes := models.PostWithEncode{
				ID:           v.ID,
				Caption:      v.Caption,
				Username:     v.Username,
				Filename:     v.Filename,
				Likes:        v.Likes,
				DatePosted:   tempDatePosted,
				UserLikesIDS: v.UserLikesIDS,
				EncodedFile:  encoded,
				IsLiked:	  isLikedString,
				DeleteButton: deleteTemp,
			}

			//Slice of posts that will be sent to be executed by HTML
			postSends = append(postSends, tempPostEncodes)

		}

		//Object for username and post data
		tempHTMLData := models.YourPostPage{
			UsernamePass:   	userTemp.Username,
			PostsWithEncode: 	postSends,
		}

		err := config.TPL.ExecuteTemplate(w, "indexLogged.gohtml",  tempHTMLData)
		if err != nil {
			log.Println(err, "Error Executing Template : handlers.go : index logged in")
		}
		return
	}

	//If not logged in logic

	//Sends Data for Kafka
	ProduceUserPageRequest("NotLoggedIn", "/")

	var postSends []models.PostWithEncode

	for _, v := range postSlice{
		var buf bytes.Buffer
		_, err := config.Bucket.DownloadToStreamByName(v.Filename, &buf)
		if err != nil {
			log.Println(err, "Error Getting File for a Post : yourPostHandlers.go : not logged in")
			//A continue here will allow other posts to be sent, but this post will not be
			continue
		}

		//turns buffer into a slice of bytes
		content := buf.Bytes()

		//encodes to base64
		encoded := base64.StdEncoding.EncodeToString(content)

		//Turns date posted into something usable
		tempDatePosted := v.DatePosted.Format("01-02-2006")

		//isLiked CSS class will always be empty, because user is not logged in
		isLikedString := ""
		//No delete button because user is logged in
		deleteTemp := false

		//Object of a single post
		tempPostEncodes := models.PostWithEncode{
			ID:           v.ID,
			Caption:      v.Caption,
			Username:     v.Username,
			Filename:     v.Filename,
			Likes:        v.Likes,
			DatePosted:   tempDatePosted,
			UserLikesIDS: v.UserLikesIDS,
			EncodedFile:  encoded,
			IsLiked:	  isLikedString,
			DeleteButton: deleteTemp,
		}

		//Slice of posts that will be sent to be executed by HTML
		postSends = append(postSends, tempPostEncodes)

	}

	//Object for username and post data
	tempHTMLData := models.YourPostPage{
		UsernamePass:   	"",
		PostsWithEncode: 	postSends,
	}

	err = config.TPL.ExecuteTemplate(w, "index.gohtml", tempHTMLData)
	if err != nil {
		log.Println(err, "Error Executing Template : handlers.go : index not logged in")
	}

}

//Only gets called if log in passes, so it will have a cookie
//Because why would you need to get the current user, without checking if they are logged in
func GetUser(w http.ResponseWriter, r *http.Request) models.User {
	// get cookie
	c, err := r.Cookie("session")

	//Gets DB ObjectID of user
	//Error for c being nil, but that shouldn't happen
	userObjectID := ServerSessions[c.Value].UserID

	//Creates Context
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//For return
	var userDB models.User

	//Find User in DB
	err = config.UserCol.FindOne(ctx, bson.D{{"_id", userObjectID}}).Decode(&userDB)
	if err != nil {
		//If this happens, just log the error
		//The user that will be returned will be NULL
		//And it might make the program wonky for that response, but it will still function correctly
		log.Println("Error Retrieving User : ", userObjectID , ": handlers.go : GetUser Func")
	}

	//Returns that users info
	return userDB
}