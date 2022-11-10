package main

type TrustedKeyCheckRequestDTO struct {
	Version   int    `json:"version"`
	PublicKey string `json:"publicKey"`
}

type TrustedKeyCheckResponseDTO struct {
	Version       int    `json:"version"`
	ProducedAt    string `json:"producedAt"`
	EndOfValidity string `json:"endOfValidity"`
	CommonName    string `json:"commonName"`
	SerialNumber  string `json:"serialNumber"`
	Status        string `json:"Status"`
}

type CertificateSigningRequest struct {
	EncodedCSR  string `json:"encodedCSR"`
	ValidAfter  string `json:"validAfter"`
	ValidBefore string `json:"validBefore"`
}

type CertificateSigningResponse struct {
	Id               int       `json:"id"`
	CertificateChain [3]string `json:"certificateChain"`
}

type TrustedKeyCheckRequest struct {
	PublicKey string `json:"PublicKey"`
}

type TrustedKeyCheckResponse struct {
	Id          int    `json:"id"`
	CreatedAt   string `json:"createdAt"`
	Description string `json:"description"`
}
