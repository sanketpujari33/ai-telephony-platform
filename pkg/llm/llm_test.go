package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextWindowPruning(t *testing.T) {
	cw := NewContextWindow("System Prompt", 4, 500, 1000)

	cw.AddUserTurn("Turn 1", 30)
	cw.AddAssistantTurn("Resp 1", 30)
	cw.AddUserTurn("Turn 2", 30)
	cw.AddAssistantTurn("Resp 2", 30)

	msgs := cw.GetMessages()
	assert.Equal(t, 5, len(msgs)) // 1 system + 4 turns

	// Exceed max turns -> should drop oldest turns
	cw.AddUserTurn("Turn 3", 30)
	msgs = cw.GetMessages()
	assert.True(t, len(msgs) <= 5)
}

func TestContextWindowSummarizeCheck(t *testing.T) {
	cw := NewContextWindow("System", 10, 100, 200)

	cw.AddUserTurn("Turn 1", 40)
	cw.AddAssistantTurn("Resp 1", 40)

	assert.True(t, cw.ShouldSummarize(0.7))  // 80/100 >= 0.7
	assert.False(t, cw.ShouldSummarize(0.9)) // 80/100 < 0.9
}
