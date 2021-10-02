// Author: Daniel TAN
// Date: 2021-10-02 22:31:49
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-02 22:33:15
// FilePath: /trinity-micro/container.go
// Description:
package trinity

import "sync"

var (
	// booting instance
	_bootingInstances []bootingInstance
)

func RegisterInstance(instanceName string, instancePool *sync.Pool) {
	newInstance := bootingInstance{
		instanceName: instanceName,
		instancePool: instancePool,
	}
	_bootingInstances = append(_bootingInstances, newInstance)
}
