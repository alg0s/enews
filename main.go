package main

/*
	ref:
	https://www.opsdash.com/blog/job-queues-in-go.html
	https://bigkevmcd.github.io/go/pgrp/context/2019/02/19/terminating-processes-in-go.html
*/

import (
	"enews/pkg/nlp/vn"
)

func closeNLPServer() {
	s := vn.NewVnNLPServer()
	s.Stop()
}

func main() {
	// extract.Extract()
	closeNLPServer()
}
