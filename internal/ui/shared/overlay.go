package shared

import "charm.land/lipgloss/v2"

// RenderOverlay wraps content in the standard overlay style (rounded border).
func RenderOverlay(content string) string {
	return OverlayStyle.Render(content)
}

// PlaceOverlay centers an overlay on screen, replacing the background view.
// Uses full terminal height minus tab bar (1) and help bar (1) for vertical centering.
func PlaceOverlay(width, height int, overlay string) string {
	contentHeight := height - 2 // account for tab bar + help bar
	if contentHeight < 10 {
		contentHeight = 10
	}
	return lipgloss.Place(width, contentHeight, lipgloss.Center, lipgloss.Center, overlay)
}
