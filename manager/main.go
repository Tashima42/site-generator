package manager

import (
	"fmt"
	"os"
)

type Options struct {
	PagesPath  string
	FolderName string
	Title      string
	HTMLPath   string
}

func AddPage(opt Options) error {
	folderPath := opt.PagesPath + "/" + opt.FolderName

	if err := os.MkdirAll(folderPath, 0700); err != nil {
		return err
	}

	pageContent := `{{define "print-style"}}
{{end}}
{{define "style"}}
{{end}}
{{define "meta"}}
{{end}}
{{define "content"}}
{{end}}`
	metaContent := fmt.Sprintf(`title: "%s"
htmlPath: "%s"`, opt.Title, opt.HTMLPath)
	if err := os.WriteFile(folderPath+"/page.html", []byte(pageContent), 0700); err != nil {
		return err
	}
	return os.WriteFile(folderPath+"/meta.yaml", []byte(metaContent), 0700)
}
