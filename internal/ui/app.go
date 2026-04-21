package ui

import (
	"time"

	"github.com/anselmo/burrow/internal/mission"
	"github.com/anselmo/burrow/internal/pet"
	"github.com/anselmo/burrow/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenNaming screen = iota
	screenHome
	screenLog
	screenRelease
	screenFarewell
	screenAbruptQuit
	screenMissionStart
	screenMissionActive
	screenMissionComplete
	screenInventory
)

type AppModel struct {
	screen          screen
	naming          NamingModel
	home            HomeModel
	log             LogModel
	release         ReleaseModel
	farewell        FarewellModel
	abruptQuit      AbruptQuitModel
	missionStart    MissionStartModel
	missionActive   MissionActiveModel
	missionComplete MissionCompleteModel
	inventory       InventoryModel
	pet             *pet.Pet
	storage         *storage.State
}

func NewAppModel(state *storage.State) AppModel {
	if state == nil {
		state = &storage.State{}
	}

	m := AppModel{storage: state}

	if state.Pet == nil {
		m.screen = screenNaming
		m.naming = NewNamingModel()
		return m
	}

	m.pet = state.Pet

	// On-launch mission check.
	if state.Mission != nil {
		if state.Mission.IsComplete() {
			result := state.Mission.Calculate()
			m.pet.GainXP(result.XP)
			if result.Item != nil {
				state.Items = append(state.Items, *result.Item)
			}
			state.Mission = nil
			_ = storage.SaveState(state)
			m.screen = screenMissionComplete
			m.missionComplete = NewMissionCompleteModel(result, m.pet)
			return m
		}
		m.screen = screenMissionActive
		m.missionActive = NewMissionActiveModel(state.Mission, m.pet)
		return m
	}

	m.screen = screenHome
	m.home = NewHomeModel(m.pet)
	return m
}

