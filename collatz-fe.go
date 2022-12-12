package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	//"github.com/go-echarts/go-echarts/v2/charts"
	//"github.com/go-echarts/go-echarts/v2/opts"
)

type resultscollatz struct {
	Value        int   `json:"value"`
	List_results []int `json:"list_results"`
}

type JsonResponse struct {
	Type    string         `json:"type"`
	Data    resultscollatz `json:"data"`
	Message string         `json:"message"`
}

var initnumber int64
var html string
var intSlice []int
var xslice []string

var allSlice []int

func main() {

	fs := http.FileServer(http.Dir("assets"))
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/graph", httpserver)
	mux.HandleFunc("/", httpserver_home)
	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}

func httpserver_home(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("www/index.html"))
	tpl.Execute(w, nil)

}

func httpserver(w http.ResponseWriter, r *http.Request) {

	var tpl = template.Must(template.ParseFiles("www/graph-header.html"))
	tpl.Execute(w, nil)

	allSlice = nil

	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := u.Query()
	nn, err := strconv.ParseInt(params.Get("nhosts"), 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	initnumber = int64(nn)

	response, err := http.Get("http://collatz-be:8081/collatz/" + params.Get("nhosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var abc JsonResponse

	json.Unmarshal(responseData, &abc)

	//json.NewDecoder(response.Body).Decode(&abc)

	fmt.Fprintf(w, "<ol> ")
	for i, p := range abc.Data.List_results {
		fmt.Fprintln(w, "<li>"+strconv.Itoa(i+1)+":"+strconv.Itoa(p)+"</li>")
	}
	fmt.Fprintf(w, "</ol> ")

	//BuildGraph(w)

	//initnumber2 := int64(nn) - 1
	//for i := initnumber2; i > 0; i-- {
	//	compute((i))

	//}

	//BuildGraphLim0(w)

	tpl = template.Must(template.ParseFiles("www/graph-footer.html"))
	tpl.Execute(w, nil)
}
