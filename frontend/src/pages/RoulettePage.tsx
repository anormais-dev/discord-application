import { AdminControls } from "../components/AdminControls";
import { ErrorToast } from "../components/ErrorToast";
import { ParticipantsPanel } from "../components/ParticipantsPanel";
import { TeamsView } from "../components/TeamsView";
import { Wheel } from "../components/Wheel";
import { useAppSelector } from "../store/hooks";
import { selectStatus } from "../store/selectors";

export function RoulettePage() {
  const status = useAppSelector(selectStatus);

  return (
    <main className="layout">
      {status === "disconnected" && <div className="banner">Conexão perdida, reconectando...</div>}
      <div className="stage">
        <Wheel />
        <AdminControls />
      </div>
      <aside className="side">
        <ParticipantsPanel />
      </aside>
      <div className="bottom">
        <TeamsView />
      </div>
      <ErrorToast />
    </main>
  );
}
