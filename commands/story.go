package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "story",
		Short: "Manage and publish stories on behalf of business accounts",
		Long: `Post, edit, repost, and delete stories on behalf of connected business accounts.

Stories require the can_manage_stories business bot right.
New story content requires uploading a local photo or video (Telegram accepts no URLs or existing file_ids here).
Stories move to the archive after the specified active period (6h, 12h, 24h, or 48h).`,
		Cmds: []methodCmd{
			{
				Use: "post", Method: "postStory", Kind: kindWrite,
				Short: "Post a story on behalf of a business account",
				Long: `Post a story on behalf of a connected business account (postStory).
--content accepts an InputStoryContent JSON object pointing at uploaded files via attach://.
Alternatively, pass --photo or --video with local files. Returns the posted Story.`,
				Example: `  tgctl story post --business-connection-id bc123 --photo ./story.jpg --active-period 86400 --caption "Daily update"
  tgctl story post --business-connection-id bc123 --video ./story.mp4 --active-period 43200
  tgctl story post --business-connection-id bc123 --content '{"type":"photo","photo":"attach://pic"}' --file pic=@./story.jpg --active-period 86400`,
				Flags: []flagSpec{
					businessConnectionFlag(),
					{Name: "active-period", Param: "active_period", Kind: flagInt, Required: true, Usage: "period after which story moves to archive in seconds (21600, 43200, 86400, or 172800)"},
					{Name: "content", Kind: flagJSON, Usage: "JSON object describing story content (InputStoryContentPhoto or InputStoryContentVideo)"},
					{Name: "caption", Usage: "story caption (0-2048 chars)"},
					parseModeFlag(),
					{Name: "caption-entities", Param: "caption_entities", Kind: flagJSON, Usage: "JSON array of MessageEntity objects"},
					{Name: "areas", Kind: flagJSON, Usage: "JSON array of StoryArea clickable interactive overlay objects"},
					{Name: "post-to-chat-page", Param: "post_to_chat_page", Kind: flagBool, Usage: "keep story accessible on profile after expiration"},
					protectContentFlag(),
				},
				Files: []fileSpec{
					{Name: "photo", Usage: "local image file for a photo story"},
					{Name: "video", Usage: "local video file for a video story"},
					{Name: "cover", Usage: "optional local image file as video cover"},
				},
				AllowFiles: true,
				PreCall:    validateStoryPost,
				Columns:    []string{"id", "chat.id", "date"},
			},
			{
				Use: "edit", Method: "editStory", Kind: kindWrite,
				Short: "Edit a story posted on behalf of a business account",
				Long: `Edit a story previously posted on behalf of a connected business account (editStory).
--content accepts an InputStoryContent JSON object pointing at uploaded files via attach://.
Alternatively, pass --photo or --video with local files. Returns the edited Story.`,
				Example: `  tgctl story edit --business-connection-id bc123 --story-id 42 --caption "Updated caption"
  tgctl story edit --business-connection-id bc123 --story-id 42 --photo ./new.jpg`,
				Flags: []flagSpec{
					businessConnectionFlag(),
					storyIDFlag(),
					{Name: "content", Kind: flagJSON, Usage: "JSON object describing new story content"},
					{Name: "caption", Usage: "new story caption"},
					parseModeFlag(),
					{Name: "caption-entities", Param: "caption_entities", Kind: flagJSON, Usage: "JSON array of MessageEntity objects"},
					{Name: "areas", Kind: flagJSON, Usage: "new JSON array of StoryArea objects"},
				},
				Files: []fileSpec{
					{Name: "photo", Usage: "local image file for a photo story"},
					{Name: "video", Usage: "local video file for a video story"},
					{Name: "cover", Usage: "optional local image file as video cover"},
				},
				AllowFiles: true,
				PreCall:    validateStoryEdit,
				Columns:    []string{"id", "chat.id", "date"},
			},
			{
				Use: "repost", Method: "repostStory", Kind: kindWrite,
				Short: "Repost a story across business accounts",
				Long: `Repost a story across managed business accounts (repostStory).
Does not require uploading new media content. Returns the reposted Story.`,
				Example: `  tgctl story repost --business-connection-id dest_bc --from-business-connection-id src_bc --story-id 42 --active-period 86400`,
				Flags: []flagSpec{
					businessConnectionFlag(),
					{Name: "from-business-connection-id", Param: "from_business_connection_id", Required: true, Usage: "business connection id that originally posted the story"},
					storyIDFlag(),
					{Name: "active-period", Param: "active_period", Kind: flagInt, Required: true, Usage: "period after which story moves to archive in seconds (21600, 43200, 86400, or 172800)"},
					{Name: "caption", Usage: "optional custom caption for the reposted story"},
					parseModeFlag(),
					{Name: "caption-entities", Param: "caption_entities", Kind: flagJSON, Usage: "JSON array of MessageEntity objects"},
					{Name: "areas", Kind: flagJSON, Usage: "JSON array of StoryArea objects"},
					{Name: "post-to-chat-page", Param: "post_to_chat_page", Kind: flagBool, Usage: "keep story accessible on profile after expiration"},
					protectContentFlag(),
				},
				PreCall: validateStoryRepost,
				Columns: []string{"id", "chat.id", "date"},
			},
			{
				Use: "delete", Method: "deleteStory", Kind: kindDestructive,
				Short: "Delete a story posted on behalf of a business account",
				Long: `Delete a story previously posted on behalf of a connected business account (deleteStory).
Returns True on success.`,
				Example: `  tgctl story delete --business-connection-id bc123 --story-id 42`,
				Flags: []flagSpec{
					businessConnectionFlag(),
					storyIDFlag(),
				},
			},
		},
	})
}

