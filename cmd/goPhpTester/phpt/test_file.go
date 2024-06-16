package phpt

type TestFile struct {
	Title      string
	GetParams  [][]string
	PostParams [][]string
	Env        map[string]string
	Ini        []string
	File       string
	Expect     string
	ExpectType string
}

func NewTestFile() *TestFile {
	return &TestFile{}
}
