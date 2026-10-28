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
package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestConverterModuleProvidesTimeConverter(t *testing.T) {
	var converter TimeConverter

	app := fxtest.New(t,
		fx.NopLogger,
		Module,
		fx.Populate(&converter),
	)

	app.RequireStart().RequireStop()

	require.NotNil(t, converter)
	assert.Equal(t, time.Unix(100, 0), converter.GetTime(100))
}
