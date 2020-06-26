# models

The models directory contains the data structures the web server uses. Some data structures are used to send and receive data from the database. Some are used to send data when executing the GoHtml templates. 

## htmlData.go

Contains 3 data structs.

#### PostPage
Fields
- Username
- Error Message

Usage
- Passed to template when an error occurs

#### YourPostPage
Fields
- Username
- Post Information

Usage
- Passed to template when a user looks at their own posts

#### FindPostPage
Fields
- Username
- UsernameToFind
- Post Information

Usage
- Passed to template when a user searches for another users posts

## kafkaData.go

Contains 1 data struct.

#### UserPageVisits
Fields
- Username
- Path requested

Usage
- Used to send data to Apache Kafka whenever a user does a request

## mongoFile.go

Contains 1 data struct.

#### MongoFile
Fields
- Object ID (object ID in MongoDB for that file)
- File Length
- File Chunk Size
- FileUploadTime
- Filename

Usage
- Used to when receiving a response from the MongoDB GridFS file database.

## posts.go

Contains 4 data structs.

#### Post
Fields
- ID (object ID in MongoDB for that post)
- Caption
- Username
- Filename
- Likes
- DatePosted
- UserLikesIDS (IDS of users who have liked that posts)

Usage
- Used to when receiving a response from the MongoDB post database collection.

#### PostWithEncode
Fields
- ID (object ID in MongoDB for that post)
- Caption
- Username
- Filename
- Likes
- DatePosted
- UserLikesIDS (IDS of users who have liked that posts)
- EncodedFile
- IsLiked
- DeleteButton

Usage
- Used to when sending post data to a gohtml template.

#### PostToDB
Fields
- Caption
- Username
- Filename
- Likes
- DatePosted
- UserLikesIDS (IDS of users who have liked that posts)

Usage
- Used to when sending a post to save to the MongoDB post collection.

#### FSFile (deprecated)
Has been moved to mongoFile.go

Fields
- Object ID (object ID in MongoDB for that file)
- File Length
- File Chunk Size
- FileUploadTime
- Filename

Usage
- Used to when receiving a response from the MongoDB GridFS file database.

## session.go

Contains 1 data struct.

#### Session
Fields
- UserID (object ID in MongoDB for that user)
- LastActivity

Usage
- Used in the ServerSession map that maps a temporary UUID to a user for state.

## user.go

Contains 2 data structs.

#### User
Fields
- ID (object ID in MongoDB for that user)
- Username
- Password (hashed in database)
- PostsIDS (IDS of posts that user has posted)
- LikesIDS (IDS of posts that user has liked)

Usage
- Used to when receiving a response from the MongoDB user database collection.

### UserToDB
- Username
- Password (hashed in database)
- PostsIDS (IDS of posts that user has posted)
- LikesIDS (IDS of posts that user has liked)

Usage
- Used to when sending a user to save to the MongoDB user collection.