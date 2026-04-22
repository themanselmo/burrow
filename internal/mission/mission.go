package mission

import (
	"math/rand"
	"time"
)

type Type string

const (
	TypeExplore Type = "explore"
	TypeSleep   Type = "sleep"
)

type Item struct {
	Name  string `json:"name"`
	Theme string `json:"theme"`
}

type Mission struct {
	StartedAt       time.Time `json:"started_at"`
	DurationMinutes int       `json:"duration_minutes"`
	Theme           string    `json:"theme"`
	Type            Type      `json:"type"`
}

type Result struct {
	Type     Type
	XP       int
	Coins    int
	Item     *Item
	Duration int
}

type Waypoint struct {
	Label    string
	FlavorIn string
}

type Theme struct {
	Name      string
	Waypoints []Waypoint
	ItemPool  []Item
}

var Themes = map[string]Theme{
	"forest_path": {
		Name: "The Forest Path",
		Waypoints: []Waypoint{
			{Label: "Meadow", FlavorIn: "stepping into the tall grass"},
			{Label: "Tall Oaks", FlavorIn: "weaving between ancient trunks"},
			{Label: "Mossy Clearing", FlavorIn: "resting on a soft mossy rock"},
			{Label: "Home", FlavorIn: "waddling back with a big smile"},
		},
		ItemPool: []Item{
			{Name: "Pinecone", Theme: "forest_path"},
			{Name: "Wildflower", Theme: "forest_path"},
			{Name: "Mossy Rock", Theme: "forest_path"},
			{Name: "Acorn", Theme: "forest_path"},
		},
	},
	"mountain_pass": {
		Name: "The Mountain Pass",
		Waypoints: []Waypoint{
			{Label: "Foothills", FlavorIn: "starting the long climb"},
			{Label: "Rocky Trail", FlavorIn: "scrambling over loose stones"},
			{Label: "Summit Ridge", FlavorIn: "catching breath at the top"},
			{Label: "Home", FlavorIn: "sliding downhill with glee"},
		},
		ItemPool: []Item{
			{Name: "Smooth Stone", Theme: "mountain_pass"},
			{Name: "Eagle Feather", Theme: "mountain_pass"},
			{Name: "Mountain Crystal", Theme: "mountain_pass"},
			{Name: "Dried Lichen", Theme: "mountain_pass"},
		},
	},
}

var FreeThemes = []string{"forest_path", "mountain_pass"}

// DreamWaypoints are used for sleep missions instead of explore themes.
var DreamWaypoints = []Waypoint{
	{Label: "Cloud Fields", FlavorIn: "floating on a warm fluffy cloud"},
	{Label: "The Warm Rock", FlavorIn: "basking under an endless golden sun"},
	{Label: "Starlit Waters", FlavorIn: "drifting through a glowing deep blue sea"},
	{Label: "Home", FlavorIn: "slowly drifting back to the surface"},
}

func RandomTheme(owned []string) Theme {
	pool := owned
	if len(pool) == 0 {
		pool = FreeThemes
	}
	return Themes[pool[rand.Intn(len(pool))]]
}

func (m *Mission) ReturnTime() time.Time {
	return m.StartedAt.Add(time.Duration(m.DurationMinutes) * time.Minute)
}

func (m *Mission) IsComplete() bool {
	return time.Now().After(m.ReturnTime())
}

func (m *Mission) ElapsedFraction() float64 {
	total := m.ReturnTime().Sub(m.StartedAt).Seconds()
	elapsed := time.Since(m.StartedAt).Seconds()
	if elapsed >= total {
		return 1.0
	}
	if elapsed < 0 {
		return 0.0
	}
	return elapsed / total
}

func (m *Mission) MinutesRemaining() int {
	rem := time.Until(m.ReturnTime())
	if rem < 0 {
		return 0
	}
	return int(rem.Minutes()) + 1
}

func EnergyCost(durationMinutes int) float64 {
	cost := float64(durationMinutes) * 0.8
	if cost > 70 {
		cost = 70
	}
	return cost
}

func (m *Mission) Calculate() Result {
	if m.Type == TypeSleep {
		return Result{Type: TypeSleep, Duration: m.DurationMinutes}
	}

	xp := m.DurationMinutes * 4
	if xp > 250 {
		xp = 250
	}

	coins := int(float64(m.DurationMinutes) * 0.4)
	if coins > 25 {
		coins = 25
	}

	var dropped *Item
	dropChance := 0.10 + float64(m.DurationMinutes)*0.001
	if dropChance > 0.25 {
		dropChance = 0.25
	}
	if rand.Float64() < dropChance {
		theme := Themes[m.Theme]
		pick := theme.ItemPool[rand.Intn(len(theme.ItemPool))]
		dropped = &pick
	}

	return Result{Type: TypeExplore, XP: xp, Coins: coins, Item: dropped, Duration: m.DurationMinutes}
}

func NewExplore(durationMinutes int, ownedThemes []string) *Mission {
	theme := RandomTheme(ownedThemes)
	key := ""
	for k, t := range Themes {
		if t.Name == theme.Name {
			key = k
			break
		}
	}
	return &Mission{
		StartedAt:       time.Now(),
		DurationMinutes: durationMinutes,
		Theme:           key,
		Type:            TypeExplore,
	}
}

func NewSleep(durationMinutes int) *Mission {
	return &Mission{
		StartedAt:       time.Now(),
		DurationMinutes: durationMinutes,
		Type:            TypeSleep,
	}
}
