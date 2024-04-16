// internal/utils/utils.go
package utils

import (
	"fmt"
	"strings"
)

func PromptLLM(description string, existingInfo map[string][]string) string {
	// Constructing the prompt string for the LLM with clear instructions.
	prompt := fmt.Sprintf(
		"Transaction description: '%s'.\n"+
			"Task:\n"+
			"- Extract relevant information from the transaction and format it as JSON.\n"+
			"- Required fields include date, payee, notes, category, and amount.\n"+
			"- Date should be in the format 'YYYY-MM-DD'.\n"+
			"- Amount should be a number with two decimal places. B really carefull to not break the value.\n"+
			"- Category should be a comma-separated list of appropriate categories.\n"+
			"- Return the response in this structured JSON format:\n"+
			"  {\n"+
			"    'date': 'YYYY-MM-DD',\n"+
			"    'payee': 'Name of the payee',\n"+
			"    'notes': 'Description of the transaction',\n"+
			"    'category': 'Appropriate categories',\n"+
			"    'amount': 'Transaction amount'\n"+
			"  }\n"+
			"- If there is no new information, return empty fields.\n"+
			"- Ensure that the categories are comma-separated and do not contain spaces.\n"+
			"- JUST RETURN THE JSON OBJECT.\n",
		description)

	// Append existing information dynamically based on what's provided.
	keys := []string{"categories", "payee", "notes"}
	for _, key := range keys {
		if items, exists := existingInfo[key]; exists && len(items) > 0 {
			joinedItems := strings.Join(items, ", ")
			prompt += fmt.Sprintf("\nExisting %s: [%s].", key, joinedItems)
		}
	}

	return strings.TrimSpace(prompt)
}
