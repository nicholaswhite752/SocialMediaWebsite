# controllers

The controller directory contains the code files that deals with fulfilling user requests. This README will give an overview of the functionality that each file completes. The code files themselves are also fully commented.

## deletePostHandlers.go 

Contains 2 functions

#### DeletePost  
When a user is logged in, their own posts have a red DELETE button on them. If a user hits that button the DeletePost function runs and does the steps needed to delete that specific post. If it is successful deleted the function returns, and the client side Javascript reloads the user's page. If not the function writeErrorForDelete runs. 

#### writeErrorForDelete  
Generic Error Message that gets sent back to web client. Since deleting a post is done with a Javascript function on the client side, you can't execute a template on success or failure because the return is sent to a Javascript XHR instead of the user's browser. So if something fails, for example a database connection, you send the word "ERROR" back to the Javascript function and then handle it client side.

## findPostsHanlder.go

Contains 2 functions.

#### FindUser  
There is a page for finding posts that appears to a user if they are logged in. It contains a search bar where they can type in a user's username.  

If this page is requested via GET request, then the web page template for this page is executed and served to the user.  

If this page is POST requested from the search, then the web page template for this page is executed with the posts for that username if that username is found. If that username is not found the template for the web page is executed with a generic error message.  

#### findUserError  
This executes if an error occurs or if a user is not found. Executes the web page template for the find posts page that includes an error message.

## handlers.go
This files contains 2 functions for generic use.

#### Index  
The Index function serves the homepage to a user by executing the web page template for the homepage. Separate templates are executed depending if the user is logged in or not logged in.

#### GetUser  
The GetUser function gets the information of the current user that is making a request by reading the cookie value of that request. This website uses cookies to track which users are making requests. This function is used in almost every code file in the controller directory.

## kafkaProducerHandler.go

Contains 1 function.

#### ProduceUserPageRequest  
This function sends the data to an Apache Kafka topic. The other functions in the controller directory call this function to send username, and request name to a topic.

## likedPostPage.go

Contains 1 function.

#### LikedPostPage  
Serves the user their liked post page when requested by executing the liked post page template.

## likeHandlers.go

Contains 3 functions.

#### LikePost  
When a user is logged in and and clicks on the clear heart icon of a post to like that post this function is called. It is called by a Javascript XHR. If it is successful liked the function returns, and the client side Javascript makes the heart icon for that post red. This keeps the page from doing a full reload, because normally the heart icon being clear/red is done on page load. If not the function writeErrorForResponse runs.

 #### UnlikePost  
 When a user is logged in and and clicks on the red heart icon of a post to unlike that post this function is called. It is called by a Javascript XHR. If it is successful unliked the function returns, and the client side Javascript makes the heart icon for that post clear. This keeps the page from doing a full reload, because normally the heart icon being clear/red is done on page load. If not the function writeErrorForResponse runs.
 
 #### writeErrorForResponse
 Returns "ERROR" to the Javascript XHR to let the client side Javascript function know that the post like/unlike was unsuccessful and to take correct action.
 
 ## loginHandlers.go
 
 Contains 2 important variables and 8 functions.
 
 #### ServerSession Map  
 Maps a UUID to a value to a value of type Session. A value of type session contains a username, and a time of last activity. This is used to keep track of active user sessions and cookies.
 
 #### sessionLength Int  
 Amount of time a cookie lasts in a users browser. This is 600 seconds for this file. 
 
 #### Login  
 This function serves the login page to the user by executing the login page template. It redirects users to the homepage if they are already logged in.
 
 #### LoginPost  
 There is a form on the login page for username and password. This function is called if a user submits that form to log in. On a login success, a cookie is created for the user, session added to the ServerSession map, and it redirects to the homepage. On a failure, it reloads the login page with an error message of what failed.
 
 #### Signup  
 If there is a GET request to the signup page, the function serves the signup page by executing the signup page template.  
 
 There is a form on the signup page for username and password. If that form is submitted in a POST request then the function tries to create an account for that user. If it is successful the user is logged in, a cookie is created for the user, session added to the ServerSession map, and it redirects to the homepage. On a failure, it reloads the login page with an error message of what failed.
 
 #### alreadyLoggedIn 
 This function checks to see if the user making a request is logged in or not. It does this by checking if the user has a cookie called "session". If there is not cookie then it returns false as the user is not logged in.
 
 The function then checks the cookies value which is a UUID value. If that UUID value from the cookie is in the ServerSession map then it updates the time of last activity to the current time, rewrites the cookie duration to the sessionLength duration, and returns true. If it is not in the ServerSession map it returns false.
 
 This function is called by almost every function in the controller directory. One reason is because the site behaves differently for users logged in then for users not logged in. Another reason is because this function updates the time of last activity and the cookie in the browser. This means for this case the cookie will expire after 10 minutes (600 seconds) since the last activity, and not just 10 minutes after creation.
 
 #### Logout
 This function logs the user out. It deletes the session cookie from the user's browser and deletes the user's session from the ServerSession map. It redirects to the cleanSessions function (below) when done.
 
 #### cleanSessions
 The session map needs to be cleaned in some cases. For example, lets say someone leaves the site without logging out. After 10 minutes the browser cookie expires, but the session is still in the ServerSession map on the server. So this function deletes all entries in the ServerSession map that are older than the session length. For this small scale site, I had it run every time someone logged out of the website. In a production environment you would run this on a set time scale interval. 
 
 #### loginError
 This function serves the login page to the user with an error message when an unexpected error occurs. An error such as bad username or password would be handled by the main login function. However a general error that the user doesn't need the source of, such as a faulty database connection, is served through this function.
 
  #### signupError
  This function serves the signup page to the user with an error message when an unexpected error occurs. An error such as bad username or password would be handled by the main signup function. However a general error that the user doesn't need the source of, such as a faulty database connection, is served through this function.
 
 ## postHandlers.go
 
 Contains 3 functions.
 
 #### UploadPageServe
 This function serves the upload post page by executing the upload post page template.
 
 #### UploadPost
 On the upload post page there is a form for uploading a photo with a caption. When that form is submitted this function is called. If the post is upload successfully, then the upload post page is reload with no error message. If the post is not uploaded successfully the upload post page is reloaded with an error message.
 
 This function deals with some of the error messages. Users could have 5 image posts that were less than 10MB per post. If the user was trying to post more than 5 times, or if they were trying to upload a non-image file, or if they were trying to upload a file more than 10MB, then this function would serve the upload post page again with an error message of what was wrong. 
 
 #### postError
 This function deals with general errors that the user does not need to know. For example, if there was a problem reading a file, then this function would serve the upload post page with a general error message telling them to retry.
 
 ## yourPostHandlers.go
 
 Contains 1 function.
 
 #### YourPosts
 This function serves the web page that shows a user their own posts by executing the your posts template and serving it to the user.
