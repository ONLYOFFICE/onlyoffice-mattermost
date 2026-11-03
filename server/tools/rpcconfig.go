/**
 *
 * (c) Copyright Ascensio System SIA 2026
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
 *
 */
package tools

import "reflect"

// RestoreFalseBools sets every nil *bool field of the struct pointed to by v
// to false.
//
// Plugins receive the server configuration over net/rpc with gob encoding,
// and gob does not transmit zero values: a *bool pointing to false arrives as
// nil. The server-side configuration always has its defaults applied, so a
// nil *bool on the plugin side can only mean "false". Without this, calling
// SetDefaults() would silently turn e.g. FileSettings.AmazonS3SSL=false into
// true and the plugin would talk HTTPS to a plain-HTTP S3 endpoint.
func RestoreFalseBools(v any) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return
	}

	boolPtr := reflect.TypeOf((*bool)(nil))
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if field.CanSet() && field.Type() == boolPtr && field.IsNil() {
			f := false
			field.Set(reflect.ValueOf(&f))
		}
	}
}
