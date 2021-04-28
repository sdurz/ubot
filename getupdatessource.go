package ubot

import (
	"context"
	"log"
	"strconv"

	"github.com/sdurz/axon"
)

// GetUpdatesSource is a ServerSource tha get updates vi long polling
// See https://core.telegram.org/bots/api#getupdates
func GetUpdatesSource(bot *Bot, ctx context.Context, updatesChan chan axon.O) {
	var nextUpdate int64 = 0
	var ok bool
	for {
		select {
		case <-ctx.Done():
			log.Println("done with getUpdatesSource")
			return
		default:
			getURL := bot.methodURL("getUpdates") + "?offset=" + strconv.FormatInt(nextUpdate, 10)
			var responseUpdates interface{}
			responseUpdates, err := bot.apiClient.GetJson(getURL)
			if err != nil {
				log.Println("Error while retrieving updates", err)
				continue
			}

			var updates axon.A
			if updates, ok = responseUpdates.([]interface{}); !ok {
				log.Fatalln("updates result not a JSON array")
			}

			if len(updates) > 0 {
				for _, update := range updates {
					var (
						updateID int64
						oUpdate  axon.O
					)
					if oUpdate, ok = update.(map[string]interface{}); !ok {
						log.Println("update not an axon.O")
						continue
					}
					if updateID, err = oUpdate.GetInteger("update_id"); err != nil {
						log.Println("update does not have an integer id")
						continue
					}
					if updateID > nextUpdate {
						nextUpdate = updateID
					}
					updatesChan <- oUpdate
				}
				nextUpdate++
			}
		}
	}
}
