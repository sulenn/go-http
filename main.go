package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	println("hello world")
	resp, err :=   http.Get("http://023.node.internetapi.cn:21030/SCIDE/SCManager?action=executeContract&arg=%7B%22commit_hash%22%3A%223f2dc2ec21876d31ffa06744f215b6a17c5d3bad%22%2C%22repo_id%22%3A+%221002604%22%7D&contractID=RepositoryDB0&operation=getCommit")
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
