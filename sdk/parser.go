/*
 * Copyright (c) 2022 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sdk

var (
	vmParser Parser
)

func SetParser(p Parser) {
	vmParser = p
}

type HookBitmap [2]uint64

type KeyVal struct {
	Key string
	Val string
}

type HttpAction interface {
	abort() bool
	getHttpResult() (*Trace, []KeyVal, error)
}

type ParseAction interface {
	abort() bool
	getParsePayloadResult() ([]*L7ProtocolInfo, error)
}

type Action interface {
	HttpAction
	ParseAction
}

func ActionAbort() Action {
	return &action{
		isAbort: true,
	}
}

func ActionAbortWithErr(err error) Action {
	return &action{
		isAbort: true,
		e:       err,
	}
}

func ActionNext() Action {
	return &action{
		isAbort: false,
	}
}

// agent will traversal to run all plugins, abort will abort the traversal, abort with no error will write the result to host.
type Parser interface {
	OnHttpReq(*HttpReqCtx) HttpAction
	OnHttpResp(*HttpRespCtx) HttpAction
	// protoNum return 0 indicate fail
	OnCheckPayload(*ParseCtx) (protoNum uint8, protoStr string)
	OnParsePayload(*ParseCtx) ParseAction
	HookEnable() HookBitmap
}

type action struct {
	e          error
	isAbort    bool
	httpResult struct {
		trace *Trace
		kv    []KeyVal
	}
	payloadResult []*L7ProtocolInfo
}

func (a *action) getHttpResult() (*Trace, []KeyVal, error) {
	return a.httpResult.trace, a.httpResult.kv, a.e
}

func (a *action) getParsePayloadResult() ([]*L7ProtocolInfo, error) {
	return a.payloadResult, a.e
}

func (e *action) abort() bool {
	return e.isAbort
}

func HttpActionAbortWithResult(trace *Trace, kv []KeyVal) HttpAction {
	return &action{
		isAbort: true,
		httpResult: struct {
			trace *Trace
			kv    []KeyVal
		}{trace: trace, kv: kv},
	}
}

func ParseActionAbortWithL7Info(info []*L7ProtocolInfo) ParseAction {
	return &action{
		isAbort:       true,
		payloadResult: info,
	}
}
