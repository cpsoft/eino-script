package main

import (
	"eino-script/engine"
	"eino-script/server"
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

		_, err = e.Invoke(in)
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
	isServer := flag.Bool("server", false, "start server")
	callback := flag.String("callback", "", "callback mode")
	flag.Parse()

	if *callback != "" {
		system, err := engine.InitSystem()
		if err != nil {
			logrus.Error(err)
		}
		defer system.Close()
	}

	if *isServer {
		server.StartServer()
	} else if *filePath != "" {
		e, err := engine.CreateEngineByFile(nil, *filePath, "flowgram")
		if err != nil {
			logrus.Error(err)
			return
		}

		if e == nil {
			logrus.Error("engine is nil")
			return
		}

		defer e.Close()

		shell(e)
	} else {
		flag.PrintDefaults()
	}

	return
}
