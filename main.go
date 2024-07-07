package main

import (
	"log"
	"os"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/teeworlds-go/protocol/messages7"
	"github.com/teeworlds-go/protocol/snapshot7"
	"github.com/teeworlds-go/protocol/teeworlds7"
)



const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	teeSprite *ebiten.Image
)

type Camera struct {
	X int
	Y int
}

type CameraOffset struct {
	X int
	Y int
}

func getImageFromFilePath(filePath string) (image.Image, error) {
    f, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    image, _, err := image.Decode(f)
    return image, err
}

func init() {
	// Preload images
	img, err := getImageFromFilePath("img/tee.png")
	if err != nil {
		panic(err)
	}
	teeSprite = ebiten.NewImageFromImage(img)
}

func getCameraOffset(camera Camera) CameraOffset {
	wc := screenWidth / 2
	hc := screenHeight / 2
	x := -camera.X + wc
	y := -camera.Y + hc
	return CameraOffset{
		X: x,
		Y: y,
	}
}

type Game struct {
	Client teeworlds7.Client
	Camera Camera
}

func (g *Game) Update() error {
	g.Client.Game.Input.Direction = 0
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Client.Game.Input.Direction = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Client.Game.Input.Direction = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.Client.SendMessage(&messages7.CtrlClose{})
		os.Exit(0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, character := range g.Client.Game.Snap.Characters {
		offset := getCameraOffset(g.Camera)

		screenX := float32(character.X) + float32(offset.X)
		screenY := float32(character.Y) + float32(offset.Y)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.069, 0.069)
		op.GeoM.Translate(float64(screenX) - 32, float64(screenY) - 32)
		screen.DrawImage(teeSprite, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("goofworlds")

	game := &Game{}
	game.Client = *teeworlds7.NewClient()
	game.Client.Name = "goofy"

	game.Client.OnSnapshot(func(snap *snapshot7.Snapshot, defaultAction teeworlds7.DefaultAction) {
		char, err := game.Client.SnapFindCharacter(game.Client.LocalClientId)
		if err == nil {
			game.Camera.X = char.X
			game.Camera.Y = char.Y
		}
	})

	go func() {
		game.Client.Connect("127.0.0.1", 8303)
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
