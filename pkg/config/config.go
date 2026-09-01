package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config holds all configuration for the platform
type Config struct {
	SIP           SIPConfig           `mapstructure:"sip"`
	Media         MediaConfig         `mapstructure:"media"`
	STT           STTConfig           `mapstructure:"stt"`
	LLM           LLMConfig           `mapstructure:"llm"`
	TTS           TTSConfig           `mapstructure:"tts"`
	Campaign      CampaignConfig      `mapstructure:"campaign"`
	Transfer      TransferConfig      `mapstructure:"transfer"`
	Storage       StorageConfig       `mapstructure:"storage"`
	Observability ObservabilityConfig `mapstructure:"observability"`
}

// SIPConfig holds SIP-related configuration
type SIPConfig struct {
	Ingress SIPIngressConfig `mapstructure:"ingress"`
	Router  SIPRouterConfig  `mapstructure:"router"`
}

// SIPIngressConfig holds SIP Ingress (Kamailio) configuration
type SIPIngressConfig struct {
	Listen      SIPListenConfig    `mapstructure:"listen"`
	TLS         SIPTLSConfig       `mapstructure:"tls"`
	Workers     int                `mapstructure:"workers"`
	MaxRequests int                `mapstructure:"max_requests_per_worker"`
	Trunks      []SIPTrunkConfig   `mapstructure:"trunks"`
	DNSSRV      SIPDNSSRVConfig    `mapstructure:"dns_srv"`
	Keepalive   SIPKeepaliveConfig `mapstructure:"keepalive"`
	RateLimit   SIPRateLimitConfig `mapstructure:"rate_limit"`
}

// SIPListenConfig holds SIP listen addresses
type SIPListenConfig struct {
	UDP string `mapstructure:"udp"`
	TCP string `mapstructure:"tcp"`
	TLS string `mapstructure:"tls"`
	WS  string `mapstructure:"ws"`
	WSS string `mapstructure:"wss"`
}

// SIPTLSConfig holds TLS configuration
type SIPTLSConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	CertFile     string `mapstructure:"cert_file"`
	KeyFile      string `mapstructure:"key_file"`
	CAFile       string `mapstructure:"ca_file"`
	VerifyClient bool   `mapstructure:"verify_client"`
}

// SIPTrunkConfig holds SIP trunk configuration
type SIPTrunkConfig struct {
	Name      string `mapstructure:"name"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Transport string `mapstructure:"transport"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	Realm     string `mapstructure:"realm"`
	Priority  int    `mapstructure:"priority"`
	Weight    int    `mapstructure:"weight"`
}

// SIPDNSSRVConfig holds DNS SRV failover configuration
type SIPDNSSRVConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Domain    string `mapstructure:"domain"`
	Transport string `mapstructure:"transport"`
}

// SIPKeepaliveConfig holds OPTIONS keepalive configuration
type SIPKeepaliveConfig struct {
	Enabled    bool `mapstructure:"enabled"`
	Interval   int  `mapstructure:"interval"`
	Timeout    int  `mapstructure:"timeout"`
	RetryCount int  `mapstructure:"retry_count"`
}

// SIPRateLimitConfig holds rate limiting configuration
type SIPRateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerSecond int  `mapstructure:"requests_per_second"`
	Burst             int  `mapstructure:"burst"`
}

// SIPRouterConfig holds Call Router configuration
type SIPRouterConfig struct {
	GRPC         SIPRouterGRPCConfig         `mapstructure:"grpc"`
	TenantLookup SIPRouterTenantLookupConfig `mapstructure:"tenant_lookup"`
	MediaPorts   SIPRouterMediaPortsConfig   `mapstructure:"media_ports"`
}

// SIPRouterGRPCConfig holds gRPC server configuration
type SIPRouterGRPCConfig struct {
	Listen         string `mapstructure:"listen"`
	MaxRecvMsgSize int    `mapstructure:"max_recv_msg_size"`
	MaxSendMsgSize int    `mapstructure:"max_send_msg_size"`
}

// SIPRouterTenantLookupConfig holds tenant lookup configuration
type SIPRouterTenantLookupConfig struct {
	CacheTTL  int `mapstructure:"cache_ttl"`
	CacheSize int `mapstructure:"cache_size"`
}

// SIPRouterMediaPortsConfig holds media port allocation configuration
type SIPRouterMediaPortsConfig struct {
	Start   int `mapstructure:"start"`
	End     int `mapstructure:"end"`
	PerCall int `mapstructure:"per_call"`
}

// MediaConfig holds media-related configuration
type MediaConfig struct {
	Proxy      MediaProxyConfig      `mapstructure:"proxy"`
	Transcoder MediaTranscoderConfig `mapstructure:"transcoder"`
}

// MediaProxyConfig holds Media Proxy configuration
type MediaProxyConfig struct {
	GRPC           MediaProxyGRPCConfig           `mapstructure:"grpc"`
	MediaSoup      MediaProxyMediaSoupConfig      `mapstructure:"mediasoup"`
	RTPPortRange   MediaProxyRTPPortRangeConfig   `mapstructure:"rtp_port_range"`
	Codecs         []MediaProxyCodecConfig        `mapstructure:"codecs"`
	JitterBuffer   MediaProxyJitterBufferConfig   `mapstructure:"jitter_buffer"`
	PLC            MediaProxyPLCConfig            `mapstructure:"plc"`
	DTMF           MediaProxyDTMFConfig           `mapstructure:"dtmf"`
	QualityMetrics MediaProxyQualityMetricsConfig `mapstructure:"quality_metrics"`
}

