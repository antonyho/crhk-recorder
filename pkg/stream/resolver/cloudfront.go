package resolver

// CloudFront cookie names
const (
	// CloudFrontCookieNamePolicy is the cookie name for CloudFront policy
	CloudFrontCookieNamePolicy = "CloudFront-Policy"

	// CloudFrontCookieNameKeyPairID is the cookie name for CloudFront key pair ID
	CloudFrontCookieNameKeyPairID = "CloudFront-Key-Pair-Id"

	// CloudFrontCookieNameSignature is the cookie name for CloudFront signature
	CloudFrontCookieNameSignature = "CloudFront-Signature"
)

type CloudfrontCookie struct {
	Policy    string
	KeyPairID string
	Signature string
}

