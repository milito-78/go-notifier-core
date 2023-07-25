package go_notifier_core

import (
	"errors"
	"github.com/golobby/container/v3"
	"log"
)

func Seed() error {
	if !initialized {
		return errors.New("dependencies not initialized")
	}

	emailUnsubReasonSeeder()
	emailUCampaignStatusSeeder()
	log.Println("Seeder runs successfully")
	return nil
}

func emailUnsubReasonSeeder() {
	var repo IEmailUnSubEventRepository
	err := container.Resolve(&repo)
	if err != nil {
		log.Fatalf("Error during resolve DI : %s", err)
		return
	}

	repo.FirstOrCreate(NewNotifierEmailUnsubscribeEvent("Bounce", NotifierEmailUnsubBounce))
	repo.FirstOrCreate(NewNotifierEmailUnsubscribeEvent("Complaint", NotifierEmailUnsubComplaint))
	repo.FirstOrCreate(NewNotifierEmailUnsubscribeEvent("Manual by Admin", NotifierEmailUnsubManualByAdmin))
	repo.FirstOrCreate(NewNotifierEmailUnsubscribeEvent("Manual by Subscriber", NotifierEmailUnsubManualBySubscriber))

}

func emailUCampaignStatusSeeder() {
	var repo IEmailStatusRepository
	err := container.Resolve(&repo)
	if err != nil {
		log.Fatalf("Error during resolve DI : %s", err)
		return
	}

	repo.FirstOrCreate(NewNotifierEmailStatus("Draft", NotifierEmailStatusDraft))
	repo.FirstOrCreate(NewNotifierEmailStatus("Queued", NotifierEmailStatusQueued))
	repo.FirstOrCreate(NewNotifierEmailStatus("Sending", NotifierEmailStatusSending))
	repo.FirstOrCreate(NewNotifierEmailStatus("Sent", NotifierEmailStatusSent))
	repo.FirstOrCreate(NewNotifierEmailStatus("Canceled", NotifierEmailStatusCanceled))
	repo.FirstOrCreate(NewNotifierEmailStatus("Failed", NotifierEmailStatusFailed))
}
