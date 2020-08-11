package vn

/*
   TODO:
        1. Check if VnCoreNLP file exists
        2. Check if Java exists
        3. Check if VnCoreNLPServer exists
        4. Add retry in case multiple clients - 1 server
*/

/*
   REF:
	https://www.digitalocean.com/community/tutorials/understanding-data-types-in-go
	https://deployeveryday.com/2019/10/21/python-idioms-go.html
*/

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	// Server Defaults
	defaultMaxHeapSize string = `-Xmx2g`
	defaultHost        string = `127.0.0.1`
	defaultPort        string = `9000`
	defaultAnnotators  string = `wseg,pos,ner,parse`
	serverJarFile      string = `pkg/nlp/vn/vncorenlp/VnCoreNLPServer.jar`
	nlpJarFile         string = `pkg/nlp/vn/vncorenlp/VnCoreNLP-1.1.1.jar`
)

// Token details of an extract entity
type Token struct {
	Index    int    `json:"index"`
	Form     string `json:"form"`
	PosTag   string `json:"posTag"`
	NerLabel string `json:"nerLabel"`
	Head     int    `json:"head"`
	DepLabel string `json:"depLabel"`
}

// Sentence is an array of Token's
type Sentence []Token

// ServerResponse is the response from VnNLPServer
type ServerResponse struct {
	Sentences []Sentence `json:"sentences"`
	Status    bool       `json:"status"`
	Error     string     `json:"error"`
	Language  string     `json:"language"`
}

// VnNLPServer handles VnCoreNLP, acting as a client, but also has ability to launch the server
// providing access to the NLP annotators
type VnNLPServer struct {
	Address     string
	Process     *os.Process
	isAlive     bool
	Port        string
	Host        string
	MaxHeapSize string
	Annotators  string
}

// IsProcessAlive returns true if VnNlpServer is alive, otherwise false
func (s *VnNLPServer) IsProcessAlive() bool {
	process, err := os.FindProcess(s.Process.Pid)
	if err != nil {
		log.Printf(`Failed to find process: %s\n`, err)
	} else {
		err := process.Signal(syscall.Signal(0))
		log.Printf(`process.Signal on pid %d returned: %v\n`, s.Process.Pid, err)
	}
	return false
}

