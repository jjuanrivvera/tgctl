# tgctl demo — rendered with VHS (https://github.com/charmbracelet/vhs).
#
# NO ACCOUNT, NO SIGNUP, NO NETWORK. The CLI talks to tools/demoserver, a local
# stand-in that answers with invented data, so what you see is the real thing:
# real tables, real JSON, real --jq, real exit paths. Only the answers are made up.
#
# Regenerate:
#   make build && go build -o bin/demoserver ./tools/demoserver && vhs .demo/demo.tape

Output assets/demo.gif

Set Shell bash
Set FontSize 15
Set Width 1200
Set Height 760
Set Padding 18
Set Theme "Dracula"
Set TypingSpeed 28ms

Hide
Type "export PATH=$PWD/bin:$PATH TGCTL_TOKEN=123456789:AAFakeTokenForTheDemoOnly"
Enter
Type "demoserver :8642 >/dev/null 2>&1 & sleep 1"
Enter
Type "alias tgctl='tgctl --base-url http://localhost:8642'; clear"
Enter
Show

Sleep 500ms
Type "tgctl bot info"
Sleep 250ms
Enter
Sleep 2000ms

Type "tgctl message send --chat @demo_channel --text 'hello from tgctl'"
Sleep 250ms
Enter
Sleep 2200ms

Type "tgctl updates get --limit 2"
Sleep 250ms
Enter
Sleep 2400ms

Type "clear"
Enter
Sleep 300ms

Type "tgctl updates get --limit 2 -o json --jq '.[].message.text'"
Sleep 250ms
Enter
Sleep 2400ms

Type "# and the request behind any command, without sending it"
Sleep 200ms
Enter
Type "tgctl message send --chat @demo_channel --text hi --dry-run"
Sleep 250ms
Enter
Sleep 2400ms
