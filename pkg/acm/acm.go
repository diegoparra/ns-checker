package acm

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
)

var domains []string

func ListCertificate() ([]string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		fmt.Println("failed to create session,", err)
		return nil, err
	}

	svc := acm.New(sess)

	params := &acm.ListCertificatesInput{
		CertificateStatuses: []*string{
			aws.String("ISSUED"), // Required
			aws.String("VALIDATION_TIMED_OUT"),
		},
		MaxItems: aws.Int64(500),
	}

	resp, err := svc.ListCertificates(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	for _, v := range resp.CertificateSummaryList {
		found, err := regexp.MatchString(`(?:linkfire|metafire|lnktest|lnk|mpo|linkfire-partners|smarturl|bio|server|amazonaws|sonymusic|lexijayde)\.(?:co|com|dk|it|tt|to|link)`, string(*v.DomainName))
		// server|amazonaws.com|linkfire-partners.com|smarturl.it|linkfire.tobis.de|*.bio.to|lexijayde.link|ads.sonymusic.com.br)`, string(*v.DomainName))
		// found, err := regexp.MatchString(`(?:([-a-z0-9]+)\.linkfire\.co|metafire.co|linkfire.com|lnk.to|lnktest.dk|mit.mpo.dk|server|lnk.tt|amazonaws.com|linkfire-partners.com|smarturl.it|linkfire.tobis.de|*.bio.to|lexijayde.link|ads.sonymusic.com.br)`, string(*v.DomainName))
		if err != nil {
			fmt.Println("Error to run regex match string", err)
		}

		if !found {
			domains = append(domains, *v.DomainName)
		}
	}

	return domains, nil
}
