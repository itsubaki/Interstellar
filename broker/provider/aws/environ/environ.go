package main

import "github.com/itsubaki/interstellar/broker"

type EnvironBroker struct {
}

func NewEnvironBroker() *EnvironBroker {
	return &EnvironBroker{}
}

func (b *EnvironBroker) Config() *broker.Config {
	return &broker.Config{
		Port: ":9083",
	}
}

func (b *EnvironBroker) Catalog() *broker.Catalog {
	return &broker.Catalog{
		Name: "aws_environ",
		Tag: []string{
			"aws",
			"environ",
		},
		Require:  []string{"aws_project"},
		Bindable: true,
		ParameterSpec: []*broker.ParamSpec{
			{Name: "project_name", Required: true},
			{Name: "environ_name", Required: true},
		},
	}
}

// ExportName is related with project_name, environ_name
// ExportValue
//  - security_group
//  - iam
func (b *EnvironBroker) Create(in *broker.CreateInput) *broker.CreateOutput {
	out := make(map[string]string)
	return &broker.CreateOutput{
		Status:  202,
		Message: "Accepted",
		Input:   in,
		Output:  out,
	}
}

func (b *EnvironBroker) Delete(in *broker.DeleteInput) *broker.DeleteOutput {
	return &broker.DeleteOutput{
		Status:  202,
		Message: "Accepted",
		Input:   in,
	}
}

func (b *EnvironBroker) Update(in *broker.UpdateInput) *broker.UpdateOutput {
	return &broker.UpdateOutput{}
}

func (b *EnvironBroker) Binding(in *broker.BindingInput) *broker.BindingOutput {
	return &broker.BindingOutput{}
}

func (b *EnvironBroker) Unbinding(in *broker.UnbindingInput) *broker.UnbindingOutput {
	return &broker.UnbindingOutput{}
}

func (b *EnvironBroker) Status(in *broker.StatusInput) *broker.StatusOutput {
	return &broker.StatusOutput{}
}