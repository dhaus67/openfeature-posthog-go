package provider

import (
	"context"
	"testing"

	"github.com/open-feature/go-sdk/openfeature"
	"github.com/posthog/posthog-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_BooleanEvaluation(t *testing.T) {
	tcs := map[string]struct {
		flag         string
		defaultValue bool
		evalCtx      openfeature.EvaluationContext
		mockSettings mockSettings
		err          bool
		res          openfeature.BooleanEvaluationDetails
	}{
		"existing boolean flag": {
			flag:         "bool-flag",
			defaultValue: false,
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "bool-flag",
					DistinctId: "12345",
				},
				res: "true",
			},
			res: openfeature.BooleanEvaluationDetails{
				Value: true,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "bool-flag",
					FlagType: openfeature.Boolean,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.TargetingMatchReason,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
		"invalid format boolean flag": {
			flag:         "bool-flag",
			defaultValue: false,
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "bool-flag",
					DistinctId: "12345",
				},
				res: "invalid boolean",
			},
			err: true,
			res: openfeature.BooleanEvaluationDetails{
				Value: false,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "bool-flag",
					FlagType: openfeature.Boolean,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.ErrorReason,
						ErrorCode:    openfeature.TypeMismatchCode,
						ErrorMessage: `"invalid boolean" is not a boolean`,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
	}

	mockClient := &mockPostHogClient{t: t}
	p := NewProvider(mockClient)
	require.NoError(t, openfeature.SetProvider(p))
	c := openfeature.NewClient("testing")

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockClient.settings = tc.mockSettings
			res, err := c.BooleanValueDetails(context.Background(), tc.flag, tc.defaultValue, tc.evalCtx)
			assert.Equal(t, tc.err, err != nil)
			assert.Equal(t, tc.res, res)
		})
	}
}

func TestProvider_FloatEvaluation(t *testing.T) {
	tcs := map[string]struct {
		flag         string
		defaultValue float64
		evalCtx      openfeature.EvaluationContext
		mockSettings mockSettings
		err          bool
		res          openfeature.FloatEvaluationDetails
	}{
		"existing float flag": {
			flag:         "float-flag",
			defaultValue: 0.10,
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "float-flag",
					DistinctId: "12345",
				},
				res: "0.55",
			},
			res: openfeature.FloatEvaluationDetails{
				Value: 0.55,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "float-flag",
					FlagType: openfeature.Float,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.TargetingMatchReason,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
		"non-existing float flag": {
			flag:         "float-flag",
			defaultValue: 0.10,
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "float-flag",
					DistinctId: "12345",
				},
				res: "false",
			},
			err: true,
			res: openfeature.FloatEvaluationDetails{
				Value: 0.10,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "float-flag",
					FlagType: openfeature.Float,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.ErrorReason,
						ErrorCode:    openfeature.FlagNotFoundCode,
						ErrorMessage: `"float-flag" not found`,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
		"invalid format float flag": {
			flag:         "float-flag",
			defaultValue: 0.10,
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "float-flag",
					DistinctId: "12345",
				},
				res: "invalid float",
			},
			err: true,
			res: openfeature.FloatEvaluationDetails{
				Value: 0.10,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "float-flag",
					FlagType: openfeature.Float,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.ErrorReason,
						ErrorCode:    openfeature.TypeMismatchCode,
						ErrorMessage: `"invalid float" is not a float`,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
	}

	mockClient := &mockPostHogClient{t: t}
	p := NewProvider(mockClient)
	require.NoError(t, openfeature.SetProvider(p))
	c := openfeature.NewClient("testing")

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockClient.settings = tc.mockSettings
			res, err := c.FloatValueDetails(context.Background(), tc.flag, tc.defaultValue, tc.evalCtx)
			assert.Equal(t, tc.err, err != nil)
			assert.Equal(t, tc.res, res)
		})
	}
}

func TestProvider_StringEvaluation(t *testing.T) {
	tcs := map[string]struct {
		flag         string
		defaultValue string
		evalCtx      openfeature.EvaluationContext
		mockSettings mockSettings
		err          bool
		res          openfeature.StringEvaluationDetails
	}{
		"existing string flag": {
			flag:         "string-flag",
			defaultValue: "test",
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "string-flag",
					DistinctId: "12345",
				},
				res: "another test",
			},
			res: openfeature.StringEvaluationDetails{
				Value: "another test",
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "string-flag",
					FlagType: openfeature.String,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.TargetingMatchReason,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
		"non-existing string flag": {
			flag:         "string-flag",
			defaultValue: "test",
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "string-flag",
					DistinctId: "12345",
				},
				res: "false",
			},
			err: true,
			res: openfeature.StringEvaluationDetails{
				Value: "test",
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "string-flag",
					FlagType: openfeature.String,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.ErrorReason,
						ErrorCode:    openfeature.FlagNotFoundCode,
						ErrorMessage: `"string-flag" not found`,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
	}

	mockClient := &mockPostHogClient{t: t}
	p := NewProvider(mockClient)
	require.NoError(t, openfeature.SetProvider(p))
	c := openfeature.NewClient("testing")

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockClient.settings = tc.mockSettings
			res, err := c.StringValueDetails(context.Background(), tc.flag, tc.defaultValue, tc.evalCtx)
			assert.Equal(t, tc.err, err != nil)
			assert.Equal(t, tc.res, res)
		})
	}
}

