package phpt

type TestFile struct {
	Title      string
	GetParams  map[string][]string
	PostParams map[string][]string
	File       string
	Expect     string
}

func NewTestFile() *TestFile {
	return &TestFile{}
}
