import type { Track } from "@/types/domain";

// This file renders registered tracks so release requests can be planned.

type TrackListProps = {
  tracks: Track[];
};

export function TrackList({ tracks }: TrackListProps) {
  return (
    <section className="panel">
      <h2>곡 목록</h2>
      {tracks.length === 0 ? (
        <p className="muted">아직 등록된 곡이 없습니다.</p>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Title</th>
                <th>ISRC</th>
                <th>Genre</th>
                <th>Duration</th>
                <th>Explicit</th>
              </tr>
            </thead>
            <tbody>
              {tracks.map((track) => (
                <tr key={track.id}>
                  <td>{track.title}</td>
                  <td>{track.isrc || "-"}</td>
                  <td>{track.genre}</td>
                  <td>{formatDuration(track.durationSeconds)}</td>
                  <td>{track.explicit ? "Yes" : "No"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

function formatDuration(seconds: number) {
  const minutes = Math.floor(seconds / 60);
  const remaining = String(seconds % 60).padStart(2, "0");
  return `${minutes}:${remaining}`;
}
