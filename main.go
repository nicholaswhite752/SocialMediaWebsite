package main

import (
	"SocialWebsite/config"
	"SocialWebsite/controllers"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"net/http"
	"time"
)


func main() {
	//Spins up go func for listening to producer events
	//This should work
	go func() {
	//	//Infinite for loop to have this run forever and check events as they come in
	//	//Good for small scale
	//	//Bad for large scale
		for {
			for e := range config.Producer.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						//Prints to file if delivery failed
						_, err := fmt.Fprintf(config.FileForKafka, "Delivery failed: %v\n", ev.TopicPartition)
						if err != nil {
							log.Println(err, "Failed to write Kafka Delivery Failed to file")
						}
					} else {
						//Prints to file if delivery success
						_, err := fmt.Fprintf(config.FileForKafka, "Delivered message to %v\n", ev.TopicPartition)
						if err != nil {
							log.Println(err, "Failed to write Kafka Delivery Success to file")
						}
					}
				}
			}
			time.Sleep(time.Second)
		}
	}()


	//Servers Main Pages
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/loginPost", controllers.LoginPost)
	http.HandleFunc("/signup", controllers.Signup)
	http.HandleFunc("/logout", controllers.Logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	//Servers Upload/Posts
	http.HandleFunc("/uploadpost", controllers.UploadPageServe)
	http.HandleFunc("/submitUploadPost", controllers.UploadPost)
	http.HandleFunc("/yourPosts", controllers.YourPosts)

	//Liking/Unliking Posts
	http.HandleFunc("/likePost", controllers.LikePost)
	http.HandleFunc("/unlikePost", controllers.UnlikePost)
	http.HandleFunc("/likedPosts", controllers.LikedPostPage)

	//DeletingPosts
	http.HandleFunc("/deletePost", controllers.DeletePost)

	//Searching for other users/posts
	http.HandleFunc("/findPosts", controllers.FindUser)

	//Session Cleaning Testing
	http.HandleFunc("/cleansess", controllers.CleanSessions)

	//Serves CSS used by HTML when executing Templates
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./templates/styles/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./templates/js/"))))

	http.ListenAndServe(":80", nil)

}

