package internal

import (
	"bytes"
	"os"
	"text/template"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

type GeneratorSettings struct {
	LogFile    bool
	NetVersion string
	Namespace  string
	ClassName  string
	ProjName   string
	ExtFunc    string
}

type Generator struct {
	LogFile               *os.File
	Setting               GeneratorSettings
	CsProjTemplate        *template.Template
	ExtensionsCsTemplate  *template.Template
	DbCtxTemplate *template.Template
	QueryFileTemplate     *template.Template
	OneQueryTemplate      *template.Template
	ManyQueryTemplate     *template.Template
	ExecQueryTemplate     *template.Template
}

func New(setting GeneratorSettings) Generator {
	var logFile *os.File = nil
	if setting.LogFile {
		lf, err := os.Create("log")
		if err != nil {
			panic(err)
		}

		logFile = lf
	}

	csprpojTemp, err := template.New("proj").Parse(CsprojTemp)
	if err != nil {
		logFile.Write([]byte(err.Error()))
		panic(err)
	}

	extensionsCsTemp, err := template.New("extensionsCs").Parse(ExtensionsCsTemp)
	if err != nil {
		logFile.Write([]byte(err.Error()))
		panic(err)
	}

	dbCtxTemp, err := template.New("dbCtx").Parse(DbCtxTemp)
	if err != nil {
		logFile.Write([]byte(err.Error()))
		panic(err)
	}

	queryFileTemp, err := template.New("queryFile").Parse(QueryFileTemp)
	if err != nil {
		logFile.Write([]byte(err.Error()))
		panic(err)
	}

	oneQueryTemp, err := template.New("oneQuery").Parse(OneQueryTemp)
	if err != nil {
		logFile.Write([]byte(err.Error()))
		panic(err)
	}

	manyQueryTemp, err := template.New("manyQuery").Parse(ManyQueryTemp)
	if err != nil {
		logFile.Write([]byte(err.Error()))
		panic(err)
	}

	execQueryTemp, err := template.New("execQuery").Parse(ExecQueryTemp)
	if err != nil {
		logFile.Write([]byte(err.Error()))
		panic(err)
	}

	return Generator{
		LogFile:               logFile,
		Setting:               setting,
		CsProjTemplate:        csprpojTemp,
		ExtensionsCsTemplate:  extensionsCsTemp,
		DbCtxTemplate: dbCtxTemp,
		QueryFileTemplate:     queryFileTemp,
		OneQueryTemplate:      oneQueryTemp,
		ManyQueryTemplate:     manyQueryTemp,
		ExecQueryTemplate:     execQueryTemp,
	}
}

func (g *Generator) Generate(req *plugin.GenerateRequest) *plugin.GenerateResponse {
	files := make([]*plugin.File, 0)

	files = append(files, &plugin.File{
		Name:     g.Setting.ProjName + ".csproj",
		Contents: g.buildCsprojFileContents(),
	})

	files = append(files, &plugin.File{
		Name:     "Extensions.cs",
		Contents: g.buildExtensionsCsFileContents(),
	})

	files = append(files, &plugin.File{
		Name:     "DbCtx.cs",
		Contents: g.buildCtxFileContents(),
	})

	files = append(files, &plugin.File{
		Name:     "Query.cs",
		Contents: g.buildQueryFileConents(req.Queries),
	})

	return &plugin.GenerateResponse{
		Files: files,
	}
}

func (g *Generator) log(msg string) {
	if g.Setting.LogFile {
		g.LogFile.Write([]byte(msg))
	}
}

func (g *Generator) buildCsprojFileContents() []byte {
	var buf bytes.Buffer
	csProjSrc := CsprojTempSrc{
		NetVersion: g.Setting.NetVersion,
	}

	err := g.CsProjTemplate.Execute(&buf, &csProjSrc)
	if err != nil {
		g.log(err.Error())
		panic("error during creation of the csproj file")
	}

	return buf.Bytes()
}

func (g *Generator) buildExtensionsCsFileContents() []byte {
	var buf bytes.Buffer
	src := ExtensionsCsSrc{
		Namespace: g.Setting.Namespace,
		ClassName: g.Setting.ClassName,
		FuncName: g.Setting.ExtFunc,
	}

	err := g.ExtensionsCsTemplate.Execute(&buf, &src)
	if err != nil {
		g.log(err.Error())
		panic("error during creation of the extensions file")
	}

	return buf.Bytes()
}

func (g *Generator) buildCtxFileContents() []byte {
	var buf bytes.Buffer
	src := DbCtxSrc{
		Namespace: g.Setting.Namespace,
		Classname: g.Setting.ClassName,
	}

	err := g.DbCtxTemplate.Execute(&buf, &src)
	if err != nil {
		g.log(err.Error())
		panic("error during creation of the db ctx file")
	}

	return buf.Bytes()
}

func (g *Generator) buildQueryFileConents(queries []*plugin.Query) []byte {
	var buf bytes.Buffer

	queryContents := make([]string, 0, len(queries))

	for _, query := range queries {
		qc := g.buildQueryContents(query)
		queryContents = append(queryContents, qc)
	}

	src := QueryFileSrc{
		Namespace: g.Setting.Namespace,
		Classname: g.Setting.ClassName,
		Queries:   queryContents,
	}

	err := g.QueryFileTemplate.Execute(&buf, &src)
	if err != nil {
		g.log(err.Error())
		panic("error during creation of the query file")
	}

	return buf.Bytes()
}

func (g *Generator) buildQueryContents(query *plugin.Query) string {
	switch query.Cmd {
	case ":one":
		return g.buildOneQueryContents(query)

	case ":many":
		return g.buildManyQueryContents(query)

	case ":exec":
		return g.buildExecQueryContents(query)
	default:
		return "unknown query"
	}
}

func (g Generator) buildOneQueryContents(query *plugin.Query) string {
	var buf bytes.Buffer

	src := QuerySrc{
		Name:   query.Name,
		Text:   query.Text,
		Result: MapColumns(query.Columns),
		Params: MapParameters(query.Params),
	}

	err := g.OneQueryTemplate.Execute(&buf, &src)
	if err != nil {
		g.log(err.Error())
		panic("error during creation of the one query")
	}

	return buf.String()
}

func (g Generator) buildManyQueryContents(query *plugin.Query) string {
	var buf bytes.Buffer

	src := QuerySrc{
		Name:   query.Name,
		Text:   query.Text,
		Result: MapColumns(query.Columns),
		Params: MapParameters(query.Params),
	}

	err := g.ManyQueryTemplate.Execute(&buf, &src)
	if err != nil {
		g.log(err.Error())
		panic("error during creation of the many query")
	}

	return buf.String()
}

func (g Generator) buildExecQueryContents(query *plugin.Query) string {
	var buf bytes.Buffer

	src := QuerySrc{
		Name:   query.Name,
		Text:   query.Text,
		Params: MapParameters(query.Params),
	}

	err := g.ExecQueryTemplate.Execute(&buf, &src)
	if err != nil {
		g.log(err.Error())
		panic("error during creation of the many query")
	}

	return buf.String()
}
