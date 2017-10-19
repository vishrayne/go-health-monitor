package monit

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// ReportCategoryTemp creates the file under /tmp
	ReportCategoryTemp int = 0
	// ReportCategoryFinal creates the file under reports folder
	ReportCategoryFinal int = 1
	// ReportFolderName points to the folder where reports are stored
	ReportFolderName string = "reports"
)

// Report representation
type Report struct {
	Name     string
	FilePath string
	Category int
	file     *os.File
}

func timeStamp(now time.Time) string {
	timeStamp := now.UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(int64(timeStamp), 10)
}

// CreateTempReport creates a new report under temp directory.
func CreateTempReport(name string) *Report {
	return createReport(name, ReportCategoryTemp)
}

// CreateSystemReport creates a new report under reports directory.
func CreateSystemReport(name string) *Report {
	return createReport(name, ReportCategoryFinal)
}

// Category should be either ReportCategoryTemp or ReportCategoryFinal
func createReport(name string, category int) *Report {
	now := time.Now()
	filename := fmt.Sprintf("%s_%s.txt", name, timeStamp(now))
	header := fmt.Sprintf("filename: %s\ncreated: %v", filename, now)

	var filePath string
	switch category {
	case ReportCategoryTemp:
		filePath = fmt.Sprintf("%s/%s", "/tmp", filename)
	case ReportCategoryFinal:
		// create report folder if not exists!
		if _, err := os.Stat(ReportFolderName); os.IsNotExist(err) {
			os.Mkdir(ReportFolderName, os.ModePerm)
		}

		filePath = fmt.Sprintf("%s/%s", ReportFolderName, filename)
		header = fmt.Sprintf(reportFormatFor(ReportCategoryFinal), "info", header)
	default:
		dealWithError("create report", errors.New("invalid value for report category"))
	}

	f, err := os.Create(filePath)
	dealWithError("create report", err)

	report := &Report{
		Name:     filename,
		FilePath: filePath,
		Category: category,
		file:     f,
	}

	report.write(header)
	return report
}

func (report *Report) write(text string) {
	_, err := report.file.WriteString(text)
	dealWithError("report", err)
}

func (report *Report) writeln(text string) {
	report.write(text + "\n")
}

func (report *Report) writeSection(title string, text string) {
	text = strings.TrimRight(text, "\r\n")
	report.writeln(fmt.Sprintf(reportFormatFor(report.Category), title, text))
}

func (report *Report) close() {
	log.Println("report closed!")
	report.file.Close()
}

func reportFormatFor(category int) string {
	switch category {
	case ReportCategoryTemp:
		return "%s:%s"
	default:
		return "\n%s\n=====================\n%s"
	}
}
