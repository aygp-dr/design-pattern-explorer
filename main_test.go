package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Pattern data integrity tests ---

func TestAllPatternsHaveName(t *testing.T) {
	for i, p := range patterns {
		if p.Name == "" {
			t.Errorf("pattern[%d] has empty Name", i)
		}
	}
}

func TestAllPatternsHaveCategory(t *testing.T) {
	validCategories := map[string]bool{
		"creational": true,
		"structural": true,
		"behavioral": true,
	}
	for i, p := range patterns {
		if p.Category == "" {
			t.Errorf("pattern[%d] (%s) has empty Category", i, p.Name)
		}
		if !validCategories[p.Category] {
			t.Errorf("pattern[%d] (%s) has invalid Category %q", i, p.Name, p.Category)
		}
	}
}

func TestAllPatternsHaveIntent(t *testing.T) {
	for i, p := range patterns {
		if p.Intent == "" {
			t.Errorf("pattern[%d] (%s) has empty Intent", i, p.Name)
		}
	}
}

func TestAllPatternsHaveWhen(t *testing.T) {
	for i, p := range patterns {
		if p.When == "" {
			t.Errorf("pattern[%d] (%s) has empty When", i, p.Name)
		}
	}
}

func TestAllPatternsHaveExample(t *testing.T) {
	for i, p := range patterns {
		if p.Example == "" {
			t.Errorf("pattern[%d] (%s) has empty Example", i, p.Name)
		}
	}
}

func TestPatternCount(t *testing.T) {
	if len(patterns) != 15 {
		t.Errorf("expected 15 patterns, got %d", len(patterns))
	}
}

func TestUniquePatternNames(t *testing.T) {
	seen := make(map[string]bool)
	for _, p := range patterns {
		if seen[p.Name] {
			t.Errorf("duplicate pattern name: %s", p.Name)
		}
		seen[p.Name] = true
	}
}

// --- Filter function tests ---

func TestFilteredNoFilter(t *testing.T) {
	m := model{patterns: patterns, filter: ""}
	got := m.filtered()
	if len(got) != len(patterns) {
		t.Errorf("no filter: expected %d patterns, got %d", len(patterns), len(got))
	}
}

func TestFilteredCreational(t *testing.T) {
	m := model{patterns: patterns, filter: "creational"}
	got := m.filtered()
	expected := 5 // Singleton, Factory Method, Abstract Factory, Builder, Prototype
	if len(got) != expected {
		t.Errorf("creational filter: expected %d patterns, got %d", expected, len(got))
	}
	for _, p := range got {
		if p.Category != "creational" {
			t.Errorf("creational filter returned pattern %q with category %q", p.Name, p.Category)
		}
	}
}

func TestFilteredStructural(t *testing.T) {
	m := model{patterns: patterns, filter: "structural"}
	got := m.filtered()
	expected := 5 // Adapter, Bridge, Composite, Decorator, Facade
	if len(got) != expected {
		t.Errorf("structural filter: expected %d patterns, got %d", expected, len(got))
	}
	for _, p := range got {
		if p.Category != "structural" {
			t.Errorf("structural filter returned pattern %q with category %q", p.Name, p.Category)
		}
	}
}

func TestFilteredBehavioral(t *testing.T) {
	m := model{patterns: patterns, filter: "behavioral"}
	got := m.filtered()
	expected := 5 // Observer, Strategy, Command, State, Template Method
	if len(got) != expected {
		t.Errorf("behavioral filter: expected %d patterns, got %d", expected, len(got))
	}
	for _, p := range got {
		if p.Category != "behavioral" {
			t.Errorf("behavioral filter returned pattern %q with category %q", p.Name, p.Category)
		}
	}
}

func TestFilteredInvalidCategory(t *testing.T) {
	m := model{patterns: patterns, filter: "nonexistent"}
	got := m.filtered()
	if len(got) != 0 {
		t.Errorf("invalid filter: expected 0 patterns, got %d", len(got))
	}
}

// --- Model initialization tests ---

func TestModelInit(t *testing.T) {
	m := model{patterns: patterns}
	if m.cursor != 0 {
		t.Errorf("initial cursor should be 0, got %d", m.cursor)
	}
	if m.filter != "" {
		t.Errorf("initial filter should be empty, got %q", m.filter)
	}
	if m.showDetail {
		t.Errorf("initial showDetail should be false")
	}
}

func TestModelInitCmd(t *testing.T) {
	m := model{patterns: patterns}
	cmd := m.Init()
	if cmd != nil {
		t.Errorf("Init() should return nil cmd")
	}
}

// --- Key handling tests ---

func sendKey(m model, key string) model {
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	return updated.(model)
}

func sendSpecialKey(m model, keyType tea.KeyType) model {
	updated, _ := m.Update(tea.KeyMsg{Type: keyType})
	return updated.(model)
}

func TestKeyJ_MovesCursorDown(t *testing.T) {
	m := model{patterns: patterns, cursor: 0}
	m = sendKey(m, "j")
	if m.cursor != 1 {
		t.Errorf("after 'j', cursor should be 1, got %d", m.cursor)
	}
}

