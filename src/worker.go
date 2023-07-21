package go_notifier_core

import (
	"errors"
	"fmt"
	"github.com/golobby/container/v3"
	"log"
	"time"
)

type (
	//IWorker used for cronjob structs
	IWorker interface {
		Run()
	}

	WorkerConfig struct {
		Duration time.Duration
		Worker   IWorker
		Name     string
	}

	WorkersList []WorkerConfig

	EmailWorker struct {
	}

	MobileWorker struct {
	}

	NotificationWorker struct {
	}
)

func (e EmailWorker) Run() {

	campaign, err := GetLatestCampaignForRun()
	if err != nil {
		log.Printf("error during run email worker : %s", err)
		return
	}

	campaign.StatusId = NotifierEmailStatusDraft
	_ = UpdateEmailCampaign(campaign)

	tags := GetEmailCampaignTags(campaign.ID)
	if len(tags) == 0 {
		log.Printf("There is no tag saved for campagin = %d", campaign.ID)
		campaign.StatusId = NotifierEmailStatusFailed
		err := UpdateEmailCampaign(campaign)
		if err != nil {
			log.Printf("Error during update campaign : %s", err)
		}
		return
	}

	queue := NewQueue("Email Queue")
	queue.StartListening()
	defer queue.CloseWorker()

	campaign.StatusId = NotifierEmailStatusSending
	_ = UpdateEmailCampaign(campaign)
	subscribers, err := GetEmailSubscribersWithTags(tags)
	if err != nil {
		log.Printf("error during get subs for tags email : %s", err)
	}

	for _, subscriber := range subscribers {
		log.Println("Subscriber id is : ", subscriber.ID)
		queue.Send(NewQueueMessage(sendEmail, NewNotifierEmailMessage(
			subscriber.Email,
			subscriber.ID,
			"campaign",
			campaign.FromEmail,
			campaign.ID,
			campaign.FromName,
			campaign.Subject,
			campaign.EmailServiceId,
			campaign.Content,
		)))
	}

	campaign.StatusId = NotifierEmailStatusSent
	_ = UpdateEmailCampaign(campaign)
}

func sendEmail(data any) error {
	message, ok := data.(*NotifierEmailMessage)
	if !ok {
		return errors.New("invalid data message to send email")
	}
	err := CheckEmailMessageExists(message)
	if err != nil {
		return nil
	}

	err = CreateEmailMessage(message)
	if err != nil {
		return err
	}

	service, err := GetEmailServiceById(message.EmailServiceId)
	if err != nil {
		log.Printf("Error during send mail (get service): %s", err)
		t := time.Now()
		message.FailedAt = &t
		er := UpdateEmailMessage(message)
		if er != nil {
			log.Printf("Error during update failed at : %s\n", er)
		}
		return err
	}

	err = handleMail(service, message)
	if err != nil {
		log.Printf("Error during send mail : %s\n", err)
		t := time.Now()
		message.FailedAt = &t
		er := UpdateEmailMessage(message)
		if er != nil {
			log.Printf("Error during update failed at : %s\n", er)
		}
		return err
	}

	t := time.Now()
	message.SentAt = &t
	err = UpdateEmailMessage(message)
	if err != nil {
		return err
	}
	return nil
}

func handleMail(service *NotifierEmailService, message *NotifierEmailMessage) error {
	var mailer Mailer
	err := container.NamedResolve(&mailer, service.Type)
	if err != nil {
		return err
	}
	mailer.SetConfig([]byte(service.Payload))
	return mailer.Send(message.FromName, message.FromEmail, message.RecipientEmail, message.Subject, message.Message)
}

func (m MobileWorker) Run() {
	panic("TODO implement")
}

func (n NotificationWorker) Run() {
	panic("TODO implement")
}

// WorkerStart starts cronjob workers
func WorkerStart(config WorkersList) {
	for _, workerConfig := range config {
		go func(c WorkerConfig) {
			fmt.Printf("Worker %s starts work\n", c.Name)
			cron := time.NewTicker(c.Duration)
			for range cron.C {
				c.Worker.Run()
			}
			fmt.Printf("Worker %s stops\n", c.Name)
		}(workerConfig)
	}
}

type Queue struct {
	name string
	recv chan *QueueMessage
	quit chan bool
}

func NewQueue(name string) *Queue {
	return &Queue{
		name: name,
	}
}

func (q *Queue) StartListening() {
	log.Printf("Initializing %s's queue...\n", q.name)
	q.recv = make(chan *QueueMessage)
	q.quit = make(chan bool)
	go q.listen()
}

func (q *Queue) Send(data *QueueMessage) {
	q.recv <- data
}

func (q *Queue) CloseWorker() {
	log.Printf("Closing %s's queue listening. Please wait...\n", q.name)
	q.quit <- true
	close(q.quit)
	close(q.recv)
}

func (q *Queue) listen() {
	log.Println("Queue starts ...")
	for {
		select {
		case tmp := <-q.recv:
			if tmp.handle() != true {
				log.Println("There is an error during handle message")
				if tmp.err != "" {
					log.Printf("Report : %s\n", tmp.err)
				}
			} else {
				log.Println("Queue handled successfully")
			}
		case _, ok := <-q.quit:
			if ok {
				log.Println("Queue stops.")
			}
			return
		}
	}
}

type QueueMessage struct {
	handler func(data interface{}) error
	data    any
	err     string
}

func (qm QueueMessage) handle() bool {
	if err := qm.handler(qm.data); err != nil {
		qm.err = err.Error()
		return false
	} else {
		return true
	}
}

func NewQueueMessage(handle func(data any) error, data interface{}) *QueueMessage {
	return &QueueMessage{handler: handle, data: data, err: ""}
}
