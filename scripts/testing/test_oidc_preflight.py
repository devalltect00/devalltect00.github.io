"""Pure identity-preflight regressions: no real JWT, credential, network or provider write."""
import base64
import importlib.util
import json
from pathlib import Path
import unittest

SPEC = importlib.util.spec_from_file_location(
    "oidc_preflight", Path(__file__).resolve().parents[1] / "content" / "oidc-preflight.py")
PREFLIGHT = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(PREFLIGHT)
NOW = 1800000000
SHA = "a" * 40


def environment():
    """Supply synthetic public context plus one recognizable, nonsecret capability canary."""
    return {"GITHUB_ACTIONS": "true", "GITHUB_EVENT_NAME": "workflow_dispatch",
            "RUNNER_ENVIRONMENT": "github-hosted", "GITHUB_SERVER_URL": "https://github.com",
            "GITHUB_API_URL": "https://api.github.com", "GITHUB_REPOSITORY": PREFLIGHT.REPOSITORY,
            "GITHUB_REPOSITORY_ID": PREFLIGHT.REPOSITORY_ID, "GITHUB_REF": "refs/heads/main",
            "GITHUB_RUN_ATTEMPT": "1", "GITHUB_SHA": SHA, "GITHUB_WORKFLOW_SHA": SHA,
            "GITHUB_WORKFLOW_REF": PREFLIGHT.WORKFLOW, "GITHUB_RUN_ID": "77",
            "ASTRO_EXPECTED_HEAD_SHA": SHA,
            "ACTIONS_ID_TOKEN_REQUEST_URL": "https://runner.actions.githubusercontent.com/issue?api-version=2.0",
            "ACTIONS_ID_TOKEN_REQUEST_TOKEN": "synthetic-private-request-canary"}


def claims():
    """Fixture claims are observations only; their synthetic signature grants no authority."""
    return {"iss": PREFLIGHT.ISSUER, "aud": PREFLIGHT.AUDIENCE,
            "sub": "repo:" + PREFLIGHT.REPOSITORY + ":environment:github-pages",
            "repository": PREFLIGHT.REPOSITORY, "repository_id": PREFLIGHT.REPOSITORY_ID,
            "ref": "refs/heads/main", "event_name": "workflow_dispatch", "run_id": "77",
            "run_attempt": "1", "sha": SHA, "workflow_sha": SHA,
            "workflow_ref": PREFLIGHT.WORKFLOW, "iat": NOW, "nbf": NOW-5, "exp": NOW+300}


def jwt(payload=None, header=None):
    """Assemble a deliberately unsigned fixture without storing a real token literal."""
    def encode(value):
        return base64.urlsafe_b64encode(json.dumps(value).encode()).decode().rstrip("=")
    return ".".join([encode(header or {"alg": "RS256", "typ": "JWT", "kid": "public-fixture-key"}),
                     encode(payload or claims()), encode("synthetic-signature")])


