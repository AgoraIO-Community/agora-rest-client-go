package req

import (
	"encoding/json"
)

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

	// Enable Multimodal Large Language Model.
	//
	//  - true: Enable
	//  - false: Disable (default)
	//
	// Enabling MLLM automatically disables ASR, LLM, and TTS. When you set this parameter to true, enable_aivad is also disabled.
	//
	// @since v0.12.0
	EnableMLLM *bool `json:"enable_mllm,omitempty"`
}

// @brief Defines the interface for TTS vendor parameters
//
// @since v0.12.0
type TTSVendorParamsInterface interface {
	// VendorParam marks this type as a valid TTS vendor parameter
	//
	// @since v0.12.0
	VendorParam()

	// GetVendorType returns the vendor type identifier for validation
	//
	// @since v0.12.0
	GetVendorType() TTSVendor
}

type TTSMinimaxVendorVoiceSettingParam struct {
	VoiceId              string  `json:"voice_id"`
	Speed                float32 `json:"speed"`
	Vol                  float32 `json:"vol"`
	Pitch                int     `json:"pitch"`
	Emotion              string  `json:"emotion"`
	LatexRender          bool    `json:"latex_render"`
	EnglishNormalization bool    `json:"english_normalization"`
}

type TTSMinimaxVendorAudioSettingParam struct {
	SampleRate int `json:"sample_rate"`
}

type PronunciationDictParam struct {
	Tone []string `json:"tone"`
}

type TimberWeightsParam struct {
	VoiceId string `json:"voice_id"`
	Weight  int    `json:"weight"`
}

// @brief Defines the Minimax vendor parameters for the Text-to-Speech (TTS) module when the agent joins the RTC channel, see
// https://platform.minimaxi.com/document/T2A%20V2 for details
//
// @since v0.7.0
type TTSMinimaxVendorParams struct {
	GroupId           string                             `json:"group_id"`
	Key               string                             `json:"key"`
	Model             string                             `json:"model"`
	VoiceSetting      *TTSMinimaxVendorVoiceSettingParam `json:"voice_setting,omitempty"`
	AudioSetting      *TTSMinimaxVendorAudioSettingParam `json:"audio_setting,omitempty"`
	PronunciationDict *PronunciationDictParam            `json:"pronunciation_dict,omitempty"`
	TimberWeights     []TimberWeightsParam               `json:"timber_weights,omitempty"`
}

func (TTSMinimaxVendorParams) VendorParam()             {}
func (TTSMinimaxVendorParams) GetVendorType() TTSVendor { return MinimaxTTSVendor }

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

func (TTSTencentVendorParams) VendorParam()             {}
func (TTSTencentVendorParams) GetVendorType() TTSVendor { return TencentTTSVendor }

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

func (TTSBytedanceVendorParams) VendorParam()             {}
func (TTSBytedanceVendorParams) GetVendorType() TTSVendor { return BytedanceTTSVendor }

type TTSMicrosoftVendorParams struct {
	// The API key used for authentication.(Required)
	Key string `json:"key"`
	// The Azure region where the speech service is hosted.(Required)
	Region string `json:"region"`
	// The identifier for the selected voice for speech synthesis.(Optional)
	VoiceName string `json:"voice_name"`
	// Indicates the speaking rate of the text.(Optional)
	//
	// The rate can be applied at the word or sentence level and should be between 0.5 and 2.0 times the original audio speed.
	Speed float32 `json:"speed"`
	// Specifies the audio volume as a number between 0.0 and 100.0, where 0.0 is the quietest and 100.0 is the loudest.
	//
	// For example, a value of 75 sets the volume to 75% of the maximum.
	//
	// The default value is 100.
	Volume float32 `json:"volume"`
	// Specifies the audio sampling rate in Hz.(Optional)
	//
	// The default value is 24000.
	SampleRate int `json:"sample_rate"`
}

func (TTSMicrosoftVendorParams) VendorParam()             {}
func (TTSMicrosoftVendorParams) GetVendorType() TTSVendor { return MicrosoftTTSVendor }

