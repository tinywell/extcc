package server

type Opt func(*config)

func WithSigner(mspid, key, cert string) Opt {
	return func(c *config) {
		c.signMSP = mspid
		c.signKey = key
		c.signCert = cert
	}
}

func WithTransient(transient map[string][]byte) Opt {
	return func(c *config) {
		c.transient = transient
	}
}
