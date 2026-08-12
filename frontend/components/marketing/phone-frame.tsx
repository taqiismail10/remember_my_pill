import { ReactNode } from "react";
import { cn } from "@/lib/utils";

type PhoneFrameProps = {
  children: ReactNode;
  className?: string;
};

/** A generic device chrome used to present original, in-brand product-preview screens. */
export function PhoneFrame({ children, className }: PhoneFrameProps) {
  return (
    <div
      className={cn(
        "relative mx-auto w-full max-w-[300px] rounded-device border-[6px] border-ink/90 bg-ink/90 shadow-lifted",
        className,
      )}
    >
      <div className="absolute left-1/2 top-0 z-10 h-5 w-28 -translate-x-1/2 rounded-b-xl bg-ink/90" />
      <div className="aspect-[9/19.5] overflow-hidden rounded-[2.25rem] bg-surface-app">
        {children}
      </div>
    </div>
  );
}
