package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/teeworlds-go/goofworlds/console"
	"github.com/teeworlds-go/goofworlds/game"
	"github.com/teeworlds-go/goofworlds/game/components"
	"github.com/teeworlds-go/protocol/messages7"
	"github.com/teeworlds-go/protocol/snapshot7"
	"github.com/teeworlds-go/protocol/teeworlds7"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("goofworlds")

	game := &game.Game{}
	game.Client = *teeworlds7.NewClient()
	game.Client.Name = "goofy"
	game.Ip = "127.0.0.1"
	game.Port = 8303
	if len(os.Args) == 2 {
		console.ExecLine(os.Args[1], game)
	} else if len(os.Args) > 2 {
		panic(fmt.Errorf("more than 1 cli arg not supported got %d", len(os.Args)-1))
	}

	game.Components = append(game.Components, &components.Chat{})

	for _, c := range game.Components {
		c.OnInit()
	}

	game.Client.OnChat(func(msg *messages7.SvChat, defaultAction teeworlds7.DefaultAction) {
		for _, c := range game.Components {
			c.OnChatMsg(msg, &game.Client)
		}
	})

	game.Client.OnSnapshot(func(snap *snapshot7.Snapshot, defaultAction teeworlds7.DefaultAction) {
		char, err := game.Client.SnapFindCharacter(game.Client.LocalClientId)
		if err == nil {
			game.Camera.X = char.X
			game.Camera.Y = char.Y
		}
	})

	go func() {
		fmt.Printf("connecting to %s:%d\n", game.Ip, game.Port)
		game.Client.Connect(game.Ip, game.Port)
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
