package main

import (
	"fmt"
	"strings"

	"github.com/diegoparra/scd/pkg/acm"
	"github.com/diegoparra/scd/pkg/r53"
)

func main() {

	acm.ListCertificate()

	// fmt.Println(domains)
	for _, v := range acm.Domains {
		fmt.Println("")
		fmt.Println("Domain: ", v)
		fmt.Println("")
		fmt.Println("Current NS: ")

		r53.GetNS(v)

		id, err := r53.GetHostedZoneID(v)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("")
		fmt.Println("HostedZoneID: ", strings.TrimPrefix(*id.Id, "/hostedzone/"))

		fmt.Println("")
		fmt.Println("Desired NS: ")

		r53.GetAwsNS(strings.TrimPrefix(*id.Id, "/hostedzone/"), v)

		fmt.Println("")
		fmt.Println("End domain")
	}
}
