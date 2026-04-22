package mission

import (
	"testing"
	"time"
)

func TestCalculate_XPScalesLinearly(t *testing.T) {
	m := &Mission{StartedAt: time.Now(), DurationMinutes: 25, Theme: "forest_path"}
	result := m.Calculate()
	if result.XP != 100 {
		t.Errorf("expected 100 XP for 25min, got %d", result.XP)
	}

	m.DurationMinutes = 45
	result = m.Calculate()
	if result.XP != 180 {
		t.Errorf("expected 180 XP for 45min, got %d", result.XP)
	}
}

func TestCalculate_XPCappedAt250(t *testing.T) {
	m := &Mission{StartedAt: time.Now(), DurationMinutes: 90, Theme: "forest_path"}
	result := m.Calculate()
	if result.XP != 250 {
		t.Errorf("expected 250 XP (cap) for 90min, got %d", result.XP)
	}
}

func TestCalculate_XPDoesNotExceedCap(t *testing.T) {
	m := &Mission{StartedAt: time.Now(), DurationMinutes: 480, Theme: "forest_path"}
	result := m.Calculate()
	if result.XP > 250 {
		t.Errorf("XP exceeded cap of 250 for 480min mission, got %d", result.XP)
	}
}

func TestIsComplete_WhenPast(t *testing.T) {
	m := &Mission{
		StartedAt:       time.Now().Add(-30 * time.Minute),
		DurationMinutes: 25,
		Theme:           "forest_path",
	}
	if !m.IsComplete() {
		t.Error("expected mission to be complete")
	}
}

func TestIsComplete_WhenFuture(t *testing.T) {
	m := &Mission{
		StartedAt:       time.Now(),
		DurationMinutes: 25,
		Theme:           "forest_path",
	}
	if m.IsComplete() {
		t.Error("expected mission to not be complete yet")
	}
}

func TestElapsedFraction_Midpoint(t *testing.T) {
	m := &Mission{
		StartedAt:       time.Now().Add(-12*time.Minute - 30*time.Second),
		DurationMinutes: 25,
		Theme:           "forest_path",
	}
	f := m.ElapsedFraction()
	if f < 0.48 || f > 0.52 {
		t.Errorf("expected ~0.5 elapsed fraction at midpoint, got %f", f)
	}
}

func TestElapsedFraction_Clamps(t *testing.T) {
	past := &Mission{
		StartedAt:       time.Now().Add(-2 * time.Hour),
		DurationMinutes: 25,
		Theme:           "forest_path",
	}
	if past.ElapsedFraction() != 1.0 {
		t.Errorf("expected elapsed fraction to clamp at 1.0, got %f", past.ElapsedFraction())
	}

	future := &Mission{
		StartedAt:       time.Now().Add(1 * time.Hour),
		DurationMinutes: 25,
		Theme:           "forest_path",
	}
	if future.ElapsedFraction() != 0.0 {
		t.Errorf("expected elapsed fraction to clamp at 0.0, got %f", future.ElapsedFraction())
	}
}

func TestMinutesRemaining_AlmostDone(t *testing.T) {
	m := &Mission{
		StartedAt:       time.Now().Add(-24 * time.Minute),
		DurationMinutes: 25,
		Theme:           "forest_path",
	}
	if m.MinutesRemaining() > 2 {
		t.Errorf("expected <= 2 minutes remaining, got %d", m.MinutesRemaining())
	}
}

func TestEnergyCost_Scales(t *testing.T) {
	if EnergyCost(25) != 20 {
		t.Errorf("expected 20 energy cost for 25min, got %f", EnergyCost(25))
	}
	if EnergyCost(45) != 36 {
		t.Errorf("expected 36 energy cost for 45min, got %f", EnergyCost(45))
	}
}

func TestEnergyCost_Cap(t *testing.T) {
	if EnergyCost(480) > 70 {
		t.Errorf("expected energy cost capped at 70, got %f", EnergyCost(480))
	}
}

func TestRandomTheme_AlwaysValid(t *testing.T) {
	for i := 0; i < 50; i++ {
		theme := RandomTheme(nil)
		if theme.Name == "" {
			t.Error("RandomTheme returned a theme with no name")
		}
		if len(theme.Waypoints) == 0 {
			t.Error("RandomTheme returned a theme with no waypoints")
		}
		if len(theme.ItemPool) == 0 {
			t.Error("RandomTheme returned a theme with no items")
		}
	}
}
