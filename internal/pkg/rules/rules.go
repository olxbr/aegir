package rules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	y2j "github.com/ghodss/yaml"
	"github.com/grupozap/aegir/internal/pkg/utils"
	livr "github.com/k33nice/go-livr"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

type RulesList struct {
	Rules []*Rule `yaml:"rules"`
}

type Rule struct {
	Name                     string           `yaml:"name"`
	Namespace                string           `yaml:"namespace"`
	ResourceType             string           `yaml:"resource_type"`
	RulesDefinitions         []RuleDefinition `yaml:"rules_definitions"`
	SlackNotificationChannel string           `yaml:"slack_notification_channel,omitempty"`
}

type RuleDefinition struct {
	Field           string     `yaml:"field"`
	FieldIsOptional bool       `yaml:"field_is_optional"`
	LivrRule        RuleObject `yaml:"livr_rule"`
}

type RuleObject struct {
	Description string          `yaml:"description"`
	RuleObj     livr.Dictionary `yaml:"rule"`
}

type Resource struct {
	Kind     map[string]interface{} `json:"kind"`
	Metadata map[string]interface{} `json:"metadata"`
	Spec     map[string]interface{} `json:"spec"`
}

var ruleStore = map[string][]*Rule{}

func RulesLoader(fp string) RulesList {
	file, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Fatalf("could not read file: %q", err)
	}
	var rules RulesList
	err = yaml.Unmarshal(file, &rules)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	return rules
}

func createKey(ns, rt string) string {
	return fmt.Sprintf("%s/%s", ns, rt)
}

func BuildRuleStore(rl *RulesList) {
	for _, rule := range rl.Rules {
		k := createKey(rule.Namespace, rule.ResourceType)
		if _, ok := ruleStore[k]; !ok {
			ruleStore[k] = []*Rule{}
		}
		ruleStore[k] = append(ruleStore[k], rule)
	}
}

func (ruledef *RuleDefinition) registerRule() *livr.Validator {
	var rule map[string]interface{}
	r, _ := yaml.Marshal(ruledef.LivrRule.RuleObj)
	j, err := y2j.YAMLToJSON(r)
	if err != nil {
		log.Fatalf("something wrong when converting YAML to JSON, error: %v", err)
	}
	err = json.Unmarshal(j, &rule)
	if err != nil {
		log.Fatalf("something wrong unmarshaling JSON to LIVR")
	}
	return livr.New(&livr.Options{LivrRules: rule})
}

func (ruledef *RuleDefinition) GetViolations(obj string) []*utils.Violation {
	violations := make([]*utils.Violation, 0)
	jp := gjson.Get(obj, ruledef.Field)
	//If field is optional, don't check if it exists.
	if !ruledef.FieldIsOptional {
		if !jp.Exists() || (jp.IsArray() && len(jp.Array()) == 0) {
			fieldNotFound := &utils.Violation{
				Description: ruledef.LivrRule.Description,
				JSONPath:    ruledef.Field,
				Message:     fmt.Sprintf("Field: %s is required", ruledef.Field),
			}
			violations = append(violations, fieldNotFound)
		}
	}
	for _, jsonobj := range GetJSONObjectByPath(obj, ruledef.Field) {
		validator := ruledef.registerRule()
		lastfield := utils.GetLastField(ruledef.Field)
		objmap := make(map[string]interface{})
		objmap[lastfield] = jsonobj.Value()
		_, err := validator.Validate(objmap)
		if err != nil {
			v := &utils.Violation{
				Description: ruledef.LivrRule.Description,
				JSONPath:    ruledef.Field,
				Object:      objmap,
				Message:     err.Error(),
			}
			violations = append(violations, v)
		}
	}
	return violations
}

func GetJSONObjectByPath(Obj, JSONPath string) []gjson.Result {
	Objects := make([]gjson.Result, 0)
	result := gjson.Get(Obj, JSONPath)
	if result.IsArray() {
		result.ForEach(func(key, value gjson.Result) bool {
			if value.IsArray() {
				Objects = append(Objects, value.Array()...)
			} else {
				Objects = append(Objects, value)
			}
			return true
		})
	} else if result.Exists() {
		Objects = append(Objects, result)
		return Objects
	}
	return Objects
}

func GetRules(ns, rt string) []*Rule {
	rs := ruleStore[createKey(ns, rt)]
	if rs == nil {
		rs = []*Rule{}
	}
	return append(rs, ruleStore[createKey("*", rt)]...)
}
