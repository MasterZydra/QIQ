package phpt

type TestFile struct {
	Title      string
	GetParams  [][]string
	PostParams [][]string
	Env        map[string]string
	File       string
	Expect     string
}

func NewTestFile() *TestFile {
	return &TestFile{}
}
