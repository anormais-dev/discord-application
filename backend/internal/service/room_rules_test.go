package service

import (
	"errors"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"

	"github.com/anormais-dev/discord-application/backend/internal/model"
)

var testRand = rand.New(rand.NewPCG(1, 2))

func newTestRoom(t *testing.T, players int, f model.Format) *model.Room {
	t.Helper()
	r := model.NewRoom()
	for i := range players {
		if err := join(r, model.Participant{ID: "u" + strconv.Itoa(i), Name: "User " + strconv.Itoa(i)}); err != nil {
			t.Fatalf("join: %v", err)
		}
		if err := setReady(r, "u"+strconv.Itoa(i), true); err != nil {
			t.Fatalf("ready: %v", err)
		}
	}
	if players > 0 {
		if err := setFormat(r, "u0", f); err != nil {
			t.Fatalf("set format: %v", err)
		}
	}
	return r
}

func spinAndCommit(t *testing.T, r *model.Room, userID string) string {
	t.Helper()
	s, err := spin(r, testRand, userID)
	if err != nil {
		t.Fatalf("spin: %v", err)
	}
	if err := commitSpin(r); err != nil {
		t.Fatalf("commit: %v", err)
	}
	return s.WinnerID
}

func expectErr(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("erro = %v, esperado %v", got, want)
	}
}

func TestFirstToJoinIsAdmin(t *testing.T) {
	r := newTestRoom(t, 3, model.Formats[0])
	if r.AdminID != "u0" {
		t.Fatalf("admin = %q, esperado u0", r.AdminID)
	}
	expectErr(t, join(r, model.Participant{ID: "u1"}), ErrAlreadyJoined)
}

func TestAdminTransfersToOldestOnline(t *testing.T) {
	r := newTestRoom(t, 3, model.Formats[0])
	if err := leave(r, "u0"); err != nil {
		t.Fatal(err)
	}
	if r.AdminID != "u1" {
		t.Fatalf("admin = %q, esperado u1", r.AdminID)
	}

	// Em drafting, desconectar mantém no time mas passa o admin.
	spinAndCommit(t, r, "u1")
	disconnect(r, "u1")
	if r.AdminID != "u2" {
		t.Fatalf("admin = %q, esperado u2", r.AdminID)
	}
	disconnect(r, "u2")
	if r.AdminID != "" || slices.ContainsFunc(r.Participants, func(p model.Participant) bool { return p.Online }) {
		t.Fatalf("sem ninguém online não deveria haver admin")
	}
	connect(r, "u1")
	if r.AdminID != "u1" {
		t.Fatalf("quem reconecta sem admin deveria virar admin, admin = %q", r.AdminID)
	}
}

func TestLobbyDisconnectRemovesParticipant(t *testing.T) {
	r := newTestRoom(t, 2, model.Formats[0])
	disconnect(r, "u1")
	if len(r.Participants) != 1 || len(r.Pool) != 1 {
		t.Fatalf("participante deveria ter saído da roda")
	}
}

func TestSpinBlockedWithoutEnoughPlayers(t *testing.T) {
	r := newTestRoom(t, 7, model.Format{Teams: 2, Size: 5})
	_, err := spin(r, testRand, "u0")
	expectErr(t, err, ErrNotEnoughPlayers)
}

func TestOnlyAdminAndSpinnersCanSpin(t *testing.T) {
	r := newTestRoom(t, 4, model.Format{Teams: 2, Size: 2})
	_, err := spin(r, testRand, "u1")
	expectErr(t, err, ErrCannotSpin)

	expectErr(t, setSpinner(r, "u1", "u2", true), ErrNotAdmin)
	if err := setSpinner(r, "u0", "u1", true); err != nil {
		t.Fatal(err)
	}
	spinAndCommit(t, r, "u1")

	expectErr(t, setFormat(r, "u1", model.Formats[0]), ErrNotAdmin)
}

func TestSpinInProgressBlocksAnotherSpin(t *testing.T) {
	r := newTestRoom(t, 2, model.Formats[0])
	if _, err := spin(r, testRand, "u0"); err != nil {
		t.Fatal(err)
	}
	_, err := spin(r, testRand, "u0")
	expectErr(t, err, ErrSpinInProgress)
}

