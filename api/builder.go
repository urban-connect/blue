package api

import (
	"encoding/json"
	"fmt"
	"regexp"

	"tinygo.org/x/bluetooth"
	"urban-connect.ch/blue/detection/filters"
)

type Filter struct {
	Kind  string          `json:"kind"`
	Props json.RawMessage `json:"props"`
}

type AllOfProps struct {
	Filters []Filter `json:"filters"`
}

type OneOfProps struct {
	Filters []Filter `json:"filters"`
}

type NameProps struct {
	Pattern string `json:"pattern"`
}

type ServiceProps struct {
	UUID string `json:"uuid"`
}

type Builder struct{}

func (builder *Builder) Build(filter Filter) (filters.Filter, error) {
	switch filter.Kind {
	case "all_of":
		f, err := builder.buildAllOf(filter.Props)

		if err != nil {
			return nil, fmt.Errorf("failed to all_of filter: %w", err)
		}

		return f, nil
	case "one_of":
		f, err := builder.buildOneOf(filter.Props)

		if err != nil {
			return nil, fmt.Errorf("failed to one_of filter: %w", err)
		}

		return f, nil
	case "name":
		f, err := builder.buildName(filter.Props)

		if err != nil {
			return nil, fmt.Errorf("failed to build name filter: %w", err)
		}

		return f, nil
	case "service":
		f, err := builder.buildService(filter.Props)

		if err != nil {
			return nil, fmt.Errorf("failed to build service filter: %w", err)
		}

		return f, nil
	default:
		return nil, fmt.Errorf("unknown filter kind: %s", filter.Kind)
	}
}

func (buidler *Builder) buildAllOf(raw json.RawMessage) (filters.Filter, error) {
	var parsed AllOfProps

	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse props: %w", err)
	}

	var args []filters.Filter

	for _, filter := range parsed.Filters {
		f, err := buidler.Build(filter)

		if err != nil {
			return nil, fmt.Errorf("failed to build an arg filter: %w", err)
		}

		args = append(args, f)
	}

	return filters.AllOf(args...), nil
}

func (buidler *Builder) buildOneOf(raw json.RawMessage) (filters.Filter, error) {
	var parsed OneOfProps

	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse props: %w", err)
	}

	var args []filters.Filter

	for _, filter := range parsed.Filters {
		f, err := buidler.Build(filter)

		if err != nil {
			return nil, fmt.Errorf("failed to build an arg filter: %w", err)
		}

		args = append(args, f)
	}

	return filters.OneOf(args...), nil
}

func (builder *Builder) buildService(raw json.RawMessage) (filters.Filter, error) {
	var parsed ServiceProps

	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse props: %w", err)
	}

	uuid, err := bluetooth.ParseUUID(parsed.UUID)

	if err != nil {
		return nil, fmt.Errorf("failed to parse UUID: %w", err)
	}

	return filters.Service(uuid), nil
}

func (builder *Builder) buildName(raw json.RawMessage) (filters.Filter, error) {
	var parsed NameProps

	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse props: %w", err)
	}

	pattern, err := regexp.Compile(parsed.Pattern)

	if err != nil {
		return nil, fmt.Errorf("failed to parse pattern: %w", err)
	}

	return filters.Name(pattern), nil
}