type TTSElevenLabsVendorParams struct {
	// The API key used for authentication.(Required)
	Key string `json:"key"`
	// Identifier of the model to be used.(Required)
	ModelId string `json:"model_id"`
	// The identifier for the selected voice for speech synthesis.(Required)
	VoiceId string `json:"voice_id"`
	// Specifies the audio sampling rate in Hz.(Optional)
	//
	// The default value is 24000.
	SampleRate int `json:"sample_rate"`
	// The stability for voice settings.(Optional)
	Stability float32 `json:"stability"`
	// Determines how closely the AI should adhere to the original voice when attempting to replicate it.
	SimilarityBoost float32 `json:"similarity_boost"`
	// Determines the style exaggeration of the voice. This setting attempts to amplify the style of the original speaker.
	//
	// It does consume additional computational resources and might increase latency if set to anything other than 0.
	Style float32 `json:"style"`
	// This setting boosts the similarity to the original speaker.
	//
	// Using this setting requires a slightly higher computational load, which in turn increases latency.
	UseSpeakerBoost bool `json:"use_speaker_boost"`
}

func (TTSElevenLabsVendorParams) VendorParam()             {}
func (TTSElevenLabsVendorParams) GetVendorType() TTSVendor { return ElevenLabsTTSVendor }

// @brief Defines the Cartesia vendor parameters for the Text-to-Speech (TTS) module when the agent joins the RTC channel
//
// @since v0.12.0
type TTSCartesiaVendorParams struct {
	APIKey  string                  `json:"api_key"`
	ModelId string                  `json:"model_id"`
	Voice   *TTSCartesiaVendorVoice `json:"voice"`
}

// @brief Defines the Cartesia vendor voice for the Text-to-Speech (TTS) module when the agent joins the RTC channel
//
// @since v0.12.0
type TTSCartesiaVendorVoice struct {
	Mode string `json:"mode"`
	Id   string `json:"id"`
}

func (TTSCartesiaVendorParams) VendorParam()             {}
func (TTSCartesiaVendorParams) GetVendorType() TTSVendor { return CartesiaTTSVendor }

// @brief Defines the OpenAI vendor parameters for the Text-to-Speech (TTS) module when the agent joins the RTC channel
//
// @since v0.12.0
type TTSOpenAIVendorParams struct {
	APIKey       string  `json:"api_key"`
	Model        string  `json:"model"`
	Voice        string  `json:"voice"`
	Instructions string  `json:"instructions"`
	Speed        float32 `json:"speed"`
}

func (TTSOpenAIVendorParams) VendorParam()             {}
func (TTSOpenAIVendorParams) GetVendorType() TTSVendor { return OpenAITTSVendor }

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
	//
	//  - TTSCartesiaVendorParams
	//
	//  - TTSOpenAIVendorParams
	//
	// @since v0.12.0
	Params TTSVendorParamsInterface `json:"params"`

	// Controls whether the TTS module skips bracketed content when reading LLM response text.
	//
	// This prevents the agent from vocalizing structural prompt information like tone indicators,
	// action descriptions, and system prompts, creating a more natural and immersive listening experience.
	//
	// Enable this feature by specifying one or more values:
	//
	// 1: Skip content in Chinese parentheses （ ）,unicode is \uFF08 \uFF09
	//
	// 2: Skip content in Chinese square brackets 【】,unicode is \u3010 \u3011
	//
	// 3: Skip content in parentheses (),unicode is \u0028 \u0029
	//
	// 4: Skip content in square brackets [ ],unicode is \u005B \u005D
	//
	// 5: Skip content in curly braces { },unicode is \u007B \u007D
	//
	// @since v0.12.0
	SkipPatterns []int `json:"skip_patterns,omitempty"`
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
	// Cartesia TTS vendor
	//
	// @since v0.12.0
	CartesiaTTSVendor TTSVendor = "cartesia"
	// OpenAI TTS vendor
	//
	// @since v0.12.0
	OpenAITTSVendor TTSVendor = "openai"
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
	// Default value is 32.
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
	//
	// Deprecated: Use [Parameters.SilenceConfig] instead
	//
	// @deprecated This field is deprecated since v0.11.0
	SilenceMessage *string `json:"silence_message,omitempty"`
	// LLM provider(Optional), supports the following settings:
	//
	// - "custom": Custom LLM provider.
	//   When you set this option, the agent includes the following fields, in addition to role and content when making requests to the custom LLM:
	//		-  turn_id: A unique identifier for each conversation turn. It starts from 0 and increments with each turn. One user-agent interaction corresponds to one turn_id.
	//		-  timestamp: The request timestamp, in milliseconds.
	// - "aliyun": Aliyun LLM provider.(Only available in China Mainland service region)
	//
	// - "bytedance": Bytedance LLM provider.(Only available in China Mainland service region)
	//
	// - "deepseek": DeepSeek LLM provider.(Only available in China Mainland service region)
	//
	// - "tencent": Tencent LLM provider.(Only available in China Mainland service region)
	//
	Vendor string `json:"vendor,omitempty"`

	// The request style for chat completion.(Optional)(Only available in global service region)
	//
	//  - "openai": OpenAI style.(Default)
	//
	//  - "gemini": Gemini style.
	//
	//  - "anthropic": Anthropic style.
	//
	//  - "dify": Dify style.
	//
	// @since v0.11.0
	Style string `json:"style,omitempty"`
}

