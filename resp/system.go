package resp

type System struct {
	APIVersion       string
	OSType           string
	Experimental     bool
	BuilderVersion   string
	NodeState        string
	ControlAvailable bool
}
