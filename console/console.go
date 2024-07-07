package console

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/teeworlds-go/goofworlds/game"
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
