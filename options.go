package gogm

type SortDirection uint

const (
	//Ascending
	ASC SortDirection = iota

	//Desending
	DESC
)

//Not supported
type SortOptions struct {
	OrderBy       []string
	sortDirection SortDirection
}

type LoadOptions struct {
	Sort  *SortOptions
	Depth int
}

type SaveOptions struct {
	Depth int
}

func NewLoadOptions() *LoadOptions {
	lo := &LoadOptions{}
	lo.Depth = 1
	return lo
}

func NewSaveOptions() *SaveOptions {
	so := &SaveOptions{}
	so.Depth = 0
	return so
}
