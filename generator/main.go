package generator

import (
	"fmt"
	"html/template"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	templatesFolder  = "templates"
	pageTemplateFile = templatesFolder + "/template.html"
)

type Options struct {
	Name            string
	DestinationPath string
	SourcePath      string
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
}

func Generate(options Options) error {
	source, err := os.ReadDir(options.SourcePath)
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
			Path:     options.SourcePath + "/" + e.Name(),
		}
		metaFile, err := os.ReadFile(meta.Path + "/" + "meta.yaml")
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
		})
	}

	for _, pm := range pagesMetaData {
		tpl, err := template.New("wrapper").ParseFiles(pageTemplateFile, pm.Path+"/page.html")
		if err != nil {
			return err
		}
		f, err := os.Create(options.DestinationPath + "/" + pm.HTMLPath)
		if err != nil {
			return err
		}
		pageData := PageData{
			HTMLTitle: options.Name,
			PageTitle: options.Name,
			MenuItems: menuItens,
		}
		if err := tpl.Execute(f, pageData); err != nil {
			return err
		}
	}
	return nil
}
