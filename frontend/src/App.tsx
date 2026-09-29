import { useAppSelector } from "./app/hooks";
import { AdminControls } from "./components/AdminControls";
import { ErrorToast } from "./components/ErrorToast";
import { ParticipantsPanel } from "./components/ParticipantsPanel";
import { TeamsView } from "./components/TeamsView";
import { Wheel } from "./components/Wheel";
import { selectFatal, selectRoom, selectStatus } from "./features/room/roomSlice";

export function App() {
  const status = useAppSelector(selectStatus);
  const room = useAppSelector(selectRoom);
  const fatal = useAppSelector(selectFatal);

  if (status === "error") {
    return (
      <div className="center">
        <div>
          <p>Não foi possível conectar ao Discord.</p>
          {fatal && <pre className="fatal">{fatal}</pre>}
        </div>
      </div>
    );
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
