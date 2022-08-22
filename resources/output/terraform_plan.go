package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"io"
)

type TfPlan struct {
	Messages []TfPlanMessage
}

type TfPlanAction string

const (
	Create  TfPlanAction = "create"
	Update               = "update"
	Replace              = "replace"
	Delete               = "delete"
)

const plannedChangedType = "planned_change"

type TfPlanMessage struct {
	Change TfPlannedChange `json:"change"`
	Type   string          `json:"type"`
}

type TfPlannedChange struct {
	Resource TfPlanResource `json:"resource"`
	Action   TfPlanAction   `json:"action"`
}

type TfPlanResource struct {
	Addr         string `json:"addr"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	ResourceKey  string `json:"resource_key"`
}

func ParsePlanFromOutput(outputJson string) (*TfPlan, error) {
	var out []TfPlanMessage
	b := bytes.NewBufferString(outputJson)
	line, err := b.ReadString('\n')
	for ; err == nil; line, err = b.ReadString('\n') {
		elem := TfPlanMessage{}
		err = json.Unmarshal([]byte(line), &elem)
		if err != nil {
			return nil, err
		}
		out = append(out, elem)
	}

	if err == io.EOF {
		return &TfPlan{Messages: out}, nil
	}

	return nil, err
}

func (t *TfPlan) MaybeGetPlannedChange(resourceRef string) *TfPlannedChange {
	for _, r := range t.Messages {
		if r.Type != plannedChangedType {
			continue
		}
		changedResource := r.Change.Resource
		if resourceRef == fmt.Sprintf("%s.%s", changedResource.ResourceType, changedResource.ResourceName) {
			return &r.Change
		}
	}

	return nil
}

func MaybeGetPlannedChageById[T any](plan *TfPlan, resourceId string) *TfPlannedChange {
	resourceRef := fmt.Sprintf("%s.%s", GetResourceName(*new(T)), resourceId)
	return plan.MaybeGetPlannedChange(resourceRef)
}

func AddToStatuses(statuses map[string]commonpb.ResourceStatus_Status, key string, plannedChange *TfPlannedChange) {
	if plannedChange != nil {
		statuses[key] = convertActionToStatus(plannedChange.Action)
	}
}

func convertActionToStatus(action TfPlanAction) commonpb.ResourceStatus_Status {
	switch action {
	case Create:
		return commonpb.ResourceStatus_NEEDS_CREATE
	case Update:
		return commonpb.ResourceStatus_NEEDS_UPDATE
	case Replace:
		return commonpb.ResourceStatus_NEEDS_RECREATE
	}

	return commonpb.ResourceStatus_UKNOWN_STATUS
}
