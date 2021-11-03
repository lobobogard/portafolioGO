package handler

import (
	"net/http"
	"time"

	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

type TimeSendEmail struct {
	StartTime  int64
	FinishTime int64
	DataEmail  []DataEmail
}
type DataEmail struct {
	Name       string
	Email      string
	StartTime  int64
	FinishTime int64
}

func Emails(DB *gorm.DB, w http.ResponseWriter, r *http.Request) (bool, TimeSendEmail) {
	userJWT := DecodeSessionUserJWT(w, r)
	var confConcurrency model.ConfConcurrency
	DB.Where("username", userJWT.Id).First(&confConcurrency)
	var timeSendEmail TimeSendEmail

	if confConcurrency.SendEmail {
		if confConcurrency.Concurrency {
			timeSendEmail = sendConcurrencyEmail()
		} else {
			timeSendEmail.StartTime = time.Now().Unix()
			timeSendEmail.DataEmail = append(timeSendEmail.DataEmail, sendEmail("gerardo", "pulido_a.@TI.com", 2), sendEmail("jennifer", "jennifer@RH.com", 5), sendEmail("william", "wallace@BH.free", 4))
			timeSendEmail.FinishTime = time.Now().Unix()
		}
	}
	return confConcurrency.SendEmail, timeSendEmail
}

func sendEmail(name string, email string, times time.Duration) DataEmail {
	var dataEmail DataEmail
	dataEmail.Name = name
	dataEmail.Email = email
	dataEmail.StartTime = time.Now().Unix()
	time.Sleep(times * time.Second)
	dataEmail.FinishTime = time.Now().Unix()
	return dataEmail
}

func sendEmailConcurrency(chanDataEmail chan<- DataEmail, name string, email string, times time.Duration) {
	var dataEmail DataEmail
	dataEmail.Name = name
	dataEmail.Email = email
	dataEmail.StartTime = time.Now().Unix()
	time.Sleep(times * time.Second)
	dataEmail.FinishTime = time.Now().Unix()
	chanDataEmail <- dataEmail
}

func sendConcurrencyEmail() TimeSendEmail {
	chanDataEmail := make(chan DataEmail)
	go sendEmailConcurrency(chanDataEmail, "gerardo", "pulido_a.@TI.com", 2)
	go sendEmailConcurrency(chanDataEmail, "jennifer", "jennifer@RH.com", 5)
	go sendEmailConcurrency(chanDataEmail, "william", "wallace@BH.free", 4)

	var timeSendEmail TimeSendEmail
	timeSendEmail.StartTime = time.Now().Unix()
	timeSendEmail.DataEmail = append(timeSendEmail.DataEmail, <-chanDataEmail, <-chanDataEmail, <-chanDataEmail)
	timeSendEmail.FinishTime = time.Now().Unix()
	return timeSendEmail
}
