// internal/utils/utils.go
package utils

import (
	"fmt"
	"strings"
)

func PromptLLM(description string, existingInfo map[string][]string) string {
	// Constructing the prompt string for the LLM with clear instructions.
	prompt := fmt.Sprintf(
		"Analyze the following transaction description to extract structured information:\n\n"+
			"Description: '%s'\n\n"+
			"Instructions:\n"+
			"1. Extract key data points relevant to a financial transaction.\n"+
			"2. Format this information into a JSON object with specific fields and formats.\n"+
			"3. Ensure all fields are accurately populated based on the description provided:\n"+
			"   - 'date' should be formatted as 'YYYY-MM-DD'.\n"+
			"   - 'payee' should clearly identify the entity involved in the transaction.\n"+
			"   - 'notes' should include any descriptive information about the transaction.\n"+
			"   - 'category' should list transaction categories, separated by commas without additional spaces.\n"+
			"			* I need you to always return the categories furthermore, I need budget related categories as well"+
			"   - 'amount' should be represented as a numeric value with two decimal places.\n\n"+
			"Format Requirements:\n"+
			"Return the JSON object with the following structure, ensuring to adhere to the format guidelines.\n"+
			"If any field does not have corresponding data in the description, return it as an empty string or null:\n\n"+
			"{\n"+
			"  'date': 'appropriate date',\n"+
			"  'payee': 'transaction entity',\n"+
			"  'notes': 'transaction details',\n"+
			"  'category': 'listed,categories',\n"+
			"  'amount': 'formatted amount'\n"+
			"}\n\n"+
			"Note: The JSON object should strictly follow this structure with no deviation in key names or data types.\n"+
			"Handle missing data gracefully by returning appropriate placeholders.\n",
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
