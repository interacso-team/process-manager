package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type variables map[string]string

type application struct {
	Command    string    `json:"-"`
	Process    *exec.Cmd `json:"-"`
	Running    bool      `json:"running"`
	Folder     string    `json:"folder"`
	Repository string    `json:"repository"`
	Name       string    `json:"name"`
	Variables  variables `json:"variables"`
}

var mapa1 map[string]*application = map[string]*application{}

func stopserver(App *application) {
	if App != nil && App.Process != nil {
		App.Process.Process.Kill()
	}

}

func deleteserver(App *application) {
	if App.Process != nil && App.Process.Process != nil {
		App.Process.Process.Kill()
	}
	delete(mapa1, App.Name)
}

func serverstatus(App *application) bool {
	return App.Running
}

func statusall() map[string]*application {
	return mapa1
}

func createapp(data appPostData) {
	stopserver(mapa1[data.Name])
	mapa1[data.Name] = &application{
		Command:    data.Command,
		Folder:     data.Folder,
		Repository: data.Repository,
		Name:       data.Name,
		Variables:  data.Variables,
	}
}

func execserver(App *application) {
	c := strings.Split(App.Command, " ")
	var env []string
	for k, v := range App.Variables {
		env = append(env, k+"="+v)
	}

	App.Process = exec.Command(c[0], c[1:]...)
	App.Process.Env = env
	App.Process.Dir = App.Folder
	App.Running = true
	stdout, _ := App.Process.StdoutPipe()
	stderr, _ := App.Process.StderrPipe()
	err := App.Process.Start()
	if err != nil {
		log.Printf("App.Process.Start() failed with '%s'\n", err)
		App.Running = false
		return

	}

	go func() {
		r := bufio.NewReader(stdout)
		line, _, _ := r.ReadLine()
		println(string(line))
	}()

	go func() {
		r := bufio.NewReader(stderr)
		line, _, _ := r.ReadLine()
		println(string(line))
	}()

	// cmd.Process.Kill()

	err = App.Process.Wait()
	if err != nil {
		log.Printf("App.Process.Start() failed with '%s'\n", err)
	}
	println("Terminado")
	App.Running = false
}

func clonerepo(data appPostData) {
	cmd := exec.Command("git", "clone", data.Repository, data.Folder)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}
