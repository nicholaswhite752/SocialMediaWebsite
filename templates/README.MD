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
(https://github.com/nicholaswhite752/SocialMediaWebsite/blob/master/PageImages/FindPostPagePreSearch.PNG)

Image of searched page

## index.gohtml

This template is executed as the main homepage of the site. This one is executed for users who are not logged in.

Image of page

## indexLogged.gohtml

This template is executed as the main homepage of the site. This one is executed for users who are logged in.

Image of Page

## login.gohtml

This template is executed when a user goes to the login page.

Image of Page

## signup.gohtml

This template is executed when a user goes to the signup page.

Image of Page

## uploadPostPage.gohtml

This template is executed when a user goes to the upload a post page.

Image of Page

## yourLikedPosts.gohtml

This template is executed when a user goes to their liked posts page.

Image of Page

## yourPosts.gohtml

This template is executed when a user goes to their posts page.

Image of Page
