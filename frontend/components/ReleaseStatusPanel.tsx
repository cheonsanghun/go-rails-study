"use client";

// This file updates overall release status and per-store delivery status.

import { useState } from "react";
import type {
  ReleaseRequest,
  ReleaseStatus,
  StoreDeliveryStatus,
  UpdateReleaseStatusPayload,
} from "@/types/domain";

type ReleaseStatusPanelProps = {
  releaseRequests: ReleaseRequest[];
  onUpdate: (
    releaseId: string,
    payload: UpdateReleaseStatusPayload,
  ) => Promise<void>;
};

const releaseStatuses: ReleaseStatus[] = [
  "submitted",
  "in_review",
  "approved",
  "delivered",
  "live",
  "rejected",
];

const storeStatuses: StoreDeliveryStatus[] = [
  "pending",
  "processing",
  "delivered",
  "live",
  "failed",
];

export function ReleaseStatusPanel({
  releaseRequests,
  onUpdate,
}: ReleaseStatusPanelProps) {
  const [releaseDrafts, setReleaseDrafts] = useState<Record<string, ReleaseStatus>>(
    {},
  );
  const [storeDrafts, setStoreDrafts] = useState<
    Record<string, StoreDeliveryStatus>
  >({});

  return (
    <section className="panel">
      <h2>배포 상태</h2>
      {releaseRequests.length === 0 ? (
        <p className="muted">아직 배포 신청이 없습니다.</p>
      ) : (
        <div className="stack">
          {releaseRequests.map((request) => (
            <article className="release-card" key={request.id}>
              <div className="workspace-header">
                <div>
                  <h3>{request.title}</h3>
                  <p className="muted">
                    Release date: {formatReleaseDate(request.releaseDate)}
                  </p>
                </div>
                <span className="status">{request.status}</span>
              </div>

              <div className="form-grid">
                <label>
                  Overall status
                  <select
                    value={releaseDrafts[request.id] ?? request.status}
                    onChange={(event) =>
                      setReleaseDrafts({
                        ...releaseDrafts,
                        [request.id]: event.target.value as ReleaseStatus,
                      })
                    }
                  >
                    {releaseStatuses.map((status) => (
                      <option key={status} value={status}>
                        {status}
                      </option>
                    ))}
                  </select>
                </label>
                <button
                  type="button"
                  onClick={() =>
                    onUpdate(request.id, {
                      status: releaseDrafts[request.id] ?? request.status,
                    })
                  }
                >
                  Update
                </button>
              </div>

              <div className="stack">
                {request.storeDeliveries.map((delivery) => {
                  const draftKey = `${request.id}:${delivery.storeId}`;
                  return (
                    <div className="store-row" key={delivery.storeId}>
                      <span>{delivery.storeName}</span>
                      <select
                        value={storeDrafts[draftKey] ?? delivery.status}
                        onChange={(event) =>
                          setStoreDrafts({
                            ...storeDrafts,
                            [draftKey]: event.target
                              .value as StoreDeliveryStatus,
                          })
                        }
                      >
                        {storeStatuses.map((status) => (
                          <option key={status} value={status}>
                            {status}
                          </option>
                        ))}
                      </select>
                      <button
                        className="secondary"
                        type="button"
                        onClick={() =>
                          onUpdate(request.id, {
                            storeId: delivery.storeId,
                            storeStatus:
                              storeDrafts[draftKey] ?? delivery.status,
                          })
                        }
                      >
                        Save
                      </button>
                    </div>
                  );
                })}
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}

function formatReleaseDate(value: string) {
  return new Intl.DateTimeFormat("ja-JP", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    timeZone: "UTC",
  }).format(new Date(value));
}
