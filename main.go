package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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

	for i, url := range result {
		i++
		fmt.Printf("#%d - Downloading %s\n", i, url)
		DownloadImages(url, strconv.Itoa(i)+"golang.jpg")
		fmt.Printf("#%d - Completed %s\n", i, url)
	}
}

// DownloadImages download the image from url and save it to data folder
func DownloadImages(url string, filepath string) error {
	CreateDirIfNotExist("data")
	// Create the file
	out, err := os.Create("data/" + filepath)
	checkError(err)
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
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
