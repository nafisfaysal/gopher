package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	// read the json file
	data, err := ioutil.ReadFile("images.json")
	if err != nil {
		fmt.Errorf("can't read json data %s", err)
		return
	}

	// unmarshal json file
	var result []string
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Errorf("error on unmarshal json %s", err)
		return
	}

	jobs := make(chan string, len(result))
	rst := make(chan string, len(result))

	for w := 1; w <= 3; w++ {
		go DownloadImages(w, jobs, rst, strconv.Itoa(w)+"golang.jpg")
	}

	for i, url := range result {
		i++
		fmt.Printf("#%d - Downloading %s\n", i, url)
		jobs <- url
		fmt.Printf("#%d - Completed %s\n", i, url)
	}

	close(jobs)

	for a := 1; a <= 5; a++ {
		<-rst
	}
}

// DownloadImages download the image from url and save it to data folder
func DownloadImages(id int, jobs <-chan string, results chan<- string, filepath string) error {
	fmt.Printf("#%d is the id of the thread number\n", id)
	for j := range jobs {
		time.Sleep(time.Second)
		CreateDirIfNotExist("data")
		// Create the file
		out, err := os.Create("data/" + filepath)
		checkError(err)
		defer out.Close()

		// Get the data
		resp, err := http.Get(j)
		if err != nil {
			fmt.Printf("got bad request: %s", err)
			return err
		}
		defer resp.Body.Close()

		// Check server response
		if resp.StatusCode != http.StatusOK {
			fmt.Errorf("bad status: %s", resp.Status)
			return err
		}

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		results <- j
	}

	return nil
}

// CreateDirIfNotExist create directory if the directory doesn't exist
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

// CheckError check the error
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
