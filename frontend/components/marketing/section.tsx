import { ReactNode } from "react";
import { cn } from "@/lib/utils";
import { Container } from "@/components/layout/container";

const tones = {
  cream: "bg-cream",
  white: "bg-surface",
  soft: "bg-cream-soft",
  green: "bg-green text-white",
} as const;

type SectionProps = {
  id?: string;
  tone?: keyof typeof tones;
  padding?: "normal" | "tight";
  "aria-labelledby"?: string;
  className?: string;
  containerClassName?: string;
  children: ReactNode;
};

export function Section({
  id,
  tone = "cream",
  padding = "normal",
  className,
  containerClassName,
  children,
  ...aria
}: SectionProps) {
  return (
    <section
      id={id}
      className={cn(
        "relative",
        padding === "tight" ? "section-pad-tight" : "section-pad",
        tones[tone],
        className,
      )}
      {...aria}
    >
      <Container className={cn("relative z-10", containerClassName)}>
        {children}
      </Container>
    </section>
  );
}
