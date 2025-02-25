package view

import (
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
	"github.com/charmbracelet/bubbles/table"
)

// UpdateTableForView updates the table model based on the current view
func UpdateTableForView(m *core.Model) {
	columns := getColumnsForView(m)
	rows := getRowsForView(m)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(6),
	)

	t.SetStyles(m.Styles.Table)
	m.Table = t
}

// getColumnsForView returns the appropriate columns for the current view
func getColumnsForView(m *core.Model) []table.Column {
	switch m.CurrentView {
	case constants.ViewProviders:
		return []table.Column{
			{Title: "Provider", Width: 30},
			{Title: "Description", Width: 50},
		}
	case constants.ViewAWSConfig:
		if m.AwsProfile == "" {
			return []table.Column{{Title: "Profile", Width: 30}}
		}
		return []table.Column{{Title: "Region", Width: 30}}
	case constants.ViewSelectService:
		return []table.Column{
			{Title: "Service", Width: 30},
			{Title: "Description", Width: 50},
		}
	case constants.ViewSelectCategory:
		return []table.Column{
			{Title: "Category", Width: 30},
			{Title: "Description", Width: 50},
		}
	case constants.ViewSelectOperation:
		return []table.Column{
			{Title: "Operation", Width: 30},
			{Title: "Description", Width: 50},
		}
	case constants.ViewApprovals:
		return []table.Column{
			{Title: "Pipeline", Width: 40},
			{Title: "Stage", Width: 30},
			{Title: "Action", Width: 20},
		}
	case constants.ViewConfirmation:
		return []table.Column{
			{Title: "Action", Width: 30},
			{Title: "Description", Width: 50},
		}
	case constants.ViewExecutingAction:
		return []table.Column{
			{Title: "Action", Width: 30},
			{Title: "Description", Width: 50},
		}
	case constants.ViewPipelineStatus:
		return []table.Column{
			{Title: "Pipeline", Width: 40},
			{Title: "Description", Width: 50},
		}
	case constants.ViewPipelineStages:
		return []table.Column{
			{Title: "Stage", Width: 30},
			{Title: "Status", Width: 20},
			{Title: "Last Updated", Width: 20},
		}
	case constants.ViewSummary:
		return []table.Column{
			{Title: "Type", Width: 30},
			{Title: "Value", Width: 50},
		}
	default:
		return []table.Column{}
	}
}

// getRowsForView returns the appropriate rows for the current view
func getRowsForView(m *core.Model) []table.Row {
	switch m.CurrentView {
	case constants.ViewProviders:
		return []table.Row{
			{"Amazon Web Services", "AWS Cloud Services"},
			{"Microsoft Azure (Coming Soon)", "Azure Cloud Platform"},
			{"Google Cloud Platform (Coming Soon)", "Google Cloud Services"},
		}
	case constants.ViewAWSConfig:
		if m.AwsProfile == "" {
			rows := make([]table.Row, len(m.Profiles)+1)
			rows[0] = table.Row{"Manual Entry"}
			for i, profile := range m.Profiles {
				rows[i+1] = table.Row{profile}
			}
			return rows
		}
		rows := make([]table.Row, len(m.Regions)+1)
		rows[0] = table.Row{"Manual Entry"}
		for i, region := range m.Regions {
			rows[i+1] = table.Row{region}
		}
		return rows
	case constants.ViewSelectService:
		return []table.Row{
			{"CodePipeline", "Continuous Delivery Service"},
		}
	case constants.ViewSelectCategory:
		return []table.Row{
			{"Workflows", "Pipeline Workflows and Approvals"},
			{"Operations (Coming Soon)", "Service Operations"},
		}
	case constants.ViewSelectOperation:
		if m.SelectedCategory != nil && m.SelectedCategory.Name == "Workflows" {
			return []table.Row{
				{"Pipeline Approvals", "Manage Pipeline Approvals"},
				{"Pipeline Status", "View Pipeline Status"},
				{"Start Pipeline", "Trigger Pipeline Execution"},
			}
		}
		return []table.Row{}
	case constants.ViewApprovals:
		rows := make([]table.Row, len(m.Approvals))
		for i, approval := range m.Approvals {
			rows[i] = table.Row{
				approval.PipelineName,
				approval.StageName,
				approval.ActionName,
			}
		}
		return rows
	case constants.ViewConfirmation:
		return []table.Row{
			{"Approve", "Approve the pipeline stage"},
			{"Reject", "Reject the pipeline stage"},
		}
	case constants.ViewExecutingAction:
		if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
			return []table.Row{
				{"Execute", "Start pipeline with latest commit"},
				{"Cancel", "Cancel and return to main menu"},
			}
		}
		action := "approve"
		if !m.ApproveAction {
			action = "reject"
		}
		return []table.Row{
			{"Execute", fmt.Sprintf("Execute %s action", action)},
			{"Cancel", "Cancel and return to main menu"},
		}
	case constants.ViewPipelineStatus:
		if m.Pipelines == nil {
			return []table.Row{}
		}
		rows := make([]table.Row, len(m.Pipelines))
		for i, pipeline := range m.Pipelines {
			rows[i] = table.Row{
				pipeline.Name,
				fmt.Sprintf("%d stages", len(pipeline.Stages)),
			}
		}
		return rows
	case constants.ViewPipelineStages:
		if m.SelectedPipeline == nil {
			return []table.Row{}
		}
		rows := make([]table.Row, len(m.SelectedPipeline.Stages))
		for i, stage := range m.SelectedPipeline.Stages {
			rows[i] = table.Row{
				stage.Name,
				stage.Status,
				stage.LastUpdated,
			}
		}
		return rows
	case constants.ViewSummary:
		if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
			if m.SelectedPipeline == nil {
				return []table.Row{}
			}
			return []table.Row{
				{"Latest Commit", "Use latest commit from source"},
				{"Manual Input", "Enter specific commit ID"},
			}
		}
		// For approval summary, don't show any rows since we're showing text input
		return []table.Row{}
	default:
		return []table.Row{}
	}
}
