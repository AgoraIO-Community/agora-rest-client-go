package req

// @brief Defines advanced feature configurations for the agent to join the RTC channel
//
// @since v0.7.0
type JoinPropertiesAdvancedFeaturesBody struct {
	// Whether to enable graceful interruption (AIVAD) (optional)
	//
	// When enabled, users can interrupt the AI at any time and respond quickly, achieving natural transitions and smooth conversations
	//
	//  - true: Enable
	//  - false: Disable (default)
	//
	EnableAIVad *bool `json:"enable_aivad,omitempty"`

	// Whether to enable the Real-time Messaging (RTM) module (optional)
	//
	// When enabled, the agent can use the capabilities provided by RTM to implement some advanced features
	//  - true: Enable
	//  - false: Disable (default)
	//
	EnableRtm *bool `json:"enable_rtm,omitempty"`
}

type TTSVendorParamsInterface interface {
	VendorParam()
}

type TTSMinimaxVendorVoiceSettingParam struct {
	VoiceId string  `json:"voice_id"`
	Speed   float32 `json:"speed"`
	Vol     float32 `json:"vol"`
	Pitch   int     `json:"pitch"`
	Emotion string  `json:"emotion"`
}

type TTSMinimaxVendorAudioSettingParam struct {
	SampleRate int `json:"sample_rate"`
}

// @brief Defines the Minimax vendor parameters for the Text-to-Speech (TTS) module when the agent joins the RTC channel, see
// https://platform.minimaxi.com/document/T2A%20V2 for details
//
// @since v0.7.0
type TTSMinimaxVendorParams struct {
	GroupId      string                            `json:"group_id"`
	Key          string                            `json:"key"`
	Model        string                            `json:"model"`
	VoiceSetting TTSMinimaxVendorVoiceSettingParam `json:"voice_setting"`
	AudioSetting TTSMinimaxVendorAudioSettingParam `json:"audio_setting"`
}

func (TTSMinimaxVendorParams) VendorParam() {}

// @brief Defines the Tencent vendor parameters for the Text-to-Speech (TTS) module when the agent joins the RTC channel, see
// https://cloud.tencent.com/document/product/1073/94308 for details
//
// @since v0.7.0
type TTSTencentVendorParams struct {
	AppId            string `json:"app_id"`
	SecretId         string `json:"secret_id"`
	SecretKey        string `json:"secret_key"`
	VoiceType        int    `json:"voice_type"`
	Volume           int    `json:"volume"`
	Speed            int    `json:"speed"`
	EmotionCategory  string `json:"emotion_category"`
	EmotionIntensity int    `json:"emotion_intensity"`
}

func (TTSTencentVendorParams) VendorParam() {}

// @brief Defines the Bytedance vendor parameters for the Text-to-Speech (TTS) module when the agent joins the RTC channel, see
// https://www.volcengine.com/docs/6561/79823 for details
//
// @since v0.7.0
type TTSBytedanceVendorParams struct {
	Token       string  `json:"token"`
	AppId       string  `json:"app_id"`
	Cluster     string  `json:"cluster"`
	VoiceType   string  `json:"voice_type"`
	SpeedRatio  float32 `json:"speed_ratio"`
	VolumeRatio float32 `json:"volume_ratio"`
	PitchRatio  float32 `json:"pitch_ratio"`
	Emotion     string  `json:"emotion"`
}

func (TTSBytedanceVendorParams) VendorParam() {}

type TTSMicrosoftVendorParams struct {
	Key       string  `json:"key"`
	Region    string  `json:"region"`
	VoiceName string  `json:"voice_name"`
	Rate      float32 `json:"rate"`
	Volume    float32 `json:"volume"`
}

func (TTSMicrosoftVendorParams) VendorParam() {}

type TTSElevenLabsVendorParams struct {
	APIKey  string `json:"api_key"`
	ModelId string `json:"model_id"`
	VoiceId string `json:"voice_id"`
}

func (TTSElevenLabsVendorParams) VendorParam() {}