// IsServerAlive returns true if a server is running, otherwise false
func (s *VnNLPServer) IsServerAlive() bool {
	var address = s.Address
	if address == "" {
		log.Panic("Missing Server Address")
		return false
	}

	resp, err := http.Get(address)
	if err != nil {
		return false
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}

// Start launches VnNLPServer and returns true if successful, otherwise false
func (s *VnNLPServer) Start(host string, port string, maxHeapSize string, annotators string) bool {
	// Assign variables
	if host == "" {
		host = defaultHost
	}
	if port == "" {
		port = defaultPort
	}
	if maxHeapSize == "" {
		maxHeapSize = defaultMaxHeapSize
	}
	if annotators == "" {
		annotators = defaultAnnotators
	}

	// 1. check if a server is running on that host
	var address = strings.Join([]string{`http://`, host, `:`, port}, ``)
	log.Println("Address of Server is: ", address)

	resp, err := http.Get(address)

	// assign constant values
	s.Address = address
	s.Port = port
	s.Host = host
	s.MaxHeapSize = maxHeapSize

	if err != nil {
		log.Println("Server Not Active: ", err)
	} else {
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Println(`Server is already running on port: `, port)

			// assign variables
			s.isAlive = true
			s.Annotators = s.getAnnotators()

			// Get the existing server process
			s.Process, err = s.findExistingServerProcess()
			if err != nil {
				log.Panic("Unable to find existing server process: ", err)
			}
			return true
		}
	}

	// 2. if no server is running, init new
	// TODO: utilize goroutine to initiate the server

	log.Println(`Starting server...`)

	serverJarFilePath, err := filepath.Abs(serverJarFile)

	if err != nil {
		log.Fatal("Unable to find VnNLPServer Jar file", err)
	}

	nlpJarFilePath, err := filepath.Abs(nlpJarFile)

	if err != nil {
		log.Fatal("Unable to find NLP Jar file", err)
	}

	var args = []string{
		maxHeapSize,
		`-jar`, serverJarFilePath, nlpJarFilePath,
		`-i`, host,
		`-p`, port,
		`-a`, annotators,
	}

	cmd := exec.Command(`java`, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()

	if err != nil {
		log.Fatal(`Unable To Start VnNLPServer: `, err)
		return false
	}
	log.Printf(`Started VnCoreNLP Server, process: %d`, cmd.Process.Pid)

	// sleep for 5s to let the server get ready
	time.Sleep(time.Second * 5)
	s.Process = cmd.Process
	s.isAlive = true
	return true
}

func (s *VnNLPServer) findExistingServerProcess() (*os.Process, error) {

	// get server Port

	// 1. Find PID
	var args = []string{
		`-t`,
		`-i`, `:` + s.Port,
		`-s`, `TCP:LISTEN`,
		`-c`, `java`,
	}
	pid, err := exec.Command(`lsof`, args...).CombinedOutput()
	if err != nil {
		log.Println("Err cmd: ", err, args)
		return nil, err
	}

	// 2. FindProcess
	intPid, err := strconv.Atoi(strings.TrimSpace(string(pid)))
	if err != nil {
		log.Println("Err find trim space: ", err, intPid)
		return nil, err
	}

	process, err := os.FindProcess(intPid)
	if err != nil {
		log.Println("Err find proc: ", err, process)
		return nil, err
	}
	return process, nil
}

// KillExistingServer kills the existing process running the VnNLPServer
func (s *VnNLPServer) killExistingServer() (bool, error) {
	// 1. find Process
	process, err := s.findExistingServerProcess()
	if err != nil {
		return false, err
	}

	// 2. kill Process
	err = process.Kill()
	if err != nil {
		log.Panic("Unable To Kill: ", err)
		return false, err
	}
	return true, nil
}

// Stop kills the process currently running VnNLPServer and returns true if successful, otherwise false
func (s *VnNLPServer) Stop() (bool, error) {
	log.Println(">>> About to stop Process: ", s.Process.Pid)
	result, err := s.killExistingServer()
	if err != nil {
		return false, err
	}
	return result, nil
}

// getAnnotators gets the annotators registered with the Server
func (s *VnNLPServer) getAnnotators() string {
	resp, err := http.Get(strings.Join([]string{s.Address, `/annotators`}, ``))
	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// process annotators
	var a string
	a = string(bodyBytes)
	a = strings.ReplaceAll(a, `"`, ``)
	a = strings.ReplaceAll(a, `[`, ``)
	a = strings.ReplaceAll(a, `]`, ``)
	return a
}

// Annotate sends a request to the server to ask for an annotation of a string and returns the response
func (s *VnNLPServer) Annotate(textInput string, annotators string) *ServerResponse {
	if annotators == "" {
		annotators = s.Annotators
	}

	// construct URL
	urlStr := strings.Join([]string{s.Address, `/handle`}, ``)
	u, _ := url.Parse(urlStr)
	query, _ := url.ParseQuery(u.RawQuery)
	query.Add(`text`, textInput)
	query.Add(`props`, annotators)
	u.RawQuery = query.Encode()

	// set content type
	contentType := `text/plain`

	resp, err := http.Post(u.String(), contentType, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	var sr ServerResponse

	if resp.StatusCode == http.StatusOK {
		// bodyBytes, err := ioutil.ReadAll((resp.Body))
		err = json.NewDecoder(resp.Body).Decode(&sr)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		return &sr
	}
	return nil
}

func (s *VnNLPServer) customAnnotate(text string, annotators string) *[]Sentence {
	var result = s.Annotate(text, annotators)
	if result.Error != "" {
		return nil
	}
	return &result.Sentences
}

// Tokenize return tokens from input string, otherwise empty string
func (s *VnNLPServer) Tokenize(text string) *[]Sentence {
	return s.customAnnotate(text, "wseg")
}

// PosTag returns POS tags from the input string, otherwise empty string
func (s *VnNLPServer) PosTag(text string) *[]Sentence {
	return s.customAnnotate(text, "wseg,pos")
}

// Ner returns NER - Named Entity Recognition from the input string, otherwise empty string
func (s *VnNLPServer) Ner(text string) *[]Sentence {
	return s.customAnnotate(text, "wseg,pos,ner")
}

// DepParse returns parsed dependencies from input string, otherwise empty string
func (s *VnNLPServer) DepParse(text string) *[]Sentence {
	return s.customAnnotate(text, "wseg,pos,ner,parse")
}

// DetectLanguage returns the detected language of the input string
func (s *VnNLPServer) DetectLanguage(text string) string {
	var result = s.Annotate(text, "lang")
	if result.Error != "" {
		return ""
	}
	return result.Language
}

// NewVnNLPServer initiates a new VnNLPServer instance
func NewVnNLPServer() *VnNLPServer {
	var s = VnNLPServer{}
	result := s.Start("", "", "", "")
	log.Println("result start: ", result)
	return &s
}

// TODO: problems arised:
// 1. Process didn't sleep to wait for startup
// 2. Killed process raises errors
