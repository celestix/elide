package elide

import (
	"elide/internal/helpers"
	"encoding/json"
)

func (c *Client) ResolveUsername(username string) (*Chat, error) {
	buf, err := json.Marshal(Body{
		Method: "resolveUsername",
		Data: ResolveUsernameBody{
			Username: username,
		},
	})
	if err != nil {
		return nil, err
	}
	buf, err = c.MakeRequest(buf)
	if err != nil {
		return nil, err
	}
	var res Response[*Chat]
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return nil, err
	}
	if res.Ok {
		(&res).Result.ID = helpers.PatchChatIdToBotApi(res.Result.ID)
	} else {
		return nil, &TelegramError{rpc: res.Error}
	}
	return res.Result, nil
}

func (c *Client) GetMessages(chatId int64, messages []int) (*Message, error) {
	buf, err := json.Marshal(Body{
		Method: "getMessages",
		Data: GetMessagesBody{
			ChatId:     chatId,
			MessageIds: messages,
		},
	})
	if err != nil {
		return nil, err
	}
	buf, err = c.MakeRequest(buf)
	if err != nil {
		return nil, err
	}
	var res Response[*Message]
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		return nil, &TelegramError{rpc: res.Error}
	}
	return res.Result, nil
}
