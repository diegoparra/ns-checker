package controllers

import (
	"fmt"
	"strings"

	"github.com/diegoparra/ns-checker/pkg/acm"
	"github.com/diegoparra/ns-checker/pkg/r53"
	"github.com/diegoparra/ns-checker/pkg/utils"
)

type domain struct {
	Domain            string
	DesiredNameServer []string
	CurrentNameServer []string
	HasFacebookCode   bool
}

var myDomain domain

var withFb int

func CheckNameServer() {

	d, err := acm.ListCertificate()
	if err != nil {
		fmt.Println(err)
	}

	// Go over Domain list
	for _, v := range d {

		myDomain.Domain = v

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
		myDomain.DesiredNameServer, err = r53.GetRecordType(strings.TrimPrefix(*id.Id, "/hostedzone/"), v, "NS")
		if err != nil {
			fmt.Println(err)
		}

		// TODO
		// We should fix this TXT record as it does not work as intended
		// Get Desired TXT
		desiredTXT, err := r53.GetRecordType(strings.TrimPrefix(*id.Id, "/hostedzone/"), v, "TXT")
		if err != nil {
			fmt.Println(err)
		}

		if len(desiredTXT) <= 0 {
			myDomain.HasFacebookCode = false
		} else {
			myDomain.HasFacebookCode = true
			withFb += 1
		}

		// fmt.Println("Printing HasFacebookCode ")
		// fmt.Println(v)
		// fmt.Println("HasFacebookCode: ", myDomain.HasFacebookCode)

		// Validate Desired NS with Current NS
		val := utils.Validate(myDomain.DesiredNameServer, myDomain.CurrentNameServer)
		if val != true {
			fmt.Println("")
			fmt.Println("Status: Error validating NS")
			fmt.Println("Domain: ", v)
			fmt.Println("Current NS: ", myDomain.CurrentNameServer)
			fmt.Println("Desired NS: ", myDomain.DesiredNameServer)
			fmt.Println("")
		}

	}

	// Show number of analyzed domains
	fmt.Println("Number of analyzed Domains: ", len(d))
	fmt.Println("Number of domains with FB code: ", withFb)
	fmt.Println("Number of domains without FB code: ", len(d)-withFb)
}
