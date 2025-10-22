package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ToggleSubtreeExpansion = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle subtree expansion recursively with Alt+Enter",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("level1")
		shell.CreateDir("level1/level2")
		shell.CreateDir("level1/level2/level3")
		shell.CreateFile("level1/level2/level3/file-deep", "deep content\n")
		shell.CreateFile("level1/file-mid", "mid content\n")
		shell.CreateDir("other")
		shell.CreateFile("other/file-other", "other content\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ level1"),
				Equals("    ▼ level2/level3"),
				Equals("      ?? file-deep"),
				Equals("    ?? file-mid"),
				Equals("  ▼ other"),
				Equals("    ?? file-other"),
			)

		// Navigate to level1 directory and press Alt+Enter to collapse recursively
		t.Views().Files().
			SelectNextItem().
			Press(keys.Files.ToggleSubtreeExpansion).
			Lines(
				Equals("▼ /"),
				Equals("  ▶ level1").IsSelected(),
				Equals("  ▼ other"),
				Equals("    ?? file-other"),
			)

		// Press Alt+Enter again on level1 to expand recursively
		t.Views().Files().
			Press(keys.Files.ToggleSubtreeExpansion).
			Lines(
				Equals("▼ /"),
				Equals("  ▼ level1").IsSelected(),
				Equals("    ▼ level2/level3"),
				Equals("      ?? file-deep"),
				Equals("    ?? file-mid"),
				Equals("  ▼ other"),
				Equals("    ?? file-other"),
			)
	},
})