func storyIDFlag() flagSpec {
	return flagSpec{
		Name:     "story-id",
		Param:    "story_id",
		Kind:     flagInt,
		Required: true,
		Usage:    "unique identifier of the story",
	}
}

func validateActivePeriod(v any) error {
	p, ok := v.(int64)
	if !ok {
		return fmt.Errorf("--active-period is required")
	}
	switch p {
	case 21600, 43200, 86400, 172800:
		return nil
	default:
		return fmt.Errorf("--active-period must be 21600 (6h), 43200 (12h), 86400 (24h), or 172800 (48h)")
	}
}

func validateStoryPost(_ *cobra.Command, params map[string]any, files map[string]string) error {
	if err := validateActivePeriod(params["active_period"]); err != nil {
		return err
	}
	return resolveStoryContent(params, files, true)
}

func validateStoryEdit(_ *cobra.Command, params map[string]any, files map[string]string) error {
	return resolveStoryContent(params, files, false)
}

func validateStoryRepost(_ *cobra.Command, params map[string]any, _ map[string]string) error {
	return validateActivePeriod(params["active_period"])
}

// resolveStoryContent resolves --photo/--video convenience flags or validates the --content JSON object.
// A story requires a fresh upload; URLs and Telegram file_ids are explicitly rejected.
func resolveStoryContent(params map[string]any, files map[string]string, required bool) error {
	// If collectFiles turned a value into a param (string), it was a URL or non-file.
	for _, f := range []string{"photo", "video", "cover"} {
		if v, ok := params[f]; ok {
			return fmt.Errorf("--%s %v: story content requires an uploadable local file — Telegram accepts no URL or file_id for stories", f, v)
		}
	}

	photoPath, hasPhoto := files["photo"]
	videoPath, hasVideo := files["video"]
	coverPath, hasCover := files["cover"]
	_, hasContent := params["content"]

	switch {
	case hasPhoto && hasVideo:
		return fmt.Errorf("pass either --photo or --video, not both")
	case hasContent && (hasPhoto || hasVideo):
		return fmt.Errorf("pass either --content or --photo/--video, not both")
	case hasPhoto:
		if hasCover {
			return fmt.Errorf("--cover applies to --video, not to --photo")
		}
		delete(files, "photo")
		files["story_content"] = photoPath
		params["content"] = map[string]any{
			"type":  "photo",
			"photo": "attach://story_content",
		}
	case hasVideo:
		delete(files, "video")
		files["story_content"] = videoPath
		contentMap := map[string]any{
			"type":  "video",
			"video": "attach://story_content",
		}
		if hasCover {
			delete(files, "cover")
			files["story_cover"] = coverPath
			contentMap["cover"] = "attach://story_cover"
		}
		params["content"] = contentMap
	case hasContent:
		v := params["content"]
		obj, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("--content must be a JSON object")
		}
		// Validate that photo or video references are attach:// references, not URLs or file_ids.
		for _, key := range []string{"photo", "video", "cover"} {
			if val, ok := obj[key].(string); ok {
				if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") || !strings.HasPrefix(val, "attach://") {
					return fmt.Errorf("story content %s %q: URLs and file_ids are not supported for stories (must be a fresh upload via attach://)", key, val)
				}
			}
		}
		for _, ref := range findAttachReferences(obj) {
			if _, hasFile := files[ref]; !hasFile {
				return fmt.Errorf("missing upload for attach reference %q (pass --file %s=@<path>)", ref, ref)
			}
		}
	default:
		if required {
			return fmt.Errorf("story content is required: pass --photo, --video, or --content")
		}
	}
	return nil
}
