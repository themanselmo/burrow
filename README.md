# Burrow

A virtual pet that lives in your terminal. Adopt a seal, keep it happy, and send it on focus missions while you get things done.

```
  ▄███▄  
 █ ◉ ◉ █ 
 █  ▲  █ 
  ▀███▀  
 ╱██ ██╲ 
```

## Features

- **Pixel art pet** rendered in Unicode/ANSI — lives right in your terminal
- **Core care loop** — feed, pet, and play to keep your pet happy and earn XP
- **Idle animations** — your pet blinks, looks around, and wanders on its own
- **Focus missions** — send your pet on a 25min, 45min, or custom timed mission while you work. Missions complete in real time, even with the terminal closed
- **XP & leveling** — interactions and missions earn XP; longer missions give more, capped to keep things fair
- **Item drops** — rare finds from missions, exclusive to each theme
- **Pet log** — every pet you've owned, recorded locally
- **Offline-first** — all state lives in `~/.burrow/`. No account required to use

## Install

```bash
go install github.com/themanselmo/burrow/cmd/burrow@latest
```

Or build from source:

```bash
git clone https://github.com/themanselmo/burrow.git
cd burrow
go build ./cmd/burrow
./burrow
```

> Requires a terminal with Unicode and 256-color support. Works great in iTerm2, Ghostty, WezTerm, and most modern terminals.

## Usage

```
burrow
```

On first launch you'll be prompted to name your seal. After that, everything is keyboard driven:

| Key | Action       |
|-----|--------------|
| F   | Feed         |
| P   | Pet          |
| A   | Play         |
| M   | Send on mission |
| I   | View items   |
| L   | Pet log      |
| R   | Release pet  |
| Q   | Quit         |

## Missions

Missions are a built-in focus timer. Send your pet on a walk or expedition, put your head down, and it'll be home when you're done.

- **25 min** → 100 XP
- **45 min** → 180 XP
- **Custom** → scales up to a cap of 250 XP (hit around 62 min)
- Small chance of a rare item drop on every completed mission

## Data

All data is stored locally in `~/.burrow/`:

| File | Contents |
|------|----------|
| `pet.json` | Active pet state, current mission |
| `log.json` | History of past pets |

## Contributing

This is an early-stage personal project. Issues and ideas welcome.

## License

MIT
