import type { Metadata } from "next";
import { StatusPage } from "@/components/status-access/status-page";

export const metadata: Metadata = {
  title: "Waitlist status",
  robots: { index: false, follow: false },
};

export default function WaitlistStatusPage() {
  return (
    <StatusPage
      enabled={process.env.NEXT_PUBLIC_STATUS_ACCESS_ENABLED === "true"}
    />
  );
}
