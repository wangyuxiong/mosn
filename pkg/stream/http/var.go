/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http

import (
	"context"
	"mosn.io/mosn/pkg/types"
	"strconv"

	"mosn.io/mosn/pkg/variable"
)

const (
	headerIndex = len(types.VarPrefixHttpHeader)
	argIndex    = len(types.VarPrefixHttpArg)
	cookieIndex = len(types.VarPrefixHttpCookie)
)

var (
	builtinVariables = []variable.Variable{
		variable.NewBasicVariable(types.VarHttpRequestMethod, nil, requestMethodGetter, nil, 0),
		variable.NewBasicVariable(types.VarHttpRequestLength, nil, requestLengthGetter, nil, 0),
	}

	prefixVariables = []variable.Variable{
		variable.NewBasicVariable(types.VarPrefixHttpHeader, nil, httpHeaderGetter, nil, 0),
		variable.NewBasicVariable(types.VarPrefixHttpArg, nil, httpArgGetter, nil, 0),
		variable.NewBasicVariable(types.VarPrefixHttpCookie, nil, httpCookieGetter, nil, 0),
	}
)

func init() {
	// register built-in variables
	for idx := range builtinVariables {
		variable.RegisterVariable(builtinVariables[idx])
	}

	// register prefix variables, like header_xxx/arg_xxx/cookie_xxx
	for idx := range prefixVariables {
		variable.RegisterPrefixVariable(prefixVariables[idx].Name(), prefixVariables[idx])
	}
}

func requestMethodGetter(ctx context.Context, value *variable.IndexedValue, data interface{}) (string, error) {
	buffers := httpBuffersByContext(ctx)
	request := &buffers.serverRequest

	return string(request.Header.Method()), nil
}

func requestLengthGetter(ctx context.Context, value *variable.IndexedValue, data interface{}) (string, error) {
	buffers := httpBuffersByContext(ctx)
	request := &buffers.serverRequest

	length := len(request.Header.Header()) + len(request.Body())
	if length == 0 {
		return variable.ValueNotFound, nil
	}

	return strconv.Itoa(length), nil
}

func httpHeaderGetter(ctx context.Context, value *variable.IndexedValue, data interface{}) (string, error) {
	buffers := httpBuffersByContext(ctx)
	request := &buffers.serverRequest

	headerName := data.(string)
	headerValue := request.Header.Peek(headerName[headerIndex:])
	// nil means no kv exists, "" means kv exists, but value is ""
	if headerValue == nil {
		return variable.ValueNotFound, nil
	}

	return string(headerValue), nil
}

func httpArgGetter(ctx context.Context, value *variable.IndexedValue, data interface{}) (string, error) {
	buffers := httpBuffersByContext(ctx)
	request := &buffers.serverRequest

	argName := data.(string)
	// TODO: support post args
	argValue := request.URI().QueryArgs().Peek(argName[argIndex:])
	// nil means no kv exists, "" means kv exists, but value is ""
	if argValue == nil {
		return variable.ValueNotFound, nil
	}

	return string(argValue), nil
}

func httpCookieGetter(ctx context.Context, value *variable.IndexedValue, data interface{}) (string, error) {
	buffers := httpBuffersByContext(ctx)
	request := &buffers.serverRequest

	cookieName := data.(string)
	cookieValue := request.Header.Cookie(cookieName[cookieIndex:])
	// nil means no kv exists, "" means kv exists, but value is ""
	if cookieValue == nil {
		return variable.ValueNotFound, nil
	}

	return string(cookieValue), nil
}
