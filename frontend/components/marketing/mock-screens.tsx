import {
  Bell,
  Calendar,
  Camera,
  Check,
  ChevronLeft,
  Heart,
  Home,
  Plus,
} from "lucide-react";
import { cn } from "@/lib/utils";

/**
 * Original, in-brand illustrations of the planned mobile screens described in
 * the PRD (10.1–10.4). These are product-preview mockups, not real
 * screenshots — every screen using them is captioned as a preview.
 */

function StatusBar() {
  return (
    <div className="flex items-center justify-between px-5 pt-3 text-[11px] font-bold text-ink/70">
      <span>9:41</span>
      <span>●●●</span>
    </div>
  );
}

function BottomNav() {
  return (
    <div className="absolute inset-x-0 bottom-0 flex items-center justify-between bg-surface px-6 py-3">
      <Home className="size-5 text-green" aria-hidden="true" />
      <Calendar className="size-5 text-ink/30" aria-hidden="true" />
      <span className="-mt-6 flex size-11 items-center justify-center rounded-full bg-coral text-white shadow-soft">
        <Plus className="size-5" aria-hidden="true" />
      </span>
      <Bell className="size-5 text-ink/30" aria-hidden="true" />
      <Heart className="size-5 text-ink/30" aria-hidden="true" />
    </div>
  );
}

function DoseRow({
  name,
  dose,
  time,
  state,
}: {
  name: string;
  dose: string;
  time: string;
  state: "done" | "due";
}) {
  return (
    <div className="flex items-center justify-between rounded-2xl bg-surface px-3.5 py-3 shadow-soft">
      <div>
        <p className="text-small font-bold text-ink">{name}</p>
        <p className="text-[11px] text-muted">
          {dose} · {time}
        </p>
      </div>
      <span
        className={cn(
          "flex size-6 items-center justify-center rounded-full border-2",
          state === "done"
            ? "border-green bg-green text-white"
            : "border-ink/15",
        )}
        aria-hidden="true"
      >
        {state === "done" && <Check className="size-3.5" strokeWidth={3} />}
      </span>
    </div>
  );
}

export function TodayScreen() {
  return (
    <div className="flex h-full flex-col gap-4 px-4 pb-20">
      <StatusBar />
      <div>
        <p className="text-h4 text-ink">Good morning, Alex</p>
        <p className="text-[11px] text-muted">Here is today&apos;s plan</p>
      </div>
      <div className="flex gap-2">
        {["Alex", "Mom"].map((person, i) => (
          <span
            key={person}
            className={cn(
              "rounded-full px-3 py-1 text-[11px] font-bold",
              i === 0 ? "bg-green text-white" : "bg-surface text-ink/70",
            )}
          >
            {person}
          </span>
        ))}
      </div>
      <div className="space-y-2">
        <p className="text-[11px] font-bold uppercase text-muted">Morning</p>
        <DoseRow name="Metformin" dose="500mg" time="8:00 AM" state="done" />
        <DoseRow name="Lisinopril" dose="10mg" time="8:00 AM" state="due" />
      </div>
      <div className="space-y-2">
        <p className="text-[11px] font-bold uppercase text-muted">Evening</p>
        <DoseRow name="Vitamin D3" dose="1 tablet" time="8:00 PM" state="due" />
      </div>
      <BottomNav />
    </div>
  );
}

export function AddMedicationScreen() {
  return (
    <div className="flex h-full flex-col gap-5 px-4 pb-6">
      <StatusBar />
      <div className="mt-1 flex items-center gap-3">
        <ChevronLeft className="size-5 text-ink/60" aria-hidden="true" />
        <p className="text-h4 text-ink">Add medication</p>
      </div>
      <span className="w-fit rounded-full bg-surface px-3 py-1 text-[11px] font-bold text-ink/70">
        For: Alex
      </span>
      <div className="grid grid-cols-2 gap-1 rounded-full bg-surface p-1">
        <span className="rounded-full bg-terracotta py-2 text-center text-[11px] font-bold text-white">
          Scan prescription
        </span>
        <span className="rounded-full py-2 text-center text-[11px] font-bold text-ink/60">
          Add manually
        </span>
      </div>
      <div className="flex flex-1 flex-col items-center justify-center gap-3 rounded-3xl border-2 border-dashed border-terracotta/40 bg-surface/60 text-center">
        <Camera className="size-9 text-terracotta" aria-hidden="true" />
        <p className="max-w-[16ch] text-[11px] text-muted">
          Center the prescription label in the frame
        </p>
      </div>
      <span className="rounded-full bg-green py-3 text-center text-small font-bold text-white">
        Capture
      </span>
    </div>
  );
}

export function ReminderScreen() {
  return (
    <div className="flex h-full flex-col gap-5 bg-green-deep px-5 py-8 text-white">
      <div className="flex items-center justify-between text-white/70">
        <span className="text-[11px] font-bold">9:41</span>
        <span className="text-[11px] font-bold">✕</span>
      </div>
      <div className="text-center">
        <p className="text-[11px] font-bold uppercase text-white/60">Due now</p>
        <p className="mt-1 text-h2 text-white">8:00 AM</p>
        <p className="mt-1 text-small text-white/70">2 medicines for Alex</p>
      </div>
      <div className="space-y-2">
        <div className="rounded-2xl bg-white/10 px-4 py-3">
          <p className="text-small font-bold">Metformin 500mg</p>
          <p className="text-[11px] text-white/60">1 tablet</p>
        </div>
        <div className="rounded-2xl bg-white/10 px-4 py-3">
          <p className="text-small font-bold">Lisinopril 10mg</p>
          <p className="text-[11px] text-white/60">1 tablet</p>
        </div>
      </div>
      <div className="mt-auto space-y-2">
        <span className="block rounded-full bg-white py-3 text-center text-small font-bold text-green-deep">
          Mark as taken
        </span>
        <span className="block rounded-full py-2 text-center text-small font-bold text-white/70">
          Snooze 10 minutes
        </span>
      </div>
    </div>
  );
}
