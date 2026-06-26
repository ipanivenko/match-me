import { useState, useRef } from "react";
import "../styles/userPhoto.css";

type Props = {
  onChange?: (file: File | null) => void;
  initialUrl?: string | null;
  placeholderEmoji?: string; 
  name?: string;
};

export default function UserPhotoField({
  onChange,
  initialUrl = null,
  placeholderEmoji = "👤",
  name = "photo",
}: Props) {

  const normalize = (u?: string | null) =>
  u && u.trim() !== "" ? u : null;

  const [preview, setPreview] = useState<string | null>(normalize(initialUrl));
  
  const inputRef = useRef<HTMLInputElement | null>(null);

  const handlePick = (f: File | null) => {
    onChange?.(f);

    if (!f) {
      setPreview(null);
      return;
    }
    const reader = new FileReader();
    reader.onload = () => setPreview(reader.result as string);
    reader.readAsDataURL(f);
  };

  return (
    <div className="field">
      <label className="label">Profile Photo</label>

      <label
        className={`photo-slot ${preview ? "" : "is-empty"}`}
        title="Click to upload"
        role="button"
        tabIndex={0}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            inputRef.current?.click();
          }
        }}>
        {preview ? (
          <img src={preview} alt="profile" />
        ) : (
          <span
            className="placeholder-emoji"
            role="img"
            aria-label="Upload photo">
            {placeholderEmoji}
          </span>
        )}

        <input
          ref={inputRef}
          type="file"
          name={name}
          accept="image/*"
          style={{ display: "none" }}
          onChange={(e) => {
            const f = e.target.files?.[0] ?? null;
            handlePick(f);
            if (inputRef.current) inputRef.current.value = "";
          }}
        />

        {preview && (
          <button
            type="button"
            className="remove"
            onClick={(e) => {
              e.preventDefault();
              e.stopPropagation(); 
              handlePick(null);
              if (inputRef.current) inputRef.current.value = "";
            }}
            aria-label="Remove photo"
            title="Remove photo">
            ×
          </button>
        )}
      </label>

      <p className="help">
        Click to upload one image. Use a square photo for best fit.
      </p>
    </div>
  );
}