func TestDraftAlternatesTeamsAndFinishes(t *testing.T) {
	r := newTestRoom(t, 5, model.Format{Teams: 2, Size: 2})
	var winners []string
	for range 4 {
		winners = append(winners, spinAndCommit(t, r, "u0"))
	}
	if r.Phase != model.PhaseFinished {
		t.Fatalf("fase = %s, esperado finished", r.Phase)
	}
	want := [][]string{{winners[0], winners[2]}, {winners[1], winners[3]}}
	for ti, team := range r.Teams {
		for pi, p := range team {
			if p.ID != want[ti][pi] {
				t.Fatalf("time %d posição %d = %s, esperado %s", ti, pi, p.ID, want[ti][pi])
			}
		}
	}
	if len(r.Pool) != 1 {
		t.Fatalf("deveria sobrar 1 fora dos times, sobrou %d", len(r.Pool))
	}
	_, err := spin(r, testRand, "u0")
	expectErr(t, err, ErrWrongPhase)
}

func TestWinnerLeavesPool(t *testing.T) {
	r := newTestRoom(t, 5, model.Format{Teams: 2, Size: 2})
	seen := map[string]bool{}
	for range 4 {
		w := spinAndCommit(t, r, "u0")
		if seen[w] {
			t.Fatalf("%s sorteado duas vezes", w)
		}
		seen[w] = true
	}
}

func TestLeaveRules(t *testing.T) {
	r := newTestRoom(t, 3, model.Formats[0])
	spinAndCommit(t, r, "u0")
	expectErr(t, leave(r, "u1"), ErrWrongPhase)
	expectErr(t, join(r, model.Participant{ID: "novo"}), ErrWrongPhase)

	spinAndCommit(t, r, "u0")
	if err := leave(r, "u1"); err != nil {
		t.Fatalf("depois de finalizada deveria poder sair: %v", err)
	}
	if len(r.Teams[0])+len(r.Teams[1]) != 2 {
		t.Fatalf("resultado dos times não deveria mudar quando alguém sai")
	}
}

func TestResetClearsEverything(t *testing.T) {
	r := newTestRoom(t, 3, model.Formats[0])
	expectErr(t, reset(r, "u0"), ErrWrongPhase)
	spinAndCommit(t, r, "u0")
	spinAndCommit(t, r, "u0")
	expectErr(t, reset(r, "u1"), ErrNotAdmin)
	if err := reset(r, "u0"); err != nil {
		t.Fatal(err)
	}
	if len(r.Participants) != 0 || r.AdminID != "" || r.Phase != model.PhaseLobby || len(r.Teams) != 0 {
		t.Fatalf("reset não limpou o estado: %+v", r)
	}
	if err := join(r, model.Participant{ID: "u1"}); err != nil || r.AdminID != "u1" {
		t.Fatalf("próximo a entrar deveria virar admin")
	}
}

func TestInvalidFormat(t *testing.T) {
	r := newTestRoom(t, 1, model.Formats[0])
	expectErr(t, setFormat(r, "u0", model.Format{Teams: 3, Size: 7}), ErrInvalidFormat)
}

func TestSingleTeamFormat(t *testing.T) {
	r := newTestRoom(t, 4, model.Format{Teams: 1, Size: 3})
	for range 3 {
		spinAndCommit(t, r, "u0")
	}
	if r.Phase != model.PhaseFinished || len(r.Teams) != 1 || len(r.Teams[0]) != 3 {
		t.Fatalf("esperado 1 time com 3, veio fase %s times %+v", r.Phase, r.Teams)
	}
	if len(r.Pool) != 1 {
		t.Fatalf("deveria sobrar 1 fora do time, sobrou %d", len(r.Pool))
	}
}

func TestLastRemainingJoinsWithoutSpin(t *testing.T) {
	r := newTestRoom(t, 4, model.Format{Teams: 2, Size: 2})
	for range 3 {
		spinAndCommit(t, r, "u0")
	}
	if r.Phase != model.PhaseFinished || len(r.Pool) != 0 {
		t.Fatalf("o último da roda deveria entrar sozinho, fase %s pool %v", r.Phase, r.Pool)
	}
	if len(r.Teams[0]) != 2 || len(r.Teams[1]) != 2 {
		t.Fatalf("times incompletos: %+v", r.Teams)
	}
}

func TestLastRemainingWaitsWhenMorePlayersThanSlots(t *testing.T) {
	r := newTestRoom(t, 3, model.Formats[0])
	spinAndCommit(t, r, "u0")
	if r.Phase != model.PhaseDrafting || len(r.Pool) != 2 {
		t.Fatalf("com 2 na roda ainda precisa girar, fase %s pool %v", r.Phase, r.Pool)
	}
}

