import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { App } from "./App";
import { store } from "./store";
import { startSession } from "./services/discord";
import { savedDevName, startDevSession } from "./services/devSession";
import { devAuth, inDiscord } from "./utils/runtime";
// Fontes servidas pelo próprio app: o CSP do Discord bloqueia o Google Fonts.
import "@fontsource/pixelify-sans/400.css";
import "@fontsource/pixelify-sans/700.css";
import "@fontsource/chakra-petch/400.css";
import "@fontsource/chakra-petch/700.css";
import "./styles.css";

if (inDiscord) {
  store.dispatch(startSession());
} else if (devAuth) {
  const name = savedDevName();
  if (name) store.dispatch(startDevSession(name));
}

createRoot(document.getElementById("main")!).render(
  <Provider store={store}>
    <App />
  </Provider>,
);
