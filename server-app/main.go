package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", employeeDetails)
	http.ListenAndServe(":8090", nil)
}

var (
	employees = []employee{{
		Name:    "arya",
		Manager: "john",
		Dept:    "stark",
		Salary:  "$60,000",
	}, {
		Name:    "Jamie",
		Manager: "Cersei",
		Dept:    "Lannister",
		Salary:  "$2,00,000",
	}, {
		Name:    "John Snow",
		Manager: "Daenerys",
		Dept:    "Targareon",
		Salary:  "$1,50,000",
	}, {
		Name:    "Tyrian",
		Manager: "none",
		Dept:    "Hand of king",
		Salary:  "$3,00,000",
	}, {
		Name:    "Daenerys",
		Manager: "Tyrian",
		Dept:    "Targareon",
		Salary:  "$4,00,000",
	},
	}
)

type employee struct {
	Name    string  `json:"name"`
	Manager string  `json:"manager"`
	Dept    string  `json:"dept"`
	Salary  string  `json:"salary"`
}

type inputReq struct {
	Input clientRequest `json:"input"`
}

type clientRequest struct {
	Method string `json:"method"`
	Path []string `json:"path"`
	User string `json:"user"`
}

func employeeDetails(writer http.ResponseWriter, reqest *http.Request) {
	reqbody, err := ioutil.ReadAll(reqest.Body)
	if err != nil {
		panic(err)
	}

	if len(reqbody) == 0 {
		emp, err := json.MarshalIndent(employees,"","  ")
		if err != nil {
			panic(err)
		}

		io.WriteString(writer, string(emp))
		return
	}

	var clientReq inputReq
	err = json.Unmarshal(reqbody, &clientReq)
	if err != nil {
		panic(err)
	}

	for _, e := range employees {
		if e.Name == clientReq.Input.User {
			fmt.Println("serving the request for ", e.Name)
			writer.Header().Set("Content-Type", "application/json")
			//writer.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(writer).Encode(e)
			if err != nil {
				panic(err)
			}
			return
		}
	}
}
