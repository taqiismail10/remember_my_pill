import { cn } from "@/lib/utils";
import { Eyebrow } from "@/components/marketing/eyebrow";

type SectionHeadingProps = {
  id?: string;
  eyebrow: string;
  title: string;
  align?: "left" | "center";
  className?: string;
  children?: React.ReactNode;
};

export function SectionHeading({
  id,
  eyebrow,
  title,
  align = "left",
  className,
  children,
}: SectionHeadingProps) {
  return (
    <div
      className={cn(
        "max-w-(--container-reading)",
        align === "center" && "mx-auto text-center",
        className,
      )}
    >
      <Eyebrow>{eyebrow}</Eyebrow>
      <h2 id={id} className="mt-3 text-h2 text-balance text-ink">
        {title}
      </h2>
      {children ? (
        <div className="mt-4 text-body-lg text-muted">{children}</div>
      ) : null}
    </div>
  );
}
