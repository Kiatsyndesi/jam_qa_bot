package mm_client

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"jam_qa_bot/internal/app/mm-client/mm_bot"
	"jam_qa_bot/internal/helpers"
	"os"
)

func StartDebuggingMMBot() {
	bot := mm_bot.NewBot()

	bot.SetupGracefulShutdownForBot()

	bot.Client = model.NewAPIv4Client(os.Getenv("MM_HTTPS_URL"))

	// Lets test to see if the mattermost server is up and running
	bot.MakeSureBotIsRunning()

	// Set bot token for using mattermost api
	bot.Client.SetOAuthToken(os.Getenv("AUTH_TOKEN"))
	bot.LoginAsTheBotUser()
	// Assign botTeam
	bot.FindBotTeam()

	// this is working
	bot.CreateBotDebuggingChannelIfNeeded()
	bot.SendMsgToDebuggingChannel("_"+os.Getenv("BOT_NAME")+" has **started** running_", "")

	// Lets start listening to some channels via the websocket!
	for {
		webSocketClient, err := model.NewWebSocketClient4(os.Getenv("MM_WSS_URL"), bot.Client.AuthToken)
		if err != nil {
			println("We failed to connect to the web socket")
			helpers.PrintError(err)
		}
		println("Connected to WS")
		webSocketClient.Listen()

		for resp := range webSocketClient.EventChannel {
			bot.HandleWebSocketResponse(resp)
		}
	}
}
