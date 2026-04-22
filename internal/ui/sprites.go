package ui

import "math/rand"

// Each animation is a sequence of frames ([][]string).
// Frames are cycled by the animation ticker in home.go.

// --- Base frames ---

var sealAwake = []string{
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealBlink1 = []string{
	`  ▄███▄  `,
	` █ ◔ ◔ █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealBlink2 = []string{
	`  ▄███▄  `,
	` █ ─ ─ █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealLookLeft = []string{
	`  ▄███▄  `,
	` █ ◉◉  █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealLookRight = []string{
	`  ▄███▄  `,
	` █  ◉◉ █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealNoseWiggle = []string{
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ω  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealBounceUp = []string{
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`         `,
}

var sealBounceDown = []string{
	`         `,
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
}

var sealHappy = []string{
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ▲  █ `,
	`  ▀███▀  `,
	` ╱██ ██╲ `,
}

var sealSad = []string{
	`  ▄███▄  `,
	` █ ╥ ╥ █ `,
	` █  ▽  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealAngry = []string{
	`  ▄███▄  `,
	` █ ◣ ◢ █ `,
	` █  ▽  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealHungry = []string{
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ω  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealSleepy = []string{
	`  ▄███▄  `,
	` █ - - █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

// Sleeping frames — z's float diagonally up and to the right from the head.
var sealAsleep1 = []string{
	`  ▄███▄  `,
	` █ ─ ─ █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealAsleep2 = []string{
	`  ▄███▄  `,
	` █ ─ ─ █z`,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealAsleep3 = []string{
	`  ▄███▄ z`,
	` █ ─ ─ █z`,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealAsleep4 = []string{
	`  ▄███▄Z `,
	` █ ─ ─ █z`,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealAsleep5 = []string{
	` Z▄███▄  `,
	` █ ─ ─ █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

var sealWalk1 = []string{
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ▲  █ `,
	`  ▀███▀  `,
	`  ▐█ ██  `,
}

var sealWalk2 = []string{
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ▲  █ `,
	`  ▀███▀  `,
	`  ██ █▌  `,
}

var sealWalk3 = []string{
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ▲  █ `,
	`  ▀███▀  `,
	` ╱█  ██  `,
}

var sealWalk4 = []string{
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ▲  █ `,
	`  ▀███▀  `,
	`  ██  █╲ `,
}

var sealWave = []string{
	`  ▄███▄  `,
	` █ ◉ ◉ █ `,
	` █  ▲  █o`,
	`  ▀███▀  `,
	` ╱██ ██╲ `,
}

var sealWaveHigh = []string{
	`  ▄███▄ o`,
	` █ ◉ ◉ █ `,
	` █  ▲  █ `,
	`  ▀███▀  `,
	` ╱██ ██╲ `,
}

var sealTiny = []string{
	`         `,
	`  ▄█▄    `,
	` █◉◉█    `,
	`  ▀█▀    `,
	`  ██     `,
}

var sealGone = []string{
	`         `,
	`         `,
	`         `,
	`         `,
	`         `,
}

var sealSleepyNod = []string{
	`  ▄███▄  `,
	` █ ─ ─ █ `,
	` █  ▿  █ `,
	`  ▀███▀  `,
	`  ██ ██  `,
}

// --- Animation sequences ---
// Each sequence plays once then returns to idle.

type Animation [][]string

var AnimIdle = Animation{sealAwake}

var AnimBlink = Animation{
	sealAwake, sealBlink1, sealBlink2, sealBlink1, sealAwake,
}

var AnimLook = Animation{
	sealAwake, sealLookLeft, sealLookLeft, sealAwake,
	sealLookRight, sealLookRight, sealAwake,
}

var AnimNoseWiggle = Animation{
	sealAwake, sealNoseWiggle, sealAwake, sealNoseWiggle, sealAwake,
}

var AnimBounce = Animation{
	sealAwake, sealBounceUp, sealAwake, sealBounceDown, sealAwake,
}

var AnimSleepyNod = Animation{
	sealSleepy, sealSleepyNod, sealSleepy, sealSleepyNod, sealSleepy,
}

var AnimWalk = Animation{
	sealWalk1, sealWalk2, sealWalk3, sealWalk4,
}

var AnimAsleep = Animation{
	sealAsleep1, sealAsleep1, sealAsleep2,
	sealAsleep3, sealAsleep4, sealAsleep5, sealAsleep1,
}

var AnimFarewell = Animation{
	sealHappy, sealWave, sealWaveHigh, sealWave,
	sealWaveHigh, sealWave, sealTiny, sealTiny, sealGone,
}

var AnimAbruptQuit = Animation{
	sealAngry, sealAngry, sealSad, sealSad,
}

// idleAnims are the pool randomly chosen from during idle.
var idleAnims = []Animation{AnimBlink, AnimBlink, AnimLook, AnimNoseWiggle, AnimBounce}

func RandomIdleAnim() Animation {
	return idleAnims[rand.Intn(len(idleAnims))]
}

// SpriteFor returns the right base sprite for a given mood, ignoring animation state.
// Used as the fallback when no animation is playing.
func SpriteFor(mood string, sleepy bool) []string {
	if sleepy {
		return sealSleepy
	}
	switch mood {
	case "happy":
		return sealHappy
	case "sad":
		return sealSad
	case "hungry":
		return sealHungry
	default:
		return sealAwake
	}
}
