package tts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTTSCacheGetSet(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cache := NewTTSCache(logger)

	prompt := "Hello, thank you for calling."
	voiceID := "eleven_multilingual_v2"
	audioData := []byte{0x00, 0x01, 0x02, 0x03}

	cache.Set(prompt, voiceID, audioData)

	retrieved, exists := cache.Get(prompt, voiceID)
	assert.True(t, exists)
	assert.Equal(t, audioData, retrieved)

	_, exists = cache.Get("nonexistent", voiceID)
	assert.False(t, exists)
}
