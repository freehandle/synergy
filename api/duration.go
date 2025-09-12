package api

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func mdToHTML(md []byte) string {
	if md == nil {
		log.Print("PANIC BUG: mdToHTML called with nil md ")
		return ""
	}
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	bytes := markdown.Render(doc, renderer)
	if len(bytes) > 0 {
		return string(bytes)
	}
	return ""
}

func PrettyDate(date time.Time) string {
	return date.Format("02 Jan 06")
}

func PrettyDuration(d time.Duration) string {
	if d.Seconds() < 60 {
		if d.Seconds() < 2 {
			return "1 segundo"
		}
		return fmt.Sprintf("%.0f segundos", d.Seconds())
	}
	if d.Minutes() < 60 {
		if d.Minutes() < 2 {
			return "1 minuto"
		}
		return fmt.Sprintf("%.0f minutos", d.Minutes())
	}
	if d.Hours() < 24 {
		if d.Hours() < 2 {
			return "1 hora"
		}
		return fmt.Sprintf("%.0f horas", d.Hours())
	}
	if d.Hours()/24.0 < 2 {
		return "1 dia"
	}
	return fmt.Sprintf("%.0f dias", d.Hours()/24.0)
}

func FileType(fileName string) string {
	pos := strings.LastIndex(fileName, ".")
	if pos < len(fileName) && pos > 0 {
		return fileName[pos+1:]
	}
	return ""
}
