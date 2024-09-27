package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"tgapiV2/internal/pg"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/rs/zerolog/log"
)

type MID struct {
	ID map[string]int64 `json:"PeerID"`
}

type HMG struct {
	Chats   []RawDial `json:"Chats"`
	Message []RawText `json:"Messages"`
}

type RawDial struct {
	Hash int64 `json:"AccessHash"`
	Id   int64 `json:"ID"`
}

type RawText struct {
	Message string           `json:"Message"`
	MsgId   int              `json:"ID"`
	Date    uint64           `json:"Date"`
	ID      map[string]int64 `json:"PeerID"`
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

func DialogsParse(ctx context.Context, client *telegram.Client, db *pg.DB, chats []int64) error {

	hs := HMG{}

	dial, err := client.API().MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{OffsetPeer: &tg.InputPeerChannel{ChannelID: chats[0]}})
	if err != nil {
		return err
	}
	dialogs, err := json.MarshalIndent(dial, "", "")
	if err != nil {
		return err
	}
	err = json.Unmarshal(dialogs, &hs)
	if err != nil {
		return err
	}
	ids := hs.Chats

	max := hs.Message

	offset := make(map[int64]int)

	dil := make(map[int64]int64)

	for i := range ids {
		dil[ids[i].Id] = ids[i].Hash
		offset[ids[i].Id] = max[i].MsgId

	}

	positive := make(map[int64]int64)

	for _, v := range chats {
		if _, ok := dil[v]; ok {
			positive[v] = dil[v]

		}

	}

	for chat, hash := range positive {
		step := offset[chat]
		fmt.Println("LOL", step)

		for {

			fmt.Println("step", step)

			if step != 0 {
				time.Sleep(1000 * time.Millisecond)
				history, err := client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
					Peer: &tg.InputPeerChannel{
						ChannelID:  chat,
						AccessHash: hash,
					},
					OffsetID: step, Limit: 100,
				})

				fmt.Println("STEP", step)
				if err != nil {
					return err
				}

				messages, err := json.MarshalIndent(history, "", "")
				if err != nil {
					return err
				}

				err = json.Unmarshal(messages, &hs)
				if err != nil {
					return err
				}
				fmt.Println("BYTE", messages)
				pgstep, err := db.AddHistPG(ctx, messages)
				if err != nil {
					return err
				}

				step = pgstep

			} else {
				break
			}
		}

	}
	fmt.Println(offset)
	return nil

}
