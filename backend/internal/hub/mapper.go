package hub

import (
	"github.com/anormais-dev/discord-application/backend/internal/dto"
	"github.com/anormais-dev/discord-application/backend/internal/roulette"
)

func toRoomResponse(r *roulette.Room) dto.RoomResponse {
	teams := make([][]dto.ParticipantResponse, len(r.Teams))
	for i, team := range r.Teams {
		teams[i] = toParticipants(team)
	}
	return dto.RoomResponse{
		Participants: toParticipants(r.Participants),
		AdminID:      r.AdminID,
		Spinners:     r.Spinners,
		Format:       toFormat(r.Format),
		Pool:         r.Pool,
		Teams:        teams,
		Phase:        string(r.Phase),
		Spin:         toSpin(r.CurrentSpin),
		Picks:        r.Picks,
	}
}

func toParticipants(ps []roulette.Participant) []dto.ParticipantResponse {
	out := make([]dto.ParticipantResponse, len(ps))
	for i, p := range ps {
		out[i] = dto.ParticipantResponse{
			ID:       p.ID,
			Name:     p.Name,
			Avatar:   p.Avatar,
			Online:   p.Online,
			JoinedAt: p.JoinedAt,
		}
	}
	return out
}

func toSpin(s *roulette.SpinState) *dto.SpinResponse {
	if s == nil {
		return nil
	}
	return &dto.SpinResponse{
		WinnerID:     s.WinnerID,
		PoolSnapshot: s.PoolSnapshot,
		TargetOffset: s.TargetOffset,
		Rotations:    s.Rotations,
		StartedAt:    s.StartedAt,
		DurationMs:   s.DurationMs,
	}
}

func toFormat(f roulette.Format) dto.Format {
	return dto.Format{Teams: f.Teams, Size: f.Size}
}

func toFormats(fs []roulette.Format) []dto.Format {
	out := make([]dto.Format, len(fs))
	for i, f := range fs {
		out[i] = toFormat(f)
	}
	return out
}
