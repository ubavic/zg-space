package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var usersCache []string
var timeRan time.Time

func getUsers() []string {

	if time.Now().Sub(timeRan).Seconds() > 2 {
		timeRan = time.Now()

		cmd := exec.Command("who")
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			panic(0)
		} else {
			lines := strings.Split(string(stdout), "\n")
			users := make([]string, 0)
			for _, x := range lines {
				fields := strings.Fields(x)

				if len(fields) != 0 {
					users = append(users, fields[0])
				}
			}

			usersCache = users
		}
	}

	return usersCache
}

func main() {

	timeRan = time.Now()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		users := getUsers()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(users)
	})

	fmt.Printf("Starting on :9990")
	http.ListenAndServe(":9990", nil)
}
