package repo

import (
	"github.com/google/uuid"
	"remindme/server/common"
)

type inMemoryEventRepo struct {
	events map[string]common.Reminder
}

func (repo *inMemoryEventRepo) Add(event common.Reminder) {
	repo.events[event.ID.String()] = event
}

func (repo *inMemoryEventRepo) List() []common.Reminder {
	eventsAsList := make([]common.Reminder, len(repo.events))
	i := 0

	for _, event := range repo.events {
		eventsAsList[i] = event
		i++
	}

	return eventsAsList
}

func (repo *inMemoryEventRepo) DeleteAll() {
	repo.events = make(map[string]common.Reminder, 0)
}

func (repo *inMemoryEventRepo) Delete(id uuid.UUID) {
	delete(repo.events, id.String())
}

func (repo *inMemoryEventRepo) Exists(id uuid.UUID) bool {
	_, found := repo.events[id.String()]
	return found
}

func NewImMemoryEventRepo() Repo {
	return &inMemoryEventRepo{events: make(map[string]common.Reminder, 0)}
}
