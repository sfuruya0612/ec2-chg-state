package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// Client ec2iface initialize
type Client struct {
	ec2iface.EC2API
}

// Instance EC2 info struct
type Instance struct {
	InstanceId string
	State      string
}

// Instances EC2 list
type Instances []Instance

func (c *Client) getInstances(tag string) (Instances, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(tag)},
			},
		},
	}

	output, err := c.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("Describe instances: %v", err)
	}

	list := Instances{}
	for _, r := range output.Reservations {
		for _, i := range r.Instances {
			if i.InstanceLifecycle != nil {
				continue
			}

			list = append(list, Instance{
				InstanceId: *i.InstanceId,
				State:      *i.State.Name,
			})
		}
	}

	return list, nil
}

func (c *Client) startEc2(id string) error {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}

	if _, err := c.StartInstances(input); err != nil {
		return fmt.Errorf("Start instances: %v", err)
	}

	return nil
}

func (c *Client) stopEc2(id string) error {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}

	if _, err := c.StopInstances(input); err != nil {
		return fmt.Errorf("Stop instances: %v", err)
	}

	return nil
}