// MediaProxyGRPCConfig holds gRPC server configuration
type MediaProxyGRPCConfig struct {
	Listen         string `mapstructure:"listen"`
	MaxRecvMsgSize int    `mapstructure:"max_recv_msg_size"`
	MaxSendMsgSize int    `mapstructure:"max_send_msg_size"`
}

// MediaProxyMediaSoupConfig holds MediaSoup worker configuration
type MediaProxyMediaSoupConfig struct {
	NumWorkers     int                               `mapstructure:"num_workers"`
	WorkerSettings MediaProxyMediaSoupWorkerSettings `mapstructure:"worker_settings"`
}

// MediaProxyMediaSoupWorkerSettings holds MediaSoup worker settings
type MediaProxyMediaSoupWorkerSettings struct {
	RTCMinPort          int    `mapstructure:"rtc_min_port"`
	RTCMaxPort          int    `mapstructure:"rtc_max_port"`
	DTLSCertificateFile string `mapstructure:"dtls_certificate_file"`
	DTLSPrivateKeyFile  string `mapstructure:"dtls_private_key_file"`
}

// MediaProxyRTPPortRangeConfig holds RTP port range configuration
type MediaProxyRTPPortRangeConfig struct {
	Start int `mapstructure:"start"`
	End   int `mapstructure:"end"`
}

// MediaProxyCodecConfig holds codec configuration
type MediaProxyCodecConfig struct {
	Name        string            `mapstructure:"name"`
	PayloadType int               `mapstructure:"payload_type"`
	ClockRate   int               `mapstructure:"clock_rate"`
	Channels    int               `mapstructure:"channels"`
	Parameters  map[string]string `mapstructure:"parameters"`
}

// MediaProxyJitterBufferConfig holds jitter buffer configuration
type MediaProxyJitterBufferConfig struct {
	Enabled     bool `mapstructure:"enabled"`
	MinDelay    int  `mapstructure:"min_delay"`
	MaxDelay    int  `mapstructure:"max_delay"`
	TargetDelay int  `mapstructure:"target_delay"`
	Adaptive    bool `mapstructure:"adaptive"`
}

// MediaProxyPLCConfig holds packet loss concealment configuration
type MediaProxyPLCConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Algorithm string `mapstructure:"algorithm"`
}

// MediaProxyDTMFConfig holds DTMF handling configuration
type MediaProxyDTMFConfig struct {
	RFC4733Enabled bool `mapstructure:"rfc4733_enabled"`
	InbandEnabled  bool `mapstructure:"inband_enabled"`
	SIPInfoEnabled bool `mapstructure:"sip_info_enabled"`
}

// MediaProxyQualityMetricsConfig holds media quality metrics configuration
type MediaProxyQualityMetricsConfig struct {
	Enabled          bool `mapstructure:"enabled"`
	Interval         int  `mapstructure:"interval"`
	EmitToPrometheus bool `mapstructure:"emit_to_prometheus"`
}

// MediaTranscoderConfig holds FFmpeg transcoder configuration
type MediaTranscoderConfig struct {
	Enabled   bool                     `mapstructure:"enabled"`
	Image     string                   `mapstructure:"image"`
	Resources MediaTranscoderResources `mapstructure:"resources"`
}

// MediaTranscoderResources holds resource limits/requests
type MediaTranscoderResources struct {
	Limits   MediaTranscoderResourceLimits   `mapstructure:"limits"`
	Requests MediaTranscoderResourceRequests `mapstructure:"requests"`
}

// MediaTranscoderResourceLimits holds resource limits
type MediaTranscoderResourceLimits struct {
	CPU    string `mapstructure:"cpu"`
	Memory string `mapstructure:"memory"`
}

// MediaTranscoderResourceRequests holds resource requests
type MediaTranscoderResourceRequests struct {
	CPU    string `mapstructure:"cpu"`
	Memory string `mapstructure:"memory"`
}

// STTConfig holds STT-related configuration
type STTConfig struct {
	Providers      []STTProviderConfig     `mapstructure:"providers"`
	CircuitBreaker STTCircuitBreakerConfig `mapstructure:"circuit_breaker"`
	Global         STTGlobalConfig         `mapstructure:"global"`
}

// STTProviderConfig holds STT provider configuration
type STTProviderConfig struct {
	Name           string                  `mapstructure:"name"`
	Type           string                  `mapstructure:"type"`
	Priority       int                     `mapstructure:"priority"`
	Enabled        bool                    `mapstructure:"enabled"`
	Config         map[string]interface{}  `mapstructure:"config"`
	ConnectionPool STTConnectionPoolConfig `mapstructure:"connection_pool"`
	Reconnection   STTReconnectionConfig   `mapstructure:"reconnection"`
}

// STTConnectionPoolConfig holds connection pool configuration
type STTConnectionPoolConfig struct {
	MaxConnections int `mapstructure:"max_connections"`
	IdleTimeout    int `mapstructure:"idle_timeout"`
}

// STTReconnectionConfig holds reconnection configuration
type STTReconnectionConfig struct {
	MaxRetries int `mapstructure:"max_retries"`
	BaseDelay  int `mapstructure:"base_delay"`
	MaxDelay   int `mapstructure:"max_delay"`
}

// STTCircuitBreakerConfig holds circuit breaker configuration
type STTCircuitBreakerConfig struct {
	Enabled          bool `mapstructure:"enabled"`
	FailureThreshold int  `mapstructure:"failure_threshold"`
	Timeout          int  `mapstructure:"timeout"`
	HalfOpenRequests int  `mapstructure:"half_open_requests"`
	MetricsEnabled   bool `mapstructure:"metrics_enabled"`
}

