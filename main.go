package main

import (
	"enews/pkg/configs"
	"enews/pkg/nlp/vn"
	"log"
)

func runNLPServer(start, stop, reset chan bool) {
	// Load VnNLP configs
	st := configs.LoadConfigs().VnNLP

	// Start VnNLP server
	var s = vn.NLPServer{
		Host:        st.Host,
		Port:        st.Port,
		Annotators:  st.Annotators,
		MaxHeapSize: st.MaxHeapSize,
	}

	// Use Select to determine when to stop the server (i.e when the program ends)
	for {
		select {
		case <-start:

			log.Println(`Starting NLP Server: `, s.Host)
			started := s.Start()

			// Let main know that the server has been started successfully
			if started == true {
				start <- started
			}

		case <-stop:
			log.Println(`Stopping NLP Server`)
			var stopped bool

			if st.VnNLPConfigs.StopAfterProgramQuit == true {
				stopped, _ = s.Stop()
			} else {
				stopped = true
			}
			stop <- stopped

		case <-reset:
			log.Println(`Restarting NLP server...`)
			s.Stop()
			restarted := s.Start()
			// Let main know that the server has been restarted successfully
			if restarted == true {
				log.Println("NLP Server restarted successfully")
				reset <- restarted
			}
		}
	}
}

func work(extract, done chan bool) {
	if <-extract {
		log.Println("Start working....")
		vn.RunExtractPipeline()
		done <- true
	}
}

// main triggers and coordinates core components concurrently
func main() {
	var start = make(chan bool)
	var stop = make(chan bool)
	var reset = make(chan bool)
	var extract = make(chan bool)
	var done = make(chan bool)

	go runNLPServer(start, stop, reset)
	go work(extract, done)

	start <- true

	for {
		select {
		case started := <-start:
			if started == true {
				extract <- started
			} else {
				log.Println("Unable to start NLP server. Shutdown.")
				return
			}
		case finished := <-done:
			log.Println("Work has finished: ", finished)
			stop <- true
		case issues := <-extract:
			log.Println("Work has server issue: ", issues)
			reset <- true
		case <-stop:
			log.Println("Shutdown.")
			return
		}
	}
}
