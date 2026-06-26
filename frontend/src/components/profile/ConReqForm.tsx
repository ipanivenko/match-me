import "bulma/css/bulma.min.css";
import "../../styles/viewProfile.css";
import { useState, useEffect } from "react";
import { get } from "../../api/client";
import type { CombinedUser } from "../../types/profile";
import { Link, useNavigate } from "react-router-dom";
import UserHeader from "../UserHeader";
import "../../styles/UserHeader.css";
import { acceptOrRejectConnection } from "../../hooks/postConnectionAction";

type CombinedUserWithId = CombinedUser & { id: string };

interface ConnectionRequestsResponse {
  user_ids: string[];
  connection_map: Record<string, string>;
}

export default function ConnectionsReqForm() {
  const navigate = useNavigate();

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [userIds, setUserIds] = useState<string[]>([]);
  const [connectionMap, setConnectionMap] = useState<Record<string, string>>({});
  const [profiles, setProfiles] = useState<CombinedUserWithId[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [busyReact, setBusyReact] = useState(false);
  const [reactError, setReactError] = useState<string | null>(null);


  useEffect(() => {
    const fetchRequests = async () => {
      try {
        setLoading(true);
        const response = await get<ConnectionRequestsResponse>("/connections/requests");
        
        if (response.user_ids && response.user_ids.length > 0) {
          setUserIds(response.user_ids);
          setConnectionMap(response.connection_map || {});
        } else {
          setError("No connection requests found");
        }
      } catch (e: any) {
        setError(e?.message || "Failed to load requests");
      } finally {
        setLoading(false);
      }
    };

    fetchRequests();
  }, []);

 
  useEffect(() => {
    if (userIds.length === 0) return;

    const fetchProfiles = async () => {
      try {
        const profilePromises = userIds.map(async (userId) => {
          const [profileData, bioData] = await Promise.all([
            get<any>(`/users/${userId}/profile`),
            get<any>(`/users/${userId}/bio`),
          ]);

          return {
            id: userId,
            ...profileData,
            ...bioData,
          } as CombinedUserWithId;
        });

        const loadedProfiles = await Promise.all(profilePromises);
        setProfiles(loadedProfiles);
      } catch (e: any) {
        console.error("Failed to load profiles:", e);
        setError("Failed to load profile details");
      }
    };

    fetchProfiles();
  }, [userIds]);

  const currentUser = profiles[currentIndex] || null;


  const handleNext = () => {
    if (currentIndex < profiles.length - 1) {
      setCurrentIndex(currentIndex + 1);
    } else {
  
      setCurrentIndex(0);
    }
  };

 
  async function handleReaction(kind: "accept" | "reject") {
    if (!currentUser || busyReact) return;

    const connectionId = connectionMap[currentUser.id];
    if (!connectionId) {
      setReactError("Connection ID not found");
      return;
    }

    setBusyReact(true);
    setReactError(null);

    try {
      await acceptOrRejectConnection(connectionId, kind);

      
      const newProfiles = profiles.filter((_, idx) => idx !== currentIndex);
      const newUserIds = userIds.filter((id) => id !== currentUser.id);
      
      setProfiles(newProfiles);
      setUserIds(newUserIds);

      
      if (newProfiles.length > 0) {
    
        if (currentIndex >= newProfiles.length) {
          setCurrentIndex(0);
        }
      
      } else {
    
        setTimeout(() => {
          alert("All requests processed! ✅");
          if (kind === "accept") {
            navigate("/connections");
          } else {
            navigate("/recommendations");
          }
        }, 300);
      }
    } catch (e: any) {
      setReactError(e?.message ?? `Failed to ${kind} connection`);
    } finally {
      setBusyReact(false);
    }
  }

  return (
    <section className="section has-background-light">
      <Link to="/connections" className="button connect is-link is-light">
        View connections
      </Link>
      <div className="recommendations-container">
        <UserHeader />

        <div
          className="buttons is-centered"
          style={{ gap: "0.5rem", marginLeft: "0.5rem" }}
        >
          <button
            className={`button is-danger ${busyReact ? "is-loading" : ""}`}
            disabled={busyReact || !currentUser}
            onClick={() => handleReaction("reject")}
          >
            Decline
          </button>
          <button
            className={`button is-primary ${busyReact ? "is-loading" : ""}`}
            disabled={busyReact || !currentUser}
            onClick={() => handleReaction("accept")}
          >
            Accept
          </button>
          <button
            className="button is-link"
            onClick={handleNext}
            disabled={profiles.length <= 1}
          >
            Next
          </button>
        </div>

        <h1 className="title has-text-centered">Your connection requests</h1>

        {profiles.length > 0 && (
          <div className="has-text-centered mb-3">
            <span className="tag is-info is-light">
              {currentIndex + 1} / {profiles.length}
            </span>
          </div>
        )}

        <div className="user-profile with-bottom-panel">
          {loading && (
            <p className="loading-text">Loading connection requests...</p>
          )}
          {error && <p className="error-text">{error}</p>}
          {reactError && <p className="error-text">{reactError}</p>}
          {!loading && !error && !currentUser && (
            <div className="box has-text-centered">
              <p className="subtitle">No connection requests found.</p>
              <button
                className="button is-primary mt-3"
                onClick={() => navigate("/recommendations")}
              >
                Find Matches
              </button>
            </div>
          )}

          {currentUser && (
            <div className="box">
              <article className="media">
                <figure className="media-left">
                  <div className="avatar-wrapper">
                    {currentUser.avatarurl ? (
                      <img src={currentUser.avatarurl} alt={`${currentUser.name} avatar`} />
                    ) : (
                      <span className="avatar-fallback">👤</span>
                    )}
                  </div>
                </figure>

                <div className="media-content">
                  <h2 className="parent-name">{currentUser.name}</h2>
                  <p className="parent-city">{currentUser.addressCity}</p>
                  {currentUser.about && <p className="parent-about">{currentUser.about}</p>}

                  {Array.isArray(currentUser.languages) &&
                    currentUser.languages.length > 0 && (
                      <>
                        <p className="has-text-weight-semibold mb-1">Languages</p>
                        <div className="tags mb-3">
                          {currentUser.languages.map((lang) => (
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
                      <span className="child-name">{currentUser.child.name}</span>
                    </p>
                    <p>
                      <strong>Gender:</strong> {currentUser.child.gender}
                    </p>
                  </div>
                  <div className="column is-half">
                    <p>
                      <strong>Age:</strong> {currentUser.child.ageYears}{" "}
                      {currentUser.child.ageYears === 1 ? "year" : "years"}
                    </p>
                  </div>

                  {currentUser.child.aboutShort && (
                    <div className="column is-full">
                      <p className="child-about">{currentUser.child.aboutShort}</p>
                    </div>
                  )}

                  {Array.isArray(currentUser.child.topInterests) &&
                    currentUser.child.topInterests.length > 0 && (
                      <div className="column is-full top-interests">
                        <p className="has-text-weight-semibold mb-1">Top interests</p>
                        <div className="tags">
                          {currentUser.child.topInterests.map((i) => (
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