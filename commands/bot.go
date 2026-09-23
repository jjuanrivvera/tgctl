package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	registerGroup(group{
		Use:   "bot",
		Short: "Inspect and configure the bot itself",
		Long:  "Read the bot's identity (getMe) and manage its name/description shown in Telegram.",
		Cmds: []methodCmd{
			{
				Use: "info", Method: "getMe", Kind: kindRead,
				Short:   "Show the authenticated bot's identity (getMe)",
				Example: "  tgctl bot info\n  tgctl bot info -o json",
				Columns: []string{"id", "username", "first_name", "can_join_groups"},
			},
			{
				Use: "set-name", Method: "setMyName", Kind: kindWrite,
				Short:   "Set the bot's name",
				Example: `  tgctl bot set-name --name "My Helper Bot"`,
				Flags: []flagSpec{
					{Name: "name", Required: true, Usage: "new bot name (0-64 chars)"},
					{Name: "language-code", Param: "language_code", Usage: "BCP-47 code this name applies to"},
				},
			},
			{
				Use: "get-name", Method: "getMyName", Kind: kindRead,
				Short:   "Get the bot's name",
				Example: "  tgctl bot get-name\n  tgctl bot get-name --language-code es",
				Flags:   []flagSpec{{Name: "language-code", Param: "language_code", Usage: "language to query"}},
			},
			{
				Use: "set-description", Method: "setMyDescription", Kind: kindWrite,
				Short:   "Set the bot's description (shown in the empty chat)",
				Example: `  tgctl bot set-description --description "I help you manage your groups."`,
				Flags: []flagSpec{
					{Name: "description", Required: true, Usage: "new description (0-512 chars)"},
					{Name: "language-code", Param: "language_code", Usage: "language this description applies to"},
				},
			},
			{
				Use: "get-description", Method: "getMyDescription", Kind: kindRead,
				Short:   "Get the bot's description",
				Example: "  tgctl bot get-description\n  tgctl bot get-description --language-code es",
				Flags:   []flagSpec{{Name: "language-code", Param: "language_code", Usage: "language to query"}},
			},
			{
				Use: "set-short-description", Method: "setMyShortDescription", Kind: kindWrite,
				Short:   "Set the bot's short description (shown on the profile page)",
				Example: `  tgctl bot set-short-description --short-description "Group management, done right."`,
				Flags: []flagSpec{
					{Name: "short-description", Param: "short_description", Usage: "new short description (0-120 chars; empty clears it)"},
					{Name: "language-code", Param: "language_code", Usage: "language this description applies to"},
				},
			},
			{
				Use: "get-short-description", Method: "getMyShortDescription", Kind: kindRead,
				Short:   "Get the bot's short description",
				Example: "  tgctl bot get-short-description\n  tgctl bot get-short-description --language-code es",
				Flags:   []flagSpec{{Name: "language-code", Param: "language_code", Usage: "language to query"}},
			},
			{
				Use: "set-admin-rights", Method: "setMyDefaultAdministratorRights", Kind: kindWrite,
				Short:   "Set the bot's default administrator rights (requested when added to a group/channel)",
				Example: `  tgctl bot set-admin-rights --rights '{"can_manage_chat":true,"can_delete_messages":true}'`,
				Flags: []flagSpec{
					{Name: "rights", Kind: flagJSON, Usage: "ChatAdministratorRights object as JSON (omit to clear)"},
					{Name: "for-channels", Param: "for_channels", Kind: flagBool, Usage: "apply to channels instead of groups/supergroups"},
				},
			},
			{
				Use: "get-admin-rights", Method: "getMyDefaultAdministratorRights", Kind: kindRead,
				Short:   "Get the bot's default administrator rights",
				Example: "  tgctl bot get-admin-rights\n  tgctl bot get-admin-rights --for-channels -o json",
				Flags: []flagSpec{
					{Name: "for-channels", Param: "for_channels", Kind: flagBool, Usage: "query the channel rights instead of group rights"},
				},
			},
			{
				Use: "set-photo", Method: "setMyProfilePhoto", Kind: kindWrite,
				Short: "Set the bot's profile photo",
				Long: `Set the bot's own profile photo (setMyProfilePhoto).

The photo must be a local file: Telegram does not accept a URL or an existing file_id
here, because a profile photo cannot be reused. Pass --photo for a static .JPG, or
--animated for an MPEG4 animation (with --main-frame-timestamp to choose the frame
Telegram shows as the still image).`,
				Example: `  tgctl bot set-photo --photo logo.jpg
  tgctl bot set-photo --animated intro.mp4 --main-frame-timestamp 1.5`,
				Files: []fileSpec{
					{Name: "photo", Usage: "local .JPG to use as the static profile photo"},
					{Name: "animated", Usage: "local MPEG4 to use as an animated profile photo"},
				},
				Flags: []flagSpec{
					{Name: "main-frame-timestamp", Param: "main_frame_timestamp", Kind: flagFloat,
						Usage: "seconds into --animated for the still frame (default 0.0)"},
				},
				PreCall: attachProfilePhoto,
			},
			{
				Use: "remove-photo", Method: "removeMyProfilePhoto", Kind: kindDestructive,
				Short:   "Remove the bot's profile photo",
				Long:    "Remove the bot's own profile photo (removeMyProfilePhoto). It cannot be restored except by uploading it again.",
				Example: `  tgctl bot remove-photo`,
			},
			{
				Use: "close", Method: "close", Kind: kindWrite,
				Short:   "Close the bot instance before moving it to another server",
				Long:    "Close the bot instance (frees server resources). Returns an error for the first 10 minutes after the bot launches.",
				Example: `  tgctl bot close`,
			},
			{
				Use: "logout", Method: "logOut", Kind: kindDestructive,
				Short:   "Log out from the cloud Bot API before running a local Bot API server",
				Long:    "Log the bot out of the cloud Bot API. After this you can use a local Bot API server; you must re-login via api.telegram.org to switch back. Returns an error for the first 10 minutes after launch.",
				Example: `  tgctl bot logout`,
			},
		},
		Extra: []func() *cobra.Command{newBotPhotoCmd},
	})
}

