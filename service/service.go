package service

import (
	"github.com/Sirupsen/logrus"
	"github.com/lutomas/PR00B121-TDD/boltdb"
	"time"
)

type Repo interface {
	Write(email string, connectTime time.Time) error
	Read(email string) (string, error)
}

type EmailService struct {
	repo Repo
}

func NewEmailService() (result *EmailService) {
	result = new(EmailService)
	result.repo = boltdb.NewBoltDb()
	return
}

func (service *EmailService) Process(action, email string) (err error) {
	writeTime, err := service.repo.Read(email)
	if err != nil {
		return err
	}

	logrus.Infof("Read: email=%s, time=%s", email, writeTime)
	return
}
