package main

import(
	"encoding/json"
	"net/http"
	"log"
	"fmt"
	"embed"
	"html/template"
	"github.com/gorilla/mux"
	"os"
	"time"
	"path/filepath"
)
//go:embed assets/*
var assetsfs embed.FS
type PageData map[string]interface{}

type JsonResp struct {
Code    int         `json:"code"`
Payload interface{} `json:"payload"`
Message string      `json:"message"`
}

func fileExists(path string) bool{

	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists

		return true

	} else if os.IsNotExist(err) {
		return false

	} else {
		return false

	}

}

func startHTTPServer(){


r:= mux.NewRouter()

if(fileExists(assets_dir)){
log.Println("Assets found on disk")
r.PathPrefix(fmt.Sprintf("/%s/",assets_dir)).Handler(http.StripPrefix(fmt.Sprintf("/%s/",assets_dir), http.FileServer(http.Dir(assets_dir))))
}else{
log.Println("Using embedded assets")
r.PathPrefix("/assets/").Handler(http.FileServer(http.FS(assetsfs)))
}


r.HandleFunc("/", httpDefaultHandler)
	r.HandleFunc("/status", httpStatusHandler)
http.Handle("/", r)

page_port := conf.GetString("page_port")

srv := &http.Server{
Handler: r,
Addr:    fmt.Sprintf(":%s", page_port),
// Good practice: enforce timeouts for servers you create!
WriteTimeout: 15 * time.Second,
ReadTimeout:  15 * time.Second,
}

log.Println("Listening on port", page_port)
log.Fatal(srv.ListenAndServe())


}

func loadTmpl(dest string) (*template.Template){

tmpl := template.New("default")
//tmpl, _ = tmpl.New("header").ParseFS(assetsfs,"assets/tmpl/header.html")
//tmpl, _ = tmpl.New("footer").ParseFS(assetsfs,"assets/tmpl/footer.html")

var err error

tmpl, err = tmpl.New(filepath.Base(dest)).Funcs(template.FuncMap{
"getTimestamp": func() int64 {
return time.Now().Unix()
},
}).ParseFS(assetsfs,dest)

if(err!=nil){

log.Println(err)
}
//log.Printf("%#v\n",tmpl.DefinedTemplates())

return tmpl
}
func httpStatusHandler(w http.ResponseWriter, r *http.Request) {
	 respond(w,200,"OK",grill.Info)
}

func httpDefaultHandler(w http.ResponseWriter, r *http.Request) {

tmpl := loadTmpl("assets/tmpl/index.html")

data := PageData{}

data["Title"] = "Home"

if(r.Header.Get("X-Requested-With") == "xmlhttprequest" || r.Header.Get("Hx-Request")=="true") {

respond(noCache(w),200,"OK",data)
return

}


tmpl.Execute(w, data)


}

func respond(w http.ResponseWriter, code int, message string, payload interface{}) {

resp := JsonResp{
Code:    code,
Payload: payload,
Message: message,
}

var jsonData []byte
jsonData, err := json.Marshal(resp)

if err != nil {
log.Println(err)
}

w.Header().Set("Content-type", "application/json")

fmt.Fprintln(w, string(jsonData))

}




func noCache(w http.ResponseWriter) (http.ResponseWriter){

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
"Expires":         epoch,
"Vary":            "accept",
"Cache-Control":   "no-cache, private, max-age=0",
"Pragma":          "no-cache",
"X-Accel-Expires": "0",
}

// Set our NoCache headers
for k, v := range noCacheHeaders {
w.Header().Set(k, v)
}

return w

}