func (m AppModel) Init() tea.Cmd {
	switch m.screen {
	case screenNaming:
		return m.naming.Init()
	case screenHome:
		return m.home.Init()
	case screenMissionActive:
		return m.missionActive.Init()
	case screenMissionComplete:
		return m.missionComplete.Init()
	}
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NameSubmittedMsg:
		m.pet = pet.New(msg.Name)
		m.storage.Pet = m.pet
		_ = storage.SaveState(m.storage)
		m.screen = screenHome
		m.home = NewHomeModel(m.pet)
		return m, m.home.Init()

	case OpenLogMsg:
		m.screen = screenLog
		m.log = NewLogModel(m.pet)
		return m, m.log.Init()

	case LogBackMsg:
		m.screen = screenHome
		return m, nil

	case OpenReleaseMsg:
		m.screen = screenRelease
		m.release = NewReleaseModel(m.pet)
		return m, m.release.Init()

	case ReleaseCancelMsg:
		m.screen = screenHome
		return m, nil

	case ReleaseConfirmedMsg:
		m.logCurrentPet()
		m.screen = screenFarewell
		m.farewell = NewFarewellModel(m.pet)
		return m, m.farewell.Init()

	case FarewellDoneMsg:
		m.pet = nil
		m.storage.Pet = nil
		m.storage.Mission = nil
		_ = storage.SaveState(m.storage)
		m.screen = screenNaming
		m.naming = NewNamingModel()
		return m, m.naming.Init()

	case OpenMissionStartMsg:
		m.screen = screenMissionStart
		m.missionStart = NewMissionStartModel()
		return m, m.missionStart.Init()

	case MissionCancelMsg:
		m.screen = screenHome
		return m, nil

	case MissionStartMsg:
		ms := mission.New(msg.DurationMinutes, nil)
		m.storage.Mission = ms
		_ = storage.SaveState(m.storage)
		m.screen = screenMissionActive
		m.missionActive = NewMissionActiveModel(ms, m.pet)
		return m, m.missionActive.Init()

	case MissionCompleteCheckMsg:
		if m.storage.Mission == nil {
			m.screen = screenHome
			return m, nil
		}
		result := m.storage.Mission.Calculate()
		m.pet.GainXP(result.XP)
		if result.Item != nil {
			m.storage.Items = append(m.storage.Items, *result.Item)
		}
		m.storage.Mission = nil
		_ = storage.SaveState(m.storage)
		m.screen = screenMissionComplete
		m.missionComplete = NewMissionCompleteModel(result, m.pet)
		return m, m.missionComplete.Init()

	case OpenInventoryMsg:
		m.screen = screenInventory
		m.inventory = NewInventoryModel(m.storage.Items)
		return m, m.inventory.Init()

	case InventoryBackMsg:
		m.screen = screenHome
		return m, nil

	case MissionDismissMsg:
		m.screen = screenHome
		m.home = NewHomeModel(m.pet)
		return m, m.home.Init()

	case QuitMsg:
		if msg.Abrupt && m.pet != nil {
			m.screen = screenAbruptQuit
			m.abruptQuit = NewAbruptQuitModel(m.pet)
			_ = storage.SaveState(m.storage)
			return m, m.abruptQuit.Init()
		}
		_ = storage.SaveState(m.storage)
		return m, tea.Quit
	}

	switch m.screen {
	case screenNaming:
		updated, cmd := m.naming.Update(msg)
		m.naming = updated.(NamingModel)
		return m, cmd

	case screenHome:
		updated, cmd := m.home.Update(msg)
		m.home = updated.(HomeModel)
		m.storage.Pet = m.home.pet
		return m, cmd

	case screenLog:
		updated, cmd := m.log.Update(msg)
		m.log = updated.(LogModel)
		return m, cmd

	case screenRelease:
		updated, cmd := m.release.Update(msg)
		m.release = updated.(ReleaseModel)
		return m, cmd

	case screenFarewell:
		updated, cmd := m.farewell.Update(msg)
		m.farewell = updated.(FarewellModel)
		return m, cmd

	case screenAbruptQuit:
		updated, cmd := m.abruptQuit.Update(msg)
		m.abruptQuit = updated.(AbruptQuitModel)
		return m, cmd

	case screenMissionStart:
		updated, cmd := m.missionStart.Update(msg)
		m.missionStart = updated.(MissionStartModel)
		return m, cmd

	case screenMissionActive:
		updated, cmd := m.missionActive.Update(msg)
		m.missionActive = updated.(MissionActiveModel)
		return m, cmd

	case screenMissionComplete:
		updated, cmd := m.missionComplete.Update(msg)
		m.missionComplete = updated.(MissionCompleteModel)
		return m, cmd

	case screenInventory:
		updated, cmd := m.inventory.Update(msg)
		m.inventory = updated.(InventoryModel)
		return m, cmd
	}

	return m, nil
}

func (m AppModel) View() string {
	switch m.screen {
	case screenNaming:
		return m.naming.View()
	case screenHome:
		return m.home.View()
	case screenLog:
		return m.log.View()
	case screenRelease:
		return m.release.View()
	case screenFarewell:
		return m.farewell.View()
	case screenAbruptQuit:
		return m.abruptQuit.View()
	case screenMissionStart:
		return m.missionStart.View()
	case screenMissionActive:
		return m.missionActive.View()
	case screenMissionComplete:
		return m.missionComplete.View()
	case screenInventory:
		return m.inventory.View()
	}
	return ""
}

func (m *AppModel) logCurrentPet() {
	if m.pet == nil {
		return
	}
	log, _ := storage.LoadLog()
	if log == nil {
		log = &storage.Log{}
	}
	now := time.Now()
	log.Entries = append([]pet.LogEntry{{
		Name:       m.pet.Name,
		Species:    m.pet.Species,
		Level:      m.pet.Level,
		OwnedSince: m.pet.OwnedSince,
		ReleasedAt: &now,
	}}, log.Entries...)
	_ = storage.SaveLog(log)
}
