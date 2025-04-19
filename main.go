package main

import (
	"eino-script/engine"
	"flag"
	"fmt"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
	"regexp"
)

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

		err = e.Invoke(in)
		if err != nil {
			logrus.Error(err)
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
	filePath := flag.String("file", "", "file path")
	flag.Parse()

	if *filePath == "" {
		flag.PrintDefaults()
		logrus.Errorf("file path is required")
	}

	system, err := engine.InitSystem()
	if err != nil {
		logrus.Error(err)
	}

	defer system.Close()

	e, err := engine.CreateEngineByFile(*filePath)
	if err != nil {
		fmt.Println(err)
	}

	if e == nil {
		fmt.Println("engine is nil")
		return
	}

	defer e.Close()

	shell(e)

	return
}
