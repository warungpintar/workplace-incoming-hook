package main

import (
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/nurza/logo"
	"github.com/warungpintar/workplaceincominghook/data"

	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
	Global variables
*/
var (
	// Logging
	l       logo.Logger
	Loggers []*logo.Logger

	// Configuration
	Thread              string // Bot's system channel
	ThreadGitlab        string
	ThreadTuleap        string
	ThreadAppCenter     string
	PushIcon            string // Push icon (Fb emoji)
	MergeIcon           string // Merge icon (Fb emoji)
	BuildIcon           string // Build icon (Fb emoji)
	BotStartMessage     string // Bot's start message
	FbAPIUrl            string // Fb API URL
	Verbose             bool   // Enable verbose mode
	ShowAllCommits      bool   // Show all commits rather than latest
	HttpTimeout         int    // Http timeout in second
	ChatType            string
	TuleapURL           string
	Port                string
	UrlNoteHookFunction string

	// Misc
	currentBuildID float64 = 0       // Current build ID
	n              string  = "%5CnX" // Encoded line return
)

type GitlabServ struct{}

/*
	Flags
*/
var (
	ConfigFile = flag.String("f", "config.json", "Configuration file")
)

const (
	Bot   int = iota
	Push  int = iota
	Merge int = iota
	Build int = iota
)

/*
	Load configuration file
*/
func LoadConf() {
	conf := struct {
		ThreadGitlab        string
		ThreadAppCenter     string
		ThreadTuleap        string
		PushIcon            string
		MergeIcon           string
		BuildIcon           string
		BotStartMessage     string
		FbAPIUrl            string
		Verbose             bool
		ShowAllCommits      bool
		HttpTimeout         float64
		ChatType            string
		TuleapURL           string
		Port                string
		UrlNoteHookFunction string
	}{}

	content, err := ioutil.ReadFile(*ConfigFile)
	if err != nil {
		l.Critical("Error: Read config file error: " + err.Error())
	}

	err = json.Unmarshal(content, &conf)
	if err != nil {
		l.Critical("Error: Parse config file error: " + err.Error())
	}

	PushIcon = conf.PushIcon
	MergeIcon = conf.MergeIcon
	BuildIcon = conf.BuildIcon
	BotStartMessage = conf.BotStartMessage
	FbAPIUrl = conf.FbAPIUrl
	Verbose = conf.Verbose
	ShowAllCommits = conf.ShowAllCommits
	HttpTimeout = int(conf.HttpTimeout)
	ChatType = conf.ChatType
	TuleapURL = conf.TuleapURL
	Port = conf.Port
	ThreadGitlab = conf.ThreadGitlab
	ThreadAppCenter = conf.ThreadAppCenter
	ThreadTuleap = conf.ThreadTuleap
	UrlNoteHookFunction = conf.UrlNoteHookFunction
}

/*
	HTTP POST request

	target:		url target
	payload:	payload to send

	Returned values:

	int:	HTTP response status code
	string:	HTTP response body
*/
func Post(target string, payload string) (int, string) {
	// Variables
	var err error          // Error catching
	var res *http.Response // HTTP response
	var req *http.Request  // HTTP request
	var body []byte        // Body response

	// Build request
	l.Debug(bytes.NewBufferString(payload))
	req, err = http.NewRequest("POST", target, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	// Do request
	client := &http.Client{}
	client.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   time.Duration(HttpTimeout) * time.Second,
			KeepAlive: time.Duration(HttpTimeout) * time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Duration(HttpTimeout) * time.Second,
	}

	res, err = client.Do(req)

	if err != nil {
		l.Error("Error : Curl POST : " + err.Error())
		if res != nil {
			return res.StatusCode, ""
		} else {
			return 0, ""
		}
	}
	defer res.Body.Close()

	// Read body
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		l.Error("Error : Curl POST body read : " + err.Error())
	}

	return res.StatusCode, string(body)
}

/*
	Encode the git commit message with replacing some special characters not allowed by the Fb API

	@param origin Git message to encode
*/

func MessageEncodeX(origin string) string {
	var result string = ""

	result = strings.Replace(origin, "%5CnX", "\\n\\n", -1)

	return result
}

