import type { Metadata } from "next";
import Link from "next/link";
import { ArrowLeft, HeartPulse } from "lucide-react";
import { BrandLink } from "@/components/brand/brand-link";

export const metadata: Metadata = {
  title: "Consumer Health Data Privacy — Pilot Draft",
  description: "Remember My Pill consumer health data pilot draft.",
};

export default function ConsumerHealthDataPrivacyPage() {
  return (
    <main className="min-h-dvh bg-cream px-6 py-10">
      <div className="mx-auto max-w-3xl">
        <BrandLink href="/" priority />
        <p className="mt-10 inline-flex items-center gap-2 rounded-full bg-amber/20 px-3 py-1 text-caption uppercase text-ink">
          <HeartPulse className="size-3.5" aria-hidden="true" />
          Pilot Draft — not production-ready
        </p>
        <h1 className="mt-4 text-h1 text-ink">Consumer health data privacy</h1>
        <div className="mt-8 space-y-6 text-body text-muted">
          <section className="rounded-card border border-border bg-surface p-6 shadow-soft">
            <h2 className="text-h3 text-ink">Sensitive pilot information</h2>
            <p className="mt-3">
              For the pilot, medication names, dosage information,
              prescription-related information, and medication schedules are
              treated as sensitive health-related information. Removing direct
              identifiers does not make those payloads anonymous.
            </p>
          </section>
          <section className="rounded-card border border-border bg-surface p-6 shadow-soft">
            <h2 className="text-h3 text-ink">
              What the pilot is designed to do
            </h2>
            <p className="mt-3">
              The pilot is designed for medication information organization,
              extraction, user confirmation or correction, reminder setup, and
              reminder notifications. It does not provide diagnosis,
              prescribing, medication-change recommendations, or a replacement
              for a doctor or pharmacist.
            </p>
          </section>
          <section className="rounded-card border border-border bg-surface p-6 shadow-soft">
            <h2 className="text-h3 text-ink">Draft status</h2>
            <p className="mt-3">
              The legal operator, business mailing address, and final external
              legal review are unresolved. This is not a final production notice
              and must be completed before a public pilot launch.
            </p>
          </section>
        </div>
        <p className="mt-8 text-small text-muted">
          Read the{" "}
          <Link
            href="/privacy"
            className="font-bold text-terracotta-deep underline"
          >
            Privacy draft
          </Link>
          .
        </p>
        <Link
          href="/"
          className="mt-8 inline-flex items-center gap-2 text-small font-bold text-terracotta-deep"
        >
          <ArrowLeft className="size-4" aria-hidden="true" />
          Back to home
        </Link>
      </div>
    </main>
  );
}
