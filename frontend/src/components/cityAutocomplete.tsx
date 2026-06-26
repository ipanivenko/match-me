import { useEffect, useRef, useState } from "react";
import "../styles/cityAutocomplete.css";

type City = {
  label: string;
  countryCode?: string;
  lat: number;
  lon: number;
  placeId?: string;
};

type Props = {
  country?: string; 
  value: City | null;
  onChange: (city: City | null) => void;
  placeholder?: string;
  onSelect?: (city: City) => void;
};

export default function CityAutocomplete({
  country = "FI",
  placeholder = "Start typing your city…",
  value,
  onChange,
  onSelect,
}: Props) {
  const [query, setQuery] = useState(value?.label ?? "");
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [items, setItems] = useState<City[]>([]);
  const boxRef = useRef<HTMLDivElement>(null);
  const apiKey = import.meta.env.VITE_GEOAPIFY_KEY as string;

  const [touched, setTouched] = useState(false);

  useEffect(() => {
    setQuery(value?.label ?? "");
    setOpen(false);
    setTouched(false);
  }, [value]);

  // Close dropdown when clicking outside
  useEffect(() => {
    const onDocClick = (e: MouseEvent) => {
      if (!boxRef.current?.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onDocClick);
    return () => document.removeEventListener("mousedown", onDocClick);
  }, []);

  // Debounced search
  useEffect(() => {
    if (!apiKey) return;
    const q = query?.trim();

    if (!touched || q.length < 2) {
      setItems([]);
      setOpen(false);
      return;
    }
    setLoading(true);
    const t = setTimeout(async () => {
      try {
        const url = new URL("https://api.geoapify.com/v1/geocode/autocomplete");
        url.searchParams.set("text", q);
        url.searchParams.set("filter", `countrycode:${country.toLowerCase()}`);
        url.searchParams.set("limit", "8");
        url.searchParams.set("type", "city");
        url.searchParams.set("apiKey", apiKey);

        const res = await fetch(url.toString());
        const data = await res.json();

        const results: City[] = (data?.features ?? [])
          .map((f: any) => {
            const props = f.properties ?? {};
            if (props.result_type && props.result_type !== "city") return null;
            const label =
              props.formatted ||
              [
                props.city,
                props.county || props.state,
                props.country_code?.toUpperCase(),
              ]
                .filter(Boolean)
                .join(", ");
            const lat = props.lat ?? f.geometry?.coordinates?.[1];
            const lon = props.lon ?? f.geometry?.coordinates?.[0];
            const cc = (props.country_code || country || "").toUpperCase();
            const placeId = props.place_id || f.properties?.place_id || "";
            if (!label || lat == null || lon == null) return null;
            return {
              label,
              countryCode: cc,
              lat: Number(lat),
              lon: Number(lon),
              placeId,
            };
          })
          .filter(Boolean) as City[];

        setItems(results);
        setOpen(touched && results.length > 0);
      } catch {
        setItems([]);
      } finally {
        setLoading(false);
        setOpen(true);
      }
    }, 300);

    return () => clearTimeout(t);
  }, [query, apiKey, country]);

  const handlePick = (item: City) => {
    setQuery(item.label);
    setOpen(false);
    setTouched(false);  
    onChange(item);
    onSelect?.(item);
  };

  return (
    <div className="field" ref={boxRef}>
      <label className="label">City</label>

      <div
        className={`dropdown ${open && items.length ? "is-active" : ""}`}
        style={{ width: "100%" }}>
        <div className="dropdown-trigger" style={{ width: "100%" }}>
          <input
            className="input"
            type="text"
            value={query}
            placeholder={placeholder}
            onFocus={() => setOpen(touched && items.length > 0)}
            onChange={(e) => {
              const next = e.target.value;
              setTouched(true);
              setQuery(next);
              if (!next.trim()) onChange(null);
              {
                onChange(null); // clearing text clears selection
                setItems([]);
                setOpen(false);
              }
            }}
            aria-haspopup="true"
            aria-controls="city-suggestions"
          />
        </div>

        <div
          className="dropdown-menu"
          id="city-suggestions"
          role="menu"
          style={{ width: "100%" }}>
          <div className="dropdown-content city-dropdown">
            {loading && (
              <div className="dropdown-item is-size-7">Searching…</div>
            )}
            {!loading && items.length === 0 && query && (
              <div className="dropdown-item is-size-7">No matches</div>
            )}
            {!loading &&
              items.map((item) => (
                <a
                  key={(item.placeId ?? "") + item.label}
                  className="dropdown-item"
                  onMouseDown={(e) => {
                    e.preventDefault();
                    handlePick(item);
                  }}>
                  {item.label}
                </a>
              ))}
          </div>
        </div>
      </div>

      <p className="help">
        Type your city. Pick from the list to save coordinates.
      </p>
    </div>
  );
}
