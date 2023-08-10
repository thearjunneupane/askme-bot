package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/thearjnep/askme-bot/config"
)

var dgBot *discordgo.Session

func Initialize() {
	var err error
	dgBot, err = discordgo.New("Bot " + config.Bot_Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dgBot.AddHandler(messageHandler)

	dgBot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = dgBot.Open()

	if err != nil {
		fmt.Println("Error Opening Connection...!\n", err.Error())
		return
	}

	fmt.Println("Bot is running...!")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, config.Bot_Prefix) {
		userMessage := strings.TrimSpace(strings.TrimPrefix(m.Content, config.Bot_Prefix))
		botResponse := getGPT3Response(userMessage)
		s.ChannelMessageSend(m.ChannelID, botResponse)
	}

}

func getGPT3Response(prompt string) string {

	data := map[string]interface{}{
		"prompt": prompt,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "An error occurred while generating response."
	}

	req, err := http.NewRequest("POST", config.Gpt_api_url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "An error occurred while generating response."
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Gpt_api_key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to GPT-3:", err)
		return "An error occurred while generating response."
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "An error occurred while generating response."
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println("Error parsing response JSON:", err)
		return "An error occurred while generating response."
	}

	botResponseValue, ok := responseMap["choices"]
	if !ok {
		return "No bot response available."
	}

	botResponse, ok := botResponseValue.([]interface{})
	if !ok || len(botResponse) == 0 {
		return "No bot response available."
	}

	firstChoice, ok := botResponse[0].(map[string]interface{})
	if !ok {
		return "Unexpected data format for the bot response."
	}

	botTextValue, ok := firstChoice["text"]
	if !ok {
		return "No bot text available."
	}

	botText, ok := botTextValue.(string)
	if !ok {
		return "Unexpected data format for the bot text."
	}

	return botText
}
