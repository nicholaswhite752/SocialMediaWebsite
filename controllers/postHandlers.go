package controllers

import (
	"crypto/sha1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"SocialWebsite/config"
	"SocialWebsite/models"
	"fmt"
	"time"
	"context"
)

//Serves the upload post page
func UploadPageServe(w http.ResponseWriter, r *http.Request) {
	//They should always be logged in at this point tbh
	if alreadyLoggedIn(w, r) {
		//Gets the user
		userTemp := GetUser(w, r)

		//Sends Data for Kafka
		ProduceUserPageRequest(userTemp.Username, "/uploadpost")

		//Creates an object to pass into the HTML template, with username and no error message
		passObj := models.PostPage{
			UsernamePass: userTemp.Username,
			ErrorMsgPass: template.HTML(""),
		}
		//Executes HTML for the page
		err := config.TPL.ExecuteTemplate(w, "uploadPostPage.gohtml",  passObj)
		if err != nil {
			log.Println(err, "Error Executing Template : postHandlers.go : UploadPageServe")
		}
		return
	}


	//Sends Data for Kafka
	ProduceUserPageRequest("NotLoggedIn", "/uploadpost")

	//If the user is not logged in, redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

//Uploads the Post
func UploadPost(w http.ResponseWriter, r *http.Request) {
	//They should always be logged in at this point
	//Because upload page will not execute if not logged in
	//However, if not logged in, redirect to login
	if !alreadyLoggedIn(w, r) {
		//Sends Data for Kafka
		ProduceUserPageRequest("NotLoggedIn", "/submitUploadPost")

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//Gets user
	userTemp := GetUser(w, r)

	//Sends Data for Kafka
	ProduceUserPageRequest(userTemp.Username, "/submitUploadPost")


	//If method is post request (submitting a post)
	if r.Method == http.MethodPost {
		//Check if user can make a post
		//Right now users can make up to 5 posts
		if len(userTemp.PostsIDS) > 5 {
			//Send an error message that user has too many posts
			errorMsg := `<div class="row justify-content-center">
              <div class="col-6 text-center">
                  <!-- Heading Text -->
                  <!-- Has custom CSS to change text size on media size -->
                  <p style="color: red;">Error: Users can make 5 posts. Please delete a post.</p>
              </div>
          </div>`
			//Creates an object to pass to executing HTML, with username and error message
			errorObj := models.PostPage{
				UsernamePass: userTemp.Username,
				ErrorMsgPass: template.HTML(errorMsg),
			}
			//Executes the template
			err := config.TPL.ExecuteTemplate(w, "uploadPostPage.gohtml",  errorObj)
			if err != nil {
				log.Println(err, "Error Executing Template : postHandlers.go : UploadPost : User More Than 5 Posts")
			}

			return
		}

		//The user has been deemed able to make a post

		//Limit File Size to 10MB
		r.Body = http.MaxBytesReader(w, r.Body, 10 * 1024 * 1024)

		//Gets file
		mf, _, err := r.FormFile("imageToServer")
		if err != nil {
			//Means a Failed Read or File Over 10 MB
			log.Println(err, "Failed to Read in FormFile")

			//Send an error message that file was too big
			errorMsg := `<div class="row justify-content-center">
              <div class="col-6 text-center">
                  <!-- Heading Text -->
                  <!-- Has custom CSS to change text size on media size -->
                  <p style="color: red;">Error: File must be 10MB or less</p>
              </div>
          </div>`
			//Creates an object to pass to executing HTML, with username and error message
			errorObj := models.PostPage{
				UsernamePass: userTemp.Username,
				ErrorMsgPass: template.HTML(errorMsg),
			}
			//Executes template
			err := config.TPL.ExecuteTemplate(w, "uploadPostPage.gohtml",  errorObj)
			if err != nil {
				log.Println(err, "Error Executing Template : postHandlers.go : UploadPost : Size More than 10MB")
			}

			return
		}

		//If no error defer the file closing
		defer mf.Close()

		//Gets Caption, limited in HTML to 300 characters
		caption := r.FormValue("captionToServer")

		//Create a buffer to read in file information, make sure file is an image
		//512 bytes because that is how much file header is
		buff := make([]byte, 512)
		//read first 512 bytes of file into buffer
		if _, err = mf.Read(buff); err != nil {
			//if error, log specific error in a file
			log.Println(err, "Error Reading into Buffer for File Info")
			//this function executes a template with a basic error message for user
			postError(w, r, userTemp)
			return
		}

		//Gets type of file
		fileTypeRes := http.DetectContentType(buff)

		//Slice of accepted file types
		xAccepetedFileTypes := []string{"image/jpeg", "image/png", "image/x-citrix-jpeg", "image/x-citrix-png", "image/x-png"}

		//Make a variable to check if file is an image
		ifImage := false

		//Ranges over accepeted file types to see if the file uploaded matches any of them
		for _, v := range xAccepetedFileTypes{
			//If it matches one type, set variable to true and break loop
			if v == fileTypeRes {
				ifImage = true
				break
			}
		}

		//If it is not an image
		//Reload the page with error msg
		if !ifImage{
			//Error message for not being an image
			errorMsg := `<div class="row justify-content-center">
              <div class="col-6 text-center">
                  <!-- Heading Text -->
                  <!-- Has custom CSS to change text size on media size -->
                  <p style="color: red;">Error: File must be a PNG or JPG/JPEG</p>
              </div>
          </div>`
			//Create an object with username and error message
			errorObj := models.PostPage{
				UsernamePass: userTemp.Username,
				ErrorMsgPass: template.HTML(errorMsg),
			}
			err := config.TPL.ExecuteTemplate(w, "uploadPostPage.gohtml",  errorObj)
			if err != nil {
				log.Println(err, "Error Executing Template : postHandlers.go : Not An Image Upload")
			}
			return
		}

		//File is of correct format, size, and user can post

		//Seek to beginning of file
		_, err = mf.Seek(0,0)
		if err != nil {
			log.Println(err, "Error Seeking to file beginning : postHandlers.go")
			postError(w, r, userTemp)
			return
		}

		//Gets bytes from file
		fileBytes, err := ioutil.ReadAll(mf)
		if err != nil {
			log.Println(err, "Error Reading File Bytes : postHandlers.go")
			postError(w, r, userTemp)
			return
		}

		//Seek to start of file again
		_, err = mf.Seek(0,0)
		if err != nil {
			log.Println(err, "Error Seeking to file beginning 2 : postHandlers.go")
			postError(w, r, userTemp)
			return
		}

		//Create a new hash function
		h := sha1.New()

		//Run the file bytes through the hash function
		_, err = io.Copy(h, mf)
		if err != nil {
			log.Println(err, "Error Copying through hash : postHandlers.go")
			postError(w, r, userTemp)
			return
		}

		//Create name for file by summing hash function
		//This creates a unique name for each unique file upload
		//If two users upload the same file, they will have the same name
		fname := fmt.Sprintf("%x", h.Sum(nil))


		//Creates context for DB searches
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		//Gets number of posts that the same filename as the post being uploaded
		numPostWithSameFilename, err := config.PostCol.CountDocuments(ctx,  bson.D{{"filename", fname}})
		if err != nil {
			log.Println(err, "Error Getting Amount of Posts with Same Filename : postHandlers.go")
			postError(w, r, userTemp)
			return
		}


		//If the file is not on the server then upload it to the server
		if numPostWithSameFilename == 0{
			//Seek back to start of file
			_, err = mf.Seek(0,0)
			if err != nil {
				log.Println(err, "Error Seeking to file beginning : postHandlers.go : Inside File not on Server")
				postError(w, r, userTemp)
				return
			}

			//Creates the fileUpload Stream
			uploadStream, err := config.Bucket.OpenUploadStream(
				fname,
			)
			if err != nil {
				log.Println(err,"Error Creating Upload Stream : postHandlers.go")
				postError(w, r, userTemp)
				return
				//os.Exit(1)
			}
			defer uploadStream.Close()

			//Writes the file data to mongodb
			_, err = uploadStream.Write(fileBytes)
			if err != nil {
				log.Println(err,"Error Writing File to MongoDB : postHandlers.go")
				postError(w, r, userTemp)
				return
			}
		}


		//If the file is already on the server because of same hash filename just send post data

		//Create a post object with info about the post
		postSend := models.PostToDB{
			Username:     userTemp.Username,
			Caption:	  caption,
			Filename:     fname,
			Likes:        0,
			DatePosted:   time.Now(),
			UserLikesIDS: nil,
		}
		//Insert new post into mongoDB
		insertedResult, err := config.PostCol.InsertOne(context.TODO(), postSend)
		if err != nil {
			log.Println(err,"Error Writing Post to MongoDB : postHandlers.go")
			postError(w, r, userTemp)
			return
		}

		//Updates Users postIDs for new post
		userTemp.PostsIDS = append(userTemp.PostsIDS, insertedResult.InsertedID.(primitive.ObjectID) )

		//Actually Doesn't return an error
		_ = config.UserCol.FindOneAndUpdate(context.TODO(), bson.D{{"_id", userTemp.ID}},  bson.D{{"$set", bson.D{{"postids", userTemp.PostsIDS}}},
		})

		//Redirect to /uploadPost which is function UploadPageServe
		http.Redirect(w, r, "/uploadpost", http.StatusSeeOther)
		return


	}

	//Redirect to /uploadPost which is function UploadPageServe
	http.Redirect(w, r, "/uploadpost", http.StatusSeeOther)
}

func postError(w http.ResponseWriter, r *http.Request ,userTemp models.User){
	//Basic Error message for users
	errorMsg := `<div class="row justify-content-center">
              <div class="col-6 text-center">
                  <!-- Heading Text -->
                  <!-- Has custom CSS to change text size on media size -->
                  <p style="color: red;">Error: Unexpected. Try Again.</p>
              </div>
          </div>`
	//Create object with username and error
	errorObj := models.PostPage{
		UsernamePass: userTemp.Username,
		ErrorMsgPass: template.HTML(errorMsg),
	}

	//Executes Template
	err := config.TPL.ExecuteTemplate(w, "uploadPostPage.gohtml",  errorObj)
	if err != nil {
		log.Println(err, "Error Executing Template : postHandlers.go : postError Function")
	}
}