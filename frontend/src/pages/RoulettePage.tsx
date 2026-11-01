import { CancelNotice } from "../components/CancelNotice";
import { CompactView } from "../components/CompactView";
import { DraftView } from "../components/DraftView";
import { ErrorToast } from "../components/ErrorToast";
import { Header } from "../components/Header";
import { LobbyControls } from "../components/LobbyControls";
import { MapModal } from "../components/MapModal";
import { PlayerWheel } from "../components/PlayerWheel";
import { ResultView } from "../components/ResultView";
import { SquadPanel } from "../components/SquadPanel";
import { useAppSelector } from "../store/hooks";
import { selectCompact, selectRoom, selectStatus } from "../store/selectors";

export function RoulettePage() {
  const status = useAppSelector(selectStatus);
  const compact = useAppSelector(selectCompact);
  const room = useAppSelector(selectRoom);

  if (compact) return <CompactView />;
  if (!room) return null;

  return (
    <>
      <Header />
      {status === "disconnected" && <div className="banner">Conexão perdida, reconectando...</div>}
      {room.phase === "lobby" && (
        <main className="lobby">
          <div className="stage">
            <PlayerWheel />
            <LobbyControls />
          </div>
          <SquadPanel />
        </main>
      )}
      {room.phase === "drafting" && <DraftView />}
      {room.phase === "finished" && <ResultView />}
      <MapModal />
      <CancelNotice />
      <ErrorToast />
    </>
  );
}
