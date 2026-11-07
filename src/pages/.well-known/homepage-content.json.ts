/** Static public build identity has no secret, activation or destination-completion authority. */
import { homepageContent, buildMarker } from "../../data/homepage.ts";
export const prerender = true;
export function GET() {
  return new Response(JSON.stringify(buildMarker(homepageContent())), {
    headers: { "Content-Type": "application/json" },
  });
}
