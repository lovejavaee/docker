package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"text/tabwriter"

	"github.com/docker/docker/api/types"
	flag "github.com/docker/docker/pkg/mflag"
)

// CmdTop displays the running processes of a container.
//
// Usage: docker top CONTAINER
func (cli *DockerCli) CmdTop(args ...string) error {
	cmd := cli.Subcmd("top", "CONTAINER [ps OPTIONS]", "Display the running processes of a container", true)
	cmd.Require(flag.Min, 1)

	cmd.ParseFlags(args, true)

	val := url.Values{}
	if cmd.NArg() > 1 {
		val.Set("ps_args", strings.Join(cmd.Args()[1:], " "))
	}

	stream, _, err := cli.call("GET", "/containers/"+cmd.Arg(0)+"/top?"+val.Encode(), nil, nil)
	if err != nil {
		return err
	}

	procList := types.ContainerProcessList{}
	if err := json.NewDecoder(stream).Decode(&procList); err != nil {
		return err
	}

	w := tabwriter.NewWriter(cli.out, 20, 1, 3, ' ', 0)
	fmt.Fprintln(w, strings.Join(procList.Titles, "\t"))

	for _, proc := range procList.Processes {
		fmt.Fprintln(w, strings.Join(proc, "\t"))
	}
	w.Flush()
	return nil
}
