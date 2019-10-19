package logger

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"fmt"
)

func Log(logGroupName string, logStreamName string, logData string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String("ap-southeast-2"),
		},
	}))

	svc := cloudwatchlogs.New(sess)

	result0, err0 := svc.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(logGroupName),
	})

	if err0 != nil {
		fmt.Println("Error", err0)
		return fmt.Errorf("Failed to describe log group")
	}

	if len(result0.LogGroups) == 0 {
		_, err1 := svc.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		})

		if err1 != nil {
			return fmt.Errorf("Failed to create log group")
		}
	}

	result1, err1 := svc.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(logGroupName),
		LogStreamNamePrefix: aws.String(logStreamName),
	})

	if err1 != nil {
		return fmt.Errorf("Failed to describe log stream")
	}

	var uploadSequenceToken *string

	if len(result1.LogStreams) == 0 {
		_, err2 := svc.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})

		if err2 != nil {
			return fmt.Errorf("Failed to create log stream")
		}
	} else {
		uploadSequenceToken = result1.LogStreams[0].UploadSequenceToken
	}

	timestamp := aws.TimeUnixMilli(time.Now())

	_, err3 := svc.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		LogEvents: []*cloudwatchlogs.InputLogEvent{
			&cloudwatchlogs.InputLogEvent{
				Message:   aws.String(logData),
				Timestamp: &timestamp,
			},
		},
		SequenceToken: uploadSequenceToken,
	})

	if err3 != nil {
		return fmt.Errorf("Failed to put log event")
	}

	return nil
}
