package main

/*
   Ref:
   https://stackoverflow.com/questions/48557810/execute-command-in-background
*/

import (
	"log"
	"os"
	"os/exec"
)

// const (
// 	maxHeapSize   = "-Xmx2g"
// 	serverJarFile = "lib/vncorenlp/VnCoreNLPServer.jar"
// 	nlpJarFile    = "lib/vncorenlp/VnCoreNLP-1.1.1.jar"
// 	endpoint      = "127.0.0.1"
// 	port          = "9000"
// 	annotators    = "wseg,pos,ner,parse"
// )

// StartVnNLPServer handles VnCoreNLP server and return the Process
// func StartVnNLPServer() *os.Process {

// 	var cmdargs = []string{
// 		maxHeapSize,
// 		`-jar`, serverJarFile, nlpJarFile,
// 		`-i`, endpoint,
// 		`-p`, port,
// 		`-a`, annotators,
// 	}

// 	cmd := exec.Command("java", cmdargs...)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	// cmd.Run()

// 	err := cmd.Start()

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)

// 	return cmd.Process
// }

// func main() {
// 	fmt.Println("Enews is browsing news with entities.")
// 	StartVnNLPServer()
// 	for i := 0; i < 10; i++ {
// 		fmt.Printf("Counting %s...", string(i))
// 	}
// }
