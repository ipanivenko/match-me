import { useId } from "react";
import "../styles/profiles.css";

type Props = {
  maxLanguages?: number;
  languages: string[];
  onChange: (langs: string[]) => void;
};

export default function UserLanguagesField({
  maxLanguages = 3,
  languages,
  onChange,
}: Props) {
  const listId = useId();


  const safe = [...languages, ...Array(Math.max(0, maxLanguages - languages.length)).fill("")]
    .slice(0, maxLanguages);

  const handleChange = (i: number, value: string) => {
    const next = [...safe];
    next[i] = value;          
    onChange(next);
  };

  const handleBlur = (i: number) => {
    const trimmed = safe[i].trim();
    if (trimmed !== safe[i]) {
      const next = [...safe];
      next[i] = trimmed;
      onChange(next);
    }
  };

  return (
    <div className="field">
      <label className="label">Languages</label>

      {Array.from({ length: maxLanguages }, (_, i) => (
        <div className="control" style={{ marginTop: i === 0 ? 0 : "0.5rem" }} key={i}>
          <input
            name={`language-${i}`}
            className="input"
            type="text"
            list={listId}
            placeholder="Type a language..."
            value={safe[i] ?? ""}
            onChange={(e) => handleChange(i, e.target.value)}
            onBlur={() => handleBlur(i)}
          />
        </div>
      ))}

      <datalist id={listId}>
        <option value="English" />
        <option value="French" />
        <option value="Russian" />
        <option value="German" />
        <option value="Spanish" />
        <option value="Italian" />
        <option value="Portuguese" />
        <option value="Chinese (Mandarin)" />
        <option value="Japanese" />
        <option value="Korean" />
        <option value="Arabic" />
        <option value="Hindi" />
        <option value="Bengali" />
        <option value="Urdu" />
        <option value="Turkish" />
        <option value="Dutch" />
        <option value="Polish" />
        <option value="Swedish" />
        <option value="Finnish" />
        <option value="Norwegian" />
        <option value="Danish" />
        <option value="Greek" />
        <option value="Hebrew" />
        <option value="Thai" />
        <option value="Vietnamese" />
        <option value="Malay" />
        <option value="Indonesian" />
        <option value="Tagalog" />
        <option value="Swahili" />
      </datalist>

      <p className="help">Choose up to {maxLanguages} languages.</p>
    </div>
  );
}