func MessageEncode(origin string) string {
	var result string = ""

	for _, e := range strings.Split(origin, "") {
		switch e {
		case "\n":
			result += "%5CnX"
		case "&":
			result += " and "
		default:
			result += e
		}
	}
	return result
}

/*
	Send a message on WorkChat

	@param channel : Targeted channel (could be personal or group)
*/
func SendWorkchatMessage(channel, message string, chattype string) {
	// Variables
	var payload string // POST data sent to Fb
	// var icon string    // Fb emoji

	// toLower(channel)
	l.Silly("toLower =", channel)
	channel = strings.ToLower(channel)
	l.Silly("toLower =", channel)

	// POST Payload formating
	payload = ""
	if chattype == "group" {
		payload += `{"recipient": { "thread_key": "` + strings.ToLower(channel) + `"} , "message": { "text": "` + message + `"}}`
	} else {
		payload += `{"recipient": { "id": "` + strings.ToLower(channel) + `"} , "message": { "text": "` + message + `"}}`
	}

	// Debug information
	if Verbose {
		l.Debug("payload =", payload)
	}

	code, body := Post(FbAPIUrl, MessageEncodeX(payload))
	if code != 200 {
		l.Error("Error post, Fb API returned:", body)
	}

	// Debug information
	if Verbose {
		l.Debug("Fb API returned:", body)
	}
}

/*
	Handler function to handle http requests for push

	@param w http.ResponseWriter
	@param r *http.Request
*/
func (s *GitlabServ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer // Buffer to get request body
	var body string         // Request body (it's a json)

	// Log
	l.Info("Request")

	// Read http request header
	gitlabEvent := r.Header.Get("X-Gitlab-Event")
	if Verbose {
		l.Debug("Gitlab Event =", gitlabEvent)
	}

	// Read service params
	serviceParam := r.URL.Query().Get("service")
	if Verbose {
		l.Debug("Service =", serviceParam)
	}

	// Read http request body and put it in a string
	buffer.ReadFrom(r.Body)
	body = buffer.String()

	// Debug information
	if Verbose {
		l.Debug("JsonString receive =", body)
	}

	switch serviceParam {
	case "tuleap":
		TaskHandler(body)
	case "appcenter":
		AppCenterHandler(body)
	default:
		switch gitlabEvent {
		case "Push Hook":
			PushHandler(body)
		case "Merge Request Hook":
			MergeHandler(body)
		case "Build Hook":
			BuildHandler(body)
		case "Note Hook":
			CommentHandler(body)
		}
	}
}

// CommentHandler call cloud function which handle comment event on gitlab
// and send to workplace bot
func CommentHandler(body string) {

	// whether UrlNoteHookFunction has been set #see on config.json
	if UrlNoteHookFunction == "" {
		l.Error("url comment service not ")
		return
	}

	var json = []byte(body)
	req, err := http.NewRequest("POST", UrlNoteHookFunction, bytes.NewBuffer(json))
	req.Header.Set("X-Gitlab-Event", "Note Hook")
	req.Header.Set("Content-Type", "application/json")

	l.Info("Call service note hook")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		l.Error("Error : call [POST] microservice failed :", err)
	}
	defer resp.Body.Close()

	l.Info("Response status :", resp.Status)
	l.Info("Response header :", resp.Header)
}

func PushHandler(body string) {
	var j data.Push
	var err error           // Error catching
	var message string = "" // Bot's message
	var date time.Time      // Time of the last commit

	// Parse json and put it in a the data.Build structure
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		// Error
		l.Error("Error : Json parser failed :", err)
	} else {
		// Ok
		// Debug information
		if Verbose {
			l.Debug("JsonObject =", j)
		}

		// Send the message

		// Date parsing (parsing result example : 18 November 2014 - 14:34)
		date, err = time.Parse("2006-01-02T15:04:05Z07:00", j.Commits[0].Timestamp)
		var dateString = date.Format("02 Jan 06 15:04")

		// Message
		lastCommit := j.Commits[len(j.Commits)-1]
		commitCount := strconv.FormatFloat(j.Total_commits_count, 'f', 0, 64)
		if ShowAllCommits {
			message += "Push on *" + j.Repository.Name + "* by *" + j.User_name + "* at *" + dateString + "* on branch *" + j.Ref + "*:" + n // First line
			message += commitCount + " commits :"                                                                                            // Second line
			for i := range j.Commits {
				c := j.Commits[i]
				message += n + "< " + c.Url + " | " + c.Id[0:7] + " >: " + "_" + MessageEncode(c.Message) + "_"
			}
		} else {
			message += "[PUSH] " + n + "Push on *" + j.Repository.Name + "* by *" + j.User_name + "* at *" + dateString + "* on branch *" + j.Ref + "*:" + n // First line
			message += "Last commit : < " + lastCommit.Url + " | " + lastCommit.Id + " > :" + n                                                              // Second line
			message += "```" + MessageEncode(lastCommit.Message) + "```"                                                                                     // Third line (last commit message)
		}
		SendWorkchatMessage(ThreadGitlab, message, ChatType)
	}
}

