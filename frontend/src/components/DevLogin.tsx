import { useState, type FormEvent } from "react";
import { useAppDispatch } from "../store/hooks";
import { startDevSession } from "../services/devSession";

export function DevLogin() {
  const dispatch = useAppDispatch();
  const [name, setName] = useState("");

  const submit = (e: FormEvent) => {
    e.preventDefault();
    if (name.trim()) dispatch(startDevSession(name.trim()));
  };

  return (
    <div className="center">
      <form className="panel dev-login" onSubmit={submit}>
        <h2>Modo de teste (fora do Discord)</h2>
        <p className="muted">Escolha um nome. Abra outras abas com nomes diferentes para simular mais pessoas.</p>
        <input autoFocus value={name} maxLength={32} onChange={(e) => setName(e.target.value)} placeholder="Seu nome" />
        <button className="primary" type="submit" disabled={!name.trim()}>
          Entrar
        </button>
      </form>
    </div>
  );
}
