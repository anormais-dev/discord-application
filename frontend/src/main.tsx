import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { App } from "./App";
import { store } from "./store";
import { startSession } from "./services/discord";
import { savedDevName, startDevSession } from "./services/devSession";
import { inDiscord } from "./utils/runtime";
import "./styles.css";

if (inDiscord) {
  store.dispatch(startSession());
} else {
  const name = savedDevName();
  if (name) store.dispatch(startDevSession(name));
}

createRoot(document.getElementById("main")!).render(
  <Provider store={store}>
    <App />
  </Provider>,
);
