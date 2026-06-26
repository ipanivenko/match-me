import { useEffect, useState } from "react";
import { get } from "../api/client";
import type { CombinedUser, UserPhoto, UserProfile } from "../types/profile";

type CombinedUserWithId = CombinedUser & { id: string };

type State = {
  loading: boolean;
  error: string | null;
  data: CombinedUserWithId | null;
  connectionMap?: Record<string, string>;
};

export function useRecCon(route: string) {
  const [state, setState] = useState<State>({
    loading: false,
    error: null,
    data: null,
  });

  useEffect(() => {
    let cancelled = false;

    (async () => {
      setState({ loading: true, error: null, data: null });
      try {
        const response = await get<any>(route);

        let ids: string[];
        let connectionMap: Record<string, string> | undefined;

        // Handle different response formats
        if (!response) {
          // Null/undefined response
          ids = [];
        } else if (Array.isArray(response)) {
          // Old format: just array of user IDs (for /recommendations)
          ids = response;
        } else if (response.user_ids && Array.isArray(response.user_ids)) {
          // New format: object with user_ids and connection_map (for /connections/requests)
          ids = response.user_ids;
          connectionMap = response.connection_map;
        } else {
          // Fallback
          ids = [];
        }

        if (!ids || ids.length === 0) {
          if (!cancelled) {
            setState({ loading: false, error: null, data: null, connectionMap });
          }
          return;
        }

        // Load ONLY the first profile
        const id = ids[0];
        const [profile, photo] = await Promise.all([
          get<UserProfile>(`/users/${id}/profile`),
          get<UserPhoto>(`/users/${id}`),
        ]);

        const merged: CombinedUserWithId = { id, ...profile, ...photo };

        if (!cancelled) {
          setState({ loading: false, error: null, data: merged, connectionMap });
        }
      } catch (err: any) {
        if (!cancelled) {
          setState({
            loading: false,
            error: err?.message ?? "Failed to load data",
            data: null,
          });
        }
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [route]);

  return state;
}