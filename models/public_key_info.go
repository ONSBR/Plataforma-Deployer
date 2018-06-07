package models

//PublicKeyInfo represents the current state of some public key installed
type PublicKeyInfo struct {
	Solution string
	Name     string
	Status   string
}

//Ok set installed status to current public key info
func (key *PublicKeyInfo) Ok() {
	key.Status = "installed"
}
