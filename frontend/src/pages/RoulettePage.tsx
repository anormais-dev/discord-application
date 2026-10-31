import { AdminControls } from "../components/AdminControls";
import { CompactView } from "../components/CompactView";
import { ErrorToast } from "../components/ErrorToast";
import { MapModal } from "../components/MapModal";
import { ParticipantsPanel } from "../components/ParticipantsPanel";
import { PlayerWheel } from "../components/PlayerWheel";
import { TeamsView } from "../components/TeamsView";
import { useAppSelector } from "../store/hooks";
import { selectCompact, selectStatus } from "../store/selectors";

export function RoulettePage() {
  const status = useAppSelector(selectStatus);
  const compact = useAppSelector(selectCompact);

  if (compact) return <CompactView />;

  return (
    <main className="layout">
      {status === "disconnected" && <div className="banner">Conexão perdida, reconectando...</div>}
      <div className="stage">
        <PlayerWheel />
        <AdminControls />
      </div>
      <aside className="side">
        <ParticipantsPanel />
      </aside>
      <div className="bottom">
        <TeamsView />
      </div>
      <MapModal />
      <ErrorToast />
    </main>
  );
}
