package r53

import (
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func GetRecordType(zoneID, zoneName, recordType string) ([]string, error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil, err
	}

	svc := route53.New(sess)

	resp, err := svc.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneID),
		StartRecordName: aws.String(zoneName),
		StartRecordType: aws.String(recordType),
	})
	if err != nil {
		fmt.Println("Failed to get the resources list", err)
		return nil, err
	}
	if len(resp.ResourceRecordSets) == 0 {
		fmt.Println("Nothing found")
	}

	// var ns []string

	// fmt.Println(resp)

	// Current working code
	// ns := make([]string, len(resp.ResourceRecordSets[0].ResourceRecords))
	// for i := range resp.ResourceRecordSets[0].ResourceRecords {
	// 	ns[i] = strings.TrimSuffix(*resp.ResourceRecordSets[0].ResourceRecords[i].Value, ".")
	// }

	// Testing new code
	ns := []string{}

	for i := range resp.ResourceRecordSets[0].ResourceRecords {
		ns = append(ns, strings.TrimSuffix(*resp.ResourceRecordSets[0].ResourceRecords[i].Value, "."))
	}

	if recordType == "TXT" {

		found, err := regexp.MatchString(`(?:facebook-domain-verification)`, string(ns[0]))
		if err != nil {
			fmt.Println("Error to run regex match string", err)
		}

		if found {
			return ns, nil
		} else {
			return []string{}, nil

		}
	}

	sort.Strings(ns)

	return ns, nil
}

func GetNS(d string) ([]string, error) {
	var nss []string
	n, _ := net.LookupNS(d)

	for _, v := range n {
		a := strings.TrimSuffix(v.Host, ".")
		nss = append(nss, a)
	}

	return nss, nil
}

func GetHostedZoneID(d string) (*route53.HostedZone, error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil, err
	}

	svc := route53.New(sess)

	params := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(d),
		MaxItems: aws.String("1"),
	}

	resp, err := svc.ListHostedZonesByName(params)
	if err != nil {
		return nil, err
	}

	zone := resp.HostedZones[0]
	return zone, nil

}
