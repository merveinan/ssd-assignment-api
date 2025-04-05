package models

// Config represents a configuration
//type Config struct {
//	ID      string `json:"id"`
//	Actions []struct {
//		Type       string `json:"type"`
//		Selector   string `json:"selector,omitempty"`
//		NewElement string `json:"newElement,omitempty"`
//		Position   string `json:"position,omitempty"`
//		Target     string `json:"target,omitempty"`
//		Element    string `json:"element,omitempty"`
//		OldValue   string `json:"oldValue,omitempty"`
//		NewValue   string `json:"newValue,omitempty"`
//	} `json:"actions"`
//}

type Config struct {
	ID      string   `yaml:"id"`
	Actions []Action `yaml:"actions"`
}

// Action represents a DOM manipulation action
type Action struct {
	Type       string `json:"type" yaml:"type"`                                 // Action type (remove, replace, insert, alter)
	Selector   string `json:"selector,omitempty" yaml:"selector,omitempty"`     // CSS selector (for remove/replace)
	NewElement string `json:"newElement,omitempty" yaml:"newElement,omitempty"` // New HTML element (for replace)
	Position   string `json:"position,omitempty" yaml:"position,omitempty"`     // Position (for insert: before/after)
	Target     string `json:"target,omitempty" yaml:"target,omitempty"`         // Target element (for insert)
	OldValue   string `json:"oldValue,omitempty" yaml:"oldValue,omitempty"`     // Old value (for alter)
	NewValue   string `json:"newValue,omitempty" yaml:"newValue,omitempty"`     // New value (for alter)
}
