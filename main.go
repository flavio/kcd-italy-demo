package main

import (
	demo "github.com/saschagrunert/demo"
	"os/exec"
)

const DEMO_POLICY = "registry://ghcr.io/kubewarden/policies/safe-labels:v0.1.6"

func main() {
	d := demo.New()
	d.Add(kwctlRun(), "kwctl demo", "kwctl demo")
	d.Run()
}

func kwctlRun() *demo.Run {
	r := demo.NewRun(
		"Running policies with kwctl",
	)

	r.Setup(cleanupKwctl)
	r.Cleanup(cleanupKwctl)

	kwctl(r)

	return r
}

func kwctl(r *demo.Run) {
	r.Step(demo.S(
		"List policies",
	), demo.S("kwctl policies"))

	r.Step(demo.S(
		"Pull a policy",
	), demo.S("kwctl pull", DEMO_POLICY))

	r.Step(demo.S(
		"List policies",
	), demo.S("kwctl policies"))

	r.Step(demo.S(
		"Inspect policy",
	), demo.S("kwctl inspect", DEMO_POLICY))

	r.Step(demo.S(
		"Demo request",
	), demo.S("bat test_data/ingress.json"))

	r.Step(demo.S(
		"Evaluate request - put constraint on 'cc-center'",
	), demo.S("kwctl -v run",
		`--settings-json '{
			"mandatory_labels": ["cc-center"],
			"constrained_labels": {
				"cc-center": "^cc-\\d+$"
			}}'`,
		"--request-path test_data/ingress.json",
		DEMO_POLICY,
		"|",
		"jq"))

	r.StepCanFail(demo.S(
		"Evaluate request - put constraint on 'cc-center' and 'owner'",
	), demo.S("kwctl -v run",
		`--settings-json '{
			"mandatory_labels": ["cc-center", "owner"],
			"constrained_labels": {
				"cc-center": "^cc-\\d+$"
			}}'`,
		"--request-path test_data/ingress.json",
		DEMO_POLICY,
		"|",
		"jq"))

	r.Step(demo.S(
		"Scaffold ClusterAdmissionPolicy Resource",
	), demo.S("kwctl manifest",
		`--settings-json '{
			"mandatory_labels": ["cc-center", "owner"],
			"constrained_labels": {
				"cc-center": "^cc-\\d+$"
			}}'`,
		"--type ClusterAdmissionPolicy",
		DEMO_POLICY,
		"|",
		"bat -l yaml"))
}

func cleanupKwctl() error {
	exec.Command("kwctl", "rm", DEMO_POLICY).Run()
	return nil
}
