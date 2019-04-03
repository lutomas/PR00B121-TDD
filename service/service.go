package service

import "time"

type Repo interface {
	Write(email string, connectTime time.Time) error
	Read(email string) (time.Time, error)
}

type EmailService struct {
	repo Repo
}