class PreflightTests(unittest.TestCase):
    """Protect the no-authority/no-leak boundary while deriving the two public policy observations."""
    def setUp(self):
        self.keys = {"keys": [{"kid": "public-fixture-key", "kty": "RSA", "alg": "RS256", "use": "sig"}]}

    def test_context_rejected_before_credentials_or_network(self):
        class Guarded(dict):
            """Any credential-variable read is a test failure on an invalid workflow context."""
            def get(self, key, default=None):
                if key.startswith("ACTIONS_"):
                    raise AssertionError("credential read before context validation")
                return super().get(key, default)
        for key, value in [("GITHUB_REF", "refs/heads/foreign"), ("GITHUB_RUN_ATTEMPT", "2"),
                           ("RUNNER_ENVIRONMENT", "self-hosted"), ("GITHUB_WORKFLOW_SHA", "b"*40),
                           ("GITHUB_WORKFLOW_REF", PREFLIGHT.REPOSITORY + "/.github/workflows/cms-pages-deploy.yml@refs/heads/main"),
                           ("GITHUB_REPOSITORY_ID", "9"), ("GITHUB_RUN_ID", "077")]:
            with self.subTest(key=key):
                env=Guarded(environment());env[key]=value
                with self.assertRaises(ValueError):
                    PREFLIGHT.collect(env, fetch=lambda *args: self.fail("network before guards"), now=NOW)

    def test_issuance_origin_and_query_rejections(self):
        for url in ["http://runner.actions.githubusercontent.com/issue", "https://evil.example/issue",
                    "https://runner.actions.githubusercontent.com.evil.example/issue",
                    "https://user@runner.actions.githubusercontent.com/issue",
                    "https://runner.actions.githubusercontent.com/issue#fragment",
                    "https://runner.actions.githubusercontent.com/issue?audience=foreign",
                    "https://runner.actions.githubusercontent.com/issue?a=1&a=2"]:
            with self.subTest(url=url), self.assertRaises(ValueError):
                PREFLIGHT.issuance_request(url)

    def test_transport_separates_capability_and_report_never_contains_jwt(self):
        env=environment();token=jwt();calls=[]
        def fetch(url, bearer=None):
            """Capture only synthetic in-memory calls to prove capability isolation."""
            calls.append((url,bearer))
            return {"value":token} if len(calls)==1 else self.keys
        result=PREFLIGHT.collect(env, fetch=fetch, now=NOW)
        self.assertIn("audience=https%3A%2F%2Fgithub.com%2Fdevalltect00", calls[0][0])
        self.assertEqual(calls[0][1],env["ACTIONS_ID_TOKEN_REQUEST_TOKEN"])
        self.assertEqual(calls[1],(PREFLIGHT.JWKS,None))
        self.assertEqual(result["issuer_origin"],"https://runner.actions.githubusercontent.com")
        self.assertEqual(result["oidc_key_id"],"public-fixture-key")
        self.assertFalse(result["authorization"])
        self.assertFalse(result["signature_verified"])
        self.assertNotIn(token,json.dumps(result))
        self.assertNotIn(env["ACTIONS_ID_TOKEN_REQUEST_TOKEN"],json.dumps(result))

    def test_wrong_claims_and_time_fail_closed(self):
        for key,value in [("iss","https://foreign.example"),("aud","foreign"),("sub","repo:foreign"),
                          ("sha","b"*40),("workflow_sha","b"*40),("run_id","78"),
                          ("run_attempt","2"),("repository_id","1"),("exp",NOW),
                          ("iat",NOW+60),("nbf",NOW+60),("exp",NOW+601),("iat",True)]:
            payload=claims();payload[key]=value
            with self.subTest(key=key,value=value),self.assertRaises(ValueError):
                PREFLIGHT.observed_identity(jwt(payload),self.keys,SHA,"77",NOW)

    def test_header_or_published_key_mismatch_rejected(self):
        for header in [{"alg":"none","kid":"public-fixture-key"},
                       {"alg":"HS256","kid":"public-fixture-key"},
                       {"alg":"RS256","kid":"foreign"},
                       {"alg":"RS256","kid":"public-fixture-key","jku":"https://foreign.example"}]:
            with self.subTest(header=header),self.assertRaises(ValueError):
                PREFLIGHT.observed_identity(jwt(header=header),self.keys,SHA,"77",NOW)
        for keys in [{"keys":[]},{"keys":self.keys["keys"]*2},
                     {"keys":[{"kid":"public-fixture-key","kty":"EC"}]}]:
            with self.subTest(keys=keys),self.assertRaises(ValueError):
                PREFLIGHT.observed_identity(jwt(),keys,SHA,"77",NOW)

    def test_ambiguous_oversized_or_nonfinite_json_rejected(self):
        for raw in [b'{"kid":"one","kid":"two"}',b'{"x":NaN}',b' '*65537,b'\xff']:
            with self.subTest(raw_type=type(raw)),self.assertRaises((ValueError,UnicodeError)):
                PREFLIGHT.parse_json(raw)


if __name__ == "__main__":
    unittest.main()