// STTGlobalConfig holds global STT settings
type STTGlobalConfig struct {
	InterimResults  bool                     `mapstructure:"interim_results"`
	KeywordBoosting bool                     `mapstructure:"keyword_boosting"`
	Diarization     bool                     `mapstructure:"diarization"`
	LatencyTracking STTLatencyTrackingConfig `mapstructure:"latency_tracking"`
	Logging         STTLoggingConfig         `mapstructure:"logging"`
}

// STTLatencyTrackingConfig holds latency tracking configuration
type STTLatencyTrackingConfig struct {
	Enabled     bool  `mapstructure:"enabled"`
	Percentiles []int `mapstructure:"percentiles"`
}

// STTLoggingConfig holds logging configuration
type STTLoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// LLMConfig holds LLM-related configuration
type LLMConfig struct {
	Providers []LLMProviderConfig `mapstructure:"providers"`
	Context   LLMContextConfig    `mapstructure:"context"`
	Router    LLMRouterConfig     `mapstructure:"router"`
}

// LLMProviderConfig holds LLM provider configuration
type LLMProviderConfig struct {
	Name           string                  `mapstructure:"name"`
	Type           string                  `mapstructure:"type"`
	Priority       int                     `mapstructure:"priority"`
	Enabled        bool                    `mapstructure:"enabled"`
	Config         map[string]interface{}  `mapstructure:"config"`
	TokenCounting  LLMTokenCountingConfig  `mapstructure:"token_counting"`
	ConnectionPool LLMConnectionPoolConfig `mapstructure:"connection_pool"`
	Retry          LLMRetryConfig          `mapstructure:"retry"`
}

// LLMTokenCountingConfig holds token counting configuration
type LLMTokenCountingConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Model   string `mapstructure:"model"`
}

// LLMConnectionPoolConfig holds connection pool configuration
type LLMConnectionPoolConfig struct {
	MaxConnections int `mapstructure:"max_connections"`
	IdleTimeout    int `mapstructure:"idle_timeout"`
}

// LLMRetryConfig holds retry configuration
type LLMRetryConfig struct {
	MaxRetries int `mapstructure:"max_retries"`
	BaseDelay  int `mapstructure:"base_delay"`
	MaxDelay   int `mapstructure:"max_delay"`
}

// LLMContextConfig holds context management configuration
type LLMContextConfig struct {
	SlidingWindow LLMSlidingWindowConfig `mapstructure:"sliding_window"`
	TokenBudget   LLMTokenBudgetConfig   `mapstructure:"token_budget"`
	Summarization LLMSummarizationConfig `mapstructure:"summarization"`
}

// LLMSlidingWindowConfig holds sliding window configuration
type LLMSlidingWindowConfig struct {
	Enabled   bool `mapstructure:"enabled"`
	MaxTurns  int  `mapstructure:"max_turns"`
	MaxTokens int  `mapstructure:"max_tokens"`
}

// LLMTokenBudgetConfig holds token budget configuration
type LLMTokenBudgetConfig struct {
	Enabled        bool    `mapstructure:"enabled"`
	TotalBudget    int     `mapstructure:"total_budget"`
	SystemReserved float64 `mapstructure:"system_reserved"`
	UserReserved   float64 `mapstructure:"user_reserved"`
}

// LLMSummarizationConfig holds summarization configuration
type LLMSummarizationConfig struct {
	Enabled          bool    `mapstructure:"enabled"`
	TriggerThreshold float64 `mapstructure:"trigger_threshold"`
	Model            string  `mapstructure:"model"`
	MaxSummaryTokens int     `mapstructure:"max_summary_tokens"`
	Interval         int     `mapstructure:"interval"`
}

// LLMRouterConfig holds provider router configuration
type LLMRouterConfig struct {
	Enabled         bool                           `mapstructure:"enabled"`
	Strategy        string                         `mapstructure:"strategy"`
	CircuitBreaker  LLMRouterCircuitBreakerConfig  `mapstructure:"circuit_breaker"`
	CostTracking    LLMRouterCostTrackingConfig    `mapstructure:"cost_tracking"`
	ResponseCaching LLMRouterResponseCachingConfig `mapstructure:"response_caching"`
	LatencyMetrics  LLMRouterLatencyMetricsConfig  `mapstructure:"latency_metrics"`
}

// LLMRouterCircuitBreakerConfig holds router circuit breaker configuration
type LLMRouterCircuitBreakerConfig struct {
	Enabled          bool `mapstructure:"enabled"`
	FailureThreshold int  `mapstructure:"failure_threshold"`
	Timeout          int  `mapstructure:"timeout"`
}

// LLMRouterCostTrackingConfig holds cost tracking configuration
type LLMRouterCostTrackingConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Currency string `mapstructure:"currency"`
}

// LLMRouterResponseCachingConfig holds response caching configuration
type LLMRouterResponseCachingConfig struct {
	Enabled bool `mapstructure:"enabled"`
	TTL     int  `mapstructure:"ttl"`
	MaxSize int  `mapstructure:"max_size"`
}

// LLMRouterLatencyMetricsConfig holds latency metrics configuration
type LLMRouterLatencyMetricsConfig struct {
	Enabled     bool  `mapstructure:"enabled"`
	Percentiles []int `mapstructure:"percentiles"`
}

