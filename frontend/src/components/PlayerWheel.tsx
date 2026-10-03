import { useAppSelector } from "../store/hooks";
import { selectParticipantsById, selectRoom } from "../store/selectors";
import { Wheel } from "./Wheel";

export function PlayerWheel() {
  const room = useAppSelector(selectRoom);
  const byId = useAppSelector(selectParticipantsById);
  if (!room) return null;

  const spin = room.spin;
  return (
    <Wheel
      spin={spin}
      ids={spin ? spin.poolSnapshot : room.pool}
      label={(id) => byId[id]?.name ?? "?"}
      empty={room.phase === "lobby" && room.participants.length === 0 && room.pool.length === 0}
      celebrate={room.phase !== "lobby"}
      resultPrefix="Sorteado: "
    />
  );
}
