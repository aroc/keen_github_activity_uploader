package client

import (
	"net/http"
	"io"
	"os"
	"compress/gzip"
	"bufio"
	"fmt"
	"encoding/json"
	"os/exec"
)

const (
	DataUrl = "http://data.githubarchive.org/"
	FileExtension = ".json.gz"
)

// The Uploader for each timeframe
type Uploader struct {
	timeframe 						 string
	archivePath 					 string
	rawJSONFilePath    		 string
	updatedJSONFilePath    string
}

// public interface
func NewUploader(timeframe string) *Uploader {
	upl := new(Uploader)
	upl.timeframe = timeframe
	upl.archivePath = ""
	upl.rawJSONFilePath = ""
	upl.updatedJSONFilePath = ""
	return upl
}

// Private methods
// *******************

func alterJSONObject(obj map[string]interface{}) map[string]interface{} {
	var keenObj = map[string]string{
		"timestamp": obj["created_at"].(string),
	}
	obj["keen"] = keenObj
	delete(obj, "created_at")
	return obj
}

func (upl Uploader) buildDownloadUrl() string {
    return DataUrl + upl.timeframe + FileExtension
}

func (upl *Uploader) downloadArchive() {
	fmt.Println("Downloading archive for " + upl.timeframe)
	path := "data_files/archives/" + upl.timeframe + FileExtension
	out, err := os.Create(path)
	defer out.Close()
	check(err)

	resp, err := http.Get(upl.buildDownloadUrl())
	defer resp.Body.Close()
	check(err)

	_, err = io.Copy(out, resp.Body)
	check(err)

	upl.archivePath = path
}

func (upl *Uploader) createJSONFile() {
	fmt.Println("Creating json file for " + upl.timeframe)
	// Open the file
	file, err := os.Open(upl.archivePath)
	check(err)
	defer file.Close()

	// Read the gzipped file
	fz, err := gzip.NewReader(file)
	check(err)
	defer fz.Close()

	// Write to a new json file
	path := "data_files/json/" + upl.timeframe + "-raw.json"
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(out, fz)
	check(err)

	upl.rawJSONFilePath = path
}

func (upl *Uploader) alterJSONFile() {
	fmt.Println("Updating json file for " + upl.timeframe)
	upl.updatedJSONFilePath = "data_files/json/" + upl.timeframe + "-updated.json"

	// Open the file
	file, err := os.Open(upl.rawJSONFilePath)
	check(err)
	defer file.Close()

	// Create the new json file
	newFile, err := os.Create(upl.updatedJSONFilePath)
  check(err)

  // Write the opening bracket of the JSON array of objects.
  _, err = newFile.WriteString("[")
  check(err)
		
	// Read the file line by line and update the JSON objects
	scanner := bufio.NewScanner(file)
	count := 0;
	for scanner.Scan() {
		if count > 0 {
			_, err = newFile.WriteString(",")
  		check(err)
		}
		var dat map[string]interface{}
		bytes := scanner.Bytes()
		err = json.Unmarshal(bytes, &dat)
		check(err)

    alteredDat := alterJSONObject(dat)
    encoded, err := json.Marshal(alteredDat)
    check(err)

    // Write to new updated JSON file
    _, err = newFile.Write(encoded)
    check(err)
    count++
	}
	
	_, err = newFile.WriteString("]")
  check(err)

	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func (upl *Uploader) uploadJSON() {
	fmt.Println("Uploading json to keen for " + upl.timeframe)

	// Create the log files
	logfile, err := os.Create("logs/results-" + upl.timeframe + ".log")
  	check(err)

	errLogfile, err := os.Create("logs/errors-" + upl.timeframe + ".log")
	check(err)

	cmd := exec.Command("keen",
		"events:add",
		"-c",
		"github_activities",
		"--batch-size",
		"2000",
		"-f",
		upl.updatedJSONFilePath,
	)

	cmd.Stderr = errLogfile
	cmd.Stdout = logfile

	err = cmd.Run()
	check(err)
}

// Public methods
// *******************

func (upl *Uploader) Run() {
	upl.downloadArchive()
	upl.createJSONFile()
	upl.alterJSONFile()
	upl.uploadJSON()
	fmt.Println("***DONE uploading to keen for " + upl.timeframe)
}