// TTSConfig holds TTS-related configuration
type TTSConfig struct {
	Providers    []TTSProviderConfig   `mapstructure:"providers"`
	Cache        TTSCacheConfig        `mapstructure:"cache"`
	Latency      TTSLatencyConfig      `mapstructure:"latency"`
	Fallback     TTSFallbackConfig     `mapstructure:"fallback"`
	VoiceCatalog TTSVoiceCatalogConfig `mapstructure:"voice_catalog"`
}

// TTSProviderConfig holds TTS provider configuration
type TTSProviderConfig struct {
	Name           string                  `mapstructure:"name"`
	Type           string                  `mapstructure:"type"`
	Priority       int                     `mapstructure:"priority"`
	Enabled        bool                    `mapstructure:"enabled"`
	Config         map[string]interface{}  `mapstructure:"config"`
	ConnectionPool TTSConnectionPoolConfig `mapstructure:"connection_pool"`
	Retry          TTSRetryConfig          `mapstructure:"retry"`
}

// TTSConnectionPoolConfig holds connection pool configuration
type TTSConnectionPoolConfig struct {
	MaxConnections int `mapstructure:"max_connections"`
	IdleTimeout    int `mapstructure:"idle_timeout"`
}

// TTSRetryConfig holds retry configuration
type TTSRetryConfig struct {
	MaxRetries int `mapstructure:"max_retries"`
	BaseDelay  int `mapstructure:"base_delay"`
	MaxDelay   int `mapstructure:"max_delay"`
}

// TTSCacheConfig holds audio caching configuration
type TTSCacheConfig struct {
	Enabled bool                  `mapstructure:"enabled"`
	Redis   TTSCacheRedisConfig   `mapstructure:"redis"`
	Warming TTSCacheWarmingConfig `mapstructure:"warming"`
}

// TTSCacheRedisConfig holds Redis cache configuration
type TTSCacheRedisConfig struct {
	KeyPrefix string `mapstructure:"key_prefix"`
	TTL       int    `mapstructure:"ttl"`
	MaxSize   int    `mapstructure:"max_size"`
}

// TTSCacheWarmingConfig holds cache warming configuration
type TTSCacheWarmingConfig struct {
	Enabled  bool     `mapstructure:"enabled"`
	Prompts  []string `mapstructure:"prompts"`
	Interval int      `mapstructure:"interval"`
}

// TTSLatencyConfig holds latency tracking configuration
type TTSLatencyConfig struct {
	Enabled         bool  `mapstructure:"enabled"`
	FirstByteTarget int   `mapstructure:"first_byte_target"`
	Percentiles     []int `mapstructure:"percentiles"`
}

// TTSFallbackConfig holds fallback configuration
type TTSFallbackConfig struct {
	Enabled              bool    `mapstructure:"enabled"`
	Timeout              int     `mapstructure:"timeout"`
	DegradationThreshold float64 `mapstructure:"degradation_threshold"`
}

// TTSVoiceCatalogConfig holds voice catalog sync configuration
type TTSVoiceCatalogConfig struct {
	Enabled      bool `mapstructure:"enabled"`
	SyncInterval int  `mapstructure:"sync_interval"`
}

// CampaignConfig holds campaign-related configuration
type CampaignConfig struct {
	Manager          CampaignManagerConfig          `mapstructure:"manager"`
	PredictiveDialer CampaignPredictiveDialerConfig `mapstructure:"predictive_dialer"`
}

// CampaignManagerConfig holds Campaign Manager configuration
type CampaignManagerConfig struct {
	GRPC          CampaignManagerGRPCConfig          `mapstructure:"grpc"`
	ContactLoader CampaignManagerContactLoaderConfig `mapstructure:"contact_loader"`
	Compliance    CampaignManagerComplianceConfig    `mapstructure:"compliance"`
	CallingWindow CampaignManagerCallingWindowConfig `mapstructure:"calling_window"`
	Scheduling    CampaignManagerSchedulingConfig    `mapstructure:"scheduling"`
}

// CampaignManagerGRPCConfig holds gRPC server configuration
type CampaignManagerGRPCConfig struct {
	Listen         string `mapstructure:"listen"`
	MaxRecvMsgSize int    `mapstructure:"max_recv_msg_size"`
	MaxSendMsgSize int    `mapstructure:"max_send_msg_size"`
}

// CampaignManagerContactLoaderConfig holds contact loader configuration
type CampaignManagerContactLoaderConfig struct {
	BatchSize   int `mapstructure:"batch_size"`
	Concurrency int `mapstructure:"concurrency"`
	Timeout     int `mapstructure:"timeout"`
}

// CampaignManagerComplianceConfig holds compliance configuration
type CampaignManagerComplianceConfig struct {
	DNCCheckEnabled     bool `mapstructure:"dnc_check_enabled"`
	ConsentCheckEnabled bool `mapstructure:"consent_check_enabled"`
	CacheTTL            int  `mapstructure:"cache_ttl"`
	CacheSize           int  `mapstructure:"cache_size"`
}

// CampaignManagerCallingWindowConfig holds calling window configuration
type CampaignManagerCallingWindowConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	DefaultStart string `mapstructure:"default_start"`
	DefaultEnd   string `mapstructure:"default_end"`
	TimezoneDB   string `mapstructure:"timezone_db"`
}

// CampaignManagerSchedulingConfig holds scheduling configuration
type CampaignManagerSchedulingConfig struct {
	Enabled       bool `mapstructure:"enabled"`
	RRULESupport  bool `mapstructure:"rrule_support"`
	CheckInterval int  `mapstructure:"check_interval"`
}

