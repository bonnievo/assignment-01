
// Copyright 2015 Bonnie Vo
// CSS 490C Tactical Software Engineering
// Professor Morris Bernstein
/*
Go program called timeserver that will serve a web page displaying 
a personalized time server. Personalized time server will support
a home page, login page, a logout page, and a time page. 
A home page should display a greeting message with the user's name.
If no user's name exisit in the user's cookies then the user should
be prompted to a login page where it should display a message for 
the user to input their name. A logout page should clear any exisiting
name in the user's cookies. A time page should display the current time
and the utc with the user's name (if cookie is available or not empty).
Default port is 8080 and displays current 
time on localhost:8080/time.
Personalized time server will also support the following URLs:
http://host:port/, http://host:port/index.html, 
http://host:port/login?name=name, http://host:port/logout

This program takes an optional command line argument flag: --port port_number
This program takes an additional flag: --V which writes the version number to 
standard output and terminates. 

Program displays error message and terminates if chosen port is already in use
*/

package main

import (
	"fmt"
	"net/http"
	"time"
	"flag"
	"sync"
)

// map to store the cookie
var cookieMap = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string] string)}


// login cookie - creates a cookie and returns a cookie
func LoginCookie(username string) http.Cookie {
	cookieValue := username // + ":" + codify.SHA(username+strconv.Itoa(rand.Intn(100000000)))
	expire := time.Now().AddDate(0, 0, 1) // expires in one day
	return http.Cookie{Name: "name", Value: cookieValue, Expires: expire, HttpOnly: true}

}

// login handler - display login form and request user's name
func loginHandler(response http.ResponseWriter, request *http.Request) {
	//request.ParseForm()
	name := request.FormValue("name")

	/*if (request.Method == "POST") {
		fmt.Println("in post method")
	} else */
	if name == "" { 
		loginPageHandler(response, request)
	} else {
		c := LoginCookie(name)
		http.SetCookie(response, &c)
		cookieMap.Lock()
		cookieMap.m["name"] = name
		cookieMap.Unlock()
		http.Redirect(response, request, "/", http.StatusFound)
	}
}

// greetings page - home page. display greeting message if there is 
// a cookie with a corresponding name. otherwise display the login form
func greetingsPage(response http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("name")
	if err != nil || cookie.Value == "" {
		http.Redirect(response, request, "/login", http.StatusFound)
		return
	}
	
	fmt.Fprintf(response,
		`<html>
			<body>
				<p>
					Greetings, %v
				</p>
			</body>
		</html>`, cookie.Value)
}

// login page handler message - display the login form 
func loginPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response,
		`<html>
			<body>
				<p>
					<form action="login">
				  		What is your name, Earthling?
			  			<input type="text" name="name" size="50">
			   			<input type="submit">
					</form>
				</p>
			</body>
		</html>`)
}

// login page again message - result of a POST request
// display modified login form
func loginPageAgain(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response,
		`<html>
			<body>
				<p>
					C'mon I need a name
					<form action="login">
				  		What is your name, Earthling?
				  			<input type="text" name="name" size="50">
			  			<input type="submit">
					</form>
				</p>
			</body>
		</html>`)	
}


// logout handler - cookie is cleared and the message "goodbye" is 
// displayed for 10 seconds. login form should be displayed after
func logoutHandler(response http.ResponseWriter, request *http.Request) {
	cookie := LoginCookie("")
	http.SetCookie(response, &cookie)
	delete(cookieMap.m, "name")
	fmt.Fprintf(response,
		`<html>
		<head>
		<META http-equiv="refresh" content="10;URL=/">
		<body>
		<p>Good-bye.</p>
		</body>
		</html>`)
}

// handler - server should generate a page that displays current time
// and current UTC. If the user is logged in, then server should 
// also display name of user
func timeHandler(response http.ResponseWriter, request *http.Request) {
	const layout = "03:04:05 PM"
	const utcLayout = "03:04:05"
	t := time.Now()
	cookie, err := request.Cookie("name")

	if err != nil || cookie.Value == "" {
		fmt.Fprintf(response,
		 	`<html>
				<head>
					<style>
						p {font-size: xx-large}
						span.time {color: red}
					</style>
				</head>
				<body>
					<p>
						The time is now <span class="time">%v</span> (%v UTC).
					</p>
				</body>
			</html>`, t.Format(layout), t.UTC().Format(utcLayout))
	} else {
		fmt.Fprintf(response,
	 	`<html>
			<head>
				<style>
					p {font-size: xx-large}
					span.time {color: red}
				</style>
			</head>
			<body>
				<p>
					The time is now <span class="time">%v</span> (%v UTC), %s.
				</p>
			</body>
		</html>`, t.Format(layout),t.UTC().Format(utcLayout), cookie.Value)
	}
}


func main() {
	portNumber := flag.Int("port", 8080, "Set port number")
	versionNumber := flag.Bool("V", false, "Current version number")
	flag.Parse()

	if *versionNumber {
		fmt.Println("The current version of this program is v2.0")
		return
	}

 	http.HandleFunc("/time", timeHandler)

 	http.HandleFunc("/", greetingsPage)
	http.HandleFunc("/index.html", greetingsPage)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%v", *portNumber) , nil)

 	fmt.Printf("Server fail: %v\n", err)
}