// attachProfilePhoto builds the InputProfilePhoto object setMyProfilePhoto expects. The method
// takes no plain file field: it takes a JSON object that POINTS at the uploaded bytes with
// "attach://<part>", so the part and the param that names it have to be built together — the
// one place in the fleet where a flag does not map 1:1 to a parameter (Bot API 10.3,
// InputProfilePhotoStatic / InputProfilePhotoAnimated).
func attachProfilePhoto(_ *cobra.Command, params map[string]any, files map[string]string) error {
	// The part name is deliberately not "photo": that is the parameter's own name, and two
	// parts with the same name is not a request Telegram can read.
	const part = "profile_photo"

	// collectFiles passes a value it could not open as a local file straight through as a
	// string param — the right default for sendPhoto (a URL or a file_id both work there) and
	// exactly wrong here, since a profile photo can only ever be a fresh upload.
	for _, name := range []string{"photo", "animated"} {
		if v, ok := params[name]; ok {
			return fmt.Errorf("--%s %v: a profile photo must be an uploadable local file — Telegram accepts no URL or file_id here", name, v)
		}
	}

	static, hasStatic := files["photo"]
	animated, hasAnimated := files["animated"]
	switch {
	case hasStatic && hasAnimated:
		return fmt.Errorf("pass either --photo or --animated, not both")
	case hasStatic:
		if _, ok := params["main_frame_timestamp"]; ok {
			return fmt.Errorf("--main-frame-timestamp applies to --animated, not to a static --photo")
		}
		delete(files, "photo")
		files[part] = static
		params["photo"] = map[string]any{"type": "static", "photo": "attach://" + part}
	case hasAnimated:
		delete(files, "animated")
		files[part] = animated
		photo := map[string]any{"type": "animated", "animation": "attach://" + part}
		if ts, ok := params["main_frame_timestamp"]; ok {
			photo["main_frame_timestamp"] = ts
			delete(params, "main_frame_timestamp")
		}
		params["photo"] = photo
	default:
		return fmt.Errorf("one of --photo <file.jpg> or --animated <file.mp4> is required")
	}
	return nil
}

// newBotPhotoCmd reports whether the bot currently has a profile photo. getUserProfilePhotos
// needs a user id and the bot's own is not something you type, so this fills it in (issue #23).
//
// The id comes from the token's non-secret prefix rather than from getMe: it costs no request,
// and it is the only source that works under --dry-run, where a call returns nothing and a
// getMe would hand back a zero-valued User whose empty id would print an invalid curl.
func newBotPhotoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "photo",
		Short:   "Show the bot's current profile photo(s)",
		Long:    "Look up the bot's own profile photos (getUserProfilePhotos for the bot's own id).",
		Example: "  tgctl bot photo\n  tgctl bot photo -o json",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := clientFromCmd(cmd)
			if err != nil {
				return err
			}
			defer func() { _ = client.Close() }()
			userID := client.BotID()
			if userID == "" {
				// Not a bot token (a custom authenticator): fall back to asking the API,
				// which is correct everywhere except a dry run, where there is no answer.
				me, err := client.GetMe(cmd.Context())
				if err != nil {
					return err
				}
				if userID = me.ID.String(); userID == "" {
					return fmt.Errorf("cannot determine the bot's own id")
				}
			}
			raw, err := client.Call(cmd.Context(), "getUserProfilePhotos",
				map[string]any{"user_id": userID, "limit": 1}, true)
			if err != nil {
				return err
			}
			if !cmd.Flags().Changed("columns") {
				if err := cmd.Flags().Set("columns", "total_count"); err != nil {
					return err
				}
			}
			return render(cmd, raw)
		},
	}
	markKind(cmd, kindRead)
	return cmd
}
