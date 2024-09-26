package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
	Hash int64 `json:"AccessHash"`
	Id   int64 `json:"ID"`
}

type RawMSG struct {
	Messages []RawText `json:"Messages"`
}

type RawText struct {
	Message string           `json:"Message"`
	MsgId   int64            `json:"ID"`
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

	ms := RawMSG{}
	dial, err := client.API().MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{OffsetPeer: &tg.InputPeerChannel{ChannelID: chats[0]}})
	if err != nil {
		fmt.Println(err)
	}
	dialogs, err := json.MarshalIndent(dial, "", "")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(dialogs, &hs)
	if err != nil {
		fmt.Println(err)
	}
	ids := hs.Chats
	err = os.WriteFile("utils/dialogs.json", dialogs, 0644)
	if err != nil {
		fmt.Println(err)
	}

	dil := make(map[int64]int64)
	for i := range ids {
		dil[ids[i].Id] = ids[i].Hash
	}

	positive := make(map[int64]int64)

	for _, v := range chats {
		if _, ok := dil[v]; ok {
			positive[v] = dil[v]
		}

	}

	for chat, hash := range positive {
		history, err := client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
			Peer: &tg.InputPeerChannel{
				ChannelID:  chat,
				AccessHash: hash,
			},
			Limit: 100,
		})

		if err != nil {
			fmt.Println(err)
		}

		messages, err := json.MarshalIndent(history, "", "")
		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal(messages, &ms)
		if err != nil {
			fmt.Println(err)
		}
		msgs := ms.Messages
		var msgid int64
		var msg string
		var id int64
		var date uint64

		for i := range msgs {
			msg = msgs[i].Message
			msgid = msgs[i].MsgId
			date = msgs[i].Date
			id = msgs[i].ID["ChannelID"]
			err := db.AddHistPG(ctx, msg, msgid, date, id)

			if err != nil {
				fmt.Println(err)
			}

		}
	}

	return nil

}
