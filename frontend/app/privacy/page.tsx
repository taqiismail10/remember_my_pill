import type { Metadata } from "next";
import Link from "next/link";
import { ArrowLeft, ShieldCheck } from "lucide-react";
import { BrandLink } from "@/components/brand/brand-link";

export const metadata: Metadata = {
  title: "Privacy — Pilot Draft",
  description: "Remember My Pill USA and Canada pilot privacy draft.",
};

const sections = [
  [
    "Pilot draft status",
    "This privacy page is a pre-release pilot draft for adults 18+ in the United States and Canada. It is not a final production privacy policy. The legal operator, business mailing address, and final external legal review remain unresolved before any public launch.",
  ],
  [
    "Pilot scope",
    "Remember My Pill is being prepared to help people organize medication information, extract prescription-related information, confirm or correct that information, set up reminders, and receive reminder notifications. It does not diagnose, prescribe, recommend medication changes, or replace doctors or pharmacists.",
  ],
  [
    "Waitlist information",
    "The waitlist asks for an email address and records the required waitlist-consent choice. Optional marketing consent is separate. The waitlist does not ask for prescriptions, medication names, dosage information, diagnoses, insurance information, or medication schedules.",
  ],
  [
    "Contact and launch blockers",
    "A privacy contact and formal legal identity will be published only after owner and legal approval. Do not rely on this draft as final legal notice. A final Terms of Use page is also required before public launch and is not yet available.",
  ],
];

export default function PrivacyPage() {
  return (
    <main className="min-h-dvh bg-cream px-6 py-10">
      <div className="mx-auto max-w-3xl">
        <BrandLink href="/" priority />
        <p className="mt-10 inline-flex items-center gap-2 rounded-full bg-amber/20 px-3 py-1 text-caption uppercase text-ink">
          <ShieldCheck className="size-3.5" aria-hidden="true" />
          Pilot Draft — not production-ready
        </p>
        <h1 className="mt-4 text-h1 text-ink">
          Privacy for the USA + Canada pilot
        </h1>
        <p className="mt-4 max-w-reading text-body-lg text-muted">
          A plain-language draft for the pre-release Remember My Pill pilot.
        </p>
        <div className="mt-10 space-y-6">
          {sections.map(([title, body]) => (
            <section
              key={title}
              className="rounded-card border border-border bg-surface p-6 shadow-soft"
            >
              <h2 className="text-h3 text-ink">{title}</h2>
              <p className="mt-3 text-body text-muted">{body}</p>
            </section>
          ))}
        </div>
        <p className="mt-8 text-small text-muted">
          For medication-related information, see the{" "}
          <Link
            href="/consumer-health-data-privacy"
            className="font-bold text-terracotta-deep underline"
          >
            Consumer Health Data Privacy draft
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
