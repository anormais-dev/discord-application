import { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../app/hooks";
import { errorCleared, selectError } from "../features/room/roomSlice";

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