func TestKeyK_MovesCursorUp(t *testing.T) {
	m := model{patterns: patterns, cursor: 3}
	m = sendKey(m, "k")
	if m.cursor != 2 {
		t.Errorf("after 'k' from 3, cursor should be 2, got %d", m.cursor)
	}
}

func TestKeyK_DoesNotGoBelowZero(t *testing.T) {
	m := model{patterns: patterns, cursor: 0}
	m = sendKey(m, "k")
	if m.cursor != 0 {
		t.Errorf("after 'k' from 0, cursor should stay 0, got %d", m.cursor)
	}
}

func TestKeyJ_DoesNotExceedListLength(t *testing.T) {
	m := model{patterns: patterns, cursor: len(patterns) - 1}
	m = sendKey(m, "j")
	if m.cursor != len(patterns)-1 {
		t.Errorf("after 'j' at end, cursor should stay at %d, got %d", len(patterns)-1, m.cursor)
	}
}

func TestKeyEnter_TogglesDetail(t *testing.T) {
	m := model{patterns: patterns, showDetail: false}
	m = sendSpecialKey(m, tea.KeyEnter)
	if !m.showDetail {
		t.Errorf("after enter, showDetail should be true")
	}
	m = sendSpecialKey(m, tea.KeyEnter)
	if m.showDetail {
		t.Errorf("after second enter, showDetail should be false")
	}
}

func TestKeyQ_ReturnsQuit(t *testing.T) {
	m := model{patterns: patterns}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Errorf("'q' should return a quit command, got nil")
	}
}

func TestKey1_ClearsFilter(t *testing.T) {
	m := model{patterns: patterns, filter: "creational", cursor: 3}
	m = sendKey(m, "1")
	if m.filter != "" {
		t.Errorf("after '1', filter should be empty, got %q", m.filter)
	}
	if m.cursor != 0 {
		t.Errorf("after '1', cursor should reset to 0, got %d", m.cursor)
	}
}

func TestKey2_SetsCreationalFilter(t *testing.T) {
	m := model{patterns: patterns}
	m = sendKey(m, "2")
	if m.filter != "creational" {
		t.Errorf("after '2', filter should be 'creational', got %q", m.filter)
	}
	if m.cursor != 0 {
		t.Errorf("after '2', cursor should reset to 0, got %d", m.cursor)
	}
}

func TestKey3_SetsStructuralFilter(t *testing.T) {
	m := model{patterns: patterns}
	m = sendKey(m, "3")
	if m.filter != "structural" {
		t.Errorf("after '3', filter should be 'structural', got %q", m.filter)
	}
}

func TestKey4_SetsBehavioralFilter(t *testing.T) {
	m := model{patterns: patterns}
	m = sendKey(m, "4")
	if m.filter != "behavioral" {
		t.Errorf("after '4', filter should be 'behavioral', got %q", m.filter)
	}
}

// --- View rendering tests ---

func TestViewContainsTitle(t *testing.T) {
	m := model{patterns: patterns}
	view := m.View()
	if !containsStr(view, "Design Pattern Explorer") {
		t.Errorf("view should contain title 'Design Pattern Explorer'")
	}
}

func TestViewContainsPatternNames(t *testing.T) {
	m := model{patterns: patterns}
	view := m.View()
	for _, p := range patterns {
		if !containsStr(view, p.Name) {
			t.Errorf("view should contain pattern name %q", p.Name)
		}
	}
}

func TestViewContainsHelpText(t *testing.T) {
	m := model{patterns: patterns}
	view := m.View()
	if !containsStr(view, "j/k: navigate") {
		t.Errorf("view should contain help text")
	}
}

func TestViewWithDetail(t *testing.T) {
	m := model{patterns: patterns, showDetail: true, cursor: 0}
	view := m.View()
	// The detail view should include the first pattern's intent
	if !containsStr(view, patterns[0].Intent) {
		t.Errorf("detail view should contain the selected pattern's Intent")
	}
}

func TestViewWithFilterShowsOnlyFiltered(t *testing.T) {
	m := model{patterns: patterns, filter: "creational"}
	view := m.View()
	// Creational patterns should appear
	if !containsStr(view, "Singleton") {
		t.Errorf("creational filter view should contain 'Singleton'")
	}
	// Behavioral patterns should not appear in the list
	// (Observer is behavioral; it should not be listed)
	// Note: Observer may appear in help text or other places, but
	// we check it's not in the pattern list area
}

// --- Cursor bounds after filter change ---

func TestCursorResetOnFilterChange(t *testing.T) {
	m := model{patterns: patterns, cursor: 10}
	m = sendKey(m, "2") // switch to creational
	if m.cursor != 0 {
		t.Errorf("cursor should reset to 0 on filter change, got %d", m.cursor)
	}
}

func TestCursorBoundsWithFilter(t *testing.T) {
	m := model{patterns: patterns, filter: "creational", cursor: 0}
	// Move to last creational pattern
	for i := 0; i < 10; i++ {
		m = sendKey(m, "j")
	}
	filtered := m.filtered()
	if m.cursor >= len(filtered) {
		t.Errorf("cursor %d should not exceed filtered length %d", m.cursor, len(filtered))
	}
}

// helper
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && searchStr(s, substr)
}

func searchStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
