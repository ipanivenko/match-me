import { get } from "../api/client";
import type { UserProfile, UserPhoto, CombinedUser } from "../types/profile";

type CombinedUserWithId = CombinedUser & { id: string };

export async function loadNextProfile(
  route: string,
  ids: string[],
  currentIndex: number,
  setIds: (ids: string[]) => void,
  setNextUser: (user: CombinedUserWithId) => void,
  setCurrentIndex: (index: number) => void,
  setLoadingNext: (loading: boolean) => void,
  setIsLast: (isLast: boolean) => void
) {
  let localIds = ids;

  try {
    // Load recommendation IDs if empty
    if (localIds.length === 0) {
      const list = await get<string[]>(route);
      setIds(list);
      if (list.length === 0) return;
      localIds = list;
    }

    const nextIndex = currentIndex + 1;
    const newIsLast = nextIndex >= localIds.length - 1;

    // update UI state
    setIsLast(newIsLast);
    if (nextIndex >= localIds.length) return;
    setLoadingNext(true);

    const id = localIds[nextIndex];
    const [profile, photo] = await Promise.all([
      get<UserProfile>(`/users/${id}/profile`),
      get<UserPhoto>(`/users/${id}`),
    ]);

    setNextUser({ id, ...profile, ...photo });
    setCurrentIndex(nextIndex);
  } catch (err) {
    console.error("Error loading next profile:", err);
  } finally {
    setLoadingNext(false);
  }
}
