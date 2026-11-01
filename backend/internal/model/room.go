package model

type Phase string

const (
	PhaseLobby    Phase = "lobby"
	PhaseDrafting Phase = "drafting"
	PhaseFinished Phase = "finished"
)

type Room struct {
	Participants []Participant
	AdminID      string
	Spinners     map[string]bool
	Ready        map[string]bool
	AutoSpin     bool
	Format       Format
	Pool         []string
	Teams        [][]Participant
	Phase        Phase
	CurrentSpin  *Spin
	Picks        int
	MapSpin      *Spin
	Map          string
	MapOpen      bool
	Cancelled    *Cancellation
}

func NewRoom() *Room {
	return &Room{
		Participants: []Participant{},
		Spinners:     map[string]bool{},
		Ready:        map[string]bool{},
		// O giro automático já vem ligado; o admin desliga se quiser girar um a um.
		AutoSpin: true,
		Format:   DefaultFormat,
		Pool:     []string{},
		Teams:    [][]Participant{},
		Phase:    PhaseLobby,
	}
}
