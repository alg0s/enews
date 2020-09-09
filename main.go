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
	for {
		select {
		case <-extract:
			log.Println("Start working....")

			switch vn.RunExtractPipeline() {
			case false:
				extract <- false
			case true:
				done <- true
			}
		case <-done:
			return
		}
	}
}

// main triggers and coordinates core components concurrently
func main() {
	var start = make(chan bool)
	var stop = make(chan bool)
	var resetServer = make(chan bool)
	var extract = make(chan bool)
	var done = make(chan bool)

	go runNLPServer(start, stop, resetServer)
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
		case issues := <-extract:
			log.Println("Work has server issue: ", issues)
			resetServer <- true
		case <-resetServer:
			log.Println("Reset successfully. Restart the pipeline...")
			extract <- true
		case finished := <-done:
			log.Println("Work has finished: ", finished)
			done <- true
			stop <- true
		case <-stop:
			log.Println("Shutdown.")
			return
		}
	}
}
