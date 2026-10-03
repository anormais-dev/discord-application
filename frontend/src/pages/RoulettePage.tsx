import { AdminControls } from "../components/AdminControls";
import { ErrorToast } from "../components/ErrorToast";
import { MapPanel } from "../components/MapPanel";
import { ParticipantsPanel } from "../components/ParticipantsPanel";
import { PlayerWheel } from "../components/PlayerWheel";
import { TeamsView } from "../components/TeamsView";
import { useAppSelector } from "../store/hooks";
import { selectStatus } from "../store/selectors";

export function RoulettePage() {
  const status = useAppSelector(selectStatus);

  return (
    <main className="layout">
      {status === "disconnected" && <div className="banner">Conexão perdida, reconectando...</div>}
      <div className="stage">
        <PlayerWheel />
        <AdminControls />
      </div>
      <aside className="side">
        <ParticipantsPanel />
        <MapPanel />
      </aside>
      <div className="bottom">
        <TeamsView />
      </div>
      <ErrorToast />
    </main>
  );
}
