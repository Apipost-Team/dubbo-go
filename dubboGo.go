package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Apipost-Team/dubbo-go/request"
	"golang.org/x/net/websocket"
)

func main() {
	var serverPort int
	flag.IntVar(&serverPort, "p", 10718, "server port， default：10718")
	flag.Parse()
	fmt.Printf("server port %d", serverPort)

	http.HandleFunc("/dubbo_send", func(w http.ResponseWriter, r *http.Request) {
		//增加代理发送功能
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))

		var dubboStruct request.DubboRequest

		err = json.Unmarshal(body, &dubboStruct)
		if err != nil {
			myHttpRespone(w, err)
			return
		}

		if dubboStruct.ApiName == "" {
			myHttpRespone(w, "api_name is empty")
			return
		}

		resp, err := dubboStruct.Send()
		if err != nil {
			myHttpRespone(w, err)
		} else {
			myHttpRespone(w, struct {
				TargetId string `json:"target_id"`
				Resp     any    `json:"resp"`
			}{TargetId: dubboStruct.TargetId, Resp: resp})
		}
	})

	http.Handle("/websocket", websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		var sendChan = make(chan string)
		go func(sendChan chan<- string, ws *websocket.Conn) {
			for {
				var body string
				if err := websocket.Message.Receive(ws, &body); err != nil {
					fmt.Println("read error")
					fmt.Println(err)
					ws.Close()
					break
				}
				var dubboStruct request.DubboRequest

				err := json.Unmarshal([]byte(body), &dubboStruct)
				if err != nil {
					msg := `{"code":501, "msg":` + errorJson(err) + `, "data":{}}`
					sendChan <- msg
					continue
				}

				if dubboStruct.ApiName == "" {
					msg := `{"code":502, "msg":"api_name is empty", "data":{"target_id":"` + dubboStruct.TargetId + `"}}`
					sendChan <- msg
					continue
				}
				resp, err := dubboStruct.Send()
				if err != nil {
					msg := `{"code":503, "msg":` + errorJson(err) + `, "data":{"target_id":"` + dubboStruct.TargetId + `"}}`
					sendChan <- msg
					continue
				}

				response, err := json.Marshal(struct {
					TargetId string `json:"target_id"`
					Resp     any    `json:"resp"`
				}{TargetId: dubboStruct.TargetId, Resp: resp})

				if err != nil {
					msg := `{"code":504, "msg":` + errorJson(err) + `, "data":{"target_id":"` + dubboStruct.TargetId + `"}}`
					sendChan <- msg
					continue
				}

				msg := `{"code":10000, "msg":"success", "data":` + string(response) + `}`
				sendChan <- msg
			}
		}(sendChan, ws)

		for {
			msg, ok := <-sendChan
			if !ok {
				break
			}
			if err := websocket.Message.Send(ws, msg); err != nil {
				fmt.Println("write")
				fmt.Println(err)
				break
			}
		}

	}))

	if err := http.ListenAndServe(":"+strconv.Itoa(serverPort), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func myHttpRespone(w http.ResponseWriter, r any) {
	err, ok := r.(error)
	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":1001, "msg":"` + err.Error() + `"}`))
		return
	}
	msg, ok := r.(string)
	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":1002, "msg":"` + msg + `"}`))
		return
	}

	response, err := json.Marshal(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":1003, "msg":"` + err.Error() + `"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"code":10000, "msg":"success", "data":` + string(response) + `}`))
}

func errorJson(err error) string {
	if err == nil {
		return "None"
	}

	msg := err.Error()
	if msg == "" {
		return "None"
	}

	jsonStr, err1 := json.Marshal(msg)
	if err1 != nil {
		return msg
	}

	return string(jsonStr)
}
