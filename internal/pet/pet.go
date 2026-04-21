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
	Name      string    `json:"name"`
	Species   Species   `json:"species"`
	Level     int       `json:"level"`
	XP        int       `json:"xp"`
	OwnedSince time.Time `json:"owned_since"`

	Hunger    float64 `json:"hunger"`
	Happiness float64 `json:"happiness"`
	Energy    float64 `json:"energy"`
	Affection float64 `json:"affection"`
}

type LogEntry struct {
	Name       string    `json:"name"`
	Species    Species   `json:"species"`
	Level      int       `json:"level"`
	OwnedSince time.Time `json:"owned_since"`
	ReleasedAt *time.Time `json:"released_at,omitempty"`
}

func New(name string) *Pet {
	return &Pet{
		Name:       name,
		Species:    SpeciesSeal,
		Level:      1,
		XP:         0,
		OwnedSince: time.Now(),
		Hunger:     80,
		Happiness:  80,
		Energy:     80,
		Affection:  80,
	}
}

func (p *Pet) Feed() {
	p.Hunger = clamp(p.Hunger+20, 0, 100)
	p.GainXP(5)
}

func (p *Pet) Pet_() {
	p.Affection = clamp(p.Affection+15, 0, 100)
	p.Happiness = clamp(p.Happiness+10, 0, 100)
	p.GainXP(5)
}

func (p *Pet) Play() {
	p.Happiness = clamp(p.Happiness+20, 0, 100)
	p.Energy = clamp(p.Energy-10, 0, 100)
	p.GainXP(8)
}

func (p *Pet) Tick() {
	p.Hunger = clamp(p.Hunger-1.5, 0, 100)
	p.Happiness = clamp(p.Happiness-0.5, 0, 100)
	p.Energy = clamp(p.Energy-0.3, 0, 100)
	p.Affection = clamp(p.Affection-0.2, 0, 100)
}

func (p *Pet) GainXP(amount int) {
	p.XP += amount
	threshold := xpForNextLevel(p.Level)
	for p.XP >= threshold {
		p.XP -= threshold
		p.Level++
		threshold = xpForNextLevel(p.Level)
	}
}

func (p *Pet) Mood() string {
	if p.Hunger < 25 {
		return "hungry"
	}
	if p.Energy < 25 {
		return "tired"
	}
	if p.Happiness < 25 {
		return "sad"
	}
	if p.Happiness > 75 && p.Affection > 75 {
		return "happy"
	}
	return "content"
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
