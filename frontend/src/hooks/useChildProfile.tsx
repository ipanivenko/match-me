import { useEffect, useState } from "react";
import { get } from "../api/client";
import type { ChildResponse } from "../types/profile";

type State = {
  loading: boolean;
  error: string | null;
  data: ChildResponse | null;
};

export function useChildProfile() {
  const [state, setState] = useState<State>({
    loading: false,
    error: null,
    data: null,
  });

  const toYMD = (val: unknown) => {
  if (!val) return "";
  const d = new Date(val as any);
  return Number.isNaN(d.getTime()) ? "" : d.toISOString().slice(0, 10);
};

  useEffect(() => {
   
    let cancelled = false;
    (async () => {
      setState((s) => ({ ...s, loading: true, error: null }));
      try {
        const [profile] = await Promise.all([
          get<ChildResponse>("/me/child"),
        ]);

        const normalized = {
          ...profile,
          birthday: toYMD((profile as any).birthday),
        };
        setState({ loading: false, error: null, data: normalized });
      } catch (err: any) {
        if (!cancelled) {
         setState({ loading: false, error: err.message ?? "Error", data: null });
        }
      }
    })();

    return () => {
      cancelled = true;
    };
  }, []);

  return state;
}
