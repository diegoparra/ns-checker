package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/diegoparra/ns-checker/pkg/acm"
	"github.com/diegoparra/ns-checker/pkg/r53"
	"github.com/diegoparra/ns-checker/pkg/utils"
)

type domain struct {
	Url               string
	DesiredNameServer []string
	CurrentNameServer []string
	HasFacebookCode   bool
}

var (
	myDomain domain
	withFb   int
)

func CheckNameServer() {

	d, err := acm.ListCertificate()
	if err != nil {
		fmt.Println(err)
	}

	// Go over Domain list
	for _, v := range d {

		myDomain.Url = v

		//Get HostedZoneID
		id, err := r53.GetHostedZoneID(v)
		if err != nil {
			fmt.Println(err)
		}

		// Get Current NS
		myDomain.CurrentNameServer, err = r53.GetNS(v)
		if err != nil {
			fmt.Println(err)
		}

		// Get Desired NS
		myDomain.DesiredNameServer, err = r53.GetRecordByType(strings.TrimPrefix(*id.Id, "/hostedzone/"), v, "NS")
		if err != nil {
			fmt.Println(err)
		}

		// Get Desired TXT
		txt, err := r53.GetRecordByType(strings.TrimPrefix(*id.Id, "/hostedzone/"), v, "TXT")
		if err != nil {
			fmt.Println(err)
		}

		if len(txt) <= 0 {
			myDomain.HasFacebookCode = false
		} else {
			myDomain.HasFacebookCode = true
			withFb += 1
		}

		// Validate Desired NS against Current NS
		val := utils.Validate(myDomain.DesiredNameServer, myDomain.CurrentNameServer)
		if val != true {
			fmt.Println("")
			fmt.Println("Status: Error validating NS")
			fmt.Println("Domain: ", myDomain.Url)
			fmt.Println("Current NS: ", myDomain.CurrentNameServer)
			fmt.Println("Desired NS: ", myDomain.DesiredNameServer)
			fmt.Println("")
		}

		j, _ := json.Marshal(myDomain)
		fmt.Println(string(j))
	}

	// Show number of analyzed domains
	fmt.Println("Number of analyzed Domains: ", len(d))
	fmt.Println("Number of domains with FB code: ", withFb)
	fmt.Println("Number of domains without FB code: ", len(d)-withFb)
}