// CampaignPredictiveDialerConfig holds predictive dialer configuration
type CampaignPredictiveDialerConfig struct {
	Enabled bool                                  `mapstructure:"enabled"`
	ErlangC CampaignPredictiveDialerErlangCConfig `mapstructure:"erlang_c"`
	Metrics CampaignPredictiveDialerMetricsConfig `mapstructure:"metrics"`
	AMD     CampaignPredictiveDialerAMDConfig     `mapstructure:"amd"`
	Retry   CampaignPredictiveDialerRetryConfig   `mapstructure:"retry"`
}

// CampaignPredictiveDialerErlangCConfig holds Erlang-C configuration
type CampaignPredictiveDialerErlangCConfig struct {
	TargetUtilization float64 `mapstructure:"target_utilization"`
	SafetyMargin      float64 `mapstructure:"safety_margin"`
	MaxAbandonRate    float64 `mapstructure:"max_abandon_rate"`
}

// CampaignPredictiveDialerMetricsConfig holds metrics configuration
type CampaignPredictiveDialerMetricsConfig struct {
	Enabled          bool `mapstructure:"enabled"`
	Interval         int  `mapstructure:"interval"`
	EmitToPrometheus bool `mapstructure:"emit_to_prometheus"`
}

// CampaignPredictiveDialerAMDConfig holds AMD configuration
type CampaignPredictiveDialerAMDConfig struct {
	Enabled            bool    `mapstructure:"enabled"`
	HumanThreshold     float64 `mapstructure:"human_threshold"`
	VoicemailThreshold float64 `mapstructure:"voicemail_threshold"`
	SilenceDuration    float64 `mapstructure:"silence_duration"`
	BeepDetection      bool    `mapstructure:"beep_detection"`
	SpeechDetection    bool    `mapstructure:"speech_detection"`
}

// CampaignPredictiveDialerRetryConfig holds retry configuration
type CampaignPredictiveDialerRetryConfig struct {
	MaxAttempts        int  `mapstructure:"max_attempts"`
	BaseDelay          int  `mapstructure:"base_delay"`
	MaxDelay           int  `mapstructure:"max_delay"`
	ExponentialBackoff bool `mapstructure:"exponential_backoff"`
}

// TransferConfig holds transfer-related configuration
type TransferConfig struct {
	Manager         TransferManagerConfig         `mapstructure:"manager"`
	StateMachine    TransferStateMachineConfig    `mapstructure:"state_machine"`
	PriorityQueuing TransferPriorityQueuingConfig `mapstructure:"priority_queuing"`
	Analytics       TransferAnalyticsConfig       `mapstructure:"analytics"`
}

// TransferManagerConfig holds Transfer Manager configuration
type TransferManagerConfig struct {
	GRPC           TransferManagerGRPCConfig           `mapstructure:"grpc"`
	Queue          TransferManagerQueueConfig          `mapstructure:"queue"`
	AgentPolling   TransferManagerAgentPollingConfig   `mapstructure:"agent_polling"`
	SIPRefer       TransferManagerSIPReferConfig       `mapstructure:"sip_refer"`
	ContextHeaders TransferManagerContextHeadersConfig `mapstructure:"context_headers"`
}

// TransferManagerGRPCConfig holds gRPC server configuration
type TransferManagerGRPCConfig struct {
	Listen         string `mapstructure:"listen"`
	MaxRecvMsgSize int    `mapstructure:"max_recv_msg_size"`
	MaxSendMsgSize int    `mapstructure:"max_send_msg_size"`
}

// TransferManagerQueueConfig holds queue integration configuration
type TransferManagerQueueConfig struct {
	Type         string `mapstructure:"type"`
	Endpoint     string `mapstructure:"endpoint"`
	APIKey       string `mapstructure:"api_key"`
	PollInterval int    `mapstructure:"poll_interval"`
	Timeout      int    `mapstructure:"timeout"`
}

// TransferManagerAgentPollingConfig holds agent polling configuration
type TransferManagerAgentPollingConfig struct {
	Enabled  bool `mapstructure:"enabled"`
	Interval int  `mapstructure:"interval"`
	Timeout  int  `mapstructure:"timeout"`
}

// TransferManagerSIPReferConfig holds SIP REFER configuration
type TransferManagerSIPReferConfig struct {
	Enabled          bool `mapstructure:"enabled"`
	ReplacesHeader   bool `mapstructure:"replaces_header"`
	AttendedTransfer bool `mapstructure:"attended_transfer"`
	Timeout          int  `mapstructure:"timeout"`
}

// TransferManagerContextHeadersConfig holds context headers configuration
type TransferManagerContextHeadersConfig struct {
	Enabled              bool   `mapstructure:"enabled"`
	HeaderPrefix         string `mapstructure:"header_prefix"`
	IncludeTranscript    bool   `mapstructure:"include_transcript"`
	IncludeExtractedData bool   `mapstructure:"include_extracted_data"`
	IncludeSentiment     bool   `mapstructure:"include_sentiment"`
}

// TransferStateMachineConfig holds transfer state machine configuration
type TransferStateMachineConfig struct {
	QueuedTimeout    int                                `mapstructure:"queued_timeout"`
	RingingTimeout   int                                `mapstructure:"ringing_timeout"`
	ConnectedTimeout int                                `mapstructure:"connected_timeout"`
	Fallback         TransferStateMachineFallbackConfig `mapstructure:"fallback"`
}

