package parser

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"gitlab.figvam.ru/figvam/tgapi/internal/pg"
)

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

type Mgs struct {
	Positive map[int64]int64
	Hs       HMG
	Client   *telegram.Client
	DB       *pg.DB
}

func New(ctx context.Context, client *telegram.Client, db *pg.DB, chats []int64) (*Mgs, error) {

	hs := HMG{}

	dial, err := client.API().MessagesGetDialogs(ctx,
		&tg.MessagesGetDialogsRequest{OffsetPeer: &tg.InputPeerChannel{ChannelID: chats[0]}, Limit: 100})

	if err != nil {
		return nil, err
	}

	dialogs, err := json.MarshalIndent(dial, "", "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dialogs, &hs)
	if err != nil {
		return nil, err
	}

	ids := hs.Chats

	max := hs.Message

	dil := make(map[int64]int64)

	for i := range max {
		dil[ids[i].Id] = ids[i].Hash
	}

	positive := make(map[int64]int64)

	for _, v := range chats {
		if _, ok := dil[v]; ok {
			positive[v] = dil[v]

		}

	}

	return &Mgs{Positive: positive, Hs: hs, Client: client, DB: db}, nil

}

func (mgs *Mgs) DialogsParse(ctx context.Context) error {

	for chat, hash := range mgs.Positive {
		rows, err := mgs.DB.CheckExist(ctx, chat)

		if err != nil {
			return err
		}

		if rows == 0 {
			step := 1
			err = mgs.HistoryAdd(ctx, step, chat, hash)
			if err != nil {
				return err
			}
		} else {
			step, err := mgs.DB.HistoryCheck(ctx, chat)
			if err != nil {
				return err
			}
			err = mgs.HistoryAdd(ctx, step, chat, hash)
			if err != nil {
				return err
			}
		}

	}
	return nil

}

func (mgs *Mgs) HistoryAdd(ctx context.Context, offset int, chat int64, hash int64) error {
	step := offset
	for {
		time.Sleep(1 * time.Second)
		history, err := mgs.Client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
			Peer: &tg.InputPeerChannel{
				ChannelID:  chat,
				AccessHash: hash,
			},
			OffsetID: step, Limit: 100, AddOffset: -100,
		})

		if err != nil {
			return err
		}

		messages, err := json.MarshalIndent(history, "", "")
		if err != nil {
			return err
		}
		err = json.Unmarshal(messages, &mgs.Hs)
		if err != nil {
			return err
		}
		if step != mgs.Hs.Message[0].MsgId {

			pgstep, err := mgs.DB.AddHistPG(ctx, messages)
			if err != nil {
				return err
			}
			step = pgstep

		} else {
			break
		}
	}
	return nil
}
