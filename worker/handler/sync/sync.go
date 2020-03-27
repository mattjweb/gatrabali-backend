package sync

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/apps4bali/gatrabali-backend/common/constant"
	"github.com/apps4bali/gatrabali-backend/common/types"
	"github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"

	"worker/handler/sync/service"
)

// Handler sync data from Miniflux to Firestore
func Handler(client *firestore.Client) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		ctx := context.Background()
		msg := c.Locals(pubsub.LocalsKey).(*pubsub.Message)

		var payload *types.SyncPayload
		if err := json.Unmarshal(msg.Message.Data, &payload); err != nil {
			c.Next(err)
			return
		}
		if payload.ID == nil || payload.Type == nil || payload.Op == nil {
			c.Next(errors.New("Invalid message payload: missing id, type or op"))
			return
		}

		switch *payload.Type {
		case constant.TypeCategory:
			if err := service.StartCategorySync(ctx, client, payload); err != nil {
				c.Next(err)
				return
			}
		case constant.TypeFeed:
			if err := service.StartFeedSync(ctx, client, payload); err != nil {
				c.Next(err)
				return
			}
		case constant.TypeEntry:
			if err := service.StartEntrySync(ctx, client, payload); err != nil {
				c.Next(err)
				return
			}
		}
		c.SendStatus(http.StatusOK)
	}
}
