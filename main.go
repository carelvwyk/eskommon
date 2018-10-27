package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	nut "github.com/robbiet480/go.nut"
)

var (
	nutHost     = flag.String("host", "127.0.0.1", "NUT server host")
	nutUsername = flag.String("username", "", "NUT user name")
	nutPassword = flag.String("password", "", "NUT user password")
)

func main() {
	flag.Parse()

	if *nutHost == "" || *nutUsername == "" || *nutPassword == "" {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}
	for {
		client, connectErr := nut.Connect(*nutHost)
		if connectErr != nil {
			log.Fatal(connectErr)
		}
		_, authenticationError := client.Authenticate(*nutUsername, *nutPassword)
		if authenticationError != nil {
			log.Fatal(authenticationError)
		}

		upsList, listErr := client.GetUPSList()
		if listErr != nil {
			log.Fatalf("GetUPSList: %v", listErr)
		}
		if len(upsList) < 1 {
			log.Fatal("No UPS found")
		}

		c, err := getBatteryCharge(upsList[0])
		if err != nil {
			log.Printf("getBatteryCharge: %v", err)
		}

		err = putBatteryStateToCloudwatch(c)
		if err != nil {
			log.Printf("putBatteryStateToCloudwatch: %v", err)
		}

		log.Printf("Logged battery charge to CW: %d%%", c)

		if _, err := client.Disconnect(); err != nil {
			log.Println(err)
		}

		time.Sleep(time.Minute)
	}
}

func getBatteryCharge(ups nut.UPS) (int64, error) {
	const batteryChargeVar = "battery.charge"
	for _, v := range ups.Variables {
		if v.Name == batteryChargeVar {
			val, ok := v.Value.(int64)
			if !ok {
				return -1,
					errors.New("UPS battery charge percentage is not an int64")
			}
			return val, nil
		}
	}
	return -1, errors.New("UPS battery charge is not available")
}

func putBatteryStateToCloudwatch(charge int64) error {
	sess := session.Must(session.NewSession())
	svc := cloudwatch.New(sess)

	m := &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("EskomMon"),
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String("UPS Charge"),
				Unit:       aws.String("Percent"),
				Value:      aws.Float64(float64(charge)),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String("UPS"),
						Value: aws.String("Mecer2000"),
					},
				},
			},
		},
	}
	_, err := svc.PutMetricData(m)
	return err
}
