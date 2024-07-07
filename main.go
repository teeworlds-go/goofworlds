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

	g := &game.Game{}
	g.Client = *teeworlds7.NewClient()
	g.Client.Name = "goofy"
	g.Ip = "127.0.0.1"
	g.Port = 8303
	if len(os.Args) == 2 {
		console.ExecLine(os.Args[1], g)
	} else if len(os.Args) > 2 {
		panic(fmt.Errorf("more than 1 cli arg not supported got %d", len(os.Args)-1))
	}

	g.ChatInp = game.NewTextInput(func(line string) {
		g.Client.SendChat(line)
		g.ChatInp.Deactivate()
	})

	g.Components = append(g.Components, &components.Chat{})

	for _, c := range g.Components {
		c.OnInit()
	}

	g.Client.OnChat(func(msg *messages7.SvChat, defaultAction teeworlds7.DefaultAction) {
		for _, c := range g.Components {
			c.OnChatMsg(msg, &g.Client)
		}
	})

	g.Client.OnSnapshot(func(snap *snapshot7.Snapshot, defaultAction teeworlds7.DefaultAction) {
		char, err := g.Client.SnapFindCharacter(g.Client.LocalClientId)
		if err == nil {
			g.Camera.X = char.X
			g.Camera.Y = char.Y
		}
	})

	go func() {
		fmt.Printf("connecting to %s:%d\n", g.Ip, g.Port)
		g.Client.Connect(g.Ip, g.Port)
	}()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
