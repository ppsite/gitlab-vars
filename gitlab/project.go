package gitlab

import (
	"encoding/json"
)

const (
  ProjectVariablesApi = "/projects/:id/variables"
)

type ProjectVariable struct {
  VariableType string       `json:"variable_type,omitempty"`
  Key string                `json:"key"`
  Value string              `json:"value"`
  Protected bool            `json:"protected,omitempty"`
  Masked bool               `json:"masked,omitempty"`
  EnvironmentScope string   `mapstructure:"environment_scope,omitempty"`
}

type ProjectVariablesList struct {
  Items []ProjectVariable
}

func (g *Gitlab) ListProjectVaribles(projectId string) (*ProjectVariablesList, error) {
  params := make(map[string]string)
  params [":id"] = projectId
  url := g.ResourceUrl(ProjectVariablesApi, params)
  varList := new(ProjectVariablesList)
  contents, err := g.buildAndExecRequest("GET", url.String(), nil)
  if err == nil {
    err = json.Unmarshal(contents, &varList.Items)
  }
  return varList, err
}


func (g *Gitlab) CreateProjectVariable(projectId string, variable *ProjectVariable) (*ProjectVariable, error) {
  params := make(map[string]string)
  params [":id"] = projectId
  url := g.ResourceUrl(ProjectVariablesApi, params)

  variableJson, err := json.Marshal(variable)
  if err != nil {
    return nil, err
  }

  var createdVariable = new(ProjectVariable)
  contents, err := g.buildAndExecRequest("POST", url.String(), variableJson)
  if err == nil {
    err = json.Unmarshal(contents, &createdVariable)
  }
  return createdVariable, err
}
