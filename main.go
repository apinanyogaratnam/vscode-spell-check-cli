package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
)

type Settings struct {
	Words []string `json:"cSpell.words"`
}

func main() {
	// Run cspell command
	cmd := exec.Command("npx", "--yes", "cspell", "lint", "src/**/*.rs")
	output, err := cmd.Output()

	if err != nil {
		log.Println("Failed to run cspell:", err)
	}

	// Parse output
	pattern := regexp.MustCompile(`Unknown word \((.*?)\)`)
	unknownWords := pattern.FindAllString(string(output), -1)

	settingsPath := ".vscode/settings.json"

	// Load existing settings or create a new one
	settings := Settings{}
	if _, err := os.Stat(settingsPath); err == nil {
		settingsFile, err := ioutil.ReadFile(settingsPath)
		if err != nil {
			log.Fatal("Failed to read settings:", err)
		}
		if err := json.Unmarshal(settingsFile, &settings); err != nil {
			log.Fatal("Failed to unmarshal settings:", err)
		}
		if settings.Words == nil {
			log.Println("Warning: 'cSpell.words' key not found in settings file. A new key will be added.")
		}
		settings.Words = append(settings.Words, unknownWords...)
	} else {
		settings.Words = unknownWords
	}

	// Remove duplicates
	settings.Words = removeDuplicates(settings.Words)

	// Create the .vscode directory if it does not exist
	if _, err := os.Stat(".vscode"); os.IsNotExist(err) {
		err = os.MkdirAll(".vscode", 0755)
		if err != nil {
			log.Fatal("Failed to create .vscode directory:", err)
		}
	}

	// Write back to settings.json
	settingsBytes, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal settings:", err)
	}
	if err := ioutil.WriteFile(settingsPath, settingsBytes, 0644); err != nil {
		log.Fatal("Failed to write settings:", err)
	}

	fmt.Println("Done!")
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if !encountered[elements[v]] {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}