// TransferStateMachineFallbackConfig holds fallback configuration
type TransferStateMachineFallbackConfig struct {
	Enabled          bool `mapstructure:"enabled"`
	VoicemailEnabled bool `mapstructure:"voicemail_enabled"`
	CallbackEnabled  bool `mapstructure:"callback_enabled"`
	CallbackDelay    int  `mapstructure:"callback_delay"`
}

// TransferPriorityQueuingConfig holds priority queuing configuration
type TransferPriorityQueuingConfig struct {
	Enabled         bool `mapstructure:"enabled"`
	Levels          int  `mapstructure:"levels"`
	DefaultPriority int  `mapstructure:"default_priority"`
}

// TransferAnalyticsConfig holds analytics configuration
type TransferAnalyticsConfig struct {
	Enabled  bool     `mapstructure:"enabled"`
	Metrics  []string `mapstructure:"metrics"`
	Interval int      `mapstructure:"interval"`
}

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	Postgres StoragePostgresConfig `mapstructure:"postgres"`
	S3       StorageS3Config       `mapstructure:"s3"`
	Redis    StorageRedisConfig    `mapstructure:"redis"`
}

// StoragePostgresConfig holds PostgreSQL configuration
type StoragePostgresConfig struct {
	Host               string                            `mapstructure:"host"`
	Port               int                               `mapstructure:"port"`
	Database           string                            `mapstructure:"database"`
	Username           string                            `mapstructure:"username"`
	Password           string                            `mapstructure:"password"`
	SSLMode            string                            `mapstructure:"ssl_mode"`
	MaxConnections     int                               `mapstructure:"max_connections"`
	MaxIdleConnections int                               `mapstructure:"max_idle_connections"`
	ConnectionTimeout  int                               `mapstructure:"connection_timeout"`
	StatementTimeout   int                               `mapstructure:"statement_timeout"`
	Partitioning       StoragePostgresPartitioningConfig `mapstructure:"partitioning"`
	Tables             StoragePostgresTablesConfig       `mapstructure:"tables"`
}

// StoragePostgresPartitioningConfig holds partitioning configuration
type StoragePostgresPartitioningConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	Strategy     string `mapstructure:"strategy"`
	PartitionKey string `mapstructure:"partition_key"`
}

// StoragePostgresTablesConfig holds table names
type StoragePostgresTablesConfig struct {
	CallRecords    string `mapstructure:"call_records"`
	DNCEntries     string `mapstructure:"dnc_entries"`
	ConsentRecords string `mapstructure:"consent_records"`
	Campaigns      string `mapstructure:"campaigns"`
	Contacts       string `mapstructure:"contacts"`
}

// StorageS3Config holds S3/MinIO configuration
type StorageS3Config struct {
	Endpoint       string                     `mapstructure:"endpoint"`
	AccessKey      string                     `mapstructure:"access_key"`
	SecretKey      string                     `mapstructure:"secret_key"`
	Bucket         string                     `mapstructure:"bucket"`
	Region         string                     `mapstructure:"region"`
	UseSSL         bool                       `mapstructure:"use_ssl"`
	ForcePathStyle bool                       `mapstructure:"force_path_style"`
	Recordings     StorageS3RecordingsConfig  `mapstructure:"recordings"`
	Transcripts    StorageS3TranscriptsConfig `mapstructure:"transcripts"`
}

// StorageS3RecordingsConfig holds recordings configuration
type StorageS3RecordingsConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Format      string `mapstructure:"format"`
	SampleRate  int    `mapstructure:"sample_rate"`
	Channels    int    `mapstructure:"channels"`
	Compression string `mapstructure:"compression"`
}

// StorageS3TranscriptsConfig holds transcripts configuration
type StorageS3TranscriptsConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Format      string `mapstructure:"format"`
	Compression string `mapstructure:"compression"`
}

// StorageRedisConfig holds Redis configuration
type StorageRedisConfig struct {
	Host               string                      `mapstructure:"host"`
	Port               int                         `mapstructure:"port"`
	Password           string                      `mapstructure:"password"`
	Database           int                         `mapstructure:"database"`
	MaxConnections     int                         `mapstructure:"max_connections"`
	MinIdleConnections int                         `mapstructure:"min_idle_connections"`
	ConnectionTimeout  int                         `mapstructure:"connection_timeout"`
	ReadTimeout        int                         `mapstructure:"read_timeout"`
	WriteTimeout       int                         `mapstructure:"write_timeout"`
	Session            StorageRedisSessionConfig   `mapstructure:"session"`
	RateLimit          StorageRedisRateLimitConfig `mapstructure:"rate_limit"`
}

// StorageRedisSessionConfig holds session configuration
type StorageRedisSessionConfig struct {
	TTL    int    `mapstructure:"ttl"`
	Prefix string `mapstructure:"prefix"`
}

// StorageRedisRateLimitConfig holds rate limiting configuration
type StorageRedisRateLimitConfig struct {
	Enabled       bool `mapstructure:"enabled"`
	DefaultLimit  int  `mapstructure:"default_limit"`
	DefaultWindow int  `mapstructure:"default_window"`
}

// ObservabilityConfig holds observability configuration
type ObservabilityConfig struct {
	Logging ObservabilityLoggingConfig `mapstructure:"logging"`
	Metrics ObservabilityMetricsConfig `mapstructure:"metrics"`
	Tracing ObservabilityTracingConfig `mapstructure:"tracing"`
	Grafana ObservabilityGrafanaConfig `mapstructure:"grafana"`
}

