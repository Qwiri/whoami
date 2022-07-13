package main

import (
	"errors"
	"fmt"
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/handlers"
	"github.com/Qwiri/gobby/pkg/validate"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/qwiri/whoami/pkg/meta"
	"strconv"
)

const (
	StateLobby gobby.State = 1 << iota
	StateSelectCharacter
	StateGame
	StateWinningScreen
)

const (
	DefaultLives = 3
)

type Meta struct {
	PackIndex int
	Selected  map[string]int // Client Name -> Selected
	Lives     map[string]int // Client Name -> Lives (3)
}

var (
	ErrLobbyFull         = errors.New("the lobby is full")
	ErrNotLobby          = errors.New("the game is not in lobby mode")
	ErrCardOutOfRange    = errors.New("card out of range")
	ErrPackOutOfRange    = errors.New("pack out of range")
	ErrPlayerRequirement = errors.New("player requirement not met")
	ErrAlreadyLobby      = errors.New("already in lobby")
)

func (m *Meta) resetCurrentGame() *Meta {
	m.Selected = make(map[string]int)
	m.Lives = make(map[string]int)
	return m
}

func (m *Meta) Win(lobby *gobby.Lobby, who *gobby.Client, reason string) {
	// send win message to all players
	lobby.BroadcastForce(gobby.NewBasicMessage("WINNER", who.Name, m.Selected[who.Name], reason))
	// set state to end screen
	lobby.ChangeState(StateWinningScreen)
	// reset game
	m.resetCurrentGame()
}

func (m *Meta) DecreaseLife(who *gobby.Client, many int) (bool, error) {
	newLives := m.Lives[who.Name] - many
	m.Lives[who.Name] = newLives

	// send new live count to player
	err := gobby.NewBasicMessageWith[int]("LIVES", newLives, DefaultLives).SendTo(who)
	if newLives > 0 {
		return false, err
	}
	return true, nil
}

func newMeta() *Meta {
	return new(Meta).resetCurrentGame()
}

type peekResponse struct {
	ID      gobby.LobbyID
	State   gobby.State
	Clients []string
}

