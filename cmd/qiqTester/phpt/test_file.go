package phpt

type TestFile struct {
	Filename   string
	Title      string
	Get        string
	Post       string
	Args       [][]string
	Env        map[string]string
	Cookie     string
	Ini        []string
	File       string
	Expect     string
	ExpectType string
}

func NewTestFile(filename string) *TestFile {
	return &TestFile{Filename: filename, Env: map[string]string{}}
}
