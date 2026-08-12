import Link from "next/link";
import { Mascot } from "@/components/brand/mascot";
import { BrandButton } from "@/components/ui/button";
import { Section } from "@/components/marketing/section";

export function FinalCta() {
  return (
    <Section tone="green" padding="tight" aria-labelledby="final-cta-title">
      <div className="flex flex-col items-center gap-6 text-center">
        <Mascot size={96} className="bg-white/10" />
        <h2 id="final-cta-title" className="max-w-lg text-h2 text-white">
          Helping you remember what matters.
        </h2>
        <p className="max-w-md text-body-lg text-white/80">
          Join the waitlist to hear from us the moment early access opens.
        </p>
        <BrandButton as={Link} href="/waitlist" variant="warm" size="lg">
          Join the waitlist
        </BrandButton>
      </div>
    </Section>
  );
}
