import { useAppSelector } from "./app/hooks";
import { AdminControls } from "./components/AdminControls";
import { ErrorToast } from "./components/ErrorToast";
import { ParticipantsPanel } from "./components/ParticipantsPanel";
import { TeamsView } from "./components/TeamsView";
import { Wheel } from "./components/Wheel";
import { selectRoom, selectStatus } from "./features/room/roomSlice";

export function App() {
  const status = useAppSelector(selectStatus);
  const room = useAppSelector(selectRoom);

  if (status === "error") {
    return <p className="center">Não foi possível conectar ao Discord.</p>;
  }
  if (!room) {
    return <p className="center">Conectando...</p>;
  }

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
