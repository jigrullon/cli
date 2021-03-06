// This file was generated by counterfeiter
package app_eventsfakes

import (
	"sync"

	"github.com/cloudfoundry/cli/cf/api/app_events"
	"github.com/cloudfoundry/cli/cf/models"
)

type FakeAppEventsRepository struct {
	RecentEventsStub        func(appGuid string, limit int64) ([]models.EventFields, error)
	recentEventsMutex       sync.RWMutex
	recentEventsArgsForCall []struct {
		appGuid string
		limit   int64
	}
	recentEventsReturns struct {
		result1 []models.EventFields
		result2 error
	}
}

func (fake *FakeAppEventsRepository) RecentEvents(appGuid string, limit int64) ([]models.EventFields, error) {
	fake.recentEventsMutex.Lock()
	fake.recentEventsArgsForCall = append(fake.recentEventsArgsForCall, struct {
		appGuid string
		limit   int64
	}{appGuid, limit})
	fake.recentEventsMutex.Unlock()
	if fake.RecentEventsStub != nil {
		return fake.RecentEventsStub(appGuid, limit)
	} else {
		return fake.recentEventsReturns.result1, fake.recentEventsReturns.result2
	}
}

func (fake *FakeAppEventsRepository) RecentEventsCallCount() int {
	fake.recentEventsMutex.RLock()
	defer fake.recentEventsMutex.RUnlock()
	return len(fake.recentEventsArgsForCall)
}

func (fake *FakeAppEventsRepository) RecentEventsArgsForCall(i int) (string, int64) {
	fake.recentEventsMutex.RLock()
	defer fake.recentEventsMutex.RUnlock()
	return fake.recentEventsArgsForCall[i].appGuid, fake.recentEventsArgsForCall[i].limit
}

func (fake *FakeAppEventsRepository) RecentEventsReturns(result1 []models.EventFields, result2 error) {
	fake.RecentEventsStub = nil
	fake.recentEventsReturns = struct {
		result1 []models.EventFields
		result2 error
	}{result1, result2}
}

var _ app_events.AppEventsRepository = new(FakeAppEventsRepository)