// @brief Defines the Text-to-Speech (TTS) module configuration for the agent to join the RTC channel
//
// @since v0.7.0
type JoinPropertiesTTSBody struct {
	// TTS vendor, see TTSVendor for details
	Vendor TTSVendor `json:"vendor"`
	// TTS vendor parameter description, see
	// 	- TTSMinimaxVendorParams
	//
	//  - TTSTencentVendorParams
	//
	// 	- TTSBytedanceVendorParams
	//
	//  - TTSMicrosoftVendorParams
	//
	//  - TTSElevenLabsVendorParams
	Params TTSVendorParamsInterface `json:"params"`
}

// @brief Defines the TTS vendor enumeration
//
// @since v0.7.0
type TTSVendor string

const (
	// Minimax TTS vendor
	MinimaxTTSVendor TTSVendor = "minimax"
	// Tencent TTS vendor
	TencentTTSVendor TTSVendor = "tencent"
	// Bytedance TTS vendor
	BytedanceTTSVendor TTSVendor = "bytedance"
	// Microsoft TTS vendor
	MicrosoftTTSVendor TTSVendor = "microsoft"
	// ElevenLabs TTS vendor
	ElevenLabsTTSVendor TTSVendor = "elevenLabs"
)

// @brief Defines the custom language model (LLM) configuration for the agent to join the RTC channel
//
// @since v0.7.0
type JoinPropertiesCustomLLMBody struct {
	// LLM callback URL (required)
	//
	// Must be compatible with the OpenAI protocol
	Url string `json:"url"`

	// LLM API key for verification (required)
	//
	// Default is empty, make sure to enable the API key in the production environment
	APIKey string `json:"api_key"`

	// A set of predefined information attached at the beginning of each LLM call to control LLM output (optional)
	//
	// Can be role settings, prompts, and answer samples, must be compatible with the OpenAI protocol
	SystemMessages []map[string]any `json:"system_messages"`

	// Additional information transmitted in the LLM message body, such as the model used, maximum token Limit, etc. (optional)
	//
	// Different LLM vendors support different configurations, see their respective LLM documentation for details
	Params map[string]any `json:"params"`

	// Number of short-term memory entries cached in the LLM (optional)
	//
	// Default value is 10
	//
	// Passing 0 means no short-term memory is cached. The agent and subscribed users will record entries separately
	MaxHistory *int `json:"max_history,omitempty"`

	// Agent greeting message (optional)
	//
	// If filled, the agent will automatically send a greeting message to the first subscribed user who joins the channel when there are no users in the remote_rtc_uids list.
	GreetingMessage string `json:"greeting_message,omitempty"`

	// Input modalities for the LLM (optional)
	//
	// 	- ["text"]: Text only (default)
	//
	//  - ["text", "image"]: Text and image, requires the selected LLM to support visual modality input
	//
	InputModalities []string `json:"input_modalities,omitempty"`

	// Output modalities for the LLM (optional)
	//
	// 	- ["text"]: Text only (default), the output text will be converted to speech by the TTS module and published to the RTC channel.
	//
	// 	- ["audio"]: Audio only. The audio will be directly published to the RTC channel.
	//
	//  - ["text", "audio"]: Text and audio. You can write your own logic to handle the LLM output as needed
	//
	OutputModalities []string `json:"output_modalities,omitempty"`

	// Failure message for the agent (optional)
	//
	// If filled, it will be returned through the TTS module when the LLM call fails.
	FailureMessage string `json:"failure_message,omitempty"`

	// Silence message for the agent (optional)
	//
	// After the agent is created and a user joins the channel,
	// the duration of the agent's non-listening, thinking, or speaking state is called the agent's silence time.
	// When the silence time reaches the set value, the agent will report the silence prompt message in llm.silence_message,
	// and then recalculate the silence time.
	//
	// When silence_timeout is set to 0, this parameter is ignored.
	SilenceMessage *string `json:"silence_message,omitempty"`
}

