package controllers

import (
	"SocialWebsite/config"
	"SocialWebsite/models"
	"context"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var ServerSessions = map[string]models.Session{} // session ID, session

const sessionLength int = 600

//Just Serves the Login Page
func Login(w http.ResponseWriter, r *http.Request) {
	//redirects to homepage if already logged in
	//Should not run unless you login and request /login route
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//Sends Data for Kafka
	ProduceUserPageRequest("NotLoggedIn", "/login")

	err := config.TPL.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		log.Println(err, "Error Executing Template : loginHandlers.go : Login")
	}
}


func LoginPost(w http.ResponseWriter, r *http.Request) {
	//redirects to homepage if already logged in
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//Sends Data for Kafka
	ProduceUserPageRequest("NotLoggedIn", "/loginPost")

	//If Method is POST
	//This should always happen, because the only place that calls this method is a POST
	if r.Method == http.MethodPost {
		//gets username and password
		//usernameToServer is the name of that form value in the HTML
		un := r.FormValue("usernameToServer")
		//passwordToServer is the name of that form value in the HTML
		p := r.FormValue("passwordToServer")


		//Context for DB searches
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		//Searches for username, returns results with a BSON cursor
		userCursor, err := config.UserCol.Find(ctx,  bson.D{{"username", un}})
		if err != nil {
			log.Println(err, "Error Pulling UserCol : loginHandlers.go : LoginPost")
			loginError(w, r)
			return
		}

		//Puts find results into a struct of Users, should only be 0 or 1 length
		var results []models.User
		for userCursor.Next(context.TODO()) {
			// create a value into which the single document can be decoded
			var userTemp models.User
			//Decode search results into an object of User
			err := userCursor.Decode(&userTemp)
			if err != nil {
				log.Println(err, "Error Decoding User Results : loginHandlers.go : LoginPost")
				loginError(w, r)
				return
			}
			//Append to slice of users
			results = append(results, userTemp)
		}

		//Length should be either 0 or 1, as that is how many results should be returned
		if len(results) > 1{
			//THIS SHOULD NEVER HAPPEN
			log.Println("Something is wrong, more than one acct returned with same username")
			loginError(w, r)
			return
		}else if len(results) == 0 {
			//Username does not Exist
			errorMsg := `<div class="row justify-content-center">
              <div class="col-6 text-center">
                  <!-- Heading Text -->
                  <!-- Has custom CSS to change text size on media size -->
                  <p style="color: red;">Error: Invalid Username Or Password</p>
              </div>
          </div>`
			//Execute Template with error message
			err = config.TPL.ExecuteTemplate(w, "login.gohtml",  template.HTML(errorMsg))
			if err != nil {
				log.Println(err, "Error Executing Template : loginHandlers.go : Login POST Invalid Username")
			}
			return
		} else{
			//Check for correct password

			err = bcrypt.CompareHashAndPassword([]byte(results[0].Password), []byte(p))
			if err != nil {
				//Passwords don't match
				//Error message for not matching passwords
				errorMsg := `<div class="row justify-content-center">
              <div class="col-6 text-center">
                  <!-- Heading Text -->
                  <!-- Has custom CSS to change text size on media size -->
                  <p style="color: red;">Error: Invalid Username Or Password</p>
              </div>
          </div>`
				err = config.TPL.ExecuteTemplate(w, "login.gohtml", template.HTML(errorMsg))
				if err != nil {
					log.Println(err, "Error Executing Template : loginHandlers.go : Login POST Passwords don't match")
				}
				return
			}

			//Passwords do match
			//Create session with a cookie
			sID, _ := uuid.NewV4()
			c := &http.Cookie{
				Name:  "session",
				Value: sID.String(),
			}
			//Set max age to 10 minutes
			//This means that it will delete from the browser after 10 minutes
			c.MaxAge = sessionLength
			http.SetCookie(w, c)

			//Puts session into Server Side Session Map
			ServerSessions[c.Value] = models.Session{
				UserID:       results[0].ID,
				LastActivity: time.Now(),
			}

			//redirects logged in user to homepage
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		}
	}

	//If for some reason there is a unlogged in GET request at this page
	//Serve them the homepage
	err := config.TPL.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		log.Println(err, "Error Executing Template : loginHandlers.go : LoginPost Bottom")
	}

}

