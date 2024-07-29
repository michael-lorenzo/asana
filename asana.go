package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Payload struct {
	Data Data `json:"data"`
}

type Data struct {
	Name      string `json:"name"`
	Assignee  string `json:"assignee"`
	Workspace string `json:"workspace"`
	Notes     string `json:"notes"`
}

//go:embed asana.toml
var emptyConfig []byte

func main() {
	generateConfig := flag.Bool("generate", false, "Generate an empty config file")
	flag.Parse()

	if *generateConfig {
		uhd, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		emptyConfigPath := filepath.Join(uhd, "asana.toml")
		f, err := os.OpenFile(emptyConfigPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Generating empty config file at", emptyConfigPath)
		_, err = f.Write(emptyConfig)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	viper.SetConfigName("asana")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	if flag.NArg() == 0 {
		log.Fatal("Task must have a name")
	}

	data := Payload{
		Data: Data{
			Name:      flag.Arg(0),
			Assignee:  "me",
			Workspace: viper.GetString("workspace_gid"),
			Notes:     flag.Arg(1),
		},
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://app.asana.com/api/1.0/tasks", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+viper.GetString("personal_access_token"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if resp.StatusCode != 201 {
		log.Fatalf("Status code %v\n", resp.StatusCode)
	}
}
