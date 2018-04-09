package prismatica_report_renderer

import (

	"github.com/flosch/pongo2"
)

type RenderTemplateProvider interface {
	ResolveTemplate(templateId string)(foundTemplate *ReportTemplate, err error)
	StoreTemplate(toStore *ReportTemplate)(err error)
	SupportsWrite()(bool)
}

type ReportTemplate struct {
	Id, Template         string
	RenderEngineTemplate *pongo2.Template
}

