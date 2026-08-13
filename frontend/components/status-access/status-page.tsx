"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Loader2 } from "lucide-react";
import { BrandButton } from "@/components/ui/button";

const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
type State = "loading" | "unavailable" | "unauthorized" | "authenticated";
type Status = {
  rank: number;
  referralCount: number;
  referralCode?: string;
  referralUrl?: string;
};

export function StatusPage({ enabled = false }: { enabled?: boolean }) {
  const [state, setState] = useState<State>(
    enabled ? "loading" : "unavailable",
  );
  const [status, setStatus] = useState<Status>();
  useEffect(() => {
    if (!enabled) return;
    void fetch(`${apiBase}/api/waitlist/status`, { credentials: "include" })
      .then(async (response) => {
        if (response.status === 404) return setState("unavailable");
        if (!response.ok) return setState("unauthorized");
        setStatus((await response.json()) as Status);
        setState("authenticated");
      })
      .catch(() => setState("unauthorized"));
  }, [enabled]);
  async function logout() {
    await fetch(`${apiBase}/api/waitlist/status/logout`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: "{}",
    });
    setState("unauthorized");
  }
  if (state === "loading")
    return (
      <Shell title="Checking your status…">
        <Loader2
          className="mx-auto size-7 animate-spin text-green"
          aria-hidden="true"
        />
      </Shell>
    );
  if (state === "unavailable")
    return (
      <Shell title="Status access is not available yet.">
        <p>
          This pre-launch feature remains unavailable while final approval is
          pending.
        </p>
      </Shell>
    );
  if (state === "unauthorized")
    return (
      <Shell title="Status access needs a new link.">
        <p>
          Your status session has expired or is unavailable. Request a new
          access link when the feature is available.
        </p>
      </Shell>
    );
  return (
    <Shell title="Your waitlist status">
      <dl className="space-y-3 text-left">
        <div>
          <dt className="text-small text-muted">Waitlist position</dt>
          <dd className="text-h3 text-ink">{status?.rank}</dd>
        </div>
        <div>
          <dt className="text-small text-muted">Successful referrals</dt>
          <dd className="text-h3 text-ink">{status?.referralCount}</dd>
        </div>
        {status?.referralCode && (
          <div>
            <dt className="text-small text-muted">Referral code</dt>
            <dd className="text-body font-bold text-ink">
              {status.referralCode}
            </dd>
          </div>
        )}
      </dl>
      <BrandButton
        type="button"
        variant="outline"
        className="mt-7"
        onClick={() => void logout()}
      >
        Sign out
      </BrandButton>
    </Shell>
  );
}

function Shell({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <main className="flex min-h-dvh items-center justify-center bg-cream px-6 py-10 text-center">
      <div className="w-full max-w-md rounded-panel bg-surface p-8 shadow-soft">
        <h1 className="text-h2 text-ink">{title}</h1>
        <div className="mt-4 text-body text-muted">{children}</div>
        <Link
          href="/waitlist"
          className="mt-7 inline-block text-small font-bold text-terracotta-deep underline"
        >
          Back to waitlist
        </Link>
      </div>
    </main>
  );
}
