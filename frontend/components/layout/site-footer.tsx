import { BrandLink } from "@/components/brand/brand-link";
import { Container } from "@/components/layout/container";

const columns = [
  {
    heading: "Product",
    links: [
      { href: "#how-it-works", label: "How it works" },
      { href: "#preview", label: "Product preview" },
      { href: "#caregivers", label: "For caregivers" },
    ],
  },
  {
    heading: "Trust",
    links: [
      { href: "#privacy", label: "Privacy" },
      { href: "#faq", label: "FAQ" },
    ],
  },
];

export function SiteFooter() {
  return (
    <footer className="border-t border-border bg-cream-soft">
      <Container className="grid gap-10 py-14 md:grid-cols-[1.3fr_1fr_1fr]">
        <div>
          <BrandLink />
          <p className="mt-4 max-w-xs text-body text-muted">
            Helping you remember what matters.
          </p>
        </div>
        {columns.map((column) => (
          <div key={column.heading}>
            <p className="text-caption uppercase text-muted">
              {column.heading}
            </p>
            <ul className="mt-3 space-y-2.5">
              {column.links.map((link) => (
                <li key={link.href}>
                  <a
                    href={link.href}
                    className="text-small font-bold text-ink/80 hover:text-terracotta-deep"
                  >
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </Container>
      <Container className="flex flex-col gap-2 border-t border-border py-6 text-small text-muted sm:flex-row sm:items-center sm:justify-between">
        <p>© {new Date().getFullYear()} Remember My Pill. Product preview.</p>
        <p>
          No medical advice or prescription processing is provided on this site.
        </p>
      </Container>
    </footer>
  );
}
