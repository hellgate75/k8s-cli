package executor

import (
	"errors"
	"fmt"
	"github.com/hellgate75/k8s-cli/common"
	"github.com/hellgate75/k8s-cli/io"
	"github.com/hellgate75/k8s-cli/model"
	"os"
	"strings"
	"sync"
)

type Executor struct {
	sync.Mutex
	baseFolder string
	request    model.CommandRequest
	internal   *model.ClusterData
}

func (c Executor) commit() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		c.Unlock()

	}()
	c.Lock()
	file := fmt.Sprintf("%s%cclusters.yaml", c.baseFolder, os.PathSeparator)
	var data = make([]byte, 0)
	data, err = io.ToYaml(c.internal)
	if err != nil {
		return err
	}
	err = io.WriteFile(file, data, true)
	if err != nil {
		return err
	}
	return err
}

func (c Executor) load() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		c.Unlock()

	}()
	c.Lock()
	file := fmt.Sprintf("%s%cclusters.yaml", c.baseFolder, os.PathSeparator)
	data, err := io.ReadFile(file)
	if err != nil {
		return err
	} else if len(data) == 0 {
		return errors.New(fmt.Sprintf("Empty file: %s", file))
	}
	err = io.FromYaml(data, &c.internal)
	if err != nil {
		return err
	}
	return nil
}
func (c Executor) Init() error {
	if _, err := os.Stat(c.baseFolder); err != nil {
		err = os.MkdirAll(c.baseFolder, 0660)
		if err != nil {
			return err
		}
	}
	file := fmt.Sprintf("%s%cclusters.yaml", c.baseFolder, os.PathSeparator)
	data, err := io.ReadFile(file)
	if err != nil || len(data) == 0 {
		err := c.commit()
		if err != nil {
			return err
		}
	} else {
		err = io.FromYaml(data, c.internal)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Executor) Execute() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Executor.Execute Error: %v", r))
		}
	}()
	switch strings.ToLower(c.request.Command) {
	case "show":
		err = c.show()
	case "details":
		err = c.details()
	case "add":
		err = c.add()
	case "remove":
		err = c.remove()
	case "verify":
		err = c.verify()
	case "prepare":
		err = c.prepare()
	case "ensure":
		err = c.ensure()
	default:
		err = errors.New(fmt.Sprintf("Unknown executor: <%s>", c.request.Command))
	}
	return err
}

func (c Executor) print(i interface{}) {
	var dt []byte
	if common.FixOutputType(c.request.Format) == "yaml" {
		dt, _ = io.ToYaml(i)
	} else {
		dt, _ = io.ToJson(i)
	}
	fmt.Printf("%s\n", string(dt))
}

func New(baseFolder string, request model.CommandRequest) Executor {
	return Executor{
		baseFolder: baseFolder,
		request:    request,
		internal: &model.ClusterData{
			Clusters: make([]model.Cluster, 0),
		},
	}
}
