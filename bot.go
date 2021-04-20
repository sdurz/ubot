// Package ubot provides types and
package ubot

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// matcherHandler encapsulates an UMatcher and the corresponding UHandler
type matcherHandler struct {
	matcher UMatcher
	handler UHandler
}

// evaluate execute the handler func if the matcher returns true
func (m *matcherHandler) evaluate(ctx context.Context, bot *Bot, message O) (result bool, err error) {
	if m.matcher(bot, message) {
		result, err = m.handler(ctx, bot, message)
	}
	return
}

// Bot is the main type of ubot.
// It implements a bot API frontend.
type Bot struct {
	Configuration         Configuration
	apiClient             ApiClient
	BotUser               UUser
	messageMHs            []matcherHandler
	editedMessageMHs      []matcherHandler
	channelPostMHs        []matcherHandler
	editedChannelPostMHs  []matcherHandler
	inlineQueryMHs        []matcherHandler
	chosenInlineResultMHs []matcherHandler
	callbackQueryMHs      []matcherHandler
	myChatMemberMHs       []matcherHandler
	chatMemberMHs         []matcherHandler
}

// NewBot creates a new Bot for the given configuration
func NewBot(configuration *Configuration) (result *Bot) {
	result = &Bot{
		Configuration: *configuration,
		apiClient:     &httpApiClient{},
	}
	return
}

// methodURL transforms a method name in the corresponding API url
func (b *Bot) methodURL(method string) (result string) {
	if method == "" {
		panic("Emtpy method")
	}
	result = "https://api.telegram.org/bot" + b.Configuration.APIToken + "/" + method
	return
}

