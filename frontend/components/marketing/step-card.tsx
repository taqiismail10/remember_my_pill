"use client";

import { motion } from "framer-motion";
import { ReactNode } from "react";
import { cn } from "@/lib/utils";

type StepCardProps = {
  index: number;
  step: string;
  title: string;
  // A rendered icon element — component references can't cross the
  // server/client boundary.
  icon: ReactNode;
  tone: "terracotta" | "amber" | "green";
  children: ReactNode;
};

const toneClasses = {
  terracotta: "bg-terracotta text-white",
  amber: "bg-amber text-white",
  green: "bg-green text-white",
} as const;

export function StepCard({
  index,
  step,
  title,
  icon,
  tone,
  children,
}: StepCardProps) {
  return (
    <motion.article
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{
        duration: 0.4,
        delay: index * 0.08,
        ease: [0.22, 1, 0.36, 1],
      }}
      className="relative flex h-full flex-col gap-4 rounded-card border border-border bg-surface p-6"
    >
      <div className="flex items-center gap-3">
        <span
          className={cn(
            "flex size-12 shrink-0 items-center justify-center rounded-2xl [&>svg]:size-6",
            toneClasses[tone],
          )}
          aria-hidden="true"
        >
          {icon}
        </span>
        <span className="text-caption uppercase text-muted">Step {step}</span>
      </div>
      <h3 className="text-h3 text-ink">{title}</h3>
      <p className="text-body text-muted">{children}</p>
    </motion.article>
  );
}