// @brief Defines the Voice Activity Detection (VAD) configuration for the agent to join the RTC channel
//
// @since v0.7.0
//
// @deprecated This field is deprecated since v0.12.0
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

// @brief Defines the Automatic Speech Recognition (ASR) vendor for the agent to join the RTC channel
//
// @since v0.12.0
type ASRVendor string

const (
	// Fengming ASR vendor
	//
	// @since v0.12.0
	ASRVendorFengming ASRVendor = "fengming"
	// Tencent ASR vendor
	//
	// @since v0.12.0
	ASRVendorTencent ASRVendor = "tencent"
	// Microsoft ASR vendor
	//
	// @since v0.12.0
	ASRVendorMicrosoft ASRVendor = "microsoft"
	// Ares ASR vendor
	//
	// @since v0.12.0
	ASRVendorAres ASRVendor = "ares"
	// Deepgram ASR vendor
	//
	// @since v0.12.0
	ASRVendorDeepgram ASRVendor = "deepgram"
)

type ASRVendorParamsInterface interface {
	VendorParam()
	GetVendorType() ASRVendor
}

// @brief Defines the Automatic Speech Recognition (ASR) Fengming vendor parameter for the agent to join the RTC channel
//
// @since v0.12.0
type ASRFengmingVendorParam struct{}

func (ASRFengmingVendorParam) VendorParam()             {}
func (ASRFengmingVendorParam) GetVendorType() ASRVendor { return ASRVendorFengming }

// @brief Defines the Automatic Speech Recognition (ASR) Ares vendor parameter for the agent to join the RTC channel
//
// @since v0.12.0
type ASRAresVendorParam struct{}

func (ASRAresVendorParam) VendorParam()             {}
func (ASRAresVendorParam) GetVendorType() ASRVendor { return ASRVendorAres }

// @brief Defines the Automatic Speech Recognition (ASR) Tencent vendor parameter for the agent to join the RTC channel
//
// @since v0.12.0
type ASRTencentVendorParam struct {
	Key             string `json:"key"`
	AppId           string `json:"app_id"`
	Secret          string `json:"secret"`
	EngineModelType string `json:"engine_model_type"`
	VoiceId         string `json:"voice_id"`
}

func (ASRTencentVendorParam) VendorParam()             {}
func (ASRTencentVendorParam) GetVendorType() ASRVendor { return ASRVendorTencent }

// @brief Defines the Automatic Speech Recognition (ASR) Microsoft vendor parameter for the agent to join the RTC channel
//
// @since v0.12.0
type ASRMicrosoftVendorParam struct {
	Key        string   `json:"key"`
	Region     string   `json:"region"`
	Language   string   `json:"language"`
	PhraseList []string `json:"phrase_list"`
}

func (ASRMicrosoftVendorParam) VendorParam()             {}
func (ASRMicrosoftVendorParam) GetVendorType() ASRVendor { return ASRVendorMicrosoft }

