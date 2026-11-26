"""Observe public workflow identity metadata without granting CMS or Pages authority.

The JWT and request capability stay in memory. This bootstrap selects metadata for
operator review; the existing CMS verifier must still verify the real signature
and exact deployment-workflow/current-Owner/attempt fences on every write.
"""
import base64
import json
import os
from pathlib import Path
import re
import time
import urllib.parse
import urllib.request

REPOSITORY = "devalltect00/devalltect00.github.io"
REPOSITORY_ID = "1389151718"
WORKFLOW = REPOSITORY + "/.github/workflows/cms-oidc-preflight.yml@refs/heads/main"
AUDIENCE = "https://github.com/devalltect00"
ISSUER = "https://token.actions.githubusercontent.com"
JWKS = ISSUER + "/.well-known/jwks"


def unique_object(pairs):
    """Reject duplicate JSON keys rather than accept ambiguous issuer metadata."""
    result = {}
    for key, value in pairs:
        if key in result:
            raise ValueError("ambiguous JSON")
        result[key] = value
    return result


def reject_constant(value):
    """Non-finite JSON numbers are not valid identity input."""
    raise ValueError("non-finite JSON")


def parse_json(raw):
    """Decode finite bounded UTF-8 JSON with one unambiguous meaning."""
    if not isinstance(raw, bytes) or not 0 < len(raw) <= 65536:
        raise ValueError("JSON size")
    return json.loads(raw.decode("utf-8"), object_pairs_hook=unique_object,
                      parse_constant=reject_constant)


class NoRedirect(urllib.request.HTTPRedirectHandler):
    """An ephemeral request capability may never follow a redirect."""
    def redirect_request(self, *args, **kwargs):
        return None


def network_json(url, bearer=None):
    """Use TLS, no proxy/cookies/redirects, finite bytes, one request and no raw-error logging."""
    headers = {"Accept": "application/json", "User-Agent": "DevalltectOIDCPreflight/1"}
    if bearer is not None:
        headers["Authorization"] = "Bearer " + bearer
    maximum = 32768 if bearer is not None else 65536
    request = urllib.request.Request(url, headers=headers)
    opener = urllib.request.build_opener(urllib.request.ProxyHandler({}), NoRedirect)
    with opener.open(request, timeout=8) as response:
        if response.status != 200 or response.headers.get("Content-Type", "").split(";")[0] != "application/json":
            raise ValueError("issuer response")
        if response.headers.get("Content-Encoding") not in (None, ""):
            raise ValueError("encoded response")
        raw = response.read(maximum + 1)
        if len(raw) > maximum:
            raise ValueError("response size")
        return parse_json(raw)


def validate_context(env):
    """Validate code/repository/first-main-run before reading any runtime credential variables."""
    expected = env.get("ASTRO_EXPECTED_HEAD_SHA", "")
    if not re.fullmatch(r"[a-f0-9]{40}", expected):
        raise ValueError("code pin")
    required = {
        "GITHUB_ACTIONS": "true", "GITHUB_EVENT_NAME": "workflow_dispatch",
        "RUNNER_ENVIRONMENT": "github-hosted", "GITHUB_SERVER_URL": "https://github.com",
        "GITHUB_API_URL": "https://api.github.com", "GITHUB_REPOSITORY": REPOSITORY,
        "GITHUB_REPOSITORY_ID": REPOSITORY_ID, "GITHUB_REF": "refs/heads/main",
        "GITHUB_RUN_ATTEMPT": "1", "GITHUB_SHA": expected,
        "GITHUB_WORKFLOW_SHA": expected, "GITHUB_WORKFLOW_REF": WORKFLOW,
    }
    if any(env.get(key) != value for key, value in required.items()):
        raise ValueError("workflow context")
    run = env.get("GITHUB_RUN_ID", "")
    if not re.fullmatch(r"[1-9][0-9]{0,18}", run) or int(run) > 9223372036854775807:
        raise ValueError("run identity")
    return expected, run


def issuance_request(raw):
    """Allow the GitHub-hosted issuance endpoint only; emit a separate query-free origin."""
    if not isinstance(raw, str) or not 1 <= len(raw) <= 4096:
        raise ValueError("issuance URL")
    parsed = urllib.parse.urlsplit(raw)
    if (parsed.scheme != "https" or parsed.username is not None or parsed.password is not None
            or parsed.fragment or parsed.port not in (None, 443) or not parsed.hostname
            or not re.fullmatch(r"[a-z0-9.-]+\.actions\.githubusercontent\.com", parsed.hostname)
            or parsed.netloc != parsed.hostname + (":443" if parsed.port == 443 else "")
            or not parsed.path.startswith("/") or parsed.path == "/"):
        raise ValueError("issuance origin")
    query = urllib.parse.parse_qsl(parsed.query, keep_blank_values=True)
    if len({key for key, _ in query}) != len(query) or any(key == "audience" for key, _ in query):
        raise ValueError("issuance query")
    endpoint = urllib.parse.urlunsplit(parsed._replace(query=urllib.parse.urlencode(query + [("audience", AUDIENCE)])))
    return endpoint, "https://" + parsed.netloc


