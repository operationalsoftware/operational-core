package andonview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AddPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	AndonIssues      []model.AndonIssueNode
	Teams            []model.Team
	SelectedPath     []int
}

func AddPage(p *AddPageProps) g.Node {

	var selectedIssue int
	if len(p.SelectedPath) > 0 {
		selectedIssue = p.SelectedPath[len(p.SelectedPath)-1]
	} else {
		selectedIssue = 0
	}

	var selectedNode *model.AndonIssueNode
	for _, node := range p.AndonIssues {
		if node.AndonIssueID == selectedIssue {
			selectedNode = &node
			break
		}
	}

	content := g.Group([]g.Node{

		h.Div(
			h.Class("select-form"),

			h.Label(
				g.Text("Issue"),
			),

			cascadingSelects(p.AndonIssues, p.SelectedPath, p.Ctx.Req.URL),
		),

		addAndonForm(&addAndonFormProps{
			values:            p.Values,
			validationErrors:  p.ValidationErrors,
			isSubmission:      p.IsSubmission,
			andonIssues:       p.AndonIssues,
			teams:             p.Teams,
			selectedIssue:     selectedIssue,
			selectedIssueNode: selectedNode,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Add New Andon Issue",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			andonIssuesBreadCrumb,
			{IconIdentifier: "plus", Title: "Add"},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonview/add_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/andonview/add_page.js"),
		},
	})
}

type addAndonFormProps struct {
	values            url.Values
	validationErrors  validate.ValidationErrors
	isSubmission      bool
	andonIssues       []model.AndonIssueNode
	teams             []model.Team
	selectedIssue     int
	selectedIssueNode *model.AndonIssueNode
}

func addAndonForm(p *addAndonFormProps) g.Node {

	issueDescriptionLabel := "Issue Description"
	issueDescriptionKey := "IssueDescription"
	issueDescriptionValue := p.values.Get(issueDescriptionKey)
	issueDescriptionError := ""
	if p.isSubmission || issueDescriptionValue != "" {
		issueDescriptionError = p.validationErrors.GetError(issueDescriptionKey, issueDescriptionLabel)
	}
	issueDescriptionHelperType := components.InputHelperTypeNone
	if issueDescriptionError != "" {
		issueDescriptionHelperType = components.InputHelperTypeError
	}

	issueIDLabel := "Issue"
	issueIDKey := "IssueID"
	issueIDValue := p.values.Get(issueIDKey)
	issueIDError := ""
	if p.isSubmission || issueIDValue != "" {
		issueIDError = p.validationErrors.GetError(issueIDKey, issueIDLabel)
	}
	issueIDHelperType := components.InputHelperTypeNone
	if issueIDError != "" {
		issueIDHelperType = components.InputHelperTypeError
	}

	// map andon issues to options for parent select
	parentSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, andonIssue := range p.andonIssues {
		intVal, _ := strconv.Atoi(issueIDValue)
		isSelected := andonIssue.AndonIssueID == intVal

		parentSelectOptions = append(parentSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", andonIssue.AndonIssueID)),
			g.Attr("data-team", nilsafe.Str(andonIssue.AssignedTeamName)),
			g.If(isSelected, h.Selected()),
			g.Text(strings.Join(andonIssue.NamePath, " > ")),
		))
	}

	assignedTeamLabel := "Assigned Team"
	assignedTeamKey := "AssignedTeam"
	assignedTeamValue := p.values.Get(assignedTeamKey)
	assignedTeamError := ""
	if p.isSubmission || assignedTeamValue != "" {
		assignedTeamError = p.validationErrors.GetError(assignedTeamKey, assignedTeamLabel)
	}
	assignedTeamHelperType := components.InputHelperTypeNone
	if assignedTeamError != "" {
		assignedTeamHelperType = components.InputHelperTypeError
	}

	sourceKey := "Source"
	sourceValue := p.values.Get(sourceKey)

	locationLabel := "Location"
	locationKey := "Location"
	locationValue := p.values.Get(locationKey)
	locationError := ""
	if p.isSubmission || locationValue != "" {
		locationError = p.validationErrors.GetError(locationKey, locationLabel)
	}
	locationHelperType := components.InputHelperTypeNone
	if locationError != "" {
		locationHelperType = components.InputHelperTypeError
	}

	teamSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, team := range p.teams {
		intVal, _ := strconv.Atoi(assignedTeamValue)
		isSelected := team.TeamID == intVal

		teamSelectOptions = append(teamSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", team.TeamID)),
			g.If(isSelected, h.Selected()),
			g.Text(team.TeamName),
		))
	}

	var assignedTeam string
	if p.selectedIssueNode != nil {
		assignedTeam = nilsafe.Str(p.selectedIssueNode.AssignedTeamName)
	}

	return components.Form(
		h.ID("add-andon-issue-form"),
		h.Method("POST"),

		h.Div(
			h.Label(

				h.Input(
					h.ID("issue-select"),
					h.Name(issueIDKey),
					h.Type("hidden"),
					h.Value(strconv.Itoa(p.selectedIssue)),
				),
			),
			g.If(
				assignedTeam != "",
				h.Option(
					h.Class("assigned-team-label"),
					h.Value(assignedTeam),
					g.Text(fmt.Sprintf("* Assigned to: %s", assignedTeam)),
				),
			),
			g.If(issueIDError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: issueIDError,
					Type:  issueIDHelperType,
				})),
			g.If(assignedTeamError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: assignedTeamError,
					Type:  assignedTeamHelperType,
				})),
		),

		h.Div(
			h.Label(
				g.Text(issueDescriptionLabel),

				h.Textarea(
					h.Name(issueDescriptionKey),
					h.Placeholder("Enter issue description"),
					h.Value(issueDescriptionValue),
					h.AutoComplete("off"),
					g.Text(issueDescriptionValue),
				),
			),
			g.If(
				issueDescriptionError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: issueDescriptionError,
					Type:  issueDescriptionHelperType,
				})),
		),

		h.Div(
			h.Label(
				g.Text(locationLabel),

				h.Input(
					h.Name(locationKey),
					h.Placeholder("Enter location"),
					h.Value(locationValue),
					h.AutoComplete("off"),
				),
			),
			g.If(
				locationError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: locationError,
					Type:  locationHelperType,
				}),
			),
		),

		h.Input(
			h.Name(assignedTeamKey),
			h.Value(assignedTeam),
			h.Type("hidden"),
		),

		h.Input(
			h.Name(sourceKey),
			h.Value(sourceValue),
			h.Type("hidden"),
		),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Add Andon Event"),
		),
	)
}

