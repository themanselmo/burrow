package pet

import (
	"testing"
)

func TestGainXP_LevelsUp(t *testing.T) {
	p := New("Testy")
	p.GainXP(100) // level 1 threshold is 100
	if p.Level != 2 {
		t.Errorf("expected level 2, got %d", p.Level)
	}
	if p.XP != 0 {
		t.Errorf("expected 0 XP after exact level-up, got %d", p.XP)
	}
}

func TestGainXP_MultipleLevel(t *testing.T) {
	p := New("Testy")
	p.GainXP(10000)
	if p.Level <= 2 {
		t.Errorf("expected multiple level-ups, got level %d", p.Level)
	}
}

func TestGainXP_Accumulates(t *testing.T) {
	p := New("Testy")
	p.GainXP(50)
	if p.Level != 1 {
		t.Errorf("expected level 1, got %d", p.Level)
	}
	if p.XP != 50 {
		t.Errorf("expected 50 XP, got %d", p.XP)
	}
}

func TestFeed_ClampsAt100(t *testing.T) {
	p := New("Testy")
	p.Hunger = 95
	p.Feed()
	if p.Hunger != 100 {
		t.Errorf("expected hunger clamped to 100, got %f", p.Hunger)
	}
}

func TestTick_DecaysAttributes(t *testing.T) {
	p := New("Testy")
	hunger := p.Hunger
	p.Tick()
	if p.Hunger >= hunger {
		t.Errorf("expected hunger to decay after tick")
	}
}

func TestTick_ClampsAtZero(t *testing.T) {
	p := New("Testy")
	p.Hunger = 0
	p.Happiness = 0
	p.Energy = 0
	p.Affection = 0
	p.Tick()
	if p.Hunger < 0 || p.Happiness < 0 || p.Energy < 0 || p.Affection < 0 {
		t.Errorf("attributes went below zero after tick")
	}
}

func TestMood_Hungry(t *testing.T) {
	p := New("Testy")
	p.Hunger = 10
	if p.Mood() != "hungry" {
		t.Errorf("expected mood hungry, got %s", p.Mood())
	}
}

func TestMood_Tired(t *testing.T) {
	p := New("Testy")
	p.Energy = 10
	if p.Mood() != "tired" {
		t.Errorf("expected mood tired, got %s", p.Mood())
	}
}

func TestMood_Happy(t *testing.T) {
	p := New("Testy")
	p.Happiness = 100
	p.Affection = 100
	if p.Mood() != "happy" {
		t.Errorf("expected mood happy, got %s", p.Mood())
	}
}

func TestMood_Sad(t *testing.T) {
	p := New("Testy")
	p.Happiness = 10
	if p.Mood() != "sad" {
		t.Errorf("expected mood sad, got %s", p.Mood())
	}
}

func TestXPProgress_ZeroAtStart(t *testing.T) {
	p := New("Testy")
	if p.XPProgress() != 0 {
		t.Errorf("expected 0 XP progress at start, got %f", p.XPProgress())
	}
}

func TestXPProgress_HalfWay(t *testing.T) {
	p := New("Testy")
	p.GainXP(50)
	progress := p.XPProgress()
	if progress < 0.49 || progress > 0.51 {
		t.Errorf("expected ~0.5 XP progress, got %f", progress)
	}
}
