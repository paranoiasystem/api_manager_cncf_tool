package envoy.authz
import future.keywords.if

import input.attributes.request.http as http_request

# https://www.envoyproxy.io/docs/envoy/v1.26.1/start/sandboxes/ext_authz.html

default allow := false

allow if is_token_valid

is_token_present if {
    http_request.headers["x-auth-token"]
}

is_token_empty if {
	http_request.headers["x-auth-token"] == ""
}

is_token_valid if {
    is_token_present
    not is_token_empty
    http_request.headers["x-auth-token"] == "peppe"
}

status_code := 200 {
  allow
} else = 401 {
  not is_token_valid
} else = 403 {
  not is_token_present
}