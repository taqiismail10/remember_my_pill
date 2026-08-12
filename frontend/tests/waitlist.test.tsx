import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import WaitlistPage from "@/app/waitlist/page";

vi.mock("next/image", () => ({
  default: (props: React.ImgHTMLAttributes<HTMLImageElement>) => (
    // eslint-disable-next-line @next/next/no-img-element
    <img alt="" {...props} />
  ),
}));

function jsonResponse(status: number, body: unknown) {
  return Promise.resolve(
    new Response(JSON.stringify(body), {
      status,
      headers: { "Content-Type": "application/json" },
    }),
  );
}

describe("/waitlist page", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("renders a single accessible email field and no name field", () => {
    render(<WaitlistPage />);
    const email = screen.getByLabelText("Email");
    expect(email).toBeInTheDocument();
    expect(email).toHaveAttribute("type", "email");
    expect(email).toHaveAttribute("autoComplete", "email");
    expect(screen.queryByLabelText("Name")).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/name/i)).not.toBeInTheDocument();
  });

  it("has a heading and a link back to the homepage", () => {
    render(<WaitlistPage />);
    expect(
      screen.getByRole("heading", {
        level: 1,
        name: /be first to know when remember my pill is ready/i,
      }),
    ).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /back to home/i })).toHaveAttribute(
      "href",
      "/",
    );
  });

  it("rejects an invalid email without calling the API", () => {
    render(<WaitlistPage />);
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "not-an-email" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Join the waitlist" }));
    expect(
      screen.getByText(/enter a valid email address/i),
    ).toBeInTheDocument();
    expect(fetch).not.toHaveBeenCalled();
  });

  it("submits only an email, never a name field", async () => {
    vi.mocked(fetch).mockReturnValue(jsonResponse(202, { status: "accepted" }));
    render(<WaitlistPage />);
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: " Ada@Example.com " },
    });
    fireEvent.click(screen.getByRole("button", { name: "Join the waitlist" }));

    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));
    const [, init] = vi.mocked(fetch).mock.calls[0];
    const body = JSON.parse(init?.body as string);
    expect(body).toEqual({ email: "ada@example.com" });
  });

  it("shows a loading state and disables the button while submitting", async () => {
    let resolveFetch: (value: Response) => void = () => {};
    vi.mocked(fetch).mockReturnValue(
      new Promise((resolve) => {
        resolveFetch = resolve;
      }),
    );
    render(<WaitlistPage />);
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "ada@example.com" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Join the waitlist" }));

    const button = await screen.findByRole("button", { name: /joining/i });
    expect(button).toBeDisabled();

    resolveFetch(await jsonResponse(202, { status: "accepted" }));
  });

  it("shows a branded success state without promising rank or referrals", async () => {
    vi.mocked(fetch).mockReturnValue(jsonResponse(202, { status: "accepted" }));
    render(<WaitlistPage />);
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "ada@example.com" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Join the waitlist" }));

    expect(
      await screen.findByRole("heading", {
        level: 1,
        name: /you're on the list/i,
      }),
    ).toBeInTheDocument();
    expect(screen.queryByText(/rank/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/referral/i)).not.toBeInTheDocument();
  });

  it("shows the same neutral success message for an accepted request", async () => {
    vi.mocked(fetch).mockReturnValue(jsonResponse(202, { status: "accepted" }));
    render(<WaitlistPage />);
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "ada@example.com" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Join the waitlist" }));

    expect(
      await screen.findByText(/if this email is eligible/i),
    ).toBeInTheDocument();
  });

  it("shows a calm retry message on a server/network error", async () => {
    vi.mocked(fetch).mockRejectedValue(new Error("network down"));
    render(<WaitlistPage />);
    fireEvent.change(screen.getByLabelText("Email"), {
      target: { value: "ada@example.com" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Join the waitlist" }));

    expect(
      await screen.findByText(/couldn't reach the waitlist service/i),
    ).toBeInTheDocument();
  });

  it("submits via the keyboard (Enter) without a mouse click", async () => {
    vi.mocked(fetch).mockReturnValue(jsonResponse(202, { status: "accepted" }));
    render(<WaitlistPage />);
    const email = screen.getByLabelText("Email");
    fireEvent.change(email, { target: { value: "ada@example.com" } });
    fireEvent.submit(email.closest("form") as HTMLFormElement);

    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));
  });
});
