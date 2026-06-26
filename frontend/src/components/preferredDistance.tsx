import { useId, useMemo } from "react";
import "../styles/preferredDistance.css";

type Props = {
  label?: string;
  value: number;          
  min?: number;
  max?: number;
  step?: number;
  onChange: (v: number) => void;
};

export default function PreferredDistanceField({
  label = "Preferred Radius (km)",
  value,
  min = 0,
  max = 150,
  step = 1,
  onChange,
}: Props) {
  const id = useId();

  const safeMin = Number.isFinite(min) ? min : 0;
  const safeMax = Number.isFinite(max) && max !== safeMin ? max : safeMin + 1;
  const clamped = Math.min(Math.max(value, safeMin), safeMax);
  const pct = useMemo(
    () => ((clamped - safeMin) / (safeMax - safeMin)) * 100,
    [clamped, safeMin, safeMax]
  );

  return (
    <div className="field">
      <label className="label" htmlFor={id}>{label}</label>
      <div className="control">
        <input
          id={id}
          className="slider"
          type="range"
          min={safeMin}
          max={safeMax}
          step={step}
          value={clamped}
          onChange={(e) => onChange(Number(e.target.value))}
          style={{ ["--pct" as any]: `${pct}%` }}
        />
      </div>
      <div className="is-flex is-justify-content-space-between is-align-items-center mt-2">
        <span className="is-size-7">{safeMin} km</span>
        <span className="tag is-primary is-light">{clamped} km</span>
        <span className="is-size-7">{safeMax} km</span>
      </div>
    </div>
  );
}
