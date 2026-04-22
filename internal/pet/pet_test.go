package pet

import (
	"testing"
)

func TestGainXP_LevelsUp(t *testing.T) {
	p := New("Testy")
	p.MoodLevel = 50 // neutral mood, no bonus
	p.GainXP(100)
	if p.Level != 2 {
		t.Errorf("expected level 2, got %d", p.Level)
	}
	if p.XP != 0 {
		t.Errorf("expected 0 XP after exact level-up, got %d", p.XP)
	}
}

func TestGainXP_MultipleLevel(t *testing.T) {
	p := New("Testy")
	p.MoodLevel = 50
	p.GainXP(10000)
	if p.Level <= 2 {
		t.Errorf("expected multiple level-ups, got level %d", p.Level)
	}
}

func TestGainXP_Accumulates(t *testing.T) {
	p := New("Testy")
	p.MoodLevel = 50
	p.GainXP(50)
	if p.Level != 1 {
		t.Errorf("expected level 1, got %d", p.Level)
	}
	if p.XP != 50 {
		t.Errorf("expected 50 XP, got %d", p.XP)
	}
}

func TestGainXP_MoodBonus(t *testing.T) {
	p := New("Testy")
	p.MoodLevel = 100
	p.GainXP(10)
	if p.XP != 12 {
		t.Errorf("expected 12 XP with mood bonus, got %d", p.XP)
	}
}

func TestEatFood_RestoresEnergy(t *testing.T) {
	p := New("Testy")
	p.Energy = 20
	p.EatFood(50)
	if p.Energy != 70 {
		t.Errorf("expected energy 70, got %f", p.Energy)
	}
}

func TestEatFood_ClampsAt100(t *testing.T) {
	p := New("Testy")
	p.Energy = 90
	p.EatFood(50)
	if p.Energy != 100 {
		t.Errorf("expected energy clamped to 100, got %f", p.Energy)
	}
}

func TestSpendCoins_Succeeds(t *testing.T) {
	p := New("Testy")
	ok := p.SpendCoins(10)
	if !ok {
		t.Error("expected SpendCoins to succeed")
	}
	if p.Coins != 5 {
		t.Errorf("expected 5 coins remaining, got %d", p.Coins)
	}
}

func TestSpendCoins_FailsWhenInsufficient(t *testing.T) {
	p := New("Testy")
	p.Coins = 3
	ok := p.SpendCoins(10)
	if ok {
		t.Error("expected SpendCoins to fail with insufficient coins")
	}
	if p.Coins != 3 {
		t.Errorf("expected coins unchanged at 3, got %d", p.Coins)
	}
}

func TestCanGoOnMission_SufficientEnergy(t *testing.T) {
	p := New("Testy")
	p.Energy = 50
	if !p.CanGoOnMission() {
		t.Error("expected CanGoOnMission true with 50 energy")
	}
}

func TestCanGoOnMission_InsufficientEnergy(t *testing.T) {
	p := New("Testy")
	p.Energy = 10
	if p.CanGoOnMission() {
		t.Error("expected CanGoOnMission false with 10 energy")
	}
}

func TestCompleteSleep_FullRestore(t *testing.T) {
	p := New("Testy")
	p.Energy = 10
	p.MoodLevel = 15
	p.CompleteSleep()
	if p.Energy != 100 {
		t.Errorf("expected energy 100 after sleep, got %f", p.Energy)
	}
	if p.MoodLevel != 100 {
		t.Errorf("expected mood 100 after sleep, got %f", p.MoodLevel)
	}
}

func TestTick_DecaysAttributes(t *testing.T) {
	p := New("Testy")
	energy := p.Energy
	mood := p.MoodLevel
	p.Tick()
	if p.Energy >= energy {
		t.Error("expected energy to decay after tick")
	}
	if p.MoodLevel >= mood {
		t.Error("expected mood to decay after tick")
	}
}

func TestTick_ClampsAtZero(t *testing.T) {
	p := New("Testy")
	p.Energy = 0
	p.MoodLevel = 0
	p.Tick()
	if p.Energy < 0 || p.MoodLevel < 0 {
		t.Error("attributes went below zero after tick")
	}
}

func TestMood_Happy(t *testing.T) {
	p := New("Testy")
	p.MoodLevel = 90
	if p.Mood() != "happy" {
		t.Errorf("expected happy, got %s", p.Mood())
	}
}

func TestMood_Content(t *testing.T) {
	p := New("Testy")
	p.MoodLevel = 60
	if p.Mood() != "content" {
		t.Errorf("expected content, got %s", p.Mood())
	}
}

func TestMood_Sad(t *testing.T) {
	p := New("Testy")
	p.MoodLevel = 20
	if p.Mood() != "sad" {
		t.Errorf("expected sad, got %s", p.Mood())
	}
}

func TestXPProgress_ZeroAtStart(t *testing.T) {
	p := New("Testy")
	if p.XPProgress() != 0 {
		t.Errorf("expected 0 XP progress at start, got %f", p.XPProgress())
	}
}
