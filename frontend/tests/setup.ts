import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import { afterEach } from "vitest";

// vitest.config.ts does not set `test.globals`, so @testing-library/react's
// automatic afterEach(cleanup) registration (which checks for a global
// afterEach) never fires. Without this, renders from separate `it` blocks
// pile up in the same document.
afterEach(cleanup);

class MockIntersectionObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}
// @ts-expect-error - jsdom does not implement IntersectionObserver
globalThis.IntersectionObserver ??= MockIntersectionObserver;

window.matchMedia ??= ((query: string) => ({
  matches: false,
  media: query,
  onchange: null,
  addListener: () => {},
  removeListener: () => {},
  addEventListener: () => {},
  removeEventListener: () => {},
  dispatchEvent: () => false,
})) as unknown as typeof window.matchMedia;
