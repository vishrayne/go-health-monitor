package monit

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	sh "github.com/codeskyblue/go-sh"
)

const (
	//AllowedLogEntries constant
	AllowedLogEntries = 200
	//DefaultRegex used for parsing apache access
	DefaultRegex = `(?P<remote_addr>[^ ]*) - - \[(?P<time_local>[^]]*)\] TIME: (?P<time>[^ ]*) "(?P<request>[^\\]*)" (?P<status>[^ ]*) (?P<size>[^ ]*)`
	//DateTimeLayout used to parse date
	DateTimeLayout = `02/Jan/2006:15:04:05 -0700`
)

/*LogTask representation*/
type accessLog struct {
	fileName       string
	regexp         string
	TotalRequest   int         `json:"total_requests"`
	InvalidRequest int         `json:"unparsable_requests"`
	StatusMap      map[int]int `json:"status_count"`
	StartTime      string      `json:"first_entry_at"`
	EndTime        string      `json:"last_entry_at"`
	Interval       string      `json:"time_interval"`
	Success        string      `json:"success_rate"`
	LogEntries     entry       `json:"entries"`
}

/*LogEntry representation*/
type logEntry struct {
	RemoteAddress string `json:"remote_addr"`
	Request       string `json:"request"`
	Size          int    `json:"response_bytes"`
	Status        int    `json:"status"`
	ElapsedTime   string `json:"elapsed_time"`
	TimeLocal     string `json:"time_local"`
}

/*Entry representation*/
type entry map[string][]logEntry

func newAccessLogParser() *accessLog {
	return &accessLog{}
}

func (al *accessLog) toJSON() string {
	return asJSON(al)
}

func (al *accessLog) parse(lineCount int, filepath string, includeEntries bool) {
	if lineCount > AllowedLogEntries {
		dealWithError("access log parser", fmt.Errorf("parsing more than %d lines is not supported, please reduce the configuration value for `line_count` to %d or less", AllowedLogEntries, AllowedLogEntries))
	}

	logs := al.tail(filepath, lineCount)
	scanner := bufio.NewScanner(bytes.NewReader(logs))
	accessLogRegexp := regexp.MustCompile(al.lookupRegexp())
	keys := accessLogRegexp.SubexpNames()
	noOfKeys := len(keys)

	if noOfKeys <= 1 {
		dealWithError("access log parser", errors.New("<empty_keyset> Please check the regexp used for parsing log"))
	}

	resultMap := make(map[string][]logEntry)
	statusMap := make(map[int]int)
	successfulCalls := 0

	for scanner.Scan() {
		text := scanner.Text()
		entriesMap := al.parseLogEntry(accessLogRegexp, keys, text)
		if entriesMap == nil || len(entriesMap) <= 0 {
			log.Printf("Skipping invalid access log entry: %s", text)
			al.InvalidRequest++
			continue
		}

		newLogEntry := logEntry{
			RemoteAddress: entriesMap[keys[1]],
			TimeLocal:     entriesMap[keys[2]],
			ElapsedTime:   entriesMap[keys[3]],
			Request:       entriesMap[keys[4]],
			Status:        al.stringToInt(entriesMap[keys[5]]),
			Size:          al.stringToInt(entriesMap[keys[6]]),
		}

		if len(al.StartTime) <= 0 {
			al.StartTime = newLogEntry.TimeLocal
		}
		al.EndTime = newLogEntry.TimeLocal

		resultMap = al.addLogEntryByStatus(resultMap, newLogEntry)
		statusMap[newLogEntry.Status]++

		if newLogEntry.Status >= 200 && newLogEntry.Status < 300 {
			successfulCalls++
		}
	}

	if includeEntries {
		al.LogEntries = resultMap
	}

	al.TotalRequest = lineCount - al.InvalidRequest
	al.StatusMap = statusMap

	al.Interval = al.timeInterval()

	successPercent := float64(successfulCalls) * 100 / float64(al.TotalRequest)
	al.Success = strconv.FormatFloat(successPercent, 'f', 2, 64)
}

func (al *accessLog) timeInterval() string {
	if len(al.StartTime) <= 0 || len(al.EndTime) <= 0 {
		log.Print("empty timstamp!")
		return ""
	}

	t1, err := time.Parse(DateTimeLayout, al.StartTime)
	if err != nil {
		log.Printf("unable to parse start time: %v", err)
		return ""
	}

	t2, err := time.Parse(DateTimeLayout, al.EndTime)
	if err != nil {
		log.Printf("unable to parse end time: %v", err)
		return ""
	}

	minutes := strconv.FormatFloat(t2.Sub(t1).Minutes(), 'f', 2, 64)
	return fmt.Sprintf("%s minutes", minutes)
}

func (al *accessLog) lookupRegexp() string {
	// if len(defaultRegexp) <= 0 {
	// 	log.Println("No used defined regex found. Falling back to default...")
	// 	defaultRegexp = DefaultRegex
	// }

	return DefaultRegex
}

func (al *accessLog) addLogEntryByStatus(resultMap map[string][]logEntry, newLogEntry logEntry) map[string][]logEntry {
	statusKey := al.intToString(newLogEntry.Status)

	var list []logEntry
	if entryList, ok := resultMap[statusKey]; ok {
		list = entryList
	}
	resultMap[statusKey] = append(list, newLogEntry)

	return resultMap
}

func (al *accessLog) tail(fileName string, lineCount int) []byte {
	output, err := sh.Command("tail", fmt.Sprintf("-n %d", lineCount), fileName).Output()
	dealWithError("access_log", err)

	return output
}

func (al *accessLog) parseLogEntry(expression *regexp.Regexp, keys []string, text string) map[string]string {
	matches := expression.FindStringSubmatch(text)
	if len(matches) <= 0 {
		return nil
	}

	parsedEntriesMap := make(map[string]string)
	for i, key := range keys {
		if i > 0 {
			parsedEntriesMap[key] = matches[i]
		}
	}
	return parsedEntriesMap
}

func (al *accessLog) stringToInt(str string) int {
	number, _ := strconv.Atoi(str)
	return number
}

func (al *accessLog) intToString(number int) string {
	return strconv.Itoa(number)
}