// @brief Defines the Voice Activity Detection (VAD) configuration for the agent to join the RTC channel
//
// @since v0.7.0
type JoinPropertiesVadBody struct {
	// Human voice duration threshold (ms), range is [120, 1200] (optional)
	//
	// Minimum duration of detected human voice signal to avoid false interruptions
	InterruptDurationMs *int `json:"interrupt_duration_ms,omitempty"`
	// Prefix padding threshold (ms), range is [0, 5000] (optional)
	//
	// Minimum duration of continuous voice required to start a new voice segment, avoiding very short sounds triggering voice activity detection
	PrefixPaddingMs *int `json:"prefix_padding_ms,omitempty"`
	// Silence duration threshold (ms), range is [0, 2000] (optional)
	//
	// Minimum duration of silence at the end of speech to ensure short pauses do not prematurely end the voice segment
	SilenceDurationMs *int `json:"silence_duration_ms,omitempty"`
	// Voice recognition sensitivity, range is (0.0,1.0) (optional)
	//
	// Determines the extent to which audio signals are considered "voice activity".
	//
	// Lower values make it easier for the agent to detect voice, higher values may ignore faint sounds.
	Threshold *float64 `json:"threshold,omitempty"`
}

// @brief Defines the Automatic Speech Recognition (ASR) configuration for the agent to join the RTC channel
//
// @since v0.7.0
type JoinPropertiesAsrBody struct {
	// Language used for interaction between the user and the agent (optional)
	//
	//  - zh-CN: Chinese (supports mixed Chinese and English) (default)
	//
	//  - en-US: English
	Language string `json:"language,omitempty"`
}

// @brief Request body for calling the Conversational AI engine Join API
//
// @since v0.7.0
type JoinPropertiesReqBody struct {
	// Token used to join the RTC channel, i.e., the dynamic key for authentication (optional). If your project has enabled App Certificate, be sure to pass the dynamic key of your project in this field
	Token string `json:"token"`
	// RTC channel name the agent joins (required)
	Channel string `json:"channel"`
	// User ID of the agent in the RTC channel (required)
	//
	// Filling "0" means a random assignment, but the Token needs to be modified accordingly
	AgentRtcUId string `json:"agent_rtc_uid"`
	// List of user IDs the agent subscribes to in the RTC channel, only subscribed users can interact with the agent (required)
	//
	// Passing "*" means subscribing to all users in the channel
	RemoteRtcUIds []string `json:"remote_rtc_uids"`
	// Whether to enable String UID (optional)
	//
	//  - true: Enable String UID
	//
	//  - false: Disable String UID (default)
	EnableStringUId *bool `json:"enable_string_uid,omitempty"`
	// Maximum idle time of the RTC channel (s) (optional)
	//
	// The time after detecting that all users specified in remote_rtc_uids have left the channel is considered idle time.
	//
	// If it exceeds the set maximum value, the agent in the channel will automatically stop and exit the channel.
	//
	// If set to 0, the agent will not stop until manually exited
	IdleTimeout *int `json:"idle_timeout,omitempty"`

	// Silence timeout (s) (optional)
	//
	// The maximum silence time of the agent (s), range is [0,60].
	//
	// After the agent is created and a user joins the channel, the duration of the agent's non-listening, thinking, or speaking state is called the agent's silence time.
	//
	// When the silence time reaches the set value, the agent will report the silence prompt message in llm.silence_message, and then recalculate the silence time.
	//
	//  - 0 (default): Do not enable this feature.
	//
	//  - (0,60]: Must also set llm.silence_message to enable the feature.
	SilenceTimeout *int `json:"silence_timeout,omitempty"`

	// Agent user ID in the RTM channel
	//
	// Only valid when advanced_features.enable_rtm is true
	AgentRtmUId *string `json:"agent_rtm_uid,omitempty"`
	// Advanced feature configurations (optional), see JoinPropertiesAdvancedFeaturesBody for details
	AdvancedFeatures *JoinPropertiesAdvancedFeaturesBody `json:"advanced_features,omitempty"`
	// Custom language model (LLM) configuration (required), see JoinPropertiesCustomLLMBody for details
	LLM *JoinPropertiesCustomLLMBody `json:"llm,omitempty"`
	// Text-to-Speech (TTS) module configuration (required), see JoinPropertiesTTSBody for details
	TTS *JoinPropertiesTTSBody `json:"tts,omitempty"`
	// Voice Activity Detection (VAD) configuration (optional), see JoinPropertiesVadBody for details
	Vad *JoinPropertiesVadBody `json:"vad,omitempty"`
	// Automatic Speech Recognition (ASR) configuration (optional), see JoinPropertiesAsrBody for details
	Asr *JoinPropertiesAsrBody `json:"asr,omitempty"`
}
