package loaders

import (
	"context"
	"eino-script/engine/types"
	"errors"
	"github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino/components/document"
)

type WebLoader struct {
	uri string
}

func (l *WebLoader) GetEinoNode() (document.Loader, error) {
	ctx := context.Background()
	loader, err := url.NewLoader(ctx, nil)
	if err != nil {
		return nil, err
	}

	//docs, err := loader.Load(ctx, document.Source{
	//	URI: l.uri,
	//})
	return loader, err
}

func CreateWebLoader(data map[string]interface{}) (types.LoaderInterface, error) {
	loader := &WebLoader{}

	loader.uri = data["uri"].(string)
	if len(loader.uri) <= 0 {
		return nil, errors.New("uri required")
	}

	return loader, nil
}
