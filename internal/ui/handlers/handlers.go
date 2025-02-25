package handlers

import (
	"context"
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/aws"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
	"github.com/HenryOwenz/cloudgate/internal/ui/navigation"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// ModelWrapper wraps a core.Model to implement the tea.Model interface
type ModelWrapper struct {
	Model *core.Model
}

// Update implements the tea.Model interface
func (m ModelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// This is just a placeholder - the actual update logic will be in the UI package
	return m, nil
}

// View implements the tea.Model interface
func (m ModelWrapper) View() string {
	// This is just a placeholder - the actual view logic will be in the UI package
	return ""
}

// Init implements the tea.Model interface
func (m ModelWrapper) Init() tea.Cmd {
	// This is just a placeholder - the actual init logic will be in the UI package
	return nil
}

// WrapModel wraps a core.Model in a ModelWrapper
func WrapModel(m *core.Model) ModelWrapper {
	return ModelWrapper{Model: m}
}

// HandleEnter processes the enter key press based on the current view
func HandleEnter(m *core.Model) (tea.Model, tea.Cmd) {
	// Special handling for manual input in AWS config view
	if m.CurrentView == constants.ViewAWSConfig && m.ManualInput {
		newModel := *m

		// Get the entered value
		value := strings.TrimSpace(m.TextInput.Value())
		if value == "" {
			// If empty, just exit manual input mode
			newModel.ManualInput = false
			newModel.ResetTextInput()
			view.UpdateTableForView(&newModel)
			return WrapModel(&newModel), nil
		}

		// Set the appropriate value based on context
		if m.AwsProfile == "" {
			// Setting profile
			newModel.AwsProfile = value
			newModel.ManualInput = false
			newModel.ResetTextInput()
			view.UpdateTableForView(&newModel)
		} else {
			// Setting region and moving to next view
			newModel.AwsRegion = value
			newModel.ManualInput = false
			newModel.ResetTextInput()
			newModel.CurrentView = constants.ViewSelectService
			view.UpdateTableForView(&newModel)
		}

		return WrapModel(&newModel), nil
	}

	// Regular view handling
	switch m.CurrentView {
	case constants.ViewProviders:
		return HandleProviderSelection(m)
	case constants.ViewAWSConfig:
		return HandleAWSConfigSelection(m)
	case constants.ViewSelectService:
		return HandleServiceSelection(m)
	case constants.ViewSelectCategory:
		return HandleCategorySelection(m)
	case constants.ViewSelectOperation:
		return HandleOperationSelection(m)
	case constants.ViewApprovals:
		return HandleApprovalSelection(m)
	case constants.ViewConfirmation:
		return HandleConfirmationSelection(m)
	case constants.ViewSummary:
		if !m.ManualInput {
			if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
				if selected := m.Table.SelectedRow(); len(selected) > 0 {
					newModel := *m
					switch selected[0] {
					case "Latest Commit":
						newModel.CurrentView = constants.ViewExecutingAction
						newModel.Summary = "" // Empty string means use latest commit
						view.UpdateTableForView(&newModel)
						return WrapModel(&newModel), nil
					case "Manual Input":
						newModel.ManualInput = true
						newModel.TextInput.Focus()
						newModel.TextInput.Placeholder = "Enter commit ID..."
						return WrapModel(&newModel), nil
					}
				}
			}
		}
		return HandleSummaryConfirmation(m)
	case constants.ViewExecutingAction:
		return HandleExecutionSelection(m)
	case constants.ViewPipelineStatus:
		if selected := m.Table.SelectedRow(); len(selected) > 0 {
			newModel := *m
			for _, pipeline := range m.Pipelines {
				if pipeline.Name == selected[0] {
					if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
						newModel.CurrentView = constants.ViewExecutingAction
						newModel.SelectedPipeline = &pipeline
						view.UpdateTableForView(&newModel)
						return WrapModel(&newModel), nil
					}
					newModel.CurrentView = constants.ViewPipelineStages
					newModel.SelectedPipeline = &pipeline
					view.UpdateTableForView(&newModel)
					return WrapModel(&newModel), nil
				}
			}
		}
	case constants.ViewPipelineStages:
		// Just view only, no action
	}
	return WrapModel(m), nil
}

