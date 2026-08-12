"use client";

import { useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { PhoneFrame } from "@/components/marketing/phone-frame";
import {
  AddMedicationScreen,
  ReminderScreen,
  TodayScreen,
} from "@/components/marketing/mock-screens";
import { cn } from "@/lib/utils";

const tabs = [
  { id: "today", label: "Today" },
  { id: "add", label: "Add med." },
  { id: "reminder", label: "Reminder" },
] as const;

type TabId = (typeof tabs)[number]["id"];

const fade = {
  initial: { opacity: 0, y: 8 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -8 },
};
const transition = { duration: 0.25, ease: [0.22, 1, 0.36, 1] as const };

export function ProductPreview() {
  const [active, setActive] = useState<TabId>("today");

  return (
    <div className="mx-auto flex max-w-sm flex-col items-center">
      <div className="mb-8 max-w-full overflow-x-auto">
        <div
          role="tablist"
          aria-label="Product preview screens"
          className="inline-flex gap-1 rounded-full bg-cream-soft p-1"
        >
          {tabs.map((tab) => (
            <button
              key={tab.id}
              role="tab"
              id={`tab-${tab.id}`}
              aria-selected={active === tab.id}
              aria-controls={`panel-${tab.id}`}
              onClick={() => setActive(tab.id)}
              className={cn(
                "min-h-11 whitespace-nowrap rounded-full px-3 text-small font-bold transition-colors duration-200 ease-calm sm:px-4",
                active === tab.id
                  ? "bg-terracotta text-white"
                  : "text-muted hover:text-ink",
              )}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      <PhoneFrame>
        <div
          role="tabpanel"
          id={`panel-${active}`}
          aria-labelledby={`tab-${active}`}
          className="relative h-full w-full"
        >
          <AnimatePresence mode="wait">
            {active === "today" && (
              <motion.div key="today" {...fade} transition={transition}>
                <TodayScreen />
              </motion.div>
            )}
            {active === "add" && (
              <motion.div key="add" {...fade} transition={transition}>
                <AddMedicationScreen />
              </motion.div>
            )}
            {active === "reminder" && (
              <motion.div key="reminder" {...fade} transition={transition}>
                <ReminderScreen />
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </PhoneFrame>
      <p className="mt-6 max-w-xs text-center text-small text-muted">
        Illustrative preview of the mobile product design. Not a live connection
        to prescriptions or health data.
      </p>
    </div>
  );
}
