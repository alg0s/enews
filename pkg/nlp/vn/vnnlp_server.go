package vn

/*
   TODO:
        1. Check if VnCoreNLP file exists
        2. Check if Java exists
        3. Check if VnCoreNLPServer exists
        4. Add retry in case multiple clients - 1 server
*/

/*
   Shortcuts:
	lsof -t -i :9000 -s TCP:LISTEN -c java

	// java -Xmx2g -jar /Users/steve/Documents/enews/pkg/nlp/vn/vncorenlp/VnCoreNLPServer.jar /Users/steve/Documents/enews/pkg/nlp/vn/vncorenlp/VnCoreNLP-1.1.1.jar -i 127.0.0.1 -p 9000 -a wseg,pos,ner,parse
*/

import (
	"io/ioutil"
	"log"
	"net/http"
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
	// serverJarFile       string = `pkg/nlp/vn/vncorenlp/VnCoreNLPServer.jar`
	serverJarFile       string = `pkg/nlp/vn/vncorenlp/VnNLPServer.jar`
	nlpJarFile          string = `pkg/nlp/vn/vncorenlp/VnCoreNLP-1.1.1.jar`
	defaultMaxURILength int    = 8192
)

// NLPServer handles VnCoreNLP, acting as a client, but also has ability to launch the server
// providing access to the NLP annotators
type NLPServer struct {
	Process     *os.Process
	Address     string
	Port        string
	Host        string
	MaxHeapSize string
	Annotators  string
	isAlive     bool
}

// isInit checks if an instance of NLPServer has been initated with required element values
func (s *NLPServer) isInit() bool {
	if s.Host == `` || s.Port == `` || s.MaxHeapSize == `` || s.Annotators == `` {
		return false
	}
	return true
}

// findExistingServerProcess finds the process of the running server
func (s *NLPServer) findExistingServerProcess() (*os.Process, error) {
	// 1. Find PID
	var args = []string{
		`-t`,
		`-i`, `:` + s.Port,
		`-s`, `TCP:LISTEN`,
		`-c`, `java`,
	}
	pid, err := exec.Command(`lsof`, args...).CombinedOutput()
	log.Println(">> PID: ", string(pid))
	if err != nil {
		log.Println("Err cmd: ", err, args)
		return nil, err
	}

	// 2. FindProcess
	intPid, err := strconv.Atoi(strings.TrimSpace(string(pid)))
	log.Println(">> intPID: ", intPid)
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

// getAnnotators gets the annotators registered with the Server
func (s *NLPServer) getAnnotators() string {
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

// KillExistingServer kills the existing process running the NLPServer
func (s *NLPServer) killExistingServer() (bool, error) {
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

func (s *NLPServer) startJavaServer() (bool, error) {
	log.Println(`Starting new VnNLP server...`)

	// 1. Get Jar files
	serverJarFilePath, err := filepath.Abs(serverJarFile)
	if err != nil {
		log.Fatal("Unable to find NLPServer Jar file", err)
	}

	nlpJarFilePath, err := filepath.Abs(nlpJarFile)
	if err != nil {
		log.Fatal("Unable to find NLP Jar file", err)
	}

	// 2. Prepare command arguments
	var args = []string{
		s.MaxHeapSize,
		`-jar`, serverJarFilePath, nlpJarFilePath,
		`-i`, s.Host,
		`-p`, s.Port,
		`-a`, s.Annotators,
	}
	log.Println("Cmd to start server: ", "java ", args)
	cmd := exec.Command(`java`, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 3. Start Server with command line
	err = cmd.Start()
	if err != nil {
		log.Fatal(`Unable To Start NLPServer: `, err)
		return false, err
	}

	serverStarted := false

	for serverStarted == false {
		log.Println(">>> Server is still starting, sleep 5s...")
		// sleep for 5s to let the server get ready
		time.Sleep(time.Second * 5)

		// check if server started
		serverStarted = s.IsServerAlive()
	}

	s.Process = cmd.Process
	s.isAlive = true
	return true, nil
}

// Start launches NLPServer and returns true if successful, otherwise false
func (s *NLPServer) Start() bool {

	if !s.isInit() {
		log.Fatal("NLPServer is not initiated yet!")
	}

	s.Address = strings.Join([]string{`http://`, s.Host, `:`, s.Port}, ``)

	// 1. check if a server is running on that host
	resp, err := http.Get(s.Address)
	if err != nil {
		log.Println("Server Not Active: ", err)
	} else {
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Println(`Server is already running on port: `, s.Port)

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

	// 2. if no server is running, init new server
	serverStarted, err := s.startJavaServer()
	if err != nil {
		log.Panic("Unable to start VnNLP Java server: ", err)
		return false
	}
	return serverStarted
}

// Stop kills the process currently running NLPServer and returns true if successful, otherwise false
func (s *NLPServer) Stop() (bool, error) {
	result, err := s.killExistingServer()
	if err != nil {
		return false, err
	}
	return result, nil
}

// GetInfo provides information about the Server
func (s *NLPServer) GetInfo() ServerInfo {
	return ServerInfo{
		Host:       s.Host,
		Port:       s.Port,
		Address:    s.Address,
		Annotators: s.Annotators,
	}
}

// IsProcessAlive returns true if NLPServer is alive, otherwise false
func (s *NLPServer) IsProcessAlive() bool {
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
func (s *NLPServer) IsServerAlive() bool {
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

// NewNLPServer initiates a new NLPServer instance
func NewNLPServer(host string, port string, annotators string, maxHeapSize string) *NLPServer {
	var s = NLPServer{
		Host:        host,
		Port:        port,
		Annotators:  annotators,
		MaxHeapSize: maxHeapSize,
	}
	s.Start()
	return &s
}
