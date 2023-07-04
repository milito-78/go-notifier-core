package domains

import "time"

type NotifierTag struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ID        uint64
}

func NewNotifierTag(name string) *NotifierTag {
	return &NotifierTag{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