func TestProvider_IntEvaluation(t *testing.T) {
	tcs := map[string]struct {
		flag         string
		defaultValue int64
		evalCtx      openfeature.EvaluationContext
		mockSettings mockSettings
		err          bool
		res          openfeature.IntEvaluationDetails
	}{
		"existing int flag": {
			flag:         "int-flag",
			defaultValue: 10,
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "int-flag",
					DistinctId: "12345",
				},
				res: "20",
			},
			res: openfeature.IntEvaluationDetails{
				Value: 20,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "int-flag",
					FlagType: openfeature.Int,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.TargetingMatchReason,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
		"non-existing int flag": {
			flag:         "int-flag",
			defaultValue: 10,
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "int-flag",
					DistinctId: "12345",
				},
				res: "false",
			},
			err: true,
			res: openfeature.IntEvaluationDetails{
				Value: 10,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "int-flag",
					FlagType: openfeature.Int,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.ErrorReason,
						ErrorCode:    openfeature.FlagNotFoundCode,
						ErrorMessage: `"int-flag" not found`,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
		"invalid format int flag": {
			flag:         "int-flag",
			defaultValue: 10,
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "int-flag",
					DistinctId: "12345",
				},
				res: "invalid int",
			},
			err: true,
			res: openfeature.IntEvaluationDetails{
				Value: 10,
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "int-flag",
					FlagType: openfeature.Int,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.ErrorReason,
						ErrorCode:    openfeature.TypeMismatchCode,
						ErrorMessage: `"invalid int" is not a int`,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
	}

	mockClient := &mockPostHogClient{t: t}
	p := NewProvider(mockClient)
	require.NoError(t, openfeature.SetProvider(p))
	c := openfeature.NewClient("testing")

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockClient.settings = tc.mockSettings
			res, err := c.IntValueDetails(context.Background(), tc.flag, tc.defaultValue, tc.evalCtx)
			assert.Equal(t, tc.err, err != nil)
			assert.Equal(t, tc.res, res)
		})
	}
}

func TestProvider_ObjectEvaluation(t *testing.T) {
	tcs := map[string]struct {
		flag         string
		defaultValue interface{}
		evalCtx      openfeature.EvaluationContext
		mockSettings mockSettings
		err          bool
		res          openfeature.InterfaceEvaluationDetails
	}{
		"existing object flag": {
			flag:         "object-flag",
			defaultValue: map[string]interface{}{"name": "john doe", "age": 35},
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "object-flag",
					DistinctId: "12345",
				},
				res: `{"name": "jane doe", "age": 52.5}`,
			},
			res: openfeature.InterfaceEvaluationDetails{
				Value: map[string]interface{}{"name": "jane doe", "age": 52.5},
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "object-flag",
					FlagType: openfeature.Object,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.TargetingMatchReason,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
		"non-existing object flag": {
			flag:         "object-flag",
			defaultValue: map[string]interface{}{"name": "john doe", "age": 35},
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "object-flag",
					DistinctId: "12345",
				},
				res: "false",
			},
			err: true,
			res: openfeature.InterfaceEvaluationDetails{
				Value: map[string]interface{}{"name": "john doe", "age": 35},
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "object-flag",
					FlagType: openfeature.Object,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.ErrorReason,
						ErrorCode:    openfeature.FlagNotFoundCode,
						ErrorMessage: `"object-flag" not found`,
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
		"invalid format object flag": {
			flag:         "object-flag",
			defaultValue: map[string]interface{}{"name": "john doe", "age": 35},
			evalCtx:      openfeature.NewEvaluationContext("12345", map[string]interface{}{}),
			mockSettings: mockSettings{
				payload: posthog.FeatureFlagPayload{
					Key:        "object-flag",
					DistinctId: "12345",
				},
				res: "{invalid json",
			},
			err: true,
			res: openfeature.InterfaceEvaluationDetails{
				Value: map[string]interface{}{"name": "john doe", "age": 35},
				EvaluationDetails: openfeature.EvaluationDetails{
					FlagKey:  "object-flag",
					FlagType: openfeature.Object,
					ResolutionDetail: openfeature.ResolutionDetail{
						Reason:       openfeature.ErrorReason,
						ErrorCode:    openfeature.TypeMismatchCode,
						ErrorMessage: "invalid JSON as flag value",
						FlagMetadata: openfeature.FlagMetadata{},
					},
				},
			},
		},
	}

	mockClient := &mockPostHogClient{t: t}
	p := NewProvider(mockClient)
	require.NoError(t, openfeature.SetProvider(p))
	c := openfeature.NewClient("testing")

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockClient.settings = tc.mockSettings
			res, err := c.ObjectValueDetails(context.Background(), tc.flag, tc.defaultValue, tc.evalCtx)
			assert.Equal(t, tc.err, err != nil)
			assert.Equal(t, tc.res, res)
		})
	}
}

type mockPostHogClient struct {
	posthog.Client

	t        *testing.T
	settings mockSettings
}

type mockSettings struct {
	payload posthog.FeatureFlagPayload
	res     interface{}
}

func (m *mockPostHogClient) GetFeatureFlag(payload posthog.FeatureFlagPayload) (interface{}, error) {
	assert.Equal(m.t, m.settings.payload, payload)
	return m.settings.res, nil
}
