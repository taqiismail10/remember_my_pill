import Link from "next/link";
import { ShieldCheck } from "lucide-react";
import { Eyebrow } from "@/components/marketing/eyebrow";
import { PhoneFrame } from "@/components/marketing/phone-frame";
import { TodayScreen } from "@/components/marketing/mock-screens";
import { BrandBadge } from "@/components/ui/badge";
import { BrandButton } from "@/components/ui/button";

export function Hero() {
  return (
    <section
      className="relative overflow-hidden section-pad"
      aria-labelledby="hero-title"
    >
      <div className="pointer-events-none absolute -right-24 -top-24 z-0 size-[32rem] rounded-full bg-terracotta/5 blur-3xl" />
      <div className="container relative z-10 grid items-center gap-12 lg:grid-cols-[1.1fr_0.9fr] lg:gap-8">
        <div>
          <div className="animate-[fade-up_0.6s_ease-calm_both]">
            <Eyebrow>Medication support, thoughtfully designed</Eyebrow>
          </div>
          <h1
            id="hero-title"
            className="mt-4 max-w-xl text-balance text-display text-ink animate-[fade-up_0.6s_ease-calm_both] [animation-delay:80ms]"
          >
            Turn confusing prescriptions into a routine you can trust.
          </h1>
          <p className="mt-5 max-w-lg text-body-lg text-muted animate-[fade-up_0.6s_ease-calm_both] [animation-delay:160ms]">
            Remember My Pill helps you snap a prescription, understand it in
            plain language, and act on a clear daily routine — designed for the
            people who take medicine, and the people who help them.
          </p>
          <div className="mt-8 flex flex-wrap items-center gap-3 animate-[fade-up_0.6s_ease-calm_both] [animation-delay:240ms]">
            <BrandButton as={Link} href="/waitlist" variant="primary" size="lg">
              Join the waitlist
            </BrandButton>
            <BrandButton
              as="a"
              href="#how-it-works"
              variant="outline"
              size="lg"
            >
              See how it works
            </BrandButton>
          </div>
          <p className="mt-6 flex items-start gap-2 text-small text-muted animate-[fade-up_0.6s_ease-calm_both] [animation-delay:320ms]">
            <ShieldCheck
              className="mt-0.5 size-4 shrink-0 text-green"
              aria-hidden="true"
            />
            Product preview — this site does not collect prescriptions,
            medication details, or health information.
          </p>
        </div>
        <div className="animate-[fade-up_0.7s_ease-calm_both] [animation-delay:200ms]">
          <div className="mb-5 flex justify-center">
            <BrandBadge tone="cream" className="shadow-soft">
              Product preview
            </BrandBadge>
          </div>
          <PhoneFrame className="max-w-[280px]">
            <TodayScreen />
          </PhoneFrame>
        </div>
      </div>
    </section>
  );
}
