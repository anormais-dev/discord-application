import { AdminControls } from "./components/AdminControls";
import { ParticipantsPanel } from "./components/ParticipantsPanel";
import { TeamsView } from "./components/TeamsView";
import { Wheel } from "./components/Wheel";

export function App() {
  return (
    <main>
      <Wheel />
      <AdminControls />
      <ParticipantsPanel />
      <TeamsView />
    </main>
  );
}
