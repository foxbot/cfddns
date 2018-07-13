package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

type ipResponse struct {
	Origin string `json:"origin"`
}

var (
	apiKey    string
	apiEmail  string
	domain    string
	subdomain string
	httpBin   string
)

func init() {
	flag.StringVar(&apiKey, "key", "", "-key=<api token>")
	flag.StringVar(&apiEmail, "email", "", "-email=<api email>")
	flag.StringVar(&domain, "domain", "", "-domain=example.com")
	flag.StringVar(&subdomain, "subdomain", "", "-subdomain=remote.example.com")
	flag.StringVar(&httpBin, "httpbin", "https://httpbin.org/ip", "-httpbin=https://httpbin.org/ip")
}

func main() {
	flag.Parse()
	if flag.NFlag() < 4 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Println("cfddns >start")

	err := run()
	if err != nil {
		log.Panicln(err)
	}

	log.Println("cfddns >done")
}

func run() error {
	r, err := http.Get(httpBin)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var ipResp = new(ipResponse)
	err = json.Unmarshal(buf, &ipResp)
	if err != nil {
		return err
	}
	ip := ipResp.Origin
	log.Println("acquired public IP: ", ip)

	api, err := cloudflare.New(apiKey, apiEmail)
	if err != nil {
		return err
	}

	zoneID, err := api.ZoneIDByName(domain)
	if err != nil {
		return err
	}

	recordParameter := cloudflare.DNSRecord{
		Name: subdomain,
	}
	records, err := api.DNSRecords(zoneID, recordParameter)
	if len(records) != 1 {
		return errors.New("got too many records back matching subdomain")
	}

	record := records[0]
	record.Content = ip

	err = api.UpdateDNSRecord(zoneID, record.ID, record)
	if err != nil {
		return err
	}

	return nil
}
