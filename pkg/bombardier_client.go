package pkg

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"os/exec"
	"strings"
	//"github.com/golang/protobuf/proto"

	"github.com/yiqinguo/armyant/pkg/models"
)

const (
	bombardierPath string = "./bombardier"
)

type bombardierClient struct {
	Exec *exec.Cmd
}

func NewBombardierClient(job models.Job) *bombardierClient {
	var args []string
	args = append(args, job.Url)
	args = append(args, fmt.Sprintf("-c%d", job.Connections))
	if job.Requests == 0 {
		args = append(args, fmt.Sprintf("-d%s", job.Duration))
	} else {
		args = append(args, fmt.Sprintf("-n%d", job.Requests))
	}
	args = append(args, fmt.Sprint("-oj"))
	args = append(args, fmt.Sprint("-pr"))
	args = append(args, fmt.Sprint("--fasthttp"))
	args = append(args, fmt.Sprint("-l"))
	command := exec.Command(bombardierPath, args...)
	fmt.Println(strings.Join(args, " "))

	return &bombardierClient{
		Exec: command,
	}
}

func (b *bombardierClient) Run() (models.Stats, error) {
	var stats models.Stats
	result, err := b.Exec.Output()
	if err != nil {
		return stats, err
	}

	err = json.Unmarshal(result, &stats)
	if err != nil {
		return stats, err
	}
	return stats, nil
}
