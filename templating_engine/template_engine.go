package templating_engine

import (
	"errors"
	"github.com/flosch/pongo2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/Project-Prismatica/prismatica_report_renderer/template_functions"
	_ "github.com/flosch/pongo2-addons"
)

type TemplatingEngine struct {
}

type ReportTemplate struct {
	Id, Template         string
	RenderEngineTemplate *pongo2.Template
}

type RenderTemplateProvider interface {
	ResolveTemplate(templateId string)(foundTemplate *ReportTemplate, err error)
	StoreTemplate(toStore *ReportTemplate)(err error)
	Shutdown()
	SupportsWrite()(bool)
}

type RenderTemplateProviderFactory func (*viper.Viper)(RenderTemplateProvider,
	error)

var (
	RegisteredTemplateProviderFactories = make(
		map[string]RenderTemplateProviderFactory)
)

func NewTemplateEngine() (engine *TemplatingEngine){
	engine = new(TemplatingEngine)

	return
}

func RegisterTemplateProvider(name string,
		provider RenderTemplateProviderFactory) {
	if oldProvider, inMap := RegisteredTemplateProviderFactories[name]; inMap {
		logrus.WithFields(logrus.Fields{
			"providerName": name,
			"oldProvider": oldProvider,
			"newProvider": provider,
		}).Warn("overwriting registered template provider")
	}
	RegisteredTemplateProviderFactories[name] = provider
}

func NewTemplate(rawTemplate string) (tpl *ReportTemplate,
		err error) {

	compiledTemplate, contextCreationError := pongo2.FromString(rawTemplate)
	if contextCreationError != nil {
		logrus.WithFields(logrus.Fields{"error": contextCreationError}).
			Error("could not create render context")
		err = contextCreationError
		return
	}

	tpl = new(ReportTemplate)
	tpl.RenderEngineTemplate = compiledTemplate
	return
}

func (s TemplatingEngine) Render (
		inputTemplate *ReportTemplate,
		contextVariables map[string]string) (result string, err error) {
	templateContext := pongo2.Context{}
	for k, v := range contextVariables {
		templateContext[k] = v
	}

	renderedString, renderingError := inputTemplate.RenderEngineTemplate.
		Execute(templateContext)
	if renderingError != nil {
		logrus.WithFields(logrus.Fields{"templateId": inputTemplate.Id,
			"error": renderingError}).
			Warn("could not render template")
		err = errors.New("rendering engine error while templating")
		return
	}

	result = renderedString
	return
}
