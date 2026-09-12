# Auto-injected credentials — redacted for security; set TESTSPRITE_AUTH_CREDENTIAL to run this locally
__AUTH_CREDENTIAL__ = __import__("os").environ.get("TESTSPRITE_AUTH_CREDENTIAL", "")
__AUTH_TYPE__ = "public"
__EXTRA_HEADERS__ = {}
__AUTH_HEADERS__ = {}
import requests

BASE_URL = "https://olx-api-xpl4.onrender.com"
VALID = {
    "title": "Valid title",
    "description": "A valid listing description.",
    "price": 1,
    "city": "Pune",
}


def expect_validation_error(name, changes, expected_field):
    payload = dict(VALID)
    payload.update(changes)
    response = requests.post(f"{BASE_URL}/listings", json=payload, timeout=20)
    assert response.status_code == 422, f"{name}: expected 422, got {response.status_code}: {response.text}"
    body = response.json()
    error = body.get("error", {})
    assert error.get("code") == "validation_failed", f"{name}: {body}"
    assert error.get("field") == expected_field, f"{name}: expected field {expected_field!r}, got {body}"


def test_all_create_listing_field_validations():
    cases = [
        ("empty title", {"title": ""}, "title"),
        ("whitespace title", {"title": "   "}, "title"),
        ("title over 100 characters", {"title": "t" * 101}, "title"),
        ("empty description", {"description": ""}, "description"),
        ("whitespace description", {"description": "   "}, "description"),
        ("description over 256 characters", {"description": "d" * 257}, "description"),
        ("empty city", {"city": ""}, "city"),
        ("whitespace city", {"city": "   "}, "city"),
        ("city over 100 characters", {"city": "c" * 101}, "city"),
        ("zero price", {"price": 0}, "price"),
        ("negative price", {"price": -1}, "price"),
    ]
    failures = []
    for name, changes, expected_field in cases:
        try:
            expect_validation_error(name, changes, expected_field)
        except AssertionError as error:
            failures.append(str(error))
    assert not failures, "\\n".join(failures)


test_all_create_listing_field_validations()