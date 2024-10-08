package extractor

const upsertTmplStruct = `

// Upsert{{.StructName}}Input struct
type Upsert{{.StructName}}Input struct {
{{- range .Fields}}
	{{- if and (not (eq (lower .Name) "uid")) (not (eq (lower .Name) "dtype"))}}
		{{- if .IsComplexType}}
			{{- if .IsSliceType}}
				{{.Name}} []string
			{{- else}}
				{{.Name}} string
			{{- end}}
		{{- else}}
			{{.Name}} {{.Type}}
		{{- end}}
	{{- end}}
{{- end}}
}

// Upsert{{$.StructName}} method
func Upsert{{$.StructName}}(input Upsert{{$.StructName}}Input) (map[string]string, error) {
	// Convert input to the original struct type
{{- if $.HasComplexFields}}
		type Temp{{$.StructName}} struct {
	{{- range .Fields}}
		{{- if .IsComplexType}}
			{{- if .IsSliceType}}
				{{.Name}} []string ` + "`{{.Tag}}`" + `
			{{- else}}
				{{.Name}} string ` + "`{{.Tag}}`" + `
			{{- end}}
		{{- else}}
			{{.Name}} {{.Type}} ` + "`{{.Tag}}`" + `
		{{- end}}
	{{- end}}
		}
		
		p := Temp{{$.StructName}}{
	{{- range .Fields}}
			{{- if not (eq (lower .Name) "uid")}}
			{{- if eq (.Name) "DType"}}
			{{.Name}}: []string{"{{$.StructName}}"},
			{{- else}}
			{{.Name}}: input.{{.Name}},
			{{- end}}
			{{- end}}
			{{- if eq (lower .Name) "uid"}}
			{{.Name}}: "_:{{lower $.StructName}}",
			{{- end}}
	{{- end}}
		}
{{- else}}
		p := {{$.StructName}}{
	{{- range .Fields}}
		{{- if not (eq (lower .Name) "uid")}}
		{{- if eq (.Name) "DType"}}
		{{.Name}}: []string{"{{$.StructName}}"},
		{{- else}}
		{{.Name}}: input.{{.Name}},
		{{- end}}
		{{- end}}
		{{- if eq (lower .Name) "uid"}}
		{{.Name}}: "_:{{lower $.StructName}}",
		{{- end}}
	{{- end}}
		}
{{- end}}

	// Convert the struct to JSON
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	
	response, err := dgraph.Execute(hostName, &dgraph.Request{
		Mutations: []*dgraph.Mutation{
			{
				SetJson: string(jsonBytes),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return response.Uids, nil
}


`

const deleteTmplStruct = `
// Delete{{.StructName}}Input struct
type Delete{{.StructName}}Input struct {
	Uid string
}

// Delete{{$.StructName}} method
func Delete{{$.StructName}}(input Delete{{$.StructName}}Input) (string, error) {
	statement := fmt.Sprintf("<%s> * * .", input.Uid)

	_, err := dgraph.Execute(hostName, &dgraph.Request{
		Mutations: []*dgraph.Mutation{
			{
				DelNquads: statement,
			},
		},
	})
	if err != nil {
		return "", err
	}

	return "success", nil
}


`

const getTmplStruct = `
// Get{{.StructName}} method
func Get{{.StructName}}(uid string) (*{{.StructName}}, error) {
	statement := ` + "`" +
	`
	query query{{.StructName}}($uid: string) {
		{{lower .StructName}}(func: uid($uid)) {
			uid
			dgraph.type
{{- range .Fields}}
			{{- if and (not (eq (lower .Name) "uid")) (not (eq (lower .Name) "dtype")) }}
			{{- if .IsComplexType}}
			{{lower .Name}} {
				uid
				dgraph.type
				expand(_all_)
			}
			{{- else}}
			{{lower .Name}}
			{{- end}}
			{{- end}}
{{- end}}
		}
	}
` + "`" + `
	variables := map[string]string{"$uid": uid}
	response, err := dgraph.Execute(hostName, &dgraph.Request{
		Query: &dgraph.Query{
			Query:    statement,
			Variables: variables,
		},
	})
	if err != nil {
		return nil, err
	}
	type {{.StructName}}Data struct {
		{{.StructName}}s []*{{.StructName}} ` + "`json:\"{{lower .StructName}}\"`" + `
	}
	var {{lower .StructName}}Data {{.StructName}}Data
	if err := json.Unmarshal([]byte(response.Json), &{{lower .StructName}}Data); err != nil {
		return nil, err
	}

	if len({{lower .StructName}}Data.{{.StructName}}s) == 0 {
		return nil, fmt.Errorf("{{.StructName}} not found")
	}

	return {{lower .StructName}}Data.{{.StructName}}s[0], nil
}


`

const queryTmplStruct = `
// Query{{.StructName}} method
func Query{{.StructName}}(limit, offset int) ([]*{{.StructName}}, error) {
	statement := ` + "`" +
	`
	query query{{.StructName}}($limit: int, $offset: int) {
		{{lower .StructName}}(func: type({{.StructName}}), first: $limit, offset: $offset) {
			uid
			dgraph.type
{{- range .Fields}}
			{{- if and (not (eq (lower .Name) "uid")) (not (eq (lower .Name) "dtype")) }}
			{{- if .IsComplexType}}
			{{lower .Name}} {
				uid
				dgraph.type
				expand(_all_)
			}
			{{- else}}
			{{lower .Name}}
			{{- end}}
			{{- end}}
{{- end}}
		}
	}
` + "`" + `
	variables := map[string]string{"$limit": fmt.Sprintf("%d", limit), "$offset": fmt.Sprintf("%d", offset)}
	response, err := dgraph.Execute(hostName, &dgraph.Request{
		Query: &dgraph.Query{
			Query:     statement,
			Variables: variables,
		},
	})
	if err != nil {
		return nil, err
	}
	type {{.StructName}}Data struct {
		{{.StructName}}s []*{{.StructName}} ` + "`json:\"{{lower .StructName}}\"`" + `
	}
	var {{lower .StructName}}Data {{.StructName}}Data
	if err := json.Unmarshal([]byte(response.Json), &{{lower .StructName}}Data); err != nil {
		return nil, err
	}

	if len({{lower .StructName}}Data.{{.StructName}}s) == 0 {
		return nil, fmt.Errorf("{{.StructName}} not found")
	}

	return {{lower .StructName}}Data.{{.StructName}}s, nil
}
`
