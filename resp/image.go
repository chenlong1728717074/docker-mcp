package resp

type Image struct {
	Repository []string
	Tag        string
	ImageID    string
	Created    int64
	Size       int64
}
