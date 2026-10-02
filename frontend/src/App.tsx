import { DevLoginPage } from "./pages/DevLoginPage";
import { ErrorPage } from "./pages/ErrorPage";
import { LoadingPage } from "./pages/LoadingPage";
import { OnlyDiscordPage } from "./pages/OnlyDiscordPage";
import { RoulettePage } from "./pages/RoulettePage";
import { useAppSelector } from "./store/hooks";
import { selectRoom, selectStatus } from "./store/selectors";
import { devAuth, inDiscord } from "./utils/runtime";

export function App() {
  const status = useAppSelector(selectStatus);
  const room = useAppSelector(selectRoom);

  if (!inDiscord && !devAuth) return <OnlyDiscordPage />;
  if (status === "error") return <ErrorPage />;
  if (!inDiscord && status === "idle") return <DevLoginPage />;
  if (!room) return <LoadingPage />;
  return <RoulettePage />;
}