// HandleProviderSelection handles the selection of a cloud provider
func HandleProviderSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		if selected[0] == "Amazon Web Services" {
			newModel := *m
			newModel.CurrentView = constants.ViewAWSConfig
			view.UpdateTableForView(&newModel)
			return WrapModel(&newModel), nil
		}
	}
	return WrapModel(m), nil
}

// HandleAWSConfigSelection handles the selection of AWS profile or region
func HandleAWSConfigSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := *m

		// Handle "Manual Entry" option
		if selected[0] == "Manual Entry" {
			newModel.ManualInput = true
			newModel.TextInput.Focus()

			// Set appropriate placeholder based on context
			if m.AwsProfile == "" {
				newModel.TextInput.Placeholder = "Enter AWS profile name..."
			} else {
				newModel.TextInput.Placeholder = "Enter AWS region..."
			}

			return WrapModel(&newModel), nil
		}

		// Handle regular selection
		if m.AwsProfile == "" {
			newModel.AwsProfile = selected[0]
			view.UpdateTableForView(&newModel)
		} else {
			newModel.AwsRegion = selected[0]
			newModel.CurrentView = constants.ViewSelectService
			view.UpdateTableForView(&newModel)
		}
		return WrapModel(&newModel), nil
	}
	return WrapModel(m), nil
}

// HandleServiceSelection handles the selection of an AWS service
func HandleServiceSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		newModel.SelectedService = &core.Service{
			Name:        selected[0],
			Description: selected[1],
		}
		newModel.CurrentView = constants.ViewSelectCategory
		view.UpdateTableForView(&newModel)
		return WrapModel(&newModel), nil
	}
	return WrapModel(m), nil
}

// HandleCategorySelection handles the selection of a service category
func HandleCategorySelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		newModel.SelectedCategory = &core.Category{
			Name:        selected[0],
			Description: selected[1],
		}
		newModel.CurrentView = constants.ViewSelectOperation
		view.UpdateTableForView(&newModel)
		return WrapModel(&newModel), nil
	}
	return WrapModel(m), nil
}

// HandleOperationSelection handles the selection of a service operation
func HandleOperationSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		newModel.SelectedOperation = &core.Operation{
			Name:        selected[0],
			Description: selected[1],
		}

		if selected[0] == "Pipeline Approvals" {
			// Start loading approvals
			newModel.IsLoading = true
			newModel.LoadingMsg = "Loading approvals..."
			return WrapModel(&newModel), FetchApprovals(m)
		} else if selected[0] == "Pipeline Status" || selected[0] == "Start Pipeline" {
			// Start loading pipeline status
			newModel.IsLoading = true
			newModel.LoadingMsg = "Loading pipelines..."
			return WrapModel(&newModel), FetchPipelineStatus(m)
		}
	}
	return WrapModel(m), nil
}

// HandleApprovalSelection handles the selection of a pipeline approval
func HandleApprovalSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		for _, approval := range m.Approvals {
			if approval.PipelineName == selected[0] &&
				approval.StageName == selected[1] &&
				approval.ActionName == selected[2] {
				newModel.SelectedApproval = &approval
				newModel.CurrentView = constants.ViewConfirmation
				view.UpdateTableForView(&newModel)
				return WrapModel(&newModel), nil
			}
		}
	}
	return WrapModel(m), nil
}

// HandleConfirmationSelection handles the confirmation of an action
func HandleConfirmationSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		if selected[0] == "Approve" {
			newModel.ApproveAction = true
			newModel.CurrentView = constants.ViewSummary
			newModel.ManualInput = true // Set manual input mode directly
			newModel.SetTextInputForApproval(true)
			newModel.ApprovalComment = "" // Reset any previous comment
		} else if selected[0] == "Reject" {
			newModel.ApproveAction = false
			newModel.CurrentView = constants.ViewSummary
			newModel.ManualInput = true // Set manual input mode directly
			newModel.SetTextInputForApproval(false)
			newModel.ApprovalComment = "" // Reset any previous comment
		}
		view.UpdateTableForView(&newModel)
		return WrapModel(&newModel), nil
	}
	return WrapModel(m), nil
}

