package utils


import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// Could've been a great additional logic but it slows down the PC (with llama3 on ollama) - Keeping it here in case summarizing functionality needed for video description.
// Helps keep consistency in description across all videos.
func Summarize(text string) string {
	summaryPrompt := "Your task is to summarize the content provided below within 200 words. Ignore emoticons, emojis and using any starting text like \"This content highlights\" etc. Directly start off with the text summary. Here's the content: " + text

	llm, err := ollama.New(ollama.WithModel("llama3"))
	if err != nil {
		log.Fatal(err)
	}
	completion, err := llms.GenerateFromSinglePrompt(
		context.TODO(),
		llm,
		summaryPrompt,
		llms.WithTemperature(0.7),
		llms.WithMinLength(150),
		llms.WithMaxLength(200),
	)
	if err != nil {
		log.Fatal(err)
	}
	return completion
}

func sanitizeTextByRemovingURLs(text string) string {
	re := regexp.MustCompile(`https?:\/\/(?:www\.)?\S+`)
	return re.ReplaceAllStringFunc(text, func(s string) string {
		words := strings.Split(s, " ")
		for i := range words {
			if i == len(words)-1 {
				break
			}
			text := strings.Join(words[i+1:], " ")
			if text != "" {
				return words[i] + " " + text
			}
		}
		return ""
	})
}

func SummarizeText(text string) string {
	if len(strings.Split(text, " ")) < 20 {
		return sanitizeTextByRemovingURLs(text) + "..."
	}
	summarySlice := strings.Split(sanitizeTextByRemovingURLs(text), " ")
	return strings.Join(summarySlice[:19], " ") + "..."
}

func ShortenViewCount(viewCount string) string {
	intViewCount, _ := strconv.Atoi(viewCount)
	if intViewCount > 1000000 {
		return fmt.Sprintf("%.1fM", float64(intViewCount)/1000000)
	} else if intViewCount > 1000 {
		return fmt.Sprintf("%.1fk", float64(intViewCount)/1000)
	} else {
		return fmt.Sprintf("%d", intViewCount)
	}
}

func ShortenStarCount(starCount string) string {
	intStarCount, _ := strconv.Atoi(starCount)
	if intStarCount >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(intStarCount)/1000000)
	} else if intStarCount >= 1000 {
		return fmt.Sprintf("%.1fk", float64(intStarCount)/1000)
	} else {
		return fmt.Sprintf("%d", intStarCount)
	}
}
