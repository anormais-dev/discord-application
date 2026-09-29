import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { App } from "./App";
import { store } from "./app/store";
import { startSession } from "./discord";
import { savedDevName, startDevSession } from "./devSession";
import { inDiscord } from "./runtime";
import "./styles.css";

if (inDiscord) {
  store.dispatch(startSession());
} else {
  const name = savedDevName();
  if (name) store.dispatch(startDevSession(name));
}

createRoot(document.getElementById("root")!).render(
  <Provider store={store}>
    <App />
  </Provider>,
);
