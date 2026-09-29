package roulette

import (
	"errors"
	"math/rand/v2"
	"strconv"
	"testing"
)

func newTestRoom(t *testing.T, players int, f Format) *Room {
	t.Helper()
	r := NewRoomWithRand(rand.New(rand.NewPCG(1, 2)))
	for i := range players {
		if err := r.Join(Participant{ID: "u" + strconv.Itoa(i), Name: "User " + strconv.Itoa(i)}); err != nil {
			t.Fatalf("join: %v", err)
		}
	}
	if players > 0 {
		if err := r.SetFormat("u0", f); err != nil {
			t.Fatalf("set format: %v", err)
		}
	}
	return r
}

func spinAndCommit(t *testing.T, r *Room, userID string) string {
	t.Helper()
	s, err := r.Spin(userID)
	if err != nil {
		t.Fatalf("spin: %v", err)
	}
	if err := r.CommitSpin(); err != nil {
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
	r := newTestRoom(t, 3, Formats[0])
	if r.AdminID != "u0" {
		t.Fatalf("admin = %q, esperado u0", r.AdminID)
	}
	expectErr(t, r.Join(Participant{ID: "u1"}), ErrAlreadyJoined)
}

func TestAdminTransfersToOldestOnline(t *testing.T) {
	r := newTestRoom(t, 3, Formats[0])
	if err := r.Leave("u0"); err != nil {
		t.Fatal(err)
	}
	if r.AdminID != "u1" {
		t.Fatalf("admin = %q, esperado u1", r.AdminID)
	}

	// Em drafting, desconectar mantém no time mas passa o admin.
	spinAndCommit(t, r, "u1")
	r.Disconnect("u1")
	if r.AdminID != "u2" {
		t.Fatalf("admin = %q, esperado u2", r.AdminID)
	}
	r.Disconnect("u2")
	if r.AdminID != "" || r.HasOnline() {
		t.Fatalf("sem ninguém online não deveria haver admin")
	}
	r.Connect("u1")
	if r.AdminID != "u1" {
		t.Fatalf("quem reconecta sem admin deveria virar admin, admin = %q", r.AdminID)
	}
}

func TestLobbyDisconnectRemovesParticipant(t *testing.T) {
	r := newTestRoom(t, 2, Formats[0])
	r.Disconnect("u1")
	if len(r.Participants) != 1 || len(r.Pool) != 1 {
		t.Fatalf("participante deveria ter saído da roda")
	}
}

func TestSpinBlockedWithoutEnoughPlayers(t *testing.T) {
	r := newTestRoom(t, 7, Format{Teams: 2, Size: 5})
	_, err := r.Spin("u0")
	expectErr(t, err, ErrNotEnoughPlayers)
}

func TestOnlyAdminAndSpinnersCanSpin(t *testing.T) {
	r := newTestRoom(t, 4, Format{Teams: 2, Size: 2})
	_, err := r.Spin("u1")
	expectErr(t, err, ErrCannotSpin)

	expectErr(t, r.SetSpinner("u1", "u2", true), ErrNotAdmin)
	if err := r.SetSpinner("u0", "u1", true); err != nil {
		t.Fatal(err)
	}
	spinAndCommit(t, r, "u1")

	expectErr(t, r.SetFormat("u1", Formats[0]), ErrNotAdmin)
}

func TestSpinInProgressBlocksAnotherSpin(t *testing.T) {
	r := newTestRoom(t, 2, Formats[0])
	if _, err := r.Spin("u0"); err != nil {
		t.Fatal(err)
	}
	_, err := r.Spin("u0")
	expectErr(t, err, ErrSpinInProgress)
}

func TestDraftAlternatesTeamsAndFinishes(t *testing.T) {
	r := newTestRoom(t, 5, Format{Teams: 2, Size: 2})
	var winners []string
	for range 4 {
		winners = append(winners, spinAndCommit(t, r, "u0"))
	}
	if r.Phase != PhaseFinished {
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
	_, err := r.Spin("u0")
	expectErr(t, err, ErrWrongPhase)
}

func TestWinnerLeavesPool(t *testing.T) {
	r := newTestRoom(t, 4, Format{Teams: 2, Size: 2})
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
	r := newTestRoom(t, 2, Formats[0])
	spinAndCommit(t, r, "u0")
	expectErr(t, r.Leave("u1"), ErrWrongPhase)
	expectErr(t, r.Join(Participant{ID: "novo"}), ErrWrongPhase)

	spinAndCommit(t, r, "u0")
	if err := r.Leave("u1"); err != nil {
		t.Fatalf("depois de finalizada deveria poder sair: %v", err)
	}
	if len(r.Teams[0])+len(r.Teams[1]) != 2 {
		t.Fatalf("resultado dos times não deveria mudar quando alguém sai")
	}
}

func TestResetClearsEverything(t *testing.T) {
	r := newTestRoom(t, 2, Formats[0])
	expectErr(t, r.Reset("u0"), ErrWrongPhase)
	spinAndCommit(t, r, "u0")
	spinAndCommit(t, r, "u0")
	expectErr(t, r.Reset("u1"), ErrNotAdmin)
	if err := r.Reset("u0"); err != nil {
		t.Fatal(err)
	}
	if len(r.Participants) != 0 || r.AdminID != "" || r.Phase != PhaseLobby || len(r.Teams) != 0 {
		t.Fatalf("reset não limpou o estado: %+v", r)
	}
	if err := r.Join(Participant{ID: "u1"}); err != nil || r.AdminID != "u1" {
		t.Fatalf("próximo a entrar deveria virar admin")
	}
}

func TestInvalidFormat(t *testing.T) {
	r := newTestRoom(t, 1, Formats[0])
	expectErr(t, r.SetFormat("u0", Format{Teams: 3, Size: 7}), ErrInvalidFormat)
}
