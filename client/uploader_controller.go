package client

import (
	"strconv"
)

// Helpers
func check(e error) {
  if (e != nil) {
    panic(e)
  }
}

// The Controllers for all the Uploaders.
type UploaderController struct {
	dataUrl       string
    targetDate    string
	fileExtension string
	uploaders     []*Uploader
}

// public interface
func NewUploaderController(targetDate string) *UploaderController {
	uplCtl := new(UploaderController)
    uplCtl.targetDate = targetDate
	uplCtl.uploaders = createUploaders(uplCtl.targetDate)
	return uplCtl
}

func createUploaders(targetDate string) []*Uploader {
    uploaders := make([]*Uploader, 0)

    for i := 1; i < 24; i++ {
        timeframe := targetDate + "-" + strconv.Itoa(i)
        uploaders = append(uploaders, NewUploader(timeframe))
    }

    return uploaders
}

func (uplCtl *UploaderController) Run() {
    for i := 0; i < len(uplCtl.uploaders); i++ {
        uplCtl.uploaders[i].Run()
    }
}