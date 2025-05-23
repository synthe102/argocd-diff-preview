package diff

import (
	"fmt"
	"strings"

	"github.com/dag-andersen/argocd-diff-preview/pkg/selector"
	"github.com/dag-andersen/argocd-diff-preview/pkg/utils"
)

// Constants for markdown template
const noAppsFoundTemplate = `
## %title%

%message%
`

// WriteNoAppsFoundMessage writes a message to the output folder when no applications are found
func WriteNoAppsFoundMessage(
	title string,
	outputFolder string,
	selectors []selector.Selector,
	changedFiles []string,
) error {
	message := getNoAppsFoundMessage(selectors, changedFiles)
	markdown := generateNoAppsFoundMarkdown(title, message)
	markdownPath := fmt.Sprintf("%s/diff.md", outputFolder)

	if err := utils.WriteFile(markdownPath, markdown); err != nil {
		return fmt.Errorf("failed to write no apps found message: %w", err)
	}

	return nil
}

// generateNoAppsFoundMarkdown generates markdown from the message
func generateNoAppsFoundMarkdown(title, message string) string {
	markdown := strings.ReplaceAll(noAppsFoundTemplate, "%title%", title)
	markdown = strings.ReplaceAll(markdown, "%message%", message)
	return strings.TrimSpace(markdown)
}

// getNoAppsFoundMessage generates an appropriate message based on selectors and changed files
func getNoAppsFoundMessage(selectors []selector.Selector, changedFiles []string) string {
	selectorString := func(s []selector.Selector) string {
		var strs []string
		for _, selector := range s {
			strs = append(strs, selector.String())
		}
		return strings.Join(strs, ",")
	}

	switch {
	case len(selectors) > 0 && len(changedFiles) > 0:
		return fmt.Sprintf(
			"Found no changed Applications that matched `%s` and watched these files: `%s`",
			selectorString(selectors),
			strings.Join(changedFiles, "`, `"),
		)
	case len(selectors) > 0:
		return fmt.Sprintf(
			"Found no changed Applications that matched `%s`",
			selectorString(selectors),
		)
	case len(changedFiles) > 0:
		return fmt.Sprintf(
			"Found no changed Applications that watched these files: `%s`",
			strings.Join(changedFiles, "`, `"),
		)
	default:
		return "Found no Applications"
	}
}
