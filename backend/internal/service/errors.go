package service

import "errors"

var (
	ErrNotAdmin         = errors.New("apenas o admin pode fazer isso")
	ErrCannotSpin       = errors.New("você não tem permissão para girar")
	ErrWrongPhase       = errors.New("ação não permitida neste momento")
	ErrAlreadyJoined    = errors.New("você já está participando")
	ErrNotParticipant   = errors.New("você não está participando")
	ErrSpinInProgress   = errors.New("a roleta já está girando")
	ErrNotEnoughPlayers = errors.New("participantes insuficientes para o formato")
	ErrInvalidFormat    = errors.New("formato inválido")
	ErrNoSpinInProgress = errors.New("nenhum giro em andamento")
)
