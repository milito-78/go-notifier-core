package go_notifier_core

import (
	"encoding/json"
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

	var campaignRepo IEmailCampaignRepository
	err := container.Resolve(&campaignRepo)
	if err != nil {
		log.Fatalf("Error during resolve : %s", err)
	}
	campaign, err := campaignRepo.GetLatestCampaign()
	if err != nil {
		time.Sleep(time.Second * 20)
		return
	}
	campaign.StatusId = 2 //TODO status picked
	_ = campaignRepo.Update(campaign)

	tags := campaignRepo.GetCampaignTags(campaign.ID)
	if len(tags) == 0 {
		log.Printf("There is no tag saved for campagin = %d", campaign.ID)
		err := campaignRepo.Update(campaign)
		if err != nil {
			log.Printf("Error during update campaign : %s", err)
		}
		return
	}

	var subRepo IEmailSubscriberRepository
	err = container.Resolve(&subRepo)
	if err != nil {
		log.Fatalf("Error during resolve : %s", err)
	}

	queue := NewQueue("Email Queue")
	queue.StartListening()
	defer queue.CloseWorker()

	campaign.StatusId = 3 // TODO sending status
	_ = campaignRepo.Update(campaign)
	subscribers := subRepo.GetUsersByTagId(tags)
	for _, subscriber := range subscribers {
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

	campaign.StatusId = 4 // TODO complete status
	_ = campaignRepo.Update(campaign)
	return
}

func sendEmail(data any) error {
	message, ok := data.(*NotifierEmailMessage)
	if !ok {
		return errors.New("invalid data message to send email")
	}

	var messageRepo IEmailMessageRepository
	err := container.Resolve(&messageRepo)
	if err != nil {
		return err
	}
	err = messageRepo.CheckMessageExists(message)
	if err != nil {
		return nil
	}

	err = messageRepo.Create(message)
	if err != nil {
		return err
	}

	var emailServiceRepo IEmailServiceRepository
	err = container.Resolve(&emailServiceRepo)
	if err != nil {
		return err
	}
	service, err := emailServiceRepo.Get(message.EmailServiceId)
	if err != nil {
		log.Printf("Error during send mail (get service): %s", err)
		t := time.Now()
		message.FailedAt = &t
		er := messageRepo.Update(message)
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
		er := messageRepo.Update(message)
		if er != nil {
			log.Printf("Error during update failed at : %s\n", er)
		}
		return err
	}

	t := time.Now()
	message.SentAt = &t
	err = messageRepo.Update(message)
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
	var tmp map[string]interface{}
	_ = json.Unmarshal([]byte(service.Payload), &tmp)
	tmp["from"] = message.FromEmail
	j, _ := json.Marshal(tmp)
	mailer.SetConfig(j)
	return mailer.Send(message.RecipientEmail, message.Subject, message.Message)
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
