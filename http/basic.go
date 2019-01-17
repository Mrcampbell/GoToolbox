package main

// http://polyglot.ninja/golang-making-http-requests/

/**
* NOTES:
* Always close the response Body
* Always use a timeout:
```
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
			Timeout: timeout,
	}
	client.Get(url)
```
*/

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

// MakeGetRequest is a simple GET call
func MakeGetRequest() {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

func MakePostRequest() {
	message := map[string]interface{}{
		"hello": "world",
		"life":  42,
		"embedded": map[string]string{
			"yes": "of course!",
		},
	}

	byteRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(byteRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)
	log.Println(result)
	log.Println(result["data"])
}

func MakeFormPostRequest() {

	formData := url.Values{
		"name": {"Mike"},
	}

	resp, err := http.PostForm("https://httpbin.org/post", formData)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result["form"])
}

func MakeFilePostRequest() {

	// Open the file
	file, err := os.Open("name.txt")
	if err != nil {
		log.Fatalln(err)
	}
	// Close the file later
	defer file.Close()

	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file_field", "name.txt")
	if err != nil {
		log.Fatalln(err)
	}

	// Copy the actual file content to the field field's writer
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Fatalln(err)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("normal_field")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = fieldWriter.Write([]byte("Value"))
	if err != nil {
		log.Fatalln(err)
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	req, err := http.NewRequest("POST", "https://httpbin.org/post", &requestBody)
	if err != nil {
		log.Fatalln(err)
	}
	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	// Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(response.Body).Decode(&result)

	log.Println(result)
}

func main() {
	MakeGetRequest()
	MakePostRequest()
	MakeFormPostRequest()
	MakeFilePostRequest()
}
