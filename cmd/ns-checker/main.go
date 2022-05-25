package main

import (
	"fmt"
	"strings"

	"github.com/diegoparra/ns-checker/pkg/acm"
	"github.com/diegoparra/ns-checker/pkg/r53"
	"github.com/diegoparra/ns-checker/pkg/utils"
)

func main() {

	d, err := acm.ListCertificate()
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range d {
		// Go over Domain list
		fmt.Println("")
		fmt.Println("Domain: ", v)

		//Get HostedZoneID
		id, err := r53.GetHostedZoneID(v)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("HostedZoneID: ", strings.TrimPrefix(*id.Id, "/hostedzone/"))

		// Get Current NS
		ns, err := r53.GetNS(v)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Current NS: ", ns)

		// Get Desired NS
		desiredNS, err := r53.GetRecordType(strings.TrimPrefix(*id.Id, "/hostedzone/"), v, "NS")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Desired NS: ", desiredNS)

		// Validate Desired NS with Current NS
		val := utils.Validate(desiredNS, ns)
		if val != true {
			fmt.Println("Domain: " + v + " not poiting to Linkfire NS's")
		} else {
			fmt.Println("Domain: " + v + " NS's working fine")
		}

		// TODO
		// We should fix this TXT record as it does not work as intended
		// Get Desired TXT
		// desiredTXT, err := r53.GetRecordType(strings.TrimPrefix(*id.Id, "/hostedzone/"), v, "TXT")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println("Current TXT: ", desiredTXT)
		// fmt.Println("")
	}

	// Show number of analyzed domains
	fmt.Println("Number of analyzed Domains: ", len(d))

}
