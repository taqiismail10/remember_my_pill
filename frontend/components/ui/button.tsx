import { ComponentPropsWithoutRef, ElementType } from "react";
import { cn } from "@/lib/utils";

const variants = {
  primary: "bg-green text-white hover:bg-green-deep",
  warm: "bg-terracotta text-white hover:bg-terracotta-deep",
  outline:
    "border-2 border-ink/15 bg-transparent text-ink hover:border-ink/30 hover:bg-ink/5",
  ghost: "bg-transparent text-ink hover:bg-ink/5",
} as const;

const sizes = {
  md: "min-h-11 px-5 text-button",
  lg: "min-h-13 px-7 text-button",
} as const;

type BrandButtonProps<T extends ElementType> = {
  as?: T;
  variant?: keyof typeof variants;
  size?: keyof typeof sizes;
} & Omit<ComponentPropsWithoutRef<T>, "as">;

export function BrandButton<T extends ElementType = "button">({
  as,
  variant = "primary",
  size = "md",
  className,
  ...props
}: BrandButtonProps<T>) {
  const Component = as ?? "button";
  return (
    <Component
      className={cn(
        "inline-flex items-center justify-center gap-2 rounded-full font-sans transition-colors duration-200 ease-calm disabled:cursor-not-allowed disabled:opacity-60",
        variants[variant],
        sizes[size],
        className,
      )}
      {...props}
    />
  );
}
