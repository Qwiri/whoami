package main

import (
	"errors"
	"fmt"
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/handlers"
	"github.com/Qwiri/gobby/pkg/validate"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/qwiri/whoami/pkg/meta"
)

const (
	StateLobby gobby.State = 1 << iota
	StateSelectCharacter
	StateGame
	StateWinningScreen
)

type Meta struct {
	PackIndex int
	Selected  map[string]int // Client Name -> Selected
}

var (
	ErrLobbyFull         = errors.New("the lobby is full")
	ErrNotLobby          = errors.New("the game is not in lobby mode")
	ErrCardOutOfRange    = errors.New("card out of range")
	ErrPlayerRequirement = errors.New("player requirement not met")
)

func main() {
	log.Infof("Packs: %v", Packs)
	log.Infof("Starting backend for whoami %s (%s@%s)",
		meta.Version, meta.GitCommit, meta.GitBranch)
	app := fiber.New()

	g := gobby.New(app)
	g.AppVersion = fmt.Sprintf("%s:%s@%s",
		meta.Version, meta.GitCommit, meta.GitBranch) // send backend version

	g.HandleAll(gobby.Handlers{
		"CHAT": handlers.Chat,
		"LIST": handlers.List,
	})

	g.MustOn(func(event *gobby.LobbyCreate) {
		event.Lobby.Meta = new(Meta)
		event.Lobby.State = StateLobby
	})

	g.MustOn(func(event *gobby.Join) {
		// do not allow more than 2 players to a game
		if len(event.Lobby.Clients) >= 2 {
			_ = gobby.NewErrorMessage(ErrLobbyFull).SendTo(event.Client)
			event.Cancel()
			return
		}
		// do not allow joining if not in lobby
		if event.Lobby.State != StateLobby {
			_ = gobby.NewErrorMessage(ErrNotLobby).SendTo(event.Client)
			event.Cancel()
			return
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

	g.MustOn(func(event *gobby.Leave) {
		// reset to lobby state if any client disconnects
		event.Lobby.ChangeState(StateLobby)
	})

	// lifecycle events
	g.Handle("START", &gobby.Handler{
		States: StateLobby | StateWinningScreen,
		Handler: func(event *gobby.Handle) error {
			// require 2 players to start
			if len(event.Lobby.Clients) != 2 {
				return event.Message.ReplyWith(event.Client, *gobby.NewErrorMessage(ErrPlayerRequirement))
			}
			// send available card selections
			m := event.Lobby.Meta.(*Meta)
			pack := Packs[m.PackIndex]
			// send available characters
			event.Lobby.BroadcastForce(gobby.NewBasicMessageWith("AVAILABLE_CHARACTERS",
				pack.Avatars...))
			event.Lobby.ChangeState(StateSelectCharacter)
			return nil
		},
	})

	g.Handle("SELECT_CHARACTER", &gobby.Handler{
		States: StateSelectCharacter,
		Validation: validate.Schemes{
			validate.Number().Min(0).As("char"),
		},
		Handler: func(event *gobby.Handle) error {
			// check if pack is selected
			m := event.Lobby.Meta.(*Meta)
			pack := Packs[m.PackIndex]

			charIndex := event.Number("char")
			if charIndex >= int64(len(pack.Avatars)) {
				return event.Message.ReplyWith(event.Client, *gobby.NewErrorMessage(ErrCardOutOfRange))
			}

			avatar := pack.Avatars[charIndex]
			m.Selected[event.Client.Name] = int(charIndex)

			_ = event.Message.ReplyWith(event.Client, *gobby.NewBasicMessage(
				"SELECT_CARD", "OK", avatar.Name, avatar.Avatar,
			))

			// check if all clients selected a card
			var missing bool
			for _, c := range event.Lobby.Clients {
				if _, ok := m.Selected[c.Name]; !ok {
					missing = true
					break
				}
			}
			// if all players selected a card, move on
			if !missing {
				event.Lobby.ChangeState(StateGame)
			}

			return nil
		},
	})

	g.Handle("GUESS", &gobby.Handler{
		States: StateGame,
		Validation: validate.Schemes{
			validate.Number().Min(0).As("char"),
		},
		Handler: func(event *gobby.Handle) error {
			m := event.Lobby.Meta.(*Meta)

			// get selection of other player
			var selected int
			for _, c := range event.Lobby.Clients {
				if c.Name != event.Client.Name {
					selected = m.Selected[c.Name]
					break
				}
			}

			if selected != int(event.Number("char")) {
				// TODO: wrong guess
			} else {
				// TODO: right guess
			}

			return nil
		},
	})

	// respond with all available packs
	g.Handle("PACKS", &gobby.Handler{
		States: StateLobby,
		Handler: func(event *gobby.Handle) error {
			return event.Message.ReplyWith(event.Client, *gobby.NewBasicMessage("PACKS", Packs))
		},
	})

	g.Handle("SELECT_PACK", &gobby.Handler{
		States: StateLobby,
		Validation: validate.Schemes{
			validate.Number().Min(0).Max(int64(len(Packs))).As("id"),
		},
		Handler: func(event *gobby.Handle) error {
			pack := int(event.Number("id"))
			event.Lobby.Meta.(*Meta).PackIndex = pack
			return event.Message.ReplyWith(event.Client, *gobby.NewBasicMessage("PACK", pack))
		},
	})

	if err := app.Listen(":8081"); err != nil {
		log.WithError(err).Warn("cannot serve")
	}
}
