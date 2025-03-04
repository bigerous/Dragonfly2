/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package job

import (
	internaljob "d7y.io/dragonfly/v2/internal/job"
	"d7y.io/dragonfly/v2/manager/config"
)

type Job struct {
	*internaljob.Job
	Preheat
}

func New(cfg *config.Config) (*Job, error) {
	j, err := internaljob.New(&internaljob.Config{
		Host:      cfg.Database.Redis.Host,
		Port:      cfg.Database.Redis.Port,
		Password:  cfg.Database.Redis.Password,
		BrokerDB:  cfg.Database.Redis.BrokerDB,
		BackendDB: cfg.Database.Redis.BackendDB,
	}, internaljob.GlobalQueue)
	if err != nil {
		return nil, err
	}

	p, err := newPreheat(j, cfg.Server.Name)
	if err != nil {
		return nil, err
	}

	return &Job{
		Job:     j,
		Preheat: p,
	}, nil
}

func (j *Job) GetGroupJobState(id string) (*internaljob.GroupJobState, error) {
	groupJobState, err := j.Job.GetGroupJobState(id)
	if err != nil {
		return nil, err
	}

	return &internaljob.GroupJobState{
		GroupUUID: groupJobState.GroupUUID,
		State:     groupJobState.State,
		CreatedAt: groupJobState.CreatedAt,
	}, nil
}
