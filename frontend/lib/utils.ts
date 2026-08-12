import { type ClassValue, clsx } from "clsx";
import { extendTailwindMerge } from "tailwind-merge";

// tailwind-merge doesn't know about the custom `--text-*` type scale defined
// in app/globals.css (@theme). Without this, e.g. `cn("text-caption",
// "text-terracotta-deep")` silently drops `text-caption` because tailwind-
// merge guesses it belongs to the text-color group, not font-size.
const twMerge = extendTailwindMerge({
  extend: {
    theme: {
      text: [
        "display",
        "h1",
        "h2",
        "h3",
        "h4",
        "body-lg",
        "body",
        "small",
        "caption",
        "button",
        "label",
      ],
    },
  },
});

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
