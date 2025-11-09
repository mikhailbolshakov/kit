package vault

const (
	CSRTypeInternal CSRType = "internal" // the private key will not be returned and cannot be retrieved later
	CSRTypeExported CSRType = "exported" // the private key will be returned in the response
	CSRTypeExisting CSRType = "existing" // we expect the key_ref parameter to use existing key material to create the CSR; kms is also supported

	SignFormatPem       SignFormat = "pem"
	SignFormatDer       SignFormat = "der"        // f der, the output is base64 encoded
	SignFormatPemBundle SignFormat = "pem_bundle" //  If pem_bundle, the certificate field will contain the certificate and, if the issuing CA is not a Vault-derived self-signed root, it will be concatenated with the certificate.
)

type CSRType string
type SignFormat string

type CSR struct {
	CSR            string `json:"csr" vault:"required"` // PEM-encoded Certificate Signing Request (CSR) to be signed by a Certificate Authority.
	KeyId          string `json:"key_id"`               // Unique identifier for the key associated with the CSR.
	PrivateKey     string `json:"private_key"`          // PEM-encoded private key corresponding to the CSR.
	PrivateKeyType string `json:"private_key_type"`     // Type of the private key (e.g., RSA, ECDSA).
}

type Role struct {
	AllowAnyName     bool     `json:"allow_any_name"`     // Specifies if clients can request any CN. Useful in some circumstances, but make sure you understand whether it is appropriate for your installation before enabling it. Note that both enforce_hostnames and allow_wildcard_certificates are still checked, which may introduce limitations on issuance with this option.
	AllowIpSans      bool     `json:"allow_ip_sans"`      // Specifies if clients can request IP Subject Alternative Names. No authorization checking is performed except to verify that the given values are valid IP addresses.
	AllowLocalhost   bool     `json:"allow_localhost"`    // Specifies if clients can request certificates for localhost as one of the requested common names. This is useful for testing and to allow clients on a single host to talk securely.
	AllowSubdomains  bool     `json:"allow_subdomains"`   // Specifies if clients can request certificates with CNs that are subdomains of the CNs allowed by the other role options. This includes wildcard subdomains. For example, an allowed_domains value of example.com with this option set to true will allow foo.example.com and bar.example.com as well as *.example.com. To restrict issuance of wildcards by this option, see allow_wildcard_certificates below. This option is redundant when using the allow_any_name option.
	AllowedDomains   []string `json:"allowed_domains"`    // Specifies the domains this role is allowed to issue certificates for. This is used with the allow_bare_domains, allow_subdomains, and allow_glob_domains options to determine the type of matching between these domains and the values of common name, DNS-typed SAN entries, and Email-typed SAN entries. When allow_any_name is used, this attribute has no effect.
	AllowedUriSans   []string `json:"allowed_uri_sans"`   // Defines allowed URI Subject Alternative Names. No authorization checking is performed except to verify that the given values are valid URIs. This can be a comma-delimited list or a JSON string slice. Values can contain glob patterns (e.g. spiffe://hostname/*).
	AllowedOtherSans []string `json:"allowed_other_sans"` // Defines allowed custom OID/UTF8-string SANs. This can be a comma-delimited list or a JSON string slice, where each element has the same format as OpenSSL: <oid>;<type>:<value>, but the only valid type is UTF8 or UTF-8. The value part of an element may be a * to allow any value with that OID. Alternatively, specifying a single * will allow any other_sans input.
	ClientFlag       bool     `json:"client_flag"`        // Specifies if certificates are flagged for client authentication use. See RFC 5280 Section 4.2.1.12 for information about the Extended Key Usage field.
	CodeSigningFlag  bool     `json:"code_signing_flag"`  // Specifies if certificates are flagged for code signing use. See RFC 5280 Section 4.2.1.12 for information about the Extended Key Usage field.
	KeyBits          int      `json:"key_bits"`           // Specifies the number of bits to use for the generated keys. Allowed values are 0 (universal default); with key_type=rsa, allowed values are: 2048 (default), 3072, 4096 or 8192; with key_type=ec, allowed values are: 224, 256 (default), 384, or 521; ignored with key_type=ed25519 or in signing operations when key_type=any.
	KeyType          string   `json:"key_type"`           // Specifies the type of key to generate for generated private keys and the type of key expected for submitted CSRs. Currently, rsa, ec, and ed25519 are supported, or when signing existing CSRs, any can be specified to allow keys of either type and with any bit size (subject to >=2048 bits for RSA keys or >= 224 for EC keys). When any is used, this role cannot generate certificates and can only be used to sign CSRs.
	Ttl              uint64   `json:"ttl"`                // Specifies the Time To Live value to be used for the validity period of the requested certificate, provided as a string duration with time suffix. Hour is the largest suffix. The value specified is strictly used for future validity. If not set, uses the system default value or the value of max_ttl, whichever is shorter. See not_after as an alternative for setting an absolute end date (rather than a relative one).
	MaxTtl           uint64   `json:"max_ttl"`            // Specifies the maximum Time To Live provided as a string duration with time suffix. Hour is the largest suffix. If not set, defaults to the system maximum lease TTL.
	ServerFlag       bool     `json:"server_flag"`        // Specifies if certificates are flagged for server authentication use. See RFC 5280 Section 4.2.1.12 for information about the Extended Key Usage field.
}

