package main

import (
	"flag"
	"log"
	"net/http"
)

// Define a string constant containing the HTML for the webpage. This consists of a h1
// header tag, and some JavaScript which fetches the JSON from our Get /v1/healthcheck
// endpoint and writes it to inside the div element.
const html = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta http-equiv="Content-Type" content="text/html;charset=UTF-8">
</head>
<body>
	<h1>Simple CORS</h1>
	<div id="ouput"></div>
	<script>
		document.addEventListener("DOMContentLoaded", function(){
			fetch("http://localhost:4000/v1/healthcheck").then(
				function (response) {
					response.text().then(function (text){
						document.getElementById("ouput").innerHTML = text;
					})
				},

				function (err) {
					document.getElementById("ouput").innerHTML = err;
				}
			);
		});
	</script>
</body>
</html>
`

func main() {
	// Make the server address configurable at runtime via a command-line flag.
	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()

	log.Printf("starting server on %s", *addr)

	// Start a HTTP server listening on the given address, which responds to all
	// requests with the webapp HTTP above.
	err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))

	log.Fatal(err)
}
