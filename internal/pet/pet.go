package pet

import (
	"math"
	"time"
)

type Species string

const (
	SpeciesSeal Species = "seal"
)

type Pet struct {
	Name       string    `json:"name"`
	Species    Species   `json:"species"`
	Level      int       `json:"level"`
	XP         int       `json:"xp"`
	OwnedSince time.Time `json:"owned_since"`
	Energy     float64   `json:"energy"`
	MoodLevel  float64   `json:"mood_level"`
	Coins      int       `json:"coins"`
}

type LogEntry struct {
	Name       string     `json:"name"`
	Species    Species    `json:"species"`
	Level      int        `json:"level"`
	OwnedSince time.Time  `json:"owned_since"`
	ReleasedAt *time.Time `json:"released_at,omitempty"`
}

func New(name string) *Pet {
	return &Pet{
		Name:       name,
		Species:    SpeciesSeal,
		Level:      1,
		XP:         0,
		OwnedSince: time.Now(),
		Energy:     80,
		MoodLevel:  80,
		Coins:      15,
	}
}

func (p *Pet) PetAction() {
	p.MoodLevel = clamp(p.MoodLevel+15, 0, 100)
	p.GainXP(5)
}

func (p *Pet) Play() {
	p.MoodLevel = clamp(p.MoodLevel+20, 0, 100)
	p.Energy = clamp(p.Energy-15, 0, 100)
	p.GainXP(8)
}

func (p *Pet) EatFood(energyRestore float64) {
	p.Energy = clamp(p.Energy+energyRestore, 0, 100)
}

func (p *Pet) SpendCoins(amount int) bool {
	if p.Coins < amount {
		return false
	}
	p.Coins -= amount
	return true
}

func (p *Pet) AddCoins(amount int) {
	p.Coins += amount
}

func (p *Pet) CanGoOnMission() bool {
	return p.Energy >= 20
}

func (p *Pet) CanAffordMission(energyCost float64) bool {
	return p.Energy >= energyCost
}

func (p *Pet) CompleteMission(xp, coins int) {
	p.GainXP(xp)
	p.AddCoins(coins)
	p.MoodLevel = clamp(p.MoodLevel+15, 0, 100)
}

func (p *Pet) CompleteSleep() {
	p.Energy = 100
	p.MoodLevel = 100
}

func (p *Pet) Tick() {
	p.Energy = clamp(p.Energy-0.5, 0, 100)
	p.MoodLevel = clamp(p.MoodLevel-0.3, 0, 100)
}

func (p *Pet) GainXP(amount int) {
	// Mood bonus: happy pet earns 20% more XP
	if p.MoodLevel > 75 {
		amount = int(float64(amount) * 1.2)
	}
	p.XP += amount
	for p.XP >= xpForNextLevel(p.Level) {
		p.XP -= xpForNextLevel(p.Level)
		p.Level++
	}
}

func (p *Pet) Mood() string {
	if p.MoodLevel > 75 {
		return "happy"
	}
	if p.MoodLevel > 40 {
		return "content"
	}
	return "sad"
}

func (p *Pet) IsSleepy() bool {
	hour := time.Now().Hour()
	return hour >= 21 || hour < 7
}

func (p *Pet) XPProgress() float64 {
	return float64(p.XP) / float64(xpForNextLevel(p.Level))
}

func xpForNextLevel(level int) int {
	return int(math.Round(100 * math.Pow(1.2, float64(level-1))))
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
