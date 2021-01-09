package pet

import (
	"sync"
	"time"
)

var (
	DefaultPetStats = []*Stat{
		{Name: "Health", Value: 100},
		{Name: "Boredom", Value: 0, Delta: 5},
		{Name: "Thirst", Value: 0, Delta: 3},
		{Name: "Hunger", Value: 0, Delta: 2},
		{Name: "Fatigue", Value: 0, Delta: 2},
	}
	StatNames          []string
	StatUpdateInterval = time.Second / 2

	CurrentPet *Pet
)

func init() {
	for _, stat := range DefaultPetStats {
		StatNames = append(StatNames, stat.Name)
	}
}

type Pet struct {
	Name   string
	IsDead bool

	StatLock           *sync.RWMutex
	Stats              []*Stat
	StatUpdateNotifier chan bool
}

func NewPet(name string) *Pet {

	p := &Pet{
		Name:               name,
		StatLock:           new(sync.RWMutex),
		Stats:              DefaultPetStats,
		StatUpdateNotifier: make(chan bool),
	}

	go func() {
		for {

			if p.IsDead {
				p.StatUpdateNotifier <- false // stop updating stats, trigger gravestone
				return
			}

			time.Sleep(StatUpdateInterval)

			var healthStatIndex int
			for i, stat := range p.Stats {
				if stat.Name == "Health" {
					healthStatIndex = i
				}
			}

			p.StatLock.Lock()
			for _, stat := range p.Stats {
				nv := stat.Value + stat.Delta
				if nv > 100 {
					nv = 100
				} else if nv < 0 {
					nv = 0
				}
				stat.Value = nv

				switch stat.Name {
				case "Thirst":
					if stat.Value == 100 {
						p.Stats[healthStatIndex].Value -= 5
					}
				case "Hunger":
					if stat.Value == 100 {
						p.Stats[healthStatIndex].Value -= 3
					}
				case "Fatigue":
					if stat.Value == 100 {
						p.Stats[healthStatIndex].Value -= 1
					}
				}

			}

			if p.Stats[healthStatIndex].Value <= 0 {
				p.IsDead = true
				p.Stats[healthStatIndex].Value = 0
			}

			p.StatLock.Unlock()

			p.StatUpdateNotifier <- true // keep updating stats
		}
	}()

	return p
}

type Stat struct {
	Name  string
	Value int
	Delta int
}
