// Copyright 2016 rrframework Author @songtianyi. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logs

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

const (
	DEFAULT_WSCHANNEL_SIZE = 10000
)

// WSWriter implements beego LoggerInterface and is used to send logs to websocket clients
type WSWriter struct {
	Level       int `json:"level"`
	ChannelSize int `json:"channelSize"`
	msgs        chan string
}

// newWSWriter create websocket writer.
func newWSWriter() Logger {
	return &WSWriter{Level: LevelTrace}
}

// Init WSWriter with json config string
func (s *WSWriter) Init(jsonConfig string) error {
	if len(jsonConfig) != 0 {
		if err := json.Unmarshal([]byte(jsonConfig), s); err != nil {
			return err
		}
		s.msgs = make(chan string, s.ChannelSize)
	} else {
		s.msgs = make(chan string, DEFAULT_WSCHANNEL_SIZE)
	}
	http.Handle("/wslogs", websocket.Handler(s.wshandler()))
	return nil
}

func (s *WSWriter) wshandler() websocket.Handler {
	return func(ws *websocket.Conn) {
		for {
			select {
			case msg := <-s.msgs:
				// read log msg from channel
				if err := websocket.Message.Send(ws, msg); err != nil {
					// TODO better way to deal this situation
					fmt.Println(err)
					s.msgs <- msg
				}
			}
		}
	}
}

// WriteMsg write message to msg channel
func (s *WSWriter) WriteMsg(when time.Time, msg string, level int) error {
	if level > s.Level {
		return nil
	}
	// write msg to channel
	msg = when.String() + " " + msg
	select {
	case s.msgs <- msg:
	default:
		// TODO msg drooped when s.msgs channel full
		return fmt.Errorf("msg %s was discarded by WSWriter, cause s.msgs channel full, current size %d", len(s.msgs))
	}
	return nil
}

// Flush implementing method. empty.
func (s *WSWriter) Flush() {
	return
}

// Destroy implementing method. empty.
func (s *WSWriter) Destroy() {
	return
}

func init() {
	Register(AdapterWS, newWSWriter)
}
