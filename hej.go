package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Choice struct {
	Message Message `json:"message"`
}

type OpenRouterResponse struct {
	Choices []Choice `json:"choices"`
}

type Fallback struct {
	Cmd  string
	Desc string
	Tags string
}

var fallbackTable = []Fallback{
	{"ls -la", "List files with details (including hidden).", "list dir files hidden"},
	{"pwd", "Print working directory.", "where am i current directory path"},
	{"df -h", "Show filesystem disk space usage.", "disk free capacity mount"},
	{"du -h --max-depth=1", "Disk usage per subdir (human-readable).", "disk usage size folder space"},
	{"find . -type f -size +1G", "Find files larger than 1GB.", "find large files size"},
	{"ps aux | grep -i <name>", "Search processes by name.", "process list running grep"},
	{"top", "Interactive process viewer.", "cpu memory processes monitor"},
	{"docker ps -a", "List containers (all).", "docker list containers"},
	{"docker logs -f <container>", "Follow container logs.", "docker logs follow"},
	{"docker exec -it <container> /bin/sh", "Shell into container.", "docker exec shell bash sh"},
	{"iptables -L -n -v", "List firewall rules.", "firewall iptables rules"},
	{"ss -tulpn", "Show listening sockets.", "ports listening networking"},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hej \"your query\"")
		fmt.Println("Note:")
		fmt.Println("Set OPENROUTER_API_KEY env var to enable online mode.")
		fmt.Println("Set OPENROUTER_MODEL env var to select the model. Defaults to `gpt-4o-mini`.")
		return
	}

	query := strings.Join(os.Args[1:], " ")

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	model := os.Getenv("OPENROUTER_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	if apiKey != "" {
		result, err := callOpenRouter(apiKey, model, query)
		if err == nil && result != "" {
			lines := strings.Split(strings.TrimSpace(result), "\n")
			if len(lines) >= 2 {
				fmt.Println(lines[0])
				fmt.Println(lines[1])
				return
			}
		}
	}

	// Fallback mode
	if fb := fallbackLookup(query); fb != nil {
		fmt.Println(fb.Cmd)
		fmt.Println(fb.Desc)
	} else {
		fmt.Println("No offline match. Try rephrasing.")
	}
}

func callOpenRouter(apiKey, model, query string) (string, error) {
	systemPrompt := `You are a Linux command helper for Ubuntu servers.
Return exactly two lines:
1) The single best shell command (bash-safe) with no backticks.
2) A one-line description.
Never add extra commentary or examples.`

	reqBody := OpenRouterRequest{
		Model: model,
		Messages: []Message{
			{"system", systemPrompt},
			{"user", query},
		},
		Temperature: 0.2,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://local.hej")
	req.Header.Set("X-Title", "hej")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error: %s", resp.Status)
	}

	data, _ := io.ReadAll(resp.Body)
	var parsed OpenRouterResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) > 0 {
		return parsed.Choices[0].Message.Content, nil
	}
	return "", nil
}

func fallbackLookup(query string) *Fallback {
	q := strings.ToLower(query)
	var best *Fallback
	bestScore := 0
	for i := range fallbackTable {
		score := 0
		for _, tag := range strings.Split(fallbackTable[i].Tags, " ") {
			if strings.Contains(q, tag) {
				score++
			}
		}
		if score > bestScore {
			bestScore = score
			best = &fallbackTable[i]
		}
	}
	if bestScore == 0 {
		return nil
	}
	return best
}
