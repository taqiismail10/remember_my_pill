import { ReactNode } from "react";
import { cn } from "@/lib/utils";

const tones = {
  terracotta: "bg-terracotta/10 text-terracotta-deep",
  green: "bg-green/10 text-green-deep",
  cream: "bg-cream-soft text-ink",
} as const;

type BrandBadgeProps = {
  tone?: keyof typeof tones;
  icon?: ReactNode;
  className?: string;
  children: ReactNode;
};

export function BrandBadge({
  tone = "terracotta",
  icon,
  className,
  children,
}: BrandBadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-label",
        tones[tone],
        className,
      )}
    >
      {icon}
      {children}
    </span>
  );
}
