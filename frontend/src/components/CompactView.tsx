import { useAppSelector } from "../store/hooks";
import { selectPermissions, selectRoom } from "../store/selectors";
import { Monkey } from "./Monkey";

// Mostrado com a Activity no PiP. O SDK não tem comando para sair do PiP, então só orienta.
export function CompactView() {
  const room = useAppSelector(selectRoom);
  const perms = useAppSelector(selectPermissions);
  const drafting = room?.phase === "drafting";

  return (
    <main className="compact">
      <Monkey mood={room?.spin || room?.mapSpin ? "dizzy" : "normal"} size={84} />
      <div>
        {drafting && (
          <div className="px compact-t">
            Escolha {Math.min(room.picks + 1, perms.slots)}/{perms.slots}
          </div>
        )}
        <div className="compact-s">{room?.spin ? "Roleta girando" : drafting ? "Sorteio em andamento" : "Roleta dos Anormais"}</div>
        <p className="compact-h">Use o botão do Discord para voltar à atividade</p>
      </div>
    </main>
  );
}
