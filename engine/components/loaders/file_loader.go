package loaders

import (
	"context"
	"errors"
	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/parser/html"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/document/parser"
)

type FileLoader struct {
	document.Loader
	uri string
}

func CreateExtParser() (parser.Parser, error) {
	ctx := context.Background()
	textParser := parser.TextParser{}

	htmlParser, _ := html.NewParser(ctx, &html.Config{})

	pdfParser, _ := pdf.NewPDFParser(ctx, &pdf.Config{})

	extParser, _ := parser.NewExtParser(ctx, &parser.ExtParserConfig{
		Parsers: map[string]parser.Parser{
			".html": htmlParser,
			".pdf":  pdfParser,
		},
		FallbackParser: textParser,
	})
	return extParser, nil
}

func CreateFileLoader(data map[string]interface{}) (document.Loader, error) {
	var err error
	loader := &FileLoader{}

	loader.uri = data["uri"].(string)
	if len(loader.uri) <= 0 {
		return nil, errors.New("uri required")
	}

	p, err := CreateExtParser()
	if err != nil {
		return nil, err
	}

	loader.Loader, err = file.NewFileLoader(context.Background(), &file.FileLoaderConfig{
		UseNameAsID: true,
		Parser:      p,
	})
	if err != nil {
		return nil, err
	}

	return loader, nil
}
