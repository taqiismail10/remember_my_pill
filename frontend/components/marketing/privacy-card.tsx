import { LucideIcon } from "lucide-react";
import { cn } from "@/lib/utils";

type PrivacyCardProps = {
  icon: LucideIcon;
  title: string;
  tone: "green" | "terracotta";
  items: string[];
};

export function PrivacyCard({
  icon: Icon,
  title,
  tone,
  items,
}: PrivacyCardProps) {
  return (
    <div className="rounded-panel border border-border bg-surface p-7">
      <span
        className={cn(
          "flex size-11 items-center justify-center rounded-2xl",
          tone === "green"
            ? "bg-green/10 text-green-deep"
            : "bg-terracotta/10 text-terracotta-deep",
        )}
      >
        <Icon className="size-5" aria-hidden="true" />
      </span>
      <h3 className="mt-4 text-h4 text-ink">{title}</h3>
      <ul className="mt-4 space-y-2.5">
        {items.map((item) => (
          <li
            key={item}
            className="flex items-start gap-2.5 text-body text-muted"
          >
            <span
              className={cn(
                "mt-2 size-1.5 shrink-0 rounded-full",
                tone === "green" ? "bg-green" : "bg-terracotta",
              )}
              aria-hidden="true"
            />
            {item}
          </li>
        ))}
      </ul>
    </div>
  );
}
