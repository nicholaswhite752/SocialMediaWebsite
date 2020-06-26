# templates

The templates directory contains gohtml files that are executed when a user requests a web page. Data is passed into gohtml templates by using the {{.}} operation.

For example if you passed a struct into a gohtml template 

GoCode
```
type User struct {
    Username string
    FavColor string
}

u1 := User{
    Username: "test",
    FavColor: "blue"
    }    

Template.ExecuteTemplate(w, "indexExample.gohtml", u1)

```

Gohtml template
```html
<html>
    <body>
        {{.Username}}
        {{.FavColor}}
    </body>
</html> 
```

The output would be a blank page with 

```
test
blue
```

The js directory contains one js file that as functions to like posts, unlike posts, and delete your own posts.

The styles directory has a custom CSS file.


## findPosts.gohtml

This template is executed when a user is searching for posts from another user.

Find Post Page Before Searching a Username:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/FindPostPagePreSearch.PNG)

Find Post Page After Searching a Username:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/FindPostPagePostSearch.PNG)

Find Post Page After Searching a Username (Zoomed Out):  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/FindPostPagePostSearchZoomedOut.PNG)

## index.gohtml

This template is executed as the main homepage of the site. This one is executed for users who are not logged in.

Index Page At Full Screen Size:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/indexNotLoggedInFullScreen.PNG)

Index Page At Smaller/Mobile Screen Size:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/indexNotLoggedInSmallScreen.PNG)

## indexLogged.gohtml

This template is executed as the main homepage of the site. This one is executed for users who are logged in.

Index Page when Logged In:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/indexWhenLoggedIn.PNG)

## login.gohtml

This template is executed when a user goes to the login page.

Login Page:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/loginPage.PNG)

Login Page with an Error:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/LoginPageWhenErrorPassedToTemplate.PNG)

## signup.gohtml

This template is executed when a user goes to the signup page.

SignUp Page:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/signupPage.PNG)

SignUp Page with an Error:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/signupPageWhenErrorPassed.PNG)

## uploadPostPage.gohtml

This template is executed when a user goes to the upload a post page.

Upload Post Page:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/uploadPostPage.PNG)

## yourLikedPosts.gohtml

This template is executed when a user goes to their liked posts page.

Liked Post Page:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/likedPostsPage.PNG)

## yourPosts.gohtml

This template is executed when a user goes to their posts page.

Your Posts Page:  
![alt-text](https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/yourPostsPage.PNG)
