package parser

import (
	"context"
	"encoding/json"
	"tgapiV2/internal/pg"

	"github.com/gotd/td/tg"
	"github.com/rs/zerolog/log"
)

type MID struct {
	ID map[string]int64 `json:"PeerID"`
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
