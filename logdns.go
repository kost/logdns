package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/miekg/dns"
)

var ip string
var VerboseLevel bool
var Port int
var myttl string
var logfilename string

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		if VerboseLevel {
			log.Printf("Query for %s: %v\n", q.Name, q.Qtype)
		}
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			rr, err := dns.NewRR(fmt.Sprintf("%s %s A %s", q.Name, myttl, ip))
			if err == nil {
				m.Answer = append(m.Answer, rr)
				if VerboseLevel {
					log.Printf("Query for %s failed: %v\n", q.Name, err)
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	if VerboseLevel {
		log.Printf("Query for \n")
	}
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func main() {
	returnip := flag.String("return", "127.0.0.1", "what address to return")
	ttl := flag.String("ttl", "3600", "Set a custom ttl for returned records")
	listen := flag.String("listen", "0.0.0.0", "listen address")
	resolve := flag.String("resolve", ".", "which domains to respond, e.g. service.")
	verbose := flag.Bool("verbose", false, "be verbose")
	flag.IntVar(&Port, "port", 53, "port to listen")
	flag.StringVar(&logfilename, "logfile", "", "Log File name")

	flag.Parse()

	VerboseLevel = *verbose
	ip = *returnip
	myttl = *ttl

	dns.HandleFunc(*resolve, handleDnsRequest)

	// start server
	listenstr := *listen + ":" + strconv.Itoa(Port)
	server := &dns.Server{Addr: listenstr, Net: "udp"}
	log.Printf("Starting at %s\n", listenstr)

	if len(logfilename) > 0 {
		file, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
	}
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server (check if you have permission to listen on specific port): %s\n ", err.Error())
	}
}