// @brief Defines the Automatic Speech Recognition (ASR) deepgram vendor parameter for the agent to join the RTC channel
//
// @since v0.12.0
type ASRDeepgramVendorParam struct {
	Url      string `json:"url"`
	Key      string `json:"key"`
	Model    string `json:"model"`
	Language string `json:"language"`
}

func (ASRDeepgramVendorParam) VendorParam()             {}
func (ASRDeepgramVendorParam) GetVendorType() ASRVendor { return ASRVendorDeepgram }

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
	// ASR vendor, see ASRVendor for details
	//
	// @since v0.12.0
	Vendor ASRVendor `json:"vendor,omitempty"`

	// ASR vendor parameter description, see
	//
	//  - ASRVendorTencentParams
	//
	//  - ASRVendorMicrosoftParams
	//
	//  - ASRVendorAresParams
	//
	//  - ASRVendorDeepgramParams
	//
	// @since v0.12.0
	Params ASRVendorParamsInterface `json:"params,omitempty"`
}

// @brief Multi-modal language model (MLLM) configuration
//
// @since v0.12.0
type JoinPropertiesMLLMBody struct {
	// The WebSocket URL for OpenAI Realtime API(Required)
	Url string `json:"url"`

	// The API key used for authentication.(Required)
	//
	// Get your API key from the OpenAI Console.
	APIKey string `json:"api_key"`

	// Array of conversation items used for short-term memory management.(Optional)
	//
	// Uses the same structure as item.content from the OpenAI Realtime API(https://platform.openai.com/docs/api-reference/realtime-client-events/conversation/item/create).
	Messages []map[string]any `json:"messages"`

	// Additional MLLM configuration parameters.(Optional)
	//
	//  - Modalities override: The modalities setting in params is overridden by `input_modalities` and `output_modalities`.
	//
	//  - Turn detection override: The turn_detection setting in params is overridden by the `turn_detection` section outside of `mllm`.
	Params map[string]any `json:"params"`

	// The number of conversation history messages to maintain.
	//
	// Cannot exceed the model's context window.
	//
	// Default value is 32.
	MaxHistory *int `json:"max_history,omitempty"`

	// Input modalities for the MLLM (optional)
	//
	// 	- ["audio"]: Audio only (default)
	//
	//  - ["audio", "text"]: Audio and text.
	//
	InputModalities []string `json:"input_modalities,omitempty"`

	// Output format options:
	//
	//  - ["text", "audio"] for both text and voice responses.
	OutputModalities []string `json:"output_modalities,omitempty"`

	// Initial message the agent speaks when a user joins the channel.
	GreetingMessage string `json:"greeting_message,omitempty"`

	// MLLM provider identifier.
	//
	// Set to `openai` for OpenAI Realtime API.
	Vendor string `json:"vendor,omitempty"`

	// API request style.
	//
	// Set to `openai` for OpenAI Realtime API format.
	Style string `json:"style,omitempty"`
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
	//
	// Deprecated: Use [Parameters.SilenceConfig] instead
	//
	// @deprecated This field is deprecated since v0.11.0
	SilenceTimeout *int `json:"silence_timeout,omitempty"`

	// Agent user ID in the RTM channel
	//
	// Only valid when advanced_features.enable_rtm is true
	//
	// Deprecated: Use AgentRtcUId instead
	//
	// @deprecated This field is deprecated since v0.11.0
	AgentRtmUId *string `json:"agent_rtm_uid,omitempty"`
	// Advanced feature configurations (optional), see JoinPropertiesAdvancedFeaturesBody for details
	AdvancedFeatures *JoinPropertiesAdvancedFeaturesBody `json:"advanced_features,omitempty"`
	// Custom language model (LLM) configuration (required), see JoinPropertiesCustomLLMBody for details
	LLM *JoinPropertiesCustomLLMBody `json:"llm,omitempty"`
	// Multimodal Large Language Model (MLLM) configuration for real-time audio and text processing(optional), see JoinPropertiesMLLMBody for details
	//
	// @since v0.12.0
	MLLM *JoinPropertiesMLLMBody `json:"mllm,omitempty"`
	// Text-to-Speech (TTS) module configuration (required), see JoinPropertiesTTSBody for details
	TTS *JoinPropertiesTTSBody `json:"tts,omitempty"`
	// Voice Activity Detection (VAD) configuration (optional), see JoinPropertiesVadBody for details
	//
	// Deprecated: Use [TurnDetection] instead
	//
	// @deprecated This field is deprecated since v0.12.0
	Vad *JoinPropertiesVadBody `json:"vad,omitempty"`
	// Automatic Speech Recognition (ASR) configuration (optional), see JoinPropertiesAsrBody for details
	Asr *JoinPropertiesAsrBody `json:"asr,omitempty"`
	// Conversation turn detection settings
	TurnDetection *TurnDetectionBody `json:"turn_detection,omitempty"`
	// Agent parameters configuration (optional), see Parameters for details
	Parameters *Parameters `json:"parameters,omitempty"`
}