func Signup(w http.ResponseWriter, r *http.Request) {
	//redirects to homepage if already logged in
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//Sends Data for Kafka
	ProduceUserPageRequest("NotLoggedIn", "/signup")

	//If method is POST
	//GET Method is for page request
	if r.Method == http.MethodPost {
		//gets username and password
		//usernameToServer is the name of that form value in the HTML
		un := r.FormValue("usernameToServer")
		//passwordToServer is the name of that form value in the HTML
		p := r.FormValue("passwordToServer")


		//Context for MondoDB Searches
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		//Gets number of posts that the same filename as the post being uploaded
		numUserWithSameUsername, err := config.UserCol.CountDocuments(ctx,  bson.D{{"username", un}})
		if err != nil {
			log.Println(err, "Error Getting Amount of Users with Same Username : loginHandlers.go : Signup POST")
			signupError(w, r)
			return
		}


		//Length should be either 0 or 1, as that is how many results should be returned
		if numUserWithSameUsername > 1{
			//THIS SHOULD NEVER HAPPEN
			log.Println("Something is wrong, more than one acct returned with same username")
			signupError(w, r)
			return
		}else if numUserWithSameUsername != 0 {
			//Username is already taken
			//Error Message
			errorMsg := `<div class="row justify-content-center">
              <div class="col-6 text-center">
                  <!-- Heading Text -->
                  <!-- Has custom CSS to change text size on media size -->
                  <p style="color: red;">Error: Username Taken</p>
              </div>
          </div>`
			//Execute template
			err = config.TPL.ExecuteTemplate(w, "signup.gohtml",  template.HTML(errorMsg))
			if err != nil {
				log.Println(err, "Error Executing Template : loginHandlers.go : Signup POST Username Taken")
			}
			return
		} else{
			//New user is being created

			//Create a hashed password using bcrypt
			passwordToDB, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
			if err != nil {
				log.Println(err, "Error Creating Hashed Password : loginHandlers.go : Signup POST")
				signupError(w, r)
				return
			}

			//Create a user that will be pushed to the mongodb database
			userToPush := models.UserToDB{
				Username: un,
				Password: string(passwordToDB),
				PostsIDS: nil,
				LikesIDS: nil,
			}

			insertResult, err := config.UserCol.InsertOne(context.TODO(), userToPush)
			if err != nil {
				log.Println(err, "Error Inserting New User : loginHandlers.go : Signup POST")
				signupError(w, r)
				return
			}

			// create session with a cookie
			sID, _ := uuid.NewV4()
			c := &http.Cookie{
				Name:  "session",
				Value: sID.String(),
			}
			//Set max age to 10 minutes
			//This means that it will delete from the browser after 10 minutes
			c.MaxAge = sessionLength
			http.SetCookie(w, c)

			//GETS OBJECT ID for session mapping
			//fmt.Println("OBJECTID: ", insertResult.InsertedID.(primitive.ObjectID))
			//Puts session into Server Side Session Map
			ServerSessions[c.Value] = models.Session{
				UserID:		  insertResult.InsertedID.(primitive.ObjectID),
				LastActivity: time.Now(),
			}

			//redirects logged in user to homepage
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		}
	}

	//Normal GET of this page
	err := config.TPL.ExecuteTemplate(w, "signup.gohtml", nil)
	if err != nil {
		log.Println(err, "Error Executing Template : loginHandlers.go : Signup Bottom")
	}

}

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	//gets cookie called session
	c, err := r.Cookie("session")
	if err != nil {
		//If no cookie, then not logged in
		return false
	}

	//gets server side session
	s, ok := ServerSessions[c.Value]
	//if value is already in server sessions, update last time of activity
	if ok {
		s.LastActivity = time.Now()
		ServerSessions[c.Value] = s
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		return ok
	} else{
		//this should hopefully never run
		//But if it does for some reason delete the cookie from the user
		c.MaxAge = -1
		c.Value = ""
		http.SetCookie(w, c)
		return ok
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	//If the user tries to logout without being logged in
	//Redirect to index
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//The user is logged in so they have cookie, but if they don't somehow
	//Redirect to homepage
	c, err := r.Cookie("session")
	if err != nil {
		//If no cookie, then not logged in
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// delete the session from the server of that user
	delete(ServerSessions, c.Value)

	// remove the cookie from users browser
	//Do this by setting MaxAge to -1 and writing cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	//Redirect to /cleansess
	//This will clean the server session table to erase any sessions older than 10 min
	//This function then redirects to index
	http.Redirect(w, r, "/cleansess", http.StatusSeeOther)

}

//Cleans Sessions when /cleansess is called
func CleanSessions(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("BEFORE CLEAN") // for demonstration purposes
	//showSessions()              // for demonstration purposes
	for k, v := range ServerSessions {
		if time.Now().Sub(v.LastActivity) > (time.Second * 600) {
			delete(ServerSessions, k)
		}
	}

	//fmt.Println("AFTER CLEAN") // for demonstration purposes
	//showSessions()             // for demonstration purposes
	//Redirects to index after cleaning the session table
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//prints sessions to terminal for testing
//func showSessions() {
//	fmt.Println("********")
//	for k, v := range ServerSessions {
//		fmt.Println(k, v.UserID)
//	}
//	fmt.Println("")
//}

//Basic Error For something went wrong
func loginError(w http.ResponseWriter, r *http.Request){
	//Basic Error message for users
	errorMsg := `<div class="row justify-content-center">
              <div class="col-6 text-center">
                  <!-- Heading Text -->
                  <!-- Has custom CSS to change text size on media size -->
                  <p style="color: red;">Error: Unexpected Error while Logging In. Try Again.</p>
              </div>
          </div>`

	err := config.TPL.ExecuteTemplate(w, "login.gohtml",  template.HTML(errorMsg))
	if err != nil {
		log.Println(err, "Error Executing Template : postHandlers.go : loginError Function")
	}

}

//Basic Error For something went wrong
func signupError(w http.ResponseWriter, r *http.Request){
	//Basic Error message for users
	errorMsg := `<div class="row justify-content-center">
              <div class="col-6 text-center">
                  <!-- Heading Text -->
                  <!-- Has custom CSS to change text size on media size -->
                  <p style="color: red;">Error: Unexpected Error while Signing Up. Try Again.</p>
              </div>
          </div>`

	err := config.TPL.ExecuteTemplate(w, "login.gohtml",  template.HTML(errorMsg))
	if err != nil {
		log.Println(err, "Error Executing Template : postHandlers.go : signupError Function")
	}

}