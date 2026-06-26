import { useEffect, useState } from "react";
import { get } from "../api/client";
import type { UserProfile } from "../types/profile";

type UserProfileWithId = UserProfile & { id: string };

type State = {
  loading: boolean;
  error: string | null;
  data: UserProfileWithId[]; 
};

export function useCon() {
  const [state, setState] = useState<State>({
    loading: false,
    error: null,
    data: [],
  });

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        setState((s) => ({ ...s, loading: true, error: null }));

        //  Get all connection IDs
        const ids = await get<string[]>("/connections");

        // Fetch all user profiles in parallel
        const profiles = await Promise.all(
          ids.map(async (id) => {
            const profile = await get<UserProfile>(`/users/${id}/profile`);
            return { ...profile, id };
          })
        );

        // Update state (if still mounted)
        if (!cancelled) {
          setState({
            loading: false,
            error: null,
            data: profiles,
          });
        }
      } catch (err) {
        console.error("[useCon] error:", err);
        if (!cancelled) {
          setState({
            loading: false,
            error: "Failed to load connections",
            data: [],
          });
        }
      }
    })();

  
    return () => {
      cancelled = true;
    };
  }, []);

  return state;
}
