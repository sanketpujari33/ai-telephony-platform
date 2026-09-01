package config

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestLoadConfig(t *testing.T) {
	ResetForTest()
	// Create a temporary config file for testing
	configContent := `
sip:
  ingress:
    listen:
      udp: "0.0.0.0:5060"
      tcp: "0.0.0.0:5060"
      tls: "0.0.0.0:5061"
    workers: 4
    trunks:
      - name: "primary"
        host: "sip.example.com"
        port: 5060
        transport: "udp"
        priority: 1
  router:
    grpc:
      listen: "0.0.0.0:9090"
    media_ports:
      start: 10000
      end: 20000
      per_call: 4

media:
  proxy:
    grpc:
      listen: "0.0.0.0:9091"
    mediasoup:
      num_workers: 4
    rtp_port_range:
      start: 10000
      end: 20000
    codecs:
      - name: "opus"
        payload_type: 111
        clock_rate: 48000
        channels: 2

stt:
  providers:
    - name: "deepgram"
      type: "deepgram"
      priority: 1
      enabled: true
      config:
        api_key: "test-key"
        model: "nova-2"
      connection_pool:
        max_connections: 10
        idle_timeout: 300
  circuit_breaker:
    enabled: true
    failure_threshold: 3
    timeout: 30

llm:
  providers:
    - name: "openai"
      type: "openai"
      priority: 1
      enabled: true
      config:
        api_key: "test-key"
        model: "gpt-4o"
      token_counting:
        enabled: true
        model: "gpt-4o"
      connection_pool:
        max_connections: 20
        idle_timeout: 300

tts:
  providers:
    - name: "elevenlabs"
      type: "elevenlabs"
      priority: 1
      enabled: true
      config:
        api_key: "test-key"
        model: "eleven_multilingual_v2"
      connection_pool:
        max_connections: 10
        idle_timeout: 300

campaign:
  manager:
    grpc:
      listen: "0.0.0.0:9092"
    compliance:
      dnc_check_enabled: true
      consent_check_enabled: true

transfer:
  manager:
    grpc:
      listen: "0.0.0.0:9093"

storage:
  postgres:
    host: "localhost"
    port: 5432
    database: "test"
    username: "test"
    password: "test"
  redis:
    host: "localhost"
    port: 6379

observability:
  logging:
    level: "info"
    format: "json"
  metrics:
    enabled: true
    prometheus:
      enabled: true
      listen: "0.0.0.0:9090"
  tracing:
    enabled: true
    otel:
      service_name: "test"
      exporter:
        type: "otlp"
        endpoint: "http://localhost:4317"
`

	// Write temp config file
	tmpFile, err := os.CreateTemp("", "platform-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Create logger
	logger, _ := zap.NewDevelopment()

	// Load config
	cfg, err := Load(tmpFile.Name(), logger)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify config loaded
	if cfg == nil {
		t.Fatal("Config is nil")
	}

	// Verify SIP config
	if cfg.SIP.Ingress.Workers != 4 {
		t.Errorf("Expected workers=4, got %d", cfg.SIP.Ingress.Workers)
	}
	if len(cfg.SIP.Ingress.Trunks) != 1 {
		t.Errorf("Expected 1 trunk, got %d", len(cfg.SIP.Ingress.Trunks))
	}

	// Verify Media config
	if cfg.Media.Proxy.GRPC.Listen != "0.0.0.0:9091" {
		t.Errorf("Expected media proxy listen=0.0.0.0:9091, got %s", cfg.Media.Proxy.GRPC.Listen)
	}

	// Verify STT config
	if len(cfg.STT.Providers) != 1 {
		t.Errorf("Expected 1 STT provider, got %d", len(cfg.STT.Providers))
	}
	if cfg.STT.Providers[0].Name != "deepgram" {
		t.Errorf("Expected STT provider name=deepgram, got %s", cfg.STT.Providers[0].Name)
	}

	// Verify LLM config
	if len(cfg.LLM.Providers) != 1 {
		t.Errorf("Expected 1 LLM provider, got %d", len(cfg.LLM.Providers))
	}
	if cfg.LLM.Providers[0].Name != "openai" {
		t.Errorf("Expected LLM provider name=openai, got %s", cfg.LLM.Providers[0].Name)
	}

	// Verify TTS config
	if len(cfg.TTS.Providers) != 1 {
		t.Errorf("Expected 1 TTS provider, got %d", len(cfg.TTS.Providers))
	}
	if cfg.TTS.Providers[0].Name != "elevenlabs" {
		t.Errorf("Expected TTS provider name=elevenlabs, got %s", cfg.TTS.Providers[0].Name)
	}

	// Verify Storage config
	if cfg.Storage.Postgres.Host != "localhost" {
		t.Errorf("Expected postgres host=localhost, got %s", cfg.Storage.Postgres.Host)
	}
	if cfg.Storage.Redis.Host != "localhost" {
		t.Errorf("Expected redis host=localhost, got %s", cfg.Storage.Redis.Host)
	}

	t.Log("Config loaded successfully")
}

func TestLoadConfigWithEnvOverrides(t *testing.T) {
	ResetForTest()
	// Set environment variables
	os.Setenv("AI_TELEPHONY_SIP_INGRESS_WORKERS", "8")
	os.Setenv("AI_TELEPHONY_STORAGE_POSTGRES_HOST", "db.example.internal")
	defer func() {
		os.Unsetenv("AI_TELEPHONY_SIP_INGRESS_WORKERS")
		os.Unsetenv("AI_TELEPHONY_STORAGE_POSTGRES_HOST")
	}()

	configContent := `
sip:
  ingress:
    listen:
      udp: "0.0.0.0:5060"
    workers: 4
    trunks:
      - name: "primary"
        host: "sip.example.com"
        port: 5060
        transport: "udp"
        priority: 1
  router:
    grpc:
      listen: "0.0.0.0:9090"
    media_ports:
      start: 10000
      end: 20000
      per_call: 4

media:
  proxy:
    grpc:
      listen: "0.0.0.0:9091"
    mediasoup:
      num_workers: 4
    rtp_port_range:
      start: 10000
      end: 20000
    codecs:
      - name: "opus"
        payload_type: 111
        clock_rate: 48000
        channels: 2

stt:
  providers:
    - name: "deepgram"
      type: "deepgram"
      priority: 1
      enabled: true
      config:
        api_key: "file-deepgram-key"
        model: "nova-2"
      connection_pool:
        max_connections: 10
        idle_timeout: 300
  circuit_breaker:
    enabled: true
    failure_threshold: 3
    timeout: 30

llm:
  providers:
    - name: "openai"
      type: "openai"
      priority: 1
      enabled: true
      config:
        api_key: "file-openai-key"
        model: "gpt-4o"
      token_counting:
        enabled: true
        model: "gpt-4o"
      connection_pool:
        max_connections: 20
        idle_timeout: 300

tts:
  providers:
    - name: "elevenlabs"
      type: "elevenlabs"
      priority: 1
      enabled: true
      config:
        api_key: "file-elevenlabs-key"
        model: "eleven_multilingual_v2"
      connection_pool:
        max_connections: 10
        idle_timeout: 300

campaign:
  manager:
    grpc:
      listen: "0.0.0.0:9092"

transfer:
  manager:
    grpc:
      listen: "0.0.0.0:9093"

storage:
  postgres:
    host: "localhost"
    port: 5432
    database: "test"
    username: "test"
    password: "test"
  redis:
    host: "localhost"
    port: 6379

observability:
  logging:
    level: "info"
    format: "json"
  metrics:
    enabled: true
    prometheus:
      enabled: true
      listen: "0.0.0.0:9090"
  tracing:
    enabled: true
    otel:
      service_name: "test"
      exporter:
        type: "otlp"
        endpoint: "http://localhost:4317"
`

	tmpFile, err := os.CreateTemp("", "platform-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	logger, _ := zap.NewDevelopment()

	cfg, err := Load(tmpFile.Name(), logger)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify environment variable override worked
	if cfg.SIP.Ingress.Workers != 8 {
		t.Errorf("Expected workers=8 (from env), got %d", cfg.SIP.Ingress.Workers)
	}

	if cfg.Storage.Postgres.Host != "db.example.internal" {
		t.Errorf("Expected postgres host=db.example.internal, got %v", cfg.Storage.Postgres.Host)
	}

	t.Log("Environment variable overrides work correctly")
}

func TestValidateConfig(t *testing.T) {
	ResetForTest()
	logger, _ := zap.NewDevelopment()

	// Test missing SIP listen addresses
	configContent := `
sip:
  ingress:
    listen: {}
    workers: 4
    trunks:
      - name: "primary"
        host: "sip.example.com"
        port: 5060
        transport: "udp"
        priority: 1
  router:
    grpc:
      listen: "0.0.0.0:9090"
    media_ports:
      start: 10000
      end: 20000
      per_call: 4

media:
  proxy:
    grpc:
      listen: "0.0.0.0:9091"
    mediasoup:
      num_workers: 4
    rtp_port_range:
      start: 10000
      end: 20000
    codecs:
      - name: "opus"
        payload_type: 111
        clock_rate: 48000
        channels: 2

stt:
  providers:
    - name: "deepgram"
      type: "deepgram"
      priority: 1
      enabled: true
      config:
        api_key: "test-key"
        model: "nova-2"
      connection_pool:
        max_connections: 10
        idle_timeout: 300
  circuit_breaker:
    enabled: true
    failure_threshold: 3
    timeout: 30

llm:
  providers:
    - name: "openai"
      type: "openai"
      priority: 1
      enabled: true
      config:
        api_key: "test-key"
        model: "gpt-4o"
      token_counting:
        enabled: true
        model: "gpt-4o"
      connection_pool:
        max_connections: 20
        idle_timeout: 300

tts:
  providers:
    - name: "elevenlabs"
      type: "elevenlabs"
      priority: 1
      enabled: true
      config:
        api_key: "test-key"
        model: "eleven_multilingual_v2"
      connection_pool:
        max_connections: 10
        idle_timeout: 300

campaign:
  manager:
    grpc:
      listen: "0.0.0.0:9092"

transfer:
  manager:
    grpc:
      listen: "0.0.0.0:9093"

storage:
  postgres:
    host: "localhost"
    port: 5432
    database: "test"
    username: "test"
    password: "test"
  redis:
    host: "localhost"
    port: 6379

observability:
  logging:
    level: "info"
    format: "json"
  metrics:
    enabled: true
    prometheus:
      enabled: true
      listen: "0.0.0.0:9090"
  tracing:
    enabled: true
    otel:
      service_name: "test"
      exporter:
        type: "otlp"
        endpoint: "http://localhost:4317"
`

	tmpFile, err := os.CreateTemp("", "platform-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	_, err = Load(tmpFile.Name(), logger)
	if err == nil {
		t.Error("Expected error for missing SIP listen addresses, got nil")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}