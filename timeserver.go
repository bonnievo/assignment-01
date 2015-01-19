
// Copyright 2015 Bonnie Vo
// CSS 490C Tactical Software Engineering
// Professor Morris Bernstein
/*
Go program called timeserver that will serve a web page displaying 
the current time of the day. Default port is 8080 and displays current 
time on localhost:8080/time. (display on time request only)
Other request displays a message 

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
)

// handler - server should generate a page that displays current time
func handler(response http.ResponseWriter, request *http.Request) {
	const layout = "03:04:05 PM"
	t := time.Now()
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
					The time is now <span class="time">%v</span>.
				</p>
			</body>
		</html>`, t.Format(layout))
}





// handler for other request - return status code 404 and display message
func handlerForOtherRequest(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(404)
	fmt.Fprintf(rw, 
		`<html>
			<body>
				<p>
					These are not the URLs you're looking for.
				</p>
			</body>
		</html>`)
}





func main() {
	portNumber := flag.Int("port", 8080, "Set port number")
	versionNumber := flag.Bool("V", false, "Current version number")
	flag.Parse()

	if *versionNumber {
		fmt.Println("The current version of this program is v1.0")
		return
	}

 	http.HandleFunc("/time", handler)
	http.HandleFunc("/", handlerForOtherRequest)

	err := http.ListenAndServe(fmt.Sprintf(":%v", *portNumber) , nil)

 	fmt.Printf("Server fail: %v\n", err)
}
