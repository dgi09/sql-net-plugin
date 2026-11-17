package main

import (
	"context"
	"encoding/json"
	"net-plugin/internal"

	"github.com/sqlc-dev/plugin-sdk-go/codegen"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

func main() {
	codegen.Run(generate)
}

type CodeGenOptions struct {
	Class      string `json:"class"`
	Namespace  string `json:"namespace"`
	NetVersion string `json:"netVersion"`
	ProjName   string `json:"projName"`
	ExtFunc    string `json:"extFunc"`
}

func generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	codeGenOptsData := req.Settings.GetCodegen().GetOptions()

	var codeGenOpts CodeGenOptions
	err := json.Unmarshal(codeGenOptsData, &codeGenOpts)
	if err != nil {
		panic(err)
	}

	genSettings := internal.GeneratorSettings{
		LogFile:    false,
		ProjName:   codeGenOpts.ProjName,
		NetVersion: codeGenOpts.NetVersion,
		Namespace:  codeGenOpts.Namespace,
		ClassName:  codeGenOpts.Class,
		ExtFunc:    codeGenOpts.ExtFunc,
	}

	generator := internal.New(genSettings)

	return generator.Generate(req), nil
}
