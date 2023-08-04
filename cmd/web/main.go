package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

/*
Define an application struct to hold the application-wide dependencies for the
web application. For now we'll only include fields for the two custom loggers, but
we'll add more to it as the build progresses.
*/
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of the application struct, containing the
	// dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// swap the route declarations to use the application struct's methods as the handlers functions.
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the other application routes as normal.
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	/*
		Initialize a new http.Server struct. We set the Addr and Handler fields so
		that the server uses the same network address and routes as before, and set
		the ErrorLog field so that the server now uses the custom errorLog logger in
		the event of any problems.
	*/
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	// Call the ListenAndServer() method on our new http.Serve struct.
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

}
