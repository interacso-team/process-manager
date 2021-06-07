package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type config struct {
	App string `json:"app,omitempty"`
}

type appPostData struct {
	Name       string    `json:"name"`
	Command    string    `json:"command"`
	Repository string    `json:"repository"`
	Folder     string    `json:"folder"`
	Variables  variables `json:"variables"`
}

func main() {

	http.HandleFunc("/stop", func(w http.ResponseWriter, peticion *http.Request) {
		decoder := json.NewDecoder(peticion.Body)
		var body config
		err := decoder.Decode(&body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		println(body.App)
		stopserver(mapa1[body.App])
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, peticion *http.Request) {
		decoder := json.NewDecoder(peticion.Body)
		var body config
		err := decoder.Decode(&body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		println(body.App)
		go execserver(mapa1[body.App])
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, peticion *http.Request) {
		decoder := json.NewDecoder(peticion.Body)
		var body config
		err := decoder.Decode(&body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		println(body.App)
		status := serverstatus(mapa1[body.App])
		if status {
			w.Write([]byte("Runnig"))
		} else {
			w.Write([]byte("Not running"))
		}
	})

	http.HandleFunc("/status-all", func(w http.ResponseWriter, peticion *http.Request) {
		status := statusall()
		b, err := json.Marshal(status)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)

	})

	http.HandleFunc("/clonerepo", func(w http.ResponseWriter, peticion *http.Request) {
		status := statusall()
		b, err := json.Marshal(status)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)

	})

	http.HandleFunc("/createapp", func(w http.ResponseWriter, peticion *http.Request) {
		decoder := json.NewDecoder(peticion.Body)
		var body appPostData
		err := decoder.Decode(&body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Println(body)
		createapp(body)
		clonerepo(body)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, peticion *http.Request) {
		decoder := json.NewDecoder(peticion.Body)
		var body appPostData
		err := decoder.Decode(&body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		println("body.App", body.Name)
		app, ok := mapa1[body.Name]
		if !ok {
			log.Println("Application not found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Application not found"))
			return
		}
		deleteserver(app)
		w.Write([]byte("OK"))
	})

	direccion := ":8080"
	fmt.Println("Servidor listo escuchando en " + direccion)
	http.ListenAndServe(direccion, nil)
}
