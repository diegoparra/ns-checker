package r53

import (
	"fmt"
	"net"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func GetAwsNS(zID string, zNAME string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		fmt.Println("failed to create session,", err)
		// return nil, err
	}

	svc := route53.New(sess)

	resp, err := svc.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zID),
		StartRecordName: aws.String(zNAME),
		StartRecordType: aws.String("NS"),
	})
	if err != nil {
		fmt.Println("Caiu no error do resp")
		fmt.Println(err)
		// return nil, err
	}
	if len(resp.ResourceRecordSets) == 0 {
		fmt.Println("Nothing found")
		// return nil, nil
	}
	ns := make([]string, len(resp.ResourceRecordSets[0].ResourceRecords))
	for i := range resp.ResourceRecordSets[0].ResourceRecords {
		ns[i] = *resp.ResourceRecordSets[0].ResourceRecords[i].Value
	}
	sort.Strings(ns)
	fmt.Println(ns)
}

func GetNS(d string) {
	var nss []string
	n, _ := net.LookupNS(d)

	for _, v := range n {
		nss = append(nss, v.Host)
	}

	// out, err := json.Marshal(n)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// for k, v := range out {
	// 	fmt.Println(k)
	// 	fmt.Println(v)
	// 	nss = append(nss)
	// }

	fmt.Println(nss)
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
