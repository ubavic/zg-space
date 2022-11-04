package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var usersCache map[string]bool
var timeRan time.Time

func getUsers() []string {
	if time.Now().Sub(timeRan).Seconds() > 2 {
		timeRan = time.Now()

		cmd := exec.Command("ps", "-ef")
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			panic(0)
		} else {
			lines := strings.Split(string(stdout), "\n")
			for _, x := range lines {
				fields := strings.Fields(x)

				if len(fields) != 0 {
					user := fields[0]

					if _, ok := usersCache[user]; ok {
						usersCache[user] = true
					}
				}
			}
		}
	}

	users := make([]string, len(usersCache))
	for k, v := range usersCache {
		if v {
			users = append(users, k)
		}
	}

	return users
}

func main() {
	cmd := exec.Command("ls", "/home")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		panic(0)
	}

	homeDirs := strings.Fields(string(stdout))

	usersCache = make(map[string]bool, len(homeDirs))

	for _, homeDir := range homeDirs {
		if homeDir != "" {
			usersCache[homeDir] = false
		}
	}

	timeRan = time.Now()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(getUsers())
	})

	fmt.Printf("Starting on :9990")
	http.ListenAndServe(":9990", nil)
}
