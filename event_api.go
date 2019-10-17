package main

import (
  "fmt"
  "time"
  "net/http"
  "log"
  "encoding/json"
  "io/ioutil"
  "strconv"

  "github.com/gorilla/mux"
)

type event struct {
  Id int `json:id`
  Title string `json:title`
  Desc  string `json:description`
}

var events = []event {
 { Id: 1, Title: "int to go lang", Desc: "book about go lang" },
}

func nextId() int {
  max := -1
  for _, event := range events {
    if max < event.Id {
      max = event.Id
    }
  }
  return max + 1
}

func getById(id int) event {
  for _, ev := range events {
    if ev.Id == id {
      return ev
    }
  }
  return event{}
}


func main() {

  const b = 10
  fmt.Println(b)

  start := time.Now()
  for ;time.Now().Sub(start) < 1000; {
    fmt.Println(time.Now().Sub(start));
  }

  log.Fatal(http.ListenAndServe(":8089", initializeRouter()))

}

func initializeRouter() *mux.Router {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", index)
  //rest api
  router.HandleFunc("/event", createEvent).Methods("POST")
  router.HandleFunc("/event", listEvents).Methods("GET")
  router.HandleFunc("/event/{id}", getEvent).Methods("GET")
  router.HandleFunc("/event", updateEvent).Methods("PUT")
  router.HandleFunc("/event", deleteEvent).Methods("DELETE")
  
  //middleware
  router.Use(loggingMiddleware)

  return router
}

func index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello world")
}

func listEvents(w http.ResponseWriter, r *http.Request) {
  fmt.Println("listing events")
  json.NewEncoder(w).Encode(events)
}

func getEvent(w http.ResponseWriter, r *http.Request) {
  eventId := mux.Vars(r)["id"]
  fmt.Println("getting event with id: " + eventId)

  id, _ := strconv.Atoi(eventId)
  json.NewEncoder(w).Encode(getById(id))
}

func createEvent(w http.ResponseWriter, r *http.Request) {

  body, err := ioutil.ReadAll(r.Body)
  
  if err != nil {
    fmt.Fprintf(w, "error while reading body of the message")
  }

  var newEvent event
  err = json.Unmarshal(body, &newEvent)

  if err != nil {
    fmt.Fprintf(w, "bad format format of message body")
  }

  newEvent.Id = nextId()
  events = append(events, newEvent)
  
  json.NewEncoder(w).Encode(newEvent)

  fmt.Printf("received message with the body: %+v\n", newEvent)
  fmt.Println("creating event")
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
  fmt.Println("updating event")
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
  fmt.Println("deleting event")
}


func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        log.Println(r)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}