type SetSignedCertificateResponse struct {
	ImportedIssuers []string          `json:"imported_issuers"` // The response will indicate what issuers and keys were created as part of this request (in the imported_issuers and imported_keys
	ImportedKeys    []string          `json:"imported_keys"`    // The response will indicate what issuers and keys were created as part of this request (in the imported_issuers and imported_keys
	Mapping         map[string]string `json:"mapping"`          // Along with a mapping field, indicating which keys belong to which issuers (including from already imported entries present in the same bundle
	ExistingIssuers []string          `json:"existing_issuers"` // The response also contains an existing_issuers and existing_keys fields, which specifies the issuer and key IDs of any entries in the bundle that already existed within this mount.
	ExistingKeys    []string          `json:"existing_keys"`    // The response also contains an existing_issuers and existing_keys fields, which specifies the issuer and key IDs of any entries in the bundle that already existed within this mount.
}

type Empty struct {
}

type CertificateResponse struct {
	Certificate    string `json:"certificate"` // Specifies the PEM-encoded Certificate.
	RevocationTime int    `json:"revocation_time"`
}

type Certificate struct {
	Certificate  string     `json:"certificate"`   // PEM-encoded certificate issued to the client.
	Expiration   uint64     `json:"expiration"`    // Expiration time of the certificate as a Unix timestamp.
	IssuerId     string     `json:"issuer_id"`     // Unique identifier of the certificate's issuer.
	IssuerName   string     `json:"issuer_name"`   // Name of the certificate's issuer.
	IssuingCa    string     `json:"issuing_ca"`    // PEM-encoded issuing Certificate Authority (CA) certificate.
	KeyId        string     `json:"key_id"`        // Unique identifier of the key associated with the certificate.
	KeyName      string     `json:"key_name"`      // Descriptive name for the key associated with the certificate.
	SerialNumber string     `json:"serial_number"` // Unique serial number assigned to the certificate.
	Format       SignFormat `json:"format"`        // Format in which the certificate is signed (e.g., PEM, DER).
	NotAfter     string     `json:"not_after"`     // Expiration date of the certificate in a human-readable format.
}

type SignRequest struct {
	CSR      string     `json:"csr"`                 // Specifies the PEM-encoded CSR.
	Name     string     `json:"name,omitempty"`      // Specifies a role. If set, the following parameters from the role will have effect: ttl, max_ttl, issuer, generate_lease, no_store, no_store_metadata and not_before_duration.
	Format   SignFormat `json:"format,omitempty"`    // Specifies the format for returned data
	TTL      string     `json:"ttl,omitempty"`       // Specifies the requested Time To Live. Cannot be greater than the engine's max_ttl value. If not provided, the engine's ttl value will be used, which defaults to system values if not explicitly set. See not_after as an alternative for setting an absolute end date (rather than a relative one).
	NotAfter string     `json:"not_after,omitempty"` // Set the Not After field of the certificate with specified date value. The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ. Supports the Y10K end date for IEEE 802.1AR-2018 standard devices, 9999-12-31T23:59:59Z.
}

