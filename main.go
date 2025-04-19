package main

import (
	"context"
	"eino-script/engine"
	"eino-script/parser"
	"flag"
	"fmt"
	"github.com/cloudwego/eino-ext/callbacks/apmplus"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"regexp"
)

func readScript() ([]byte, error) {
	filePath := flag.String("file", "", "file path")
	flag.Parse()

	if *filePath == "" {
		flag.PrintDefaults()
		return nil, fmt.Errorf("file path is required")
	}

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("open file error:", err)
		return nil, fmt.Errorf("open file error, %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file error, %s", err)
	}
	return data, nil
}

func shell(e *engine.Engine) {
	var msg string
	var assistantMessage string

	var chat_history = make([]*schema.Message, 0)
	for {
		fmt.Print(">>> ")
		_, err := fmt.Scanln(&msg)
		if err != nil {
			break
		}
		if msg == "exit" {
			return
		}

		in := map[string]any{"message": msg}
		if assistantMessage != "" {
			re := regexp.MustCompile(`(?s)<think>.*?</think>`)
			assistantMessage = re.ReplaceAllString(assistantMessage, "")
			chat_history = append(chat_history, schema.AssistantMessage(assistantMessage, nil))
			assistantMessage = ""
			in["chat_history"] = chat_history
		}

		//fmt.Printf("%v\n", in)

		err = e.Invoke(in)
		if err != nil {
			return
		}

		//err = e.Stream(in)
		//if err != nil {
		//	fmt.Println("stream error:", err)
		//	return
		//}
		//
		//chat_history = append(chat_history, schema.UserMessage(msg))
		//defer e.Close()
		//
		//for {
		//	chunks, err := e.Recv()
		//	if err != nil {
		//		if err == io.EOF {
		//			fmt.Println()
		//			break
		//		}
		//		fmt.Println(err)
		//		return
		//
		//	}
		//	//for _, chunk := range chunks {
		//	//	assistantMessage += chunk.Content
		//	//	fmt.Printf(chunk.Content)
		//	//}
		//	assistantMessage += chunks.Content
		//	fmt.Printf(chunks.Content)
		//}
	}
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	data, err := readScript()
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg, err := parser.Parser(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	cbh, showdown, err := apmplus.NewApmplusHandler(&apmplus.Config{
		Host:        "apmplus-cn-beijing.volces.com:4317",
		AppKey:      "fad265ce90e3613e299af8ef6f4af03c",
		ServiceName: "eino-app",
		Release:     "release/v0.0.1",
	})
	if err != nil {
		fmt.Println(err)
	}

	callbacks.AppendGlobalHandlers(cbh)

	e, err := engine.CreateEngine(cfg)
	if err != nil {
		fmt.Println(err)
	}

	if e == nil {
		fmt.Println("engine is nil")
		return
	}

	defer e.Close()

	shell(e)

	showdown(context.Background())
	return
}
