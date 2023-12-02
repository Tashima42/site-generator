package generator

import (
	"embed"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	templatesFolder  = "tmp/templates"
	pagesFolder      = "tmp/pages"
	pageTemplateFile = templatesFolder + "/template.html"
)

type Options struct {
	Name            string
	DestinationPath string
}

type PageData struct {
	HTMLTitle string
	PageTitle string
	MenuItems []MenuItem
}

type MetaData struct {
	PageTitle string `yaml:"title"`
	FileName  string `yaml:"fileName"`
	HTMLPath  string `yaml:"htmlPath"`
	Path      string
}

type MenuItem struct {
	HTMLPageName string
	PageName     string
	Current      bool
}

//go:embed tmp
var fs embed.FS

func Generate(opt Options) error {
	source, err := fs.ReadDir("tmp/pages")
	if err != nil {
		return err
	}

	var pagesMetaData []MetaData
	var menuItens []MenuItem

	for _, e := range source {
		if !e.IsDir() {
			return fmt.Errorf("%s is not a directory, review your pages folder", e.Name())
		}
		meta := MetaData{
			FileName: e.Name(),
			Path:     pagesFolder + "/" + e.Name(),
		}
		metaFile, err := fs.ReadFile(meta.Path + "/" + "meta.yaml")
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(metaFile, &meta); err != nil {
			return err
		}
		pagesMetaData = append(pagesMetaData, meta)
		menuItens = append(menuItens, MenuItem{
			HTMLPageName: meta.HTMLPath,
			PageName:     meta.PageTitle,
			Current:      false,
		})
	}

	if err := cleanFolder(opt.DestinationPath); err != nil {
		return err
	}

	for i, pm := range pagesMetaData {
		tpl, err := template.New("wrapper").ParseFS(fs, pageTemplateFile, pm.Path+"/page.html")
		if err != nil {
			return err
		}
		f, err := os.Create(opt.DestinationPath + "/" + pm.HTMLPath)
		if err != nil {
			return err
		}
		defer f.Close()
		menuItens[i].Current = true
		pageData := PageData{
			HTMLTitle: opt.Name,
			PageTitle: opt.Name,
			MenuItems: menuItens,
		}
		if err := tpl.Execute(f, pageData); err != nil {
			return err
		}
		menuItens[i].Current = false
	}
	return nil
}

func cleanFolder(folder string) error {
	if err := os.RemoveAll(folder); err != nil {
		return err
	}
	return os.MkdirAll(folder, 0700)
}