type IssueRequest struct {
	CommonName string     `json:"common_name,omitempty"` // Specifies the requested CN for the certificate
	Format     SignFormat `json:"format,omitempty"`      // Specifies the format for returned data
	TTL        string     `json:"ttl,omitempty"`         // Specifies the requested Time To Live. Cannot be greater than the engine's max_ttl value. If not provided, the engine's ttl value will be used, which defaults to system values if not explicitly set. See not_after as an alternative for setting an absolute end date (rather than a relative one).
	NotAfter   string     `json:"not_after,omitempty"`   // Set the Not After field of the certificate with specified date value. The value format should be given in UTC format YYYY-MM-ddTHH:MM:SSZ. Supports the Y10K end date for IEEE 802.1AR-2018 standard devices, 9999-12-31T23:59:59Z.
}

type IssueResponse struct {
	Certificate    string   `json:"certificate"`      // PEM-encoded certificate issued to the client.
	Expiration     uint64   `json:"expiration"`       // Expiration time of the certificate as a Unix timestamp.
	IssuingCa      string   `json:"issuing_ca"`       // PEM-encoded issuing Certificate Authority (CA) certificate.
	CaChain        []string `json:"ca_chain"`         // Chain of CA certificates leading up to the root certificate.
	SerialNumber   string   `json:"serial_number"`    // Unique serial number assigned to the certificate.
	PrivateKey     string   `json:"private_key"`      // PEM-encoded private key associated with the certificate.
	PrivateKeyType string   `json:"private_key_type"` // Type of the private key (e.g., RSA, ECDSA).
}

type RoleRequest struct {
	Name            string   `json:"-"`
	AllowedDomains  []string `json:"allowed_domains,omitempty"`  // Specifies the domains this role is allowed to issue certificates for. This is used with the allow_bare_domains, allow_subdomains, and allow_glob_domains options to determine the type of matching between these domains and the values of common name, DNS-typed SAN entries, and Email-typed SAN entries. When allow_any_name is used, this attribute has no effect.
	AllowSubdomains *bool    `json:"allow_subdomains,omitempty"` // Specifies if clients can request certificates with CNs that are subdomains of the CNs allowed by the other role options. This includes wildcard subdomains. For example, an allowed_domains value of example.com with this option set to true will allow foo.example.com and bar.example.com as well as *.example.com. To restrict issuance of wildcards by this option, see allow_wildcard_certificates below. This option is redundant when using the allow_any_name option.
	MaxTtl          string   `json:"max_ttl,omitempty"`          // Specifies the maximum Time To Live provided as a string duration with time suffix. Hour is the largest suffix. If not set, defaults to the system maximum lease TTL.
}

type CSRRequest struct {
	Type CSRType `json:"-"`
	CertificateRequest
}

type CertificateRequest struct {
	CommonName         string     `json:"common_name"`            // Specifies the requested CN for the certificate. If more than one common_name is desired, specify the alternative names in the alt_names list.
	AlternativeNames   string     `json:"alt_names,omitempty"`    // Specifies the requested Subject Alternative Names, in a comma-delimited list. These can be host names or email addresses; they will be parsed into their respective fields.
	Country            []string   `json:"country,omitempty"`      // Specifies the C (Country) values in the subject field of the resulting CSR.
	Province           []string   `json:"province,omitempty"`     // Specifies the ST (Province) values in the subject field of the resulting CSR.
	Locality           []string   `json:"locality,omitempty"`     // Specifies the L (Locality) values in the subject field of the resulting CSR.
	Organization       []string   `json:"organization,omitempty"` // Specifies the O (Organization) values in the subject field of the resulting CSR.
	OrganizationalUnit []string   `json:"ou,omitempty"`           // Specifies the OU (OrganizationalUnit) values in the subject field of the resulting CSR.
	Format             SignFormat `json:"format,omitempty"`       // Specifies the format for returned data
}

type SetSignedCertificateRequest struct {
	Certificate string `json:"certificate"` // Specifies the PEM-encoded Certificate.
}
