package server

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"

	"github.com/Project-Prismatica/prismatica_report_renderer"
)

func getCurrentTimestamp()(*timestamp.Timestamp) {
	currentTime := time.Now()
	return &timestamp.Timestamp{
		Seconds:currentTime.Unix(),
	}
}

func TestXpathSingleResultTemplate(t *testing.T) {
	ctx := context.Background()
	req := &prismatica_report_renderer.RenderRequest{}

	registrationResponse, err := registerTestTemplate(xpathSingleResultTemplate)
	assert.Nil(t, err)

	req.Timestamp = getCurrentTimestamp()
	req.TemplateId = registrationResponse.TemplateId

	res, err := cli.Render(ctx, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Result,"z contents")
}

func TestXpathMultiResultTemplate(t *testing.T) {
	ctx := context.Background()
	req := &prismatica_report_renderer.RenderRequest{}

	registrationResponse, err := registerTestTemplate(xpathMultiResultTemplate)
	assert.Nil(t, err)

	req.Timestamp = getCurrentTimestamp()
	req.TemplateId = registrationResponse.TemplateId

	res, err := cli.Render(ctx, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Result,"a contents")
	assert.Contains(t, res.Result,"b contents")
}

func TestRenderKnownTemplate(t *testing.T) {
	ctx := context.Background()
	req := &prismatica_report_renderer.RenderRequest{}

	registrationResponse, err := registerTestTemplate(noParameterRenderTemplate)
	assert.Nil(t, err)

	req.Timestamp = getCurrentTimestamp()
	req.TemplateId = registrationResponse.TemplateId

	res, err := cli.Render(ctx, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Result,"42")
}

func TestRenderUnknownTemplateId(t *testing.T) {
	ctx := context.Background()
	req := &prismatica_report_renderer.RenderRequest{}

	req.Timestamp = getCurrentTimestamp()
	req.TemplateId = "asdf"

	res, err := cli.Render(ctx, req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unknown template Id")
}

func TestRenderNoTimestamp(t *testing.T) {
	ctx := context.Background()
	req := &prismatica_report_renderer.RenderRequest{}

	res, err := cli.Render(ctx, req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no timestamp")
}