// getUpdatesSource
func (b *Bot) getUpdatesSource(ctx context.Context, updatesChan chan O) {
	var nextUpdate int64 = 0
	var ok bool
	for {
		select {
		case <-ctx.Done():
			log.Println("done with getUpdatesSource")
			return
		default:
			getURL := b.methodURL("getUpdates") + "?offset=" + strconv.FormatInt(nextUpdate, 10)
			var responseUpdates interface{}
			responseUpdates, err := b.apiClient.GetJson(getURL)
			if err != nil {
				log.Println("Error while retrieving updates", err)
				continue
			}

			var updates A
			if updates, ok = responseUpdates.([]interface{}); !ok {
				log.Fatalln("updates result not a JSON array")
			}

			if len(updates) > 0 {
				for _, update := range updates {
					var (
						updateID int64
						oUpdate  O
					)
					if oUpdate, ok = update.(map[string]interface{}); !ok {
						log.Println("update not an O")
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

func (b *Bot) serverSource(ctx context.Context, updates chan O) {

	serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		log.Println(string(body))

		var (
			rawUpdate interface{}
			update    O
			err       error
			ok        bool
		)
		if err = json.Unmarshal(body, rawUpdate); err != nil {
			log.Printf("Error decoding body: %v", err)
			http.Error(w, "can't decode body", http.StatusBadRequest)
			return
		}

		if update, ok = rawUpdate.(map[string]interface{}); !ok {
			log.Printf("Error decoding body: %v", err)
			return
		}

		err = b.process(ctx, &update)

		if err != nil {
			log.Println("Update processing error: ", err)
		}
	})

	mux := http.NewServeMux()
	mux.Handle("/ubot/"+b.Configuration.APIToken, serverHandler)
	srv := &http.Server{
		Addr:    b.Configuration.ServerPort,
		Handler: mux,
	}

	if b.Configuration.WebhookUrl == "" {
		log.Fatal("empty webhook url")
	}

	go http.ListenAndServe(b.Configuration.ServerPort, mux)
	ok, err := b.SetWebhook(map[string]interface{}{
		"url": b.Configuration.WebhookUrl,
	})
	if !ok || err != nil {
		log.Fatal("Can't set webhook")
		return
	}

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Println("Server stopped")
}

func (b *Bot) AddChatMemberHandler(matcher UMatcher, handler UHandler) {
	b.chatMemberMHs = append(b.chatMemberMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) Forever(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	source := b.getUpdatesSource
	if b.Configuration.LongPoll {
		source = b.getUpdatesSource
	}

	if user, err := b.GetMe(); err == nil {
		b.BotUser = *user
	} else {
		log.Fatal(err)
	}

	updates := make(chan O)
	go source(ctx, updates)

	semaphore := make(chan int, b.Configuration.WorkerNo)
	for {
		select {
		case <-ctx.Done():
			log.Println("forever is over")
			return nil
		case update := <-updates:
			semaphore <- 1
			go func() {
				b.process(ctx, update)
				<-semaphore
			}()
		}
	}
}

func (b *Bot) process(ctx context.Context, update O) (err error) {
	var (
		ok              bool
		stop            bool
		matcherHandlers []matcherHandler
		rawPayload      interface{}
		payload         O
	)
	if rawPayload, ok = update["message"]; ok {
		matcherHandlers = b.messageMHs
	} else if rawPayload, ok = update["edited_message"]; ok {
		matcherHandlers = b.editedMessageMHs
	} else if rawPayload, ok = update["channel_post"]; ok {
		matcherHandlers = b.channelPostMHs
	} else if rawPayload, ok = update["edited_channel_post"]; ok {
		matcherHandlers = b.editedChannelPostMHs
	} else if rawPayload, ok = update["callback_query"]; ok {
		matcherHandlers = b.callbackQueryMHs
	} else if rawPayload, ok = update["inline_query"]; ok {
		matcherHandlers = b.inlineQueryMHs
	} else if rawPayload, ok = update["chosen_inline_result"]; ok {
		matcherHandlers = b.chosenInlineResultMHs
	} else {
		err = errors.New("update without data")
	}
	if payload, ok = rawPayload.(map[string]interface{}); !ok {
		err = errors.New("update payload not an O")
		return
	}
	if err == nil {
		for _, ha := range matcherHandlers {
			if stop, err = ha.evaluate(ctx, b, payload); err != nil {
				log.Println(err)
				break
			}
			if stop {
				break
			}
		}
	}
	return
}

func (b *Bot) doGet(method string, data O) (interface{}, error) {
	return b.apiClient.GetJson(b.methodURL(method))
}

func (b *Bot) doPost(method string, request O) (interface{}, error) {
	return b.apiClient.PostJson(b.methodURL(method), request)
}

func (b *Bot) doPostMultipart(method string, request O) (interface{}, error) {
	return b.apiClient.PostMultipart(b.methodURL(method), request)
}

func (b *Bot) AddMessageHandler(matcher UMatcher, handler UHandler) {
	b.messageMHs = append(b.messageMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) AddChannelPostHandler(matcher UMatcher, handler UHandler) {
	b.channelPostMHs = append(b.channelPostMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) AddEditedMessageHandler(matcher UMatcher, handler UHandler) {
	b.editedMessageMHs = append(b.editedMessageMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) AddEditedChannelPostHandler(matcher UMatcher, handler UHandler) {
	b.editedChannelPostMHs = append(b.editedChannelPostMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) AddCallbackQueryHandler(matcher UMatcher, handler UHandler) {
	b.callbackQueryMHs = append(b.callbackQueryMHs, matcherHandler{matcher: matcher, handler: handler})
}

func (b *Bot) AddMyChatMemberHandler(matcher UMatcher, handler UHandler) {
	b.myChatMemberMHs = append(b.myChatMemberMHs, matcherHandler{matcher: matcher, handler: handler})
}

// API Mappings
func (b *Bot) GetMe() (result *UUser, err error) {
	result = &UUser{}
	_, err = b.doGet("getMe", nil)
	return
}

func (b *Bot) LogOut() (err error) {
	_, err = b.doGet("logOut", nil)
	return
}

func (b *Bot) Close() (err error) {
	_, err = b.doGet("close", nil)
	return
}

// SendMessage sends a message
func (b *Bot) SendMessage(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPost("sendMessage", request); err != nil {
		result = response.(map[string]interface{})
	}
	return
}

// ForwardMessage forwads a message
func (b *Bot) ForwardMessage(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPost("forwardMessage", request); err != nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendPhoto sends a photo
func (b *Bot) SendPhoto(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendPhoto", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendVideo sends a video
func (b *Bot) SendVideo(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendVideo", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendAudio sends an audio
func (b *Bot) SendAudio(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendAudio", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendAnimation sends an animation
func (b *Bot) SendAnimation(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendAnimation", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendVoice sends a voice
func (b *Bot) SendVoice(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPostMultipart("sendVoice", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// SendLocation sends a location
func (b *Bot) SendLocation(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPost("sendLocation", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

// TODO: complete Send* messages
// SetWebhook implements setWebhook from Telegram Bot API
func (b *Bot) SetWebhook(request O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("setWebhook", request); err == nil {
		result = response.(bool)
	}
	return
}

// DeleteWebhook implements deleteWebhook from Telegram Bot API
func (b *Bot) DeleteWebhook(request O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("deleteWebhook", request); err == nil {
		result = response.(bool)
	}
	return
}

func (b *Bot) GetWebhookInfo() (result O, err error) {
	var response interface{}
	if response, err = b.doGet("getWebhookInfo", result); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) GetFile(fileId string) (result O, err error) {
	var response interface{}
	if response, err = b.doGet("getFile?file_id="+fileId, nil); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) GetUserProfilePhotos(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPost("getUserProfilesPhotos", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) SendDice(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPost("sendDice", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) SendMediaGroup(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPost("sendMediaGroup", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) SendChatAction(request O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("sendChatAction", request); err == nil {
		result = response.(bool)
	}
	return
}

func (b *Bot) AnswerCallbackQuery(request O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("answerCallbackQuery", request); err == nil {
		result = response.(bool)
	}
	return
}

func (b *Bot) PinChatMessage(request O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("answerCallbackQuery", request); err == nil {
		result = response.(bool)
	}
	return
}

func (b *Bot) UnpinChatMessage(request O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("answerCallbackQuery", request); err == nil {
		result = response.(bool)
	}
	return
}

func (b *Bot) UnpinAllChatMessages(request O) (result bool, err error) {
	var response interface{}
	if response, err = b.doPost("answerCallbackQuery", request); err == nil {
		result = response.(bool)
	}
	return
}

func (b *Bot) GetChat(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPost("getChat", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) LeaveChat(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPost("leaveChat", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) GetChatMembersCount(request O) (result int64, err error) {
	var response interface{}
	if response, err = b.doPost("getChat", request); err == nil {
		result = response.(int64)
	}
	return
}

func (b *Bot) SetMyCommands(request O) (result O, err error) {
	var response interface{}
	if response, err = b.doPost("setMyCommands", request); err == nil {
		result = response.(map[string]interface{})
	}
	return
}

func (b *Bot) GetMyCommands() (result O, err error) {
	var response interface{}
	if response, err = b.doGet("getMyCommands", result); err == nil {
		result = response.(map[string]interface{})
	}
	return
}
