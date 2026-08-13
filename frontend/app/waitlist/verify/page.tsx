import type { Metadata } from "next";
import { StatusVerify } from "@/components/status-access/status-verify";

export const metadata: Metadata = {
  title: "Verify status access",
  robots: { index: false, follow: false },
  other: { referrer: "no-referrer" },
};

export default function VerifyStatusPage() {
  return (
    <StatusVerify
      enabled={process.env.NEXT_PUBLIC_STATUS_ACCESS_ENABLED === "true"}
    />
  );
}
