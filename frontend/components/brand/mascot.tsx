import Image from "next/image";
import { cn } from "@/lib/utils";

type MascotProps = {
  size?: number;
  alt?: string;
  className?: string;
  priority?: boolean;
};

/** Official 2D mascot artwork from the RMP Brand Identity file. */
export function Mascot({
  size = 160,
  alt = "",
  className,
  priority,
}: MascotProps) {
  return (
    <Image
      src="/brand/mascot/mascot-wave.png"
      alt={alt}
      width={size}
      height={size}
      priority={priority}
      className={cn("rounded-full", className)}
    />
  );
}
