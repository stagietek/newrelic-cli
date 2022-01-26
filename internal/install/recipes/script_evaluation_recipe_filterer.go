package recipes

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type ScriptEvaluationRecipeFilterer struct {
	recipeExecutor execution.RecipeExecutor
	installStatus  *execution.InstallStatus
}

func NewScriptEvaluationRecipeFilterer(installStatus *execution.InstallStatus) *ScriptEvaluationRecipeFilterer {
	recipeExecutor := execution.NewShRecipeExecutor()

	return &ScriptEvaluationRecipeFilterer{
		recipeExecutor: recipeExecutor,
		installStatus:  installStatus,
	}
}

func (f *ScriptEvaluationRecipeFilterer) Filter(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) bool {
	if err := f.recipeExecutor.ExecutePreInstall(ctx, *r, types.RecipeVars{}); err != nil {
		log.Tracef("recipe %s failed script evaluation %s", r.Name, err)

		fmt.Printf("\nScriptEvaluationRecipeFilterer - Incoming:                %+v \n", err)

		var metadata map[string]interface{}
		if e, ok := err.(*types.CustomStdError); ok {
			fmt.Printf("\nScriptEvaluationRecipeFilterer - Metadata:   %+v \n", e.Metadata)
		}

		fmt.Print("\n **************************** \n\n")

		if utils.IsExitStatusCode(132, err) {
			event := execution.RecipeStatusEvent{
				Recipe:   *r,
				Metadata: metadata,
			}
			f.installStatus.RecipeDetected(*r, event)
		}

		return true
	}

	return false
}

func (f *ScriptEvaluationRecipeFilterer) CheckCompatibility(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) error {
	err := f.recipeExecutor.ExecutePreInstall(ctx, *r, types.RecipeVars{})

	if err != nil {
		var metadata map[string]interface{}
		var message string
		if e, ok := err.(*types.CustomStdError); ok {
			metadata = e.Metadata
		} else {
			message = err.Error()
		}

		fmt.Printf("\nScriptEvaluationRecipeFilterer::CheckCompatibility - Metadata:   %+v \n", metadata)

		if utils.IsExitStatusCode(132, err) {
			event := execution.RecipeStatusEvent{
				Recipe:   *r,
				Msg:      message,
				Metadata: metadata,
			}
			f.installStatus.RecipeDetected(*r, event)
		}
	}

	return err
}