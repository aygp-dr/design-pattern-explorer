# design-pattern-explorer

A terminal UI (TUI) for interactively exploring Gang of Four design patterns. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Install

```bash
go install github.com/aygp-dr/design-pattern-explorer@latest
```

## Run

```bash
go run main.go
```

## Controls

- `j/k` or arrow keys: navigate patterns
- `Enter` or `Space`: toggle detail view
- `1-4`: filter by category (All, Creational, Structural, Behavioral)
- `q`: quit

## Patterns Included

15 GoF patterns: Singleton, Factory Method, Abstract Factory, Builder, Prototype, Adapter, Bridge, Composite, Decorator, Facade, Observer, Strategy, Command, State, Template Method.
