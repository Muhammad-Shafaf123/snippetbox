package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.01 (Project Setup and Creating
	// a Module) so that the import statement looks like this:
	// "{your-module-path}/internal/models". If you can't remember what module path you
	// used, you can find it at the top of the go.mod file.
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.shafaf.net/internal/models"
)

/*
Define an application struct to hold the application-wide dependencies for the
web application. For now we'll only include fields for the two custom loggers, but
we'll add more to it as the build progresses.
*/
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:root@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	/*
		To keep the main() function tidy I've put the code for creating a connection
		pool into the separate openDB() function below. We pass openDB() the DSN
		from the command-line flag.
	*/

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exit.
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Listening on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

	// // swap the route declarations to use the application struct's methods as the handlers functions.
	// mux := http.NewServeMux()

	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// // Register the other application routes as normal.
	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet/view", app.snippetView)
	// mux.HandleFunc("/snippet/create", app.snippetCreate)

	// /*
	// 	Initialize a new http.Server struct. We set the Addr and Handler fields so
	// 	that the server uses the same network address and routes as before, and set
	// 	the ErrorLog field so that the server now uses the custom errorLog logger in
	// 	the event of any problems.
	// */
	// srv := &http.Server{
	// 	Addr:     *addr,
	// 	ErrorLog: errorLog,
	// 	Handler:  app.routes(),
	// }

	// infoLog.Printf("Starting server on %s", *addr)
	// // Call the ListenAndServer() method on our new http.Serve struct.
	// err := srv.ListenAndServe()
	// errorLog.Fatal(err)

}

// The openDB() function wraps sql.open() and returns sql.DB connection pool.
// for a given DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
