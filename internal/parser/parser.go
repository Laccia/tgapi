package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"tgapiV2/internal/pg"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/rs/zerolog/log"
)

type MID struct {
	ID map[string]int64 `json:"PeerID"`
}

type HMG struct {
	Chats []RawDial `json:"Chats"`
}

type RawDial struct {
	Hash int `json:"AccessHash"`
}

func MessageParse(ctx context.Context, update tg.MessageClass, db *pg.DB, chats []int64) error {

	newmsg, err := json.MarshalIndent(update, "", "")
	if err != nil {
		return err
	}

	st := MID{}

	err = json.Unmarshal(newmsg, &st)

	if err != nil {
		return err
	}

	id := st.ID["ChannelID"]

	for _, v := range chats {
		if v == id {
			err = db.AddMsgPG(ctx, newmsg, id)
			if err != nil {
				log.Err(err)
			}
			break
		}
	}

	return nil
}

func DialogsParse(ctx context.Context, client *telegram.Client) error {
	dial, err := client.API().MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{OffsetPeer: &tg.InputPeerChannel{ChannelID: 1195476893}, Limit: 1})
	if err != nil {
		return err
	}
	dialogs, err := json.MarshalIndent(dial, "", "")
	if err != nil {
		return err
	}

	hs := HMG{}

	err = json.Unmarshal(dialogs, &hs)
	if err != nil {
		return err
	}

	hash := hs.Chats[0].Hash

	fmt.Println(hash)
	return nil
}
