package phpt

type TestFile struct {
	Filename   string
	Title      string
	GetParams  [][]string
	PostParams [][]string
	Args       [][]string
	Env        map[string]string
	Ini        []string
	File       string
	Expect     string
	ExpectType string
}

func NewTestFile(filename string) *TestFile {
	return &TestFile{Filename: filename}
}
