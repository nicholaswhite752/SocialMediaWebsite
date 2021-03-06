package controllers

import (
	"SocialWebsite/config"
	"SocialWebsite/models"
	"bytes"
	"context"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"sort"
	"time"
)

func LikedPostPage(w http.ResponseWriter, r *http.Request) {
	//Go to login if not logged in
	if !alreadyLoggedIn(w, r) {
		//Sends Data for Kafka
		ProduceUserPageRequest("NotLoggedIn", "/likedPosts")

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//Gets Current User
	userTemp := GetUser(w, r)

	//Sends Data for Kafka
	ProduceUserPageRequest(userTemp.Username, "/likedPosts")

	//Context for MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var postSlice []models.Post

	//Gets posts that user likes and puts them in post slice
	for _, v := range userTemp.LikesIDS {
		var tempPost models.Post
		err := config.PostCol.FindOne(ctx,  bson.D{{"_id", v}}).Decode(&tempPost)
		if err != nil {
			log.Println(err, "Error Getting a Post : likedPostPage.go")
			continue
		}
		//If post is found and decodes correctly then append to found Post Slice
		postSlice = append(postSlice, tempPost)
	}

	//sorts by date posted
	sort.Slice(postSlice, func(i, j int) bool {
		return postSlice[i].DatePosted.After(postSlice[j].DatePosted)
	})

	var postSends []models.PostWithEncode
	//
	for _, v := range postSlice{
		var buf bytes.Buffer
		_, err := config.Bucket.DownloadToStreamByName(v.Filename, &buf)
		if err != nil {
			log.Println(err, "Error Getting File for a Post :  likedPostPage.go")
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

	//Executes the template, with the current user's posts
	err := config.TPL.ExecuteTemplate(w, "yourLikedPosts.gohtml", tempHTMLData)
	if err != nil {
		log.Println(err, "Error Executing Template : likedPostPage.go")
	}


}
