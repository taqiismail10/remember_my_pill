import type { Metadata } from "next";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import { BrandLink } from "@/components/brand/brand-link";
import { WaitlistEmailForm } from "@/components/waitlist/waitlist-email-form";

export const metadata: Metadata = {
  title: "Join the waitlist",
  description:
    "Join the Remember My Pill launch waitlist with just your email. We'll let you know the moment early access opens.",
};

export default function WaitlistPage() {
  return (
    <div className="flex min-h-dvh flex-col bg-cream">
      <div className="flex flex-1 flex-col items-center justify-center px-6 py-10">
        <div className="mb-6 animate-[fade-up_0.5s_ease-calm_both]">
          <BrandLink href="/" priority />
        </div>

        <div className="w-full max-w-md animate-[fade-up_0.6s_ease-calm_both] [animation-delay:80ms]">
          <WaitlistEmailForm />
        </div>

        <Link
          href="/"
          className="mt-8 inline-flex items-center gap-1.5 text-small font-bold text-ink/70 transition-colors hover:text-terracotta-deep animate-[fade-up_0.6s_ease-calm_both] [animation-delay:160ms]"
        >
          <ArrowLeft className="size-4" aria-hidden="true" />
          Back to home
        </Link>
      </div>
    </div>
  );
}
