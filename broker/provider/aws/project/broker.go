package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/itsubaki/env"
	"github.com/itsubaki/interstellar/broker"
	"github.com/itsubaki/interstellar/broker/repository"
	"github.com/itsubaki/required"
)

type ProjectBroker struct {
	config   *broker.Config
	catalog  *broker.Catalog
	required *required.Map
	template string
	instance repository.InstanceRepository
}

func NewProjectBroker() (*ProjectBroker, error) {
	catalog := &broker.Catalog{
		Name: "aws_project",
		Tag: []string{
			"aws",
			"project",
		},
		Bindable: false,
		ParameterSpec: []broker.ParamSpec{
			{Name: "project_name", Required: true},
			{Name: "region", Required: true},
			{Name: "domain", Required: true},
			{Name: "cidr", Required: false},
		},
	}

	config := &broker.Config{
		Port:     env.GetValue("PORT", ":8080"),
		Template: env.GetValue("TEMPLATE", "./template.yml"),
	}

	req := required.NewMap()
	for _, p := range catalog.ParameterSpec {
		req.Put(p.Name, p.Required)
	}

	f, err := ioutil.ReadFile(config.Template)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	b := &ProjectBroker{
		config:   config,
		catalog:  catalog,
		required: req,
		template: string(f),
		instance: repository.InstanceRepository{},
	}

	return b, nil
}

func (b *ProjectBroker) Config() *broker.Config {
	return b.config
}

func (b *ProjectBroker) Catalog() *broker.Catalog {
	return b.catalog
}

func (b *ProjectBroker) Create(in *broker.CreateInput) *broker.CreateOutput {
	if err := b.required.Satisfy(in.Parameter); err != nil {
		return &broker.CreateOutput{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("invalid parameter: %v", err),
		}
	}

	sess := session.Must(session.NewSession())
	cfn := cloudformation.New(sess, &aws.Config{Region: aws.String(in.Parameter["region"])})

	param := []*cloudformation.Parameter{
		{ParameterKey: aws.String("ProjectName"), ParameterValue: aws.String(in.Parameter["project_name"])},
		{ParameterKey: aws.String("DomainName"), ParameterValue: aws.String(in.Parameter["domain"])},
		{ParameterKey: aws.String("Region"), ParameterValue: aws.String(in.Parameter["region"])},
	}

	name := fmt.Sprintf("%s-%s", in.Parameter["project_name"], in.InstanceID)
	input := &cloudformation.CreateStackInput{
		StackName:    &name,
		Parameters:   param,
		TemplateBody: &b.template,
	}

	if _, err := cfn.CreateStack(input); err != nil {
		return &broker.CreateOutput{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("create stack: %v", err),
		}
	}

	i := &broker.Instance{
		InstanceID: in.InstanceID,
		Parameter:  in.Parameter,
	}

	b.instance.Insert(i)

	go func() {
		input := &cloudformation.DescribeStacksInput{
			StackName: &name,
		}

		if err := cfn.WaitUntilStackCreateComplete(input); err != nil {
			log.Printf("wait until stack create complete %s: %v", in.InstanceID, err)
			return
		}

		desc, err := cfn.DescribeStacks(input)
		if err != nil {
			log.Printf("desctibe stack %s: %v", in.InstanceID, err)
			return
		}

		out := make(map[string]string)
		for i := range desc.Stacks[0].Outputs {
			o := desc.Stacks[0].Outputs[i]
			out[*o.OutputKey] = *o.OutputValue
		}

		i.Status = *desc.Stacks[0].StackStatus
		i.Output = out

		if err := b.instance.Update(i); err != nil {
			log.Printf("update instance_id=%s: %v", in.InstanceID, err)
		}
	}()

	return &broker.CreateOutput{
		Status:   http.StatusAccepted,
		Message:  "Accepted",
		Instance: i,
	}
}

func (b *ProjectBroker) Delete(in *broker.DeleteInput) *broker.DeleteOutput {
	return &broker.DeleteOutput{
		Status:  http.StatusAccepted,
		Message: "Accepted",
	}
}

func (b *ProjectBroker) Update(in *broker.UpdateInput) *broker.UpdateOutput {
	return &broker.UpdateOutput{
		Status:  http.StatusAccepted,
		Message: "Accepted",
	}
}

func (b *ProjectBroker) Binding(in *broker.BindingInput) *broker.BindingOutput {
	return &broker.BindingOutput{
		Status:  http.StatusAccepted,
		Message: "Accepted",
	}
}

func (b *ProjectBroker) Unbinding(in *broker.UnbindingInput) *broker.UnbindingOutput {
	return &broker.UnbindingOutput{
		Status:  http.StatusAccepted,
		Message: "Accepted",
	}
}

func (b *ProjectBroker) Describe(in *broker.DescribeInput) *broker.DescribeOutput {
	i, ok := b.instance.FindByID(in.InstanceID)
	if !ok {
		return &broker.DescribeOutput{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("instance=%s not found", in.InstanceID),
		}
	}

	return &broker.DescribeOutput{
		Status:   http.StatusOK,
		Instance: i,
	}
}
