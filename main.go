package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/teeworlds-go/goofworlds/game"
	"github.com/teeworlds-go/protocol/snapshot7"
	"github.com/teeworlds-go/protocol/teeworlds7"
)

func ExecLine(line string, game *game.Game) {
	if strings.HasPrefix(line, "connect ") {
		fullIp := strings.Split(line, " ")[1]
		game.Ip = strings.Split(fullIp, ":")[0]
		var err error
		game.Port, err = strconv.Atoi(strings.Split(fullIp, ":")[1])
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("unknown command: %s\n", line)
		os.Exit(1)
	}
}

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("goofworlds")

	game := &game.Game{}
	game.Client = *teeworlds7.NewClient()
	game.Client.Name = "goofy"
	game.Ip = "127.0.0.1"
	game.Port = 8303
	if len(os.Args) == 2 {
		ExecLine(os.Args[1], game)
	} else if len(os.Args) > 2 {
		panic(fmt.Errorf("more than 1 cli arg not supported got %d", len(os.Args)-1))
	}

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
