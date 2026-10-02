import { DevLoginPage } from "./pages/DevLoginPage";
import { ErrorPage } from "./pages/ErrorPage";
import { LoadingPage } from "./pages/LoadingPage";
import { RoulettePage } from "./pages/RoulettePage";
import { useAppSelector } from "./store/hooks";
import { selectRoom, selectStatus } from "./store/selectors";
import { inDiscord } from "./utils/runtime";

export function App() {
  const status = useAppSelector(selectStatus);
  const room = useAppSelector(selectRoom);

  if (status === "error") return <ErrorPage />;
  if (!inDiscord && status === "idle") return <DevLoginPage />;
  if (!room) return <LoadingPage />;
  return <RoulettePage />;
}
