package dns

import (
	"fmt"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) String() string {
	return "Error " + e.Code + ": " + e.Message
}

type ErrorResponse struct {
	Errors []Error `xml:"Errors>Error"`
}

func (er ErrorResponse) String() string {
	e := ""
	for _, v := range er.Errors {
		e += v.String() + "\n"
	}
	return e
}

type Tag struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

func (d *InstanceSetItem) String() string {
	s := "InstanceSetItem: " + d.InstanceId + "\n\t" +
		"imageId: " + d.ImageId + " state: " + d.StateName + " host: " + d.DnsName + "\n"
	return s
}

type InstanceSetItem struct {
	InstanceId        string `xml:"instanceId"`
	ImageId           string `xml:"imageId"`
	StateCode         int    `xml:"instanceState>code"`
	StateName         string `xml:"instanceState>name"`
	CurrentStateCode  int    `xml:"currentState>code"`
	CurrentStateName  string `xml:"currentState>name"`
	PreviousStateCode int    `xml:"previousState>code"`
	PreviousStateName string `xml:"previousState>name"`
	DnsName           string `xml:"dnsName"`
	KeyName           string `xml:"keyName"`
	InstanceType      string `xml:"instanceType"`
	LaunchTime        string `xml:"launchTime"`
	IpAddress         string `xml:"ipAddress"`
	TagSet            []Tag  `xml:"tagSet>item"`
	Hypervisor        string `xml:"hypervisor"`
}

func (d *ReservationSetItem) String() string {
	s := "ReservationSetItem: " + d.ReservationId + "\n\t" +
		"requesterId: " + d.RequesterId + " ownerId: " + d.OwnerId + "\n"
	for _, v := range d.Items {
		s += fmt.Sprint(&v)
	}
	return s
}

type ReservationSetItem struct {
	ReservationId string            `xml:"reservationId"`
	RequesterId   string            `xml:"requesterId"`
	OwnerId       string            `xml:"ownerId"`
	Items         []InstanceSetItem `xml:"instancesSet>item"`
}

func (d *DescribeInstancesResponse) String() string {
	s := "DescribeInstancesResponse: " + d.RequestId + "\n"
	for _, v := range d.Items {
		s += fmt.Sprint(&v)
	}
	return s
}

type DescribeInstancesResponse struct {
	RequestId string               `xml:"requestId"`
	Items     []ReservationSetItem `xml:"reservationSet>item"`
}

// code: 16 running, 32 shutting-down, 48 terminated
type TerminateInstancesResponse struct {
	RequestId string            `xml:"requestId"`
	Items     []InstanceSetItem `xml:"instancesSet>item"`
}

type RunInstancesResponse struct {
	RequestId     string            `xml:"requestId"`
	ReservationId string            `xml:"reservationId"`
	OwnerId       string            `xml:"ownerId"`
	RequesterId   string            `xml:"requesterId"`
	Items         []InstanceSetItem `xml:"instancesSet>item"`
}
