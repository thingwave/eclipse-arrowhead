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