/*
	Handler function to handle http requests for merge

	@param body string
*/
func MergeHandler(body string) {
	var j data.Merge
	var err error           // Error catching
	var message string = "" // Bot's message
	var date time.Time      // Time of the last commit

	// Parse json and put it in a the data.Build structure
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		// Error
		l.Error("Error : Json parser failed :", err)
	} else {
		// Ok
		// Debug information
		if Verbose {
			l.Debug("JsonObject =", j)
		}

		// Send the message

		// Date parsing (parsing result example : 18 November 2014 - 14:34)
		date, err = time.Parse("2006-01-02 15:04:05 UTC", j.Object_attributes.Created_at)
		var dateString = date.Format("02 Jan 06 15:04")

		// Message
		message += "[MERGE REQUEST " + strings.ToUpper(j.Object_attributes.State) + "] " + n + "Target : *" + j.Object_attributes.Target.Name + "/" + j.Object_attributes.Target_branch + "* Source : *" + j.Object_attributes.Source.Name + "/" + j.Object_attributes.Source_branch + "* at *" + dateString + "* " + n // First line
		message += "Description: " + MessageEncode(j.Object_attributes.Description)                                                                                                                                                                                                                                     // Third line (last commit message)
		SendWorkchatMessage(ThreadGitlab, message, ChatType)
	}
}

/*
	Handler function to handle http requests for build

	@param body string
*/
func BuildHandler(body string) {
	var j data.Build
	var err error           // Error catching
	var message string = "" // Bot's message
	var date time.Time      // Time of the last commit

	// Parse json and put it in a the data.Build structure
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		// Error
		l.Error("Error : Json parser failed :", err)
	} else {
		// Ok
		// Debug information
		if Verbose {
			l.Debug("JsonObject =", j)
		}

		// Test if the message is already sent
		if currentBuildID < j.Build_id {
			// Not sent
			currentBuildID = j.Build_id // Update current build ID

			// Send the message

			// Date parsing (parsing result example : 18 November 2014 - 14:34)
			date, err = time.Parse("2006-01-02T15:04:05Z07:00", j.Push_data.Commits[0].Timestamp)
			var dateString = strconv.Itoa(date.Day()) + " " + date.Month().String() + " " + strconv.Itoa(date.Year()) +
				" - " + strconv.Itoa(date.Hour()) + ":" + strconv.Itoa(date.Minute())

			// Message
			lastCommit := j.Push_data.Commits[len(j.Push_data.Commits)-1]
			message += "[BUILD] " + n + strings.ToUpper(j.Build_status) + " : Push on *" + j.Push_data.Repository.Name + "* by *" + j.Push_data.User_name + "* at *" + dateString + "* on branch *" + j.Ref + "*:" + n // First line
			message += "Last commit : <" + lastCommit.Url + "|" + lastCommit.Id + "> :" + n                                                                                                                            // Second line
			message += "```" + MessageEncode(lastCommit.Message) + "```"                                                                                                                                               // Third line (last commit message)
			SendWorkchatMessage(ThreadGitlab, message, ChatType)
		} else {
			// Already sent
			// Do nothing
		}
	}

}

/*
	Handler function to handle http requests for build

	@param body string
*/

