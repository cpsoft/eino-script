package loaders

import (
	"errors"
	"github.com/cloudwego/eino/components/document"
)

type LoaderCreateFunc func(data map[string]interface{}) (document.Loader, error)

var loaderRegistry = map[string]LoaderCreateFunc{
	"web":  CreateWebLoader,
	"file": CreateFileLoader,
}

func CreateLoaderNode(data map[string]interface{}) (document.Loader, error) {
	loaderType, ok := data["type"].(string)
	if !ok {
		return nil, errors.New("loaderType is not a string")
	}

	loaderFunc, ok := loaderRegistry[loaderType]
	if !ok {
		return nil, errors.New("loaderType is not a string")
	}

	return loaderFunc(data)
}
