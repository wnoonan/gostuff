package jobbr

import (
	"fmt"
	"strconv"

	"github.com/pterm/pterm"
)

type Outputter interface {
	// RootOpts displays the root options
	RootOpts() error
	// QueueJobs displays selectable jobs
	QueueJobs() error
	// ShowQueue displays all jobs in the queue
	ShowQueue() error
	// Show Jobs displays all jobs
	ShowJobs() error
	// ShowJob displays a single job and its details
	ShowJob(job *Job) error
	// CreateJob creates a new job
	CreateJob() error
}

type out struct {
	jobber Jobber
}

// NewOutput creates a new outputter
func NewOutput(jobber Jobber) Outputter {
	return &out{
		jobber: jobber,
	}
}

func (d *out) RootOpts() error {
	var opts []string

	optActions := map[string]func() error{
		"Queue Jobs": d.QueueJobs,
		"Show Queue": d.ShowQueue,
		"Create Job": d.CreateJob,
	}

	for opt := range optActions {
		opts = append(opts, opt)
	}

	selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(opts).Show()
	if err != nil {
		return fmt.Errorf("failed to select option: %w", err)
	}

	return optActions[selectedOption]()
}

func (d *out) QueueJobs() error {
	var opts []string

	jobs, err := d.jobber.Jobs()
	if err != nil {
		return fmt.Errorf("failed to get jobs: %w", err)
	}

	for _, job := range *jobs {
		if !job.Executed {
			opts = append(opts, job.Name)
		}
	}

	printer := pterm.DefaultInteractiveMultiselect.
		WithOptions(opts).
		WithFilter(false).
		WithCheckmark(&pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")})

	selectedOptions, err := printer.Show()
	if err != nil {
		return fmt.Errorf("failed to select jobs: %w", err)
	}

	for _, opt := range selectedOptions {
		for _, job := range *jobs {
			if job.Name == opt {
				d.jobber.Enqueue(&job)
			}
		}
	}

	return d.ShowQueue()
}

func (d *out) ShowQueue() error {
	queue, err := d.jobber.Queue()
	if err != nil {
		return fmt.Errorf("failed to get queue: %w", err)
	}

	jobs := queue.Jobs
	multi := pterm.DefaultMultiPrinter

	for _, job := range jobs {
		_, err := pterm.DefaultSpinner.WithWriter(multi.NewWriter()).Start(job.Name)
		if err != nil {
			return fmt.Errorf("failed to start spinner while viewing queue: %w", err)
		}
	}

	multi.Stop()

	return d.RootOpts()
}

func (d *out) ShowJobs() error {
	jobs, err := d.jobber.Jobs()
	if err != nil {
		return fmt.Errorf("failed to get jobs: %w", err)
	}

	var opts []string
	optJobs := map[string]*Job{}

	for _, job := range *jobs {
		opts = append(opts, job.Name)
		optJobs[job.Name] = &job
	}

	selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(opts).Show()
	if err != nil {
		return fmt.Errorf("failed to select option: %w", err)
	}

	return d.ShowJob(optJobs[selectedOption])
}

func (d *out) ShowJob(job *Job) error {
	tableData := pterm.TableData{
		{"ID", job.ID},
		{"Name", job.Name},
		{"Success", strconv.FormatBool(job.Success)},
		{"Command", job.Command},
		{"Output", job.Output},
		{"Error", job.Error},
		{"Created", job.Created.String()},
		{"Completed", job.Completed.String()},
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return d.RootOpts()
}

func (d *out) CreateJob() error {
	jobNameInput := pterm.DefaultInteractiveTextInput.WithDefaultText("Enter Job Name")
	jobName, err := jobNameInput.Show()
	if err != nil {
		return fmt.Errorf("failed to get job name: %w", err)
	}

	pterm.Println()
	pterm.Info.Printfln("Job Name: %s", jobName)

	jobCommandInput := pterm.DefaultInteractiveTextInput.
		WithDefaultText("Enter Command to Execute").
		WithDefaultValue("`echo 'Hello, World!'`").
		WithMultiLine()

	cmd, err := jobCommandInput.Show()
	if err != nil {
		return fmt.Errorf("failed to get job command: %w", err)
	}

	pterm.Println()
	pterm.Info.Printfln("Job Command: %s", cmd)

	err = d.jobber.CreateJob(
		&Job{
			Name:    jobName,
			Command: cmd,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return d.RootOpts()
}
