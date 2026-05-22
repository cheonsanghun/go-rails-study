"use client";

// This file renders the track registration form and returns typed payloads upward.

import { FormEvent, useState } from "react";
import type { Artist, CreateTrackPayload } from "@/types/domain";

type TrackFormProps = {
  artist: Artist;
  onCreate: (payload: CreateTrackPayload) => Promise<void>;
};

const initialForm: CreateTrackPayload = {
  title: "",
  isrc: "",
  durationSeconds: 180,
  genre: "",
  language: "ja",
  explicit: false,
  audioFileName: "",
};

export function TrackForm({ artist, onCreate }: TrackFormProps) {
  const [form, setForm] = useState<CreateTrackPayload>({
    ...initialForm,
    genre: artist.primaryGenre,
  });
  const [saving, setSaving] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSaving(true);
    await onCreate(form);
    setForm({ ...initialForm, genre: artist.primaryGenre });
    setSaving(false);
  }

  return (
    <section className="panel">
      <h2>곡 등록</h2>
      <form className="form-grid" onSubmit={handleSubmit}>
        <label className="full-span">
          Track title
          <input
            required
            value={form.title}
            onChange={(event) =>
              setForm({ ...form, title: event.target.value })
            }
            placeholder="New single title"
          />
        </label>

        <label>
          ISRC
          <input
            value={form.isrc}
            onChange={(event) => setForm({ ...form, isrc: event.target.value })}
            placeholder="JPA1B2600123"
          />
        </label>

        <label>
          Duration seconds
          <input
            min={1}
            required
            type="number"
            value={form.durationSeconds}
            onChange={(event) =>
              setForm({ ...form, durationSeconds: Number(event.target.value) })
            }
          />
        </label>

        <label>
          Genre
          <input
            value={form.genre}
            onChange={(event) =>
              setForm({ ...form, genre: event.target.value })
            }
          />
        </label>

        <label>
          Language
          <select
            value={form.language}
            onChange={(event) =>
              setForm({ ...form, language: event.target.value })
            }
          >
            <option value="ja">Japanese</option>
            <option value="en">English</option>
            <option value="ko">Korean</option>
            <option value="instrumental">Instrumental</option>
          </select>
        </label>

        <label className="full-span">
          Audio file name
          <input
            value={form.audioFileName}
            onChange={(event) =>
              setForm({ ...form, audioFileName: event.target.value })
            }
            placeholder="track-master.wav"
          />
        </label>

        <label className="checkbox-row full-span">
          <input
            checked={form.explicit}
            type="checkbox"
            onChange={(event) =>
              setForm({ ...form, explicit: event.target.checked })
            }
          />
          Explicit lyrics
        </label>

        <button className="full-span" disabled={saving} type="submit">
          {saving ? "Saving..." : "Register track"}
        </button>
      </form>
    </section>
  );
}