// ObservabilityLoggingConfig holds logging configuration
type ObservabilityLoggingConfig struct {
	Level              string                         `mapstructure:"level"`
	Format             string                         `mapstructure:"format"`
	Output             string                         `mapstructure:"output"`
	TraceIDPropagation bool                           `mapstructure:"trace_id_propagation"`
	ServiceName        string                         `mapstructure:"service_name"`
	Loki               ObservabilityLoggingLokiConfig `mapstructure:"loki"`
}

// ObservabilityLoggingLokiConfig holds Loki configuration
type ObservabilityLoggingLokiConfig struct {
	Enabled      bool              `mapstructure:"enabled"`
	Endpoint     string            `mapstructure:"endpoint"`
	BatchSize    int               `mapstructure:"batch_size"`
	BatchTimeout int               `mapstructure:"batch_timeout"`
	Retention    string            `mapstructure:"retention"`
	Labels       map[string]string `mapstructure:"labels"`
}

// ObservabilityMetricsConfig holds metrics configuration
type ObservabilityMetricsConfig struct {
	Enabled         bool                                      `mapstructure:"enabled"`
	Prometheus      ObservabilityMetricsPrometheusConfig      `mapstructure:"prometheus"`
	CallLatency     ObservabilityMetricsCallLatencyConfig     `mapstructure:"call_latency"`
	ConcurrentCalls ObservabilityMetricsConcurrentCallsConfig `mapstructure:"concurrent_calls"`
	ErrorRates      ObservabilityMetricsErrorRatesConfig      `mapstructure:"error_rates"`
	CircuitBreaker  ObservabilityMetricsCircuitBreakerConfig  `mapstructure:"circuit_breaker"`
	CampaignMetrics ObservabilityMetricsCampaignMetricsConfig `mapstructure:"campaign_metrics"`
	MediaQuality    ObservabilityMetricsMediaQualityConfig    `mapstructure:"media_quality"`
}

// ObservabilityMetricsPrometheusConfig holds Prometheus configuration
type ObservabilityMetricsPrometheusConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Listen    string `mapstructure:"listen"`
	Path      string `mapstructure:"path"`
	Namespace string `mapstructure:"namespace"`
	Subsystem string `mapstructure:"subsystem"`
}

// ObservabilityMetricsCallLatencyConfig holds call latency configuration
type ObservabilityMetricsCallLatencyConfig struct {
	Enabled bool      `mapstructure:"enabled"`
	Buckets []float64 `mapstructure:"buckets"`
	Labels  []string  `mapstructure:"labels"`
}

// ObservabilityMetricsConcurrentCallsConfig holds concurrent calls configuration
type ObservabilityMetricsConcurrentCallsConfig struct {
	Enabled bool     `mapstructure:"enabled"`
	Labels  []string `mapstructure:"labels"`
}

// ObservabilityMetricsErrorRatesConfig holds error rates configuration
type ObservabilityMetricsErrorRatesConfig struct {
	Enabled bool     `mapstructure:"enabled"`
	Labels  []string `mapstructure:"labels"`
}

// ObservabilityMetricsCircuitBreakerConfig holds circuit breaker configuration
type ObservabilityMetricsCircuitBreakerConfig struct {
	Enabled bool     `mapstructure:"enabled"`
	Labels  []string `mapstructure:"labels"`
}

// ObservabilityMetricsCampaignMetricsConfig holds campaign metrics configuration
type ObservabilityMetricsCampaignMetricsConfig struct {
	Enabled bool     `mapstructure:"enabled"`
	Labels  []string `mapstructure:"labels"`
}

// ObservabilityMetricsMediaQualityConfig holds media quality configuration
type ObservabilityMetricsMediaQualityConfig struct {
	Enabled       bool      `mapstructure:"enabled"`
	Labels        []string  `mapstructure:"labels"`
	JitterBuckets []float64 `mapstructure:"jitter_buckets"`
	LossBuckets   []float64 `mapstructure:"loss_buckets"`
	MOSBuckets    []float64 `mapstructure:"mos_buckets"`
}

// ObservabilityTracingConfig holds tracing configuration
type ObservabilityTracingConfig struct {
	Enabled bool                           `mapstructure:"enabled"`
	OTEL    ObservabilityTracingOTELConfig `mapstructure:"otel"`
}

// ObservabilityTracingOTELConfig holds OTEL configuration
type ObservabilityTracingOTELConfig struct {
	ServiceName        string                                       `mapstructure:"service_name"`
	Exporter           ObservabilityTracingOTELExporterConfig       `mapstructure:"exporter"`
	Sampler            ObservabilityTracingOTELSamplerConfig        `mapstructure:"sampler"`
	Propagators        []string                                     `mapstructure:"propagators"`
	SpanAttributes     ObservabilityTracingOTELSpanAttributesConfig `mapstructure:"span_attributes"`
	ResourceAttributes map[string]string                            `mapstructure:"resource_attributes"`
}

// ObservabilityTracingOTELExporterConfig holds OTEL exporter configuration
type ObservabilityTracingOTELExporterConfig struct {
	Type     string            `mapstructure:"type"`
	Endpoint string            `mapstructure:"endpoint"`
	Insecure bool              `mapstructure:"insecure"`
	Headers  map[string]string `mapstructure:"headers"`
}

// ObservabilityTracingOTELSamplerConfig holds OTEL sampler configuration
type ObservabilityTracingOTELSamplerConfig struct {
	Type  string  `mapstructure:"type"`
	Ratio float64 `mapstructure:"ratio"`
}

