"use client";

import { FormEvent, useId, useState } from "react";
import Link from "next/link";
import { CheckCircle2, Loader2, ShieldCheck } from "lucide-react";
import { Mascot } from "@/components/brand/mascot";
import { BrandButton } from "@/components/ui/button";
import { Eyebrow } from "@/components/marketing/eyebrow";
import { cn } from "@/lib/utils";

const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

type Status = "idle" | "loading" | "success" | "error";

const errorMessages: Record<string, string> = {
  WAITLIST_VALIDATION_ERROR: "Enter a valid email address.",
  WAITLIST_RATE_LIMITED: "Please wait a moment before trying again.",
  WAITLIST_INTERNAL_ERROR:
    "We couldn't join the waitlist right now. Please try again.",
};

/**
 * Email-only waitlist signup for the dedicated /waitlist page. Sends only
 * `{ email }` — the API's `name` field is optional, so this never fakes a
 * placeholder name.
 */
export function WaitlistEmailForm() {
  const [status, setStatus] = useState<Status>("idle");
  const [message, setMessage] = useState("");
  const [email, setEmail] = useState("");
  const [invalid, setInvalid] = useState(false);
  const [waitlistConsent, setWaitlistConsent] = useState(false);
  const [marketingConsent, setMarketingConsent] = useState(false);
  const [consentInvalid, setConsentInvalid] = useState(false);
  const emailId = useId();
  const statusId = useId();

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (status === "loading") return;

    // Client-side honeypot. A populated value does not submit from this form;
    // the API also accepts this documented field for server-side bot handling.
    const honeypot = (
      event.currentTarget.elements.namedItem("company") as HTMLInputElement
    )?.value;
    if (honeypot) return;

    const normalizedEmail = email.trim().toLowerCase();
    if (!/^\S+@\S+\.\S+$/.test(normalizedEmail)) {
      setStatus("error");
      setInvalid(true);
      setMessage("Enter a valid email address.");
      return;
    }
    if (!waitlistConsent) {
      setStatus("error");
      setConsentInvalid(true);
      setMessage("Please confirm the required waitlist consent.");
      return;
    }

    setInvalid(false);
    setConsentInvalid(false);
    setStatus("loading");
    setMessage("Joining the waitlist…");
    try {
      const response = await fetch(`${apiBase}/api/waitlist`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: normalizedEmail,
          consent: true,
          consentVersion: "waitlist-consent-v1",
          marketingConsent,
          marketingConsentVersion: marketingConsent
            ? "marketing-consent-v1"
            : null,
        }),
      });
      if (response.status === 202) {
        setStatus("success");
        setMessage(
          "Thanks — if this email is eligible, your waitlist request has been received.",
        );
        return;
      }
      const body = await response.json().catch(() => null);
      const code = body?.error?.code as string | undefined;
      setStatus("error");
      setMessage(
        (code && errorMessages[code]) ??
          body?.error?.message ??
          "We couldn't join the waitlist. Please try again.",
      );
    } catch {
      setStatus("error");
      setMessage(
        "We couldn't reach the waitlist service. Check your connection and try again.",
      );
    }
  }

  if (status === "success") {
    return (
      <div className="flex flex-col items-center gap-3 text-center animate-[fade-up_0.5s_ease-calm_both]">
        <span className="flex size-14 items-center justify-center rounded-full bg-green/10 text-green-deep">
          <CheckCircle2 className="size-8" aria-hidden="true" />
        </span>
        <h1 className="text-h2 text-ink">You&apos;re on the list.</h1>
        <p role="status" className="max-w-sm text-body text-muted">
          {message}
        </p>
        <Mascot size={64} className="mt-1" />
      </div>
    );
  }

  const isLoading = status === "loading";

  return (
    <div className="w-full text-center">
      <Eyebrow>Early access</Eyebrow>
      <h1 className="mt-3 text-h2 text-balance text-ink">
        Be first to know when Remember My Pill is ready.
      </h1>
      <p className="mt-3 text-body text-muted">
        Join the launch waitlist for the USA + Canada adult pilot. This is a
        pilot draft; launch terms and legal details are not yet final.
      </p>

      <form onSubmit={submit} noValidate className="mt-6 text-left">
        <div className="hidden" aria-hidden="true">
          <label htmlFor="company">Company</label>
          <input
            id="company"
            name="company"
            type="text"
            tabIndex={-1}
            autoComplete="off"
          />
        </div>

        <label htmlFor={emailId} className="text-label text-ink">
          Email
        </label>
        <input
          id={emailId}
          name="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          maxLength={255}
          type="email"
          inputMode="email"
          autoComplete="email"
          placeholder="you@example.com"
          required
          aria-invalid={invalid}
          aria-describedby={statusId}
          disabled={isLoading}
          className={cn(
            "mt-1.5 min-h-12 w-full rounded-control border bg-surface px-4 text-body text-ink outline-none transition-colors",
            invalid ? "border-error" : "border-border focus:border-terracotta",
          )}
        />

        <fieldset
          className="mt-4 space-y-3"
          aria-describedby={consentInvalid ? statusId : undefined}
        >
          <legend className="sr-only">Email consents</legend>
          <label className="flex cursor-pointer items-start gap-3 rounded-control border border-border bg-surface p-3.5 text-small text-ink">
            <input
              type="checkbox"
              checked={waitlistConsent}
              onChange={(event) => setWaitlistConsent(event.target.checked)}
              disabled={isLoading}
              aria-invalid={consentInvalid}
              className="mt-0.5 size-4 shrink-0 accent-green"
            />
            <span>
              I agree to join the Remember My Pill waitlist and receive emails
              necessary to manage my waitlist participation and notify me about
              the RMP pilot and launch. I understand that I can withdraw from
              the waitlist at any time. See the{" "}
              <Link
                href="/privacy"
                className="font-bold text-terracotta-deep underline"
              >
                Privacy Policy
              </Link>
              .
            </span>
          </label>
          <label className="flex cursor-pointer items-start gap-3 rounded-control border border-border bg-surface p-3.5 text-small text-ink">
            <input
              type="checkbox"
              checked={marketingConsent}
              onChange={(event) => setMarketingConsent(event.target.checked)}
              disabled={isLoading}
              className="mt-0.5 size-4 shrink-0 accent-green"
            />
            <span>
              I would also like to receive Remember My Pill product news,
              feature updates, surveys, educational content, promotions, and
              special offers by email. I understand that I can unsubscribe from
              marketing emails at any time.
            </span>
          </label>
        </fieldset>

        <BrandButton
          type="submit"
          variant="primary"
          size="lg"
          disabled={isLoading}
          className="mt-4 w-full"
        >
          {isLoading && (
            <Loader2 className="size-4 animate-spin" aria-hidden="true" />
          )}
          {isLoading ? "Joining…" : "Join the waitlist"}
        </BrandButton>

        <p
          id={statusId}
          role="status"
          aria-live="polite"
          className={cn(
            "mt-3 min-h-5 text-center text-small",
            status === "error" && "font-bold text-error",
          )}
        >
          {message}
        </p>
      </form>

      <p className="mt-4 flex items-start justify-center gap-2 text-left text-small text-muted">
        <ShieldCheck
          className="mt-0.5 size-4 shrink-0 text-green"
          aria-hidden="true"
        />
        The waitlist stores only your email and consent choices. Do not submit
        prescriptions, medication details, or other health information here.
      </p>
    </div>
  );
}
