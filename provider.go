// Copyright 2024 Daniel Haus <dhaus67>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openfeatureposthog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/open-feature/go-sdk/openfeature"
	"github.com/posthog/posthog-go"
)

// All supported keys in the evaluation context.
const (
	DistinctIDContextKey = openfeature.TargetingKey
	GroupsContextKey     = "groups"
	PropertiesContextKey = "properties"
)

var (
	_ openfeature.FeatureProvider = (*Provider)(nil)

	errMissingTargetKey  = errors.New("missing target key in evaluation context")
	errInvalidGroups     = errors.New("invalid groups in evaluation context")
	errInvalidProperties = errors.New("invalid properties in evaluation context")
)

type PostHogProperties struct {
	GroupProperties  map[string]posthog.Properties
	PersonProperties posthog.Properties
}

type Provider struct {
	client posthog.Client
}

// NewProvider creates a new PostHog provider.
func NewProvider(client posthog.Client) *Provider {
	return &Provider{
		client: client,
	}
}

// Metadata returns the providers metadata.
func (p *Provider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{
		Name: "PostHog",
	}
}

// Hooks returns the list of hooks of the provider. The PostHog provider does not have any hooks, an empty slice is returned.
func (p *Provider) Hooks() []openfeature.Hook {
	return []openfeature.Hook{}
}

func (p *Provider) BooleanEvaluation(_ context.Context, flag string, defaultValue bool, evalCtx openfeature.FlattenedContext) openfeature.BoolResolutionDetail {
	payload, err := translateFeatureFlagPayload(evalCtx, flag)
	if err != nil {
		if errors.Is(err, errMissingTargetKey) {
			return openfeature.BoolResolutionDetail{
				Value: defaultValue,
				ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
					ResolutionError: openfeature.NewTargetingKeyMissingResolutionError(err.Error()),
					Reason:          openfeature.ErrorReason,
				},
			}
		}

		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	res, found, err := p.getFeatureFlag(payload)
	if err != nil {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if !found {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(fmt.Sprintf("%q not found", flag)),
				Reason:          openfeature.DefaultReason,
			},
		}
	}

	parsedValue, resolutionErr := parseFlagValue[bool](res)
	if resolutionErr != nil {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: *resolutionErr,
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	return openfeature.BoolResolutionDetail{
		Value: parsedValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason: openfeature.TargetingMatchReason,
		},
	}
}

func (p *Provider) FloatEvaluation(_ context.Context, flag string, defaultValue float64, evalCtx openfeature.FlattenedContext) openfeature.FloatResolutionDetail {
	payload, err := translateFeatureFlagPayload(evalCtx, flag)
	if err != nil {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	res, found, err := p.getFeatureFlag(payload)
	if err != nil {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if !found {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(fmt.Sprintf("%q not found", flag)),
				Reason:          openfeature.DefaultReason,
			},
		}
	}

	parsedValue, resolutionErr := parseFlagValue[float64](res)
	if resolutionErr != nil {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: *resolutionErr,
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	return openfeature.FloatResolutionDetail{
		Value: parsedValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason: openfeature.TargetingMatchReason,
		},
	}
}

func (p *Provider) IntEvaluation(_ context.Context, flag string, defaultValue int64, evalCtx openfeature.FlattenedContext) openfeature.IntResolutionDetail {
	payload, err := translateFeatureFlagPayload(evalCtx, flag)
	if err != nil {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	res, found, err := p.getFeatureFlag(payload)
	if err != nil {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if !found {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(fmt.Sprintf("%q not found", flag)),
				Reason:          openfeature.DefaultReason,
			},
		}
	}

	parsedValue, resolutionErr := parseFlagValue[int64](res)
	if resolutionErr != nil {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: *resolutionErr,
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	return openfeature.IntResolutionDetail{
		Value: parsedValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason: openfeature.TargetingMatchReason,
		},
	}
}

