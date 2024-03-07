package workerid

import (
	"fmt"
	"sync"
)

type Client struct {
	workerIDGen *IDGen
}

var client *Client
var initOnce sync.Once

const DefaultWorkerID = 1

func GetWorkerID() (int64, error) {
	initOnce.Do(func() {
		client = &Client{
			workerIDGen: NewWorkerIDGen(),
		}
	})
	return client.GetWorkerID()
}

func (c *Client) GetWorkerID() (int64, error) {
	workerID, err := c.genWorkerID()
	if err != nil {
		fmt.Println("getWorkerID err:", err)
		return DefaultWorkerID, nil
	}
	return workerID, nil
}

func (c *Client) genWorkerID() (int64, error) {
	return c.workerIDGen.GenWorkerID()
}
