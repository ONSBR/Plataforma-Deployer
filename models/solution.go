package models

//Solution is root of apps on platform
type Solution struct {
	BaseModel
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

//NewSolution creates a new solution object
func NewSolution() *Solution {
	solution := new(Solution)
	solution.Metadata = Metadata{
		Type:        "system",
		ChangeTrack: "create",
	}
	return solution
}
