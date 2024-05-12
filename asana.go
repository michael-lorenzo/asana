package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	PersonalAccessToken string `toml:"personal_access_token"`
	WorkspaceGid        string `toml:"workspace_gid"`
}

type Payload struct {
	Data Data `json:"data"`
}
type Data struct {
	Name      string `json:"name"`
	Assignee  string `json:"assignee"`
	Workspace string `json:"workspace"`
	Notes     string `json:"notes"`
}

func main() {
	f, err := os.Open("asana.toml")
	if err != nil {
		log.Fatalln("Error opening config file: asana.toml")
	}
	d := toml.NewDecoder(f)
	var c Config
	_, err = d.Decode(&c)
	if err != nil {
		log.Fatalln("Error bad formatting in config file: asana.toml")
	}

	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalln("Task must have a name")
	}

	data := Payload{
		Data: Data{
			Name:      flag.Arg(0),
			Assignee:  "me",
			Workspace: c.WorkspaceGid,
			Notes:     flag.Arg(1),
		},
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://app.asana.com/api/1.0/tasks", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.PersonalAccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error connecting to Asana API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		log.Fatalf("Error status code %d\n", resp.StatusCode)
	}
}