func cascadingSelects(andonIssues []model.AndonIssueNode, selectedNodePath []int, baseURL *url.URL) g.Node {
	type NodeMap map[int][]model.AndonIssueNode

	// Build a map of parentID -> children
	childrenMap := make(NodeMap)
	for _, node := range andonIssues {
		if node.ParentID != nil {
			childrenMap[nilsafe.Int(node.ParentID)] = append(childrenMap[nilsafe.Int(node.ParentID)], node)
		} else {
			childrenMap[0] = append(childrenMap[0], node) // root level
		}
	}

	var selects []g.Node
	currentParentID := 0

	for level, selectedID := range selectedNodePath {
		selects = append(selects, issueSelect(level, childrenMap[currentParentID], selectedID, baseURL))
		currentParentID = selectedID
	}

	if nextLevelNodes, ok := childrenMap[currentParentID]; ok && len(nextLevelNodes) > 0 {
		selects = append(selects, issueSelect(len(selectedNodePath), nextLevelNodes, 0, baseURL))
	}

	return g.Group(selects)
}

func issueSelect(level int, nodes []model.AndonIssueNode, selectedID int, baseURL *url.URL) g.Node {
	newQuery := url.Values{}
	for key, vals := range baseURL.Query() {
		if strings.HasPrefix(key, "Node[") {
			var i int
			_, err := fmt.Sscanf(key, "Node[%d]", &i)
			if err == nil {
				if i >= level {
					continue
				}
			}
		}
		for _, v := range vals {
			newQuery.Add(key, v)
		}
	}

	name := fmt.Sprintf("Node[%d]", level)

	var hidden []g.Node
	for i := 0; i < level; i++ {
		k := fmt.Sprintf("Node[%d]", i)
		v := newQuery.Get(k)
		if v != "" {
			hidden = append(hidden, h.Input(
				h.Type("hidden"),
				h.Name(k),
				h.Value(v),
			))
		}
	}

	for key, vals := range newQuery {
		if key == name {
			continue
		}
		for _, v := range vals {
			hidden = append(hidden, h.Input(
				h.Type("hidden"),
				h.Name(key),
				h.Value(v),
			))
		}
	}

	return components.Form(
		h.Method("GET"),
		g.Group(hidden),

		h.Select(
			h.Name(name),
			g.Attr("onchange", "this.form.submit()"),
			h.Option(h.Value(""), g.Text("– Select –")),
			g.Group(func() []g.Node {
				var opts []g.Node
				for _, n := range nodes {

					val := fmt.Sprintf("%d", n.AndonIssueID)
					opt := h.Option(
						h.Value(val),
						g.Text(strings.Join(n.NamePath, " > ")),
					)
					if n.AndonIssueID == selectedID {
						opt = h.Option(
							h.Value(val),
							h.Selected(),
							g.Text(strings.Join(n.NamePath, " > ")),
						)
					}
					opts = append(opts, opt)
				}
				return opts
			}()),
		),
	)
}
