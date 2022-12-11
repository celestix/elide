package api

import (
	"context"
	"elide/internal/api/config"
	"elide/internal/api/methods"
	"elide/internal/api/server"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/anonyindian/gotgproto"
	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/gotgproto/sessionMaker"
	"github.com/gotd/td/telegram"
)

const tgConnectionTimeout = 10

type Engine struct {
	appId                 int
	apiHash               string
	token                 string
	server                *server.Server
	establishedConnection bool
}

func CreateEngine(protocol, port, appId int, apiHash, token string) *Engine {
	defer log.Println("[Elide][Engine]: Created a new instance...")
	return &Engine{
		appId:   appId,
		apiHash: apiHash,
		token:   token,
		server:  server.New(protocol, port),
	}
}

func (e *Engine) establishTelegramConn() {
	gotgproto.StartClient(&gotgproto.ClientHelper{
		AppID:      e.appId,
		ApiHash:    e.apiHash,
		Session:    sessionMaker.NewSession("elide-tg", sessionMaker.Session),
		BotToken:   e.token,
		Dispatcher: dispatcher.MakeDispatcher(),
		TaskFunc: func(contx context.Context, client *telegram.Client) error {
			go func() {
				for {
					if gotgproto.Sender != nil {
						e.establishedConnection = true
						break
					}
				}
				config.Debugln("[Elide-DEBUG][Engine]: Established a connection to telegram...")
				e.server.Context = ext.NewContext(contx, gotgproto.Api, gotgproto.Self, gotgproto.Sender, nil)
				config.Debugln("[Elide-DEBUG][Engine][Server]: Populated server with pseudo context")
			}()
			return nil
		},
	})
}

func (e *Engine) Run() {
	log.Println("[Elide][Engine]: Establishing a connection to telegram...")
	go e.establishTelegramConn()
	i := 0
	for !e.establishedConnection {
		if i == tgConnectionTimeout {
			log.Println("[Elide][Engine]: Failed to establish connection to telegram: Timeout")
			os.Exit(1)
		}
		// Wait for telegram to connect
		time.Sleep(1 * time.Second)
		i++
	}
	log.Println("[Elide][Engine]: Loading method handlers...")
	e.loadHandlers()
	log.Println("[Elide][Engine]: Starting server...")
	e.server.Start()
}

func (e *Engine) loadHandlers() {
	e.server.AddHandler("echo", func(_ *ext.Context, body json.RawMessage) (any, error) {
		if body == nil {
			return nil, errors.New("no data provided")
		}
		return body, nil
	})
	e.server.AddHandler("resolveUsername", methods.ResolveUsername)
	e.server.AddHandler("deleteMessages", methods.DeleteMessages)
	e.server.AddHandler("getMessages", methods.GetMessages)
}
