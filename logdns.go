package main

import (
	"fmt"
	"log"
	"strconv"
	"flag"

	"github.com/miekg/dns"
)

var ip string
var VerboseLevel bool
var Port int

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		if VerboseLevel {
			log.Printf("Query for %s: %v\n", q.Name, q.Qtype)
		}
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
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
        listen := flag.String("listen", "0.0.0.0", "listen address")
	resolve := flag.String("resolve", ".", "which domains to respond, e.g. service.")
        verbose := flag.Bool("verbose", false, "be verbose")
        flag.IntVar(&Port, "port", 53, "port to listen")

        flag.Parse()

	VerboseLevel=*verbose
	ip=*returnip

	dns.HandleFunc(*resolve, handleDnsRequest)

	// start server
	listenstr := *listen + ":" + strconv.Itoa(Port)
	server := &dns.Server{Addr: listenstr, Net: "udp"}
	log.Printf("Starting at %s\n", listenstr)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server (check if you have permission to listen on specific port): %s\n ", err.Error())
	}
}