func main() {
	log.Infof("Packs: %v", Packs)
	log.Infof("Starting backend for whoami %s (%s@%s)",
		meta.Version, meta.GitCommit, meta.GitBranch)

	app := fiber.New()
	app.Use(cors.New())

	g := gobby.New(app)
	g.AppVersion = fmt.Sprintf("%s:%s@%s",
		meta.Version, meta.GitCommit, meta.GitBranch) // send backend version

	app.Get(g.Prefix+"peek/:lobby_id", func(ctx *fiber.Ctx) error {
		lobby, ok := g.Lobbies[gobby.LobbyID(utils.CopyString(ctx.Params("lobby_id")))]
		if !ok {
			return fiber.NewError(fiber.StatusNotFound, "lobby not found")
		}
		var clients = make([]string, 0, len(lobby.Clients))
		for _, c := range lobby.Clients {
			clients = append(clients, c.Name)
		}
		return ctx.JSON(peekResponse{
			ID:      lobby.ID,
			State:   lobby.State,
			Clients: clients,
		})
	})

	// default handlers: CHAT (for in-game chat) and LIST (to list all players in a lobby)
	g.HandleAll(gobby.Handlers{
		"CHAT": handlers.Chat,
		"LIST": handlers.List,
	})

	g.MustOn(
		// init lobby after a new lobby has been created
		func(event *gobby.LobbyCreateEvent) {
			event.Lobby.Meta = newMeta()
			event.Lobby.ChangeState(StateLobby)
		},
		// reset current lobby to lobby mode if any client disconnects
		func(event *gobby.LeaveEvent) {
			// resetCurrentGame player selection on player leave and change state to lobby
			event.Lobby.Meta.(*Meta).resetCurrentGame()
			event.Lobby.ChangeState(StateLobby)
		},
		// check if the lobby can be joined and cancel the join process if not
		func(event *gobby.PreJoinEvent) {
			log.Infof("PreJoinEvent: %s", event.Client.Name)
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
		},
		// send some meta to the player after joining a lobby, like
		// - available pack
		// - selected pack
		func(event *gobby.PostJoinEvent) {
			log.Infof("PostJoinEvent: %s", event.Client.Name)
			// send available packs
			_ = gobby.NewBasicMessageWith[[]*CardPack]("PACKS", Packs).
				SendTo(event.Client)

			// send selected pack
			_ = gobby.NewBasicMessageWith[int]("SELECTED_PACK_CHANGED", event.Lobby.Meta.(*Meta).PackIndex).
				SendTo(event.Client)
		},
	)

	// lifecycle events
	g.Handle("START", &gobby.Handler{
		States: StateLobby,
		Handler: func(event *gobby.Handle) error {
			// check that 2 players are in the lobby
			if len(event.Lobby.Clients) != 2 {
				return event.Message.ReplyWith(event.Client, *gobby.NewErrorMessage(ErrPlayerRequirement))
			}

			m := event.Lobby.Meta.(*Meta)
			m.resetCurrentGame()

			// send game meta to clients:
			// lives
			for _, c := range event.Lobby.Clients {
				m.Lives[c.Name] = DefaultLives
				_, _ = m.DecreaseLife(c, 0) // send reset lives
			}
			// selectable characters
			pack := Packs[m.PackIndex]
			event.Lobby.BroadcastForce(gobby.NewBasicMessageWith("AVAILABLE_CHARACTERS", pack.Avatars...))

			// change state to select character
			event.Lobby.ChangeState(StateSelectCharacter)
			return nil
		},
	})

	g.Handle("CANCEL", &gobby.Handler{
		Handler: func(event *gobby.Handle) error {
			if event.Lobby.State == StateLobby {
				return event.Message.ReplyError(event.Client.Socket, ErrAlreadyLobby)
			}
			event.Lobby.Meta.(*Meta).resetCurrentGame()
			event.Lobby.ChangeState(StateLobby)
			return event.Message.ReplyWith(event.Client, *gobby.NewBasicMessage("CANCEL", "ok"))
		},
	})

	g.Handle("SELECT_CHARACTER", &gobby.Handler{
		States: StateSelectCharacter,
		Validation: validate.Schemes{
			validate.String().As("char"),
		},
		Handler: func(event *gobby.Handle) error {
			char, err := strconv.Atoi(event.String("char"))
			if err != nil {
				return event.Message.ReplyError(event.Client.Socket, err)
			}

			// check if pack is selected
			m := event.Lobby.Meta.(*Meta)
			pack := Packs[m.PackIndex]

			if char < 0 || char >= len(pack.Avatars) {
				return event.Message.ReplyWith(event.Client, *gobby.NewErrorMessage(ErrCardOutOfRange))
			}

			avatar := pack.Avatars[char]
			m.Selected[event.Client.Name] = char

			_ = event.Message.ReplyWith(event.Client, *gobby.NewBasicMessage(
				"SELECT_CHARACTER", "OK", avatar.Name, avatar.Avatar,
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
			validate.String().As("char"),
		},
		Handler: func(event *gobby.Handle) error {
			char, err := strconv.Atoi(event.String("char"))
			if err != nil {
				return event.Message.ReplyError(event.Client.Socket, err)
			}

			m := event.Lobby.Meta.(*Meta)

			// get selection of other player
			var other *gobby.Client
			for _, c := range event.Lobby.Clients {
				if c.Name != event.Client.Name {
					other = c
					break
				}
			}

			if m.Selected[other.Name] != char {
				err = event.Message.ReplyWith(event.Client, *gobby.NewBasicMessage("GUESS", "wrong"))
				if died, err := m.DecreaseLife(event.Client, 1); err != nil {
					return err
				} else if died {
					m.Win(event.Lobby, other, "LIVES")
				}
				return err
			}

			err = event.Message.ReplyWith(event.Client, *gobby.NewBasicMessage("GUESS", "correct"))

			// current player wins because the player guessed right
			m.Win(event.Lobby, event.Client, "GUESS")

			// resetCurrentGame selected char
			m.resetCurrentGame()

			// broadcast round end
			return err
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
			validate.String().As("id"),
		},
		Handler: func(event *gobby.Handle) error {
			// convert id string to number
			packStr := event.String("id")
			if pack, err := strconv.Atoi(packStr); err != nil {
				return event.Message.ReplyError(event.Client.Socket, err)
			} else {
				if pack < 0 || pack >= len(Packs) {
					return event.Message.ReplyError(event.Client.Socket, ErrPackOutOfRange)
				}
				event.Lobby.Meta.(*Meta).PackIndex = pack
				event.Lobby.BroadcastForce(gobby.NewBasicMessageWith[int]("SELECTED_PACK_CHANGED", pack))
				return event.Message.ReplyWith(event.Client, *gobby.NewBasicMessage("SELECT_PACK", pack))
			}
		},
	})

	if err := app.Listen(":8081"); err != nil {
		log.WithError(err).Warn("cannot serve")
	}
}
