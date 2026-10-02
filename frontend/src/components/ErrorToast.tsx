import { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import { errorCleared } from "../store/roomSlice";
import { selectError } from "../store/selectors";

const TOAST_MS = 3000;

export function ErrorToast() {
  const dispatch = useAppDispatch();
  const error = useAppSelector(selectError);

  useEffect(() => {
    if (!error) return;
    const id = setTimeout(() => dispatch(errorCleared()), TOAST_MS);
    return () => clearTimeout(id);
  }, [error, dispatch]);

  if (!error) return null;
  return (
    <div className="toast" role="alert">
      {error}
    </div>
  );
}
