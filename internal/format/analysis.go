package format

import (
	"bytes"
	"fmt"

	"github.com/mitchellh/colorstring"
	"github.com/padok-team/tfautomv/internal/tfautomv"
)

func Analysis(analysis *tfautomv.Analysis) string {
	var analysisBuf bytes.Buffer

	analysisBuf.WriteString(colorstring.Color("[bold][cyan]Analysis"))
	analysisBuf.WriteByte('\n')

	for _, createdResources := range analysis.CreatedByType {
		for _, created := range createdResources {

			// Display the resource planned for creation.

			analysisBuf.WriteByte('\n')
			analysisBuf.WriteString(colorstring.Color(fmt.Sprintf("[bold]%s[reset]", created.Address)))
			analysisBuf.WriteByte('\n')

			var resourceBuf bytes.Buffer

			// List all resources planned for destruction that matched.

			for _, comp := range analysis.Comparisons[created] {
				if comp.IsMatch() {
					resourceBuf.WriteString(colorstring.Color("[bold][green]Match: "))
					resourceBuf.WriteString(comp.Destroyed.Address)
					resourceBuf.WriteByte('\n')

					var diffBuf bytes.Buffer
					for _, attr := range comp.IgnoredAttributes {
						diffBuf.WriteString(colorstring.Color(fmt.Sprintf("[yellow]~ [reset]%s (some differences are ignored)", attr)))
						diffBuf.WriteByte('\n')
					}
					if diffBuf.Len() > 0 {
						resourceBuf.WriteString(WithLeftRule(&diffBuf, "green"))
					}
				}
			}

			// List all resources planned for destruction that mismatched and
			// include mismatching attributes.

			for _, comp := range analysis.Comparisons[created] {
				// Matching resources were already displayed above.
				if comp.IsMatch() {
					continue
				}

				resourceBuf.WriteString(colorstring.Color("[bold][red]Mismatch: "))
				resourceBuf.WriteString(comp.Destroyed.Address)
				resourceBuf.WriteByte('\n')

				var diffBuf bytes.Buffer
				for _, attr := range comp.MismatchingAttributes {
					diffBuf.WriteString(colorstring.Color(fmt.Sprintf("[green]+ [reset]%s = %#v", attr, created.Attributes[attr])))
					diffBuf.WriteByte('\n')
					diffBuf.WriteString(colorstring.Color(fmt.Sprintf("[red]- [reset]%s = %#v", attr, comp.Destroyed.Attributes[attr])))
					diffBuf.WriteByte('\n')
				}
				resourceBuf.WriteString(WithLeftRule(&diffBuf, "red"))
			}

			analysisBuf.WriteString(WithLeftRule(&resourceBuf, "white"))
		}
	}

	if len(analysis.CreatedByType) == 0 {
		analysisBuf.WriteString("\nNo resources are planned for creation.")
	}

	return WithLeftRule(&analysisBuf, "cyan")
}
