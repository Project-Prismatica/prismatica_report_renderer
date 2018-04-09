package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"

	"github.com/Project-Prismatica/prismatica_report_renderer"
)

var (
	noParameterRenderTemplate = prismatica_report_renderer.ReportTemplate{
		Id: "noParameterRenderTemplate",
		Template: `the answer to life, the universe, and everything is {{ 41 + 1}}`,
	}

	xpathSingleResultTemplate = prismatica_report_renderer.ReportTemplate{
		Id: "xpathTemplate",
		Template: `
			{{ "<z>z contents</z><b>b contents</b>" | xpath:"/z" }}
		`,
	}

	xpathMultiResultTemplate = prismatica_report_renderer.ReportTemplate{
		Id: "xpathTemplate",
		Template: `
			{% for result in "<a>a contents</a><b>b contents</b>" | xpath:"/a|/b" %}
				{{ result }}
			{% endfor %}
		`,
	}
)

func registerTestTemplate(template prismatica_report_renderer.ReportTemplate) (
		*prismatica_report_renderer.TemplateRegistrationResponse, error){
	ctx := context.Background()
	req := &prismatica_report_renderer.TemplateRegistrationRequest{}

	req.Timestamp = getCurrentTimestamp()
	req.Template = template.Template

	return cli.RegisterTemplate(ctx, req)
}

func TestRegisterTemplate(t *testing.T) {
	res, err := registerTestTemplate(noParameterRenderTemplate)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
