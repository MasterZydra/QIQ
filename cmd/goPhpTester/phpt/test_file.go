package phpt

type TestFile struct {
	Title      string
	GetParams  [][]string
	PostParams [][]string
	File       string
	Expect     string
}

func NewTestFile() *TestFile {
	return &TestFile{}
}
