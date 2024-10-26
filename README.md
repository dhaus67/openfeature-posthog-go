# PostHog OpenFeature Provider for Go

This repository provides an implementation of an [OpenFeature provider](https://openfeature.dev/docs/reference/concepts/provider) for [PostHog](https://posthog.com/).

## Getting started

This simple snippet should be enough to get you started to use the provider:
```go
import (
 	"fmt"
  "context"

	"github.com/dhaus67/openfeature-posthog-go"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/posthog/posthog-go"
)

func main() {
	// Start by creating a PostHog client with your desired configuration.
	client, err := posthog.NewWithConifg("<your api key>", posthog.Config{})
	if err != nil {
		panic(err)
	}

	// Create the provider and register it.
	openfeature.SetProvider(provider.NewProvider(client))

	client := openfeature.NewClient("my-client")

	// The targeting key is required with this provider. It is used to evaluate with PostHog
	// whether for the specific user the feature is enabled or not.
	evalCtx := openfeature.NewEvaluationContext("<distinct-user-id>", map[string]interface{}{})

	secretFeature, err := client.BooleanValue(context.Background(), "secret", false, evalCtx)
	if err != nil {
		panic(err)
	}

	if secretFeature {
		fmt.Println("Secret feature is enabled")
	}
}
```

In addition to the targeting key, it is also possible to specify additional values to filter on
for the PostHog user: `groups`, `groupProperties`, and `personProperties`.

The documentation for [the PostHog Go SDK has a rich documentation about these use-cases](https://posthog.com/docs/libraries/go#advanced-overriding-server-properties).
