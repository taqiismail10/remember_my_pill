"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { ArrowLeft, Loader2 } from "lucide-react";
import { BrandButton } from "@/components/ui/button";

const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

type VerifyState = "loading" | "unavailable" | "invalid";

export function StatusVerify({ enabled = false }: { enabled?: boolean }) {
  const [state, setState] = useState<VerifyState>(() =>
    enabled ? "loading" : "unavailable",
  );

  useEffect(() => {
    const fragment = new URLSearchParams(window.location.hash.slice(1));
    const token = fragment.get("v");
    // Remove the fragment before any network work or user navigation. The raw
    // value remains only in this ephemeral closure while the POST is pending.
    window.history.replaceState(null, "", "/waitlist/verify");
    if (!enabled) {
      return;
    }
    if (!token) {
      queueMicrotask(() => setState("invalid"));
      return;
    }
    void fetch(`${apiBase}/api/waitlist/status-access/exchange`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ verificationToken: token }),
    })
      .then((response) => {
        if (!response.ok) throw new Error("verification failed");
        window.location.replace("/waitlist/status");
      })
      .catch(() => setState("invalid"));
  }, [enabled]);

  if (state === "loading") {
    return (
      <StatusShell title="Verifying your access…">
        <Loader2
          className="mx-auto size-7 animate-spin text-green"
          aria-hidden="true"
        />
      </StatusShell>
    );
  }
  if (state === "unavailable") {
    return (
      <StatusShell title="Status access is not available yet.">
        <p>
          This pre-launch feature remains unavailable while final approval is
          pending.
        </p>
      </StatusShell>
    );
  }
  return (
    <StatusShell title="We couldn’t verify that link.">
      <p>
        This link may have expired or already been used. You can try again once
        status access is available.
      </p>
    </StatusShell>
  );
}

function StatusShell({
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
        <BrandButton
          as={Link}
          href="/waitlist"
          variant="outline"
          className="mt-7"
        >
          <ArrowLeft className="size-4" aria-hidden="true" />
          Back to waitlist
        </BrandButton>
      </div>
    </main>
  );
}
