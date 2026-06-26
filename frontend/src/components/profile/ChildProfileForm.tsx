import "bulma/css/bulma.min.css";
import "../../styles/profiles.css";
import { useChildProfile } from "../../hooks/useChildProfile";
import { useState, useEffect, useRef } from "react";
import { buildChildPayload } from "./updateChildProfile";
import type { ChildFields } from "./updateChildProfile";
import { saveProfile } from "../../hooks/patchUser";
import UserHeader from "../UserHeader";
import "../../styles/UserHeader.css";

export default function ChildProfileForm() {
  const { data } = useChildProfile();

  const [initialized, setInitialized] = useState(false);
  

  // ---- Local editable state ----
  const [name, setName] = useState("");
  const [birthday, setBirthday] = useState<string>("");
  const [gender, setGender] = useState("");
  const [about_short, setAbout_short] = useState("");
  const [interests, setInterests] = useState<string[]>([]);
  const [activity_level, setActivity_level] = useState("");
  const [limitations, setLimitations] = useState<string[]>([]);
  const [allergies, setAllergies] = useState<string[]>([]);
  const [play_styles, setPlay_styles] = useState<string[]>([]);

  const [interests_weight, setInterestsWeight] = useState<number>(1);
  const [activity_level_weight, setActivityLevelWeight] = useState<number>(2);
  const [limitations_weight, setLimitationsWeight] = useState<number>(3);
  const [allergies_weight, setAllergiesWeight] = useState<number>(3);
  const [play_styles_weight, setPlayStylesWeight] = useState<number>(1);
  const [max_age_difference, setMaxAgeDifference] = useState<number>(2);

  // ---- Refs for originals
  const originalName = useRef("");
  const originalBirthday = useRef("");
  const originalGender = useRef("");
  const originalAboutShort = useRef("");
  const originalInterests = useRef<string[]>([]);
  const originalActivityLevel = useRef("");
  const originalLimitations = useRef<string[]>([]);
  const originalAllergies = useRef<string[]>([]);
  const originalPlayStyles = useRef<string[]>([]);

  const originalInterestsWeight = useRef(1);
  const originalActivityLevelWeight = useRef(2);
  const originalLimitationsWeight = useRef(3);
  const originalAllergiesWeight = useRef(3);
  const originalPlayStylesWeight = useRef(1);
  const originalMaxAgeDifference = useRef(2);

  // Initialize from API once
  useEffect(() => {
    if (!data || initialized) return;

    setName(data.name ?? "");
    setBirthday(data.birthday ?? "");
    setGender(data.gender ?? "");
    setAbout_short(data.about_short ?? "");
    setInterests(data.interests ?? []);
    setActivity_level(data.activity_level ?? "");
    setLimitations(data.limitations ?? []);
    setAllergies(data.allergies ?? []);
    setPlay_styles(data.play_styles ?? []);
    setInterestsWeight(data.interests_weight ?? 1);
    setActivityLevelWeight(data.activity_level_weight ?? 2);
    setLimitationsWeight(data.limitations_weight ?? 3);
    setAllergiesWeight(data.allergies_weight ?? 3);
    setPlayStylesWeight(data.play_styles_weight ?? 1);
    setMaxAgeDifference(data.max_age_difference ?? 2);
    setInitialized(true);

    originalName.current = data.name ?? "";
    originalBirthday.current = data.birthday ?? "";
    originalGender.current = data.gender ?? "";
    originalAboutShort.current = data.about_short ?? "";
    originalInterests.current = data.interests ?? [];
    originalActivityLevel.current = data.activity_level ?? "";
    originalLimitations.current = data.limitations ?? [];
    originalAllergies.current = data.allergies ?? [];
    originalPlayStyles.current = data.play_styles ?? [];
    originalInterestsWeight.current = data.interests_weight ?? 1;
    originalActivityLevelWeight.current = data.activity_level_weight ?? 2;
    originalLimitationsWeight.current = data.limitations_weight ?? 3;
    originalAllergiesWeight.current = data.allergies_weight ?? 3;
    originalPlayStylesWeight.current = data.play_styles_weight ?? 1;
    originalMaxAgeDifference.current = data.max_age_difference ?? 2;
  }, [data, initialized]);

  async function handleSaveChild() {
    // Build current/original snapshots
    const current: ChildFields = {
      name,
      birthday,
      gender,
      about_short,
      interests,
      activity_level,
      limitations,
      allergies,
      play_styles,
      interests_weight,
      activity_level_weight,
      limitations_weight,
      allergies_weight,
      play_styles_weight,
      max_age_difference,
    };

    const original: ChildFields = {
      name: originalName.current,
      birthday: originalBirthday.current,
      gender: originalGender.current,
      about_short: originalAboutShort.current,
      interests: originalInterests.current,
      activity_level: originalActivityLevel.current,
      limitations: originalLimitations.current,
      allergies: originalAllergies.current,
      play_styles: originalPlayStyles.current,
      interests_weight: originalInterestsWeight.current,
      activity_level_weight: originalActivityLevelWeight.current,
      limitations_weight: originalLimitationsWeight.current,
      allergies_weight: originalAllergiesWeight.current,
      play_styles_weight: originalPlayStylesWeight.current,
      max_age_difference: originalMaxAgeDifference.current,
    };

    // Compute payload
    const payload = buildChildPayload(current, original);
    if (Object.keys(payload).length === 0) return;

    //setSaving(true);
    try {
      await saveProfile(payload, "/me/child");

      // Sync refs after success so hasChanges becomes false on next render
      if ("name" in payload) originalName.current = name.trim();
      if ("birthday" in payload) originalBirthday.current = birthday;
      if ("gender" in payload) originalGender.current = gender;
      if ("about_short" in payload)
        originalAboutShort.current = about_short.trim();
      if ("activity_level" in payload)
        originalActivityLevel.current = activity_level;

      if ("interests" in payload) originalInterests.current = [...interests];
      if ("limitations" in payload)
        originalLimitations.current = [...limitations];
      if ("allergies" in payload) originalAllergies.current = [...allergies];
      if ("play_styles" in payload)
        originalPlayStyles.current = [...play_styles];
      if ("interests_weight" in payload)
        originalInterestsWeight.current = interests_weight;
      if ("activity_level_weight" in payload)
        originalActivityLevelWeight.current = activity_level_weight;
      if ("limitations_weight" in payload)
        originalLimitationsWeight.current = limitations_weight;
      if ("allergies_weight" in payload)
        originalAllergiesWeight.current = allergies_weight;
      if ("play_styles_weight" in payload)
        originalPlayStylesWeight.current = play_styles_weight;
      if ("max_age_difference" in payload)
        originalMaxAgeDifference.current = max_age_difference;
    } catch (e) {
      console.error(e);
    }
  }
  return (
    <section className="section has-background-light">
         <UserHeader />
           <div className="container">
        <h1 className="title has-text-centered">Child Profile</h1>

        <form className="child-profile with-bottom-panel">
          {/* Name */}
          <div className="field">
            <label className="label">Name</label>
            <div className="control">
              <input
                className="input"
                type="text"
                name="name"
                placeholder="Enter name"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
          </div>

          {/* Birthday */}
          <div className="field is-horizontal">
            <div className="field-body">
              {/* Birthday */}
              <div className="field">
                <label className="label">Birthday</label>
                <div className="control">
                  <input
                    className="input"
                    type="date"
                    name="birthday"
                    value={birthday}
                    max={new Date().toISOString().split("T")[0]}
                    onChange={(e) => setBirthday(e.target.value)}
                  />
                </div>
              </div>

              {/* Max Age Difference */}
              <div className="field">
                <label className="label">Max Age Difference</label>
                <div className="control rate">
                  <input
                    className="input"
                    type="number"
                    min="0"
                    max="5"
                    value={max_age_difference}
                    onChange={(e) =>
                     setMaxAgeDifference(Number(e.target.value))
                     }
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Gender */}
          <div className="field">
            <label className="label">Gender</label>
            <div className="control">
              <div className="select">
                <select
                  name="gender"
                  value={gender}
                  onChange={(e) => setGender(e.target.value)}>
                  <option value="">Select gender</option>
                  <option value="male">Boy</option>
                  <option value="female">Girl</option>
                  <option value="other">Other</option>
                </select>
              </div>
            </div>
          </div>

          {/* About_short */}
          <div className="field">
            <label className="label">About short</label>
            <div className="control">
              <textarea
                className="textarea"
                name="about_short"
                placeholder="Short description"
                value={about_short}
                onChange={(e) => setAbout_short(e.target.value)}
              />
            </div>
          </div>

          {/* Interests */}
          <div className="field is-horizontal">
            <div className="field-body">
              <div className="field">
                <label className="label">Interests</label>
                <div className="control">
                  <input
                    className="input"
                    type="text"
                    name="interests"
                    placeholder="e.g. football, drawing"
                    value={interests.join(", ")}
                    onChange={(e) => setInterests([e.target.value])}
                    onBlur={(e) =>
                      setInterests(
                        e.target.value
                          .split(",")
                          .map((s) => s.trim())
                          .filter(Boolean)
                      )
                    }
                  />
                </div>
              </div>

              <div className="field">
                <label className="label">Importancy of match</label>
                <div className="control">
                  <div className="select rate">
                    <select
                      name="interests_weight"
                      value={interests_weight}
                      onChange={(e) =>
                        setInterestsWeight(Number(e.target.value))
                      }>
                      <option value="0">0</option>
                      <option value="1">1</option>
                      <option value="2">2</option>
                      <option value="3">3</option>
                      <option value="4">4</option>
                      <option value="5">5</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Activity level */}
          <div className="field is-horizontal">
            <div className="field-body">
              <div className="field">
                <label className="label">Activity level</label>
                <div className="control">
                  <div className="select">
                    <select
                      name="activity_level"
                      value={activity_level}
                      onChange={(e) => setActivity_level(e.target.value)}>
                      <option value="">Select level</option>
                      <option value="low">Low</option>
                      <option value="medium">Medium</option>
                      <option value="high">High</option>
                    </select>
                  </div>
                </div>
              </div>

              <div className="field">
                <label className="label">Importancy of match</label>
                <div className="control">
                  <div className="select rate">
                    <select
                      name="activity_level_weight"
                      value={activity_level_weight}
                      onChange={(e) =>
                        setActivityLevelWeight(Number(e.target.value))
                      }>
                     <option value="0">0</option>
                      <option value="1">1</option>
                      <option value="2">2</option>
                      <option value="3">3</option>
                      <option value="4">4</option>
                      <option value="5">5</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Limitations */}
          <div className="field is-horizontal">
            <div className="field-body">
              <div className="field">
                <label className="label">Limitations</label>
                <div className="control">
                  <input
                    className="input"
                    type="text"
                    name="limitations"
                    placeholder="e.g. no climbing, no dairy"
                    value={limitations.join(", ")}
                    onChange={(e) => setLimitations([e.target.value])}
                    onBlur={(e) =>
                      setLimitations(
                        e.target.value
                          .split(",")
                          .map((s) => s.trim())
                          .filter(Boolean)
                      )
                    }
                  />
                </div>
              </div>

              <div className="field">
                <label className="label">Importancy of match</label>
                <div className="control">
                  <div className="select rate">
                    <select
                      name="limitations_weight"
                      value={limitations_weight}
                      onChange={(e) =>
                        setLimitationsWeight(Number(e.target.value))
                      }>
                      <option value="0">0</option>
                      <option value="1">1</option>
                      <option value="2">2</option>
                      <option value="3">3</option>
                      <option value="4">4</option>
                      <option value="5">5</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Allergies */}
          <div className="field is-horizontal">
            <div className="field-body">
              <div className="field">
                <label className="label">Allergies</label>
                <div className="control">
                  <input
                    className="input"
                    type="text"
                    name="allergies"
                    placeholder="e.g. peanuts, pollen or NONE"
                    value={allergies.join(", ")}
                    onChange={(e) => setAllergies([e.target.value])}
                    onBlur={(e) =>
                      setAllergies(
                        e.target.value
                          .split(",")
                          .map((s) => s.trim())
                          .filter(Boolean)
                      )
                    }
                  />
                </div>
              </div>

              <div className="field">
                <label className="label">Importancy of match</label>
                <div className="control">
                  <div className="select rate">
                    <select
                      name="allergies_weight"
                      value={allergies_weight}
                      onChange={(e) =>
                        setAllergiesWeight(Number(e.target.value))
                      }>
                      <option value="0">0</option>
                      <option value="1">1</option>
                      <option value="2">2</option>
                      <option value="3">3</option>
                      <option value="4">4</option>
                      <option value="5">5</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Play styles */}
          <div className="field is-horizontal">
            <div className="field-body">
              <div className="field">
                <label className="label">Play styles</label>
                <div className="control">
                  <input
                    className="input"
                    type="text"
                    name="play_styles"
                    placeholder="e.g. role play, building, puzzles"
                    value={play_styles.join(", ")}
                    onChange={(e) => setPlay_styles([e.target.value])}
                    onBlur={(e) =>
                      setPlay_styles(
                        e.target.value
                          .split(",")
                          .map((s) => s.trim())
                          .filter(Boolean)
                      )
                    }
                  />
                </div>
              </div>

              <div className="field">
                <label className="label">Importancy of match</label>
                <div className="control">
                  <div className="select rate">
                    <select
                      name="play_styles_weight"
                      value={play_styles_weight}
                      onChange={(e) =>
                        setPlayStylesWeight(Number(e.target.value))
                      }>
                      <option value="0">0</option>
                      <option value="1">1</option>
                      <option value="2">2</option>
                      <option value="3">3</option>
                      <option value="4">4</option>
                      <option value="5">5</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div className="field">
            <div className="control">
              <button
                className="button is-primary"
                type="button"
                onClick={handleSaveChild}>
                Save changes
              </button>
            </div>
          </div>
        </form>
      </div>
    </section>
  );
}
