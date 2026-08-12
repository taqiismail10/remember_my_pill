import Image from "next/image";
import Link from "next/link";
import { cn } from "@/lib/utils";

/**
 * Reproduces the Figma "Primary Logo" tile: the cream R/pill mark on a
 * terracotta rounded-square, exactly as shown on the RMP Pallete 1 board.
 */
export function BrandLink({
  className,
  priority = false,
  href = "#top",
}: {
  className?: string;
  priority?: boolean;
  /** Same-page anchor (default "#top") or a route like "/" for other pages. */
  href?: string;
}) {
  const classes = cn(
    "inline-flex items-center gap-2.5 text-h4 font-extrabold text-terracotta-deep",
    className,
  );
  const mark = (
    <>
      <span className="flex size-9 shrink-0 items-center justify-center rounded-[0.6rem] bg-terracotta p-[7px]">
        <Image
          src="/brand/r-pill-mark.svg"
          alt=""
          width={32}
          height={37}
          priority={priority}
          className="h-full w-auto"
        />
      </span>
      <span>Remember My Pill</span>
    </>
  );

  if (href.startsWith("#")) {
    return (
      <a href={href} aria-label="Remember My Pill home" className={classes}>
        {mark}
      </a>
    );
  }
  return (
    <Link href={href} aria-label="Remember My Pill home" className={classes}>
      {mark}
    </Link>
  );
}
