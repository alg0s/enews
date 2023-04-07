package main

import (
	"enews/pkg/configs"
	"enews/pkg/db"
	"enews/pkg/nlp/vn"
	"log"
)

func runNLPServer(start, stop, reset chan bool) {
	// Load VnNLP configs
	conf := configs.LoadConfigs().VnNLP

	// Start VnNLP server
	var s = vn.NLPServer{
		Host:        conf.Host,
		Port:        conf.Port,
		Annotators:  conf.Annotators,
		MaxHeapSize: conf.MaxHeapSize,
	}

	// Use Select to determine when to stop the server (i.e when the program ends)
	for {
		select {
		case <-start:

			log.Println(`Starting NLP Server: `, s.Host)
			// started := s.Start()

			// Let main know that the server has been started successfully
			// if started == true {
			// 	start <- started
			// }
			start <- true

		case <-stop:
			log.Println(`Stopping NLP Server`)
			var stopped bool

			if conf.VnNLPConfigs.StopAfterProgramQuit == true {
				stopped, err := s.Stop()
				if err != nil {
					log.Println("Unable to stop VnNLP Server: ", err)
				}
				log.Println("Server stopped? ", stopped)
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

func work(db *db.DB, extract, done chan bool) {
	for {
		select {
		case <-extract:
			log.Println("Start working....")
			vn.RunExtractor(db)
			done <- true
		case <-done:
			return
		}
	}
}

// main triggers and coordinates core components concurrently
func main() {
	// Start a connection with the database
	dbconn, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	//
	var start = make(chan bool)
	var stop = make(chan bool)
	var resetServer = make(chan bool)
	var extract = make(chan bool)
	var done = make(chan bool)

	go runNLPServer(start, stop, resetServer)
	go work(dbconn, extract, done)

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
			done <- true
			stop <- true
		case <-stop:
			log.Println("Shutdown.")
			return
		}
	}
}
