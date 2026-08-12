"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Menu, X } from "lucide-react";
import { BrandLink } from "@/components/brand/brand-link";
import { BrandButton } from "@/components/ui/button";
import { Container } from "@/components/layout/container";

const links = [
  { href: "#how-it-works", label: "How it works" },
  { href: "#preview", label: "Product preview" },
  { href: "#privacy", label: "Privacy" },
  { href: "#faq", label: "FAQ" },
];

export function SiteHeader() {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (!open) return;
    function onKey(event: KeyboardEvent) {
      if (event.key === "Escape") setOpen(false);
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [open]);

  return (
    <header className="sticky top-0 z-50 border-b border-border/70 bg-cream/85 backdrop-blur-md">
      <Container className="flex h-[4.5rem] items-center justify-between">
        <BrandLink priority />
        <nav aria-label="Primary" className="hidden items-center gap-7 lg:flex">
          {links.map((link) => (
            <a
              key={link.href}
              href={link.href}
              className="text-small font-bold text-ink/80 transition-colors hover:text-terracotta-deep"
            >
              {link.label}
            </a>
          ))}
        </nav>
        <div className="hidden lg:block">
          <BrandButton as={Link} href="/waitlist" variant="primary">
            Join the waitlist
          </BrandButton>
        </div>
        <button
          type="button"
          className="flex size-11 items-center justify-center rounded-full text-ink lg:hidden"
          aria-expanded={open}
          aria-controls="mobile-nav"
          aria-label={open ? "Close menu" : "Open menu"}
          onClick={() => setOpen((v) => !v)}
        >
          {open ? <X className="size-6" /> : <Menu className="size-6" />}
        </button>
      </Container>
      {open && (
        <nav
          id="mobile-nav"
          aria-label="Mobile"
          className="border-t border-border/70 bg-cream lg:hidden"
        >
          <Container className="flex flex-col gap-1 py-4">
            {links.map((link) => (
              <a
                key={link.href}
                href={link.href}
                onClick={() => setOpen(false)}
                className="rounded-lg px-2 py-3 text-body font-bold text-ink"
              >
                {link.label}
              </a>
            ))}
            <BrandButton
              as={Link}
              href="/waitlist"
              variant="primary"
              className="mt-2 w-full"
              onClick={() => setOpen(false)}
            >
              Join the waitlist
            </BrandButton>
          </Container>
        </nav>
      )}
    </header>
  );
}
