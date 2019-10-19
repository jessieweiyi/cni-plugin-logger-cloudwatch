package logger

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"fmt"
)

type CloudWatchLogger struct {
	LogGroupName         string
	LogStreamName        string
	CloudWatchLogsClient *cloudwatchlogs.CloudWatchLogs
}

func NewCloudWatchLogger(logGroupName string, containerID string, ifName string) (*CloudWatchLogger, error) {
	logStreamName := fmt.Sprintf("/%s/%s", containerID, ifName)

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
		return nil, fmt.Errorf("Failed to describe log group")
	}

	if len(result0.LogGroups) == 0 {
		_, err1 := svc.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		})

		if err1 != nil {
			return nil, fmt.Errorf("Failed to create log group")
		}
	}

	return &CloudWatchLogger{
		LogGroupName:         logGroupName,
		LogStreamName:        logStreamName,
		CloudWatchLogsClient: svc,
	}, nil
}

func (cwl *CloudWatchLogger) Log(cniLogData []byte) error {
	result1, err1 := cwl.CloudWatchLogsClient.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(cwl.LogGroupName),
		LogStreamNamePrefix: aws.String(cwl.LogStreamName),
	})

	if err1 != nil {
		return fmt.Errorf("Failed to describe log stream")
	}

	var uploadSequenceToken *string

	if len(result1.LogStreams) == 0 {
		_, err2 := cwl.CloudWatchLogsClient.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(cwl.LogGroupName),
			LogStreamName: aws.String(cwl.LogStreamName),
		})

		if err2 != nil {
			return fmt.Errorf("Failed to create log stream")
		}
	} else {
		uploadSequenceToken = result1.LogStreams[0].UploadSequenceToken
	}

	timestamp := aws.TimeUnixMilli(time.Now())

	_, err3 := cwl.CloudWatchLogsClient.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(cwl.LogGroupName),
		LogStreamName: aws.String(cwl.LogStreamName),
		LogEvents: []*cloudwatchlogs.InputLogEvent{
			&cloudwatchlogs.InputLogEvent{
				Message:   aws.String(string(cniLogData)),
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
