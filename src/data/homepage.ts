/** Astro frontmatter reads one pinned build snapshot; browsers never contact the CMS or import server adapters. */
import { projects } from "./projects.ts";
import { site } from "./site.ts";
import { readPrepared, selectRelease } from "../server/cms/cache.ts";
import type { GatewaySnapshot } from "../server/cms/snapshot.ts";
import type { ContactLink, Project } from "../types/site.ts";
export interface HomepageContent {
  source: "authored" | "cms-build-preview";
  projects: Project[];
  profile: GatewaySnapshot["profile"];
  contacts: ContactLink[];
  snapshot?: GatewaySnapshot;
}
let memo: HomepageContent | undefined;
/** Authored identity is loaded only without explicit CMS opt-in; a missing CMS profile intentionally has no authored contacts. */
export function homepageContent(): HomepageContent {
  if (memo) return memo;
  const id = selectRelease(process.env.HOMEPAGE_CMS_RELEASE);
  if (id) {
    const snapshot = readPrepared(process.cwd(), id);
    memo = {
      source: "cms-build-preview",
      projects: snapshot.projects,
      profile: snapshot.profile,
      contacts: snapshot.contacts,
      snapshot,
    };
  } else
    memo = {
      source: "authored",
      projects,
      profile: null,
      contacts: [
        { kind: "email", href: "mailto:" + site.email, icon: "mail" },
        { kind: "github", href: site.githubUrl, icon: "github" },
        { kind: "gitlab", href: site.gitlabUrl, icon: "gitlab" },
        { kind: "linkedin", href: site.linkedInUrl, icon: "linkedin" },
      ],
    };
  return memo;
}
/** Public metadata records prepared build identity only; it cannot acknowledge deployment or complete a CMS job. */
export function buildMarker(state: HomepageContent) {
  return state.snapshot
    ? {
        schema_version: 1,
        source: state.source,
        scope: "astro-project-profile-preview-v1",
        release_id: state.snapshot.releaseID,
        sequence: state.snapshot.sequence,
        snapshot_digest: state.snapshot.snapshotDigest,
        projection_digest: state.snapshot.projectionDigest,
        consumed_kinds: ["project", "profile"],
      }
    : { schema_version: 1, source: "authored", consumed_kinds: [] };
}
