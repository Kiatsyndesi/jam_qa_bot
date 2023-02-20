package mm_bot

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"jam_qa_bot/internal/helpers"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

type Bot struct {
	Client           *model.Client4
	WebSocketClient  *model.WebSocketClient
	BotUser          *model.User
	BotTeam          *model.Team
	DebuggingChannel *model.Channel
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) FindBotTeam() {
	if team, resp := b.Client.GetTeamByName(os.Getenv("TEAM_NAME"), ""); resp.Error != nil {
		println("We failed to get the initial load")
		println("or we do not appear to be a member of the team '" + os.Getenv("TEAM_NAME") + "'")
		helpers.PrintError(resp.Error)
		os.Exit(1)
	} else {
		b.BotTeam = team
	}
}

func (b *Bot) LoginAsTheBotUser() {
	if user, resp := b.Client.GetMe(os.Getenv("BOT_ETAG")); resp.Error != nil {
		println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		helpers.PrintError(resp.Error)
		os.Exit(1)
	} else {
		b.BotUser = user
	}
}

func (b *Bot) MakeSureBotIsRunning() {
	if props, resp := b.Client.GetOldClientConfig(""); resp.Error != nil {
		println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		helpers.PrintError(resp.Error)
		os.Exit(1)
	} else {
		println("Server detected and is running version " + props["Version"])
	}
}

func (b *Bot) SetupGracefulShutdownForBot() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if b.WebSocketClient != nil {
				b.WebSocketClient.Close()
			}

			b.SendMsgToDebuggingChannel("_"+os.Getenv("BOT_NAME")+" has **stopped** running_", "")
			os.Exit(0)
		}
	}()
}

func (b *Bot) CreateBotDebuggingChannelIfNeeded() {
	if rchannel, resp := b.Client.GetChannelByName(os.Getenv("DEBUG_CHANNEL_NAME"), b.BotTeam.Id, ""); resp.Error != nil {
		println("We failed to get the channels")
		helpers.PrintError(resp.Error)
	} else {
		b.DebuggingChannel = rchannel
		return
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = os.Getenv("DEBUG_CHANNEL_NAME")
	channel.DisplayName = "Debugging For Sample Bot"
	channel.Purpose = "This is used as a test channel for logging bot debug messages"
	channel.Type = model.CHANNEL_OPEN
	channel.TeamId = b.BotTeam.Id
	if rchannel, resp := b.Client.CreateChannel(channel); resp.Error != nil {
		println("We failed to create the channel " + os.Getenv("DEBUG_CHANNEL_NAME"))
		helpers.PrintError(resp.Error)
	} else {
		b.DebuggingChannel = rchannel
		println("Looks like this might be the first run so we've created the channel " + os.Getenv("DEBUG_CHANNEL_NAME"))
	}
}

func (b *Bot) SendMsgToDebuggingChannel(msg string, replyToId string) {
	post := &model.Post{}
	post.ChannelId = b.DebuggingChannel.Id
	post.Message = msg

	post.RootId = replyToId

	if _, resp := b.Client.CreatePost(post); resp.Error != nil {
		println("We failed to send a message to the logging channel")
		helpers.PrintError(resp.Error)
	}
}

func (b *Bot) HandleMsgFromDebuggingChannel(event *model.WebSocketEvent) {
	// If this isn't the debugging channel then lets ingore it
	if event.GetBroadcast().ChannelId != b.DebuggingChannel.Id {
		return
	}

	// Lets only responded to messaged posted events
	if event.EventType() != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	println("responding to debugging channel msg")

	post := model.PostFromJson(strings.NewReader(event.GetData()["post"].(string)))
	if post != nil {

		// ignore my events
		if post.UserId == b.BotUser.Id {
			return
		}

		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)alive(?:$|\W)`, post.Message); matched {
			b.SendMsgToDebuggingChannel("Yes I'm running", post.Id)
			return
		}
	}

	b.SendMsgToDebuggingChannel("I did not understand you!", post.Id)
}

func (b *Bot) HandleWebSocketResponse(event *model.WebSocketEvent) {
	b.HandleMsgFromDebuggingChannel(event)
}
