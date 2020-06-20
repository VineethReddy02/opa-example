# OPA example

This repository contains two services i.e client & server. On request invocation, we query OPA to decide whether the request is allowed to server or to be denyed by the client itself.

#### Running client:

```
./bin/client
 ```

#### Running server:

```
./bin/server
```

#### Rego policy file

The below policy enforces the rule only employee can view own details

```
package employees.authz

default allow = false

allow {
	some name
	input.method == "GET"
	input.path = ["employee", name]
	name == input.user
}

```

#### Running Open Policy Agent:

```
curl -L -o opa https://openpolicyagent.org/downloads/latest/opa_linux_amd64
chmod 755 ./opa
opa run --server test.rego 
```

#### Now server should show the data it holds:

Visiting http://localhost:8090/ should show the below data.

```
[
  {
    "name": "arya",
    "manager": "john",
    "dept": "stark",
    "salary": "$60,000"
  },
  {
    "name": "Jamie",
    "manager": "Cersei",
    "dept": "Lannister",
    "salary": "$2,00,000"
  },
  {
    "name": "John Snow",
    "manager": "Daenerys",
    "dept": "Targareon",
    "salary": "$1,50,000"
  },
  {
    "name": "Tyrian",
    "manager": "none",
    "dept": "Hand of king",
    "salary": "$3,00,000"
  },
  {
    "name": "Daenerys",
    "manager": "Tyrian",
    "dept": "Targareon",
    "salary": "$4,00,000"
  }
]

```

Now doing a GET request to the client application i.e http://localhost:8089/employee/arya with request body as

```
{
        "user": "arya"
}
```

Will result in response

```
{
  "name": "arya",
  "manager": "john",
  "dept": "stark",
  "salary": "$60,000"
}
```

cheers!
