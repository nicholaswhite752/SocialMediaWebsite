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

func FindUser(w http.ResponseWriter, r *http.Request) {
	//Go to login if not logged in
	if !alreadyLoggedIn(w, r) {
		//Sends Data for Kafka
		ProduceUserPageRequest("NotLoggedIn", "/findPosts")

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//Gets current user
	userTemp := GetUser(w, r)

	//Sends Data for Kafka
	ProduceUserPageRequest(userTemp.Username, "/findPosts")

	if r.Method == http.MethodPost {
		//Gets inputted username
		unToFind := r.FormValue("userToServer")

		//Create context for DB search
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		//USE THE CURSOR because you will need that found user's info anyway
		//Searches for username, returns results with a BSON cursor
		userCursor, err := config.UserCol.Find(ctx,  bson.D{{"username", unToFind}})
		if err != nil {
			log.Println(err, "Error Pulling Searched User : findPostsHandler.go : FindUser")
			findUserError(w, r, userTemp)
			return
		}

		//Puts find results into a struct of Users
		var results []models.User
		for userCursor.Next(context.TODO()) {
			// create an object to be decoded into
			var userTemp models.User
			err := userCursor.Decode(&userTemp)
			if err != nil {
				log.Println(err, "Error Decoding User : findPostsHandler.go : FindUser")
				continue
			}
			//If user is found and decodes, then append it to user slice
			results = append(results, userTemp)
		}


		//There should be only one or zero results, if more something is wrong but case is still handled

		//No results found for that username
		if len(results) == 0{
			//Passes current user's name and that a user was not found for the query
			passedData := models.FindPostPage{
				UsernamePass:   userTemp.Username,
				UserToFind:      "User Not Found: " + unToFind,
				PostsWithEncode: nil,
			}
			//Execute template with error
			err := config.TPL.ExecuteTemplate(w, "findPosts.gohtml",  passedData)
			if err != nil {
				log.Println(err, "Error Executing Template : findPostsHandler.go : Len Results 0")
			}

			return
		} else if len(results) == 1 {
			userFound := results[0]

			//Case if you search for yourself
			//Redirect to yourPosts
			if userFound.ID == userTemp.ID{
				http.Redirect(w, r, "/yourPosts", http.StatusSeeOther)
				return
			}

			var postSlice []models.Post

			//Gets posts by found user, puts into post slice
			for _, v := range userFound.PostsIDS {
				var tempPost models.Post
				err := config.PostCol.FindOne(ctx, bson.D{{"_id", v}}).Decode(&tempPost)
				if err != nil {
					log.Println(err, "Error Getting a Post : findPostsHandlers.go")
					continue
				}
				//If post is found and decodes correctly then append to found Post Slice
				postSlice = append(postSlice, tempPost)
			}

			//Puts posts in order of Date
			//Posts that are more recent appear first
			sort.Slice(postSlice, func(i, j int) bool {
				return postSlice[i].DatePosted.After(postSlice[j].DatePosted)
			})

			var postSends []models.PostWithEncode
			//fmt.Println(postSlice)
			for _, v := range postSlice {
				var buf bytes.Buffer
				_, err := config.Bucket.DownloadToStreamByName(v.Filename, &buf)
				if err != nil {
					log.Println(err, "Error Getting File for a Post : findPostsHandlers.go")
					//A continue here will allow other posts to be sent, but this post will not be
					continue
				}

				//turns buffer into a slice of bytes
				content := buf.Bytes()

				//encodes to base64
				encoded := base64.StdEncoding.EncodeToString(content)

				//Turns date posted into something usable
				tempDatePosted := v.DatePosted.Format("01-02-2006")

				//Tells if current user likes the post
				isLikedString := ""
				//Ranges through the users liked posts array
				//If the id of the post is in the user likes, then give the attribute liked
				//It is passed into the template as a custom css class
				for _, v2 := range userTemp.LikesIDS {
					if v2 == v.ID {
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
					IsLiked:      isLikedString,
					DeleteButton: deleteTemp,
				}

				//Slice of posts that will be sent to be executed by HTML
				postSends = append(postSends, tempPostEncodes)

			}//Ends for loop

			//Object for username and post data
			tempHTMLData := models.FindPostPage{
				UsernamePass:    userTemp.Username,
				UserToFind:      userFound.Username,
				PostsWithEncode: postSends,
			}

			//Executes the template, with the searched user's posts
			err := config.TPL.ExecuteTemplate(w, "findPosts.gohtml", tempHTMLData)
			if err != nil {
				log.Println(err, "Error Executing Template : findPostsHandler.go : Len Results 1")
			}
			return

		}else {
			//If this happens something is just really wrong with the database
			passedData := models.FindPostPage{
				UsernamePass:   userTemp.Username,
				UserToFind:      "User Not Found: " + "SOMETHING REALLY WRONG WITH DB",
				PostsWithEncode: nil,
			}

			//Executes the template, with error
			err := config.TPL.ExecuteTemplate(w, "findPosts.gohtml",  passedData)
			if err != nil {
				log.Println(err, "Error Executing Template : findPostsHandler.go : Len Results Not 0 OR 1")
			}
			return
		}


	}

	//For case of just arriving on page
	passedData := models.FindPostPage{
		UsernamePass:   userTemp.Username,
		UserToFind:      "",
		PostsWithEncode: nil,
	}

	err := config.TPL.ExecuteTemplate(w, "findPosts.gohtml",  passedData)
	if err != nil {
		log.Println(err, "Error Executing Template : findPostsHandler.go : Bottom for GET")
	}

}

//Basic Error For something went wrong
func findUserError(w http.ResponseWriter, r *http.Request, userTemp models.User){
	//For case of just arriving on page
	passedData := models.FindPostPage{
		UsernamePass:   userTemp.Username,
		UserToFind:      "ERROR Finding User. Please Try Again",
		PostsWithEncode: nil,
	}

	err := config.TPL.ExecuteTemplate(w, "findPosts.gohtml",  passedData)
	if err != nil {
		log.Println(err, "Error Executing Template : findPostsHandler.go : findUserError")
	}

}