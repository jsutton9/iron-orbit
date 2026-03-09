package materials

type MaterialType int
const (
	HOFuel MaterialType = 1
)
type Material struct {
	Type MaterialType
	Mass float64
}
