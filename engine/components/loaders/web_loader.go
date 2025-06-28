package loaders

import (
	"context"
	"errors"
	"github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino/components/document"
)

type WebLoader struct {
	document.Loader
	uri string
}

func CreateWebLoader(data map[string]interface{}) (document.Loader, error) {
	var err error
	loader := &WebLoader{}

	loader.uri = data["uri"].(string)
	if len(loader.uri) <= 0 {
		return nil, errors.New("uri required")
	}

	loader.Loader, err = url.NewLoader(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return loader, nil
}