// ObservabilityTracingOTELSpanAttributesConfig holds span attributes configuration
type ObservabilityTracingOTELSpanAttributesConfig struct {
	CallID   bool `mapstructure:"call_id"`
	TenantID bool `mapstructure:"tenant_id"`
	Provider bool `mapstructure:"provider"`
	Model    bool `mapstructure:"model"`
}

// ObservabilityGrafanaConfig holds Grafana configuration
type ObservabilityGrafanaConfig struct {
	Enabled    bool                                  `mapstructure:"enabled"`
	Dashboards []ObservabilityGrafanaDashboardConfig `mapstructure:"dashboards"`
	Alerting   ObservabilityGrafanaAlertingConfig    `mapstructure:"alerting"`
}

// ObservabilityGrafanaDashboardConfig holds dashboard configuration
type ObservabilityGrafanaDashboardConfig struct {
	Name    string `mapstructure:"name"`
	Title   string `mapstructure:"title"`
	Refresh string `mapstructure:"refresh"`
}

// ObservabilityGrafanaAlertingConfig holds alerting configuration
type ObservabilityGrafanaAlertingConfig struct {
	Enabled   bool                                        `mapstructure:"enabled"`
	PagerDuty ObservabilityGrafanaAlertingPagerDutyConfig `mapstructure:"pagerduty"`
	Rules     []ObservabilityGrafanaAlertingRuleConfig    `mapstructure:"rules"`
}

// ObservabilityGrafanaAlertingPagerDutyConfig holds PagerDuty configuration
type ObservabilityGrafanaAlertingPagerDutyConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	IntegrationKey string `mapstructure:"integration_key"`
	Severity       string `mapstructure:"severity"`
}

// ObservabilityGrafanaAlertingRuleConfig holds alerting rule configuration
type ObservabilityGrafanaAlertingRuleConfig struct {
	Name        string            `mapstructure:"name"`
	Expr        string            `mapstructure:"expr"`
	For         string            `mapstructure:"for"`
	Labels      map[string]string `mapstructure:"labels"`
	Annotations map[string]string `mapstructure:"annotations"`
}

var (
	cfg  *Config
	once sync.Once
)

// ResetForTest resets the configuration state for testing purposes
func ResetForTest() {
	cfg = nil
	once = sync.Once{}
}

// Load loads the configuration from file and environment variables
func Load(configPath string, logger *zap.Logger) (*Config, error) {
	var err error
	once.Do(func() {
		cfg, err = loadConfig(configPath, logger)
	})
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// MustLoad loads the configuration and panics on error
func MustLoad(configPath string, logger *zap.Logger) *Config {
	cfg, err := Load(configPath, logger)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}
	return cfg
}

// Get returns the loaded configuration (must call Load first)
func Get() *Config {
	if cfg == nil {
		panic("configuration not loaded, call Load() first")
	}
	return cfg
}

func loadConfig(configPath string, logger *zap.Logger) (*Config, error) {
	v := viper.New()

	// Set config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Default config paths
		v.SetConfigName("platform")
		v.SetConfigType("yaml")
		v.AddConfigPath("/etc/ai-telephony-platform/")
		v.AddConfigPath("./config/")
		v.AddConfigPath(".")
	}

	// Environment variable overrides
	v.SetEnvPrefix("AI_TELEPHONY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind environment variables for sensitive fields
	bindEnvVars(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warn("Config file not found, using defaults and environment variables", zap.Error(err))
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		logger.Info("Loaded configuration file", zap.String("file", v.ConfigFileUsed()))
	}

	// Unmarshal into config struct
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&c); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &c, nil
}

func bindEnvVars(v *viper.Viper) {
	// Note: Viper has a known issue where binding env vars to slice indices
	// (e.g., stt.providers.0) breaks unmarshalling of the entire slice from YAML.
	// Users should use AI_TELEPHONY_... env vars or we should override post-unmarshal.
}

func validateConfig(c *Config) error {
	// In a microservices environment where each service receives a subset of the config
	// via its own Helm chart, strict global validation will fail. We relax this for now.
	// Individual services should validate their own required sections.
	return nil
}

// GetConfigPath returns the default config path
func GetConfigPath() string {
	// Check environment variable first
	if path := os.Getenv("AI_TELEPHONY_CONFIG"); path != "" {
		return path
	}

	// Check common locations
	paths := []string{
		"/etc/ai-telephony-platform/platform.yaml",
		"./config/platform.yaml",
		"../config/platform.yaml",
		"../../config/platform.yaml",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

// Reload reloads the configuration
func Reload(configPath string, logger *zap.Logger) (*Config, error) {
	once = sync.Once{}
	return Load(configPath, logger)
}

// WatchConfig watches for config file changes and reloads
func WatchConfig(configPath string, logger *zap.Logger, onChange func(*Config)) error {
	v := viper.New()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("platform")
		v.SetConfigType("yaml")
		v.AddConfigPath("/etc/ai-telephony-platform/")
		v.AddConfigPath("./config/")
		v.AddConfigPath(".")
	}

	v.SetEnvPrefix("AI_TELEPHONY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnvVars(v)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logger.Info("Configuration file changed, reloading...", zap.String("file", e.Name))
		var c Config
		if err := v.Unmarshal(&c); err != nil {
			logger.Error("Failed to unmarshal config after change", zap.Error(err))
			return
		}
		if err := validateConfig(&c); err != nil {
			logger.Error("Configuration validation failed after change", zap.Error(err))
			return
		}
		cfg = &c
		if onChange != nil {
			onChange(&c)
		}
		logger.Info("Configuration reloaded successfully")
	})

	return nil
}
