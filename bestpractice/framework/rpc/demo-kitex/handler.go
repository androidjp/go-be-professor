// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"

	"kdemo/kitex_gen/api"
	expert "kdemo/kitex_gen/api/expert"
)

// HelloImpl implements the last service interface defined in the IDL.
type HelloImpl struct{}

// Echo implements the HelloImpl interface.
func (s *HelloImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	// TODO: Your code here...
	resp = &api.Response{Message: req.Message}
	return
}

// QuestionSearch implements the ExpertServiceImpl interface.
func (s *ExpertServiceImpl) QuestionSearch(ctx context.Context, req *expert.QuestionSearchReqBody) (resp *expert.QuestionSearchRes, err error) {
	// TODO: Your code here...
	return
}

// Hello implements the ExpertServiceImpl interface.
func (s *ExpertServiceImpl) Hello(ctx context.Context, req *expert.HelloReqBody) (resp *expert.HelloRes, err error) {
	// TODO: Your code here...
	return
}
