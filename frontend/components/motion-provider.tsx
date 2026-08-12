"use client";

import { MotionConfig } from "framer-motion";
import { ReactNode } from "react";

/** Ensures every Framer Motion animation in the tree respects prefers-reduced-motion. */
export function MotionProvider({ children }: { children: ReactNode }) {
  return <MotionConfig reducedMotion="user">{children}</MotionConfig>;
}
