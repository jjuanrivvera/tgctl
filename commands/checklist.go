package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "checklist",
		Short: "Send and edit interactive checklists",
		Long:  "Send and edit interactive checklists on behalf of a connected business account.",
		Cmds: []methodCmd{
			{
				Use: "send", Method: "sendChecklist", Kind: kindWrite,
				Short:   "Send a checklist on behalf of a business account",
				Long:    "Send an interactive checklist on behalf of a connected business account. Returns the sent Message.",
				Example: `  tgctl checklist send --business-connection-id bc123 --chat 12345678 --checklist '{"title":"Tasks","tasks":[{"id":1,"text":"First task"}]}'`,
				Flags: []flagSpec{
					businessConnectionFlag(),
					chatFlag(),
					{Name: "checklist", Kind: flagJSON, Required: true, Usage: "JSON object describing the checklist to send (InputChecklist)"},
					threadFlag(),
					silentFlag(),
					protectContentFlag(),
					{Name: "message-effect-id", Param: "message_effect_id", Usage: "unique identifier of the message effect"},
					{Name: "reply-parameters", Param: "reply_parameters", Kind: flagJSON, Usage: "JSON object describing message to reply to"},
					{Name: "reply-markup", Param: "reply_markup", Kind: flagJSON, Usage: "JSON object for an inline keyboard"},
				},
				PreCall: requireValidInputChecklist,
			},
			{
				Use: "edit", Method: "editMessageChecklist", Kind: kindWrite,
				Short:   "Edit a checklist sent on behalf of a business account",
				Long:    "Edit an interactive checklist message sent on behalf of a connected business account. Returns the edited Message.",
				Example: `  tgctl checklist edit --business-connection-id bc123 --chat 12345678 --message-id 42 --checklist '{"title":"Updated Tasks","tasks":[{"id":1,"text":"First task"}]}'`,
				Flags: []flagSpec{
					businessConnectionFlag(),
					chatFlag(),
					messageIDFlag(),
					{Name: "checklist", Kind: flagJSON, Required: true, Usage: "JSON object describing the updated checklist (InputChecklist)"},
					{Name: "reply-markup", Param: "reply_markup", Kind: flagJSON, Usage: "JSON object for an inline keyboard"},
				},
				PreCall: requireValidInputChecklist,
			},
		},
	})
}

// requireValidInputChecklist ensures the InputChecklist object conforms to Bot API requirements
// (non-empty title up to 255 chars, and a non-empty array of valid tasks with id and text).
func requireValidInputChecklist(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	v, ok := params["checklist"]
	if !ok {
		return nil
	}
	obj, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("--checklist must be a JSON object")
	}
	title, ok := obj["title"].(string)
	if !ok || len([]rune(title)) < 1 || len([]rune(title)) > 255 {
		return fmt.Errorf("--checklist title must be between 1 and 255 characters")
	}
	tasksRaw, ok := obj["tasks"].([]any)
	if !ok || len(tasksRaw) < 1 {
		return fmt.Errorf("--checklist tasks must be a non-empty JSON array")
	}
	for i, t := range tasksRaw {
		taskMap, ok := t.(map[string]any)
		if !ok {
			return fmt.Errorf("--checklist task at index %d must be an object", i)
		}
		if _, hasID := taskMap["id"]; !hasID {
			return fmt.Errorf("--checklist task at index %d must have an id", i)
		}
		txt, ok := taskMap["text"].(string)
		if !ok || txt == "" {
			return fmt.Errorf("--checklist task at index %d must have non-empty text", i)
		}
	}
	return nil
}
