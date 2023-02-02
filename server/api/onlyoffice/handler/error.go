/**
 *
 * (c) Copyright Ascensio System SIA 2023
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
package handler

import (
	"errors"
	"fmt"
)

var ErrHandlerAlreadyRegistered = errors.New(_OnlyofficeLoggerPrefix + "handler with this code has already been registered")
var ErrInvalidUserID = errors.New(_OnlyofficeLoggerPrefix + "invalid callback user")

type FilePersistenceError struct {
	FileID string
	Reason string
}

func (e *FilePersistenceError) Error() string {
	return fmt.Sprintf("[ONLYOFFICE Filestore]: file %s could not be saved: %s", e.FileID, e.Reason)
}

type FileNotFoundError struct {
	FileID string
	Reason string
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf(_OnlyofficeLoggerPrefix+"file %s not found: %s", e.FileID, e.Reason)
}

type InvalidFileDownloadUrlError struct {
	FileID string
}

func (e *InvalidFileDownloadUrlError) Error() string {
	return fmt.Sprintf(_OnlyofficeLoggerPrefix+"could not find a callback file %s url", e.FileID)
}

type CallbackHandlerDoesNotExistError struct {
	Code int
}

func (e *CallbackHandlerDoesNotExistError) Error() string {
	return fmt.Sprintf(_OnlyofficeLoggerPrefix+"callback handler for code (%d) does not exist", e.Code)
}
