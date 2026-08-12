import Link from "next/link";
import {
  Bell,
  BellRing,
  Camera,
  ClipboardList,
  Eye,
  FileQuestion,
  Lock,
  ListChecks,
  Mail,
  ShieldCheck,
  ShieldOff,
  Sparkles,
  UserPlus,
  Users,
} from "lucide-react";
import { SiteHeader } from "@/components/layout/site-header";
import { Eyebrow } from "@/components/marketing/eyebrow";
import { SiteFooter } from "@/components/layout/site-footer";
import { Section } from "@/components/marketing/section";
import { SectionHeading } from "@/components/marketing/section-heading";
import { Hero } from "@/components/marketing/hero";
import { StepCard } from "@/components/marketing/step-card";
import { FeatureCard } from "@/components/marketing/feature-card";
import { PrivacyCard } from "@/components/marketing/privacy-card";
import { ProductPreview } from "@/components/marketing/product-preview";
import { FaqAccordion } from "@/components/marketing/faq-accordion";
import { FinalCta } from "@/components/marketing/final-cta";
import { Mascot } from "@/components/brand/mascot";
import { BrandButton } from "@/components/ui/button";

const howItWorks = [
  {
    icon: UserPlus,
    title: "Join the waitlist",
    description: "Share your email — nothing about your health.",
  },
  {
    icon: Mail,
    title: "We keep you posted",
    description: "No rank or referral games yet — just a note when it matters.",
  },
  {
    icon: Sparkles,
    title: "Get early access",
    description: "We'll email you the moment the mobile product opens up.",
  },
];

