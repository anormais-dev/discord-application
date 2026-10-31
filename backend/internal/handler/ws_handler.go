package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"github.com/anormais-dev/discord-application/backend/internal/dto"
	"github.com/anormais-dev/discord-application/backend/internal/model"
	"github.com/anormais-dev/discord-application/backend/internal/service"
	"github.com/anormais-dev/discord-application/backend/pkg/discord"
)

const (
	authTimeout  = 10 * time.Second
	writeTimeout = 5 * time.Second
	// Proxies no caminho (Discord, Cloudflare, Traefik) derrubam WebSocket ocioso;
	// o ping mantém o tráfego e detecta conexão que morreu sem fechar.
	pingInterval = 20 * time.Second
	pingTimeout  = 10 * time.Second
)

type WSHandler struct {
	auth  *service.AuthService
	rooms *service.RoomService
	maps  *service.MapService
}

type client struct {
	user *discord.User
	send chan []byte
}

func NewWSHandler(auth *service.AuthService, rooms *service.RoomService, maps *service.MapService) *WSHandler {
	return &WSHandler{auth: auth, rooms: rooms, maps: maps}
}

// ServeWS trata GET /api/ws?instance=<instanceId>&guild=<guildId>&channel=<channelId>.
func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if q.Get("instance") == "" {
		http.Error(w, "instance obrigatório", http.StatusBadRequest)
		return
	}
	instance := roomKey(q.Get("instance"), q.Get("guild"), q.Get("channel"))
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*.discordsays.com"},
	})
	if err != nil {
		return
	}
	defer conn.CloseNow()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user, err := h.authenticate(ctx, conn)
	if err != nil {
		log.Printf("ws auth: %v", err)
		conn.Close(websocket.StatusPolicyViolation, "autenticação falhou")
		return
	}

	c := &client{user: user, send: make(chan []byte, 32)}
	h.rooms.Attach(instance, c)
	defer h.rooms.Detach(instance, c)

	go writeLoop(ctx, cancel, conn, c)
	go keepAlive(ctx, cancel, conn)

	for {
		var msg dto.WSRequest
		if err := wsjson.Read(ctx, conn, &msg); err != nil {
			return
		}
		if err := h.handle(ctx, instance, c, msg); err != nil {
			c.enqueue(mustMarshal(dto.ErrorMessage{Type: "error", Message: err.Error()}))
		}
	}
}

// roomKey prende a sala ao canal de voz: cada call tem a sua, e reabrir a Activity
// na mesma call volta para a mesma sala. Sem canal (dev), fica a instance.
func roomKey(instance, guild, channel string) string {
	if channel == "" {
		return instance
	}
	return "channel:" + guild + ":" + channel
}

func (h *WSHandler) authenticate(ctx context.Context, conn *websocket.Conn) (*discord.User, error) {
	ctx, cancel := context.WithTimeout(ctx, authTimeout)
	defer cancel()
	var msg dto.WSRequest
	if err := wsjson.Read(ctx, conn, &msg); err != nil {
		return nil, err
	}
	if msg.Type != "auth" || msg.AccessToken == "" {
		return nil, errors.New("primeira mensagem precisa ser auth")
	}
	return h.auth.CurrentUser(ctx, msg.AccessToken, msg.GuildID)
}

func (h *WSHandler) handle(ctx context.Context, instance string, c *client, msg dto.WSRequest) error {
	uid := c.user.ID
	switch msg.Type {
	case "join":
		return h.rooms.Join(instance, model.Participant{
			ID:     uid,
			Name:   service.CleanName(c.user.Nick, c.user.GlobalName, c.user.Username),
			Avatar: c.user.AvatarURL(),
		})
	case "leave":
		return h.rooms.Leave(instance, uid)
	case "set_format":
		return h.rooms.SetFormat(instance, uid, model.Format{Teams: msg.Format.Teams, Size: msg.Format.Size})
	case "set_spinner":
		return h.rooms.SetSpinner(instance, uid, msg.UserID, msg.Enabled)
	case "set_ready":
		return h.rooms.SetReady(instance, uid, msg.Enabled)
	case "set_auto_spin":
		return h.rooms.SetAutoSpin(instance, uid, msg.Enabled)
	case "spin":
		return h.rooms.Spin(instance, uid)
	case "set_map_open":
		return h.rooms.SetMapOpen(instance, uid, msg.Enabled)
	case "spin_map":
		maps, err := h.maps.CompetitiveMaps(ctx)
		if err != nil {
			return err
		}
		return h.rooms.SpinMap(instance, uid, maps)
	case "reset":
		return h.rooms.Reset(instance, uid)
	default:
		return errors.New("mensagem desconhecida")
	}
}

func (c *client) UserID() string {
	return c.user.ID
}

func (c *client) SendState(room *model.Room) {
	c.enqueue(mustMarshal(dto.StateMessage{
		Type:    "state",
		Room:    toRoomResponse(room),
		You:     c.user.ID,
		Formats: toFormats(model.Formats),
	}))
}

// enqueue nunca bloqueia; um cliente lento demais perde mensagens e recebe o próximo estado completo.
func (c *client) enqueue(b []byte) {
	select {
	case c.send <- b:
	default:
	}
}

func writeLoop(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, c *client) {
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case b := <-c.send:
			wctx, wcancel := context.WithTimeout(ctx, writeTimeout)
			err := conn.Write(wctx, websocket.MessageText, b)
			wcancel()
			if err != nil {
				return
			}
		}
	}
}

// keepAlive depende do loop de leitura rodando, que é quem recebe o pong.
func keepAlive(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pctx, pcancel := context.WithTimeout(ctx, pingTimeout)
			err := conn.Ping(pctx)
			pcancel()
			if err != nil {
				cancel()
				return
			}
		}
	}
}

func mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
