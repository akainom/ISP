package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	baseUrl    = "http://webdav:webdav123@localhost:3000/"
	destUrl    = "http://localhost:3000/"
	baseFolder = "A"
	workFolder = "B"
)

func request(method, url string, body io.Reader, h map[string]string) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println(method, "error: ", err)
	}
	for k, v := range h {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(method, "error: ", err)
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	fmt.Println("=== ", method, url, " ===")
	fmt.Println("status: ", resp.Status)
	fmt.Println(string(b), "\n")
}

func main() {
	request("MKCOL", baseUrl+baseFolder+"/", nil, nil)
	request("PUT", baseUrl+baseFolder+"/hello.txt",
		bytes.NewBuffer([]byte("Hello WebDAV")),
		map[string]string{"Content-Type": "text/plain"},
	)
	request("GET", baseUrl+baseFolder+"/hello.txt", nil, nil)
	request("PROPFIND", baseUrl+baseFolder+"/",
		bytes.NewBuffer([]byte(`<?xml version="1.0"?>
<propfind xmlns="DAV:">
    <allprop/>
</propfind>`)),
		map[string]string{"Depth": "1", "Content-Type": "application/xml"},
	)
	request("PROPPATCH", baseUrl+baseFolder+"/hello.txt",
		bytes.NewBuffer([]byte(`<?xml version="1.0"?>
<propertyupdate xmlns="DAV:">
    <set>
        <prop>
            <description>test file</description>
        </prop>
    </set>
</propertyupdate>`)),
		map[string]string{"Content-Type": "application/xml"},
	)
	request("MKCOL", baseUrl+workFolder+"/", nil, nil)
	request("COPY", baseUrl+baseFolder+"/hello.txt",
		nil,
		map[string]string{"Destination": destUrl + workFolder + "/hello_copy.txt"},
	)
	request("MOVE", baseUrl+workFolder+"/hello_copy.txt",
		nil,
		map[string]string{"Destination": destUrl + baseFolder + "/hello_moved.txt"},
	)
	request("DELETE", baseUrl+baseFolder+"/", nil, nil)
	request("DELETE", baseUrl+workFolder+"/", nil, nil)
}
