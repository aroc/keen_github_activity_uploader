package main

import (
    "fmt"
    "os"
    "time"
    "strconv"
    "github.com/aroc/keen_github_uploader/client"
)

// Helpers
func check(e error) {
  if (e != nil) {
    panic(e)
  }
}

func yesterday() string {
    negativeOneday, err := time.ParseDuration("-24h")
    check(err)
    yesterday := time.Now().Add(negativeOneday)
    return formatTimeString(yesterday)
}

func formatTimeString(userTime time.Time) string {
    year := strconv.Itoa(userTime.Year())
    month := fmt.Sprintf("%02v", strconv.Itoa(int(userTime.Month())))
    day := fmt.Sprintf("%02v", strconv.Itoa(userTime.Day()))
    return year + "-" + month + "-" + day
}

func checkTime(userTime string) {
    layout := "2006-01-02"
    _, err := time.Parse(layout, userTime)
    check(err)
}

func main() {
    var targetTime string
    
    if len(os.Args) > 1 {
        userTime := os.Args[1]
        checkTime(userTime)
        targetTime = userTime
    } else  {
        targetTime = yesterday()
    }

    uploader := client.NewUploaderController(targetTime)
    uploader.Run()
	fmt.Println("Uploading complete")
}