func TestSpinNeedsEveryoneReady(t *testing.T) {
	r := newTestRoom(t, 2, model.Formats[0])
	if err := setReady(r, "u1", false); err != nil {
		t.Fatal(err)
	}
	_, err := spin(r, testRand, "u0")
	expectErr(t, err, ErrNotAllReady)

	if err := join(r, model.Participant{ID: "novo"}); err != nil {
		t.Fatal(err)
	}
	setReady(r, "u1", true)
	_, err = spin(r, testRand, "u0")
	expectErr(t, err, ErrNotAllReady)

	if err := leave(r, "novo"); err != nil {
		t.Fatal(err)
	}
	spinAndCommit(t, r, "u0")
}

func TestReadyRules(t *testing.T) {
	r := newTestRoom(t, 3, model.Formats[0])
	expectErr(t, setReady(r, "ninguem", true), ErrNotParticipant)
	disconnect(r, "u2")
	if r.Ready["u2"] {
		t.Fatalf("quem sai da sala não deveria continuar pronto")
	}
	spinAndCommit(t, r, "u0")
	expectErr(t, setReady(r, "u0", false), ErrWrongPhase)
}

func TestAutoSpinRules(t *testing.T) {
	r := newTestRoom(t, 5, model.Format{Teams: 2, Size: 2})
	expectErr(t, setAutoSpin(r, "u1", true), ErrCannotSpin)
	if err := setAutoSpin(r, "u0", true); err != nil {
		t.Fatal(err)
	}
	if autoSpin(r, testRand) {
		t.Fatalf("o primeiro giro precisa ser manual")
	}
	spinAndCommit(t, r, "u0")
	if !autoSpin(r, testRand) || r.CurrentSpin == nil {
		t.Fatalf("depois do primeiro giro deveria girar sozinho")
	}
	if autoSpin(r, testRand) {
		t.Fatalf("não deveria girar com outro giro em andamento")
	}
	commitSpin(r)

	setAutoSpin(r, "u0", false)
	if autoSpin(r, testRand) {
		t.Fatalf("com o automático desligado não deveria girar")
	}
}

func TestSpinMapOnlyAdmin(t *testing.T) {
	r := newTestRoom(t, 2, model.Format{Teams: 1, Size: 2})
	if _, err := spinMap(r, testRand, "u1", []string{"Ascent"}); !errors.Is(err, ErrNotAdmin) {
		t.Fatalf("erro = %v, esperado ErrNotAdmin", err)
	}
	if r.MapSpin != nil {
		t.Fatalf("não-admin não deveria iniciar giro de mapa")
	}
}

func TestSpinMapWorksInAnyPhaseWithoutTouchingTeams(t *testing.T) {
	maps := []string{"Ascent", "Bind", "Haven"}
	r := newTestRoom(t, 2, model.Format{Teams: 1, Size: 2})

	for _, phase := range []model.Phase{model.PhaseLobby, model.PhaseDrafting, model.PhaseFinished} {
		if phase == model.PhaseDrafting {
			if _, err := spin(r, testRand, "u0"); err != nil {
				t.Fatalf("spin: %v", err)
			}
		}
		if phase == model.PhaseFinished {
			commitSpin(r)
		}
		if r.Phase != phase {
			t.Fatalf("fase = %s, esperado %s", r.Phase, phase)
		}
		teamSpin := r.CurrentSpin
		s, err := spinMap(r, testRand, "u0", maps)
		if err != nil {
			t.Fatalf("%s: spinMap: %v", phase, err)
		}
		if _, err := spinMap(r, testRand, "u0", maps); !errors.Is(err, ErrSpinInProgress) {
			t.Fatalf("%s: erro = %v, esperado ErrSpinInProgress", phase, err)
		}
		if err := commitMapSpin(r); err != nil {
			t.Fatalf("%s: commitMapSpin: %v", phase, err)
		}
		if r.Map != s.WinnerID || !slices.Contains(maps, r.Map) || r.MapSpin != nil {
			t.Fatalf("%s: mapa = %q, giro = %+v", phase, r.Map, r.MapSpin)
		}
		if r.CurrentSpin != teamSpin {
			t.Fatalf("%s: giro de mapa não deveria mexer no giro dos times", phase)
		}
	}
}

func TestSpinMapWithoutMaps(t *testing.T) {
	r := newTestRoom(t, 1, model.Format{Teams: 1, Size: 1})
	if _, err := spinMap(r, testRand, "u0", nil); !errors.Is(err, ErrMapsUnavailable) {
		t.Fatalf("erro = %v, esperado ErrMapsUnavailable", err)
	}
}
