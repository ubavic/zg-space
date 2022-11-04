package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var homeDirs []string
var usersCache []string
var timeRan time.Time

func contains(s []string, str string) bool {
	for _, val := range s {
		if val == str {
			return true
		}
	}

	return false
}

func getUsers() []string {
	if time.Now().Sub(timeRan).Seconds() > 1 {
		timeRan = time.Now()

		cmd := exec.Command("ps", "-ef")
		stdout, err := cmd.Output()

		users := make([]string, len(homeDirs))

		if err != nil {
			fmt.Println(err.Error())
			panic(0)
		} else {
			lines := strings.Split(string(stdout), "\n")

			for _, x := range lines {
				fields := strings.Fields(x)

				if len(fields) == 0 {
					continue
				}

				user := fields[0]

				var regularUser = contains(homeDirs, user)
				var alreadyOnList = contains(users, user)

				if regularUser && !alreadyOnList {
					users = append(users, user)
				}
			}

			usersCache = users
		}
	}

	return usersCache
}

func main() {
	cmd := exec.Command("ls", "/home")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		panic(0)
	}

	homeDirs_ := strings.Fields(string(stdout))
	usersCache = make([]string, len(homeDirs))

	for _, homeDir := range homeDirs_ {
		if homeDir != "" {
			homeDirs = append(homeDirs, homeDir)
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
