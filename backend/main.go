package main

import (
	"errors"
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/handlers"
	"github.com/Qwiri/gobby/pkg/validate"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
)

const (
	StateLobby gobby.State = 1 << iota
)

type Meta struct {
	Pack     *CardPack
	Selected map[string]int // Client Name -> Selected
}

func main() {
	log.Infof("Packs: %v", Packs)
	app := fiber.New()
	g := gobby.New(app)

	g.HandleAll(gobby.Handlers{
		"CHAT": handlers.Chat,
		"LIST": handlers.List,
	})

	g.MustOn(func(event *gobby.Join) {
		if len(event.Lobby.Clients) >= 2 {
			// TODO: cancel join if there are more than 2 players
			return
		}
		// send selected pack
		if err := gobby.NewBasicMessage("PACK", event.Lobby.Meta.(*Meta).Pack).SendTo(event.Client); err != nil {
			log.WithError(err).Warnf("cannot send current pack to %v", event.Client.Name)
		}
		// send user list
		if err := handlers.List.Handler(&gobby.Handle{
			Lobby:   event.Lobby,
			Client:  event.Client,
			Message: gobby.NewBasicMessage("LIST"),
		}); err != nil {
			return
		}
	})

	g.MustOn(func(event *gobby.LobbyCreate) {
		// create empty meta
		event.Lobby.Meta = &Meta{}
	})

	// respond with all available packs
	g.Handle("PACKS", &gobby.Handler{
		Roles:  gobby.RoleLeader,
		States: StateLobby,
		Handler: func(event *gobby.Handle) error {
			return event.Message.ReplyWith(event.Client, *gobby.NewBasicMessage("PACKS", Packs))
		},
	})

	g.Handle("SELECT_PACK", &gobby.Handler{
		Roles:  gobby.RoleLeader,
		States: StateLobby,
		Validation: validate.Schemes{
			validate.Number().Min(0).Max(int64(len(Packs))).As("id"),
		},
		Handler: func(event *gobby.Handle) error {
			pack := Packs[event.Number("id")]
			// update pack
			event.Lobby.Meta.(*Meta).Pack = pack
			// broadcast pack change
			return event.Lobby.Broadcast(gobby.NewBasicMessage("PACK", pack))
		},
	})

	var (
		ErrPackNotSelected = errors.New("pack not selected")
		ErrCardOutOfRange  = errors.New("card out of range")
	)
	g.Handle("SELECT_CARD", &gobby.Handler{
		States: StateLobby,
		Validation: validate.Schemes{
			validate.Number().Min(0).As("card"),
		},
		Handler: func(event *gobby.Handle) error {
			// check if pack is selected
			meta := event.Lobby.Meta.(*Meta)
			if meta.Pack == nil {
				return event.Message.ReplyWith(event.Client, *gobby.NewErrorMessage(ErrPackNotSelected))
			}
			cardIndex := event.Number("card")
			if cardIndex >= int64(len(meta.Pack.Avatars)) {
				return event.Message.ReplyWith(event.Client, *gobby.NewErrorMessage(ErrCardOutOfRange))
			}
			avatar := meta.Pack.Avatars[cardIndex]
			meta.Selected[event.Client.Name] = int(cardIndex)
			return event.Message.ReplyWith(event.Client, *gobby.NewBasicMessage(
				"SELECT_CARD", "OK", avatar.Name, avatar.Avatar,
			))
		},
	})

	if err := app.Listen(":8081"); err != nil {
		log.WithError(err).Warn("cannot serve")
	}
}
