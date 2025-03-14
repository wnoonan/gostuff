package jobbr

import (
	"fmt"
	"sort"
	"strconv"
	"time"

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
	// Delete Job deletes a job
	DeleteJob() error
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
		"1. Queue Jobs":  d.QueueJobs,
		"2. Show Queue":  d.ShowQueue,
		"3. Clear Queue": d.ClearQueue,
		"4. Dequeue Job": d.DequeueJob,
		"5. Create Job":  d.CreateJob,
		"6. Show Jobs":   d.ShowJobs,
		"7. Delete Job":  d.DeleteJob,
	}

	for opt := range optActions {
		opts = append(opts, opt)
	}

	sort.Strings(opts)

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

	if len(jobs) == 0 {
		pterm.Info.Println("No jobs available")
		return d.RootOpts()
	}

	for _, job := range jobs {
		if !job.Executed {
			opts = append(opts, job.Name)
		}
	}

	if len(opts) == 0 {
		pterm.Info.Println("No jobs available to queue")
		return d.RootOpts()
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
		for _, job := range jobs {
			if job.Name == opt {
				d.jobber.Enqueue(job)
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

	if len(queue) == 0 {
		pterm.Info.Println("No jobs in queue")
		return d.RootOpts()
	}

	jobsAndSpinners := map[string]*pterm.SpinnerPrinter{}

	multi := pterm.DefaultMultiPrinter

	for _, job := range queue {
		spinner, err := pterm.DefaultSpinner.WithWriter(multi.NewWriter()).Start(job.Name)
		if err != nil {
			return fmt.Errorf("failed to start spinner while viewing queue: %w", err)
		}
		jobsAndSpinners[job.ID] = spinner
		if job.Executed && job.Success {
			jobsAndSpinners[job.ID].Success(job.Name)
		} else if job.Executed && !job.Success {
			jobsAndSpinners[job.ID].Fail(job.Name)
		}
	}

	multi.Start()

	for {
		time.Sleep(1 * time.Second)
		queue, err := d.jobber.Queue()
		if err != nil || len(queue) == 0 {
			break
		}

		for _, job := range queue {
			if job.Executed && job.Success {
				jobsAndSpinners[job.ID].Success(job.Name)
			} else if job.Executed && !job.Success {
				jobsAndSpinners[job.ID].Fail(job.Name)
			}
		}

		allExecuted := true
		for _, job := range queue {
			if !job.Executed {
				allExecuted = false
			}
		}

		if allExecuted {
			break
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

	if len(jobs) == 0 {
		pterm.Info.Println("No jobs available")
		return d.RootOpts()
	}

	var opts []string
	optJobs := map[string]*Job{}

	for _, job := range jobs {
		opts = append(opts, job.Name)
		optJobs[job.Name] = job
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
		{"Executed", strconv.FormatBool(job.Executed)},
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

	pterm.Info.Printf("Job %s created", jobName)

	return d.RootOpts()
}

func (d *out) ClearQueue() error {
	err := d.jobber.ClearQueue()
	if err != nil {
		return fmt.Errorf("failed to clear queue: %w", err)
	}

	pterm.Info.Println("Queue cleared")

	return d.RootOpts()
}

func (d *out) DequeueJob() error {
	jobs, err := d.jobber.Queue()
	if err != nil {
		return fmt.Errorf("failed to get jobs: %w", err)
	}

	if len(jobs) == 0 {
		pterm.Info.Println("No jobs in queue")
		return d.RootOpts()
	}

	var opts []string
	optJobs := map[string]*Job{}

	for _, job := range jobs {
		opts = append(opts, job.Name)
		optJobs[job.Name] = job
	}

	selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(opts).Show()
	if err != nil {
		return fmt.Errorf("failed to select option: %w", err)
	}

	_, err = d.jobber.Dequeue(optJobs[selectedOption])
	if err != nil {
		return fmt.Errorf("failed to dequeue job: %w", err)
	}

	pterm.Info.Printf("Job %s dequeued", selectedOption)

	return d.RootOpts()
}

func (d *out) DeleteJob() error {
	jobs, err := d.jobber.Jobs()
	if err != nil {
		return fmt.Errorf("failed to get jobs: %w", err)
	}

	if len(jobs) == 0 {
		pterm.Info.Println("No jobs available")
		return d.RootOpts()
	}

	var opts []string
	optJobs := map[string]*Job{}

	for _, job := range jobs {
		opts = append(opts, job.Name)
		optJobs[job.Name] = job
	}

	selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(opts).Show()
	if err != nil {
		return fmt.Errorf("failed to select option: %w", err)
	}

	err = d.jobber.DeleteJob(optJobs[selectedOption])
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	pterm.Info.Printf("Job %s deleted", selectedOption)

	return d.RootOpts()
}
