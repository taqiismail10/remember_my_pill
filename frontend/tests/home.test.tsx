import { render, screen, within } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import Home from "@/app/page";

vi.mock("next/image", () => ({
  default: (props: React.ImgHTMLAttributes<HTMLImageElement>) => (
    // eslint-disable-next-line @next/next/no-img-element
    <img alt="" {...props} />
  ),
}));

describe("marketing homepage", () => {
  it("renders the hero, primary navigation, and privacy notice", () => {
    render(<Home />);
    expect(
      screen.getByRole("heading", {
        level: 1,
        name: /turn confusing prescriptions into a routine you can trust/i,
      }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("navigation", { name: "Primary" }),
    ).toBeInTheDocument();
    expect(
      screen.getByText(/does not collect prescriptions, medication details/i),
    ).toBeInTheDocument();
  });

  it("routes the waitlist section's CTA to the dedicated /waitlist page instead of an inline form", () => {
    render(<Home />);
    const waitlistHeading = screen.getByRole("heading", {
      level: 2,
      name: "Join the waitlist",
    });
    const section = waitlistHeading.closest("section");
    expect(section).not.toBeNull();
    const scoped = within(section as HTMLElement);
    expect(
      scoped.getByRole("link", { name: "Join the waitlist" }),
    ).toHaveAttribute("href", "/waitlist");
    // Email-only signup happens on /waitlist, not inline on the landing page.
    expect(scoped.queryByLabelText("Name")).not.toBeInTheDocument();
    expect(scoped.queryByLabelText("Email")).not.toBeInTheDocument();
  });

  it("routes every high-intent CTA (header, hero, final CTA) to /waitlist", () => {
    render(<Home />);
    const ctas = screen.getAllByRole("link", { name: "Join the waitlist" });
    expect(ctas.length).toBeGreaterThanOrEqual(3);
    for (const cta of ctas) {
      expect(cta).toHaveAttribute("href", "/waitlist");
    }
  });

  it("addresses referral links and waitlist rank honestly in the FAQ", () => {
    render(<Home />);
    expect(
      screen.getByRole("button", {
        name: /will i get a referral link or waitlist rank/i,
      }),
    ).toBeInTheDocument();
  });
});
