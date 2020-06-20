package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type result struct {
	Result employees `json:"result"`
}

type employees struct {
	Employees authz `json:"employees"`
}

type authz struct {
	Authz allow `json:"authz"`
}

type allow struct {
	Allow bool `json:"allow"`
}

type inputReq struct {
	Input clientRequest `json:"input"`
}

type clientRequest struct {
	Method string `json:"method"`
	Path []string `json:"path"`
	User string `json:"user"`
}

type employee struct {
	Name    string  `json:"name"`
	Manager string  `json:"manager"`
	Dept    string  `json:"dept"`
	Salary  string  `json:"salary"`
}


func main() {
	router := mux.NewRouter()
	router.HandleFunc("/employee/{id}", getEmployeeDetails).Methods(http.MethodGet)
	http.ListenAndServe(":8089", router)
}

func getOPADecision( r *http.Request) (result, *inputReq) {
	url := "http://localhost:8181/v1/data"
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var input clientRequest
	err = json.Unmarshal(reqbody, &input)
	if err != nil {
		panic(err)
	}

	opaReq := &inputReq{
		Input: clientRequest{
			Method: "GET",
			Path:   path,
			User:   input.User,
		},
	}

	req01, err := json.Marshal(opaReq)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(req01))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var res result
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println("Request to OPA server:", opaReq)
	return res, opaReq
}

func getEmployeeDetails(w http.ResponseWriter, r *http.Request) {
	serverURL := "http://localhost:8090/"
	res, clientReq := getOPADecision(r)
	fmt.Println("response Body:", res.Result.Employees.Authz.Allow)

	if res.Result.Employees.Authz.Allow {
		// request the server for details
		b, err := json.Marshal(clientReq)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("GET", serverURL, bytes.NewBuffer(b))
		if err != nil {
			panic(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Print(string(body))

		var empResponse employee
		err = json.Unmarshal(body, &empResponse)
		if err != nil {
			panic(err)
		}

		io.WriteString(w,string(body))
	} else {
		// fail the request as the user doesn't have permissions
		io.WriteString(w, fmt.Sprintf("failed to process request for %s", clientReq.Input.User))
		fmt.Printf("failed to process request %s", clientReq.Input.User)
	}
}
