# tgctl demo — rendered with VHS (https://github.com/charmbracelet/vhs).
#
# NO ACCOUNT, NO NETWORK, NO SETUP. The token below is fake and never leaves the
# recording: every command runs with --dry-run, which prints the exact request the
# CLI would send and redacts the credential on its way out. That is real output —
# the request is what you would actually get — without a bot, a signup or a keyring.
#
# Regenerate:
#   make build && vhs .demo/demo.tape

Output assets/demo.gif

Set Shell bash
Set FontSize 15
Set Width 1200
Set Height 720
Set Padding 18
Set Theme "Dracula"
Set TypingSpeed 28ms

Hide
Type "export PATH=$PWD/bin:$PATH; export TGCTL_TOKEN=123456789:AAFakeTokenForTheDemoOnly; clear"
Enter
Show

Sleep 500ms
Type "tgctl --help | head -18"
Sleep 250ms
Enter
Sleep 2600ms

Type "clear"
Enter
Sleep 300ms

Type "# every command can show the request it would send"
Sleep 200ms
Enter
Type "tgctl message send --chat @channel --text 'hello from tgctl' --dry-run"
Sleep 250ms
Enter
Sleep 2200ms

Type "tgctl updates get --limit 5 -o json --dry-run"
Sleep 250ms
Enter
Sleep 2200ms

Type "tgctl message dice --chat @channel --emoji 🎲 --dry-run"
Sleep 250ms
Enter
Sleep 2400ms

Type "clear"
Enter
Sleep 300ms

Type "tgctl version"
Sleep 250ms
Enter
Sleep 1800ms
