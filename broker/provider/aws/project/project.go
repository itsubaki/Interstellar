package main

import "github.com/itsubaki/interstellar/broker"

type ProjectBroker struct {
}

func NewProjectBroker() *ProjectBroker {
	return &ProjectBroker{}
}

func (b *ProjectBroker) Config() *broker.Config {
	return &broker.Config{
		Port: ":9084",
	}
}

func (b *ProjectBroker) Catalog() *broker.Catalog {
	return &broker.Catalog{
		Name: "aws_project",
		Tag: []string{
			"aws",
			"project",
		},
		Require:  []string{},
		Bindable: false,
		ParameterSpec: []*broker.ParamSpec{
			{Name: "integration_role_arn", Required: true},
			{Name: "project_name", Required: true},
			{Name: "cidr", Required: true},
			{Name: "domain", Required: true},
		},
	}
}

// create s3bucket, vpc, subnet, certificate, hostedzone
// ExportName is related with project_name
// ExportValue
//  - integration_role_arn
//  - cidr
//  - subnet
//  - domain
//  - bucket_name
func (b *ProjectBroker) Create(in *broker.CreateInput) *broker.CreateOutput {
	out := make(map[string]string)
	out["nameserver"] = "ns-1,ns-2,ns-3,ns-4"
	out["bucket_log"] = "s3://log.${project_name}.${domain}"
	out["bucket_deploy"] = "s3://deploy.${project_name}.${domain}"

	return &broker.CreateOutput{
		Status:  202,
		Message: "Accepted",
		Input:   in,
		Output:  out,
	}
}

func (b *ProjectBroker) Delete(in *broker.DeleteInput) *broker.DeleteOutput {
	return &broker.DeleteOutput{
		Status:  202,
		Message: "Accepted",
		Input:   in,
	}
}

func (b *ProjectBroker) Update(in *broker.UpdateInput) *broker.UpdateOutput {
	return &broker.UpdateOutput{}
}

func (b *ProjectBroker) Binding(in *broker.BindingInput) *broker.BindingOutput {
	return &broker.BindingOutput{}
}

func (b *ProjectBroker) Unbinding(in *broker.UnbindingInput) *broker.UnbindingOutput {
	return &broker.UnbindingOutput{}
}

func (b *ProjectBroker) Status(in *broker.StatusInput) *broker.StatusOutput {
	return &broker.StatusOutput{}
}
