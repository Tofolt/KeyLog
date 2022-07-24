package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	//TODO add man page output
	var (
		defaultPath      = ".\\"
		pathToStoreFlag  = flag.String("path", defaultPath, "Path to store log files. Default is current directory")
		portToListenFlag = flag.Int("port", 1337, "Listen Port")
		tokenFlag        = flag.String("token", "", "Bot token to send files to Telegram")
	)

	flag.Parse()

	if *tokenFlag != "" {
		BotInit(*tokenFlag)
	}

	Listen(*tokenFlag, *pathToStoreFlag, *portToListenFlag)
}

func MakeConfigFile() {

	type Config struct {
		TelegramChatID int64 `json:"telegram_chat_ID"`
	}

	config := Config{TelegramChatID: chatID}
	configByte, _ := json.Marshal(config)
	err := os.WriteFile("./config.json", configByte, 0644)

	if err != nil {
		return
	}
}

func ReceiveBody(r *http.Request) []byte {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	return buf
}

func WriteFileToRoot(r *http.Request, path string) {
	filePath := CreateFilePath(path)
	bornFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(bornFile, string(ReceiveBody(r)))
	defer bornFile.Close()
	log.Println("DATA RECEIVED")
	log.Println("FILE WRITTEN", "[", filePath, "]")

}

func CreateFilePath(path string) string {
	t := time.Now()
	filePath := path + t.Format("2006-01-02 15-04") + ".txt"
	return filePath
}

func Listen(tokenFlag string, pathToStoreFlag string, portToListenFlag int) {
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		WriteFileToRoot(r, pathToStoreFlag)
		if tokenFlag != "" {
			TelegramBotSend(CreateFilePath(pathToStoreFlag))
		}

	})
	err := http.ListenAndServe(":"+strconv.Itoa(portToListenFlag), nil)
	if err != nil {
		return
	}
}