// @brief Conversation turn detection settings
//
// @since v0.11.0
type TurnDetectionBody struct {
	// Turn detection mechanism.(Optional)
	//
	//  - "agora_vad": Agora VAD.(Default)
	//
	//  - "server_vad": The model detects the start and end of speech based on audio volume and responds at the end of user speech.
	// 				 	Only available when mllm is enabled and OpenAI is selected.
	//
	//  - "semantic_vad": Uses a turn detection model in conjunction with VAD to semantically estimate whether the user has finished speaking,
	// 					then dynamically sets a timeout based on this probability for more natural conversations.
	// 					Only available when mllm is enabled and OpenAI is selected.
	//
	// @since v0.12.0
	Type string `json:"type"`

	// When the agent is interacting (speaking or thinking), the mode of human voice interrupting the agent's behavior, support the following values:
	//
	//  - "interrupt"(Default): Interrupt mode, human voice immediately interrupts the agent's interaction.
	//	               The agent will terminate the current interaction and directly process the human voice input.
	//
	//  - "append": Append mode, human voice does not interrupt the agent. (Default)
	//				The agent will process the human voice request after the current interaction ends.
	//
	//  - "ignore": Ignore mode, the agent ignores the human voice request.
	//				If the agent is speaking or thinking and receives human voice during the process,
	//				the agent will directly ignore and discard the human voice request, not storing it in the context.
	InterruptMode string `json:"interrupt_mode,omitempty"`

	// The amount of time in milliseconds that the user's voice must exceed the VAD threshold before an interruption is triggered.(Optional)
	//
	// Default value is 160.
	//
	// @since v0.12.0
	InterruptDurationMs *int `json:"interrupt_duration_ms,omitempty"`

	// The extra forward padding time in milliseconds before the processing system starts to process the speech input. This padding helps capture the beginning of the speech.
	//
	// Default value is 800.
	//
	// @since v0.12.0
	PrefixPaddingMs *int `json:"prefix_padding_ms,omitempty"`

	// The duration of audio silence in milliseconds.(Optional)
	//
	// If no voice activity is detected during this period, the agent assumes that the user has stopped speaking.
	//
	// Default value is 480.
	//
	// @since v0.12.0
	SilenceDurationMs *int `json:"silence_duration_ms,omitempty"`

	// Identification sensitivity determines the level of sound in the audio signal that is considered voice activity.(Optional)
	//
	// Lower values make it easier for the agent to detect speech, and higher values ignore weak sounds.
	//
	// The value range is (0.0, 1.0).
	//
	// Default value is 0.5.
	//
	// @since v0.12.0
	Threshold *float64 `json:"threshold,omitempty"`

	// Whether to automatically generate a response when a VAD stop event occurs.(Optional)
	//
	// Only available in server_vad and semantic_vad modes when using OpenAI Realtime API.
	//
	// Default value is true.
	//
	// @since v0.12.0
	CreateResponse *bool `json:"create_response,omitempty"`

	// Whether to automatically interrupt any ongoing response when a VAD start event occurs.
	//
	// Only available in server_vad and semantic_vad modes when using OpenAI Realtime API.
	//
	// Default value is true.
	//
	// @since v0.12.0
	InterruptResponse *bool `json:"interrupt_response,omitempty"`

	// The eagerness of the model to respond(Optional):
	//
	//  - "auto": Equivalent to medium(Default)
	//
	//  - "low": Wait longer for the user to continue speaking
	//
	//  - "high": Respond more quickly
	//
	// Only available in semantic_vad mode when using OpenAI Realtime API.
	//
	// @since v0.12.0
	Eagerness *string `json:"eagerness,omitempty"`
}

