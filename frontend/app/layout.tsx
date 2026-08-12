import type { Metadata } from "next";
import { Manrope } from "next/font/google";
import { MotionProvider } from "@/components/motion-provider";
import "./globals.css";

const manrope = Manrope({ subsets: ["latin"], variable: "--font-manrope" });

const siteUrl = process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000";

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  title: {
    default: "Remember My Pill | Medication support, thoughtfully designed",
    template: "%s | Remember My Pill",
  },
  description:
    "Remember My Pill is a privacy-first medication support product in development. Snap a prescription, understand it in plain language, and act on a clear daily routine. Join the waitlist for early access.",
  icons: {
    icon: "/brand/r-pill-mark.svg",
  },
  openGraph: {
    type: "website",
    siteName: "Remember My Pill",
    title: "Remember My Pill | Medication support, thoughtfully designed",
    description:
      "A privacy-first medication support product in development. Explore the product preview and join the waitlist for early access.",
    images: ["/product-preview/app-icon.png"],
  },
  twitter: {
    card: "summary_large_image",
    title: "Remember My Pill | Medication support, thoughtfully designed",
    description:
      "A privacy-first medication support product in development. Explore the product preview and join the waitlist for early access.",
    images: ["/product-preview/app-icon.png"],
  },
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body className={`${manrope.variable} font-sans`}>
        <MotionProvider>{children}</MotionProvider>
      </body>
    </html>
  );
}
