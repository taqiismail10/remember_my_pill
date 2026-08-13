import { render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { StatusVerify } from "@/components/status-access/status-verify";
import { StatusPage } from "@/components/status-access/status-page";

describe("status-access browser foundation", () => {
  beforeEach(() => vi.stubGlobal("fetch", vi.fn()));
  afterEach(() => {
    vi.unstubAllGlobals();
    window.history.replaceState(null, "", "/");
  });

  it("cleans a verification fragment without sending it while disabled", async () => {
    window.history.replaceState(
      null,
      "",
      "/waitlist/verify#v=short-lived-token",
    );
    render(<StatusVerify />);
    await screen.findByText(/not available yet/i);
    expect(window.location.hash).toBe("");
    expect(fetch).not.toHaveBeenCalled();
  });

  it("posts the fragment token in JSON only when explicitly enabled", async () => {
    vi.mocked(fetch).mockResolvedValue(
      new Response(`{"error":"unavailable"}`, { status: 401 }),
    );
    window.history.replaceState(
      null,
      "",
      "/waitlist/verify#v=short-lived-token",
    );
    render(<StatusVerify enabled />);
    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));
    const [url, init] = vi.mocked(fetch).mock.calls[0];
    expect(url).toContain("/api/waitlist/status-access/exchange");
    expect(String(url)).not.toContain("verificationToken");
    expect(JSON.parse(init?.body as string)).toEqual({
      verificationToken: "short-lived-token",
    });
    expect(window.location.hash).toBe("");
  });

  it("keeps the status page unavailable by default without credential storage", () => {
    render(<StatusPage />);
    expect(screen.getByText(/not available yet/i)).toBeInTheDocument();
    expect(fetch).not.toHaveBeenCalled();
  });

  it("renders only approved status fields when enabled", async () => {
    vi.mocked(fetch).mockResolvedValue(
      new Response(`{"rank":12,"referralCount":0,"referralCode":"RMP12"}`, {
        status: 200,
      }),
    );
    render(<StatusPage enabled />);
    expect(await screen.findByText("12")).toBeInTheDocument();
    expect(screen.getByText("RMP12")).toBeInTheDocument();
    expect(screen.queryByText(/email/i)).not.toBeInTheDocument();
  });
});
