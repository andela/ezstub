# ezstub

Pronounced "easy stub", is an easy way to stub API endpoints with the use of a simple yaml configuration file.

### Usage 
Create a `ezstub.yaml` file.
```yaml
title: Stubs for some test API
port: 8080
endpoints:
- url: /users
  description: Get all users
  response:
    headers:
     - key: Content-Type
       value: application/json
    data: WwogICAgewogICAgICAgICJpZCIgOiAxLAogICAgICAgICJuYW1lIjogIkpvaG4gU25vdyIKICAgIH0sCiAgICB7CiAgICAgICAgImlkIiA6IDIsCiAgICAgICAgIm5hbWUiOiAiTG9yZXAgSXBzdW0iCiAgICB9Cl0=
    status: 200
```
Start ezstub
```sh
$ ezstub -c ezstub.yaml
Stubs for some test API
ezstub listening on :8080
```
Test the endpoint
```sh
$ curl http://localhost:8080/users
[
    {
        "id" : 1,
        "name": "John Snow"
    },
    {
        "id" : 2,
        "name": "Lorep Ipsum"
    }
]
```
### Docs
The yaml configuration format.

#### Top level
* title [string]: title for the API configuration.
* port [int]: port the server should listen on
* host [string]: ip/hostname the server should bind the ip to.
* endpoints [array]: array of [endpoints](#endpoint).

#### endpoint 
* url [string]: endpoint url.
* description [string]: endpoint description.
* validation [array]: array of [validations](#validation).
* response: [response](#response).

#### validation
Requests missing the following [key-values](#key-value) will get a 403 (Forbidden) response.
* headers [array]: Request headers. Array of [key-values](#key-value). 
* params [array]: Form/query parameters. Array of [key-values](#key-value). 

#### response
 One of `data` or `file` should be used.
 * data [string]: base64 encoding of a string.
 * file [string]: file path
 * status [int]: status code. 

#### key-value
* key [string]: key
* value [string]: value
