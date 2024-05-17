package mcserver

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestPingServer(t *testing.T) {
	vanillaServer := default_servers["vanilla"]
	vanillaServer.PrivateHost = "fdaa:1:8e0d:a7b:83:8f2b:d0f1:2"
	serverInfo, err := PingServer(context.Background(), vanillaServer, time.Second*5)
	if err != nil {
		fmt.Println("err...", err)
	} else if serverInfo != nil {
		fmt.Println(serverInfo.Description)
		fmt.Printf("%d/%d\n", serverInfo.Players.Online, serverInfo.Players.Max)
		fmt.Println(serverInfo.Version.Name, serverInfo.Version.Protocol)
	} else {
		fmt.Println("No response!")
	}

	assert.Equal(t, false, true)
}
