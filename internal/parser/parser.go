package parser

import (
	"context"
	"encoding/json"
	"fmt"

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
	Offset   map[int64]int
	Hs       HMG
	Client   *telegram.Client
	DB       *pg.DB
}

func NewHistory(ctx context.Context, client *telegram.Client, db *pg.DB, chats []int64) (*Mgs, error) {

	hs := HMG{}

	dial, err := client.API().MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{OffsetPeer: &tg.InputPeerChannel{ChannelID: chats[0]}})
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

	return &Mgs{Positive: positive, Offset: offset, Hs: hs, Client: client, DB: db}, nil

}

func (mgs *Mgs) DialogsParse(ctx context.Context) error {

	for chat, hash := range mgs.Positive {
		rows, err := mgs.DB.CheckExist(ctx, chat)
		fmt.Println(rows)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if rows == 0 {
			err = AllHistoryADD(ctx, mgs.Offset, chat, hash, mgs.Hs, mgs.Client, mgs.DB)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			step, err := mgs.DB.HistoryCheck(ctx, chat)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = HistoryAdd(ctx, step, chat, hash, mgs.Hs, mgs.Client, mgs.DB)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

	}
	return nil

}

func AllHistoryADD(ctx context.Context, offset map[int64]int, chat int64, hash int64, hs HMG, client *telegram.Client, db *pg.DB) error {
	step := offset[chat]
	step = step + 1
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
			pgstep, err := db.AddAllHistPG(ctx, messages)
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

func HistoryAdd(ctx context.Context, offset int, chat int64, hash int64, hs HMG, client *telegram.Client, db *pg.DB) error {
	step := offset
	for {

		fmt.Println("step", step)

		if step != 0 {
			time.Sleep(1000 * time.Millisecond)
			history, err := client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
				Peer: &tg.InputPeerChannel{
					ChannelID:  chat,
					AccessHash: hash,
				},
				MinID: step, Limit: 100,
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
	return nil
}
