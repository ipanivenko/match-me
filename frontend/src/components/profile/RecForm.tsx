import "bulma/css/bulma.min.css";
import "../../styles/viewProfile.css";
import { useState } from "react";
import { useLogout } from "../../auth/useLogout";
import { useRecCon } from "../../hooks/useRecCon";
import { loadNextProfile } from "../../hooks/loadNextProfile";
import { reactToUser } from "../../hooks/postReaction";
import type { CombinedUser } from "../../types/profile";
import UserHeader from "../UserHeader";
import "../../styles/UserHeader.css";

type CombinedUserWithId = CombinedUser & { id: string };


export default function RecommendationsForm() {
  const route = "/recommendations";
  const logout = useLogout();

  const { loading, error, data } = useRecCon(route);

  const [ids, setIds] = useState<string[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [loadingNext, setLoadingNext] = useState(false);
  const [busyReact, setBusyReact] = useState(false);
  const [nextUser, setNextUser] = useState<CombinedUserWithId | null>(null);
  const [reactError, setReactError] = useState<string | null>(null);
  const [isLast, setIsLast] = useState(false);

  // Prefer the nextUser (after pressing Next), otherwise show the initial one
  const user =
    (nextUser as CombinedUserWithId | null) ??
    (data as CombinedUserWithId | null);

  const handleNext = () =>{
     if (loadingNext || isLast) return;
    loadNextProfile(
      route,
      ids,
      currentIndex,
      setIds,
      setNextUser,
      setCurrentIndex,
      setLoadingNext,
      setIsLast
    );
  }

  async function handleReaction(kind: "like" | "dislike") {
    if (!user || busyReact || loadingNext) return;
    setBusyReact(true);
    setReactError(null);

  
    const prevUser = user;
    handleNext();

    try {
      await reactToUser(route, prevUser.id, kind);
    } catch (e: any) {
      setReactError(e?.message ?? "Failed to send reaction");
    } finally {
      setBusyReact(false);
    }
  }

  return (
    <section className="section has-background-light">
       <UserHeader />
      <div className="recommendations-container">
        <button className="button is-dark logout" onClick={() => logout()}>
          Log out
        </button>

        <div
          className="buttons is-centered"

          style={{ gap: "0.5rem", marginLeft: "0.5rem" }}>
          <button
            className={`button is-danger ${busyReact ? "is-loading" : ""}`}
            disabled={busyReact || loadingNext || !user}
            onClick={() => handleReaction("dislike")}>
            Dismiss
          </button>
          <button
            className={`button is-primary ${busyReact ? "is-loading" : ""}`}
            disabled={busyReact || loadingNext || !user}
            onClick={() => handleReaction("like")}>
            Connect
          </button>
          <button
            className={`button is-link ${loadingNext ? "is-loading" : ""}`}
            onClick={handleNext}
           disabled={loadingNext || isLast}>
            Next
          </button>
        </div>

        <h1 className="title has-text-centered">Your recommendations</h1>

        <div className="user-profile with-bottom-panel">
          {loading && (
            <p className="loading-text">Loading recommendations...</p>
          )}
          {error && <p className="error-text">{error}</p>}
          {reactError && <p className="error-text">{reactError}</p>}
          {!loading && !error && !user && <p>No recommendations found.</p>}

          {user && (
            <div className="box">
              <article className="media">
                <figure className="media-left">
                  <div className="avatar-wrapper">
                    {user.avatarurl ? (
                      <img src={user.avatarurl} alt={`${user.name} avatar`} />
                    ) : (
                      <span className="avatar-fallback">👤</span>
                    )}
                  </div>
                </figure>

                <div className="media-content">
                  <h2 className="parent-name">{user.name}</h2>
                  <p className="parent-city">{user.addressCity}</p>
                  {user.about && <p className="parent-about">{user.about}</p>}

                  {Array.isArray(user.languages) &&
                    user.languages.length > 0 && (
                      <>
                        <p className="has-text-weight-semibold mb-1">
                          Languages
                        </p>
                        <div className="tags mb-3">
                          {user.languages.map((lang) => (
                            <span key={lang} className="tag is-info is-light">
                              {lang.toUpperCase()}
                            </span>
                          ))}
                        </div>
                      </>
                    )}
                </div>
              </article>

              <hr className="my-4" />

              <div className="content child-block">
                <h3 className="title is-5">Child</h3>
                <div className="columns is-multiline">
                  <div className="column is-half">
                    <p>
                      <strong>Name:</strong>{" "}
                      <span className="child-name">{user.child.name}</span>
                    </p>
                    <p>
                      <strong>Gender:</strong> {user.child.gender}
                    </p>
                  </div>
                  <div className="column is-half">
                    <p>
                      <strong>Age:</strong> {user.child.ageYears}{" "}
                      {user.child.ageYears === 1 ? "year" : "years"}
                    </p>
                  </div>

                  {user.child.aboutShort && (
                    <div className="column is-full">
                      <p className="child-about">{user.child.aboutShort}</p>
                    </div>
                  )}

                  {Array.isArray(user.child.topInterests) &&
                    user.child.topInterests.length > 0 && (
                      <div className="column is-full top-interests">
                        <p className="has-text-weight-semibold mb-1">
                          Top interests
                        </p>
                        <div className="tags">
                          {user.child.topInterests.map((i) => (
                            <span key={i} className="tag is-success is-light">
                              {i}
                            </span>
                          ))}
                        </div>
                      </div>
                    )}
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </section>
  );
}
