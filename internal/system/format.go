package system

import (
	"context"

	log "github.com/sirupsen/logrus"
	"hurracloud.io/jawhar/internal/agent"
	pb "hurracloud.io/jawhar/internal/agent/proto"
)

type FormatJob struct {
	status     string
	deviceType string
	deviceFile string
	filesystem string
}

func CreateFormatJob(deviceType string, deviceFile string, filesystem string) *FormatJob {
	return &FormatJob{
		deviceType: deviceType,
		deviceFile: deviceFile,
		filesystem: filesystem,
	}
}

func (j *FormatJob) Run() {
	j.status = "started"
	_, err := agent.Client.FormatDrive(context.Background(), &pb.FormatDriveRequest{DeviceFile: j.deviceFile})
	if err != nil {
		log.Error("Agent Client Failed to call FormatDrive: ", err)
		j.status = "error"
		return // TODO Notify user of error
	}
	j.status = "completed"
}

func (j *FormatJob) Status() string {
	return j.status
}
