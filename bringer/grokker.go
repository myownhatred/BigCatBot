package bringer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"Guenhwyvar/config"
	"Guenhwyvar/lib/mlog"
)

type grokMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type grokRequest struct {
	Messages    []grokMessage `json:"messages"`
	Model       string        `json:"model"`
	Stream      bool          `json:"stream"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type grokResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails struct {
			TextTokens   int `json:"text_tokens"`
			AudioTokens  int `json:"audio_tokens"`
			ImageTokens  int `json:"image_tokens"`
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
		CompletionTokensDetails struct {
			ReasoningTokens          int `json:"reasoning_tokens"`
			AudioTokens              int `json:"audio_tokens"`
			AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
			RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
		} `json:"completion_tokens_details"`
		NumSourcesUsed int `json:"num_sources_used"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

// AnalysisResult contains both summary and topic starters
type AnalysisResult struct {
	Summary       string `json:"summary"`
	TopicStarters string `json:"topic_starters"`
}

// GrokClient handles Grok API interactions
type GrokClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// type for contextual talks with grok
type Conversation struct {
	messages    []grokMessage
	model       string
	temperature float64
	maxTokens   int
}

// NewConversation creates a new Conversation instance.
// You can optionally provide a system prompt as the first message.
func NewConversation(model string, temperature float64, maxTokens int, systemPrompt string) *Conversation {
	conv := &Conversation{
		model:       model,
		temperature: temperature,
		maxTokens:   maxTokens,
	}
	if systemPrompt != "" {
		conv.messages = append(conv.messages, grokMessage{Role: "system", Content: systemPrompt})
	}
	return conv
}

// AppendMessage appends a message to the conversation history without sending it.
// Use this for manual context management if needed.
func (conv *Conversation) AppendMessage(role string, content string) error {
	if role != "user" && role != "system" && role != "assistant" {
		return fmt.Errorf("invalid role: %s; must be 'user', 'system', or 'assistant'", role)
	}
	conv.messages = append(conv.messages, grokMessage{Role: role, Content: content})
	return nil
}

// GetHistory returns the current message history.
func (conv *Conversation) GetHistory() []grokMessage {
	return append([]grokMessage(nil), conv.messages...) // Copy to avoid mutation
}

// ClearHistory clears the conversation history, keeping only the system prompt if present.
func (conv *Conversation) ClearHistory() {
	if len(conv.messages) > 0 && conv.messages[0].Role == "system" {
		conv.messages = conv.messages[:1]
	} else {
		conv.messages = nil
	}
}

// NewGrokClient creates a new Grok API client
func NewGrokClient(apiKey string) *GrokClient {
	return &GrokClient{
		APIKey:  apiKey,
		BaseURL: "https://api.x.ai/v1/chat/completions",
		Client:  &http.Client{Timeout: 600 * time.Second},
	}
}

func (c *GrokClient) ChatCompletion(req *grokRequest) (*grokResponse, error) {
	req.Stream = false // Enforce non-streaming for this implementation

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.BaseURL // Assuming endpoint is /chat/completions, similar to OpenAI

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	var grokResp grokResponse
	if err := json.NewDecoder(resp.Body).Decode(&grokResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &grokResp, nil
}

func (gc *GrokClient) callGenericGrokAPI(prompt string) (string, string, error) {
	request := grokRequest{
		Messages: []grokMessage{
			{
				Role:    "system",
				Content: "ты помошник в чате весёлых друзей, отвечай правдиво честно и прямо, можешь использовать пошлый юмор",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:       "grok-4-1-fast-non-reasoning-latest",
		Stream:      false,
		MaxTokens:   7000,
		Temperature: 0.8,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", gc.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+gc.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := gc.Client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var grokResp grokResponse
	if err := json.Unmarshal(body, &grokResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(grokResp.Choices) == 0 {
		return "", "", fmt.Errorf("no response choices returned")
	}

	// Calculate costs
	inputCost := float64(grokResp.Usage.PromptTokensDetails.TextTokens+grokResp.Usage.PromptTokensDetails.ImageTokens)*3.00/1_000_000 +
		float64(grokResp.Usage.PromptTokensDetails.CachedTokens)*0.75/1_000_000
	outputCost := float64(grokResp.Usage.CompletionTokens) * 15.00 / 1_000_000
	totalCost := inputCost + outputCost

	output := fmt.Sprintf(`
	Usage Data:
  Prompt Tokens: %d
  Completion Tokens: %d
  Total Tokens: %d
  Prompt Tokens Details:
    Text Tokens: %d
    Audio Tokens: %d
    Image Tokens: %d
    Cached Tokens: %d
  Completion Tokens Details:
    Reasoning Tokens: %d
    Audio Tokens: %d
    Accepted Prediction Tokens: %d
    Rejected Prediction Tokens: %d
  Num Sources Used: %d
System Fingerprint: %s
Cost Breakdown:
  Input Cost: $%.6f
  Output Cost: $%.6f
  Total Cost: $%.6f`,
		grokResp.Usage.PromptTokens,
		grokResp.Usage.CompletionTokens,
		grokResp.Usage.TotalTokens,
		grokResp.Usage.PromptTokensDetails.TextTokens,
		grokResp.Usage.PromptTokensDetails.AudioTokens,
		grokResp.Usage.PromptTokensDetails.ImageTokens,
		grokResp.Usage.PromptTokensDetails.CachedTokens,
		grokResp.Usage.CompletionTokensDetails.ReasoningTokens,
		grokResp.Usage.CompletionTokensDetails.AudioTokens,
		grokResp.Usage.CompletionTokensDetails.AcceptedPredictionTokens,
		grokResp.Usage.CompletionTokensDetails.RejectedPredictionTokens,
		grokResp.Usage.NumSourcesUsed,
		grokResp.SystemFingerprint,
		inputCost,
		outputCost,
		totalCost,
	)

	return grokResp.Choices[0].Message.Content, output, nil

}

func (gc *GrokClient) callBiogenGrokAPI(prompt string) (string, string, error) {
	prompt = "создай краткую (не больше 10 предложений) биогафию для данного персонажа:" + prompt
	request := grokRequest{
		Messages: []grokMessage{
			{
				Role:    "system",
				Content: fmt.Sprintf("ты писатель с уклоном в черный юмор и гэги на грани: %s и эпический размах в сочетании с мелкими деталями для контраста. твоя задача писать короткие биографии по листкам характеристик персонажей днд, придумывай становление и мотивацию персонажа и какую-то фишку ", strings.Join(tags(), ", ")),
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:       "grok-4-latest",
		Stream:      false,
		MaxTokens:   3000,
		Temperature: 0.9,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", gc.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+gc.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := gc.Client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var grokResp grokResponse
	if err := json.Unmarshal(body, &grokResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(grokResp.Choices) == 0 {
		return "", "", fmt.Errorf("no response choices returned")
	}

	// Calculate costs
	inputCost := float64(grokResp.Usage.PromptTokensDetails.TextTokens+grokResp.Usage.PromptTokensDetails.ImageTokens)*3.00/1_000_000 +
		float64(grokResp.Usage.PromptTokensDetails.CachedTokens)*0.75/1_000_000
	outputCost := float64(grokResp.Usage.CompletionTokens) * 15.00 / 1_000_000
	totalCost := inputCost + outputCost

	output := fmt.Sprintf(`
	Usage Data:
  Prompt Tokens: %d
  Completion Tokens: %d
  Total Tokens: %d
  Prompt Tokens Details:
    Text Tokens: %d
    Audio Tokens: %d
    Image Tokens: %d
    Cached Tokens: %d
  Completion Tokens Details:
    Reasoning Tokens: %d
    Audio Tokens: %d
    Accepted Prediction Tokens: %d
    Rejected Prediction Tokens: %d
  Num Sources Used: %d
System Fingerprint: %s
Cost Breakdown:
  Input Cost: $%.6f
  Output Cost: $%.6f
  Total Cost: $%.6f`,
		grokResp.Usage.PromptTokens,
		grokResp.Usage.CompletionTokens,
		grokResp.Usage.TotalTokens,
		grokResp.Usage.PromptTokensDetails.TextTokens,
		grokResp.Usage.PromptTokensDetails.AudioTokens,
		grokResp.Usage.PromptTokensDetails.ImageTokens,
		grokResp.Usage.PromptTokensDetails.CachedTokens,
		grokResp.Usage.CompletionTokensDetails.ReasoningTokens,
		grokResp.Usage.CompletionTokensDetails.AudioTokens,
		grokResp.Usage.CompletionTokensDetails.AcceptedPredictionTokens,
		grokResp.Usage.CompletionTokensDetails.RejectedPredictionTokens,
		grokResp.Usage.NumSourcesUsed,
		grokResp.SystemFingerprint,
		inputCost,
		outputCost,
		totalCost,
	)

	return grokResp.Choices[0].Message.Content, output, nil

}

func (gc *GrokClient) callParGrokAPI(prompt, role string, temp float64) (string, string, error) {
	request := grokRequest{
		Messages: []grokMessage{
			{
				Role:    "system",
				Content: role,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:       "grok-4-latest",
		Stream:      false,
		MaxTokens:   3000,
		Temperature: temp,
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", gc.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+gc.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := gc.Client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var grokResp grokResponse
	if err := json.Unmarshal(body, &grokResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(grokResp.Choices) == 0 {
		return "", "", fmt.Errorf("no response choices returned")
	}

	// Calculate costs
	inputCost := float64(grokResp.Usage.PromptTokensDetails.TextTokens+grokResp.Usage.PromptTokensDetails.ImageTokens)*3.00/1_000_000 +
		float64(grokResp.Usage.PromptTokensDetails.CachedTokens)*0.75/1_000_000
	outputCost := float64(grokResp.Usage.CompletionTokens) * 15.00 / 1_000_000
	totalCost := inputCost + outputCost

	output := fmt.Sprintf(`
	Usage Data:
  Prompt Tokens: %d
  Completion Tokens: %d
  Total Tokens: %d
  Prompt Tokens Details:
    Text Tokens: %d
    Audio Tokens: %d
    Image Tokens: %d
    Cached Tokens: %d
  Completion Tokens Details:
    Reasoning Tokens: %d
    Audio Tokens: %d
    Accepted Prediction Tokens: %d
    Rejected Prediction Tokens: %d
  Num Sources Used: %d
System Fingerprint: %s
Cost Breakdown:
  Input Cost: $%.6f
  Output Cost: $%.6f
  Total Cost: $%.6f`,
		grokResp.Usage.PromptTokens,
		grokResp.Usage.CompletionTokens,
		grokResp.Usage.TotalTokens,
		grokResp.Usage.PromptTokensDetails.TextTokens,
		grokResp.Usage.PromptTokensDetails.AudioTokens,
		grokResp.Usage.PromptTokensDetails.ImageTokens,
		grokResp.Usage.PromptTokensDetails.CachedTokens,
		grokResp.Usage.CompletionTokensDetails.ReasoningTokens,
		grokResp.Usage.CompletionTokensDetails.AudioTokens,
		grokResp.Usage.CompletionTokensDetails.AcceptedPredictionTokens,
		grokResp.Usage.CompletionTokensDetails.RejectedPredictionTokens,
		grokResp.Usage.NumSourcesUsed,
		grokResp.SystemFingerprint,
		inputCost,
		outputCost,
		totalCost,
	)

	return grokResp.Choices[0].Message.Content, output, nil

}

// callGrokAPI makes a request to Grok API
func (gc *GrokClient) callGrokAPI(prompt string) (string, error) {
	request := grokRequest{
		Messages: []grokMessage{
			{
				Role:    "system",
				Content: "Ты няшный анализатор логов чата. Испольуешь UwU-лексику и когда используешь мат потом себя поправляешь няшно.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:       "grok-4-latest",
		Stream:      false,
		MaxTokens:   8000,
		Temperature: 0.5,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", gc.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+gc.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := gc.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var grokResp grokResponse
	if err := json.Unmarshal(body, &grokResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(grokResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return grokResp.Choices[0].Message.Content, nil
}

// formatMessagesForAnalysis converts messages to a readable format
func formatMessagesForAnalysis(messages []mlog.Mlog, maxMessages int) string {
	// Take the most recent messages if there are too many
	if len(messages) > maxMessages {
		messages = messages[len(messages)-maxMessages:]
	}

	var result bytes.Buffer
	result.WriteString("Telegram Chat Messages:\n\n")

	for _, msg := range messages {
		timestamp := msg.Timestamp
		formattedTime := timestamp.Format("2006-01-02 15:04:05")
		result.WriteString(fmt.Sprintf("[%s] %s: %s (ID: %d)\n",
			formattedTime, msg.Sender, msg.Content, msg.MessageId))
	}

	return result.String()
}

// AnalyzeChatSummary generates a comprehensive summary of the chat
func (gc *GrokClient) AnalyzeChatSummary(messages []mlog.Mlog) (string, error) {
	formattedMessages := formatMessagesForAnalysis(messages, 3000)

	prompt := fmt.Sprintf(`
  Используй неформальный стиль, возможно слегка пошлый и весёлый. Отчёт должен быть на русском языке.
  Но некоторые термины можно вставлять на английском если это смешно. 
  Проанализируй лог чата, выдели темы на которые велись разговоры и краткий итог каждого
  Можно разбить на временные интервалы - утро, день, вечер, ночь
  Для обширных тем выделяй ник топикстартера и отметь ID с которого началась тема. В половине случаев после ника топикстартера добавляй "хуй, блять"
  перед ID ставь \"msg:\/\/\"  
  Освещай темы по возможности кратко в 2-4 предложений, но если тема обширная, можно осветить её более подробноы
  Выводы предоставь без технической информации, в простом человеческом формате
	Если обсуждалось что-то крутое, вставляй в конец отчёта о крутом слово "ебём".
  Не нужно писать приамбулу в духе - я составил отчёт и всё такое - сразу текст отчёта без лишних слов 
  Утро: обсудили то-то и то-то
  День: были такие-то темы
  Нельзя превышать объём отчёта в 4096 символов, но лучше делать его более кратким, в идеале на одну печатную страницу

Chat data:
%s`, formattedMessages)

	return gc.callGrokAPI(prompt)
}

// AnalyzeChat performs complete analysis including both summary and topic identification
func (gc *GrokClient) AnalyzeChat(messages []mlog.Mlog) (*AnalysisResult, error) {
	fmt.Printf("Analyzing %d messages...\n", len(messages))

	// Generate summary
	fmt.Println("Generating summary...")
	summary, err := gc.AnalyzeChatSummary(messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	return &AnalysisResult{
		Summary:       summary,
		TopicStarters: "",
	}, nil
}

// simple ansver
func (g *GrokGrokker) SimpleAnswer(prompt string) (string, string, error) {
	result, debug, err := g.Grik.callGenericGrokAPI(prompt)
	if err != nil {
		return "", "", fmt.Errorf("ошибка апизации зопроса: %w", err)
	}

	return result, debug, nil
}

func (g *GrokGrokker) DNDBiogen(prompt string) (string, string, error) {
	result, debug, err := g.GrikDND.callBiogenGrokAPI(prompt)
	if err != nil {
		return "", "", fmt.Errorf("ошибка апизации биогена: %w", err)
	}

	return result, debug, nil
}

func (g *GrokGrokker) GenGrok(prompt, role string, temp float64) (string, string, error) {
	result, debug, err := g.GrikDND.callParGrokAPI(prompt, role, temp)
	if err != nil {
		return "", "", fmt.Errorf("ошибка вызова грока: %w", err)
	}
	return result, debug, nil
}

func (g *GrokGrokker) GetOrCreateConversation(chatID int64, model string, temperature float64, maxTokens int, systemPrompt string) *Conversation {
	g.mu.Lock()
	defer g.mu.Unlock()
	if conv, ok := g.convs[chatID]; ok {
		return conv
	}
	conv := NewConversation(model, temperature, maxTokens, systemPrompt)
	g.convs[chatID] = conv
	return conv
}

func (g *GrokGrokker) SendMessageInConversation(chatID int64, content string) (string, error) {
	conv := g.GetOrCreateConversation(chatID, "grok-4-latest", 0.2, 8000, "") // Defaults; load from config if needed

	if err := conv.AppendMessage("user", content); err != nil {
		return "", err
	}

	req := &grokRequest{
		Messages:    conv.messages,
		Model:       conv.model,
		Temperature: conv.temperature,
		MaxTokens:   conv.maxTokens,
	}

	client := g.Grik

	resp, err := client.ChatCompletion(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	assistantContent := resp.Choices[0].Message.Content
	if err := conv.AppendMessage("assistant", assistantContent); err != nil {
		return "", err
	}

	return assistantContent, nil
}

// ClearConversation clears the history for a specific chatID.
func (g *GrokGrokker) ClearConversation(chatID int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if conv, ok := g.convs[chatID]; ok {
		conv.ClearHistory()
	}
}

// DeleteConversation removes the conversation for a chatID.
func (g *GrokGrokker) DeleteConversation(chatID int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.convs, chatID)
}

type GrokGrokker struct {
	Grik    *GrokClient
	GrikDND *GrokClient
	convs   map[int64]*Conversation
	mu      sync.Mutex
}

func NewGrokker(c *config.AppConfig) *GrokGrokker {
	keys := c.Grok.Grokeys
	return &GrokGrokker{
		Grik:    NewGrokClient(keys[1]),
		GrikDND: NewGrokClient(keys[0]),
		convs:   make(map[int64]*Conversation),
	}
}

func (g *GrokGrokker) AnalChatDay(filename string) (string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var messages []mlog.Mlog

	err = json.Unmarshal(file, &messages)
	if err != nil {
		return "", err
	}

	// Analyze chat
	result, err := g.Grik.AnalyzeChat(messages)
	if err != nil {
		return "", fmt.Errorf("ошибка анализатции чата: %v", err)
	}

	return result.Summary, nil

	// Optionally save to JSON file
	//jsonResult, err := json.MarshalIndent(result, "", "  ")
	//if err == nil {
	//  fmt.Printf("\nJSON Result:\n%s\n", string(jsonResult))
	//}

}

func tags() []string {
	arr := []string{"расизм", "геноцид", "боллбастинг",
		"фистинг", "зоофилия", "отсос", "ебля в попец", "зрада",
		"неудачи украины", "украинские проститутки",
		"использование исковерканных имён украинских политиков и общественных деятелей",
		"заманивание на хуйню", "карманный биллиард", "катафрактарии", "ачинск",
		"пердеж", "литургия"}

	n := len(arr)
	size := 8

	indices := rand.Perm(n)[:size]

	sort.Ints(indices)

	result := make([]string, size)
	for i, idx := range indices {
		result[i] = arr[idx]
	}
	return result
}
