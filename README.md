# uBot
uBot is a minimalistic Telegram BOT API library for Golang that aims to be complete, idiomatic and extensible.

## Overview
uBot is a bot framework for Telegram Bot API that I'm writing to support my own bot implementations. 
Instead of providing a full API mapping, uBot relies on a minimal wrapping of Golang JSON marshal/unmarshal functionalities, with JSON accessor methods. It also offers a simple but extensible mechanism for message routing and filtering that allows the developers to write, reuse and compose their own matchers and handlers.

The method of Telegram bot API are mappped one on one to the Bot object, the messages that are sent to the server are to be composed as specified on Telegram's reference. There's no need to memorize an additional
layer nor coding conventions (in turn there's no guarantee that the sent messages are well formed, beware of 400's error).

Strengths or uBot are:
- Minimal footprint
- Easily extensible
- It doesn't create any additional abstraction layer over Telegram bot API.
- Intuitive JSON response handling
- Pluggable update source for receiving updates (in progress)
- Context and WaitGroup aware bots
- Custom method invocations. Invoke latest API methods even if the current uBot version isn't up to date.

weaknesses are:
- Not completely type safe
- Not as proven as other libraries
- API still not fully tested (help needed)
  

This package has been written in the process of learning Golang, critics and contributions are welcome.

## Get started
A basic bot that can receive a message and send a reply, getting updates via long poll:

```golang

func main() (result *ubot.Bot, err error) {  
	bot := ubot.NewBot(ubot.Configuration{APIToken: "<yourAPIToken>", LongPoll: true})

	bot.AddMessageHandler(ubot.Any,
	), func(ctx context.Context, bot *ubot.Bot, message ubot.O) (done bool, err error) {
        var chatID int64
		if chatID, err = message.GetInteger("from.id"); err == nil {
		    _, err = bot.SendMessage(O{"chat_id": chatID, "text": "I got your message"})
        }
		return
	})

    var wg sync.WaitGroup
    go bot.Forever(context.BackgroundContext(), &wg)
    wg.Add(1)

    wg.Wait()
}

```

## UMatcher and UHandler

An UMatcher is a func that is executed to check wheter an update is to be handled:

```golang
type UMatcher func(*Bot, O) bool
```

An UHandler is a func that actually handles the incoming payload:
```golang
type UHandler func(context.Context, *Bot, O) (bool, error)
```

UMatcher(s) can be reused and composed, _uBot_ provides quite a few boolean operators that help to compose simpler matchers.

## Caveats
Methods mapping is still not complete.

## License
uBot is distributed under MIT.