export default function Home() {
  return (
    <div id="top">
      <SiteHeader />
      <main>
        <Hero />

        <Section tone="white" aria-labelledby="problem-title">
          <SectionHeading
            id="problem-title"
            eyebrow="The problem"
            title="Prescription instructions shouldn't add to the mental load."
          >
            Handwritten notes, clinical shorthand, and multiple schedules can
            make daily routines feel needlessly difficult — especially when
            you&apos;re managing medicine for more than one person.
          </SectionHeading>
          <div className="mt-10 grid gap-5 sm:grid-cols-3">
            <FeatureCard
              icon={<FileQuestion />}
              title="Hard to interpret"
              delay={0}
            >
              Clear information is easier to act on than a prescription full of
              abbreviations and shorthand.
            </FeatureCard>
            <FeatureCard
              icon={<Users />}
              title="Hard to coordinate"
              delay={0.08}
            >
              A routine can involve several medicines and the people who help
              manage them.
            </FeatureCard>
            <FeatureCard icon={<Bell />} title="Hard to remember" delay={0.16}>
              A calm, grouped plan makes the next dose more obvious — for you or
              someone you care for.
            </FeatureCard>
          </div>
        </Section>

        <Section id="how-it-works" aria-labelledby="how-it-works-title">
          <SectionHeading
            id="how-it-works-title"
            eyebrow="How it works"
            title="Snap. Understand. Act."
          >
            The product is organized around three simple steps, with a
            confirmation-first approach to prescription information.
          </SectionHeading>
          <div className="mt-10 grid gap-5 md:grid-cols-3">
            <StepCard
              index={0}
              step="1"
              title="Snap"
              icon={<Camera />}
              tone="terracotta"
            >
              Capture or upload a prescription in the mobile product — no typing
              required.
            </StepCard>
            <StepCard
              index={1}
              step="2"
              title="Understand"
              icon={<Sparkles />}
              tone="amber"
            >
              Review structured, plain-language instructions before anything is
              saved to your routine.
            </StepCard>
            <StepCard
              index={2}
              step="3"
              title="Act"
              icon={<BellRing />}
              tone="green"
            >
              Organize grouped reminders by time of day and mark doses as taken.
            </StepCard>
          </div>
        </Section>

        <Section id="preview" tone="white" aria-labelledby="preview-title">
          <SectionHeading
            id="preview-title"
            eyebrow="Product preview"
            title="A warm, clear foundation for the product ahead."
            align="center"
          >
            Explore the planned home screen, prescription capture, and reminder
            experience.
          </SectionHeading>
          <div className="mt-12">
            <ProductPreview />
          </div>
        </Section>

        <Section id="caregivers" aria-labelledby="caregiver-title">
          <div className="grid items-center gap-10 lg:grid-cols-2">
            <div>
              <SectionHeading
                id="caregiver-title"
                eyebrow="For caregivers"
                title="Built for the people who manage medicine for someone else."
              >
                Whether you&apos;re looking after a parent, a partner, or your
                own routine, Remember My Pill is designed to keep each
                person&apos;s schedule clear and separate.
              </SectionHeading>
              <ul className="mt-6 space-y-3">
                {[
                  "Switch between the people you support from one home screen.",
                  "Keep each person's medicines and schedule grouped on their own.",
                  "See what's due next without digging through paper notes.",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-3 text-body text-muted"
                  >
                    <span
                      className="mt-2 size-1.5 shrink-0 rounded-full bg-green"
                      aria-hidden="true"
                    />
                    {item}
                  </li>
                ))}
              </ul>
            </div>
            <div className="flex justify-center rounded-panel bg-cream-soft p-10">
              <Mascot size={200} />
            </div>
          </div>
        </Section>

        <Section tone="white" aria-labelledby="benefits-title">
          <SectionHeading
            id="benefits-title"
            eyebrow="Why it helps"
            title="Everything is designed around calm, clear routines."
            align="center"
          />
          <div className="mt-10 grid gap-5 md:grid-cols-3">
            <FeatureCard icon={<Sparkles />} title="Simplicity" delay={0}>
              Plain-language instructions instead of clinical shorthand.
            </FeatureCard>
            <FeatureCard
              icon={<ListChecks />}
              title="Organization"
              tone="green"
              delay={0.06}
            >
              Medicines grouped into a daily routine, not a loose list.
            </FeatureCard>
            <FeatureCard icon={<Eye />} title="Clarity" delay={0.12}>
              A confirmation-first review before anything is saved.
            </FeatureCard>
            <FeatureCard
              icon={<Users />}
              title="Caregiver support"
              tone="green"
              delay={0.18}
            >
              Multi-person schedules, kept clearly separated.
            </FeatureCard>
            <FeatureCard icon={<Lock />} title="Privacy" delay={0.24}>
              Data-minimizing by design, from the waitlist onward.
            </FeatureCard>
          </div>
        </Section>

        <Section id="privacy" aria-labelledby="privacy-title">
          <SectionHeading
            id="privacy-title"
            eyebrow="Privacy and trust"
            title="We only ask for what the waitlist actually needs."
            align="center"
          >
            Remember My Pill is built with a data-minimizing, privacy-first
            posture. Final legal review is still required before launch.
          </SectionHeading>
          <div className="mt-10 grid gap-5 md:grid-cols-2">
            <PrivacyCard
              icon={ClipboardList}
              title="What this waitlist collects"
              tone="green"
              items={[
                "Your email, so we can let you know when early access opens.",
              ]}
            />
            <PrivacyCard
              icon={ShieldOff}
              title="What we never collect here"
              tone="terracotta"
              items={[
                "Prescriptions, medication names, or dosages.",
                "Diagnoses, insurance details, or medical record numbers.",
                "Behavioral advertising trackers.",
              ]}
            />
          </div>
        </Section>

        <Section tone="white" aria-labelledby="waitlist-info-title">
          <SectionHeading
            id="waitlist-info-title"
            eyebrow="Joining the waitlist"
            title="What happens after you sign up."
            align="center"
          />
          <div className="mt-10 grid gap-6 sm:grid-cols-3">
            {howItWorks.map((step, index) => (
              <div key={step.title} className="text-center">
                <span className="mx-auto flex size-12 items-center justify-center rounded-2xl bg-terracotta/10 text-terracotta-deep">
                  <step.icon className="size-6" aria-hidden="true" />
                </span>
                <p className="mt-4 text-h4 text-ink">
                  {index + 1}. {step.title}
                </p>
                <p className="mt-2 text-body text-muted">{step.description}</p>
              </div>
            ))}
          </div>
        </Section>

        <Section id="waitlist" aria-labelledby="waitlist-title">
          <div className="mx-auto max-w-xl text-center">
            <Eyebrow>Early access</Eyebrow>
            <h2 id="waitlist-title" className="mt-3 text-h2 text-ink">
              Join the waitlist
            </h2>
            <p className="mt-4 text-body-lg text-muted">
              Be the first to know when Remember My Pill opens for early access.
            </p>
          </div>
          <div className="mx-auto mt-10 max-w-xl rounded-panel bg-surface p-8 text-center shadow-soft sm:p-10">
            <p className="text-body text-muted">
              Just your email — nothing else. Takes a few seconds.
            </p>
            <BrandButton
              as={Link}
              href="/waitlist"
              variant="primary"
              size="lg"
              className="mt-5"
            >
              Join the waitlist
            </BrandButton>
            <p className="mt-4 flex items-start justify-center gap-2 text-left text-small text-muted">
              <ShieldCheck
                className="mt-0.5 size-4 shrink-0 text-green"
                aria-hidden="true"
              />
              We only ask for your email. Please don&apos;t include medication,
              prescription, or other health information — this waitlist never
              collects it.
            </p>
          </div>
        </Section>

        <Section id="faq" tone="white" aria-labelledby="faq-title">
          <SectionHeading
            id="faq-title"
            eyebrow="FAQ"
            title="Good to know"
            align="center"
          />
          <div className="mx-auto mt-8 max-w-2xl">
            <FaqAccordion />
          </div>
        </Section>

        <FinalCta />
      </main>
      <SiteFooter />
    </div>
  );
}