// HandleSummaryConfirmation handles the confirmation of a summary
func HandleSummaryConfirmation(m *core.Model) (tea.Model, tea.Cmd) {
	if m.SelectedApproval != nil {
		// For approval actions, check if we have a comment
		if strings.TrimSpace(m.TextInput.Value()) == "" {
			// No comment yet, keep the text input focused
			newModel := *m
			newModel.TextInput.Focus()
			return WrapModel(&newModel), nil
		}

		// We have a comment, store it and proceed to execution confirmation screen
		newModel := *m
		newModel.ApprovalComment = m.TextInput.Value() // Explicitly store the comment
		newModel.CurrentView = constants.ViewExecutingAction
		newModel.ManualInput = false
		view.UpdateTableForView(&newModel)
		return WrapModel(&newModel), nil
	} else if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
		if m.ManualInput && strings.TrimSpace(m.TextInput.Value()) == "" {
			// No commit ID yet, keep the text input focused
			newModel := *m
			newModel.TextInput.Focus()
			return WrapModel(&newModel), nil
		}

		// We have a commit ID or using latest, store it and proceed to execution confirmation screen
		newModel := *m
		if m.ManualInput {
			newModel.CommitID = m.TextInput.Value() // Explicitly store the commit ID
			newModel.ManualCommitID = true
		}
		newModel.CurrentView = constants.ViewExecutingAction
		newModel.ManualInput = false
		view.UpdateTableForView(&newModel)
		return WrapModel(&newModel), nil
	}

	return WrapModel(m), nil
}

// HandleExecutionSelection handles the selection of an execution action
func HandleExecutionSelection(m *core.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		if selected[0] == "Execute" {
			newModel := *m
			newModel.IsLoading = true
			if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
				newModel.LoadingMsg = "Starting pipeline..."
				return WrapModel(&newModel), ExecutePipeline(m)
			} else if m.SelectedApproval != nil {
				newModel.LoadingMsg = "Executing approval action..."
				return WrapModel(&newModel), ExecuteApproval(m)
			}
		} else if selected[0] == "Cancel" {
			return WrapModel(navigation.NavigateBack(m)), nil
		}
	}
	return WrapModel(m), nil
}

// Async operations

// FetchApprovals fetches pipeline approvals
func FetchApprovals(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		provider, err := aws.New(ctx, m.AwsProfile, m.AwsRegion)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		approvals, err := provider.GetPendingApprovals(ctx)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		return core.ApprovalsMsg{
			Provider:  provider,
			Approvals: approvals,
		}
	}
}

// FetchPipelineStatus fetches pipeline status
func FetchPipelineStatus(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		provider, err := aws.New(ctx, m.AwsProfile, m.AwsRegion)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		pipelines, err := provider.GetPipelineStatus(ctx)
		if err != nil {
			return core.ErrMsg{Err: err}
		}

		return core.PipelineStatusMsg{
			Provider:  provider,
			Pipelines: pipelines,
		}
	}
}

// ExecuteApproval executes an approval action
func ExecuteApproval(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		provider, err := aws.New(ctx, m.AwsProfile, m.AwsRegion)
		if err != nil {
			return core.ApprovalResultMsg{Err: err}
		}

		err = provider.PutApprovalResult(ctx, *m.SelectedApproval, m.ApproveAction, m.ApprovalComment)
		return core.ApprovalResultMsg{Err: err}
	}
}

// ExecutePipeline executes a pipeline
func ExecutePipeline(m *core.Model) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		provider, err := aws.New(ctx, m.AwsProfile, m.AwsRegion)
		if err != nil {
			return core.PipelineExecutionMsg{Err: err}
		}

		commitID := ""
		if m.ManualCommitID {
			commitID = m.CommitID
		}

		err = provider.StartPipelineExecution(ctx, m.SelectedPipeline.Name, commitID)
		return core.PipelineExecutionMsg{Err: err}
	}
}
