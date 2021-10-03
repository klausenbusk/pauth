package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

var (
	pamExecEnvironmentVariables = [...]string{"PAM_RHOST", "PAM_RUSER", "PAM_SERVICE", "PAM_TTY", "PAM_USER", "PAM_TYPE"}
)

func auth(server string, authUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	c, _, err := websocket.Dial(ctx, server, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close(websocket.StatusInternalError, "")

	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Write(ctx, websocket.MessageText, []byte(id.String())); err != nil {
		log.Fatal(err)
	}

	variables := make(map[string]string)
	for _, env := range pamExecEnvironmentVariables {
		if val, ok := os.LookupEnv(env); ok {
			variables[env] = val
		}
	}

	if err := c.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("%s,%s", authUUID, variables))); err != nil {
		log.Fatal(err)
	}

	typ, buf, err := c.Read(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if typ == websocket.MessageBinary {
		c.Close(websocket.StatusUnsupportedData, "")
		log.Fatal("Invalid message")
	}
	arr := strings.SplitN(string(buf), ",", 2)
	if len(arr) != 2 {
		log.Fatal("Invalid message")
	}
	if arr[1] == "deny" {
		// FIXME: error shouldn't be used for this?
		log.Fatal("Auth denied")
	}
	if arr[1] != "accept" {
		log.Fatal("Invalid message")
	}

	c.Close(websocket.StatusNormalClosure, "")
	return nil
}