// @brief Fixed parameters
//
// @since v0.11.0
type FixedParams struct {
	// Silence configuration for the agent
	SilenceConfig *SilenceConfig `json:"silence_config,omitempty"`

	// Agent data transmission channel(Optional):
	//
	//  - "rtm": Use RTM transmission. This configuration takes effect only when advanced_features.enable_rtm is true.
	//
	//  - "datastream": Use RTC data stream transport.(Default)
	//
	// @since v0.12.0
	DataChannel *string `json:"data_channel,omitempty"`

	// Whether to receive agent performance data(Optional):
	//
	//  - true: Receive agent performance data.
	//
	//  - false: Do not receive agent performance data.(Default)
	//
	// This setting only takes effect when advanced_features.enable_rtm is true.
	//
	// @since v0.12.0
	EnableMetrics bool `json:"enable_metrics,omitempty"`

	// Whether to receive agent error events(Optional):
	//
	//  - true: Receive agent error events.
	//
	//  - false: Do not receive agent error events.(Default)
	//
	// This setting only takes effect when advanced_features.enable_rtm is true
	//
	// @since v0.12.0
	EnableErrorMessage bool `json:"enable_error_message,omitempty"`
}

// @brief Silence configuration for the agent
//
// @since v0.11.0
type SilenceConfig struct {
	// Agent maximum silence time (ms).(Optional)
	//
	// After the agent is created and a user joins the channel,
	// the duration of the agent's non-listening, thinking, or speaking state is called the agent's silence time.
	//
	// When the silence time reaches the set value, the agent will report the silence prompt message.
	//
	// This feature can be used to let the agent remind users when users are inactive.
	//
	// Set 0: Do not enable this feature.
	//
	// Set to (0,60000]: Must also set content to enable normal reporting of silence prompts, otherwise the setting is invalid.
	TimeoutMs *int `json:"timeout_ms,omitempty"`

	// When the silence time reaches the set value, the agent will take the following actions(Optional):
	//
	//  - "speak": Use TTS module to report the silence message (Default)
	//
	//  - "think": Append the silence message to the end of the context and pass it to LLM
	Action *string `json:"action,omitempty"`

	// Content of the silence message (Optional)
	//
	// The content will be used in different ways according to the settings in the action.
	Content *string `json:"content,omitempty"`
}

// @brief Agent parameters configuration
//
// @note Parameters that contains both extra data and fixed data. The same key in extra data and fixed data will be merged.
//
// @since v0.11.0
type Parameters struct {
	// Extra parameters for flexible key-value pairs
	ExtraParams map[string]any `json:"-"`
	// Fixed parameters for type-safe parameters
	FixedParams *FixedParams `json:"-"`
}

// MarshalJSON implements custom JSON marshaling
func (p *Parameters) MarshalJSON() ([]byte, error) {
	// Create a map to hold the merged data
	merged := make(map[string]any)

	// Add fixed parameters if present
	if p.FixedParams != nil {
		structBytes, err := json.Marshal(p.FixedParams)
		if err != nil {
			return nil, err
		}
		var structMap map[string]any
		if err := json.Unmarshal(structBytes, &structMap); err != nil {
			return nil, err
		}
		for k, v := range structMap {
			merged[k] = v
		}
	}

	// Add extra parameters if present
	if p.ExtraParams != nil {
		for k, v := range p.ExtraParams {
			merged[k] = v
		}
	}

	return json.Marshal(merged)
}

// UnmarshalJSON implements custom JSON unmarshaling
func (p *Parameters) UnmarshalJSON(data []byte) error {
	var mapData map[string]any
	if err := json.Unmarshal(data, &mapData); err != nil {
		return err
	}
	p.ExtraParams = mapData
	return nil
}
