// web checks a list of URLs and outputs information about
// availability to a file
//
// Copyright (c) 2015 John Graham-Cumming

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	interval := flag.Duration("interval", 1*time.Second,
		"Interval between checks")
	output := flag.String("output", "", "Name of output file")
	flag.Parse()

	if *output == "" {
		log.Fatalf("The -output parameter cannot be empty")
	}

	urls := flag.Args()
	if len(urls) == 0 {
		log.Fatalf("The command-line must contain some URLs")
	}

	var client http.Client

	var reqs = make([]*http.Request, len(urls))
	for i, url := range urls {
		var err error
		if reqs[i], err = http.NewRequest("GET", url, nil); err != nil {
			log.Fatalf("Failed to create HTTP request: %s", err)
		}

		reqs[i].Header.Set("User-Agent",
			"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:42.0) Gecko/20100101 Firefox/42.0")
	}

	for {
		out, err := os.OpenFile(*output, os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0600)
		if err == nil {
			for _, req := range reqs {
				resp, err := client.Do(req)
				errorString := ""
				if err != nil {
					errorString = fmt.Sprintf("%s", err)
				}
				fmt.Fprintf(out, "%s,%s,%s,%s\n",
					time.Now().Format(time.RFC3339Nano), req.URL, resp.Status,
					errorString)
				resp.Body.Close()
			}
			out.Close()
		} else {
			log.Printf("Error opening output file %s: %s", *output, err)
		}

		time.Sleep(*interval)
	}
}