def jwt_object(segment):
    """Parse only bounded canonical base64url segments; this is not signature verification."""
    if not isinstance(segment, str) or not re.fullmatch(r"[A-Za-z0-9_-]{1,32768}", segment):
        raise ValueError("JWT segment")
    raw = base64.urlsafe_b64decode(segment + "=" * (-len(segment) % 4))
    if base64.urlsafe_b64encode(raw).decode().rstrip("=") != segment:
        raise ValueError("noncanonical JWT segment")
    value = parse_json(raw)
    if not isinstance(value, dict):
        raise ValueError("JWT object")
    return value


def observed_identity(token, keys, expected, run, now):
    """Return only allowlisted public metadata, never token bytes, signature, arbitrary claims or authority."""
    if not isinstance(token, str) or not 1 <= len(token) <= 16384:
        raise ValueError("JWT size")
    segments = token.split(".")
    if len(segments) != 3 or not re.fullmatch(r"[A-Za-z0-9_-]{1,8192}", segments[2]):
        raise ValueError("JWT shape")
    header, claims = jwt_object(segments[0]), jwt_object(segments[1])
    kid = header.get("kid", "")
    if (set(header) - {"typ", "alg", "kid", "x5t"} or header.get("alg") != "RS256"
            or header.get("typ", "JWT") != "JWT" or not isinstance(kid, str) or not re.fullmatch(r"[A-Za-z0-9_-]{1,128}", kid)):
        raise ValueError("JWT header")
    fixed = {"iss": ISSUER, "aud": AUDIENCE, "sub": "repo:" + REPOSITORY + ":environment:github-pages",
             "repository": REPOSITORY, "repository_id": REPOSITORY_ID, "ref": "refs/heads/main",
             "event_name": "workflow_dispatch", "run_id": run, "run_attempt": "1",
             "sha": expected, "workflow_sha": expected, "workflow_ref": WORKFLOW}
    if any(claims.get(key) != value for key, value in fixed.items()):
        raise ValueError("JWT context")
    if any(type(claims.get(key)) is not int for key in ["iat", "nbf", "exp"]):
        raise ValueError("JWT time type")
    if (claims["iat"] > now + 30 or claims["nbf"] > now + 30 or claims["exp"] <= now
            or not 0 < claims["exp"] - claims["iat"] <= 600 or now - claims["iat"] > 600):
        raise ValueError("JWT time")
    if not isinstance(keys, dict) or not isinstance(keys.get("keys"), list) or not 1 <= len(keys["keys"]) <= 32:
        raise ValueError("JWKS")
    matches = [key for key in keys["keys"] if isinstance(key, dict) and key.get("kid") == kid
               and key.get("kty") == "RSA" and key.get("alg", "RS256") == "RS256" and key.get("use", "sig") == "sig"]
    if len(matches) != 1:
        raise ValueError("published key")
    return {"schema_version": 1, "scope": "astro-oidc-preflight-metadata-v1",
            "repository": REPOSITORY, "repository_id": REPOSITORY_ID, "head_sha": expected,
            "run_id": run, "workflow_ref": WORKFLOW, "audience": AUDIENCE, "oidc_key_id": kid,
            "published_key_match": True, "signature_verified": False, "authorization": False}


def collect(env, fetch=network_json, now=None):
    """Request one ephemeral identity after guards, then read public keys without forwarding the capability."""
    expected, run = validate_context(env)
    endpoint, origin = issuance_request(env.get("ACTIONS_ID_TOKEN_REQUEST_URL", ""))
    bearer = env.get("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "")
    if not isinstance(bearer, str) or not re.fullmatch(r"[A-Za-z0-9._~-]{1,16384}", bearer):
        raise ValueError("request capability")
    issued = fetch(endpoint, bearer)
    keys = fetch(JWKS)
    report = observed_identity(issued.get("value"), keys, expected, run, time.time() if now is None else now)
    report["issuer_origin"] = origin
    return report


def main():
    """Write public metadata only to one exclusive runner-temp file; failures disclose no upstream strings."""
    try:
        report = collect(os.environ)
        parent = Path(os.environ["RUNNER_TEMP"])
        if not parent.is_absolute() or not parent.is_dir():
            raise ValueError("runner directory")
        file = parent / "astro-oidc-preflight.json"
        descriptor = os.open(file, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
        with os.fdopen(descriptor, "w", encoding="utf-8", newline="\n") as output:
            json.dump(report, output, indent=2)
            output.write("\n")
        print("OIDC_PREFLIGHT_METADATA_SAVED; no Pages deployment or CMS admission")
    except Exception:
        raise SystemExit("OIDC_PREFLIGHT_UNVERIFIED") from None


if __name__ == "__main__":
    main()
