import { useAppSelector } from "../store/hooks";
import { selectFatal } from "../store/selectors";

export function ErrorPage() {
  const fatal = useAppSelector(selectFatal);

  return (
    <div className="center">
      <div>
        <p>Não foi possível conectar ao Discord.</p>
        {fatal && <pre className="fatal">{fatal}</pre>}
      </div>
    </div>
  );
}