func (p *Provider) ObjectEvaluation(_ context.Context, flag string, defaultValue interface{}, evalCtx openfeature.FlattenedContext) openfeature.InterfaceResolutionDetail {
	payload, err := translateFeatureFlagPayload(evalCtx, flag)
	if err != nil {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	res, found, err := p.getFeatureFlag(payload)
	if err != nil {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if !found {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(fmt.Sprintf("%q not found", flag)),
				Reason:          openfeature.DefaultReason,
			},
		}
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(res.(string)), &obj); err != nil {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewTypeMismatchResolutionError("invalid JSON as flag value"),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	return openfeature.InterfaceResolutionDetail{
		Value: obj,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason: openfeature.TargetingMatchReason,
		},
	}
}

func (p *Provider) StringEvaluation(_ context.Context, flag string, defaultValue string, evalCtx openfeature.FlattenedContext) openfeature.StringResolutionDetail {
	payload, err := translateFeatureFlagPayload(evalCtx, flag)
	if err != nil {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	res, found, err := p.getFeatureFlag(payload)
	if err != nil {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewGeneralResolutionError(err.Error()),
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	if !found {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(fmt.Sprintf("%q not found", flag)),
				Reason:          openfeature.DefaultReason,
			},
		}
	}

	parsedValue, resolutionErr := parseFlagValue[string](res)
	if resolutionErr != nil {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: *resolutionErr,
				Reason:          openfeature.ErrorReason,
			},
		}
	}

	return openfeature.StringResolutionDetail{
		Value: parsedValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			Reason: openfeature.TargetingMatchReason,
		},
	}
}

func (p *Provider) getFeatureFlag(payload posthog.FeatureFlagPayload) (interface{}, bool, error) {
	res, err := p.client.GetFeatureFlag(payload)
	if err != nil {
		return res, false, err
	}

	// In case the flag could not be found, the client will return false for the flag value.
	// This is the only time where it will return a boolean value instead of a string value.
	if res, ok := res.(bool); ok {
		return res, false, nil
	}

	return res, true, nil
}

func translateFeatureFlagPayload(evalCtx openfeature.FlattenedContext, key string) (posthog.FeatureFlagPayload, error) {
	distinctID, ok := evalCtx[DistinctIDContextKey]
	if !ok {
		return posthog.FeatureFlagPayload{}, errMissingTargetKey
	}

	groups, ok := evalCtx[GroupsContextKey]
	var postHogGroups posthog.Groups
	if ok {
		postHogGroups, ok = groups.(posthog.Groups)
		if !ok {
			return posthog.FeatureFlagPayload{}, errInvalidGroups
		}
	}

	context, ok := evalCtx[PropertiesContextKey]
	var postHogCtx PostHogProperties
	if ok {
		postHogCtx, ok = context.(PostHogProperties)
		if !ok {
			return posthog.FeatureFlagPayload{}, errInvalidProperties
		}
	}

	return posthog.FeatureFlagPayload{
		Key:              key,
		DistinctId:       distinctID.(string),
		Groups:           postHogGroups,
		PersonProperties: postHogCtx.PersonProperties,
		GroupProperties:  postHogCtx.GroupProperties,
	}, nil
}

func parseFlagValue[T any](v interface{}) (T, *openfeature.ResolutionError) {
	// The API response is always a string.
	s, ok := v.(string)
	if !ok {
		err := openfeature.NewTypeMismatchResolutionError(fmt.Sprintf("%v is not a string", v))
		return *new(T), &err
	}

	switch any(*new(T)).(type) {
	case bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			err := openfeature.NewTypeMismatchResolutionError(fmt.Sprintf("%q is not a boolean", v))
			return *new(T), &err
		}
		return any(b).(T), nil
	case float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			err := openfeature.NewTypeMismatchResolutionError(fmt.Sprintf("%q is not a float", v))
			return *new(T), &err
		}
		return any(f).(T), nil
	case int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			err := openfeature.NewTypeMismatchResolutionError(fmt.Sprintf("%q is not a int", v))
			return *new(T), &err
		}
		return any(i).(T), nil
	case string:
		return any(v).(T), nil
	default:
		err := openfeature.NewTypeMismatchResolutionError(fmt.Sprintf("unsupported type %T", *new(T)))
		return *new(T), &err
	}
}
