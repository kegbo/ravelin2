package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
)

type Data struct {
	WebsiteUrl         string	`json:"websiteUrl"`
	SessionId          string	`json:"sessionId"`
	ResizeFrom         Dimension	`json:"resizeFrom"`
	ResizeTo           Dimension	`json:"resizeTo"`
	CopyAndPaste       map[string]bool `json:"copyAndPaste"` // map[fieldId]true
	FormCompletionTime int  `json:"formCompletionTime"`  // Seconds
}

type Dimension struct {
	Width  int64	`json:"width"`
	Height int64	`json:"height"`
}


type Payload struct {
	EventType  string	`json:"type"`
	Resize ResizePayload	`json:"resize,omitempty"`
	Copyandpaste CopyandpastePayload	`json:"copyandpaste,omitempty"`
	Timetaken TimeTakenPayload	`json:"timetaken,omitempty"`
}

 type ResizePayload struct {
	EventType string `json:"eventType"`
	WebsiteUrl string  `json:"websiteUrl"`
	SessionId string  `json:"sessionId"`
	OldWidth int64  `json:"oldWidth"`
	OldHeight int64  `json:"oldHeight"`
	NewWidth int64  `json:"newWidth"`
	NewHeight int64  `json:"newHeight"`
}

type CopyandpastePayload struct {
	EventType string `json:"eventType"`
	WebsiteUrl string  `json:"websiteUrl"`
	SessionId string  `json:"sessionId"`
	Pasted bool  `json:"pasted"`
	FormId string  `json:formid`
}

type TimeTakenPayload struct {
	EventType string `json:"eventType"`
	WebsiteUrl string  `json:"websiteUrl"`
	SessionId string  `json:"sessionId"`
	Time int  `json:"time"`
	
}
func main() {

	
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit",formHandler)
	http.HandleFunc("/data",Handler)
	fmt.Println("Server now running on localhost:8080")
	fmt.Println(`Try running: curl -X POST -d '{"hello":"test123"}' http://localhost:8080/helloworld`)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type helloWorldRequest struct {
	Hello string `json:"hello"`
}

var d = Data{}

func formHandler(w http.ResponseWriter, r *http.Request) {
	
	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read body"))
		return
	}

	fmt.Println(d)
	w.Write([]byte("Recieved"))
}

func Handler(w http.ResponseWriter, r *http.Request) {


	body, err := ioutil.ReadAll(r.Body)
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read body"))
		return
	}
	var req Payload
	err = json.Unmarshal(body, &req);
	
	if  err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to unmarshal JSON request"))
		return
	}

	
	fmt.Println(req.EventType)
	var ls = req.EventType
	switch ls {
			case "resize":
				var re ResizePayload
				re = req.Resize
				
				d.WebsiteUrl = re.WebsiteUrl
				d.SessionId = re.SessionId
				d.ResizeFrom.width = re.OldWidth
				d.ResizeFrom.height = re.OldHeight
				d.ResizeTo.width = re.NewWidth
				d.ResizeTo.height = re.NewHeight
				fmt.Println(d)
			case "copypaste":
				var re CopyandpastePayload
				re = req.Copyandpaste
				fmt.Println(re.EventType)
				d.WebsiteUrl = re.WebsiteUrl
				d.SessionId = re.SessionId
				var cp = map[string]bool{
					re.FormId :	re.Pasted,
				}
				d.CopyAndPaste = cp
				fmt.Println(d)
			case "timetaken":
				var re TimeTakenPayload
				re = req.Timetaken
				fmt.Println(re.EventType)
				d.WebsiteUrl = re.WebsiteUrl
				d.SessionId = re.SessionId
				d.FormCompletionTime = re.Time
				fmt.Println(d)
			default:
				w.WriteHeader(http.StatusBadRequest)
				return
		}
}



func indexHandler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("../client/index.html")
    t.Execute(w,"Index")
    w.WriteHeader(http.StatusOK)

}


func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read body"))
		return
	}

	req := &helloWorldRequest{}

	if err = json.Unmarshal(body, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to unmarshal JSON request"))
		return
	}

	log.Printf("Request received %+v", req)

	w.WriteHeader(http.StatusOK)
}