func TaskHandler(body string) {
	var j data.Tptask
	var Task data.TuleapTask
	var err error           // Error catching
	var message string = "" // Bot's message
	var date time.Time      // Time of the last commit

	// Parse json and put it in a the data.Build structure
	payload := strings.Split(body, "payload=")
	parsedValue, _ := url.QueryUnescape(payload[1])

	err = json.Unmarshal([]byte(parsedValue), &j)
	if err != nil {
		// Error
		l.Error("Error : Json parser failed :", err)
	} else {
		// Ok
		// Debug information
		if Verbose {
			l.Debug("JsonObject =", j)
		}
		Task.Name = j.User.RealName
		for _, val := range j.Current.Values {
			switch val.Label {
			case "Task title":
				Task.TaskTitle = val.Value.(string)
			case "Status":
				Task.Status = val.VValues[0].Label
			case "Links":
				if len(val.Reverse_Links) > 0 {
					Task.ProjectURL = TuleapURL + "projects/" + string(val.Reverse_Links[0].Tracker.Project.ID)
					Task.ProjectName = val.Reverse_Links[0].Tracker.Project.Label
				}
			case "Artifact ID":
				Task.TaskID = strconv.FormatFloat(val.Value.(float64), 'f', 0, 64)
				Task.TrackerURL = TuleapURL + "plugins/tracker/?aid=" + Task.TaskID
			case "Submitted on":
				Task.SubmittedOn = val.Value.(string)
			case "Details":
				Task.Details = strip.StripTags(val.Value.(string))
			case "Type":
				Task.Type = val.VValues[0].Label
			}
		}
		for _, val := range j.Previous.Values {
			if val.Label == "Status" {
				Task.OldStatus = val.VValues[0].Label
			}
		}

		if Task.Status != Task.OldStatus {
			date, err = time.Parse(time.RFC3339, Task.SubmittedOn)
			var dateString = date.Format("02 Jan 06 15:04")

			// Message
			message += "Move Task *#" + Task.TaskID + "* [" + Task.TaskTitle + "] on Project *" + Task.ProjectName + "* by *" + Task.Name + "* at *" + dateString + "* from *" + Task.OldStatus + "* to *" + Task.Status + "*" + n // First line
			message += "Description: " + MessageEncode(Task.Details) + n                                                                                                                                                           // Third line (last commit message)
			message += "Task URL : " + Task.TrackerURL
			SendWorkchatMessage(ThreadTuleap, message, ChatType)
		}
	}
}

/*
	Handler function to handle http requests for appcenter

	@param body string
*/
func AppCenterHandler(body string) {
	var j data.AppCenter
	var err error           // Error catching
	var message string = "" // Bot's message
	var date time.Time      // Time of the last commit

	// Parse json and put it in a the data.Build structure
	err = json.Unmarshal([]byte(body), &j)
	if err != nil {
		// Error
		l.Error("Error : Json parser failed :", err)
	} else {
		// Ok
		// Debug information
		if Verbose {
			l.Debug("JsonObject =", j)
		}

		// Send the message

		// Date parsing (parsing result example : 18 November 2014 - 14:34)
		date, err = time.Parse(time.RFC3339, j.SentAt)
		var dateString = date.Format("02 Jan 06 15:04")

		// Message
		message = ""
		if j.Reason == "" {
			message += "*" + j.AppName + "* (" + j.OS + ") Branch *" + j.Branch + "*" + n // First line
			message += "Build *#" + j.BuildID + "* [" + j.BuildStatus + "] on " + dateString + n
			message += "URL: " + j.BuildLink // Third line (last commit message)
		} else {
			message += "*Crash!!! [" + j.AppDisplayName + "]*" + n // First line
			message += "Reason *" + j.Reason + "* [" + j.Name + "] on " + dateString + n
			message += "URL: " + j.Url // Third line (last commit message)
		}

		SendWorkchatMessage(ThreadAppCenter, message, ChatType)
	}
}

/*
	Main function
*/
func main() {
	flag.Parse()                                             // Parse flags
	l.AddTransport(logo.Console).AddColor(logo.ConsoleColor) // Configure Logger
	l.EnableAllLevels()                                      // Configure Logger
	LoadConf()                                               // Load configuration
	l.Info(BotStartMessage)                                  // Logging
	http.ListenAndServe(":"+Port, &GitlabServ{})             // Run HTTP server for push hook
}
