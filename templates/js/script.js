function like(idName){
    //Gets Element of the Like Button
    var wholeCard = document.getElementById(idName);
    var cardBod = wholeCard.children[1];
    var likesPar = cardBod.children[1];
    var likesNumber = likesPar.children[0];
    var likesIcon = likesPar.children[1];


    //If the post is currently liked
    if(likesIcon.classList.contains("liked")) {
        //Creating a new HTTP Request
        var xhr = new XMLHttpRequest();
        //Setting request to Post on path /unlikePost and async
        xhr.open("POST", "/unlikePost", true);

        //When the state of xhr changes
        xhr.onreadystatechange = function() {
            //If the response is normal
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                //Server sends back the word true if everything worked correctly
                //Server sends back the word notLogged if not logged in
                //Server sends back Error (else statement) if an error occurred
                var isTrueSet = (xhr.responseText.toString().trim() === "true");
                var isNotLogged = (xhr.responseText.toString().trim() === "notLogged");
                if (isTrueSet){
                    //If true toggle the red background on the like button to default
                    likesIcon.classList.toggle("liked");
                    //Decreases number the web client sees
                    likesNumber.textContent = Number(likesNumber.textContent) - 1;
                }
                else if (isNotLogged){
                    //Alerts that user is not logged in
                    alert("You are not logged in. Redirecting.");
                    //Redirects to login page
                    setTimeout(() => {  window.location.replace("/login"); }, 200);
                }
                else {
                    //Alerts that some error occurred
                    alert("Failed to Unlike Post. Please Try Again.");
                    //Reloads page for user to retry
                    setTimeout(() => {  location.reload(); }, 200);
                }
            }
        };

        //Sends XHR with ID of post to be unliked as the request body
        xhr.send(idName);
    }
    else{
        //Post started with not being liked
        //Creates HTTP request
        var xhr = new XMLHttpRequest();
        //Setting request to Post on path /unlikePost and async
        xhr.open("POST", "/likePost", true);

        //When the state of xhr changes
        xhr.onreadystatechange = function() {
            //If the response is normal
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                //Server sends back the word true if everything worked correctly
                //Server sends back the word notLogged if not logged in
                //Server sends back Error (else statement) if an error occurred
                var isTrueSet = (xhr.responseText.toString().trim() === "true");
                var isNotLogged = (xhr.responseText.toString().trim() === "notLogged");
                if (isTrueSet){
                    //If true toggle the red background on the like button to being on
                    likesIcon.classList.toggle("liked");
                    //Increases number the web client sees
                    likesNumber.textContent = Number(likesNumber.textContent) + 1;
                }
                else if (isNotLogged){
                    //Alerts that user is not logged in
                    alert("You are not logged in. Redirecting.");
                    //Redirects to login page
                    setTimeout(() => {  window.location.replace("/login"); }, 200);
                }
                else {
                    //Alerts that some error occurred
                    alert("Failed to Like Post. Please Try Again.");
                    //Reloads page for user to retry
                    setTimeout(() => {  location.reload(); }, 200);
                }
            }
        };

        //Sends XHR with ID of post to be Liked as the request body
        xhr.send(idName);
    }
}

function deletePost(idName){
    //Creates HTTP request
    var xhr = new XMLHttpRequest();
    //Setting request to Post on path /unlikePost and async
    xhr.open("POST", "/deletePost", true);

    //When state of XHR changes
    xhr.onreadystatechange = function() {
        //If the response is normal
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            //Server sends back the word true if everything worked correctly
            //Server sends back the word notLogged if not logged in
            //Server sends back Error (else statement) if an error occurred
            var isTrueSet = (xhr.responseText.toString().trim() === "true");
            var isNotLogged = (xhr.responseText.toString().trim() === "notLogged");
            if (isTrueSet){
                //If true is sent from server, then reloads page
                setTimeout(() => {  location.reload(); }, 200);
            }
            else if (isNotLogged){
                //If not logged in is sent from server, alerts
                alert("You are not logged in. Redirecting.");
                //Then redirects to login
                setTimeout(() => {  window.location.replace("/login"); }, 200);
            }
            else {
                //An error occurred while deleting the post, alerts
                alert("Failed to Delete Post. Please Try Again.");
                //Reloads the page to let user try again
                setTimeout(() => {  location.reload(); }, 200);
            }
        }
    };

    //Sends XHR with ID of post to be Deleted as the request body
    xhr.send(idName);

}