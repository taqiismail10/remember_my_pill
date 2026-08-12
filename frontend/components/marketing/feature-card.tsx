"use client";

import { motion } from "framer-motion";
import { ReactNode } from "react";
import { cn } from "@/lib/utils";

type FeatureCardProps = {
  // A rendered icon element (e.g. <FileQuestion />), not a component
  // reference — component functions can't cross the server/client boundary.
  icon: ReactNode;
  title: string;
  tone?: "terracotta" | "green";
  delay?: number;
  className?: string;
  children: ReactNode;
};

export function FeatureCard({
  icon,
  title,
  tone = "terracotta",
  delay = 0,
  className,
  children,
}: FeatureCardProps) {
  return (
    <motion.article
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, delay, ease: [0.22, 1, 0.36, 1] }}
      className={cn(
        "flex h-full flex-col gap-3 rounded-card border border-border bg-cream-soft p-6",
        className,
      )}
    >
      <span
        className={cn(
          "flex size-10 items-center justify-center rounded-xl [&>svg]:size-5",
          tone === "terracotta"
            ? "bg-terracotta/10 text-terracotta-deep"
            : "bg-green/10 text-green-deep",
        )}
        aria-hidden="true"
      >
        {icon}
      </span>
      <h3 className="text-h4 text-ink">{title}</h3>
      <p className="text-body text-muted">{children}</p>
    </motion.article>
  );
}
