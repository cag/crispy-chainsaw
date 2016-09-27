package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/negroni"
)

type Logger struct{}

var (
	methodcol   = color.New(color.Bold).SprintFunc()
	pathcol     = color.New(color.Italic).SprintFunc()
	timecol     = pathcol
	infocol     = color.New(color.Bold, color.FgCyan).SprintFunc()
	successcol  = color.New(color.Bold, color.FgGreen).SprintFunc()
	redirectcol = color.New(color.Bold, color.FgYellow).SprintFunc()
	clierrcol   = color.New(color.Bold, color.FgRed).SprintFunc()
	serverrcol  = color.New(color.Bold, color.FgMagenta).SprintFunc()
)

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	res := rw.(negroni.ResponseWriter)
	status := res.Status()
	var statstr string
	switch {
	case 200 <= status && status < 300:
		statstr = successcol(status)
	case 300 <= status && status < 400:
		statstr = redirectcol(status)
	case 400 <= status && status < 500:
		statstr = clierrcol(status)
	case 500 <= status:
		statstr = serverrcol(status)
	default:
		statstr = infocol(status)
	}
	log.Printf("%v %v %vB → %v %v %vB",
		methodcol(r.Method), pathcol(r.URL.Path), r.ContentLength,
		statstr, timecol(time.Since(start)), res.Size())
}

func main() {
	portcolor := color.New(color.Bold, color.FgBlue).SprintFunc()

	var port int
	flag.IntVar(&port, "p", 8080, "port to serve on")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/whyy", func(w http.ResponseWriter, req *http.Request) {
		// fmt.Fprintf(w, "Welcome to the home page!")
		panic("I'M PANICKING")
	})

	n := negroni.New(
		&Logger{},
		negroni.NewRecovery(),
		negroni.NewStatic(http.Dir("public")))

	n.UseHandler(mux)

	log.Printf("Serving on HTTP port: %v\n", portcolor(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), n))
}
